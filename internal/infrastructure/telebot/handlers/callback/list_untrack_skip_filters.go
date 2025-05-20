package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type ListUntrackSkipFilters struct{}

func (l ListUntrackSkipFilters) Handle(ctx context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.SkipFiltersHandler.Handle"

	listSession := hctx.Session.(*models.ListUntrackSession)

	logger := slog.With(
		"operation", op,
		"chat_id", hctx.ChatID,
		"tags", listSession.AvailableTags,
		"filter name", listSession.FilterName,
		"filter value", listSession.FilterValue,
	)
	serviceContext := context.WithValue(ctx, models.ContextKeyChatID, hctx.ChatID)

	links, err := hctx.ScrapperService.GetLinks(serviceContext)
	if err != nil {
		logger.Error("failed to get links", "error", err)
		return fmt.Errorf("%s: failed to get links: %w", op, err)
	}

	selectedLinks := ValidateSelectedLinks(
		links,
		listSession.SelectedTags,
		listSession.FilterName,
		listSession.FilterValue,
	)
	if err := DisplayListUntrackResults(
		selectedLinks,
		hctx.ChatID,
		listSession,
		hctx.MessageService,
		hctx.SessionService,
	); err != nil {
		return err
	}

	logger.Info("edited message to user", "message_id", *listSession.LastBotMessageID)

	return nil
}
