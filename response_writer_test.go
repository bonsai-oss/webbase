package webbase

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseWriterStatusCode(t *testing.T) {
	for _, testCase := range []struct {
		name         string
		handler      http.HandlerFunc
		expectedCode int
	}{
		{
			name:         "StatusCodeOKWhenHandlerDoesNotSet",
			handler:      http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}),
			expectedCode: http.StatusOK,
		},
		{
			name: "StatusCodeReflectsHandlerSetCode",
			handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
			}),
			expectedCode: http.StatusNotFound,
		},
	} {
		t.Run(testCase.name, func(t *testing.T) {
			testRecorder := httptest.NewRecorder()
			responseWriter := newResponseWriter(testRecorder)

			testCase.handler.ServeHTTP(responseWriter, nil)

			assert.Equal(t, testCase.expectedCode, responseWriter.statusCode)
		})
	}
}
