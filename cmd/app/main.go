package main

import (
	"CRUDQueue/internal/config"
	handlers "CRUDQueue/internal/handler/queue"
	"CRUDQueue/internal/repo/InMemoryRepo"
	"CRUDQueue/internal/service/queue"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"sync"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cnf := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	var mutex sync.RWMutex
	repository := InMemoryRepo.New(&mutex, logger)
	service := queue.New(repository, logger)

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/create", handlers.Create(service, logger))
	r.Put("/add", handlers.Add(service, logger))
	r.Put("/next", handlers.NextUser(service, logger))

	server := &http.Server{
		Addr:        cnf.HttpServer.Address,
		ReadTimeout: cnf.HttpServer.Timeout,
		IdleTimeout: cnf.HttpServer.IdleTimeout,
		Handler:     r,
	}

	if err := server.ListenAndServe(); err != nil {
		logger.Error(fmt.Sprintf("Error while starting server. Error: %s", err.Error()))
		os.Exit(1)
	}
}
