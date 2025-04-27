package sessionmanager

import (
	"context"
	"fmt"
	"sync"
)

// MockUserStateManager is a mock implementation of UserStateManager using sync.Map.
// It provides thread-safe operations for managing user states in memory.
type UserSessionManagerInMem[S UserSession] struct {
	sessions sync.Map // key: int64, value: S

	userSessionFabric NewUserSession[S]
}

type ErrSessionManagement struct {
	Reason string
}

func (e ErrSessionManagement) Error() string {
	return fmt.Sprintf("session fetch error %s", e.Reason)
}

func NewUserSessionManagerInMem[S UserSession](usf NewUserSession[S]) *UserSessionManagerInMem[S] {
	return &UserSessionManagerInMem[S]{userSessionFabric: usf}
}

// InitState initializes a new session with NewUserSession(userID).
func (usm *UserSessionManagerInMem[S]) InitSession(ctx context.Context, userID int64) (S, error) {
	var zero S
	newSession, err := usm.userSessionFabric(userID)
	if err != nil {
		return zero, err
	}
	usm.sessions.Store(userID, newSession)

	return newSession, nil
}

func (m *UserSessionManagerInMem[S]) Fetch(
	ctx context.Context, userID int64,
) (S, error) {
	var result S
	if value, ok := m.sessions.Load(userID); ok {
		if state, ok := value.(S); ok {
			return state, nil
		}
		return result, ErrSessionManagement{
			Reason: fmt.Sprintf("user session type assertion failed for user ID %d", userID),
		}
	}

	return result, ErrSessionManagement{
		Reason: fmt.Sprintf("user session not found for user ID %d", userID),
	}
}

// DropActive removes the state for the specified user ID.
func (m *UserSessionManagerInMem[S]) DropActive(ctx context.Context, userID int64) error {
	m.sessions.Delete(userID)
	return nil
}

// Set updates the session for the specified userID and version
// Returns the new session if successful, or an error
func (m *UserSessionManagerInMem[S]) Set(
	ctx context.Context, userID int64, newSession S, prevVersion int,
) (S, error) {
	oldValue, loaded := m.sessions.Load(userID)
	if !loaded {
		return newSession, ErrSessionManagement{Reason: "session not found"}
	}
	var oldSession S = oldValue.(S)
	if oldSession.GetVersion() != prevVersion {
		return oldSession, ErrSessionManagement{Reason: "version mismatch"}
	}

	swapped := m.sessions.CompareAndSwap(userID, oldSession, newSession)
	if !swapped {
		err := ErrSessionManagement{
			Reason: fmt.Sprintf("failed to Set for user ID %d", userID),
		}
		return oldSession, err
	}

	return newSession, nil
}

type SimpleUserSession struct {
	UserID  int64
	State   string
	Version int
}

func NewSimpleUserSession(userID int64) (*SimpleUserSession, error) {
	return &SimpleUserSession{
		UserID:  userID,
		State:   "",
		Version: 1,
	}, nil
}

func (s SimpleUserSession) GetVersion() int {
	return s.Version
}
