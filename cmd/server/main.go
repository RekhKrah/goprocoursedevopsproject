package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/rekh/_temp/goprocoursedevopsproject/cmd/server/handlers"
	"net/http"
)

const (
	port = ":8080"
)

func main() {
	r := chi.NewRouter()
	r.Post("/update/{type}/{name}/{value}", handlers.UpdateMetrics)

	r.Get("/", handlers.GetMetrics)
	r.Get("/value/{type}/{name}", handlers.GetMetric)

	http.ListenAndServe(port, r)
}
