package callback

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type ListUntrackUntrack struct{}

func (l ListUntrackUntrack) Handle(ctx context.Context, hctx *handlers.HandlerContext) error {
	const op = "callback.ListUntrackUntrack.Handle"

	session := hctx.Session.(*models.ListUntrackSession)

	logger := slog.With(slog.String("operation", op), slog.Int64("chatID", hctx.ChatID))

	untrackLink := hctx.Update.CallbackQuery.Data[len(models.PrefixUntrack):]

	logger.Info("Received untrack link callback", slog.String("link", untrackLink))

	serviceContext := context.WithValue(
		ctx,
		models.ContextKeyChatID,
		hctx.ChatID,
	)
	serviceContext = context.WithValue(
		serviceContext,
		models.ContextLinkURL,
		untrackLink,
	)

	err := hctx.ScrapperService.DeleteLink(serviceContext)
	if err != nil {
		logger.Error("Error deleting link", slog.String("link", untrackLink), slog.Any("error", err))
		return fmt.Errorf("%s: unable to delete link: %w", op, err)
	}

	logger.Info("Link deleted successfully", slog.String("link", untrackLink))
	logger.Info("User session", slog.Any("session", session))

	// Update the UI to show the untracked link
	if _, err := hctx.MessageService.EditMessage(
		hctx.ChatID,
		*session.LastBotMessageID,
		ui.IconChecked+" Successfully untracked link:\n• "+untrackLink,
		nil,
	); err != nil {
		logger.Error("Error editing message", slog.String("link", untrackLink), slog.Any("error", err))
		return fmt.Errorf("%s: unable to edit message: %w", op, err)
	}

	logger.Info("Update UI for user", slog.String("link", untrackLink))

	return nil
}
