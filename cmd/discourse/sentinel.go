package main

import (
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/SimonRichardson/discourse/pkg/cluster"
	"github.com/SimonRichardson/discourse/pkg/cluster/members"
	"github.com/SimonRichardson/discourse/pkg/sentinel"
	"github.com/SimonRichardson/discourse/pkg/status"
	"github.com/SimonRichardson/flagset"
	"github.com/SimonRichardson/gexec"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/pkg/errors"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	defaultClusterReplicationFactor = 5
	defaultMetricsRegistration      = true
)

const (
	SentinelPeerType members.PeerType = "peertype:sentinel"
)

func runSentinel(args []string) error {
	// flags for the sentinel command
	var (
		flags = flagset.NewFlagSet("sentinel", flag.ExitOnError)

		debug                    = flags.Bool("debug", false, "debug logging")
		debugCluster             = flags.Bool("debug.cluster", false, "debug cluster logging")
		apiAddr                  = flags.String("api", defaultAPIAddr, "listen address for query API")
		clusterBindAddr          = flags.String("cluster", defaultClusterAddr, "listen address for cluster")
		clusterAdvertiseAddr     = flags.String("cluster.advertise-addr", "", "optional, explicit address to advertise in cluster")
		clusterReplicationFactor = flags.Int("cluster.replication.factor", defaultClusterReplicationFactor, "replication factor for node configuration")
		metricsRegistration      = flags.Bool("metrics.registration", defaultMetricsRegistration, "Registration of metrics on launch")

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

	// Instrumentation
	connectedClients := prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: "coherence",
		Name:      "connected_clients",
		Help:      "Number of currently connected clients by modality.",
	}, []string{"modality"})
	apiDuration := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "coherence",
		Name:      "api_request_duration_seconds",
		Help:      "API request duration in seconds.",
		Buckets:   prometheus.DefBuckets,
	}, []string{"method", "path", "status_code"})

	if *metricsRegistration {
		prometheus.MustRegister(
			connectedClients,
			apiDuration,
		)
	}

	// Parse API addresses.
	var apiNetwork string
	var apiHost string
	var apiPort int
	{
		var err error
		apiNetwork, _, apiHost, apiPort, err = cluster.ParseAddr(*apiAddr, defaultAPIPort)
		if err != nil {
			return err
		}
	}

	apiListener, err := net.Listen(apiNetwork, net.JoinHostPort(apiHost, strconv.Itoa(apiPort)))
	if err != nil {
		return err
	}
	level.Debug(logger).Log("API", fmt.Sprintf("%s://%s", apiNetwork, net.JoinHostPort(apiHost, strconv.Itoa(apiPort))))

	// Parse cluster comms addresses.
	var chp cluster.HostPorts
	{
		var err error
		chp, err = cluster.CalculateHostPorts(
			*clusterBindAddr, *clusterAdvertiseAddr,
			defaultClusterPort, clusterPeers, logger,
		)
		if err != nil {
			return errors.Wrap(err, "calculating cluster hosts and ports")
		}
	}

	peer, err := configureRemoteCache(*debugCluster,
		logger,
		*clusterReplicationFactor,
		*apiAddr, defaultAPIPort,
		chp.BindHost, chp.BindPort,
		chp.AdvertiseHost, chp.AdvertisePort,
		clusterPeers.Slice(),
	)
	if err != nil {
		return err
	}

	// Execution group.
	g := gexec.NewGroup()
	gexec.Block(g)
	{
		cancel := make(chan struct{})
		g.Add(func() error {
			if _, err := peer.Join(); err != nil {
				return err
			}
			<-cancel
			return peer.Leave()
		}, func(error) {
			close(cancel)
		})
	}
	{
		g.Add(func() error {
			mux := http.NewServeMux()
			mux.Handle("/sentinel/", http.StripPrefix("/sentinel", sentinel.NewAPI(
				peer,
				log.With(logger, "component", "store_api"),
				connectedClients.WithLabelValues("api"),
				apiDuration,
			)))
			mux.Handle("/status/", status.NewAPI(
				log.With(logger, "component", "status_api"),
				connectedClients.WithLabelValues("status"),
				apiDuration,
			))

			registerMetrics(mux)
			registerProfile(mux)

			return http.Serve(apiListener, mux)
		}, func(error) {
			apiListener.Close()
		})
	}
	gexec.Interrupt(g)
	return g.Run()
}
