package telegomiddleware

import (
	"errors"
	"runtime"
	"sync"
)

// MockUserStateManager is a mock implementation of UserStateManager using sync.Map.
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
		return zero, errors.New("user state type assertion failed")
	}
	var zero S
	return zero, errors.New("user state not found")
}

func (m *MockUserStateManager[S]) InitState(userID int64) (S, error) {
	var newState S
	m.states.Store(userID, newState)
	return newState, nil
}

func (m *MockUserStateManager[S]) DropActiveState(userID int64) {
	m.states.Delete(userID)
	return
}

func (m *MockUserStateManager[S]) Set(userID int64, state S) (S, error) {
	var oldState S

	var ok bool
	for range 10 {
		if ok := m.states.CompareAndSwap(userID, oldState, state); ok {
			return state, nil
		}
		// cast loop
		runtime.Gosched()
	}
	if !ok {
		return oldState, errors.New("failed to set user state")
	}

	return state, nil
}
