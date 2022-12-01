package main

import (
	"context"
	"log"
	"sync"

	"github.com/fkocharli/metricity/internal/handlers"
	"github.com/fkocharli/metricity/internal/server"
)

const (
	address = "127.0.0.1:8080"
)

func main() {

	repo := handlers.NewRepository()
	handler := handlers.NewHandler(repo)

	serv := server.New(address, handler.Mux)

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
