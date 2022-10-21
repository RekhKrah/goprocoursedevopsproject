package handlers

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/rekh/_temp/goprocoursedevopsproject/internal"
)

const (
	contentType = "text/plain"
)

var memStorage internal.MemStorage

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := memStorage.GetMetrics()

	var list string
	for _, v := range metrics {
		list += fmt.Sprintf("%v = %v\n", v.Name, v.Value)
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, list)
}

func GetMetric(w http.ResponseWriter, r *http.Request) {
	mt := chi.URLParam(r, "type")
	if len(mt) == 0 {
		http.Error(w, "Request is incorrect", http.StatusBadRequest)
		return
	}

	mn := chi.URLParam(r, "name")
	if len(mn) == 0 {
		http.Error(w, "Metric is non found", http.StatusNotFound)
		return
	}

	result, err := memStorage.GetMetricValue(mt, mn)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, result)
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Content-Type") != contentType {
		http.Error(w, fmt.Sprintf("Only %v Content-Type header is allowed", contentType), http.StatusBadRequest)
		return
	}

	url := strings.Split(r.URL.String(), "/")[2:]

	l := len(url)
	if l == 1 {
		http.Error(w, "metric name is not found", http.StatusNotFound)
		return
	}
	if l == 2 {
		if len(url[1]) == 0 {
			http.Error(w, "metric name is not found", http.StatusNotFound)
			return
		}
		http.Error(w, "metric value is not found", http.StatusBadRequest)
		return
	}

	mt := url[0]
	mn := url[1]
	mv := url[2]

	if mt != "gauge" && mt != "counter" {
		http.Error(w, "Incorrect metric type", http.StatusNotImplemented)
		return
	}

	_, err := strconv.ParseFloat(mv, 64)

	if err != nil {
		http.Error(w, "Incorrect metric value", http.StatusBadRequest)
		return
	}

	data := internal.MetricsURL{
		Type:  mt,
		Name:  mn,
		Value: mv,
	}
	err = memStorage.Update(data)

	status := http.StatusOK
	message := "data accepted"
	if err != nil {
		status = http.StatusBadRequest
		message = err.Error()
	}

	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintln(w, message)
}
