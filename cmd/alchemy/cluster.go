package main

import (
	"github.com/SimonRichardson/alchemy/pkg/cluster"
	"github.com/SimonRichardson/alchemy/pkg/cluster/members"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pborman/uuid"
	"github.com/pkg/errors"
)

func configureRemoteCache(debugCluster bool,
	logger log.Logger,
	replicationFactor int,
	apiAddr string, apiPort int,
	bindAddrHost string, bindAddrPort int,
	advertiseAddrHost string, advertiseAddrPort int,
	peers []string,
) (cluster.Peer, error) {
	clusterMembersConfig, err := members.Build(
		members.WithPeerType(RegistryPeerType),
		members.WithNodeName(uuid.New()),
		members.WithAPIAddrPort(apiAddr, apiPort),
		members.WithBindAddrPort(bindAddrHost, bindAddrPort),
		members.WithAdvertiseAddrPort(advertiseAddrHost, advertiseAddrPort),
		members.WithExisting(peers),
		members.WithLogOutput(membersLogOutput{
			output: debugCluster,
			logger: log.With(logger, "component", "cluster"),
		}),
	)
	if err != nil {
		return nil, errors.Wrap(err, "members remote config")
	}

	clusterMembers, err := members.NewRealMembers(clusterMembersConfig, log.With(logger, "component", "members"))
	if err != nil {
		return nil, errors.Wrap(err, "members remote")
	}

	return cluster.NewPeer(clusterMembers, log.With(logger, "component", "peer")), nil
}

type membersLogOutput struct {
	output bool
	logger log.Logger
}

func (m membersLogOutput) Write(b []byte) (int, error) {
	if m.output {
		level.Debug(m.logger).Log("fwd_msg", string(b))
		return len(b), nil
	}
	return 0, nil
}
