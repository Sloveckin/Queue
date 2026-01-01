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

	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

			// для preflight-запросов
			if r.Method == http.MethodOptions {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	r.Use(middleware.RequestID)
	r.Use(middleware.Recoverer)

	r.Post("/create", handlers.Create(service, logger))
	r.Put("/join", handlers.Join(service, logger))
	r.Put("/next", handlers.NextUser(service, logger))
	r.Get("/queues/{uuid}/ws", handlers.HandleRoom(service, logger))

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
