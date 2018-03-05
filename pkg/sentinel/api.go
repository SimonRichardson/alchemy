package sentinel

import (
	"net/http"
	"strconv"
	"time"

	"github.com/SimonRichardson/discourse/pkg/api"
	"github.com/SimonRichardson/discourse/pkg/cluster"
	"github.com/SimonRichardson/discourse/pkg/cluster/members"
	"github.com/SimonRichardson/discourse/pkg/metrics"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

// These are the sentinel API URL paths.
const (
	APIPathServicesQuery = "/services"
)

const (
	defaultContentType = "application/json"
)

// API wraps a sentinel and provides a basic HTTP API.
type API struct {
	handler  http.Handler
	peer     cluster.Peer
	logger   log.Logger
	clients  metrics.Gauge
	duration metrics.HistogramVec
	errors   api.Error
}

// NewAPI creates a API with the correct dependencies.
// The API is an http.Handler and can ServeHTTP.
//
//     GET /services
//         Returns the current list of all services according to the sentinel.
//
//     GET /services?type={type}
//         Returns the current list of services according to the sentinel that
//         correspond to the type.
//         Returns 400 Bad Request if the type is in an invalid format.
//         Returns 404 Not Found if the type doesn't exist.
//
func NewAPI(peer cluster.Peer,
	logger log.Logger,
	clients metrics.Gauge,
	duration metrics.HistogramVec,
) *API {
	api := &API{
		peer:     peer,
		logger:   logger,
		clients:  clients,
		duration: duration,
		errors:   api.NewError(logger),
	}
	{
		router := mux.NewRouter().StrictSlash(true)
		router.Methods("GET").Path(APIPathServicesQuery).HandlerFunc(api.handleServices)
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

func (a *API) handleServices(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	var params ServicesParams
	if err := params.DecodeFrom(r.Header, mux.Vars(r)); err != nil {
		a.errors.BadRequest(w, r, err.Error())
		return
	}

	services, err := a.peer.Current(params.Type)
	if err != nil {
		a.errors.InternalServerError(w, r, err.Error())
		return
	}

	if list, ok := services[params.Type]; !ok || len(list) == 0 {
		a.errors.NotFound(w, r)
		return
	}
}

// ServicesParams handles
type ServicesParams struct {
	Type members.PeerType
}

// DecodeFrom populates a ServicesParams from a Request.
func (p *ServicesParams) DecodeFrom(headers http.Header, vars map[string]string) (err error) {
	if accept := headers.Get("Accept"); accept != "" && accept != defaultContentType {
		return errors.Errorf("expected %q content-type, got %q", defaultContentType, accept)
	}

	if typ, ok := vars["type"]; !ok || typ == "" {
		p.Type = cluster.PeerTypeAny
	} else {
		p.Type, err = cluster.ParsePeerType(typ)
	}
	return
}
