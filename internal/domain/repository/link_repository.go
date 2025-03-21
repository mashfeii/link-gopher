package repository

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
)

type LinkRepository interface {
	AddLink(ctx context.Context, link *models.Link) (*int64, error)
	DeleteLink(ctx context.Context, chatID int64, url string) (*models.Link, error)
	GetLinksByUser(ctx context.Context, chatID int64) ([]models.Link, error)
	GetLinkByURL(ctx context.Context, chatID int64, url string) (*models.Link, error)
	GetAllActiveLinks(ctx context.Context) ([]models.Link, error)
}
