package telebot

import (
	"fmt"
	"slices"
	"strings"

	scrapperclient "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func buildLinksKeyBoard(links []scrapperclient.LinkResponse, linkType string) tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()
	currentRow := []tgbotapi.InlineKeyboardButton{}

	for _, link := range links {
		shortenURL := pkg.ShortenURL("https://" + *link.Url)

		var linkButton tgbotapi.InlineKeyboardButton

		switch linkType {
		case "list":
			linkButton = tgbotapi.NewInlineKeyboardButtonURL(
				shortenURL,
				*link.Url,
			)
		default:
			linkButton = tgbotapi.NewInlineKeyboardButtonData(
				shortenURL,
				"untrack "+*link.Url,
			)
		}

		currentRow = append(currentRow, linkButton)
		if len(currentRow) == 3 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)
	}

	return markup
}

func buildTagsKeyboard(tags, selectedTags []string) tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()
	currentRow := []tgbotapi.InlineKeyboardButton{}

	for _, tag := range tags {
		emoji := "◻️"
		if slices.Contains(selectedTags, tag) {
			emoji = "✅"
		}

		tagButton := tgbotapi.NewInlineKeyboardButtonData(
			fmt.Sprintf("%s %s", emoji, tag),
			fmt.Sprintf("tag_toggle %s", tag),
		)

		currentRow = append(currentRow, tagButton)
		if len(currentRow) == 3 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)
	}

	doneTags := "done_tags"

	markup.InlineKeyboard = append(markup.InlineKeyboard,
		[]tgbotapi.InlineKeyboardButton{
			{Text: "✅ Done", CallbackData: &doneTags},
		},
	)

	return markup
}

func buildStringsKeyboard(variants []string, data string) tgbotapi.InlineKeyboardMarkup {
	markup := tgbotapi.NewInlineKeyboardMarkup()
	currentRow := []tgbotapi.InlineKeyboardButton{}

	for _, variant := range variants {
		button := tgbotapi.NewInlineKeyboardButtonData(
			variant,
			fmt.Sprintf("%s %s", data, variant),
		)

		currentRow = append(currentRow, button)

		if len(currentRow) == 3 {
			markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)
			currentRow = nil
		}
	}

	if len(currentRow) > 0 {
		markup.InlineKeyboard = append(markup.InlineKeyboard, currentRow)
	}

	return markup
}

func buildFiltersKeyboard(variants []string) tgbotapi.InlineKeyboardMarkup {
	markup := buildStringsKeyboard(variants, "select_key")
	skipFilters := skipFilters
	returnTags := "return_tags"

	markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "◀️ Step back", CallbackData: &returnTags},
		{Text: "🚫 Skip", CallbackData: &skipFilters},
	})

	return markup
}

func buildFiltersValueKeyboard(variants []string) tgbotapi.InlineKeyboardMarkup {
	markup := buildStringsKeyboard(variants, "select_value")

	skipFilters := skipFilters
	returnFilters := "return_filters"

	markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "◀️ Step back", CallbackData: &returnFilters},
		{Text: "🚫 Skip", CallbackData: &skipFilters},
	})

	return markup
}

func (bot *BotClient) updateTagKeyboard(chatID int64, session *UserSession) error {
	newKeyboard := buildTagsKeyboard(session.AvailableTags, session.SelectedTags)

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

func (bot *BotClient) sendFilterKeyKeyboard(chatID int64, session *UserSession) error {
	availableFilters := []string{}

	for _, filter := range session.AvailableFilters {
		filterName := strings.Split(filter, ":")[0]
		if !slices.Contains(availableFilters, filterName) {
			availableFilters = append(availableFilters, filterName)
		}
	}

	newKeyboard := buildFiltersKeyboard(availableFilters)

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

func (bot *BotClient) sendFilterValueKeyboard(chatID int64, session *UserSession) error {
	availableFilters := []string{}

	for _, filter := range session.AvailableFilters {
		keyAndValue := strings.Split(filter, ":")
		if keyAndValue[0] == session.CurrentFilterName &&
			!slices.Contains(availableFilters, keyAndValue[1]) {
			availableFilters = append(availableFilters, keyAndValue[1])
		}
	}

	newKeyboard := buildFiltersValueKeyboard(availableFilters)

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
