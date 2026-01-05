package main

import (
	"CRUDQueue/internal/config"
	handlers "CRUDQueue/internal/handler/queue"
	HubInMemoryRepo "CRUDQueue/internal/repo/hub/InMemoryRepo"
	ServiceHub "CRUDQueue/internal/service/hub"
	"fmt"
	"log/slog"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	cnf := config.MustLoad()
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	hubRepository := HubInMemoryRepo.New(logger)
	hubService := ServiceHub.New(hubRepository, logger)

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

	r.Post("/create", handlers.Create(hubService, logger))
	r.Put("/join", handlers.Join(hubService, logger))
	r.Get("/queues/{uuid}/ws", handlers.HandleRoom(hubService, logger))

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
