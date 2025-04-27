package sessionmanager

import "context"

type userSessionContextKey struct{}

type UserSession interface {
	GetVersion() int
}

type NewUserSession[S UserSession] func(userID int64) (S, error)

type UserSessionManager[S UserSession] interface {
	InitSession(ctx context.Context, userID int64) (S, error)
	Fetch(ctx context.Context, userID int64) (S, error)
	Set(ctx context.Context, userID int64, newSession S, prevVersion int) (S, error)
	DropActive(ctx context.Context, userID int64) error
}

