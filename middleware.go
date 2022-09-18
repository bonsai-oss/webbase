package webbase

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
)

func sentryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		transaction := sentry.TransactionFromContext(r.Context())
		requestID := uuid.New()

		transaction.SetTag("Request-ID", requestID.String())
		transaction.SetTag("Function-Name", FunctionName)

		w.Header().Set("X-Request-ID", requestID.String())
		next.ServeHTTP(w, r)
	})
}
