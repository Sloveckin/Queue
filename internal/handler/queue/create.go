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
	Response
	Uuid *uuid.UUID `json:"uuid,omitempty"`
}

type createService interface {
	Create(name, password string) (uuid.UUID, error)
}

func Create(service createService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "Handler queue.Create"

		var req createRequest
		err := render.DecodeJSON(r.Body, &req)
		if err != nil {
			logger.Info(fmt.Sprintf("%s: Error decoding request: %s", op, err.Error()))

			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Error(err.Error()))
			return
		}

		validate := validator.New()
		err = validate.Struct(req)
		if err != nil {
			logger.Info(fmt.Sprintf("%s: Error: %s", op, err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			render.JSON(w, r, Error(err.Error()))
			return
		}

		req.Name = strings.TrimSpace(req.Name)
		req.Password = strings.TrimSpace(req.Password)

		h, err := service.Create(req.Name, req.Password)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			render.JSON(w, r, Error(err.Error()))
			return
		}

		w.WriteHeader(http.StatusCreated)
		render.JSON(w, r, createResponse{
			Response: Ok(),
			Uuid:     &h,
		})
	}
}
