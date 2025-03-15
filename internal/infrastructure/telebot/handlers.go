package telebot

import (
	"context"
	"fmt"
	"net/http"

	scrapperclient "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *BotClient) handleCommand(update *tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	switch update.Message.Command() {
	case "track":
		bot.createTrackSession(chatID, StateWaitingURL)
	case "untrack", "list":
		return bot.createListUntrackSession(chatID, update.Message.Command())
	case "start":
		return bot.handleStart(update)
	case "help":
		return bot.handleHelp(update.Message.Chat.ID)
	case "cancel":
		return bot.handleCancel(update.Message.Chat.ID)
	default:
		return bot.handleUnknown(update.Message.Chat.ID)
	}

	return nil
}

func (bot *BotClient) handleStart(update *tgbotapi.Update) error {
	chatID := update.Message.Chat.ID

	resp, err := bot.scrapperClient.PostTgChatId(context.TODO(), chatID)
	if err != nil {
		return fmt.Errorf("unable to register user: %w", err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode == http.StatusAlreadyReported {
		message := tgbotapi.NewMessage(chatID, "😅 You are already registered! Feel free to use the bot.")
		_, _ = bot.bot.Send(message)

		return bot.handleHelp(chatID)
	}

	message := tgbotapi.NewMessage(chatID, "🥳 Successfully registered!")

	if _, err := bot.bot.Send(message); err != nil {
		return err
	}

	return nil
}

func (bot *BotClient) handleHelp(chatID int64) error {
	availableCommands, err := bot.bot.GetMyCommands()
	if err != nil {
		return err
	}

	message := tgbotapi.NewMessage(chatID, "🤖 Available commands:\n")
	for _, command := range availableCommands {
		message.Text += fmt.Sprintf("/%s: %s\n", command.Command, command.Description)
	}

	if _, err := bot.bot.Send(message); err != nil {
		return err
	}

	return nil
}

func (bot *BotClient) handleUnknown(chatID int64) error {
	message := tgbotapi.NewMessage(chatID, "❌ Unknown command, use /help to get the list of commands")

	if _, err := bot.bot.Send(message); err != nil {
		return err
	}

	return nil
}

func (bot *BotClient) saveTracking(chatID int64, session *UserSession, saveType string) error {
	_, err := bot.scrapperClient.PostLinksWithResponse(context.TODO(), &scrapperclient.PostLinksParams{
		TgChatId: chatID,
	}, scrapperclient.PostLinksJSONRequestBody{
		Filters: &session.Filters,
		Tags:    &session.Tags,
		Link:    &session.URL,
	})
	if err != nil {
		return fmt.Errorf("unable to save tracking: %w", err)
	}

	return bot.handleList(chatID, session, saveType)
}

func (bot *BotClient) handleList(chatID int64, session *UserSession, saveType string) error {
	resp, err := bot.scrapperClient.GetLinksWithResponse(context.TODO(), &scrapperclient.GetLinksParams{
		TgChatId: chatID,
	})
	if err != nil {
		return fmt.Errorf("unable to get links: %w", err)
	}

	switch saveType {
	case "edit":
		edit := tgbotapi.NewEditMessageTextAndMarkup(
			chatID,
			*session.LastMessageID,
			"🚀 Tracking saved!",
			buildLinksKeyBoard(*resp.JSON200.Links, "list"),
		)

		if _, err := bot.bot.Send(edit); err != nil {
			return fmt.Errorf("unable to send message: %w", err)
		}
	default:
		message := tgbotapi.NewMessage(chatID, "🚀 Tracking saved!")
		message.ReplyMarkup = buildLinksKeyBoard(*resp.JSON200.Links, "list")

		if _, err := bot.bot.Send(message); err != nil {
			return fmt.Errorf("unable to send message: %w", err)
		}
	}

	return nil
}

func (bot *BotClient) handleCancel(chatID int64) error {
	if session := bot.getSession(chatID); session == nil {
		_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "You don't have any active sessions"))
		return nil
	}

	bot.clearSession(chatID)

	if _, err := bot.bot.Send(tgbotapi.NewMessage(chatID, "🚮 Operation canceled")); err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

func bindKeyboardMessage(chatID int64, messageText string, data [][]string) tgbotapi.MessageConfig {
	message := tgbotapi.NewMessage(chatID, messageText)
	buttonsRow := []tgbotapi.InlineKeyboardButton{}

	for _, row := range data {
		buttonsRow = append(buttonsRow, tgbotapi.InlineKeyboardButton{
			Text:         row[0],
			CallbackData: &row[1],
		})
	}

	message.ReplyMarkup = tgbotapi.NewInlineKeyboardMarkup(buttonsRow)

	return message
}
