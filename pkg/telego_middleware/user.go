package telegomiddleware

import (
	"context"
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

type userStateContextKey struct{}

type UserState interface{}

type UserStateManager[S UserState] interface {
	Fetch(userId int64) (S, error)
	InitState(userId int64) (S, error)
}

type UserStateMiddleware[S UserState] struct {
	repo UserStateManager[S]
}

func NewUserStateMiddleware[S UserState](repo UserStateManager[S]) *UserStateMiddleware[S] {
	return &UserStateMiddleware[S]{repo: repo}
}

func (um *UserStateMiddleware[S]) WithUserState(ctx context.Context, state S) context.Context {
	return context.WithValue(ctx, userStateContextKey{}, state)
}

func (um *UserStateMiddleware[S]) Middleware() th.Middleware {
	return func(bot *telego.Bot, update telego.Update, next th.Handler) {
		userId := update.Message.From.ID
		state, err := um.repo.Fetch(userId)
		if err != nil {
			state, err = um.repo.InitState(userId)
			if err != nil {
				slog.ErrorContext(
					update.Context(),
					"FetchState is missed, InitState failed",
					"error", err,
				)
				return
			}
		}

		ctx := um.WithUserState(update.Context(), state)

		next(bot, update.WithContext(ctx))
	}
}
