package webbase

import (
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path", "functionName"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status", "functionName"},
)

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_response_time_seconds",
	Help:    "Duration of HTTP requests.",
	Buckets: prometheus.ExponentialBuckets(0.005, 2, 15),
}, []string{"path", "functionName"})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.WithLabelValues(path, FunctionName))
		rw := NewResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.WithLabelValues(strconv.Itoa(statusCode), FunctionName).Inc()
		totalRequests.WithLabelValues(path, FunctionName).Inc()

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, statusCode, timer.ObserveDuration())
	})
}

func NewResponseWriter(w http.ResponseWriter) *responseWriter {
	return &responseWriter{w, http.StatusOK}
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	prometheus.Register(httpDuration)
}
