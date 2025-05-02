package commands

import (
	"context"
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CancelHandler struct {
	*BaseCommandHandler
	sessionService core.SessionManager
	messageService core.MessageService
}

func NewCancelHandler(
	api *tgbotapi.BotAPI,
	sessionService core.SessionManager,
	messageService core.MessageService,
) *CancelHandler {
	return &CancelHandler{
		BaseCommandHandler: NewBaseCommandHandler(api, core.CommandNameCancel),
		sessionService:     sessionService,
		messageService:     messageService,
	}
}

func (h *CancelHandler) Handle(_ context.Context, update *tgbotapi.Update) error {
	const op = "commands.CancelHandler.Handle"

	chatID := update.Message.Chat.ID

	if session := h.sessionService.Get(chatID); session == nil {
		return fmt.Errorf("%s: %w", op, errors.NewErrNoActiveSession(chatID))
	}

	h.sessionService.Clear(chatID)

	_, err := h.messageService.Send(chatID, "🚮 Operation canceled", nil)

	return err
}

func (h *CancelHandler) GetDescription() string {
	return "cancel the current operation (session)"
}
