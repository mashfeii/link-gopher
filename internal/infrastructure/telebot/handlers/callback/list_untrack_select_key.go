package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type ListUntrackSelectKey struct{}

func (l ListUntrackSelectKey) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.ListUntrackSelectKey.Handle"

	listSession := hctx.Session.(*models.ListUntrackSession)
	listSession.FilterName = hctx.Update.CallbackData()[len(models.PrefixSelectKey):]

	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", hctx.ChatID))
	logger.Info("User selected filter key", slog.String("filter key", listSession.FilterName))

	if _, err := hctx.MessageService.EditMessage(
		hctx.ChatID,
		*listSession.LastBotMessageID,
		ui.IconSearch+" Please select the filter value",
		ui.BuildFiltersValueKeyboard(listSession.AvailableFilters[listSession.FilterName]),
	); err != nil {
		logger.Error("Failed to update message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to update message: %w", op, err)
	}

	logger.Info("Tags done, handling filter selection")

	return nil
}
