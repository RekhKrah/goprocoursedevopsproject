package main

import (
	"fmt"
	"math/rand"
	"net/http"
	"reflect"
	"regexp"
	"runtime"
	"time"
)

func main() {
	m := Metrics{}

	poll := time.After(pollInterval)
	report := time.After(reportInterval)

	for {
		select {
		case <-poll:
			getMetrics(&m)
			pollCount += 1
			m.PollCount = pollCount
			m.RandomValue = genRandomGauge()
			fmt.Println(m)
			poll = time.After(pollInterval)

		case <-report:
			err := sendMetrics(m)
			if err != nil {
				panic(err)
			}
			report = time.After(reportInterval)
		}
	}
}

type gauge float64
type counter int64

type Metrics struct {
	Alloc         gauge
	BuckHashSys   gauge
	Frees         gauge
	GCCPUFraction gauge
	GCSys         gauge
	HeapAlloc     gauge
	HeapIdle      gauge
	HeapInuse     gauge
	HeapObjects   gauge
	HeapReleased  gauge
	HeapSys       gauge
	LastGC        gauge
	Lookups       gauge
	MCacheInuse   gauge
	MCacheSys     gauge
	MSpanInuse    gauge
	MSpanSys      gauge
	Mallocs       gauge
	NextGC        gauge
	NumForcedGC   gauge
	NumGC         gauge
	OtherSys      gauge
	PauseTotalNs  gauge
	StackInuse    gauge
	StackSys      gauge
	Sys           gauge
	TotalAlloc    gauge
	PollCount     counter
	RandomValue   gauge
}

var pollCount counter = 0

const (
	pollInterval   = 2 * time.Second
	reportInterval = 10 * time.Second

	host = "127.0.0.1"
	port = "8080"
	//method      = http.MethodPost
	contentType = "text/plain"
)

func genURL(host, port, metricType, metricName, metricValue string) string {
	return fmt.Sprintf("http://%v:%v/update/%v/%v/%v",
		host, port, metricType, metricName, metricValue)
}

func genRandomGauge() gauge {
	return gauge(rand.Float64())
}

func getMetrics(metrics *Metrics) {
	m := runtime.MemStats{}
	runtime.ReadMemStats(&m)
	metrics.Alloc = gauge(m.Alloc)
	metrics.BuckHashSys = gauge(m.BuckHashSys)
	metrics.Frees = gauge(m.Frees)
	metrics.GCCPUFraction = gauge(m.GCCPUFraction)
	metrics.GCSys = gauge(m.GCSys)
	metrics.HeapAlloc = gauge(m.HeapAlloc)
	metrics.HeapIdle = gauge(m.HeapIdle)
	metrics.HeapInuse = gauge(m.HeapInuse)
	metrics.HeapObjects = gauge(m.HeapObjects)
	metrics.HeapReleased = gauge(m.HeapReleased)
	metrics.HeapSys = gauge(m.HeapSys)
	metrics.LastGC = gauge(m.LastGC)
	metrics.Lookups = gauge(m.Lookups)
	metrics.MCacheInuse = gauge(m.MCacheInuse)
	metrics.MCacheSys = gauge(m.MCacheSys)
	metrics.MSpanInuse = gauge(m.MSpanInuse)
	metrics.MSpanSys = gauge(m.MSpanSys)
	metrics.Mallocs = gauge(m.Mallocs)
	metrics.NextGC = gauge(m.NextGC)
	metrics.NumForcedGC = gauge(m.NumForcedGC)
	metrics.NumGC = gauge(m.NumGC)
	metrics.OtherSys = gauge(m.OtherSys)
	metrics.PauseTotalNs = gauge(m.PauseTotalNs)
	metrics.StackInuse = gauge(m.StackInuse)
	metrics.StackSys = gauge(m.StackSys)
	metrics.Sys = gauge(m.Sys)
	metrics.TotalAlloc = gauge(m.TotalAlloc)
}

func sendMetrics(metrics Metrics) error {
	v := reflect.ValueOf(metrics)

	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		value := v.FieldByName(key)

		re := regexp.MustCompile(`\w+$`)
		valueType := re.FindString(value.Type().String())

		fmt.Printf("******************%v: %v - %v \n", key, value, valueType)

		url := genURL(host, port, valueType, key, fmt.Sprintf("%v", value))
		fmt.Println(url)

		_, err := http.Post(url, contentType, nil)
		if err != nil {
			return err
		}
	}
	return nil
}
