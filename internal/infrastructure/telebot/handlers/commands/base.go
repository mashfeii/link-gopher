package commands

import (
	"context"
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

// General implementation of a command handler interface.
type BaseCommandHandler struct {
	botAPI *tgbotapi.BotAPI
	name   string
}

func NewBaseCommandHandler(api *tgbotapi.BotAPI, name string) *BaseCommandHandler {
	return &BaseCommandHandler{
		botAPI: api,
		name:   name,
	}
}

func (h *BaseCommandHandler) CanHandle(_ context.Context, update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.Command() == h.name
}

func (h *BaseCommandHandler) Handle(_ context.Context, _ *tgbotapi.Update) error {
	const op = "commands.BaseCommandHandler.Handle"
	return fmt.Errorf("%s: not implemented", op)
}

func (h *BaseCommandHandler) GetName() string {
	return h.name
}
