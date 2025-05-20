package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type TrackReturnTags struct{}

func (h *TrackReturnTags) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.TrackReturnTags.HandleCallback"

	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", hctx.ChatID))
	trackSession := hctx.Session.(*models.TrackSession)

	logger.Info("retrieve session for callback update")

	if trackSession.LastBotMessageID == nil {
		logger.Warn("no last bot message ID found")
		return fmt.Errorf("%s: unable to find last bot message ID", op)
	}

	_, err := hctx.MessageService.EditMessage(
		hctx.ChatID,
		*trackSession.LastBotMessageID,
		ui.IconTag+" Enter tags (_space-separated_):",
		ui.GetBackSkipKeyboard(models.CallbackReturnTags, models.CallbackSkipTags),
	)
	if err != nil {
		logger.Error("failed to edit message", slog.Any("error", err))
		return fmt.Errorf("%s: unable to edit message: %w", op, err)
	}

	logger.Info("updated UI back to tags input")

	trackSession.State = models.StateTrackInputTags
	hctx.SessionService.Set(hctx.ChatID, trackSession)

	logger.Info("saved the state for track session")

	return nil
}
