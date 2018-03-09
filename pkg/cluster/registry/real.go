package registry

import (
	"errors"
	"fmt"
	"sync"

	"github.com/SimonRichardson/alchemy/pkg/cluster/hashring"
	"github.com/SimonRichardson/alchemy/pkg/cluster/members"
)

var (
	// ErrNoHashRingFound states if a hashring is not found
	ErrNoHashRingFound = errors.New("no hashring found")
)

type real struct {
	mtx               sync.RWMutex
	hashRings         map[string]*hashring.HashRing
	keys              map[string]map[string]Key
	hashFn            func([]byte) uint32
	replicationFactor int
}

func New(hashFn func([]byte) uint32, replicationFactor int) Registry {
	return &real{
		hashRings:         make(map[string]*hashring.HashRing),
		keys:              make(map[string]Key),
		hashFn:            hashFn,
		replicationFactor: replicationFactor,
	}
}

func (r *real) Add(key Key) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	keyType := key.Type()
	if _, ok := r.hashRings[keyType]; !ok {
		r.hashRings[keyType] = hashring.New(r.hashFn, r.replicationFactor)
	}

	if r.hashRings[keyType].Add(key.Address()) {
		r.keys[key.Name()] = key
		return true
	}

	return false
}

func (r *real) Remove(key Key) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var (
		keyType = key.Type()
		addr    = key.Address()
	)
	if _, ok := r.hashRings[keyType]; ok {
		r.hashRings[keyType].Remove(addr)
	}
	if keys, ok := r.keys[addr]; ok {
		delete(keys, key.Name())
	}
	return true
}

func (r *real) Update(key Key) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var (
		keyType = key.Type()
		addr    = key.Address()
	)
	if _, ok := r.hashRings[keyType]; !ok || (ok && !r.hashRings[keyType].Contains(addr)) {
		return false
	}

	if _, ok := r.keys[addr]; !ok {
		return false
	}

	name := key.Name()
	if _, ok := r.keys[addr][name]; !ok {
		return false
	}
	r.keys[addr][name] = key

	return true
}

func (r *real) Locate(keyType string, val string) (Key, bool) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	if ring, ok := r.hashRings[keyType]; ok {
		if addr, ok := ring.Lookup(val); ok {
			r.keys[addr]
		}
	}
}

func (r *real) Info(s string) (Info, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	hashRing, ok := r.hashRings[s]
	if !ok {
		return Info{}, ErrNoHashRingFound
	}

	checksum, err := hashRing.Checksum()
	if err != nil {
		return Info{}, err
	}

	hashes := make(map[string]string)
	if err := hashRing.Walk(func(hash, addr string) error {
		hashes[hash] = addr
		return nil
	}); err != nil {
		return Info{}, err
	}

	keys := make(map[string][]Key)
	for _, v := range hashes {
		if k := r.getKeysByAddress(v); len(k) > 0 {
			keys[v] = append(keys[v], k...)
		}
	}

	return Info{
		Checksum: fmt.Sprintf("%08x", checksum),
		Hashes:   hashes,
		Keys:     keys,
	}, nil
}

func (r *real) Checksum(s string) (string, error) {
	r.mtx.RLock()
	defer r.mtx.RUnlock()

	hashRing, ok := r.hashRings[s]
	if !ok {
		return "", ErrNoHashRingFound
	}

	checksum, err := hashRing.Checksum()
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%08x", checksum), nil
}

func (r *real) getKeysByAddress(addr string) (res []Key) {
	if keys, ok := r.keys[addr]; ok {
		for _, v := range keys {
			res = append(res, v)
		}
	}
	return
}

type key struct {
	member members.Member
}

func NewMemberKey(member members.Member) Key {
	return &key{member}
}

func (k *key) Name() string {
	return k.member.Name()
}

func (k *key) Type() string {
	return k.member.PeerType().String()
}

func (k *key) Address() string {
	return k.member.Address()
}

func (k *key) Tags() map[string]string {
	return k.member.Tags()
}
