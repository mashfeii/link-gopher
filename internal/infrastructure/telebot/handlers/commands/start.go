package commands

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type StartHandler struct{}

func (h *StartHandler) Handle(ctx context.Context, hctx *handlers.HandlerContext) error {
	chatID := hctx.ChatID
	serviceContext := context.WithValue(ctx, models.ContextKeyChatID, chatID)

	if err := hctx.ScrapperService.RegisterUser(serviceContext); err != nil {
		return err
	}

	_, err := hctx.MessageService.Send(
		chatID,
		"🥳 Successfully registered!",
		&models.SendMessageOptions{
			ReplyMarkup: ui.NewReplyKeyboardBuilder().Button("/help - Get list of commands").Build(),
		},
	)

	return err
}
