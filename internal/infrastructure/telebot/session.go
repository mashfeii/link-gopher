package telebot

import (
	"fmt"
	"log/slog"
	"slices"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type (
	State       int
	SessionType int
)

const (
	SessionTypeTrack = iota
	SessionTypeListUntrack
)

const (
	StateWaitingURL = iota
	StateWaitingTags
	StateWaitingFilters
	StateWaitingTagsSelection
	StateWaitingFiltersSelection
)

type UserSession struct {
	Type    SessionType
	State   State
	Command string

	URL     string
	Tags    []string
	Filters []string

	SelectedTags       []string
	AvailableTags      []string
	AvailableFilters   []string
	CurrentFilterName  string
	CurrentFilterValue string
	LastMessageID      *int
	LastUserMessageID  *int
}

func (bot *BotClient) createTypedSession(sessionType SessionType, command string, filters, tags []string) *UserSession {
	return &UserSession{
		Type:             sessionType,
		State:            StateWaitingTagsSelection,
		Command:          command,
		AvailableFilters: filters,
		AvailableTags:    tags,
	}
}

func (s *UserSession) toggleTag(tag string) {
	for i, selectedTag := range s.SelectedTags {
		if selectedTag == tag {
			s.SelectedTags = slices.Delete(s.SelectedTags, i, i+1)
			return
		}
	}

	s.SelectedTags = append(s.SelectedTags, tag)
}

func (bot *BotClient) getSession(chatID int64) *UserSession {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	return bot.sessions[chatID]
}

func (bot *BotClient) setSession(chatID int64, session *UserSession) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	bot.sessions[chatID] = session
}

func (bot *BotClient) createTrackSession(chatID int64, initialState State) {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	if _, err := bot.handleAuthorizationCheck(chatID); err != nil {
		slog.Error("unable to validate user", "error", err)
		return
	}

	session := &UserSession{
		Type:  SessionTypeTrack,
		State: initialState,
	}
	message, _ := bot.bot.Send(tgbotapi.NewMessage(chatID, "🔗 Please enter the URL you want to track"))
	session.LastMessageID = &message.MessageID
	bot.sessions[chatID] = session
}

func (bot *BotClient) createListUntrackSession(chatID int64, command string) error {
	bot.mu.Lock()
	defer bot.mu.Unlock()

	resp, err := bot.handleAuthorizationCheck(chatID)
	if err != nil {
		return err
	}

	links, err := validateLinksResponse(chatID, bot, resp)
	if err != nil || links == nil {
		return err
	}

	availableFilters, availableTags := bot.extractFiltersAndTags(*links)
	session := bot.createTypedSession(SessionTypeListUntrack, command, availableFilters, availableTags)

	return bot.handleTypedSession(chatID, session)
}

func (bot *BotClient) handleSession(update *tgbotapi.Update, session *UserSession) error {
	chatID := update.Message.Chat.ID
	text := strings.TrimPrefix(update.Message.Text, "https://")

	if update.Message.IsCommand() {
		if update.Message.Command() != "cancel" {
			bot.clearSession(chatID)
			_, _ = bot.bot.Send(tgbotapi.NewMessage(chatID, "Okey, let's go with another command 🤖"))
		}

		return bot.handleCommand(update)
	}

	switch session.State {
	case StateWaitingURL:
		ok, err := bot.validateLink(chatID, text)
		if err != nil {
			return err
		}

		if !ok {
			message, _ := bot.bot.Send(tgbotapi.NewMessage(chatID,
				"❌ This link is already being tracked. Please enter another one or use /untrack to remove it."))
			session.LastMessageID = &message.MessageID
			bot.setSession(chatID, session)

			return nil
		}

		session.URL = text
		session.State = StateWaitingTags

		message := bindKeyboardMessage(chatID, "📋 Enter tags (space-separated):",
			[][]string{{"◀️ Step back", "return_url"}, {"🚫 Skip", "skip_tags"}})

		messageSend, err := bot.bot.Send(message)
		if err != nil {
			return fmt.Errorf("unable to send message: %w", err)
		}

		session.LastMessageID = &messageSend.MessageID
		bot.setSession(chatID, session)
	case StateWaitingTags:
		session.Tags = strings.Fields(text)
		session.State = StateWaitingFilters

		message := bindKeyboardMessage(chatID, "📋 Enter filters (space-separated key:value pairs):",
			[][]string{{"◀️ Step back", returnTags}, {"🚫 Skip", "skip_filters"}})

		messageSend, err := bot.bot.Send(message)
		if err != nil {
			return err
		}

		session.LastMessageID = &messageSend.MessageID
		bot.setSession(chatID, session)
	case StateWaitingFilters:
		session.Filters = strings.Fields(text)
		if err := bot.saveTracking(chatID, session, "new"); err != nil {
			return fmt.Errorf("unable to save tracking: %w", err)
		}

		bot.clearSession(chatID)
	default:
		request := tgbotapi.NewDeleteMessage(chatID, *session.LastUserMessageID)
		if _, err := bot.bot.Request(request); err != nil {
			return fmt.Errorf("unable to delete message: %w", err)
		}
	}

	return nil
}

func (bot *BotClient) handleTypedSession(chatID int64, session *UserSession) error {
	switch {
	case len(session.AvailableTags) == 0 && len(session.AvailableFilters) == 0:
		return bot.finalizeListUntrackSession(chatID, session)

	case len(session.AvailableTags) == 0:
		session.State = StateWaitingFiltersSelection
		if err := bot.sendFilterKeyKeyboard(chatID, session); err != nil {
			return fmt.Errorf("unable to send filter key keyboard: %w", err)
		}

	default:
		msg := tgbotapi.NewMessage(chatID, "🏷 Select tags to filter:")
		msg.ReplyMarkup = buildTagsKeyboard(session.AvailableTags, session.SelectedTags)

		sentMsg, err := bot.bot.Send(msg)
		if err != nil {
			return fmt.Errorf("failed to send message: %w", err)
		}

		session.LastMessageID = &sentMsg.MessageID
	}

	bot.sessions[chatID] = session

	return nil
}

func (bot *BotClient) clearSession(chatID int64) {
	// BUG: mutex switch leads to lock
	delete(bot.sessions, chatID)
}

func (bot *BotClient) handleInvalidURL(chatID int64, session *UserSession) {
	message := tgbotapi.NewMessage(
		chatID,
		`❌ Invalid URL. Try another one:
  ✅ github.com/golang/go
  ✅ stackoverflow.com/questions/17333517/how-to-compile-a-program-in-go-language`,
	)

	if _, err := bot.bot.Send(message); err != nil {
		slog.Error("unable to send message", slog.Any("error", err))
	}

	session.State = StateWaitingURL
	bot.setSession(chatID, session)
}
