package queue

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/go-chi/render"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type createRequest struct {
	Name     string `json:"name" validate:"required"`
	Password string `json:"password" validate:"required"`
}

type createResponse struct {
	Status string     `json:"status"`
	Error  string     `json:"error,omitempty"`
	Uuid   *uuid.UUID `json:"uuid,omitempty"`
}

type createService interface {
	CreateQueue(name, password *string) (uuid.UUID, error)
}

func Create(service createService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handler queue.Create"

		var req createRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Info(fmt.Sprintf("%s: Error decoding request: %s", op, err.Error()))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, createResponse{
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

		req.Name = strings.TrimSpace(req.Name)
		req.Password = strings.TrimSpace(req.Password)

		u, err := service.CreateQueue(&req.Name, &req.Password)
		if err != nil {
			// Мне кажется, что тут должен быть другой статус код
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, createResponse{
				Status: "Error",
				Error:  err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, createResponse{
			Status: "OK",
			Uuid:   &u,
		})
	}
}
