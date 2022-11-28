package handlers

import (
	"net/http"
	"regexp"
	"strconv"
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

	n, m, er := valid(r.URL.Path, w)
	if er != nil {
		return
	}

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

	n, m, er := valid(r.URL.Path, w)
	if er != nil {
		return
	}

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

func valid(s string, w http.ResponseWriter) (string, string, error) {
	p := strings.Split(s, "/")
	matched, err := regexp.MatchString(`/update/(gauge|counter)/[A-Za-z0-9]+/[0-9]`, s)
	if err != nil || !matched {
		if p[2] != "counter" && p[2] != "gauge" {
			w.WriteHeader(http.StatusNotImplemented)
			return "", "", err
		}

		if _, err := strconv.Atoi(p[4]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return "", "", err
		}
		w.WriteHeader(http.StatusNotFound)
		return "", "", err
	}

	return p[3], p[4], nil
}
