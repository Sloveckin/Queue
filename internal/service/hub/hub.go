package ServiceHub

import (
	"CRUDQueue/internal/hub"
	"CRUDQueue/internal/queue"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type hubRepository interface {
	CreateHub(name, password string, queue *queue.Queue) (uuid.UUID, error)
	GetHub(uuid uuid.UUID) (*hub.Hub, error)
}

type HubServiceImpl struct {
	logger *slog.Logger
	repo   hubRepository
}

func New(repo hubRepository, logger *slog.Logger) *HubServiceImpl {
	return &HubServiceImpl{
		logger: logger,
		repo:   repo,
	}
}

func (h *HubServiceImpl) Create(name, password string) (uuid.UUID, error) {
	const op = "Hubservice.Create"

	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		h.logger.Info(fmt.Sprintf("%s. Error=%s", op, err.Error()))
		return uuid.Nil, err
	}
	hashPassword := string(bytes)

	q := queue.New()

	u, err := h.repo.CreateHub(name, hashPassword, q)
	if err != nil {
		h.logger.Info(fmt.Sprintf("%s. Error=%s", op, err.Error()))
		return uuid.Nil, err
	}

	h.logger.Info(fmt.Sprintf("%s. Create hub with id=%s", op, u))

	return u, err
}

func (h *HubServiceImpl) AddUser(uuid uuid.UUID, name *string) error {
	const op = "Hubservice.AddUser"

	hub, err := h.repo.GetHub(uuid)
	if err != nil {
		h.logger.Info(fmt.Sprintf("%s. Error=%s", op, err.Error()))
		return err
	}

	err = hub.Queue.Add(name)
	if err != nil {
		h.logger.Info(fmt.Sprintf("%s. Error=%s", op, err.Error()))
		return err
	}

	// TODO: Need to check error
	state, _ := json.Marshal(hub.Queue.List)

	hub.Broadcast <- state

	h.logger.Info(fmt.Sprintf("%s. Added user %s", op, *name))
	return nil
}

func (h *HubServiceImpl) GetHub(uuid uuid.UUID) (*hub.Hub, error) {
	const op = "Hubservice.GetHub"

	hub, err := h.repo.GetHub(uuid)
	if err != nil {
		h.logger.Info(fmt.Sprintf("%s. Error=%s", op, err.Error()))
		return nil, err
	}

	h.logger.Info(fmt.Sprintf("%s. Get hub with id=%s", op, uuid.String()))
	return hub, nil
}
