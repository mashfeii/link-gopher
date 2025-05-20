package router

import (
	"context"
	"fmt"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/middleware"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler struct {
	Name        string
	Description string
	Handler     handlers.Handler
}

type CallbackHanlder struct {
	Prefix  string
	Handler handlers.Handler
}

type Router struct {
	commandHandlers  map[models.CommandName]CommandHandler
	sessionHandlers  map[models.State]handlers.Handler
	callbackHandlers []CallbackHanlder
	sessionManager   core.SessionManager
	messageService   core.MessageService
	scrapperService  core.ScrapperService
	botAPI           *tgbotapi.BotAPI
}

func NewRouter(
	botAPI *tgbotapi.BotAPI,
	sessionManager core.SessionManager,
	messageService core.MessageService,
	scrapperService core.ScrapperService,
) *Router {
	return &Router{
		commandHandlers: make(map[models.CommandName]CommandHandler),
		sessionHandlers: make(map[models.State]handlers.Handler),
		sessionManager:  sessionManager,
		messageService:  messageService,
		scrapperService: scrapperService,
		botAPI:          botAPI,
	}
}

func (r *Router) RegisterWithTelegram() error {
	commands := make([]tgbotapi.BotCommand, 0, len(r.commandHandlers))

	for _, handler := range r.commandHandlers {
		commands = append(commands, tgbotapi.BotCommand{
			Command:     handler.Name,
			Description: handler.Description,
		})
	}

	if _, err := r.botAPI.Request(tgbotapi.NewSetMyCommands(commands...)); err != nil {
		return fmt.Errorf("failed to register commands with Telegram: %w", err)
	}

	return nil
}

func (r *Router) RegisterCommand(
	name models.CommandName,
	description string,
	handler handlers.Handler,
	middlewares ...middleware.Middleware,
) {
	wrapped := handler

	for _, m := range middlewares {
		wrapped = m(wrapped)
	}

	r.commandHandlers[name] = CommandHandler{
		Name:        name.String(),
		Description: description,
		Handler:     wrapped,
	}
}

func (r *Router) RegisterSessionHandler(state models.State, handler handlers.Handler, middlewares ...middleware.Middleware) {
	wrapped := handler

	for _, m := range middlewares {
		wrapped = m(wrapped)
	}

	r.sessionHandlers[state] = wrapped
}

func (r *Router) RegisterCallbackHandler(prefix string, handler handlers.Handler, middlewares ...middleware.Middleware) {
	wrapped := handler

	for _, m := range middlewares {
		wrapped = m(wrapped)
	}

	r.callbackHandlers = append(r.callbackHandlers, CallbackHanlder{
		Prefix:  prefix,
		Handler: wrapped,
	})
}

func (r *Router) Handle(ctx context.Context, update *tgbotapi.Update, chatID int64) error {
	hctx := &handlers.HandlerContext{
		ChatID:          chatID,
		Update:          update,
		Session:         r.sessionManager.Get(chatID),
		MessageService:  r.messageService,
		SessionService:  r.sessionManager,
		ScrapperService: r.scrapperService,
		BotAPI:          r.botAPI,
	}

	if update.CallbackQuery != nil {
		return r.handleCallback(ctx, hctx)
	}

	if update.Message != nil {
		if update.Message.IsCommand() {
			return r.handleCommand(ctx, hctx)
		}

		if hctx.Session != nil {
			return r.handleSession(ctx, hctx)
		}
	}

	return nil
}

func (r *Router) handleCommand(ctx context.Context, hctx *handlers.HandlerContext) error {
	const op = "router.handleCommand"

	command := models.CommandName(hctx.Update.Message.Command())
	if commandInfo, ok := r.commandHandlers[command]; ok {
		if err := commandInfo.Handler.Handle(ctx, hctx); err != nil {
			r.sessionManager.Clear(hctx.Update.Message.Chat.ID)

			return err
		}

		return nil
	}

	return fmt.Errorf("%s:%w", op, errors.NewErrErrUnknownHandler(string(command)))
}

func (r *Router) handleCallback(ctx context.Context, hctx *handlers.HandlerContext) error {
	const op = "router.handleCallback"

	callbackName := hctx.Update.CallbackQuery.Data

	for _, handler := range r.callbackHandlers {
		if strings.HasPrefix(callbackName, handler.Prefix) {
			if err := handler.Handler.Handle(ctx, hctx); err != nil {
				_, _ = r.messageService.EditMarkup(hctx.ChatID, *hctx.Session.GetLastBotMessageID(), nil)
				r.sessionManager.Clear(hctx.ChatID)

				return fmt.Errorf("%s: error handling callback: %w", op, err)
			}

			return nil
		}
	}

	r.sessionManager.Clear(hctx.Update.CallbackQuery.Message.Chat.ID)

	return fmt.Errorf("%s: callback handler not found for data: %s", op, hctx.Update.CallbackQuery.Data)
}

func (r *Router) handleSession(ctx context.Context, hctx *handlers.HandlerContext) error {
	const op = "router.handleSession"

	if hctx.Update.Message.IsCommand() {
		return r.handleCommand(ctx, hctx)
	}

	if handler, ok := r.sessionHandlers[hctx.Session.GetState()]; ok {
		if err := handler.Handle(ctx, hctx); err != nil {
			r.sessionManager.Clear(hctx.Update.Message.Chat.ID)

			return fmt.Errorf("%s: error handling session: %w", op, err)
		}

		return nil
	}

	r.sessionManager.Clear(hctx.Update.Message.Chat.ID)

	return fmt.Errorf("%s: session handler not found for state: %s", op, hctx.Session.GetState())
}
