package queue

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type addRequest struct {
	Uuid uuid.UUID `json:"uuid" validate:"required,uuid"`
	Name string    `json:"name" validate:"required,min=2,max=32"`
}

type addResponse struct {
	Status string `json:"status" validate:"required"`
	Error  string `json:"error,omitempty"`
}

type addService interface {
	AddUser(uuid uuid.UUID, name *string) error
}

func Add(service addService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handler queue.Add"

		var req addRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Info(fmt.Sprintf("%s: Error decoding request: %s", op, err.Error()))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, addResponse{
				Status: "error",
				Error:  err.Error(),
			})
			return
		}

		validate := validator.New()

		err = validate.Struct(req)
		if err != nil {
			logger.Info(fmt.Sprintf("%s: Error: %s", op, err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, createResponse{
				Status: "Error",
				Error:  err.Error(),
			})

			return
		}

		err = service.AddUser(req.Uuid, &req.Name)
		if err != nil {
			logger.Info(fmt.Sprintf("%s: Error adding item: %s", op, err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, createResponse{
				Status: "Error",
				Error:  err.Error(),
			})

			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, addResponse{
			Status: "OK",
		})
	}
}
