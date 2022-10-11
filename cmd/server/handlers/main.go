package handlers

import (
	"fmt"
	"github.com/rekh/_temp/goprocoursedevopsproject/cmd/server/storage"
	"net/http"
)

const (
	contentType = "text/plain"
)

var memStorage = storage.MemStorage{}

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request are allowed", http.StatusMethodNotAllowed)
		return
	}
	if r.Header.Get("Content-Type") != contentType {
		http.Error(w, fmt.Sprintf("Only %v Content-Type header is allowed", contentType), http.StatusBadRequest)
		return
	}

	url := r.URL.String()

	err := memStorage.Update(url)

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
