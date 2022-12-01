package storage

var metricsList = []string{"Alloc", "BuckHashSys", "Frees", "GCCPUFraction", "GCSys", "HeapAlloc", "HeapIdle", "HeapInuse", "HeapObjects", "HeapReleased", "HeapSys", "LastGC", "Lookups", "MCacheInuse", "MCacheSys", "MSpanInuse", "MSpanSys", "Mallocs", "NextGC", "NumForcedGC", "NumGC", "OtherSys", "PauseTotalNs", "StackInuse", "StackSys", "Sys", "TotalAlloc"}

type Repository interface {
	UpdateGaugeMetrics(name, value string) error
	UpdateCounterMetrics(name, value string) error
	//	Get() string
}

func NewRepository() Repository {
	gaugeDefault := make(GaugeMetrics)
	for _, v := range metricsList {
		gaugeDefault[v] = gauge(0)
	}

	return &MemStorage{
		GaugeMetrics:   gaugeDefault,
		CounterMetrics: make(CounterMetrics),
	}

}
