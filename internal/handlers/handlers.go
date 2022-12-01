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
			Path:    "/update/",
			Handler: http.HandlerFunc(sh.update),
		},
	}

	return sh

}

func (s *ServerHandlers) update(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	x, n, m, er := valid(r.URL.Path, w)
	if er != nil {
		return
	}

	switch x {
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
	}

	w.Header().Add("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

}

func valid(s string, w http.ResponseWriter) (string, string, string, error) {
	p := strings.Split(s, "/")
	matched, err := regexp.MatchString(`/update/(gauge|counter)/[A-Za-z0-9]+/[0-9]`, s)
	if err != nil || !matched {

		if p[2] != "counter" && p[2] != "gauge" {
			w.WriteHeader(http.StatusNotImplemented)
			return "", "", "", err
		}

		if len(p) < 5 {
			w.WriteHeader(http.StatusNotFound)
			return "", "", "", err
		}

		if _, err := strconv.Atoi(p[4]); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return "", "", "", err
		}
		w.WriteHeader(http.StatusNotFound)
		return "", "", "", err
	}

	return p[1], p[2], p[3], nil
}
