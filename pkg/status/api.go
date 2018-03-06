package status

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/SimonRichardson/alchemy/pkg/api"
	"github.com/SimonRichardson/alchemy/pkg/metrics"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
)

// These are the status API URL paths.
const (
	APIPathLivenessQuery  = "/status/health"
	APIPathReadinessQuery = "/status/ready"
)

// API for Liveness and provides a basic HTTP API.
// Note that this is an artificially restricted API.
//
//     GET /status/health
//         Returns the current health of the server.
//         Returns 500 Internal Server Error if the health of the server has
//         deteriorated.
//
//     GET /status/ready
//         Performs a check of the server to see if everything is operational
//         to perform requests.
//         Returns 504 Gateway Timeout Error if the server isn't ready.
//         Returns 500 Internal Server Error if the health of the server has
//         deteriorated.
//
type API struct {
	handler  http.Handler
	logger   log.Logger
	clients  metrics.Gauge
	duration metrics.HistogramVec
	errors   api.Error
}

// NewAPI creates a API with the correct dependencies.
// The API is an http.Handler and can ServeHTTP.
func NewAPI(logger log.Logger,
	clients metrics.Gauge,
	duration metrics.HistogramVec,
) *API {
	api := &API{
		logger:   logger,
		clients:  clients,
		duration: duration,
		errors:   api.NewError(logger),
	}
	{
		router := mux.NewRouter().StrictSlash(true)
		router.Methods("GET").Path(APIPathLivenessQuery).HandlerFunc(api.handleLiveness)
		router.Methods("GET").Path(APIPathReadinessQuery).HandlerFunc(api.handleReadiness)
		router.NotFoundHandler = http.HandlerFunc(api.errors.NotFound)
		api.handler = router
	}
	return api
}

func (a *API) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	level.Info(a.logger).Log("method", r.Method, "url", r.URL.String())

	iw := &interceptingWriter{http.StatusOK, w}
	w = iw

	// Metrics
	a.clients.Inc()
	defer a.clients.Dec()

	defer func(begin time.Time) {
		a.duration.WithLabelValues(
			r.Method,
			r.URL.Path,
			strconv.Itoa(iw.code),
		).Observe(time.Since(begin).Seconds())
	}(time.Now())

	a.handler.ServeHTTP(w, r)
}

func (a *API) handleLiveness(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	a.respond(w, r)
}

func (a *API) handleReadiness(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	a.respond(w, r)
}

func (a *API) respond(w http.ResponseWriter, r *http.Request) {
	if accept := r.Header.Get("Accept"); accept != "" && accept != "application/json" {
		a.errors.BadRequest(w, r, "invalid accept header")
		return
	}

	if err := json.NewEncoder(w).Encode(struct{}{}); err != nil {
		a.errors.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
