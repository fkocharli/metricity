package handlers

import (
	"net/http"
	"strings"

	"github.com/fkocharli/metricity/internal/server"
	"github.com/fkocharli/metricity/internal/storage"
)

type ServerHandlers struct {
	Routes     []server.Route
	Repository storage.Repository
}

func NewHandler(r storage.Repository) ServerHandlers {

	sh := ServerHandlers{
		Repository: r,
	}

	sh.Routes = []server.Route{
		{
			Path:    "/update/gauge/",
			Handler: http.HandlerFunc(sh.updateGauge),
		},
		{
			Path:    "/update/counter/",
			Handler: http.HandlerFunc(sh.updateCounter),
		},
		// {
		// 	Path:    "/get",
		// 	Handler: http.HandlerFunc(sh.get),
		// },
	}

	return sh

}

func (s *ServerHandlers) updateGauge(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	n, m := getMetricsFromPath(r.URL.Path)

	err := s.Repository.UpdateGaugeMetrics(n, m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

func (s *ServerHandlers) updateCounter(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	n, m := getMetricsFromPath(r.URL.Path)

	err := s.Repository.UpdateCounterMetrics(n, m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

// func (s *ServerHandlers) get(w http.ResponseWriter, r *http.Request) {
// 	w.Write([]byte(s.Repository.Get()))
// }

func getMetricsFromPath(path string) (name, value string) {
	metrics := strings.Split(path, "/")

	name = metrics[3]
	value = metrics[4]
	return
}
