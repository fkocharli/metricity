package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

type MockStorageType struct {
	MockField string
}

func (m *MockStorageType) UpdateGaugeMetrics(name, value string) error {

	return nil
}

func (m *MockStorageType) UpdateCounterMetrics(name, value string) error {
	return nil
}

func TestHandlers(t *testing.T) {
	type want struct {
		contenType string
		statusCode int
	}

	type request struct {
		path    string
		method  string
		handler http.HandlerFunc
	}

	mockRepo := MockStorageType{}

	handler := NewHandler(&mockRepo)

	tests := []struct {
		name    string
		request request
		want    want
	}{
		{
			name: "Update Gauge Metrics",
			request: request{
				path:    "/update/gauge/Sys/4.063232e+06",
				method:  http.MethodPost,
				handler: handler.update,
			},
			want: want{
				contenType: "text/plain",
				statusCode: http.StatusOK,
			},
		},
		{
			name: "Update Counter Metrics",
			request: request{
				path:    "/update/counter/Counter/10",
				method:  http.MethodPost,
				handler: handler.update,
			},
			want: want{
				contenType: "text/plain",
				statusCode: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			request := httptest.NewRequest(tt.request.method, tt.request.path, nil)

			tt.request.handler(w, request)

			results := w.Result()
			defer results.Body.Close()

			assert.Equal(t, tt.want.contenType, results.Header.Get("Content-Type"))
			assert.Equal(t, tt.want.statusCode, results.StatusCode)
		})
	}
}
