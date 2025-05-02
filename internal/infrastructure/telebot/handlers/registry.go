package handlers

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type CommandsRegistry struct {
	botAPI   *tgbotapi.BotAPI
	commands map[string]CommandHandler
}

func NewCommandsRegistry(botAPI *tgbotapi.BotAPI) *CommandsRegistry {
	return &CommandsRegistry{
		botAPI:   botAPI,
		commands: make(map[string]CommandHandler),
	}
}

func (r *CommandsRegistry) Register(handler CommandHandler) {
	r.commands[handler.GetName()] = handler
}

func (r *CommandsRegistry) RegisterWithTelegram() error {
	commands := make([]tgbotapi.BotCommand, 0, len(r.commands))

	for name, handler := range r.commands {
		commands = append(commands, tgbotapi.BotCommand{
			Command:     name,
			Description: handler.GetDescription(),
		})
	}

	_, err := r.botAPI.Request(tgbotapi.NewSetMyCommands(commands...))

	return err
}

func (r *CommandsRegistry) GetHandler(name string) (CommandHandler, bool) {
	handler, ok := r.commands[name]
	return handler, ok
}

func (r *CommandsRegistry) GetHandlers() []CommandHandler {
	commands := make([]CommandHandler, 0, len(r.commands))
	for _, handler := range r.commands {
		commands = append(commands, handler)
	}
	return commands
}
