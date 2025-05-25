package commands

import (
	"context"
	"fmt"
	"slices"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/samber/lo"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers/callback"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/session"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type ListUntrackHandler struct{}

func (h *ListUntrackHandler) Handle(ctx context.Context, hctx *handlers.HandlerContext) error {
	const op = "commands.ListUntrackHandler.Handle"

	links, err := hctx.ScrapperService.GetLinks(context.WithValue(ctx, models.ContextKeyChatID, hctx.ChatID))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if len(links) == 0 {
		return fmt.Errorf("%s: no links found: %w", op, errors.NewErrNoLinksFound(hctx.ChatID))
	}

	tags := make(map[string]struct{})
	filters := make(map[string][]string)

	for _, link := range links {
		for _, tag := range link.Tags {
			tags[tag] = struct{}{}
		}

		for _, filter := range link.Filters {
			keyValue := strings.SplitN(filter, ":", 2)

			if len(keyValue) != 2 {
				continue
			}

			if slices.Contains(filters[keyValue[0]], keyValue[1]) {
				continue
			}

			filters[keyValue[0]] = append(filters[keyValue[0]], keyValue[1])
		}
	}

	newSession := session.CreateListUntrackSession(
		models.CommandName(hctx.Update.Message.Command()), nil, &hctx.Update.Message.MessageID,
	)
	newSession.AvailableTags = lo.Keys(tags)
	newSession.AvailableFilters = filters

	hctx.SessionService.Set(hctx.ChatID, newSession)

	if len(tags) == 0 && len(filters) == 0 {
		return callback.DisplayListUntrackResults(
			links,
			hctx.ChatID,
			newSession,
			hctx.MessageService,
			hctx.SessionService,
		)
	} else if len(tags) == 0 {
		hctx.Session = newSession
		return callback.ListUntrackSkipFilters{}.Handle(ctx, hctx)
	}

	message, _ := hctx.MessageService.Send(
		hctx.ChatID,
		ui.IconTag+" Select tags to filter:",
		&models.SendMessageOptions{
			ParseMode:   tgbotapi.ModeMarkdown,
			ReplyMarkup: ui.BuildTagsKeyboard(lo.Keys(tags), nil),
		},
	)

	newSession.LastBotMessageID = &message.MessageID

	hctx.SessionService.Set(hctx.ChatID, newSession)

	return nil
}
