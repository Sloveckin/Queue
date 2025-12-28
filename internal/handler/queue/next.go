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

type nextRequest struct {
	Uuid     uuid.UUID `json:"uuid" validate:"required,uuid"`
	Password string    `json:"password" validate:"required"`
}

type nextResponse struct {
	Status string `json:"status" validate:"required"`
	Error  string `json:"error,omitempty"`
}

type nextService interface {
	NextUser(uuid uuid.UUID, password *string) error
}

func NextUser(service nextService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handler.User.NextUser"

		var req nextRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Info(fmt.Sprintf("%s: Error decoding request: %s", op, err.Error()))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, nextResponse{
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
			render.JSON(w, r, nextResponse{
				Status: "Error",
				Error:  err.Error(),
			})

			return
		}

		req.Password = strings.TrimSpace(req.Password)

		err = service.NextUser(req.Uuid, &req.Password)
		if err != nil {
			// Мне кажется, что тут должен быть другой статус код
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, nextResponse{
				Status: "Error",
				Error:  err.Error(),
			})
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, nextResponse{
			Status: "OK",
		})

	}
}
