package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type TrackReturnURL struct{}

func (h *TrackReturnURL) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.TrackReturnURL.HandleCallback"

	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", hctx.ChatID))
	trackSession := hctx.Session.(*models.TrackSession)

	if trackSession.LastBotMessageID == nil {
		logger.Warn("no last bot message ID found")
		return fmt.Errorf("%s: unable to find last bot message ID", op)
	}

	_, err := hctx.MessageService.EditMarkup(hctx.ChatID, *trackSession.LastBotMessageID, nil)
	if err != nil {
		logger.Error("failed to remove keyboard", slog.Any("error", err))
		return fmt.Errorf("%s: unable to remove keyboard: %w", op, err)
	}

	_, err = hctx.MessageService.EditText(hctx.ChatID, *trackSession.LastBotMessageID, "🔗 Please enter the URL you want to track")
	if err != nil {
		logger.Error("failed to edit message", slog.Any("error", err))
		return fmt.Errorf("%s: unable to edit message: %w", op, err)
	}

	logger.Info("updated UI back to URL input")

	trackSession.State = models.StateTrackInputURL
	hctx.SessionService.Set(hctx.ChatID, trackSession)

	logger.Info("saved the state for track session")

	return nil
}
