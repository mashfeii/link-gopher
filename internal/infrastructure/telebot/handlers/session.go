package handlers

import (
	"context"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type BaseSessionHandler struct {
	BotAPI       *tgbotapi.BotAPI
	SessionTypes []models.SessionType
}

func (b BaseSessionHandler) CanHandle(ctx context.Context, update *tgbotapi.Update) bool {
	//TODO implement me
	panic("implement me")
}

func (b BaseSessionHandler) Handle(ctx context.Context, update *tgbotapi.Update) error {
	//TODO implement me
	panic("implement me")
}

func (b BaseSessionHandler) CanHandleState(state models.State) bool {
	//TODO implement me
	panic("implement me")
}

func (b BaseSessionHandler) HandleSession(ctx context.Context, update *tgbotapi.Update, session *models.UserSession) error {
	//TODO implement me
	panic("implement me")
}
