package webbase

import (
	"log"
	"net/http"
	"strconv"

	"github.com/bonsai-oss/mux"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/prometheus/client_golang/prometheus"
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
		if rw.statusCode == 0 {
			rw.WriteHeader(http.StatusOK)
		}

		if transaction != nil {
			statusCode := rw.statusCode
			transaction.SetTag("Status-Code", strconv.Itoa(statusCode))
			transaction.Status = sentry.HTTPtoSpanStatus(statusCode)
		}
	})
}

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
		if rw.statusCode == 0 {
			rw.WriteHeader(http.StatusOK)
		}
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
