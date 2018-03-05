package hashring

import (
	"fmt"
	"strings"
	"sync"
)

type hashFn func(string) int

// HashRing stores strings on a consistent hash ring. HashRing internally uses
// Red-Black Tree to achieve O(log N) lookup and insertion time.
type HashRing struct {
	mtx               sync.RWMutex
	hashFn            hashFn
	replicationFactor int
	hosts             map[string]struct{}
	tree              *RBTree
}

// New creates a new HashRing with a replication factor
func New(hashFn func([]byte) uint32, replicationFactor int) *HashRing {
	return &HashRing{
		hashFn: func(s string) int {
			return int(hashFn([]byte(s)))
		},
		replicationFactor: replicationFactor,
		hosts:             make(map[string]struct{}, 0),
		tree:              NewRBTree(),
	}
}

// Add a host and replicate it around the hashring according to the replication
// factor.
// Returns true if an insertion happens for all replicated points
func (r *HashRing) Add(host string) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.hosts[host]; ok {
		return false
	}

	r.hosts[host] = struct{}{}

	added := true
	for i := 0; i < r.replicationFactor; i++ {
		key := fmt.Sprintf("%s%d", host, i)
		added = added && r.tree.Insert(r.hashFn(key), host)
	}
	return added
}

// Remove a host from the hashring including all the subsequent replicated
// hosts.
// Returns true if a deletion happens to all the replicated points
func (r *HashRing) Remove(host string) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	if _, ok := r.hosts[host]; !ok {
		return false
	}

	removed := true
	for i := 0; i < r.replicationFactor; i++ {
		key := fmt.Sprintf("%s%d", host, i)
		removed = removed && r.tree.Delete(r.hashFn(key))
	}

	delete(r.hosts, host)

	return removed
}

// Lookup returns the owner of the given key and whether the HashRing contains
// the key at all.
func (r *HashRing) Lookup(key string) (string, bool) {
	if s := r.LookupN(key, 1); len(s) > 0 {
		return s[0], true
	}
	return "", false
}

// LookupN returns the N servers that own the given key. Duplicates in the form
// of virtual nodes are skipped to maintain a list of unique servers. If there
// are less servers than N, we simply return all existing servers.
func (r *HashRing) LookupN(key string, n int) []string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.tree.LookupNUniqueAt(n, r.hashFn(key))
}

// Contains checks to see if a key is already in the ring.
// Returns true if a key is found with in the ring.
func (r *HashRing) Contains(key string) bool {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if _, ok := r.hosts[key]; ok {
		return true
	}

	_, ok := r.tree.Search(r.hashFn(key))
	return ok
}

// Hosts returns the hosts in a slice.
func (r *HashRing) Hosts() []string {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var (
		idx int
		res = make([]string, len(r.hosts))
	)
	for k := range r.hosts {
		res[idx] = k
		idx++
	}
	return res
}

// Walk iterates over each node in the hashring, because of the replication
// factor, the number of nodes you walk over will be many.
// If an error is returned whilst walking the nodes, it will stop walking
// immediately and return that error.
func (r *HashRing) Walk(fn func(string, string) error) error {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return r.tree.Walk(func(n *RBNode) error {
		return fn(fmt.Sprintf("%08x", n.key), n.value)
	})
}

// Len returns the number of unique hosts
func (r *HashRing) Len() int {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	return len(r.hosts)
}

// Checksum the hashring to verify if there have been any changes
func (r *HashRing) Checksum() (uint32, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	var values []string
	if err := r.tree.Walk(func(n *RBNode) error {
		values = append(values, fmt.Sprintf("%s:%s", n.nodeType.String(), n.value))
		return nil
	}); err != nil {
		return 0, err
	}

	return uint32(r.hashFn(strings.Join(values, ";"))), nil
}
