package HubInMemoryRepo

import (
	"CRUDQueue/internal/hub"
	"CRUDQueue/internal/queue"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

var (
	HubNotFound = errors.New("hub not found")
)

type RepositoryImpl struct {
	logger *slog.Logger
	hubs   map[uuid.UUID]*hub.Hub
	mutex  sync.RWMutex
}

func New(logger *slog.Logger) *RepositoryImpl {
	return &RepositoryImpl{
		logger: logger,
		hubs:   make(map[uuid.UUID]*hub.Hub),
		mutex:  sync.RWMutex{},
	}
}

func (r *RepositoryImpl) CreateHub(name, password string, queue *queue.Queue) (uuid.UUID, error) {
	const op = "HubInmemoryRepo.CreateHub"

	u, err := uuid.NewRandom()
	if err != nil {
		r.logger.Error(fmt.Sprintf("%s. Error=%s", op, err.Error()))
		return uuid.Nil, err
	}

	h := hub.New(u, name, password, queue)
	r.mutex.Lock()
	r.hubs[u] = h
	r.mutex.Unlock()

	go h.Listen()

	return u, nil
}

func (r *RepositoryImpl) GetHub(uuid uuid.UUID) (*hub.Hub, error) {
	const op = "HubInMemoryRepo.GetHub"

	r.mutex.Lock()
	h, ok := r.hubs[uuid]
	r.mutex.Unlock()

	if !ok {
		r.logger.Error(fmt.Sprintf("%s. Hub with uuid=%s not found", op, uuid.String()))
		return nil, HubNotFound
	}

	r.logger.Info(fmt.Sprintf("%s. Get hub with uuid=%s", op, uuid.String()))
	return h, nil
}
