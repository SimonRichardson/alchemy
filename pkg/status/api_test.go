package status

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	metricMocks "github.com/SimonRichardson/alchemy/pkg/metrics/mocks"
	"github.com/go-kit/kit/log"
	"github.com/golang/mock/gomock"
)

func TestAPI(t *testing.T) {
	t.Parallel()

	t.Run("liveness", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			api      = NewAPI(log.NewNopLogger(), clients, duration)
			server   = httptest.NewServer(api)
		)
		defer server.Close()

		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("GET", "/status/health", "200").Return(observer).Times(1)
		observer.EXPECT().Observe(MatchAnyFloat64()).Times(1)

		response, err := http.Get(fmt.Sprintf("%s/status/health", server.URL))
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := http.StatusOK, response.StatusCode; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})

	t.Run("readiness", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		var (
			clients  = metricMocks.NewMockGauge(ctrl)
			duration = metricMocks.NewMockHistogramVec(ctrl)
			observer = metricMocks.NewMockObserver(ctrl)
			api      = NewAPI(log.NewNopLogger(), clients, duration)
			server   = httptest.NewServer(api)
		)
		defer server.Close()

		clients.EXPECT().Inc().Times(1)
		clients.EXPECT().Dec().Times(1)

		duration.EXPECT().WithLabelValues("GET", "/status/ready", "200").Return(observer).Times(1)
		observer.EXPECT().Observe(MatchAnyFloat64()).Times(1)

		response, err := http.Get(fmt.Sprintf("%s/status/ready", server.URL))
		if err != nil {
			t.Fatal(err)
		}

		if expected, actual := http.StatusOK, response.StatusCode; expected != actual {
			t.Errorf("expected: %d, actual: %d", expected, actual)
		}
	})
}

type float64AnyMatcher struct{}

func (float64AnyMatcher) Matches(x interface{}) bool {
	_, ok := x.(float64)
	return ok
}

func (float64AnyMatcher) String() string {
	return "is float64"
}

// MatchAnyFloat64 checks to see if the value it's trying to match is a float64
func MatchAnyFloat64() gomock.Matcher { return float64AnyMatcher{} }
