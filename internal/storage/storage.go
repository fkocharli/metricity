package storage

import (
	"errors"
	"fmt"
	"strconv"
)

type (
	gauge   float64
	counter int64
)

type GaugeMetrics map[string]gauge
type CounterMetrics map[string]counter

type MemStorage struct {
	GaugeMetrics   GaugeMetrics
	CounterMetrics CounterMetrics
}

var metricsList = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

func NewStorage() *MemStorage {
	gaugeDefault := make(GaugeMetrics)
	for _, v := range metricsList {
		gaugeDefault[v] = gauge(0)
	}

	return &MemStorage{
		GaugeMetrics:   gaugeDefault,
		CounterMetrics: make(CounterMetrics),
	}
}

func (m *MemStorage) UpdateGaugeMetrics(name, value string) error {
	g, err := strconv.ParseFloat(value, 64)
	if err != nil {
		return fmt.Errorf("unable to parse value to gauge. value: %v, error: %v", value, err)
	}

	m.GaugeMetrics[name] = gauge(g)
	return nil
}

func (m *MemStorage) UpdateCounterMetrics(name, value string) error {
	g, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("unable to parse value to counter. value: %v, error: %v", value, err)
	}

	m.CounterMetrics[name] += counter(g)
	return nil
}

func (m *MemStorage) GetGaugeMetrics(name string) (string, error) {
	v, ok := m.GaugeMetrics[name]
	if !ok {
		return "", errors.New("Metric Not Found")
	}
	return fmt.Sprintf("%f", v), nil
}

func (m *MemStorage) GetCounterMetrics(name string) (string, error) {
	v, ok := m.CounterMetrics[name]
	if !ok {
		return "", errors.New("Metric Not Found")
	}
	return fmt.Sprintf("%d", v), nil
}
