package ui

import (
	"fmt"
	"slices"
	"sort"

	scrapperclient "github.com/es-debug/backend-academy-2024-go-template/internal/api/openapi/v1/clients/scrapper"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type InlineKeyboardBuilder struct {
	buttonsPerRow int
	rows          [][]tgbotapi.InlineKeyboardButton
}

type ReplyKeyboardBuilder struct {
	buttons [][]tgbotapi.KeyboardButton
}

func NewInlineKeyboardBuilder() *InlineKeyboardBuilder {
	initialKeyboard := make([][]tgbotapi.InlineKeyboardButton, 1)
	initialKeyboard[0] = make([]tgbotapi.InlineKeyboardButton, 0)

	return &InlineKeyboardBuilder{
		rows: initialKeyboard,
	}
}

func NewReplyKeyboardBuilder() *ReplyKeyboardBuilder {
	initialKeyboard := make([][]tgbotapi.KeyboardButton, 1)
	initialKeyboard[0] = make([]tgbotapi.KeyboardButton, 0)

	return &ReplyKeyboardBuilder{
		buttons: initialKeyboard,
	}
}

func (b *ReplyKeyboardBuilder) Button(text string) *ReplyKeyboardBuilder {
	b.buttons = append(b.buttons, []tgbotapi.KeyboardButton{
		{Text: text},
	})

	return b
}

func (b *ReplyKeyboardBuilder) Build() tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewOneTimeReplyKeyboard(b.buttons...)
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

func GetBackSkipKeyboard(backCallback, skipCallback string) tgbotapi.InlineKeyboardMarkup {
	return NewInlineKeyboardBuilder().
		SetButtonsPerRow(2).
		DataButton(IconArrowLeft+" Step back", backCallback).
		DataButton("Skip "+IconArrowRight, skipCallback).
		Build()
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
				fmt.Sprintf("%s%s", models.PrefixUntrack, *link.Url),
			)
		}
	}

	return builder.Build()
}

func BuildTagsKeyboard(tags, selectedTags []string) tgbotapi.InlineKeyboardMarkup {
	builder := NewInlineKeyboardBuilder().SetButtonsPerRow(2)

	sort.Strings(tags)

	for _, tag := range tags {
		emoji := IconCheckbox
		if slices.Contains(selectedTags, tag) {
			emoji = IconChecked
		}

		builder.DataButton(
			fmt.Sprintf("%s %s", emoji, tag),
			fmt.Sprintf("%s%s", models.PrefixTagToggle, tag),
		)
	}

	builder.SetButtonsPerRow(1).DataButton(
		"Done", models.CallbackListDoneTags,
	)

	return builder.Build()
}

func BuildStringsKeyboard(variants []string, data string) tgbotapi.InlineKeyboardMarkup {
	builder := NewInlineKeyboardBuilder().SetButtonsPerRow(2)

	for _, variant := range variants {
		builder.DataButton(variant, fmt.Sprintf("%s%s", data, variant))
	}

	return builder.Build()
}

func BuildFiltersKeyboard(variants []string) tgbotapi.InlineKeyboardMarkup {
	markup := BuildStringsKeyboard(variants, models.PrefixSelectKey)

	skipFilters := models.CallbackListSkipFilters
	returnTags := models.CallbackListReturnTags

	markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: IconArrowLeft + " Step back", CallbackData: &returnTags},
		{Text: "Skip " + IconArrowRight, CallbackData: &skipFilters},
	})

	return markup
}

func BuildFiltersValueKeyboard(variants []string) tgbotapi.InlineKeyboardMarkup {
	markup := BuildStringsKeyboard(variants, models.PrefixSelectValue)

	skipFilters := models.CallbackListSkipFilters
	returnFilters := models.CallbackListReturnFilters

	markup.InlineKeyboard = append(markup.InlineKeyboard, []tgbotapi.InlineKeyboardButton{
		{Text: IconArrowLeft + " Step back", CallbackData: &returnFilters},
		{Text: "Skip " + IconArrowRight, CallbackData: &skipFilters},
	})

	return markup
}
