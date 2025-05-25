package service

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/repository"
)

type UserService interface {
	RegisterUser(ctx context.Context, chatID int64) error
	DeleteUser(ctx context.Context, chatID int64) error
	GetUser(ctx context.Context, chatID int64) (*models.User, error)
}

type LinkMeta struct {
	Tags    []string
	Filters []string
}

type UserServiceImpl struct {
	usersRepo repository.UserRepository
}

func NewUserService(usersRepo repository.UserRepository) *UserServiceImpl {
	return &UserServiceImpl{usersRepo: usersRepo}
}

func (s *UserServiceImpl) RegisterUser(ctx context.Context, chatID int64) error {
	return s.usersRepo.AddUser(ctx, &models.User{ChatID: chatID})
}

func (s *UserServiceImpl) DeleteUser(ctx context.Context, chatID int64) error {
	return s.usersRepo.DeleteUser(ctx, chatID)
}

func (s *UserServiceImpl) GetUser(ctx context.Context, chatID int64) (*models.User, error) {
	return s.usersRepo.GetUser(ctx, chatID)
}
