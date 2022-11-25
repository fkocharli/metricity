package main

import (
	"context"
	"log"
	"sync"

	"github.com/fkocharli/metricity/internal/handlers"
	"github.com/fkocharli/metricity/internal/server"
	"github.com/fkocharli/metricity/internal/storage"
)

const (
	address = "127.0.0.1:8080"
)

func main() {

	repo := storage.NewRepository()
	handler := handlers.NewHandler(repo)

	mux := server.NewRouter(handler.Routes)

	serv := server.New(address, mux)

	ctx, cancel := context.WithCancel(context.Background())

	group := sync.WaitGroup{}

	group.Add(1)
	go func() {
		defer group.Done()
		if err := serv.Run(ctx); err != nil {
			log.Printf("server run error: %v", err)
			cancel()
		}
	}()

	group.Wait()
}
