package ui

import (
	"fmt"
	"slices"

	scrapperclient "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	untrackPrefix         = "untrack "
	tagTogglePrefix       = "tag_toggle "
	selectKeyPrefix       = "select_key"
	selectValuePrefix     = "select_value"
	doneTagsCallback      = "done_tags"
	skipFiltersCallback   = "skip_filters"
	returnTagsCallback    = "return_tags"
	returnFiltersCallback = "return_filters"
)

type InlineKeyboardBuilder struct {
	buttonsPerRow int
	rows          [][]tgbotapi.InlineKeyboardButton
}

func NewInlineKeyboardBuilder() *InlineKeyboardBuilder {
	initialKeyboard := make([][]tgbotapi.InlineKeyboardButton, 1)
	initialKeyboard[0] = make([]tgbotapi.InlineKeyboardButton, 0)

	return &InlineKeyboardBuilder{
		rows: initialKeyboard,
	}
}

func (b *InlineKeyboardBuilder) SetButtonsPerRow(buttonsPerRow int) *InlineKeyboardBuilder {
	b.buttonsPerRow = buttonsPerRow
	return b
}

func (b *InlineKeyboardBuilder) checkAutoRow() {
	if b.buttonsPerRow <= 0 {
		return
	}

	currentRow := len(b.rows) - 1
	if len(b.rows[currentRow]) >= b.buttonsPerRow {
		b.Row()
	}
}

func (b *InlineKeyboardBuilder) Row() *InlineKeyboardBuilder {
	b.rows = append(b.rows, []tgbotapi.InlineKeyboardButton{})
	return b
}

func (b *InlineKeyboardBuilder) DataButton(text, callback string) *InlineKeyboardBuilder {
	b.checkAutoRow()
	currentRow := len(b.rows) - 1
	button := tgbotapi.NewInlineKeyboardButtonData(text, callback)
	b.rows[currentRow] = append(b.rows[currentRow], button)

	return b
}

func (b *InlineKeyboardBuilder) LinkButton(text, url string) *InlineKeyboardBuilder {
	b.checkAutoRow()
	currentRow := len(b.rows) - 1
	button := tgbotapi.NewInlineKeyboardButtonURL(text, url)
	b.rows[currentRow] = append(b.rows[currentRow], button)

	return b
}

func (b *InlineKeyboardBuilder) Build() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(b.rows...)
}

func BuildLinksKeyBoard(links []scrapperclient.LinkResponse, linkType string) tgbotapi.InlineKeyboardMarkup {
	builder := NewInlineKeyboardBuilder().SetButtonsPerRow(2)

	for _, link := range links {
		shortenURL := pkg.ShortenURL("https://" + *link.Url)

		switch linkType {
		case "list":
			builder.LinkButton(shortenURL, *link.Url)
		default:
			builder.DataButton(
				shortenURL,
				fmt.Sprintf("%s %s", untrackPrefix, *link.Url),
			)
		}
	}

	return builder.Build()
}

func BuildTagsKeyboard(tags, selectedTags []string) tgbotapi.InlineKeyboardMarkup {
	builder := NewInlineKeyboardBuilder().SetButtonsPerRow(2)

	for _, tag := range tags {
		emoji := "◻️"
		if slices.Contains(selectedTags, tag) {
			emoji = "✅"
		}

		builder.DataButton(
			fmt.Sprintf("%s %s", emoji, tag),
			fmt.Sprintf("%s %s", tagTogglePrefix, tag),
		)
	}

	builder.SetButtonsPerRow(1).DataButton(
		"✅ Done", doneTagsCallback,
	)

	return builder.Build()
}

func BuildStringsKeyboard(variants []string, data string) tgbotapi.InlineKeyboardMarkup {
	builder := NewInlineKeyboardBuilder().SetButtonsPerRow(2)

	for _, variant := range variants {
		builder.DataButton(variant, fmt.Sprintf("%s %s", data, variant))
	}

	return builder.Build()
}

func BuildFiltersKeyboard(variants []string) tgbotapi.InlineKeyboardMarkup {
	markup := BuildStringsKeyboard(variants, selectKeyPrefix)

	skipFilters := skipFiltersCallback
	returnTags := returnTagsCallback

	markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "◀️ Step back", CallbackData: &returnTags},
		{Text: "🚫 Skip", CallbackData: &skipFilters},
	})

	return markup
}

func BuildFiltersValueKeyboard(variants []string) tgbotapi.InlineKeyboardMarkup {
	markup := BuildStringsKeyboard(variants, selectValuePrefix)

	skipFilters := skipFiltersCallback
	returnFilters := returnFiltersCallback

	markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: "◀️ Step back", CallbackData: &returnFilters},
		{Text: "🚫 Skip", CallbackData: &skipFilters},
	})

	return markup
}
