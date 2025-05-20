package commands

import (
	"context"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
)

type CancelHandler struct{}

func (h *CancelHandler) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "commands.CancelHandler.Handle"

	chatID := hctx.ChatID

	session := hctx.SessionService.Get(chatID)

	if messageID := session.GetLastBotMessageID(); messageID != nil {
		_, err := hctx.MessageService.EditMarkup(chatID, *messageID, nil)
		if err != nil {
			slog.Warn("failed to edit message", slog.String("operation", op), slog.Any("error", err))
		}
	}

	hctx.SessionService.Clear(chatID)

	_, err := hctx.MessageService.Send(chatID, "🚮 Operation canceled", nil)

	return err
}
