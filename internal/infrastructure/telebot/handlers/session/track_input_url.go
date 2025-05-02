package session

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TrackInputURL struct {
	messageService  core.MessageService
	sessionService  core.SessionManager
	scrapperService core.ScrapperService
}

func NewTrackInputURL(
	messageService core.MessageService,
	sessionService core.SessionManager,
	scrapperService core.ScrapperService,
) *TrackInputURL {
	return &TrackInputURL{
		messageService:  messageService,
		sessionService:  sessionService,
		scrapperService: scrapperService,
	}
}

func (h *TrackInputURL) CanHandleState(state models.State) bool {
	return state == models.StateTrackInputURL
}

func (h *TrackInputURL) HandleSession(
	ctx context.Context,
	update *tgbotapi.Update,
	session models.Session,
) error {
	const op = "handlers.TrackInputURL.HandleSession"

	// Check if the session is of type TrackSession
	trackSession, ok := session.(*models.TrackSession)
	if !ok {
		return fmt.Errorf("%s: %w", op, errors.NewErrInvalidSessionType())
	}

	// Update the last message from user
	trackSession.LastUserMessageID = &update.Message.MessageID
	// Retrieve the user ID from the update and entered URL
	URL := update.Message.Text
	chatID := update.Message.Chat.ID
	serviceContext := context.WithValue(ctx, core.ContextKeyChatID, chatID)

	// Check if the user exists
	links, err := h.scrapperService.GetLinks(serviceContext)
	if err != nil {
		slog.Error("failed to get links", slog.String("operation", op), slog.String("url", URL), slog.Any("error", err))
		return fmt.Errorf("%s: failed to get links: %w", op, err)
	}

	// Check for link duplicate
	for _, link := range links {
		if link.URL != URL {
			continue
		}

		slog.Warn("duplicate URL", slog.String("operation", op), slog.String("url", URL))

		_, err := h.messageService.EditText(
			chatID,
			*trackSession.LastBotMessageID,
			"❗️ This link already exists in your list. Please enter a different URL.",
		)
		if err != nil {
			slog.Error("failed to edit message", slog.String("operation", op), slog.String("url", URL), slog.Any("error", err))
			return fmt.Errorf("%s: failed to edit message: %w", op, err)
		}

		if _, err := h.messageService.DeleteMessage(chatID, *trackSession.LastUserMessageID); err != nil {
			slog.Warn("failed to delete message", slog.String("operation", op), slog.String("url", URL), slog.Any("error", err))
			return fmt.Errorf("%s: failed to delete message: %w", op, err)
		}

		return nil
	}

	// Validate link format
	_, _, gitErr := pkg.ValidateGithubURL(URL)
	_, stackErr := pkg.ValidateStackOverflowURL(URL)

	if gitErr != nil && stackErr != nil {
		slog.Warn("invalid URL format", slog.String("operation", op), slog.String("url", URL), slog.Any("error", gitErr))

		_, err := h.messageService.EditText(
			chatID,
			*trackSession.LastBotMessageID,
			"❗️ Invalid URL format. Please provide a valid GitHub Repo or StackOverflow Question URL.",
		)
		if err != nil {
			slog.Error("failed to edit message", slog.String("operation", op), slog.String("url", URL), slog.Any("error", err))
			return fmt.Errorf("%s: failed to edit message: %w", op, err)
		}

		if _, err := h.messageService.DeleteMessage(chatID, *trackSession.LastUserMessageID); err != nil {
			slog.Warn("failed to delete message", slog.String("operation", op), slog.String("url", URL), slog.Any("error", err))
			return fmt.Errorf("%s: failed to delete message: %w", op, err)
		}

		return nil
	}

	// Update UI
	keyboard := ui.NewInlineKeyboardBuilder().
		SetButtonsPerRow(2).
		DataButton("◀️ Step back", core.CallbackReturnURL).
		DataButton("🚫 Skip", core.CallbackSkipTags).
		Build()

	message, err := h.messageService.Send(
		chatID,
		"📋 Enter tags (_space-separated_):",
		&models.SendMessageOptions{
			ParseMode:   tgbotapi.ModeMarkdown,
			ReplyMarkup: keyboard,
		},
	)
	if err != nil {
		slog.Error("failed to send message", slog.String("operation", op), slog.String("url", URL), slog.Any("error", err))
		return fmt.Errorf("%s: failed to send message: %w", op, err)
	}

	// Update current session
	trackSession.URL = URL
	trackSession.State = models.StateTrackInputTags
	trackSession.LastBotMessageID = &message.MessageID

	slog.Info("Saved the link for track session", slog.String("operation", op), slog.String("url", URL), slog.Int64("chatID", chatID))

	h.sessionService.Set(chatID, trackSession)

	return nil
}
