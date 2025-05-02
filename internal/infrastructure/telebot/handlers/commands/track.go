package commands

import (
	"context"
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/session"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TrackHandler struct {
	*BaseCommandHandler
	scrapperService core.ScrapperService
	messageService  core.MessageService
	sessionService  core.SessionManager
}

func NewTrackHandler(
	api *tgbotapi.BotAPI,
	scrapperService core.ScrapperService,
	messageService core.MessageService,
	sessionService core.SessionManager,
) *TrackHandler {
	return &TrackHandler{
		BaseCommandHandler: NewBaseCommandHandler(api, core.CommandNameTrack),
		scrapperService:    scrapperService,
		messageService:     messageService,
		sessionService:     sessionService,
	}
}

func (h *TrackHandler) Handle(ctx context.Context, update *tgbotapi.Update) error {
	const op = "commands.TrackHandler.Handle"

	// Translate id to service
	chatID := update.Message.Chat.ID
	serviceContext := context.WithValue(ctx, core.ContextKeyChatID, chatID)

	// Check if the user is already registered
	if _, err := h.scrapperService.GetUser(serviceContext); err != nil {
		return fmt.Errorf("%s: unable to get user: %w", op, err)
	}

	// First action is to ask for the URL
	message, _ := h.messageService.Send(chatID, "🔗 Please enter the URL you want to track", nil)

	// Create typed session
	h.sessionService.Set(chatID, session.CreateTrackSession(message.MessageID, update.Message.MessageID))

	return nil
}

func (h *TrackHandler) GetDescription() string {
	return "start tracking a link for changes"
}
