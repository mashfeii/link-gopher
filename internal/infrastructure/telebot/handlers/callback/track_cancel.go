package callback

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type TrackCancel struct{}

func (h *TrackCancel) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.TrackCancel.HandleCallback"

	chatID := hctx.Update.CallbackQuery.From.ID
	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", chatID))

	trackSession := hctx.Session.(*models.TrackSession)

	logger = logger.With(
		slog.String("filters", strings.Join(trackSession.Filters, ",")),
		slog.String("tags", strings.Join(trackSession.Tags, ",")),
		slog.String("url", trackSession.URL),
	)

	if trackSession.LastBotMessageID != nil {
		if _, err := hctx.MessageService.EditMarkup(chatID, *trackSession.LastBotMessageID, nil); err != nil {
			logger.Warn("failed to edit message", slog.Any("error", err))
		}
	}

	_, err := hctx.MessageService.Send(
		chatID,
		ui.IconCross+" Tracking canceled, start again with /track",
		nil,
	)
	if err != nil {
		logger.Error("failed to send message", slog.Any("error", err))

		return fmt.Errorf("%s: unable to send message: %w", op, err)
	}

	logger.Info("Link tracking canceled")

	hctx.SessionService.Clear(chatID)

	return nil
}
