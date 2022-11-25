package storage

import (
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

func (m *MemStorage) UpdateGaugeMetrics(name, value string) error {
	g, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("Unable to parse value to gauge. Value: %v, Error: %v", value, err)
	}

	m.GaugeMetrics[name] = gauge(g)
	return nil
}

func (m *MemStorage) UpdateCounterMetrics(name, value string) error {
	g, err := strconv.Atoi(value)
	if err != nil {
		return fmt.Errorf("Unable to parse value to counter. Value: %v, Error: %v", value, err)
	}

	m.CounterMetrics[name] += counter(g)
	return nil
}

// func (m *MemStorage) Get() string {
// 	return fmt.Sprint(m)
// }
