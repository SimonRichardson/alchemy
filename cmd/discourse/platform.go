package main

import (
	"flag"
	"os"

	"github.com/SimonRichardson/flagset"
	"github.com/SimonRichardson/gexec"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

func runPlatform(args []string) error {
	// flags for the platform command
	var (
		flags = flagset.NewFlagSet("platform", flag.ExitOnError)

		debug = flags.Bool("debug", false, "debug logging")
	)

	flags.Usage = usageFor(flags, "platform [flags]")
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
