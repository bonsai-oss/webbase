package webbase

import (
	"time"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/gorilla/mux"
)

// NewRouter creates a new router with the default webbase middleware
func NewRouter() *mux.Router {
	sentryHandler := sentryhttp.New(sentryhttp.Options{
		WaitForDelivery: false,
		Timeout:         1 * time.Second,
	})

	router := mux.NewRouter()
	router.Use(sentryHandler.Handle, prometheusMiddleware, sentryMiddleware)

	return router
}
