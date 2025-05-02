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
	GetName() string
	GetDescription() string
}

type SessionHandler interface {
	CanHandleState(state models.State) bool
	HandleSession(ctx context.Context, update *tgbotapi.Update, session models.Session) error
}

type CallbackHandler interface {
	CanHandleCallback(callbackData string) bool
	HandleCallback(ctx context.Context, update *tgbotapi.Update) error
}
