package telebot

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	scrapper_client "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	listCommand    = "list"
	untrackCommand = "untrack"
	skipFilters    = "skip_filters"
	returnTags     = "return_tags"
)

func (bot *BotClient) handleCallback(update *tgbotapi.Update) error {
	chatID := update.CallbackQuery.Message.Chat.ID
	data := update.CallbackQuery.Data
	session := bot.getSession(chatID)

	callback := tgbotapi.NewCallback(update.CallbackQuery.ID, "⏳ Processing...")
	if _, err := bot.bot.Request(callback); err != nil {
		slog.Error("callback failed", slog.Any("error", err))
	}

	if session == nil {
		return errors.NewErrNoActiveSession(chatID)
	}

	switch session.Type {
	case SessionTypeTrack:
		return bot.handleSessionCallback(chatID, data, session)
	case SessionTypeListUntrack:
		return bot.handleListUntrackCallback(chatID, data, session)
	default:
		return errors.NewErrUnknownSessionType(chatID)
	}
}

func (bot *BotClient) handleSessionCallback(chatID int64, data string, session *UserSession) error {
	switch data {
	case "return_url":
		session.State = StateWaitingURL
		session.Filters = nil
		bot.setSession(chatID, session)

		edit := tgbotapi.NewEditMessageText(
			chatID,
			*session.LastMessageID,
			"🔗 Please enter the URL you want to track",
		)

		if _, err := bot.bot.Send(edit); err != nil {
			return err
		}

	case returnTags:
		session.State = StateWaitingTags
		session.Filters = nil
		bot.setSession(chatID, session)

		skipButton := tgbotapi.NewInlineKeyboardButtonData("🚫 Skip", "skip_tags")
		returnButton := tgbotapi.NewInlineKeyboardButtonData("◀️ Step back", "return_url")

		edit := tgbotapi.NewEditMessageTextAndMarkup(
			chatID,
			*session.LastMessageID,
			"📋 Enter tags (space-separated):",
			tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{returnButton, skipButton}),
		)

		if _, err := bot.bot.Send(edit); err != nil {
			return err
		}

	case "skip_tags":
		session.State = StateWaitingFilters

		skipBtn := tgbotapi.NewInlineKeyboardButtonData("🚫 Skip", "skip_filters")
		returnButton := tgbotapi.NewInlineKeyboardButtonData("◀️ Step back", returnTags)

		edit := tgbotapi.NewEditMessageTextAndMarkup(
			chatID,
			*session.LastMessageID,
			"📋 Enter filters (space-separated key:value pairs):",
			tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{returnButton, skipBtn}),
		)

		bot.setSession(chatID, session)

		if _, err := bot.bot.Send(edit); err != nil {
			return err
		}
	case skipFilters:
		if err := bot.saveTracking(chatID, session, "edit"); err != nil {
			_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "Failed to save tracking."))
			return err
		}

		bot.clearSession(chatID)
	}

	return nil
}

func (bot *BotClient) handleListUntrackCallback(chatID int64, data string, session *UserSession) error {
	switch {
	case strings.HasPrefix(data, "tag_toggle "):
		tag := strings.TrimPrefix(data, "tag_toggle ")
		session.toggleTag(tag)

		if err := bot.updateTagKeyboard(chatID, session); err != nil {
			return err
		}

	case data == "done_tags":
		if len(session.AvailableFilters) == 0 {
			return bot.finalizeListUntrackSession(chatID, session)
		}

		session.State = StateWaitingFiltersSelection
		if err := bot.sendFilterKeyKeyboard(chatID, session); err != nil {
			return err
		}
	case strings.HasPrefix(data, "select_key "):
		key := strings.TrimPrefix(data, "select_key ")
		session.CurrentFilterName = key

		if err := bot.sendFilterValueKeyboard(chatID, session); err != nil {
			return err
		}

	case strings.HasPrefix(data, "select_value "):
		value := strings.TrimPrefix(data, "select_value ")
		session.CurrentFilterValue = value

		return bot.finalizeListUntrackSession(chatID, session)

	case data == "return_filters":
		session.CurrentFilterName = ""
		session.State = StateWaitingFiltersSelection

		if err := bot.sendFilterKeyKeyboard(chatID, session); err != nil {
			return err
		}

	case data == returnTags:
		session.State = StateWaitingTagsSelection

		if err := bot.updateTagKeyboard(chatID, session); err != nil {
			return err
		}
	case data == skipFilters:
		return bot.finalizeListUntrackSession(chatID, session)
	case strings.HasPrefix(data, "untrack "):
		return bot.handleUntrackCallback(chatID, data)
	}

	return nil
}

func (bot *BotClient) handleUntrackCallback(chatID int64, data string) error {
	URL := strings.TrimPrefix(data, "untrack ")

	resp, err := bot.scrapperClient.DeleteLinksWithResponse(context.TODO(),
		&scrapper_client.DeleteLinksParams{TgChatId: chatID},
		scrapper_client.DeleteLinksJSONRequestBody{Link: &URL},
	)
	if err != nil {
		return fmt.Errorf("unable to untrack link: %w", err)
	}

	if resp.JSON400 != nil {
		return fmt.Errorf("invalid request: %s", *resp.JSON400.Description)
	}

	if resp.StatusCode() == http.StatusOK {
		_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "✅ Link untracked successfully."))

		bot.clearSession(chatID)

		return nil
	}

	return nil
}

func (bot *BotClient) finalizeListUntrackSession(chatID int64, session *UserSession) error {
	linksResponse, err := bot.scrapperClient.GetLinksWithResponse(context.TODO(),
		&scrapper_client.GetLinksParams{TgChatId: chatID},
	)
	if err != nil {
		return fmt.Errorf("unable to finilize session: %w", err)
	}

	if err := handleErrorResponse(linksResponse); err != nil {
		return err
	}

	links, err := validateLinksResponse(chatID, bot, linksResponse)
	if err != nil || links == nil {
		return err
	}

	filter := session.CurrentFilterName + ":" + session.CurrentFilterValue
	filterValid := regexp.MustCompile(`^\w+:\w+$`).MatchString(filter)
	hasTags := len(session.SelectedTags) > 0

	selectedLinks := bot.filterLinks(*links, hasTags, session.SelectedTags, filterValid, filter)
	if len(selectedLinks) == 0 {
		_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "😔 No links found, try to change tags or filters"))
		return nil
	}

	if len(selectedLinks) == 0 {
		return err
	}

	message := tgbotapi.NewMessage(chatID, "")
	message.Text += fmt.Sprintf("Found %d link(s):\n", len(selectedLinks))
	message.ReplyMarkup = buildLinksKeyBoard(selectedLinks, session.Command)

	if _, err = bot.bot.Send(message); err != nil {
		return fmt.Errorf("unable to send message: %w", err)
	}

	if session.Command == listCommand {
		bot.deleteMessage(chatID, session.LastMessageID)
		bot.clearSession(chatID)
	}

	return nil
}

func (bot *BotClient) deleteMessage(chatID int64, messageID *int) {
	if messageID == nil {
		return
	}

	request := tgbotapi.NewDeleteMessage(chatID, *messageID)

	if _, err := bot.bot.Request(request); err != nil {
		slog.Error("message deletion failed", slog.Any("error", err))
	}
}
