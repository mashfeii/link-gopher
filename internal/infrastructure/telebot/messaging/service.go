package messaging

import (
	standarderrors "errors"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
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
	edit := tgbotapi.NewEditMessageTextAndMarkup(
		chatID,
		messageID,
		text,
		markup.(tgbotapi.InlineKeyboardMarkup),
	)

	message, err := s.bot.Send(edit)
	if err != nil {
		return s.handleEditError(err, chatID, messageID, text)
	}

	return message, nil
}

func (s *MessageService) EditText(chatID int64, messageID int, text string) (tgbotapi.Message, error) {
	edit := tgbotapi.NewEditMessageText(chatID, messageID, text)

	message, err := s.bot.Send(edit)
	if err != nil {
		return s.handleEditError(err, chatID, messageID, text)
	}

	return message, nil
}

func (s *MessageService) EditMarkup(chatID int64, messageID int, markup any) (tgbotapi.Message, error) {
	if markup == nil {
		markup = tgbotapi.NewInlineKeyboardMarkup([][]tgbotapi.InlineKeyboardButton{}...)
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
	if standarderrors.As(err, &errors.ErrUserNotFound{}) {
		return "You are not registered yet. Use /start to register."
	}

	if standarderrors.As(err, &errors.ErrNoActiveSession{}) {
		return "😅 You are already in the main menu."
	}

	if standarderrors.As(err, &errors.ErrUserAlreadyExists{}) {
		return "😅 You are already registered! Feel free to use the bot."
	}

	if standarderrors.As(err, &errors.ErrLinkAlreadyExists{}) {
		return "❌ This link already exists in your list. Please provide a different one."
	}

	if standarderrors.As(err, &errors.ErrInvalidURL{}) {
		return "❌ Invalid URL format. Please provide a valid URL."
	}

	if standarderrors.As(err, &errors.ErrInvalidFilterFormat{}) {
		return "❌ Invalid filter format. Please provide a valid filter."
	}

	return "❌ Something went wrong. Please try again later."
}
