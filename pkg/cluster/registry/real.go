package registry

import (
	"sync"

	"github.com/SimonRichardson/alchemy/pkg/cluster/hashring"
	"github.com/SimonRichardson/alchemy/pkg/cluster/members"
)

type real struct {
	mtx               sync.Mutex
	hashrings         map[string]*hashring.HashRing
	keys              map[Address]map[string]Key
	hashFn            func([]byte) uint32
	replicationFactor int
}

func New(hashFn func([]byte) uint32, replicationFactor int) Registry {
	return &real{
		hashrings:         make(map[string]*hashring.HashRing),
		keys:              make(map[Address]map[string]Key),
		hashFn:            hashFn,
		replicationFactor: replicationFactor,
	}
}

func (r *real) Add(key Key) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	keyType := key.Type()
	if _, ok := r.hashrings[keyType]; !ok {
		r.hashrings[keyType] = hashring.New(r.hashFn, r.replicationFactor)
	}

	var (
		addr = key.Address()
		res  = r.hashrings[keyType].Add(addr.HostPort())
	)
	if _, ok := r.keys[addr]; !ok {
		r.keys[addr] = make(map[string]Key)
	}
	r.keys[addr][key.Name()] = key

	return res
}

func (r *real) Remove(key Key) bool {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	var (
		keyType = key.Type()
		addr    = key.Address()
	)
	if _, ok := r.hashrings[keyType]; ok {
		r.hashrings[keyType].Remove(addr.HostPort())
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
	if _, ok := r.hashrings[keyType]; !ok || (ok && !r.hashrings[keyType].Contains(addr.HostPort())) {
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

func (k *key) Address() Address {
	return Address(k.member.Address())
}

func (k *key) Tags() map[string]string {
	return k.member.Tags()
}
