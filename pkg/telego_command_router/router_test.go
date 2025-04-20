package telegocommandrouter_test

import (
	"testing"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	tcr "github.com/gitrus/digikeeper-bot/pkg/telego_command_router"
)

type MockBotHandler struct {
	mock.Mock
}

func (m *MockBotHandler) Group(predicates ...th.Predicate) *th.HandlerGroup {
	m.Called(predicates)
	return &th.HandlerGroup{}
}

func (m *MockBotHandler) Handle(handler th.Handler, predicates ...th.Predicate) {
	m.Called(handler, predicates)
}

func TestNewCommandHandlerGroup(t *testing.T) {
	chg := tcr.NewCommandHandlerGroup()
	assert.NotNil(t, chg, "CommandHandlerGroup should not be nil")
}

func TestRegisterCommand(t *testing.T) {
	chg := tcr.NewCommandHandlerGroup()

	// Create a simple test handler
	testHandler := func(ctx *th.Context, update telego.Update) error { return nil }

	// Register a command
	chg.RegisterCommand("test", testHandler, "Test command description")

	commandToHandler := chg.GetRegisteredCommandsInfo()

	// Verify our test command is in the map
	regCmd, ok := commandToHandler["test"]
	assert.True(t, ok, "The 'test' command should be registered")
	assert.Equal(t, "Test command description", regCmd.Description, "Description should match")
}

func TestBindCommandHandlerGroup(t *testing.T) {
	chg := tcr.NewCommandHandlerGroup()

	mockBotHandler := new(MockBotHandler)

	mockBotHandler.On("Group", mock.MatchedBy(func(predicates []th.Predicate) bool {
		return len(predicates) == 1
	})).Return(mockBotHandler)

	testHandler := func(ctx *th.Context, update telego.Update) error { return nil }

	chg.RegisterCommand("test", testHandler, "Test command description")

	mockBotHandler.On("Handle",
		mock.AnythingOfType("telegohandler.Handler"),
		mock.MatchedBy(func(predicates []th.Predicate) bool {
			return len(predicates) == 1
		})).Return().Times(3) // Called for test, help, and unknown commands

	chg.BindCommandsToHandler(mockBotHandler)

	mockBotHandler.AssertCalled(t, "Group", mock.MatchedBy(func(predicates []th.Predicate) bool {
		return len(predicates) == 1
	}))

	mockBotHandler.AssertNumberOfCalls(t, "Handle", 3)

	calls := mockBotHandler.Calls

	assert.Equal(t, 3, len(calls), "Should have exactly 3 calls to Handle")

	for _, call := range calls {
		assert.Equal(t, "Handle", call.Method, "All calls should be to the Handle method")
	}

	for i, call := range calls {
		assert.Equal(t, 2, len(call.Arguments),
			"Call %d should have exactly 2 arguments (handler and predicates)", i)
	}

	for i, call := range calls {
		assert.NotNil(t, call.Arguments[0], "Call %d should have a non-nil handler", i)
	}

	for i, call := range calls {
		predicates := call.Arguments[1].([]th.Predicate)
		assert.Equal(t, 1, len(predicates),
			"Call %d should have exactly 1 predicate", i)
	}
}
