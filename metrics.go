package webbase

import (
	"github.com/prometheus/client_golang/prometheus"
)

const (
	labelPath         = "path"
	labelFunctionName = "functionName"
	labelStatus       = "status"
)

var (
	totalRequests = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	}, []string{labelPath, labelFunctionName})

	responseStatus = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	}, []string{labelPath, labelStatus, labelFunctionName})

	httpDuration = prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_response_time_seconds",
		Help:    "Duration of HTTP requests.",
		Buckets: prometheus.ExponentialBuckets(0.005, 2, 15),
	}, []string{labelPath, labelFunctionName})
)

func init() {
	prometheus.MustRegister(totalRequests)
	prometheus.MustRegister(responseStatus)
	prometheus.MustRegister(httpDuration)
}
