package session

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TrackInputTags struct {
	messageService core.MessageService
	sessionService core.SessionManager
}

func NewTrackInputTags(
	messageService core.MessageService,
	sessionService core.SessionManager,
) *TrackInputTags {
	return &TrackInputTags{
		messageService: messageService,
		sessionService: sessionService,
	}
}

func (h *TrackInputTags) CanHandleState(state models.State) bool {
	return state == models.StateTrackInputTags
}

func (h *TrackInputTags) HandleSession(
	_ context.Context,
	update *tgbotapi.Update,
	session models.Session,
) error {
	const op = "handlers.TrackInputTags.HandleSession"

	trackSession, ok := session.(*models.TrackSession)
	if !ok {
		slog.Error("failed to cast session", slog.String("operation", op), slog.Any("session", session))
		return fmt.Errorf("%s: %w", op, errors.NewErrInvalidSessionType())
	}

	// Retrieve the user ID from the update and entered tags
	chatID := update.Message.Chat.ID
	tags := strings.Fields(update.Message.Text)

	trackSession.Tags = tags
	trackSession.State = models.StateTrackInputFilters

	// Update the UI
	keyboard := ui.NewInlineKeyboardBuilder().
		SetButtonsPerRow(2).
		DataButton("◀️ Step back", core.CallbackReturnTags).
		DataButton("🚫 Skip", core.CallbackSkipFilters).
		Build()

	message, err := h.messageService.Send(
		chatID,
		"📋 Enter filters (_space-separated_ `key:value` _pairs_):",
		&models.SendMessageOptions{
			ParseMode:   tgbotapi.ModeMarkdown,
			ReplyMarkup: keyboard,
		},
	)
	if err != nil {
		slog.Error("failed to send message", slog.String("operation", op), slog.String("tags", strings.Join(tags, " ")), slog.Any("error", err))
		return fmt.Errorf("%s: unable to send message: %w", op, err)
	}

	slog.Info("Saved the tags for track session",
		slog.String("operation", op),
		slog.Any("filters", tags),
		slog.Int64("chatID", chatID),
	)

	// Update current session
	trackSession.LastBotMessageID = &message.MessageID
	trackSession.LastUserMessageID = &update.Message.MessageID
	h.sessionService.Set(chatID, trackSession)

	return nil
}
