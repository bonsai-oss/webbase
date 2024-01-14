package webbase

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestSentryMiddleware(t *testing.T) {
	testCases := []struct {
		name           string
		handler        http.HandlerFunc
		expectedStatus sentry.SpanStatus
	}{
		{
			name:           "SetsCorrectTags",
			handler:        http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectedStatus: sentry.SpanStatusOK,
		},
		{
			name: "SetsCorrectStatusCode",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			expectedStatus: sentry.SpanStatusNotFound,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			rr := httptest.NewServer(sentryMiddleware(testCase.handler))
			defer rr.Close()

			request, _ := http.NewRequest(http.MethodGet, rr.URL, nil)
			response, _ := http.DefaultClient.Do(request)

			requestId := response.Header.Get("X-Request-ID")
			_, parseError := uuid.Parse(requestId)
			assert.NoError(t, parseError)
		})
	}
}
