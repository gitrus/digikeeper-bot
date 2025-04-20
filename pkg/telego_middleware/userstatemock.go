package telegomiddleware

import (
	"fmt"
	"runtime"
	"sync"
)

// MockUserStateManager is a mock implementation of UserStateManager using sync.Map.
// It provides thread-safe operations for managing user states in memory.
type MockUserStateManager[S UserState] struct {
	states sync.Map // key: int64, value: S
}

func NewMockUserStateManager[S UserState]() *MockUserStateManager[S] {
	return &MockUserStateManager[S]{}
}

func (m *MockUserStateManager[S]) Fetch(userID int64) (S, error) {
	if value, ok := m.states.Load(userID); ok {
		if state, ok := value.(S); ok {
			return state, nil
		}
		var zero S
		return zero, fmt.Errorf("user state type assertion failed for user ID %d", userID)
	}
	var zero S
	return zero, fmt.Errorf("user state not found for user ID %d", userID)
}

// InitState initializes a new state for the specified user ID.
// It creates a zero value of type S and stores it in the map.
func (m *MockUserStateManager[S]) InitState(userID int64) (S, error) {
	var newState S
	m.states.Store(userID, newState)
	return newState, nil
}

// DropActiveState removes the state for the specified user ID.
// This is useful for cleaning up after a session ends.
func (m *MockUserStateManager[S]) DropActiveState(userID int64) {
	m.states.Delete(userID)
}

// Set updates the state for the specified user ID.
// It uses CompareAndSwap to ensure thread safety and retries up to maxRetries times.
// Returns the new state if successful, or an error if it fails after maxRetries attempts.
func (m *MockUserStateManager[S]) Set(userID int64, state S) (S, error) {
	var oldState S

	var ok bool
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		if ok = m.states.CompareAndSwap(userID, oldState, state); ok {
			return state, nil
		}
		// Yield to other goroutines to avoid busy-waiting
		runtime.Gosched()
	}
	if !ok {
		return oldState, fmt.Errorf("failed to set user state for user ID %d after %d retries", userID, maxRetries)
	}

	return state, nil
}
