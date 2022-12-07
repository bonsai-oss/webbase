package main

import (
	"fmt"
	"net/http"

	"github.com/bonsai-oss/webbase"
)

func functionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
	fmt.Fprintln(w, r.RemoteAddr)
}

func main() {
	webbase.Serve("example", functionHandler)
}
