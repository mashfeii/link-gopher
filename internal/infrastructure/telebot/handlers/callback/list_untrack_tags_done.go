package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/samber/lo"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type ListUntrackTagsDone struct{}

func (l ListUntrackTagsDone) Handle(ctx context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.ListUntrackTagsDone.Handle"

	listSession := hctx.Session.(*models.ListUntrackSession)

	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", hctx.ChatID))
	logger.Info("Tags done, processing with filter")

	if len(listSession.AvailableFilters) == 0 {
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

		return DisplayListUntrackResults(
			selectedLinks,
			hctx.ChatID,
			listSession,
			hctx.MessageService,
			hctx.SessionService,
		)
	}

	if _, err := hctx.MessageService.EditMessage(
		hctx.ChatID,
		*listSession.LastBotMessageID,
		ui.IconSearch+" Please select the filter name you want to apply",
		ui.BuildFiltersKeyboard(lo.Keys(listSession.AvailableFilters)),
	); err != nil {
		logger.Error("Failed to update message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to update message: %w", op, err)
	}

	logger.Info("Tags done, handling filter selection")

	return nil
}
