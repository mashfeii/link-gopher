package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type TrackSkipTags struct{}

func (h *TrackSkipTags) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.TrackSkipTags.HandleCallback"

	chatID := hctx.Update.CallbackQuery.From.ID
	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", chatID))

	trackSession := hctx.Session.(*models.TrackSession)

	if trackSession.LastBotMessageID == nil {
		logger.Warn("no last bot message ID found")
		return fmt.Errorf("%s: unable to find last bot message ID", op)
	}

	_, err := hctx.MessageService.EditMessage(
		chatID,
		*trackSession.LastBotMessageID,
		ui.IconSearch+" Enter filters (_space-separated_ `key:value` _pairs_):",
		ui.GetBackSkipKeyboard(models.CallbackReturnTags, models.CallbackSkipFilters),
	)
	if err != nil {
		logger.Error("failed to send message", slog.Any("error", err))
		return fmt.Errorf("%s: unable to send message: %w", op, err)
	}

	// Update current session
	trackSession.State = models.StateTrackInputFilters
	hctx.SessionService.Set(chatID, trackSession)

	logger.Info("updated UI forward to filters input")

	return nil
}
