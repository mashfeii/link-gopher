package handlers

import (
	"context"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Handler interface {
	CanHandle(ctx context.Context, update *tgbotapi.Update) bool
	Handle(ctx context.Context, update *tgbotapi.Update) error
}

type CommandHandler interface {
	Handler
	GetCommandName() string
}

type SessionHandler interface {
	Handler
	CanHandleState(state models.State) bool
	HandleSession(ctx context.Context, update *tgbotapi.Update, session *models.UserSession) error
}
