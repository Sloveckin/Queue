package queue

import (
	"CRUDQueue/internal/queue"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var (
	PasswordNotCorrect = errors.New("password not correct")
)

type Repository interface {
	CreateQueue(name, password *string) (uuid.UUID, error)
	GetQueue(id uuid.UUID) (*queue.Queue, error)
}

type ServiceQueueImpl struct {
	logger *slog.Logger
	repo   Repository
}

func New(repo Repository, logger *slog.Logger) *ServiceQueueImpl {
	return &ServiceQueueImpl{
		logger: logger,
		repo:   repo,
	}
}

func (s *ServiceQueueImpl) CreateQueue(name, password *string) (uuid.UUID, error) {
	const op = "QueueService.CreateQueue"

	bytes, err := bcrypt.GenerateFromPassword([]byte(*password), 10)
	if err != nil {
		return uuid.Nil, err
	}

	hashPassword := string(bytes)
	u, err := s.repo.CreateQueue(name, &hashPassword)
	if err != nil {
		return uuid.Nil, err
	}

	s.logger.Info(fmt.Sprintf("%s: Create queue with name=%s and password=%s", op, u.String(), hashPassword))

	return u, nil
}

func (s *ServiceQueueImpl) AddUser(uuid uuid.UUID, name *string) error {
	const op = "QueueService.Add"
	q, err := s.repo.GetQueue(uuid)
	if err != nil {
		return err
	}

	err = q.Add(name)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s: Error: %s", op, err.Error()))
		return err
	}

	s.logger.Info(fmt.Sprintf("%s: Uuid=%s. Add new user with name=%s", op, uuid.String(), *name))
	return nil
}

func (s *ServiceQueueImpl) NextUser(uuid uuid.UUID, password *string) error {
	const op = "QueueService.NextUser"
	q, err := s.repo.GetQueue(uuid)
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s: Error: %s", op, err.Error()))
		return err
	}

	err = bcrypt.CompareHashAndPassword([]byte(q.Password), []byte(*password))
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s: Error: %s", op, PasswordNotCorrect.Error()))
		return PasswordNotCorrect
	}

	err = q.Next()
	if err != nil {
		s.logger.Info(fmt.Sprintf("%s: Error: %s", op, err.Error()))
		return err
	}

	return nil
}
