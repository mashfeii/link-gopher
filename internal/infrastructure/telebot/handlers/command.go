package handlers

import (
	"context"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseCommandHandler struct {
	BotAPI      *tgbotapi.BotAPI
	CommandName string
}

func (h *BaseCommandHandler) Handle(ctx context.Context, update *tgbotapi.Update) error {
	//TODO implement me
	panic("implement me")
}

func (h *BaseCommandHandler) CanHandle(_ context.Context, update *tgbotapi.Update) bool {
	return update.Message != nil && update.Message.IsCommand() && update.Message.Command() == h.CommandName
}

func (h *BaseCommandHandler) GetCommandName() string {
	return h.CommandName
}
