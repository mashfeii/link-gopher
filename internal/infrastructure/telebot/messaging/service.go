package messaging

import (
	standarderrors "errors"
	"fmt"
	"log/slog"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/ui"
)

type MessageService struct {
	bot *tgbotapi.BotAPI
}

const editMessageError = "Bad Request: message is not modified: specified new message content and reply markup are exactly the same as a current content and reply markup of the message" //nolint:lll // ignore line length

func NewMessageService(bot *tgbotapi.BotAPI) *MessageService {
	return &MessageService{bot: bot}
}

func (s *MessageService) Send(chatID int64, text string, opts *models.SendMessageOptions) (tgbotapi.Message, error) {
	msg := tgbotapi.NewMessage(chatID, text)

	if opts != nil {
		msg.ReplyMarkup = opts.ReplyMarkup
		msg.ParseMode = opts.ParseMode
	}

	return s.bot.Send(msg)
}

func (s *MessageService) SendError(chatID int64, err error) (tgbotapi.Message, error) {
	slog.Warn("Error occurred", "error", err)
	text := s.getErrorMessage(err)

	return s.Send(chatID, text, nil)
}

func (s *MessageService) handleEditError(err error, chatID int64, messageID int, text string) (tgbotapi.Message, error) {
	if err != nil && err.Error() == editMessageError {
		return tgbotapi.Message{
			MessageID: messageID,
			Chat:      &tgbotapi.Chat{ID: chatID},
			Text:      text,
		}, nil
	}

	return tgbotapi.Message{}, err
}

func (s *MessageService) EditMessage(chatID int64, messageID int, text string, markup any) (tgbotapi.Message, error) {
	if markup == nil {
		markup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{})
	}

	edit := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		text,
		markup.(tgbotapi.InlineKeyboardMarkup),
	)
	edit.ParseMode = tgbotapi.ModeMarkdown

	message, err := s.bot.Send(edit)
	if err != nil {
		return s.handleEditError(err, chatID, messageID, text)
	}

	return message, nil
}

func (s *MessageService) EditText(chatID int64, messageID int, text string) (tgbotapi.Message, error) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)
	edit.ParseMode = tgbotapi.ModeMarkdown

	message, err := s.bot.Send(edit)
	if err != nil {
		return s.handleEditError(err, chatID, messageID, text)
	}

	return message, nil
}

func (s *MessageService) EditMarkup(chatID int64, messageID int, markup any) (tgbotapi.Message, error) {
	if markup == nil {
		markup = tgbotapi.NewInlineKeyboardMarkup([]tgbotapi.InlineKeyboardButton{})
	}

	edit := tgbotapi.NewEditMessageReplyMarkup(
		chatID,
		messageID,
		markup.(tgbotapi.InlineKeyboardMarkup),
	)

	message, err := s.bot.Send(edit)
	if err != nil {
		return s.handleEditError(err, chatID, messageID, "")
	}

	return message, nil
}

func (s *MessageService) DeleteMessage(chatID int64, messageID int) (*tgbotapi.APIResponse, error) {
	deleteRequest := tgbotapi.NewDeleteMessage(chatID, messageID)

	return s.bot.Request(deleteRequest)
}

func (s *MessageService) getErrorMessage(err error) string {
	if standarderrors.As(err, &errors.ErrErrUnknownHandler{}) {
		return ui.IconCross + " Unknown command. Use /help to see the list of available commands."
	}

	if standarderrors.As(err, &errors.ErrUserNotFound{}) {
		return "You are not registered yet. Use /start to register."
	}

	if standarderrors.As(err, &errors.ErrNoActiveSession{}) {
		return "😅 You are already in the main menu."
	}

	if standarderrors.As(err, &errors.ErrNoLinksFound{}) {
		return ui.IconCross + " No links found. Please add a link first."
	}

	if standarderrors.As(err, &errors.ErrUserAlreadyExists{}) {
		return "😅 You are already registered! Feel free to use the bot."
	}

	if standarderrors.As(err, &errors.ErrLinkAlreadyExists{}) {
		return ui.IconCross + " This link already exists in your list. Please provide a different one."
	}

	if standarderrors.As(err, &errors.ErrInvalidURL{}) {
		return fmt.Sprintf("%s Invalid URL. Try another one:\n• github.com/golang/go\n• stackoverflow.com/questions/17333517/how-to-compile-a-program-in-go-language", //nolint:lll // ignore line length
			ui.IconCross,
		)
	}

	if standarderrors.As(err, &errors.ErrInvalidFilterFormat{}) {
		return ui.IconCross + " Invalid filter format. Please provide a valid filter."
	}

	return ui.IconCross + " Something went wrong. Please try again later."
}
