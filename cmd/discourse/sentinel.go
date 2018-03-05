package main

import (
	"flag"
	"os"

	"github.com/SimonRichardson/flagset"
	"github.com/SimonRichardson/gexec"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func runSentinel(args []string) error {
	// flags for the sentinel command
	var (
		flags = flagset.NewFlagSet("sentinel", flag.ExitOnError)

		debug                = flags.Bool("debug", false, "debug logging")
		clusterBindAddr      = flags.String("cluster", defaultClusterAddr, "listen address for cluster")
		clusterAdvertiseAddr = flags.String("cluster.advertise-addr", "", "optional, explicit address to advertise in cluster")

		clusterPeers stringSlice
	)

	flags.Var(&clusterPeers, "peer", "cluster peer host:port (repeatable)")
	flags.Usage = usageFor(flags, "sentinel [flags]")
	if err := flags.Parse(args); err != nil {
		return nil
	}

	// Setup the logger.
	var logger log.Logger
	{
		logLevel := level.AllowInfo()
		if *debug {
			logLevel = level.AllowAll()
		}
		logger = log.NewLogfmtLogger(os.Stdout)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = level.NewFilter(logger, logLevel)
	}

	// Execution group.
	g := gexec.NewGroup()
	gexec.Block(g)
	gexec.Interrupt(g)
	return g.Run()
}
