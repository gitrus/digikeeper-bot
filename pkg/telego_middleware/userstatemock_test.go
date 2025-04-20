package telegomiddleware_test

import (
	"testing"

	"github.com/gitrus/digikeeper-bot/pkg/telego_middleware"
	"github.com/stretchr/testify/assert"
)

func TestMockUserStateManager(t *testing.T) {
	// Create mock manager
	manager := telegomiddleware.NewMockUserStateManager[string]()
	assert.NotNil(t, manager)

	// Test InitState - should create a new state
	userID := int64(123)
	state, err := manager.InitState(userID)
	assert.NoError(t, err)
	assert.Equal(t, "", state) // Default string state is empty string

	// Test Fetch - should now return the state we just created
	fetchedState, err := manager.Fetch(userID)
	assert.NoError(t, err)
	assert.Equal(t, state, fetchedState)

	// Test Set - should update the state
	newState := "active"
	updatedState, err := manager.Set(userID, newState)
	assert.NoError(t, err)
	assert.Equal(t, newState, updatedState)

	// Test Fetch again - should return the updated state
	fetchedState, err = manager.Fetch(userID)
	assert.NoError(t, err)
	assert.Equal(t, newState, fetchedState)

	// Test DropActiveState - should remove the state
	manager.DropActiveState(userID)

	// Test Fetch after drop - should return an error
	_, err = manager.Fetch(userID)
	assert.Error(t, err)

	// Test Fetch for non-existent user - should return an error
	_, err = manager.Fetch(int64(456))
	assert.Error(t, err)
}