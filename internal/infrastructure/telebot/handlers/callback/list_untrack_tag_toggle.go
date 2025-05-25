package callback

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/samber/lo"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type ListUntrackTagToggle struct{}

func (l ListUntrackTagToggle) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.TagToggleHandler.Handle"

	callbackData := hctx.Update.CallbackQuery.Data

	tag := strings.TrimPrefix(callbackData, models.PrefixTagToggle)
	if tag == callbackData {
		return fmt.Errorf("%s: invalid callback data: %s", op, callbackData)
	}

	listSession := hctx.Session.(*models.ListUntrackSession)

	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", hctx.ChatID), slog.String("tag", tag))
	logger.Info("Toggling tag")

	if lo.Contains(listSession.SelectedTags, tag) {
		listSession.SelectedTags = lo.Without(listSession.SelectedTags, tag)
	} else {
		listSession.SelectedTags = append(listSession.SelectedTags, tag)
	}

	if _, err := hctx.MessageService.EditMarkup(
		hctx.ChatID,
		*listSession.LastBotMessageID,
		ui.BuildTagsKeyboard(listSession.AvailableTags, listSession.SelectedTags),
	); err != nil {
		logger.Error("Failed to update message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to update message: %w", op, err)
	}

	hctx.SessionService.Set(hctx.ChatID, listSession)
	logger.Info("Tag toggled and session updated")

	return nil
}
