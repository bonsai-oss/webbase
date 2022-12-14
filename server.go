package webbase

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/bonsai-oss/mux"
	"github.com/getsentry/sentry-go"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const (
	DefaultFunctionAddress = ":8080"
	DefaultMetricsAddress  = ":8081"
)

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
func ServeFunction(name string, functionHandler http.HandlerFunc, options ...serveOption) {
	router := NewRouter()
	router.PathPrefix("/").HandlerFunc(functionHandler)
	ServeRouter(name, router, options...)
}

// ServeRouter starts a server with a set of routes
func ServeRouter(name string, router *mux.Router, options ...serveOption) {
	config := serveConfiguration{
		webListenAddress:     DefaultFunctionAddress,
		serviceListenAddress: DefaultMetricsAddress,
	}
	for _, option := range options {
		if optionApplyError := option(&config); optionApplyError != nil {
			log.Fatalf("failed to apply serveOption: %v", optionApplyError)
		}
	}

	// initialize sentry connection
	sentry.Init(sentry.ClientOptions{
		TracesSampleRate: 1.0,
		EnableTracing:    true,
		Debug:            config.sentryDebug,
		Transport:        sentry.NewHTTPSyncTransport(),
	})
	defer sentry.Flush(2 * time.Second)
	FunctionName = name
	go serviceListener(config.serviceListenAddress)
	err := http.ListenAndServe(config.webListenAddress, router)
	log.Fatal(err)
}

// serviceListener starts a server to listen for metrics and health checks
func serviceListener(listenAddress string) {
	prometheus.DefaultRegisterer.Register(prometheus.NewProcessCollector(prometheus.ProcessCollectorOpts{}))
	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprint(writer, "ok")
	})
	http.ListenAndServe(listenAddress, nil)
}
