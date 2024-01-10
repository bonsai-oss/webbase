package main

import (
	"fmt"
	"net/http"

	"github.com/bonsai-oss/mux"
	"github.com/getsentry/sentry-go"

	"github.com/bonsai-oss/webbase"
)

func functionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	foo := vars["foo"]
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, foo)
}

func main() {
	router := webbase.NewRouter()
	router.Path("/{foo}/").Methods(http.MethodGet).HandlerFunc(functionHandler)
	webbase.ServeRouter("example", router,
		webbase.WithWebListenAddress("127.0.0.1:8080"),
		webbase.WithServiceListenAddress("127.0.0.1:8081"),
		webbase.WithSentryClientOptions(sentry.ClientOptions{
			TracesSampleRate: 1.0,
			SampleRate:       1.0,
			Debug:            true,
			Environment:      "development",
		}),
	)
}
