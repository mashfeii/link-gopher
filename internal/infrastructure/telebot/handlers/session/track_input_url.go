package session

import (
	"context"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
	"github.com/es-debug/backend-academy-2024-go-template/pkg"
)

type TrackInputURL struct {
	hctx *handlers.HandlerContext
}

func (h *TrackInputURL) Handle(
	ctx context.Context,
	hctx *handlers.HandlerContext,
) error {
	const op = "handlers.TrackInputURL.HandleSession"

	h.hctx = hctx

	trackSession := hctx.Session.(*models.TrackSession)

	URL := hctx.Update.Message.Text
	trackSession.LastUserMessageID = &hctx.Update.Message.MessageID

	logger := h.createLogger(op, URL, hctx.ChatID)
	serviceContext := context.WithValue(ctx, models.ContextKeyChatID, hctx.ChatID)

	if ok, err := h.checkForDuplicateLink(serviceContext, trackSession, URL, hctx.ChatID, logger); !ok || err != nil {
		return err
	}

	if ok, err := h.validateURL(URL, trackSession, hctx.ChatID, logger); !ok || err != nil {
		return err
	}

	if err := h.updateUIAndSession(hctx.ChatID, URL, trackSession, logger); err != nil {
		return err
	}

	return nil
}

func (h *TrackInputURL) createLogger(op, url string, chatID int64) *slog.Logger {
	return slog.With(
		slog.String("operation", op),
		slog.String("url", url),
		slog.Int64("chatID", chatID),
	)
}

func (h *TrackInputURL) checkForDuplicateLink(
	ctx context.Context,
	trackSession *models.TrackSession,
	url string,
	chatID int64,
	logger *slog.Logger,
) (bool, error) {
	const op = "handlers.TrackInputURL.checkForDuplicateLink"

	links, err := h.hctx.ScrapperService.GetLinks(ctx)
	if err != nil {
		logger.Error("failed to get links", slog.Any("error", err))
		return false, fmt.Errorf("%s: failed to get links: %w", op, err)
	}

	for _, link := range links {
		if link.URL != url {
			continue
		}

		logger.Warn("Duplicate link found")

		return false, h.handleDuplicateLink(chatID, trackSession, logger)
	}

	return true, nil
}

func (h *TrackInputURL) handleDuplicateLink(
	chatID int64,
	trackSession *models.TrackSession,
	logger *slog.Logger,
) error {
	const op = "handlers.TrackInputURL.handleDuplicateLink"

	_, err := h.hctx.MessageService.EditText(
		chatID,
		*trackSession.LastBotMessageID,
		"❗️ This link already exists in your list. Please enter a different URL.",
	)
	if err != nil {
		logger.Error("failed to edit message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to edit message: %w", op, err)
	}

	if _, err := h.hctx.MessageService.DeleteMessage(chatID, *trackSession.LastUserMessageID); err != nil {
		logger.Warn("failed to delete message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to delete message: %w", op, err)
	}

	return nil
}

func (h *TrackInputURL) validateURL(
	url string,
	trackSession *models.TrackSession,
	chatID int64,
	logger *slog.Logger,
) (bool, error) {
	const op = "handlers.TrackInputURL.validateURL"

	_, _, gitErr := pkg.ValidateGithubURL(url)
	_, stackErr := pkg.ValidateStackOverflowURL(url)

	if gitErr != nil && stackErr != nil {
		logger.Warn("invalid URL format", slog.String("opeartion", op), slog.Any("error", stackErr))
		return false, h.handleInvalidURL(chatID, trackSession, logger)
	}

	return true, nil
}

func (h *TrackInputURL) handleInvalidURL(
	chatID int64,
	trackSession *models.TrackSession,
	logger *slog.Logger,
) error {
	const op = "handlers.TrackInputURL.handleInvalidURL"

	_, err := h.hctx.MessageService.EditText(
		chatID,
		*trackSession.LastBotMessageID,
		fmt.Sprintf("%s Invalid URL. Try another one:\n• github.com/golang/go\n• stackoverflow.com/questions/17333517/how-to-compile-a-program-in-go-language", //nolint:lll // ignore line length
			ui.IconCross,
		),
	)
	if err != nil {
		logger.Error("failed to edit message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to edit message: %w", op, err)
	}

	if _, err := h.hctx.MessageService.DeleteMessage(chatID, *trackSession.LastUserMessageID); err != nil {
		logger.Warn("failed to delete message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to delete message: %w", op, err)
	}

	return nil
}

func (h *TrackInputURL) updateUIAndSession(
	chatID int64,
	url string,
	trackSession *models.TrackSession,
	logger *slog.Logger,
) error {
	const op = "handlers.TrackInputURL.updateUIAndSession"

	message, err := h.hctx.MessageService.Send(
		chatID,
		ui.IconTag+" Enter tags (_space-separated_):",
		&models.SendMessageOptions{
			ParseMode:   tgbotapi.ModeMarkdown,
			ReplyMarkup: ui.GetBackSkipKeyboard(models.CallbackReturnURL, models.CallbackSkipTags),
		},
	)
	if err != nil {
		logger.Error("failed to send message", slog.Any("error", err))
		return fmt.Errorf("%s: failed to send message: %w", op, err)
	}

	trackSession.URL = url
	trackSession.State = models.StateTrackInputTags
	trackSession.LastBotMessageID = &message.MessageID

	logger.Info("Get the link for track session")

	h.hctx.SessionService.Set(chatID, trackSession)

	return nil
}
