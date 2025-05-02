package middleware

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Wrapper func(handlers.CommandHandler) handlers.CommandHandler

func WithAuth(scrapperService core.ScrapperService, messageService core.MessageService) Wrapper {
	return func(next handlers.CommandHandler) handlers.CommandHandler {
		return &authDecorator{
			next:            next,
			scrapperService: scrapperService,
			messageService:  messageService,
		}
	}
}

type authDecorator struct {
	next            handlers.CommandHandler
	scrapperService core.ScrapperService
	messageService  core.MessageService
}

func (d *authDecorator) GetName() string        { return d.next.GetName() }
func (d *authDecorator) GetDescription() string { return d.next.GetDescription() }
func (d *authDecorator) CanHandle(ctx context.Context, update *tgbotapi.Update) bool {
	return d.next.CanHandle(ctx, update)
}

func (d *authDecorator) Handle(ctx context.Context, update *tgbotapi.Update) error {
	const op = "authDecorator.Handle"

	chatID := update.Message.Chat.ID
	serviceContext := context.WithValue(ctx, core.ContextKeyChatID, chatID)

	ok, err := d.scrapperService.GetUser(serviceContext)
	if !ok {
		slog.Error("user not found", slog.String("operation", op), slog.Int64("chat_id", chatID))
		return fmt.Errorf("%s: %w", op, errors.NewErrUserNotFound(chatID))
	}

	if err != nil {
		slog.Error("failed to get user", slog.String("operation", op), slog.Int64("chat_id", chatID), slog.Any("error", err))
		return fmt.Errorf("%s: %w", op, err)
	}

	return d.next.Handle(ctx, update)
}
