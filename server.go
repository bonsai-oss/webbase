package webbase

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const FunctionAddress = ":8080"
const MetricsAddress = ":8081"

var (
	FunctionName string
)

// Serve is a helper function to start a server with a single function handler
//
// Deprecated: Use ServeFunction instead
func Serve(name string, functionHandler http.HandlerFunc) {
	ServeFunction(name, functionHandler)
}

// ServeFunction starts a server with a single handler function
func ServeFunction(name string, functionHandler http.HandlerFunc) {
	router := NewRouter()
	router.PathPrefix("/").HandlerFunc(functionHandler)
	ServeRouter(name, router)
}

// ServeRouter starts a server with a set of routes
func ServeRouter(name string, router *mux.Router) {
	// initialize sentry connection
	sentry.Init(sentry.ClientOptions{
		TracesSampleRate: 1.0,
		Debug:            true,
		Transport:        sentry.NewHTTPSyncTransport(),
	})
	defer sentry.Flush(2 * time.Second)
	FunctionName = name
	go serviceListener()
	err := http.ListenAndServe(FunctionAddress, router)
	log.Fatal(err)
}

// serviceListener starts a server to listen for metrics and health checks
func serviceListener() {
	prometheus.DefaultRegisterer.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "ok")
	})
	http.ListenAndServe(MetricsAddress, nil)
}
