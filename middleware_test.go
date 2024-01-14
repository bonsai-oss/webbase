package webbase

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/bonsai-oss/mux"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
)

func TestSentryMiddleware(t *testing.T) {
	testCases := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus sentry.SpanStatus
	}{
		{
			name:           "SetsCorrectTags",
			handler:        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectedStatus: sentry.SpanStatusOK,
		},
		{
			name: "SetsCorrectStatusCode",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			expectedStatus: sentry.SpanStatusNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rr := httptest.NewServer(sentryMiddleware(testCase.handler))
			defer rr.Close()

			request, _ := http.NewRequest(http.MethodGet, rr.URL, nil)
			response, _ := http.DefaultClient.Do(request)

			requestId := response.Header.Get("X-Request-ID")
			_, parseError := uuid.Parse(requestId)
			assert.NoError(t, parseError)
		})
	}
}

func TestPrometheusMiddleware(t *testing.T) {
	tests := []struct {
		name           string
		method         string
		url            string
		expectedStatus int
		expectedPath   string
	}{
		{
			name:           "GET request",
			method:         http.MethodGet,
			url:            "/testpath",
			expectedStatus: http.StatusOK,
			expectedPath:   "/testpath",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Create a request
			req, err := http.NewRequest(tc.method, tc.url, nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder to record the response
			rr := httptest.NewRecorder()

			// Create a router and register the middleware
			r := mux.NewRouter()
			r.Handle(tc.expectedPath, prometheusMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tc.expectedStatus)
			})))

			// Serve the HTTP request
			r.ServeHTTP(rr, req)

			// Check if the status code is as expected
			assert.Equal(t, tc.expectedStatus, rr.Code)

			foundMetricsCount := map[string]int{
				"http_requests_total":        0,
				"http_response_time_seconds": 0,
				"response_status":            0,
				"go_info":                    0,
				"go_gc_duration_seconds":     0,
			}

			gatheredMetrics, _ := prometheus.DefaultGatherer.Gather()
			for _, metricFamily := range gatheredMetrics {
				foundMetricsCount[metricFamily.GetName()]++
			}

			for metricName, count := range foundMetricsCount {
				assert.Equal(t, 1, count, "metric %s not found", metricName)
			}
		})
	}
}
