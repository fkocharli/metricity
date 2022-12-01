package server

import (
	"context"
	"net/http"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

type Server struct {
	server *http.Server
}

func New(address string, handler *chi.Mux) *Server {
	return &Server{
		server: &http.Server{
			Handler: handler,
			Addr:    address,
		},
	}
}

func NewRouter() *chi.Mux {

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.RealIP)
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	return r
}

func (s *Server) Run(ctx context.Context) (err error) {

	errChan := make(chan error, 1)

	group := &sync.WaitGroup{}
	group.Add(1)
	go func() {
		defer group.Done()
		if err := s.server.ListenAndServe(); err != nil {
			errChan <- err
		}
	}()

	select {
	case <-ctx.Done():
		err = s.server.Shutdown(ctx)
	case err = <-errChan:
	}

	group.Wait()

	return err
}
