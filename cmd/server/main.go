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
	http.ListenAndServe(port, nil)
}
