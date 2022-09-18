package webbase

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const FunctionAddress = ":8080"
const MetricsAddress = ":8081"

var (
	FunctionName string
)

func Serve(name string, functionHandler func(w http.ResponseWriter, r *http.Request)) {
	// initialize sentry connection
	sentry.Init(sentry.ClientOptions{
		TracesSampleRate: 1.0,
		Debug:            true,
		Transport:        sentry.NewHTTPSyncTransport(),
	})
	defer sentry.Flush(2 * time.Second)
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		WaitForDelivery: false,
		Timeout:         1 * time.Second,
	})

	FunctionName = name
	router := mux.NewRouter()
	router.Use(sentryHandler.Handle, prometheusMiddleware, sentryMiddleware)

	go func() {
		http.Handle("/metrics", promhttp.Handler())
		http.HandleFunc("/health", func(writer http.ResponseWriter, request *http.Request) {
			fmt.Fprint(writer, "ok")
		})
		http.ListenAndServe(MetricsAddress, nil)
	}()
	router.PathPrefix("/function").HandlerFunc(functionHandler)
	err := http.ListenAndServe(FunctionAddress, router)
	log.Fatal(err)
}
