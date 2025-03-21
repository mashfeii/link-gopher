package storage

import (
	"context"
	"sync"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
)

type InMemoryUserRepository struct {
	users map[int64]struct{}
	mu    sync.RWMutex
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[int64]struct{}),
	}
}

func (u *InMemoryUserRepository) AddUser(_ context.Context, user *models.User) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.users[user.ChatID]; ok {
		return errors.NewErrUserAlreadyExists(user.ChatID)
	}

	u.users[user.ChatID] = struct{}{}

	return nil
}

func (u *InMemoryUserRepository) GetUser(_ context.Context, chatID int64) (*models.User, error) {
	u.mu.RLock()
	defer u.mu.RUnlock()

	if _, ok := u.users[chatID]; !ok {
		return nil, errors.NewErrUserNotFound()
	}

	return &models.User{ChatID: chatID}, nil
}

func (u *InMemoryUserRepository) DeleteUser(_ context.Context, chatID int64) error {
	u.mu.Lock()
	defer u.mu.Unlock()

	if _, ok := u.users[chatID]; !ok {
		return errors.NewErrUserNotFound()
	}

	delete(u.users, chatID)

	return nil
}
