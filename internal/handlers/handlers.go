package handlers

import (
	"net/http"

	"github.com/fkocharli/metricity/internal/server"
	"github.com/fkocharli/metricity/internal/storage"
	"github.com/go-chi/chi/v5"
)

type ServerHandlers struct {
	*chi.Mux
	Repository Repository
}

type Repository interface {
	UpdateGaugeMetrics(name, value string) error
	UpdateCounterMetrics(name, value string) error
}

func NewRepository() Repository {
	return storage.NewStorage()
}

func NewHandler(r Repository) *ServerHandlers {

	sh := &ServerHandlers{
		Mux:        server.NewRouter(),
		Repository: r,
	}

	sh.Mux.Route("/update", func(r chi.Router) {
		r.Post("/counter/{metricname}/{metricvalue}", sh.updateCounter)
		r.Post("/gauge/{metricname}/{metricvalue}", sh.updateGauge)
	})

	return sh

}

func (s *ServerHandlers) updateGauge(w http.ResponseWriter, r *http.Request) {

	n := chi.URLParam(r, "metricname")
	m := chi.URLParam(r, "metricvalue")

	err := s.Repository.UpdateGaugeMetrics(n, m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}

func (s *ServerHandlers) updateCounter(w http.ResponseWriter, r *http.Request) {

	n := chi.URLParam(r, "metricname")
	m := chi.URLParam(r, "metricvalue")

	err := s.Repository.UpdateCounterMetrics(n, m)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}
