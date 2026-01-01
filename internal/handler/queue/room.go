package queue

import (
	"CRUDQueue/internal/queue"
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

type roomResponse struct {
	Names []string `json:"names"`
}

type queueService interface {
	GetQueue(uuid uuid.UUID) (*queue.Queue, error)
	CheckExist(uuid uuid.UUID) (bool, error)
}

func HandleRoom(service queueService, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "handlers.HandleRoom"
		logger.Info("New connection into the room")

		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			logger.Warn("Failed to upgrade to websocket connection: " + err.Error())
			return
		}

		defer conn.Close()

		id := chi.URLParam(r, "uuid")
		logger.Info(fmt.Sprintf("%s. Get uuid=%s", op, id))

		u, err := uuid.Parse(id)
		if err != nil {
			logger.Info("Invalid uuid. Error: " + err.Error())
			return
		}

		exists, err := service.CheckExist(u)
		if err != nil {
			logger.Info("Error checking existence: " + err.Error())
			return
		}

		if !exists {
			logger.Info("Room does not exist")
			return
		}

		q, err := service.GetQueue(u)
		if err != nil {
			logger.Info("Error getting queue: " + err.Error())
		}

		var response roomResponse

		curNode := q.List.Front()
		for i := 0; i < q.List.Len(); i++ {
			response.Names = append(response.Names, curNode.Value.(string))
			curNode = curNode.Next()
		}

		conn.WriteJSON(response)
	}
}
