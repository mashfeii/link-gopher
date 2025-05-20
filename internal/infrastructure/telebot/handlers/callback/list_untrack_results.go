package callback

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/core"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func ValidateSelectedLinks(links []models.LinkData, tags []string, filterName, filterValue string) []models.LinkData {
	selectedLinks := make([]models.LinkData, 0)

	for _, link := range links {
		ok := true

		for _, tag := range tags {
			if !slices.Contains(link.Tags, tag) {
				ok = false
				break
			}
		}

		if !ok {
			continue
		}

		if filterName == "" || filterValue == "" {
			selectedLinks = append(selectedLinks, link)
			continue
		}

		ok = false

		for _, filter := range link.Filters {
			filterNameValue := strings.SplitN(filter, ":", 2)
			if filterNameValue[0] == filterName && filterNameValue[1] == filterValue {
				ok = true
				break
			}
		}

		if ok {
			selectedLinks = append(selectedLinks, link)
		}
	}

	return selectedLinks
}

func DisplayListUntrackResults(
	links []models.LinkData,
	chatID int64,
	session *models.ListUntrackSession,
	messageService core.MessageService,
	sessionManager core.SessionManager,
) error {
	const op = "callback.DisplayListUntrackResults"

	logger := slog.With("operation", op, "chat_id", chatID)

	var message string

	if len(links) == 0 {
		logger.Info("no links found with selected tags")

		message = "No links found with the selected filters."
	} else {
		logger.Info("found links with selected tags", "links", links)

		message = "*Found links with the selected filters:*\n"
	}

	keyboard := ui.NewInlineKeyboardBuilder().SetButtonsPerRow(1)

	for _, link := range links {
		if session.Command == models.CommandNameList {
			keyboard.LinkButton(pkg.ShortenURL(link.URL), link.URL)
		} else {
			keyboard.DataButton(pkg.ShortenURL(link.URL), fmt.Sprintf("%s%s", models.PrefixUntrack, link.URL))
		}
	}

	if session.LastBotMessageID == nil {
		if _, err := messageService.Send(
			chatID,
			message,
			&models.SendMessageOptions{
				ReplyMarkup: keyboard.Build(),
				ParseMode:   tgbotapi.ModeMarkdown,
			},
		); err != nil {
			logger.Error("failed to edit message", "error", err)
			return fmt.Errorf("%s: failed to edit message: %w", op, err)
		}

		logger.Info("sent result message to user")
	} else {
		if _, err := messageService.EditMessage(
			chatID,
			*session.LastBotMessageID,
			message,
			keyboard.Build(),
		); err != nil {
			logger.Error("failed to edit message", "error", err)
			return fmt.Errorf("%s: failed to edit message: %w", op, err)
		}

		logger.Info("edited message with result to user", "message_id", *session.LastBotMessageID)
	}

	if session.Command == models.CommandNameList {
		sessionManager.Clear(chatID)

		logger.Info("cleared session for user", "chat_id", chatID)
	}

	return nil
}
