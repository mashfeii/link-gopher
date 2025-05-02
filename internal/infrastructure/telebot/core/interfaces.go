package core

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type SessionManager interface {
	Get(chatID int64) models.Session
	Set(chatID int64, session models.Session)
	Clear(chatID int64)
}

type ScrapperService interface {
	RegisterUser(ctx context.Context) error
	GetUser(ctx context.Context) (bool, error)
	SaveLink(ctx context.Context, linkData models.LinkData) error
	GetLinks(ctx context.Context) ([]models.LinkData, error)
	DeleteLink(ctx context.Context) error
}

type MessageService interface {
	Send(chatID int64, text string, opts *models.SendMessageOptions) (tgbotapi.Message, error)
	SendError(chatID int64, err error) (tgbotapi.Message, error)
	EditMessage(chatID int64, messageID int, text string, markup any) (tgbotapi.Message, error)
	EditText(chatID int64, messageID int, text string) (tgbotapi.Message, error)
	EditMarkup(chatID int64, messageID int, markup any) (tgbotapi.Message, error)
	DeleteMessage(chatID int64, messageID int) (*tgbotapi.APIResponse, error)
}
