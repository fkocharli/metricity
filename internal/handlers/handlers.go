package handlers

import (
	"html/template"
	"net/http"
	"strconv"

	"github.com/fkocharli/metricity/internal/server"
	"github.com/fkocharli/metricity/internal/storage"
	"github.com/go-chi/chi/v5"
)

const temp = `
<html>
	<ul>
		{{range $key, $value := .}}
			<li><strong>{{$key}}:</strong> {{$value}}</li>
		{{end}}
	</ul>
</html>
`

type ServerHandlers struct {
	*chi.Mux
	Repository Repository
}

type Repository interface {
	UpdateGaugeMetrics(name, value string) error
	UpdateCounterMetrics(name, value string) error
	GetGaugeMetrics(name string) (string, error)
	GetCounterMetrics(name string) (string, error)
	GetAll() map[string]string
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
	sh.Mux.Get("/value/{type}/{metricname}", sh.get)
	sh.Mux.Get("/", sh.home)
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

func (s *ServerHandlers) get(w http.ResponseWriter, r *http.Request) {
	t := chi.URLParam(r, "type")
	n := chi.URLParam(r, "metricname")

	var value string

	switch t {
	case "counter":
		v, err := s.Repository.GetCounterMetrics(n)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		value = v
	case "gauge":
		v, err := s.Repository.GetGaugeMetrics(n)
		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		value = v
	default:
		w.WriteHeader(http.StatusNotImplemented)
		return
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	w.Write([]byte(value))
}

func (s *ServerHandlers) home(w http.ResponseWriter, r *http.Request) {
	t := template.Must(template.New("").Parse(temp))
	data := s.Repository.GetAll()
	w.Header().Add("Content-Type", "text/html")
	t.Execute(w, data)
	w.WriteHeader(http.StatusOK)
}
