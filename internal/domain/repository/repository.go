package repository

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
)

type UnifiedRepository interface {
	CreateUser(ctx context.Context, chatID int64) error
	GetUser(ctx context.Context, chatID int64) (*models.User, error)
	DeleteUser(ctx context.Context, chatID int64) (*models.User, error)

	AddLink(ctx context.Context, chatID int64, link *models.Link) error
	DeleteLink(ctx context.Context, chatID int64, url string) (*models.Link, error)
	GetLinks(ctx context.Context, chatID int64) ([]models.Link, error)
	GetAllActiveLinks(ctx context.Context) ([]models.Link, error)

	AddTagToLink(ctx context.Context, chatID int64, url string, tag string) error
	DeleteTagFromLink(ctx context.Context, chatID int64, url string, tag string) error

	UpdateLinkFilters(ctx context.Context, chatID int64, url string, filters map[string][]string) error
}
