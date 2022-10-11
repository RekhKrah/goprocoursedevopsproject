package handlers

import (
	"fmt"
	"github.com/rekh/_temp/goprocoursedevopsproject/cmd/server/storage"
	"net/http"
	"regexp"
	"strconv"
)

const (
	contentType = "text/plain"
	//urlParserRegexp = `update/(\S+)/(\S+)/(\S+)$`
	urlParserRegexp = `[^\/\n]+`
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
	re := regexp.MustCompile(urlParserRegexp)
	//result := strings.Split(strings.TrimPrefix(url, "/"), "/")
	result := re.FindAllStringSubmatch(url, -1)
	data := storage.MetricsUrl{}

	if len(result) == 2 {
		http.Error(w, "Не указано имя метрики", http.StatusNotFound)
		return
	}
	if result[1][0] != "gauge" && result[1][0] != "counter" {
		http.Error(w, "Incorrect metric type", http.StatusNotImplemented)
		return
	}

	data.Type = result[1][0]

	if len(result) < 4 {
		http.Error(w, "metric value not found", http.StatusBadRequest)
		return
	}

	data.Name = result[2][0]

	_, err := strconv.ParseFloat(result[3][0], 64)

	if err != nil {
		http.Error(w, "Incorrect metric value", http.StatusBadRequest)
		return
	}
	data.Value = result[3][0]

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

//type MetricsUrl struct {
//	Name  string
//	Type  string
//	Value string
//}

//func parseMetricsURL(url string) (MetricsUrl, error) {
//	re := regexp.MustCompile(urlParserRegexp)
//	result := re.FindStringSubmatch(url)
//
//	return MetricsUrl{
//		Name:  result[2],
//		Type:  result[1],
//		Value: result[3],
//	}, nil
//}
