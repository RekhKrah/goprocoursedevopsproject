package handlers

import (
	"fmt"
	"github.com/go-chi/chi/v5"
	"github.com/rekh/_temp/goprocoursedevopsproject/cmd/server/storage"
	"net/http"
	"strconv"
)

const (
	contentType = "text/plain"
	//urlParserRegexp = `update/(\S+)/(\S+)/(\S+)$`
	urlParserRegexp = `[^\/\n]+`
)

var memStorage storage.MemStorage

func GetMetrics(w http.ResponseWriter, r *http.Request) {
	metrics := memStorage.GetMetrics()

	var list string
	for _, v := range metrics {
		list += fmt.Sprintf("<br>%v = %v</br>", v.Name, v.Value)
	}
	page := fmt.Sprintf("<html><body>%v</body></html>", list)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	fmt.Fprintln(w, page)
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

	//url := r.URL.String()
	//re := regexp.MustCompile(urlParserRegexp)
	//result := re.FindAllStringSubmatch(url, -1)
	data := storage.MetricsUrl{}

	mt := chi.URLParam(r, "type")
	mn := chi.URLParam(r, "name")
	mv := chi.URLParam(r, "value")

	fmt.Println(mt + mn + mv)

	if len(mn) == 0 {
		http.Error(w, "Не указано имя метрики", http.StatusNotFound)
		return
	}
	if mt != "gauge" && mt != "counter" {
		http.Error(w, "Incorrect metric type", http.StatusNotImplemented)
		return
	}

	data.Type = mt

	if len(mv) == 0 {
		http.Error(w, "metric value not found", http.StatusBadRequest)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Only POST request are allowed", http.StatusMethodNotAllowed)
		return
	}

	data.Name = mn

	_, err := strconv.ParseFloat(mv, 64)

	if err != nil {
		http.Error(w, "Incorrect metric value", http.StatusBadRequest)
		return
	}
	data.Value = mv
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
