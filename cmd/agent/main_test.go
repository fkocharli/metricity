package main

import (
	"math/rand"
	"testing"

	"github.com/stretchr/testify/assert"
)

var listMetrics = []string{"Alloc", "BuckHashSys", "Frees", "RandomValue"}

func Test_collectMetrics(t *testing.T) {
	listMetrics := []string{"Alloc", "BuckHashSys", "Frees", "RandomValue"}

	mValue := &MetricValues{
		Gauge:     make(map[string]gauge),
		PollCount: 0,
	}

	r := gauge(rand.Float64())
	for _, v := range listMetrics {
		if v == "RandomValue" {
			mValue.Gauge[v] = r
			continue
		}
		mValue.Gauge[v] = 0
	}

	collectMetrics(mValue)

	for _, v := range listMetrics {
		if v == "RandomValue" {
			assert.NotEqual(t, r, v)
			assert.NotNil(t, v)
			continue
		}
		assert.NotNil(t, mValue.Gauge[v])

	}

	assert.NotNil(t, mValue.PollCount)

}

func Test_newMetricValues(t *testing.T) {

	tests := []struct {
		name string
		args []string
		want *MetricValues
	}{
		struct {
			name string
			args []string
			want *MetricValues
		}{

			name: "Testing with fields",
			args: listMetrics,
			want: &MetricValues{
				Gauge: map[string]gauge{
					"Alloc":       0,
					"BuckHashSys": 0,
					"Frees":       0,
					"RandomValue": 0,
				},
				PollCount: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := newMetricValues(tt.args)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestRunAgent(t *testing.T) {
	type args struct {
		metrics *MetricValues
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			RunAgent(tt.args.metrics)
		})
	}
}
