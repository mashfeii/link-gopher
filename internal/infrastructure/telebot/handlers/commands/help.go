package commands

import (
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type HelpHandler struct {
	*BaseCommandHandler
	messageService core.MessageService
}

func NewHelpHandler(api *tgbotapi.BotAPI, messageService core.MessageService) *HelpHandler {
	return &HelpHandler{
		BaseCommandHandler: NewBaseCommandHandler(api, core.CommandNameHelp),
		messageService:     messageService,
	}
}

func (help *HelpHandler) GetDescription() string {
	return "displays a list of available commands"
}
