package main

import (
	"github.com/rekh/_temp/goprocoursedevopsproject/cmd/server/handlers"
	"net/http"
)

const (
	port = ":8080"
)

func main() {
	http.HandleFunc("/update/", handlers.UpdateMetrics)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "", http.StatusNotFound)
	})
	http.ListenAndServe(port, nil)
}
