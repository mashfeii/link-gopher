package telebot

import (
	errdefault "errors"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	scrapper_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BotClient struct {
	bot            *tgbotapi.BotAPI
	scrapperClient *scrapper_client.ClientWithResponses
	sessions       map[int64]*UserSession
	mu             sync.Mutex
}

func NewBotClient(cfg *config.Config) (*BotClient, error) {
	bot, err := tgbotapi.NewBotAPI(cfg.Secret.BotToken)
	if err != nil {
		return nil, fmt.Errorf("unable to create bot: %w", err)
	}

	bot.Debug = cfg.Serving.Debug

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	scrapperClient, err := scrapper_client.NewClientWithResponses(
		fmt.Sprintf("http://%s:%d", cfg.Serving.Host, cfg.Serving.ScrapperPort),
		scrapper_client.WithHTTPClient(client),
	)
	if err != nil {
		return nil, err
	}

	return &BotClient{
		bot:            bot,
		scrapperClient: scrapperClient,
		sessions:       make(map[int64]*UserSession),
		mu:             sync.Mutex{},
	}, nil
}

func (bot *BotClient) GetBot() *tgbotapi.BotAPI {
	return bot.bot
}

func (bot *BotClient) Run() {
	updates := bot.bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	})

	for update := range updates {
		if update.CallbackQuery != nil {
			slog.Info("received callback query", slog.Any("data", update.CallbackQuery.Data))

			if err := bot.handleCallback(&update); err != nil {
				slog.Error("unable to process the callback", slog.Any("error", err))
				bot.errorMessage(update.CallbackQuery.Message.Chat.ID)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		session := bot.getSession(chatID)

		if session != nil {
			slog.Info("received message", slog.Any("message", update.Message.Text))
			session.LastUserMessageID = &update.Message.MessageID
			bot.setSession(chatID, session)

			err := bot.handleSession(&update, session)
			if err != nil {
				if errdefault.As(err, &errors.ErrInvalidURL{}) {
					bot.handleInvalidURL(chatID, session)
					continue
				}

				slog.Error("unable to process the session", slog.Any("error", err))
				bot.errorMessage(chatID)
			}

			continue
		}

		if update.Message.IsCommand() {
			slog.Info("received command", slog.Any("command", update.Message.Command()))

			err := bot.handleCommand(&update)
			if err != nil {
				slog.Error("unable to process the command", slog.Any("error", err))
				bot.errorMessage(chatID)
			}

			continue
		}

		_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "😅 I don't understand you. Use /help to see available commands."))
	}
}

func (bot *BotClient) errorMessage(chatID int64) {
	_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "❌ Something went wrong. Please try again later."))
	bot.clearSession(chatID)
}
