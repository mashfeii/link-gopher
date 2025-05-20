package callback

import (
	"context"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type ListUntrackReturnTags struct{}

func (l ListUntrackReturnTags) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.ReturnTagsHandler.Handle"

	listSession := hctx.Session.(*models.ListUntrackSession)
	listSession.SelectedTags = nil

	logger := slog.With("operation", op, "chat_id", hctx.ChatID, "tags", listSession.AvailableTags)

	_, err := hctx.MessageService.EditMessage(
		hctx.ChatID,
		*listSession.LastBotMessageID,
		ui.IconTag+" Select tags to filter:",
		ui.BuildTagsKeyboard(listSession.AvailableTags, nil),
	)
	if err != nil {
		logger.Warn("failed to edit message", "error", err)
	}

	hctx.SessionService.Set(hctx.ChatID, listSession)

	logger.Info("edited message to user", "message_id", *listSession.LastBotMessageID)

	return nil
}
