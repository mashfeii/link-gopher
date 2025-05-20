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

type TrackSkipFilters struct{}

func (h *TrackSkipFilters) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.TrackSkipFilters.HandleCallback"

	chatID := hctx.Update.CallbackQuery.From.ID
	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", chatID))

	trackSession := hctx.Session.(*models.TrackSession)

	if trackSession.LastBotMessageID == nil {
		logger.Warn("no last bot message ID found")
		return fmt.Errorf("%s: unable to find last bot message ID", op)
	}

	message := fmt.Sprintf(ui.IconChecked+" Ready to track:\n\n• _URL_: %s\n• _Tags_: %s\n• _Filters_: %s\n\n*Confirm?*",
		trackSession.URL,
		strings.Join(trackSession.Tags, ", "),
		strings.Join(trackSession.Filters, ", "),
	)

	keyboard := ui.NewInlineKeyboardBuilder().
		SetButtonsPerRow(2).
		DataButton(ui.IconCross+" Cancel", models.CallbackCancelTrack).
		DataButton(ui.IconChecked+" Confirm", models.CallbackConfirmTrack).
		Build()

	_, err := hctx.MessageService.EditMessage(
		chatID,
		*trackSession.LastBotMessageID,
		message,
		keyboard,
	)
	if err != nil {
		logger.Error("failed to send message", slog.Any("error", err))
		return fmt.Errorf("%s: unable to send message: %w", op, err)
	}

	// Update current session
	trackSession.State = models.StateTrackConfirm
	hctx.SessionService.Set(chatID, trackSession)

	logger.Info("updated UI forward to filters input")

	return nil
}
