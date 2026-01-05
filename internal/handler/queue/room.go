package queue

import (
	"CRUDQueue/internal/hub"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

type hubService interface {
	GetHub(uuid uuid.UUID) (*hub.Hub, error)
}

func HandleRoom(service hubService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.HandleRoom"
		logger.Info("New connection into the room")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Warn("Failed to upgrade to websocket connection: " + err.Error())
			return
		}

		id := chi.URLParam(r, "uuid")
		logger.Info(fmt.Sprintf("%s. Get uuid=%s", op, id))

		u, err := uuid.Parse(id)
		if err != nil {
			logger.Info("Invalid uuid. Error: " + err.Error())
			return
		}

		h, err := service.GetHub(u)
		if err != nil {
			logger.Info("Error getting queue: " + err.Error())
		}

		h.Register <- conn

		names := []string{}
		for e := h.Queue.List.Front(); e != nil; e = e.Next() {
			names = append(names, e.Value.(string))
		}

		state, _ := json.Marshal(names)

		h.Broadcast <- state

		if err != nil {
			logger.Warn("Failed to write to websocket connection: " + err.Error())
			return
		}

		go func() {
			defer func() {
				h.Unregister <- conn
			}()

			for {
				if _, _, err := conn.ReadMessage(); err != nil {
					break
				}
			}
		}()

	}
}
