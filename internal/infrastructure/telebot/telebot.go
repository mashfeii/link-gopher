package telebot

import (
	errdefault "errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"
	"time"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/service"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/es-debug/backend-academy-2024-go-template/config"
	scrapperclient "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/session"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type BotClient struct {
	bot            *tgbotapi.BotAPI
	scrapperClient *scrapperclient.ClientWithResponses
	apiService     core.ScrapperService
	sessionManager core.SessionManager
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

	scrapperClient, err := scrapperclient.NewClientWithResponses(
		fmt.Sprintf("http://%s:%d", cfg.Serving.Host, cfg.Serving.ScrapperPort),
		scrapperclient.WithHTTPClient(client),
	)
	if err != nil {
		return nil, err
	}

	return &BotClient{
		bot:            bot,
		scrapperClient: scrapperClient,
		apiService:     service.NewScrapperClientWrapper(scrapperClient),
		sessionManager: session.NewSessionManager(),
	}, nil
}

func (bot *BotClient) GetBot() *tgbotapi.BotAPI {
	return bot.bot
}

func (bot *BotClient) registerCommands() error {
	commands := []tgbotapi.BotCommand{
		{Command: "start", Description: " user registration"},
		{Command: "help", Description: " displays a list of available commands"},
		{Command: "track", Description: "start tracking the link"},
		{Command: "untrack", Description: "stop tracking the link"},
		{Command: "list", Description: "list all tracked links"},
		{Command: "cancel", Description: "cancel current operation"},
	}

	commandsConfig := tgbotapi.NewSetMyCommands(commands...)
	_, err := bot.bot.Request(commandsConfig)

	return err
}

func (bot *BotClient) Run() {
	if err := bot.registerCommands(); err != nil {
		slog.Error("unable to generate commands", slog.Any("error", err))
		return
	}

	updates := bot.bot.GetUpdatesChan(tgbotapi.UpdateConfig{
		Offset:  0,
		Timeout: 60,
	})

	for update := range updates {
		if update.CallbackQuery != nil {
			slog.Info("received callback query", slog.Any("data", update.CallbackQuery.Data))

			if err := bot.handleCallback(&update); err != nil {
				slog.Error("unable to process the callback", slog.Any("error", err))
				bot.errorMessage(update.CallbackQuery.Message.Chat.ID, err)
			}

			continue
		}

		if update.Message == nil {
			continue
		}

		chatID := update.Message.Chat.ID
		userSession := bot.sessionManager.Get(chatID)

		if userSession != nil {
			slog.Info("received message", slog.Any("message", update.Message.Text))
			userSession.LastUserMessageID = &update.Message.MessageID
			bot.sessionManager.Set(chatID, userSession)

			err := bot.handleSession(&update, userSession)
			if err != nil {
				if errdefault.As(err, &errors.ErrInvalidURL{}) {
					bot.handleInvalidURL(chatID, userSession)
					continue
				}

				slog.Error("unable to process the user session", slog.Any("error", err))
				bot.errorMessage(chatID, err)
			}

			continue
		}

		if update.Message.IsCommand() {
			slog.Info("received command", slog.Any("command", update.Message.Command()))

			err := bot.handleCommand(&update)
			if err != nil {
				slog.Error("unable to process the command", slog.Any("error", err))
				bot.errorMessage(chatID, err)
			}

			continue
		}

		_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "😅 I don't understand you. Use /help to see available commands."))
	}
}

func (bot *BotClient) errorMessage(chatID int64, err error) {
	message := tgbotapi.NewMessage(chatID, "")

	switch err.(type) {
	case errors.ErrUserNotFound:
		message.Text = "😔 You are not registered yet, use /start to register"
	default:
		message.Text = "❌ Something went wrong. Please try again later."
	}

	_, _ = bot.bot.Send(message)
	bot.sessionManager.Clear(chatID)
}

func (bot *BotClient) updateTagKeyboard(chatID int64, session *models.UserSession) error {
	newKeyboard := ui.BuildTagsKeyboard(session.AvailableTags, session.SelectedTags)

	edit := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		*session.LastMessageID,
		newKeyboard,
	)

	if _, err := bot.bot.Send(edit); err != nil {
		return err
	}

	return nil
}

func (bot *BotClient) sendFilterKeyKeyboard(chatID int64, session *models.UserSession) error {
	availableFilters := []string{}

	for _, filter := range session.AvailableFilters {
		filterName := strings.Split(filter, ":")[0]
		if !slices.Contains(availableFilters, filterName) {
			availableFilters = append(availableFilters, filterName)
		}
	}

	newKeyboard := ui.BuildFiltersKeyboard(availableFilters)

	if session.LastMessageID == nil {
		message := tgbotapi.NewMessage(chatID, "📋 Please, select filter name:")
		message.ReplyMarkup = newKeyboard
		sentMsg, err := bot.bot.Send(message)
		session.LastMessageID = &sentMsg.MessageID

		return err
	}

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		*session.LastMessageID,
		"📋 Please, select filter name:",
		newKeyboard,
	)

	if _, err := bot.bot.Send(edit); err != nil {
		return err
	}

	return nil
}

func (bot *BotClient) sendFilterValueKeyboard(chatID int64, session *models.UserSession) error {
	availableFilters := []string{}

	for _, filter := range session.AvailableFilters {
		keyAndValue := strings.Split(filter, ":")
		if keyAndValue[0] == session.CurrentFilterName &&
			!slices.Contains(availableFilters, keyAndValue[1]) {
			availableFilters = append(availableFilters, keyAndValue[1])
		}
	}

	newKeyboard := ui.BuildFiltersValueKeyboard(availableFilters)

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		*session.LastMessageID,
		"📋 Please, select filter value:",
		newKeyboard,
	)

	if _, err := bot.bot.Send(edit); err != nil {
		return err
	}

	return nil
}
