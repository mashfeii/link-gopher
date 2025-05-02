package router

import (
	"context"
	"fmt"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/middleware"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Router struct {
	commandRegistry  *handlers.CommandsRegistry
	sessionManager   core.SessionManager
	sessionHandlers  []handlers.SessionHandler
	callbackHandlers []handlers.CallbackHandler
}

func NewRouter(
	botAPI *tgbotapi.BotAPI,
	sessionManager core.SessionManager,
) *Router {
	return &Router{
		commandRegistry: handlers.NewCommandsRegistry(botAPI),
		sessionManager:  sessionManager,
	}
}

func (r *Router) Init() error {
	if err := r.commandRegistry.RegisterWithTelegram(); err != nil {
		return fmt.Errorf("failed to register commands with Telegram: %w", err)
	}

	return nil
}

func (r *Router) RegisterCommand(handler handlers.CommandHandler, middlewares ...middleware.Wrapper) {
	wrapped := handler

	for _, m := range middlewares {
		wrapped = m(wrapped)
	}

	r.commandRegistry.Register(wrapped)
}

func (r *Router) RegisterSessionHandler(handler handlers.SessionHandler) {
	r.sessionHandlers = append(r.sessionHandlers, handler)
}

func (r *Router) RegisterCallbackHandler(handler handlers.CallbackHandler) {
	r.callbackHandlers = append(r.callbackHandlers, handler)
}

func (r *Router) Handle(ctx context.Context, update *tgbotapi.Update) error {
	if update.CallbackQuery != nil {
		return r.handleCallback(ctx, update)
	}

	if update.Message == nil {
		return nil
	}

	chatID := update.Message.Chat.ID

	if session := r.sessionManager.Get(chatID); session != nil {
		return r.handleSession(ctx, update, session)
	}

	if update.Message.IsCommand() {
		return r.handleCommand(ctx, update)
	}

	return nil
}

func (r *Router) handleCommand(ctx context.Context, update *tgbotapi.Update) error {
	const op = "router.handleCommand"

	command := update.Message.Command()
	if handler, ok := r.commandRegistry.GetHandler(command); ok {
		return handler.Handle(ctx, update)
	}

	return fmt.Errorf("%s: command handler not found for command: %s", op, command)
}

func (r *Router) handleCallback(ctx context.Context, update *tgbotapi.Update) error {
	const op = "router.handleCallback"

	for _, handler := range r.callbackHandlers {
		if handler.CanHandleCallback(update.CallbackQuery.Data) {
			return handler.HandleCallback(ctx, update)
		}
	}

	return fmt.Errorf("%s: callback handler not found for data: %s", op, update.CallbackQuery.Data)
}

func (r *Router) handleSession(ctx context.Context, update *tgbotapi.Update, session models.Session) error {
	const op = "router.handleSession"

	if update.Message.IsCommand() {
		return r.handleCommand(ctx, update)
	}

	for _, handler := range r.sessionHandlers {
		if handler.CanHandleState(session.GetState()) {
			return handler.HandleSession(ctx, update, session)
		}
	}

	return fmt.Errorf("%s: session handler not found for state: %s", op, session.GetState())
}
