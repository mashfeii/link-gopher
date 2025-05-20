package commands

import (
	"context"
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
)

type HelpHandler struct{}

func (h *HelpHandler) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	commands, err := hctx.BotAPI.GetMyCommands()
	if err != nil {
		return fmt.Errorf("failed to get bot commands: %w", err)
	}

	var helpMessage string
	for _, command := range commands {
		helpMessage += fmt.Sprintf("/%s - %s\n", command.Command, command.Description)
	}

	_, err = hctx.MessageService.Send(hctx.ChatID, helpMessage, nil)
	if err != nil {
		return fmt.Errorf("failed to send help message: %w", err)
	}

	return nil
}
