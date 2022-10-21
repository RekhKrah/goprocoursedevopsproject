package agent

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"os"
	"reflect"
	"runtime"
	"strings"
	"time"
)

type Gauge float64
type Counter int64

type Metrics struct {
	Alloc         Gauge
	BuckHashSys   Gauge
	Frees         Gauge
	GCCPUFraction Gauge
	GCSys         Gauge
	HeapAlloc     Gauge
	HeapIdle      Gauge
	HeapInuse     Gauge
	HeapObjects   Gauge
	HeapReleased  Gauge
	HeapSys       Gauge
	LastGC        Gauge
	Lookups       Gauge
	MCacheInuse   Gauge
	MCacheSys     Gauge
	MSpanInuse    Gauge
	MSpanSys      Gauge
	Mallocs       Gauge
	NextGC        Gauge
	NumForcedGC   Gauge
	NumGC         Gauge
	OtherSys      Gauge
	PauseTotalNs  Gauge
	StackInuse    Gauge
	StackSys      Gauge
	Sys           Gauge
	TotalAlloc    Gauge
	PollCount     Counter
	RandomValue   Gauge
}

var cfg = config{}

func (m *Metrics) Get() {
	checkConfig()

	ms := runtime.MemStats{}
	runtime.ReadMemStats(&ms)
	m.Alloc = Gauge(ms.Alloc)
	m.BuckHashSys = Gauge(ms.BuckHashSys)
	m.Frees = Gauge(ms.Frees)
	m.GCCPUFraction = Gauge(ms.GCCPUFraction)
	m.GCSys = Gauge(ms.GCSys)
	m.HeapAlloc = Gauge(ms.HeapAlloc)
	m.HeapIdle = Gauge(ms.HeapIdle)
	m.HeapInuse = Gauge(ms.HeapInuse)
	m.HeapObjects = Gauge(ms.HeapObjects)
	m.HeapReleased = Gauge(ms.HeapReleased)
	m.HeapSys = Gauge(ms.HeapSys)
	m.LastGC = Gauge(ms.LastGC)
	m.Lookups = Gauge(ms.Lookups)
	m.MCacheInuse = Gauge(ms.MCacheInuse)
	m.MCacheSys = Gauge(ms.MCacheSys)
	m.MSpanInuse = Gauge(ms.MSpanInuse)
	m.MSpanSys = Gauge(ms.MSpanSys)
	m.Mallocs = Gauge(ms.Mallocs)
	m.NextGC = Gauge(ms.NextGC)
	m.NumForcedGC = Gauge(ms.NumForcedGC)
	m.NumGC = Gauge(ms.NumGC)
	m.OtherSys = Gauge(ms.OtherSys)
	m.PauseTotalNs = Gauge(ms.PauseTotalNs)
	m.StackInuse = Gauge(ms.StackInuse)
	m.StackSys = Gauge(ms.StackSys)
	m.Sys = Gauge(ms.Sys)
	m.TotalAlloc = Gauge(ms.TotalAlloc)
	m.PollCount += 1
	m.RandomValue = Gauge(rand.Float64())
}

func (m Metrics) Send() error {
	checkConfig()

	v := reflect.ValueOf(m)

	for i := 0; i < v.NumField(); i++ {
		key := v.Type().Field(i).Name
		value := v.FieldByName(key)

		valueType :=
			strings.ToLower(
				strings.Split(
					value.Type().String(), ".")[1])

		fmt.Printf("******************%v: %v - %v \n", key, value, valueType)

		url := fmt.Sprintf("http://%v:%v/update/%v/%v/%v",
			cfg.Host, cfg.Port, valueType, key, fmt.Sprintf("%v", value))
		fmt.Println(url)

		res, err := http.Post(url, cfg.ContentType, nil)
		if err != nil {
			return err
		}
		res.Body.Close() //? зачем его закрывать, если можно не считывать?
	}
	return nil
}

type config struct {
	Intervals struct {
		Poll   time.Duration
		Report time.Duration
	}
	Host        string
	Port        string
	ContentType string
}

func checkConfig() {
	if cfg.Host == "" || cfg.Port == "" || cfg.Intervals.Poll == 0 || cfg.Intervals.Report == 0 {
		cfg = loadConfig()
	}
}

func loadConfig() config {
	file, err := os.Open("configs/agent_config.json")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	cfg := config{}
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println(err)
	}
	return cfg
}

func GetConfig() config {
	checkConfig()
	return cfg
}
