package handlers

import (
	"net/http"
	"strconv"

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

	sh.Mux.Post("/update/{type}/{metricname}/{metricvalue}", sh.update)

	return sh

}

func (s *ServerHandlers) update(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "metricname")
	m := chi.URLParam(r, "metricvalue")

	if _, err := strconv.ParseFloat(m, 64); err != nil {

		w.WriteHeader(http.StatusBadRequest)
		return
	}

	switch t {
	case "counter":
		err := s.Repository.UpdateCounterMetrics(n, m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	case "gauge":
		err := s.Repository.UpdateGaugeMetrics(n, m)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)
}
