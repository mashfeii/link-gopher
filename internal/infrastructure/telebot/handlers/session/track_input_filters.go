package session

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type TrackInputFilters struct{}

func (h *TrackInputFilters) Handle(
	_ context.Context,
	hctx *handlers.HandlerContext,
) error {
	const op = "handlers.TrackInputFilters.HandleSession"

	trackSession := hctx.Session.(*models.TrackSession)

	filters := lo.Uniq(strings.Fields(hctx.Update.Message.Text))

	trackSession.Filters = filters
	trackSession.State = models.StateTrackConfirm

	logger := slog.With(
		slog.String("operation", op),
		slog.Int64("chatID", hctx.ChatID),
		slog.String("filters", strings.Join(filters, ", ")),
	)

	if _, err := hctx.MessageService.EditMarkup(hctx.ChatID, *trackSession.LastBotMessageID, nil); err != nil {
		logger.Warn("failed to edit message", slog.Any("error", err))
	}

	logger.Info("Saved the filters for track session")

	trackSession.LastUserMessageID = &hctx.Update.Message.MessageID
	hctx.SessionService.Set(hctx.ChatID, trackSession)

	message := fmt.Sprintf(ui.IconCheckbox+" Ready to track:\n\n• _URL_: %s\n• _Tags_: %s\n• _Filters_: %s\n\n*Confirm?*",
		trackSession.URL,
		strings.Join(trackSession.Tags, ", "),
		strings.Join(trackSession.Filters, ", "),
	)

	keyboard := ui.NewInlineKeyboardBuilder().
		SetButtonsPerRow(2).
		DataButton(ui.IconCross+" Cancel", models.CallbackCancelTrack).
		DataButton(ui.IconChecked+" Confirm", models.CallbackConfirmTrack).
		Build()

	msg, err := hctx.MessageService.Send(hctx.ChatID, message, &models.SendMessageOptions{
		ParseMode:   tgbotapi.ModeMarkdown,
		ReplyMarkup: keyboard,
	})
	if err != nil {
		return fmt.Errorf("%s: failed to send confirmation message: %w", op, err)
	}

	trackSession.LastBotMessageID = &msg.MessageID
	hctx.SessionService.Set(hctx.ChatID, trackSession)

	return nil
}
