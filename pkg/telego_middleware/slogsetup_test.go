package telegomiddleware

import (
	"log/slog"
	"sync"
	"testing"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Test the firstNRunes function
func TestFirstNRunes(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		n        int
		expected string
	}{
		{
			name:     "empty string",
			input:    "",
			n:        5,
			expected: "",
		},
		{
			name:     "string shorter than n",
			input:    "hello",
			n:        10,
			expected: "hello",
		},
		{
			name:     "string longer than n",
			input:    "hello world",
			n:        5,
			expected: "hello",
		},
		{
			name:     "string with unicode characters",
			input:    "こんにちは世界",
			n:        3,
			expected: "こんに",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := FirstNRunes(tc.input, tc.n)
			assert.Equal(t, tc.expected, result)
		})
	}
}

// Test that AddSlogAttrs returns a function and checks its behavior
func TestAddSlogAttrs(t *testing.T) {
	handler := AddSlogAttrs()
	assert.NotNil(t, handler, "AddSlogAttrs should return a non-nil handler")

	handlerType := assert.IsType(
		t,
		(th.Handler)(nil),
		handler,
		"AddSlogAttrs should return a th.Handler",
	)
	assert.True(t, handlerType, "Handler should be of type th.Handler")
}

func TestAddSlogAttrsHandle(t *testing.T) {
	token := "1234567890:aaaabbbbaaaabbbbaaaabbbbaaaabbbbccc"
	bot, err := telego.NewBot(token)
	require.NoError(t, err)
	updates := make(chan telego.Update, 10)

	bh, err := th.NewBotHandler(bot, updates)
	require.NoError(t, err)

	wg := sync.WaitGroup{}

	wg.Add(1)
	handlerCalled := false
	handler := func(ctx *th.Context, msg telego.Message) error {
		defer wg.Done()
		handlerCalled = true

		attrs := loggingctx.GetLogAttrs(ctx)
		assert.NotEmpty(t, attrs, "attrs persistence error")
		assert.Len(t, attrs, 5, "expected 5 attributes in loggingctx")

		attrsMap := make(map[string]any)
		for _, attr := range attrs {
			slogAttr := attr.(slog.Attr)
			attrsMap[slogAttr.Key] = slogAttr.Value.Any()
		}

		updateID, _ := attrsMap["update_id"]
		assert.Equal(t, int64(999), updateID, "update_id should match expected value")

		messageID, _ := attrsMap["message_id"]
		assert.Equal(t, int64(123), messageID, "message_id should match expected value")

		textFirst10, _ := attrsMap["text_first10"]
		assert.Equal(t, "Test messa", textFirst10, "text_first10 should match expected value")

		chatID, _ := attrsMap["chat_id"]
		assert.Equal(t, int64(456), chatID, "chat_id should match expected value")

		userID, _ := attrsMap["user_id"]
		assert.Equal(t, int64(789), userID, "user_id should match expected value")

		return nil
	}

	bh.Use(AddSlogAttrs())
	bh.HandleMessage(handler)

	go bh.Start()

	testUpdate := telego.Update{
		UpdateID: 999,
		Message: &telego.Message{
			MessageID: 123,
			Text:      "Test message",
			Chat:      telego.Chat{ID: 456},
			From:      &telego.User{ID: 789},
		},
	}
	updates <- testUpdate

	close(updates)
	wg.Wait()
	bh.Stop()

	assert.True(t, handlerCalled, "Handler should have been called")
}
