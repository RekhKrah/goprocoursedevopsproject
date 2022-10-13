package storage

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
)

type gauge float64
type counter int64

type Value struct {
	valueType string
	value     string
}
type MemStorage struct {
	storage map[string]Value
}

//type MemStorage struct {
//Alloc         gauge
//BuckHashSys   gauge
//Frees         gauge
//GCCPUFraction gauge
//GCSys         gauge
//HeapAlloc     gauge
//HeapIdle      gauge
//HeapInuse     gauge
//HeapObjects   gauge
//HeapReleased  gauge
//HeapSys       gauge
//LastGC        gauge
//Lookups       gauge
//MCacheInuse   gauge
//MCacheSys     gauge
//MSpanInuse    gauge
//MSpanSys      gauge
//Mallocs       gauge
//NextGC        gauge
//NumForcedGC   gauge
//NumGC         gauge
//OtherSys      gauge
//PauseTotalNs  gauge
//StackInuse    gauge
//StackSys      gauge
//Sys           gauge
//TotalAlloc    gauge
//PollCount     counter
//RandomValue   gauge
//}

type MetricsUrl struct {
	Name  string
	Type  string
	Value string
}

//type storage interface {
//	Update()
//}
type metricShort struct {
	Name  string
	Value string
}

func (m *MemStorage) GetMetrics() (list []metricShort) {
	for name, data := range m.storage {
		list = append(list, metricShort{
			Name:  name,
			Value: data.value,
		})
	}
	return
}

func (m *MemStorage) GetMetricValue(metricType, metricName string) (string, error) {
	metric, ok := m.storage[metricName]
	if !ok {
		return "", errors.New("Metric if not found")
	}
	if metric.valueType != metricType {
		return "", errors.Errorf("Metric with type %v is not found", metricType)
	}
	return metric.value, nil
}

func (m *MemStorage) Update(data MetricsUrl) error {
	if len(m.storage) == 0 {
		m.storage = make(map[string]Value)
	}

	v, ok := m.storage[data.Name]

	if !ok {
		m.storage[data.Name] = Value{data.Type, data.Value}
		return nil
	}

	if data.Type != m.storage[data.Name].valueType {
		return errors.Errorf("Типы полученой (%v) и имеющейся (%v) метрик не совпадают",
			data.Type, m.storage[data.Name].valueType)
	}

	if data.Type == "counter" {
		vint, err := strconv.Atoi(v.value)
		if err != nil {
			return err
		}
		if vint >= 0 {
			m.storage[data.Name] = Value{"counter", string(rune(vint + 1))}
			return nil
		}
	}
	m.storage[data.Name] = Value{data.Type, data.Value}

	return nil
}

func setValueByName(v interface{}, field string, newval interface{}) {
	r := reflect.ValueOf(v).Elem().FieldByName(field)
	r.Set(reflect.ValueOf(newval))
}
