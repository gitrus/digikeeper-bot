package sessionmanager_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/gitrus/digikeeper-bot/pkg/sessionmanager"
)

type MockSession struct {
	mock.Mock
	State   string
	version int
}

func (s *MockSession) GetVersion() int {
	return s.version
}

func NewMockSession(userID int64) (*MockSession, error) {
	return &MockSession{}, nil
}

func TestUserSessionManagerInMem_interfact(t *testing.T) {
	manager := sessionmanager.NewUserSessionManagerInMem[*MockSession](NewMockSession)
	assert.NotNil(t, manager)

	// act
	_, ok := interface{}(manager).(sessionmanager.UserSessionManager[*MockSession])
	//assert
	assert.True(t, ok, "manager should implement UserSessionManager interface")
}

func TestUserSessionManagerInMem_FetchSet(t *testing.T) {
	ctx := context.Background()
	manager := sessionmanager.NewUserSessionManagerInMem[*MockSession](NewMockSession)
	assert.NotNil(t, manager)
	userID := int64(123)
	state, err := manager.InitSession(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, &MockSession{}, state)

	// act
	fetchedState, err := manager.Fetch(ctx, userID)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, state, fetchedState)

	// act
	newSession := &MockSession{State: "action", version: 1}
	updatedState, err := manager.Set(ctx, userID, newSession, 0)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, newSession, updatedState)

	// act
	fetchedState, err = manager.Fetch(ctx, userID)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, newSession, fetchedState)
}

func TestUserSessionManagerInMem_DropActive(t *testing.T) {
	ctx := context.Background()
	manager := sessionmanager.NewUserSessionManagerInMem[*MockSession](NewMockSession)
	assert.NotNil(t, manager)
	userID := int64(123)
	state, err := manager.InitSession(ctx, userID)
	assert.NoError(t, err)
	assert.Equal(t, &MockSession{}, state)

	// act
	newSession := &MockSession{State: "action", version: 1}
	updatedState, err := manager.Set(ctx, userID, newSession, 0)

	// assert
	assert.NoError(t, err)
	assert.Equal(t, newSession, updatedState)

	// act
	manager.DropActive(ctx, userID)

	// assert
	_, err = manager.Fetch(ctx, userID)
	assert.Error(t, err)
	assert.Equal(t, sessionmanager.ErrSessionManagement{
		Reason: "user session not found for user ID 123",
	}, err)
}

func TestUserSessionManagerInMem_FetchEmpty(t *testing.T) {
	ctx := context.Background()
	manager := sessionmanager.NewUserSessionManagerInMem[*MockSession](NewMockSession)
	assert.NotNil(t, manager)

	manager.Set(ctx, 455, &MockSession{}, 1)

	//act
	_, err := manager.Fetch(ctx, 456)

	// assert
	assert.Error(t, err)
	assert.Equal(t, sessionmanager.ErrSessionManagement{
		Reason: "user session not found for user ID 456",
	}, err)
}
