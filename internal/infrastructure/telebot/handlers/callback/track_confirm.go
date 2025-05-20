package callback

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type TrackConfirm struct{}

func (h *TrackConfirm) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.TrackConfirm.HandleCallback"

	chatID := hctx.Update.CallbackQuery.From.ID
	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", chatID))

	trackSession := hctx.Session.(*models.TrackSession)

	logger = logger.With(
		slog.String("filters", strings.Join(trackSession.Filters, ",")),
		slog.String("tags", strings.Join(trackSession.Tags, ",")),
		slog.String("url", trackSession.URL),
	)
	serviceContext := context.WithValue(context.Background(), models.ContextKeyChatID, chatID)

	err := hctx.ScrapperService.SaveLink(serviceContext, models.LinkData{
		URL:     trackSession.URL,
		Tags:    trackSession.Tags,
		Filters: trackSession.Filters,
	})
	if err != nil {
		logger.Error("failed to save link", slog.Any("error", err))

		return fmt.Errorf("%s: unable to save link: %w", op, err)
	}

	if trackSession.LastBotMessageID != nil {
		if _, err := hctx.MessageService.EditMarkup(chatID, *trackSession.LastBotMessageID, nil); err != nil {
			logger.Warn("failed to edit message", slog.Any("error", err))
		}
	}

	_, err = hctx.MessageService.Send(
		chatID,
		"🚀 Tracking saved! You will receive notifications when the link is updated.",
		nil,
	)
	if err != nil {
		logger.Error("failed to send message", slog.Any("error", err))

		return fmt.Errorf("%s: unable to send message: %w", op, err)
	}

	logger.Info("Link is saved successfully")

	hctx.SessionService.Clear(chatID)

	return nil
}
