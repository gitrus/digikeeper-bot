package telegomiddleware

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	session "github.com/gitrus/digikeeper-bot/pkg/sessionmanager"
)

type userSessionContextKey struct{}

type UserSessionMiddleware[S session.UserSession] struct {
	repo session.UserSessionManager[S]
}

func NewUserSessionMiddleware[S session.UserSession](repo session.UserSessionManager[S]) *UserSessionMiddleware[S] {
	return &UserSessionMiddleware[S]{repo: repo}
}

func (um *UserSessionMiddleware[S]) WithUserState(ctx context.Context, state S) context.Context {
	return context.WithValue(ctx, userSessionContextKey{}, state)
}

func (um *UserSessionMiddleware[S]) Middleware() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		userID := update.Message.From.ID
		state, err := um.repo.Fetch(ctx, userID)
		if err != nil {
			state, err = um.repo.InitSession(ctx, userID)
			if err != nil {
				slog.ErrorContext(
					ctx,
					"FetchState is missed, InitState failed",
					"error", err,
				)
				return err
			}
		}

		innerCtx := um.WithUserState(ctx.Context(), state)

		return ctx.WithContext(innerCtx).Next(update)
	}
}
