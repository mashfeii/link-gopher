package telebot

import (
	"fmt"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
	"log/slog"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func (bot *BotClient) createTypedSession(sessionType models.SessionType, command string, filters, tags []string) *models.UserSession {
	return &models.UserSession{
		Type:             sessionType,
		State:            models.StateWaitingTagsSelection,
		Command:          command,
		AvailableFilters: filters,
		AvailableTags:    tags,
	}
}

func (bot *BotClient) createTrackSession(chatID int64, initialState models.State) {
	if _, err := bot.handleAuthorizationCheck(chatID); err != nil {
		bot.errorMessage(chatID, err)
		return
	}

	session := &models.UserSession{
		Type:  models.SessionTypeTrack,
		State: initialState,
	}
	message, _ := bot.bot.Send(tgbotapi.NewMessage(chatID, "🔗 Please enter the URL you want to track"))
	session.LastMessageID = &message.MessageID
	bot.sessionManager.Set(chatID, session)
}

func (bot *BotClient) createListUntrackSession(chatID int64, command string) error {
	slog.Info("Bot client: createListUntrackSession", slog.Any("chatID", chatID), slog.Any("command", command))

	resp, err := bot.handleAuthorizationCheck(chatID)
	if err != nil {
		return err
	}

	slog.Info("Bot Client: createListUntrackSession, authorized")

	links, err := validateLinksResponse(chatID, bot, resp)
	if err != nil || links == nil {
		return err
	}

	slog.Info("Bot Client: createListUntrackSession", slog.Any("links", links))

	availableFilters, availableTags := bot.extractFiltersAndTags(*links)

	slog.Info("Bot Client: createListUntrackSession",
		slog.Any("availableFilters", availableFilters),
		slog.Any("availableTags", availableTags),
	)

	session := bot.createTypedSession(models.SessionTypeListUntrack, command, availableFilters, availableTags)
	bot.sessionManager.Set(chatID, session)

	return bot.handleTypedSession(chatID, session)
}

func (bot *BotClient) handleSession(update *tgbotapi.Update, session *models.UserSession) error {
	chatID := update.Message.Chat.ID
	text := strings.TrimPrefix(update.Message.Text, "https://")

	if update.Message.IsCommand() {
		if update.Message.Command() != "cancel" {
			bot.sessionManager.Clear(chatID)
			_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "Okey, let's go with another command 🤖"))
		}

		return bot.handleCommand(update)
	}

	switch session.State {
	case models.StateWaitingURL:
		ok, err := bot.validateLink(chatID, text)
		if err != nil {
			return err
		}

		if !ok {
			message, _ := bot.bot.Send(tgbotapi.NewMessage(chatID,
				"❌ This link is already being tracked. Please enter another one or use /untrack to remove it."))
			session.LastMessageID = &message.MessageID
			bot.sessionManager.Set(chatID, session)

			return nil
		}

		session.URL = text
		session.State = models.StateWaitingTags

		message := bindKeyboardMessage(chatID, "📋 Enter tags (space-separated):",
			[][]string{{"◀️ Step back", "return_url"}, {"🚫 Skip", "skip_tags"}})

		messageSend, err := bot.bot.Send(message)
		if err != nil {
			return fmt.Errorf("unable to send message: %w", err)
		}

		session.LastMessageID = &messageSend.MessageID
		bot.sessionManager.Set(chatID, session)
	case models.StateWaitingTags:
		session.Tags = strings.Fields(text)
		session.State = models.StateWaitingFilters

		message := bindKeyboardMessage(chatID, "📋 Enter filters (space-separated key:value pairs):",
			[][]string{{"◀️ Step back", returnTags}, {"🚫 Skip", "skip_filters"}})

		messageSend, err := bot.bot.Send(message)
		if err != nil {
			return err
		}

		session.LastMessageID = &messageSend.MessageID
		bot.sessionManager.Set(chatID, session)
	case models.StateWaitingFilters:
		session.Filters = strings.Fields(text)
		if err := bot.saveTracking(chatID, session, "new"); err != nil {
			return fmt.Errorf("unable to save tracking: %w", err)
		}

		bot.sessionManager.Clear(chatID)
	default:
		request := tgbotapi.NewDeleteMessage(chatID, *session.LastUserMessageID)
		if _, err := bot.bot.Request(request); err != nil {
			return fmt.Errorf("unable to delete message: %w", err)
		}
	}

	return nil
}

func (bot *BotClient) handleTypedSession(chatID int64, session *models.UserSession) error {
	switch {
	case len(session.AvailableTags) == 0 && len(session.AvailableFilters) == 0:
		return bot.finalizeListUntrackSession(chatID, session)

	case len(session.AvailableTags) == 0:
		session.State = models.StateWaitingFiltersSelection
		if err := bot.sendFilterKeyKeyboard(chatID, session); err != nil {
			return fmt.Errorf("unable to send filter key keyboard: %w", err)
		}

	default:
		msg := tgbotapi.NewMessage(chatID, "🏷 Select tags to filter:")
		msg.ReplyMarkup = ui.BuildTagsKeyboard(session.AvailableTags, session.SelectedTags)

		sentMsg, err := bot.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		session.LastMessageID = &sentMsg.MessageID
	}

  bot.sessionManager.Set(chatID, session)

	return nil
}

func (bot *BotClient) handleInvalidURL(chatID int64, session *models.UserSession) {
	message := tgbotapi.NewMessage(
		chatID,
		`❌ Invalid URL. Try another one:
  ✅ github.com/golang/go
  ✅ stackoverflow.com/questions/17333517/how-to-compile-a-program-in-go-language`,
	)

	if _, err := bot.bot.Send(message); err != nil {
		slog.Error("unable to send message", slog.Any("error", err))
	}

	session.State = models.StateWaitingURL
	bot.sessionManager.Set(chatID, session)
}
