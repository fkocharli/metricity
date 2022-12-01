package storage

import (
	"errors"
	"fmt"
	"strconv"
	"sync"
)

type (
	gauge   float64
	counter int64
)

type GaugeMetrics map[string]gauge
type CounterMetrics map[string]counter

type MemStorage struct {
	GaugeMetrics        GaugeMetrics
	GaugeMetricsMutex   *sync.RWMutex
	CounterMetrics      CounterMetrics
	CounterMetricsMutex *sync.RWMutex
}

var metricsList = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

func NewStorage() *MemStorage {
	gaugeDefault := make(GaugeMetrics)
	for _, v := range metricsList {
		gaugeDefault[v] = gauge(0)
	}

	return &MemStorage{
		GaugeMetrics:        gaugeDefault,
		GaugeMetricsMutex:   &sync.RWMutex{},
		CounterMetrics:      make(CounterMetrics),
		CounterMetricsMutex: &sync.RWMutex{},
	}
}

func (m *MemStorage) UpdateGaugeMetrics(name, value string) error {
	g, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("unable to parse value to gauge. value: %v, error: %v", value, err)
	}

	m.GaugeMetricsMutex.Lock()
	defer m.GaugeMetricsMutex.Unlock()

	m.GaugeMetrics[name] = gauge(g)
	return nil
}

func (m *MemStorage) UpdateCounterMetrics(name, value string) error {
	g, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("unable to parse value to counter. value: %v, error: %v", value, err)
	}
	m.CounterMetricsMutex.Lock()
	defer m.CounterMetricsMutex.Unlock()

	m.CounterMetrics[name] += counter(g)
	return nil
}

func (m *MemStorage) GetGaugeMetrics(name string) (string, error) {
	m.GaugeMetricsMutex.RLock()
	defer m.GaugeMetricsMutex.RUnlock()

	v, ok := m.GaugeMetrics[name]
	if !ok {
		return "", errors.New("metric not found")
	}
	return fmt.Sprintf("%.3f", v), nil
}

func (m *MemStorage) GetCounterMetrics(name string) (string, error) {
	m.CounterMetricsMutex.RLock()
	defer m.CounterMetricsMutex.RUnlock()

	v, ok := m.CounterMetrics[name]
	if !ok {
		return "", errors.New("m	etric not found")
	}
	return fmt.Sprintf("%d", v), nil
}

func (m *MemStorage) GetAll() map[string]string {
	m.GaugeMetricsMutex.RLock()
	defer m.GaugeMetricsMutex.RUnlock()

	m.CounterMetricsMutex.RLock()
	defer m.CounterMetricsMutex.RUnlock()

	res := make(map[string]string)

	for k, v := range m.GaugeMetrics {
		res[k] = fmt.Sprintf("%.3f", v)
	}

	for k, v := range m.CounterMetrics {
		res[k] = fmt.Sprintf("%d", v)
	}

	return res
}
