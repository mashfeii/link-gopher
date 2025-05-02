package commands

import (
	"context"
	"errors"
	"fmt"

	clienterrors "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type StartHandler struct {
	*BaseCommandHandler
	scrapperService core.ScrapperService
	messageService  core.MessageService
}

func NewStartHandler(api *tgbotapi.BotAPI, scrapperService core.ScrapperService, messageService core.MessageService) *StartHandler {
	return &StartHandler{
		BaseCommandHandler: NewBaseCommandHandler(api, core.CommandNameStart),
		scrapperService:    scrapperService,
		messageService:     messageService,
	}
}

func (h *StartHandler) Handle(ctx context.Context, update *tgbotapi.Update) error {
	const op = "commands.StartHandler.Handle"

	chatID := update.Message.Chat.ID
	serviceContext := context.WithValue(ctx, core.ContextKeyChatID, chatID)

	err := h.scrapperService.RegisterUser(serviceContext)
	if err != nil {
		if errors.Is(err, clienterrors.ErrUserAlreadyExists{}) {
			_, err := h.messageService.Send(chatID, "😅 You are already registered! Feel free to use the bot.", nil)
			return err
		}

		return fmt.Errorf("%s: unable to register user: %w", op, err)
	}

	_, err = h.messageService.Send(chatID, "🥳 Successfully registered!", nil)

	return err
}

func (h *StartHandler) GetDescription() string {
	return "user registration"
}
