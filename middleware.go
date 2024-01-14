package webbase

import (
	"net/http"
	"strconv"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
)

func sentryMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		transaction := sentry.TransactionFromContext(r.Context())
		requestID := uuid.New()

		if transaction != nil {
			transaction.SetTag("Request-ID", requestID.String())
			transaction.SetTag("Function-Name", FunctionName)
		}

		w.Header().Set("X-Request-ID", requestID.String())
		rw := &responseWriter{ResponseWriter: w}

		next.ServeHTTP(rw, r)

		if transaction != nil {
			statusCode := rw.statusCode
			transaction.SetTag("Status-Code", strconv.Itoa(statusCode))
			transaction.Status = sentry.HTTPtoSpanStatus(statusCode)
		}
	})
}
