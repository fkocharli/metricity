package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"reflect"
	"runtime"
	"time"
)

type (
	gauge   float64
	counter int64
)

type MetricValues struct {
	Gauge     map[string]gauge
	PollCount counter
}

const (
	pollInterval   = time.Duration(2 * time.Second)
	reportInterval = time.Duration(10 * time.Second)
	address        = "127.0.0.1"
	port           = "8080"
)

var (
	metricList   = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}
	poolticket   = time.NewTicker(pollInterval)
	reportTicker = time.NewTicker(reportInterval)
)

func main() {
	currentMetricsValue := newMetricValues(metricList)

	RunAgent(currentMetricsValue)

}

func collectMetrics(metricList *MetricValues) {

	var ms runtime.MemStats
	runtime.ReadMemStats(&ms)
	v := reflect.ValueOf(&ms).Elem()

	metricList.PollCount++

	for k := range metricList.Gauge {
		if k == "RandomValue" {
			metricList.Gauge[k] = gauge(rand.Float64())
		}
		for i := 0; i < v.NumField(); i++ {
			field := v.Type().Field(i).Name
			typ := v.Type().Field(i).Type.Kind()

			if field == k && typ == reflect.Uint64 {
				val := v.Field(i).Uint()
				metricList.Gauge[k] = gauge(val)
			}
		}
	}
}

func sendMetrics(request *http.Request, client http.Client) error {
	request.Header.Add("Content-Type", "text/plain")
	res, err := client.Do(request)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	_, err = io.Copy(io.Discard, res.Body)
	if err != nil {
		return err
	}
	return nil
}

func newMetricValues(metricList []string) *MetricValues {
	m := &MetricValues{
		Gauge:     make(map[string]gauge),
		PollCount: 0,
	}
	for _, v := range metricList {
		m.Gauge[v] = 0
	}
	return m
}

func RunAgent(metrics *MetricValues) {
	client := http.Client{
		Timeout: 10 * time.Second,
	}
	ctx := context.Background()

	for {
		select {
		case <-poolticket.C:
			collectMetrics(metrics)
		case <-reportTicker.C:
			for k, v := range metrics.Gauge {
				url := fmt.Sprintf("http://%s:%s/update/gauge/%s/%s", address, port, k, fmt.Sprint(v))
				req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
				if err != nil {
					log.Printf("Unable send metric for URL: %s \n Error: %s", url, err)
				}
				err = sendMetrics(req, client)
				if err != nil {
					log.Printf("Unable send metric for URL: %s \n Error: %s", url, err)
				}
			}
			url := fmt.Sprintf("http://%s:%s/update/counter/PollCount/%s", address, port, fmt.Sprint(metrics.PollCount))
			req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, nil)
			if err != nil {
				log.Printf("Unable send metric for URL: %s \n Error: %s", url, err)
			}
			err = sendMetrics(req, client)
			if err != nil {
				log.Printf("Unable send metric for URL: %s \n Error: %s", url, err)
			}

		}
	}
}
