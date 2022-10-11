package storage

import (
	"github.com/pkg/errors"
	"reflect"
	"regexp"
	"strconv"
)

type gauge float64
type counter int64

type MemStorage struct {
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

type storage interface {
	Update()
}

const (
	urlParserRegexp = `update/(\S+)/(\S+)/(\S+)$`
)

func (m *MemStorage) Update(url string) error {
	data, err := parseMetricsURL(url)
	if err != nil {
		return err
	}

	if data.Name == "PollCount" {
		pc, _ := strconv.Atoi(data.Value)

		if m.PollCount >= 0 {
			m.PollCount += counter(pc)
			return nil
		}
		m.PollCount = counter(pc)
	}

	f, err := strconv.ParseFloat(data.Value, 64)
	if err != nil {
		return errors.Errorf("can't convert value (%v) of %v into float64", data.Value, data.Name)
	}
	setValueByName(m, data.Name, gauge(f))

	return nil
}

func setValueByName(v interface{}, field string, newval interface{}) {
	r := reflect.ValueOf(v).Elem().FieldByName(field)
	r.Set(reflect.ValueOf(newval))
}

type metricsUrl struct {
	Name  string
	Type  string
	Value string
}

func parseMetricsURL(url string) (metricsUrl, error) {
	re := regexp.MustCompile(urlParserRegexp)
	result := re.FindStringSubmatch(url)

	if len(result) < 4 {
		return metricsUrl{}, errors.Errorf("Can't parse url '%v'", url)
	}

	return metricsUrl{
		Name:  result[2],
		Type:  result[1],
		Value: result[3],
	}, nil
}
