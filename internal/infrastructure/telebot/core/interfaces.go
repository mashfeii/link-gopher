package core

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type SessionManager interface {
	Get(chatID int64) *models.UserSession
	Set(chatID int64, session *models.UserSession)
	Clear(chatID int64)
}

type ScrapperService interface {
	RegisterUser(ctx context.Context, chatID int64) error
	SaveLink(ctx context.Context, chatID int64, linkData models.LinkData) error
	GetLinks(ctx context.Context, chatID int64) ([]models.LinkData, error)
	DeleteLink(ctx context.Context, chatID int64, url string) error
}
