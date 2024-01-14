package webbase

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bonsai-oss/mux"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelPath         = "path"
	labelFunctionName = "functionName"
	labelStatus       = "status"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{labelPath, labelFunctionName},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{labelPath, labelStatus, labelFunctionName},
)

var httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
	Name:    "http_response_time_seconds",
	Help:    "Duration of HTTP requests.",
	Buckets: prometheus.ExponentialBuckets(0.005, 2, 15),
}, []string{labelPath, labelFunctionName})

func prometheusMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		route := mux.CurrentRoute(r)
		path, _ := route.GetPathTemplate()

		timer := prometheus.NewTimer(httpDuration.With(prometheus.Labels{
			labelPath:         path,
			labelFunctionName: FunctionName,
		}))
		rw := newResponseWriter(w)
		next.ServeHTTP(rw, r)

		statusCode := rw.statusCode

		responseStatus.With(prometheus.Labels{
			labelPath:         path,
			labelStatus:       strconv.Itoa(statusCode),
			labelFunctionName: FunctionName,
		}).Inc()
		totalRequests.With(prometheus.Labels{
			labelPath:         path,
			labelFunctionName: FunctionName,
		}).Inc()

		log.Printf("%s %s %d %s", r.Method, r.URL.Path, statusCode, timer.ObserveDuration())
	})
}

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(responseStatus)
	prometheus.MustRegister(httpDuration)
}
