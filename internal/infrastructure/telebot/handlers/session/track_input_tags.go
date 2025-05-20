package session

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/samber/lo"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TrackInputTags struct{}

func (h *TrackInputTags) Handle(
	_ context.Context,
	hctx *handlers.HandlerContext,
) error {
	const op = "handlers.TrackInputTags.HandleSession"

	trackSession := hctx.Session.(*models.TrackSession)

	// Retrieve the user ID from the update and entered tags
	tags := lo.Uniq(strings.Fields(hctx.Update.Message.Text))

	trackSession.Tags = tags
	trackSession.State = models.StateTrackInputFilters

	logger := slog.With(
		slog.String("operation", op),
		slog.Int64("chatID", hctx.ChatID),
		slog.String("tags", strings.Join(tags, ", ")),
	)

	if _, err := hctx.MessageService.EditMarkup(hctx.ChatID, *trackSession.LastBotMessageID, nil); err != nil {
		logger.Warn("failed to edit message", slog.Any("error", err))
	}

	// Update the UI
	message, err := hctx.MessageService.Send(
		hctx.ChatID,
		ui.IconSearch+" Enter filters (_space-separated_ `key:value` _pairs_):",
		&models.SendMessageOptions{
			ParseMode:   tgbotapi.ModeMarkdown,
			ReplyMarkup: ui.GetBackSkipKeyboard(models.CallbackReturnTags, models.CallbackSkipFilters),
		},
	)
	if err != nil {
		logger.Error("failed to send message", slog.Any("error", err))
		return fmt.Errorf("%s: unable to send message: %w", op, err)
	}

	logger.Info("Saved the tags for track session")

	// Update current session
	trackSession.LastBotMessageID = &message.MessageID
	trackSession.LastUserMessageID = &hctx.Update.Message.MessageID
	hctx.SessionService.Set(hctx.ChatID, trackSession)

	return nil
}
