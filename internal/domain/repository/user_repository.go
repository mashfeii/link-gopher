package repository

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
)

type UserRepository interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUser(ctx context.Context, chatID int64) (*models.User, error)
	DeleteUser(ctx context.Context, chatID int64) error
}
