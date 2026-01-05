package hub

import (
	"CRUDQueue/internal/queue"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Hub struct {
	Uuid       uuid.UUID
	Name       string
	Password   string
	Queue      *queue.Queue
	Users      map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

func New(u uuid.UUID, name, password string, queue *queue.Queue) *Hub {
	return &Hub{
		Uuid:       u,
		Name:       name,
		Password:   password,
		Queue:      queue,
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		Broadcast:  make(chan []byte),
		Users:      make(map[*websocket.Conn]bool),
	}
}

func (hub *Hub) Listen() {
	for {
		select {
		case conn := <-hub.Register:
			hub.Users[conn] = true
		case conn := <-hub.Unregister:
			delete(hub.Users, conn)
		case message := <-hub.Broadcast:
			for conn := range hub.Users {
				// TODO: Check error
				conn.WriteMessage(websocket.TextMessage, message)
			}
		}
	}
}
