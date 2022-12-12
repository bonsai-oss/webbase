package webbase

import (
	"time"

	"github.com/bonsai-oss/mux"
	sentryhttp "github.com/getsentry/sentry-go/http"
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
