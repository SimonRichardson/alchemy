package registry

import "github.com/SimonRichardson/alchemy/pkg/cluster/hashring"
import "github.com/SimonRichardson/alchemy/pkg/cluster/members"
import "sync"

type real struct {
	mtx               sync.Mutex
	hashrings         map[members.PeerType]*hashring.HashRing
	members           map[string][]members.Member
	hashFn            func([]byte) uint32
	replicationFactor int
}

func (r *real) Add(m []members.Member) {
	r.mtx.Lock()
	defer r.mtx.Unlock()

	for _, v := range m {
		peerType := v.PeerType()
		if _, ok := r.hashrings[peerType]; !ok {
			r.hashrings[peerType] = hashring.New(r.hashFn, r.replicationFactor)
		}

		addr := v.Address()
		if ok := r.hashrings[peerType].Add(addr); ok {
			r.members[addr] = append(r.members[addr], v)
		}
	}
}

func (r *real) Remove(m []members.Member) {

}

func (r *real) Update(m []members.Member) {

}
