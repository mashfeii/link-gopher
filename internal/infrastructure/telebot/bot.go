package telebot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers/commands"
	sessionHandlers "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers/session"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/messaging"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/middleware"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/router"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	scrapperclient "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/session"
)

type BotClient struct {
	bot            *tgbotapi.BotAPI
	router         *router.Router
	apiService     core.ScrapperService
	sessionManager core.SessionManager
	messageService core.MessageService
}

func (bot *BotClient) GetBotAPI() *tgbotapi.BotAPI {
	return bot.bot
}

func NewBotClient(cfg *config.Config) (*BotClient, error) {
	const op = "telebot.NewBotClient"

	// Initialize the Telegram bot API
	apiWrapper, err := tgbotapi.NewBotAPI(cfg.Secret.BotToken)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	apiWrapper.Debug = cfg.Serving.Debug

	slog.Info("authorized on account",
		slog.String("operation", op),
		slog.String("bot_name", apiWrapper.Self.UserName),
	)

	// Set up the HTTP client with a timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}
	// Create a new scrapper client with the HTTP client
	scrapperClient, err := scrapperclient.NewClientWithResponses(
		fmt.Sprintf("http://%s:%d", cfg.Serving.Host, cfg.Serving.ScrapperPort),
		scrapperclient.WithHTTPClient(client),
	)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Create bot dependencies
	sessionManager := session.NewSessionManager()
	apiService := service.NewScrapperClientWrapper(scrapperClient)
	messageService := messaging.NewMessageService(apiWrapper)

	slog.Info("creating bot client",
		slog.String("operation", op),
		slog.String("scrapper_client", fmt.Sprintf("http://%s:%d", cfg.Serving.Host, cfg.Serving.ScrapperPort)),
		slog.String("scrapper_client_timeout", fmt.Sprintf("%f", client.Timeout.Seconds())),
	)

	bot := &BotClient{
		bot:            apiWrapper,
		apiService:     apiService,
		sessionManager: sessionManager,
		messageService: messageService,
	}

	// Register defined handlers for commands
	baseRouter := router.NewRouter(apiWrapper, sessionManager)
	withAuth := middleware.WithAuth(apiService, messageService)
	baseRouter.RegisterCommand(commands.NewStartHandler(apiWrapper, apiService, messageService))
	baseRouter.RegisterCommand(commands.NewHelpHandler(apiWrapper, messageService))
	baseRouter.RegisterCommand(commands.NewCancelHandler(apiWrapper, sessionManager, messageService))

	// Register defined protected handlers for commands
	baseRouter.RegisterCommand(
		commands.NewTrackHandler(apiWrapper, apiService, messageService, sessionManager),
		withAuth,
	)

	// Register defined handlers for sessions
	baseRouter.RegisterSessionHandler(sessionHandlers.NewTrackInputURL(messageService, sessionManager, apiService))
	baseRouter.RegisterSessionHandler(sessionHandlers.NewTrackInputTags(messageService, sessionManager))

	bot.router = baseRouter

	return bot, nil
}

func (bot *BotClient) Run() {
	const op = "telebot.BotClient.Run"

	slog.Info("starting bot client",
		slog.String("operation", op),
		slog.String("bot_name", bot.bot.Self.UserName),
		slog.String("bot_id", fmt.Sprintf("%d", bot.bot.Self.ID)),
	)

	if err := bot.router.Init(); err != nil {
		slog.Error("failed to initialize router", slog.String("operation", op), slog.Any("error", err))
		return
	}

	updates := bot.bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	})

	ctx := context.Background()

	for update := range updates {
		slog.Info("received update",
			slog.String("operation", op),
			slog.String("chat_id", fmt.Sprintf("%d", update.Message.Chat.ID)),
			slog.Bool("callback_query", update.CallbackQuery != nil),
			slog.Bool("command", update.Message.IsCommand()),
			slog.Bool("message", update.Message != nil),
		)

		if err := bot.router.Handle(ctx, &update); err != nil {
			if _, err := bot.messageService.SendError(update.Message.Chat.ID, err); err != nil {
				slog.Error("failed to send error message", slog.String("operation", op), slog.Any("error", err))
			}
		}

		slog.Info("update handled")
	}
}
