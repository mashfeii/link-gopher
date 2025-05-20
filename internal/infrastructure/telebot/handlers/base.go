package handlers

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HandlerContext struct {
	ChatID          int64
	Update          *tgbotapi.Update
	Session         models.Session
	BotAPI          *tgbotapi.BotAPI
	MessageService  core.MessageService
	SessionService  core.SessionManager
	ScrapperService core.ScrapperService
}

type Handler interface {
	Handle(ctx context.Context, hctx *HandlerContext) error
}

type HandlerFunc func(ctx context.Context, hctx *HandlerContext) error

func (f HandlerFunc) Handle(ctx context.Context, hctx *HandlerContext) error {
	return f(ctx, hctx)
}
