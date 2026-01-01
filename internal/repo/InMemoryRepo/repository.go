package InMemoryRepo

import (
	"CRUDQueue/internal/queue"
	"errors"
	"fmt"
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

var (
	NoSuchQueue = errors.New("no such queue")
)

type RepositoryImpl struct {
	logger *slog.Logger
	mutex  *sync.RWMutex
	data   map[uuid.UUID]*queue.Queue
}

func New(mutex *sync.RWMutex, logger *slog.Logger) *RepositoryImpl {
	return &RepositoryImpl{
		logger: logger,
		mutex:  mutex,
		data:   make(map[uuid.UUID]*queue.Queue),
	}
}

func (r *RepositoryImpl) CreateQueue(name, password *string) (uuid.UUID, error) {
	const op = "InMemoryRepo.Create"
	r.mutex.Lock()
	defer r.mutex.Unlock()

	u := uuid.New()

	r.logger.Info("U: " + u.String())

	r.data[u] = queue.New(name, password)

	r.logger.Info(fmt.Sprintf("%s: UUID=%s", op, u.String()))

	return u, nil
}

func (r *RepositoryImpl) GetQueue(id uuid.UUID) (*queue.Queue, error) {
	const op = "InMemoryRepo.Get"

	r.mutex.Lock()
	defer r.mutex.Unlock()

	q, ok := r.data[id]
	if !ok {
		r.logger.Info(fmt.Sprintf("%s: No such queue with UUID=%s", op, id.String()))
		return nil, NoSuchQueue
	}

	r.logger.Info(fmt.Sprintf("%s: UUID=%s", op, id.String()))
	return q, nil
}
