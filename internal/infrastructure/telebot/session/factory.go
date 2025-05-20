package session

import "github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"

func CreateTrackSession(botMessage, userMessage *int) *models.TrackSession {
	return &models.TrackSession{
		BaseSession: models.BaseSession{
			State:             models.StateTrackInputURL,
			LastBotMessageID:  botMessage,
			LastUserMessageID: userMessage,
		},
	}
}

func CreateListUntrackSession(command models.CommandName, botMessage, userMessage *int) *models.ListUntrackSession {
	return &models.ListUntrackSession{
		Command: command,
		BaseSession: models.BaseSession{
			State:             models.StateListUntrackSelectTags,
			LastBotMessageID:  botMessage,
			LastUserMessageID: userMessage,
		},
	}
}
