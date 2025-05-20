package commands

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/session"
)

type TrackHandler struct{}

func (h *TrackHandler) Handle(_ context.Context, hctx *handlers.HandlerContext) error {
	// First action is to ask for the URL
	message, _ := hctx.MessageService.Send(hctx.ChatID, "🔗 Please enter the URL you want to track", nil)

	// Create typed session
	hctx.SessionService.Set(hctx.ChatID, session.CreateTrackSession(&message.MessageID, &hctx.Update.Message.MessageID))

	return nil
}
