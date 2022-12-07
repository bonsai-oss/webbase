package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"

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
	webbase.ServeRouter("example", router)
}
