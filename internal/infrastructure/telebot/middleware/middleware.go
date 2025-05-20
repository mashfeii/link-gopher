package middleware

import (
	"context"

	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/errors"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/handlers"
	"github.com/es-debug/backend-academy-2024-go-template/internal/infrastructure/telebot/models"
)

type (
	Middleware func(handlers.Handler) handlers.Handler
)

func SessionMiddleware() Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(ctx context.Context, hctx *handlers.HandlerContext) error {
			if hctx.Session == nil {
				return errors.NewErrNoActiveSession(hctx.Update.Message.Chat.ID)
			}

			return next.Handle(ctx, hctx)
		})
	}
}

func TrackSessionMiddleware() Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(ctx context.Context, hctx *handlers.HandlerContext) error {
			if hctx.Session == nil {
				return errors.NewErrNoActiveSession(hctx.Update.Message.Chat.ID)
			}

			_, ok := hctx.Session.(*models.TrackSession)
			if !ok {
				return errors.NewErrInvalidSessionType()
			}

			return next.Handle(ctx, hctx)
		})
	}
}

func ListUntrackSessionMiddleware() Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(ctx context.Context, hctx *handlers.HandlerContext) error {
			if hctx.Session == nil {
				return errors.NewErrNoActiveSession(hctx.Update.Message.Chat.ID)
			}

			_, ok := hctx.Session.(*models.ListUntrackSession)
			if !ok {
				return errors.NewErrInvalidSessionType()
			}

			return next.Handle(ctx, hctx)
		})
	}
}

func AuthMiddleware() Middleware {
	return func(next handlers.Handler) handlers.Handler {
		return handlers.HandlerFunc(func(ctx context.Context, hctx *handlers.HandlerContext) error {
			serviceContext := context.WithValue(ctx, models.ContextKeyChatID, hctx.ChatID)

			if _, err := hctx.ScrapperService.GetUser(serviceContext); err != nil {
				return err
			}

			return next.Handle(ctx, hctx)
		})
	}
}
