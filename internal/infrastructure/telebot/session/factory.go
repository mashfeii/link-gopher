package session

import "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"

func CreateTrackSession(botMessage, userMessage int) *models.TrackSession {
	return &models.TrackSession{
		BaseSession: models.BaseSession{
			State:             models.StateTrackInputURL,
			LastBotMessageID:  &botMessage,
			LastUserMessageID: &userMessage,
		},
	}
}

func CreateListUntrackSession(botMessage, userMessage int) *models.ListUntrackSession {
	return &models.ListUntrackSession{
		BaseSession: models.BaseSession{
			State:             models.StateListUntrackSelectTags,
			LastBotMessageID:  &botMessage,
			LastUserMessageID: &userMessage,
		},
	}
}
