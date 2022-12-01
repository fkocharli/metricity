package server

import (
	"context"
	"net/http"
	"sync"
)

type Server struct {
	server *http.Server
}

type Route struct {
	Path    string
	Handler http.HandlerFunc
}

func New(address string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Handler: handler,
			Addr:    address,
		},
	}
}

func NewRouter(r []Route) *http.ServeMux {
	mux := http.NewServeMux()
	for _, v := range r {
		mux.HandleFunc(v.Path, v.Handler)
	}

	return mux
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
