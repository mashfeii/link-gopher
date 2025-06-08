package repository

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/domain/models"
)

type UpdateRepository interface {
	SaveUpdate(ctx context.Context, update *models.Update) error
	GetUpdateByLink(ctx context.Context, linkID int64, limit, offset int) ([]models.Update, error)
}
