package loggingctx_test

import (
	"context"
	"log/slog"
	"testing"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/stretchr/testify/assert"
)

func TestAddLogAttr(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		value    interface{}
		setupFn  func(context.Context) context.Context
		expected int
	}{
		{
			name:     "add first attribute",
			key:      "test-key",
			value:    "test-value",
			setupFn:  func(ctx context.Context) context.Context { return ctx },
			expected: 1,
		},
		{
			name:  "override existing attribute",
			key:   "test-key",
			value: "test-value2",
			setupFn: func(ctx context.Context) context.Context {
				return loggingctx.AddLogAttr(ctx, "test-key", "test-value")
			},
			expected: 1,
		},
		{
			name:  "add additional attribute",
			key:   "test-key",
			value: 64,
			setupFn: func(ctx context.Context) context.Context {
				return loggingctx.AddLogAttr(ctx, "another-key", "another-value")
			},
			expected: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Setup
			ctx := context.Background()
			ctx = tc.setupFn(ctx)

			// Execute
			newCtx := loggingctx.AddLogAttr(ctx, tc.key, tc.value)

			// Verify using GetLogAttrs
			attrs := loggingctx.GetLogAttrs(newCtx)
			assert.Equal(t, tc.expected, len(attrs))

			// Check if the new attribute exists
			found := false
			for _, anyAttr := range attrs {
				if anyAttr.(slog.Attr).Key == tc.key {
					found = true
					switch v := tc.value.(type) {
					case int:
						if intVal, ok := anyAttr.(slog.Attr).Value.Any().(int64); ok {
							assert.Equal(t, int64(v), intVal)
						} else {
							assert.Equal(t, v, anyAttr.(slog.Attr).Value.Any())
						}
					default:
						assert.Equal(t, tc.value, anyAttr.(slog.Attr).Value.Any())
					}
					break
				}
			}

			if !found && tc.key != "" {
				t.Errorf("Attribute with key %s not found", tc.key)
			}
		})
	}
}

func TestGetLogAttrsWithDifferentContexts(t *testing.T) {
	ctx1 := context.Background()
	ctx1 = loggingctx.AddLogAttr(ctx1, "key1", "value1")

	ctx2 := context.Background()
	ctx2 = loggingctx.AddLogAttr(ctx2, "key2", "value2")

	attrs1 := loggingctx.GetLogAttrs(ctx1)
	attrs2 := loggingctx.GetLogAttrs(ctx2)

	assert.Equal(t, 1, len(attrs1))
	assert.Equal(t, 1, len(attrs2))
	assert.Equal(t, "value1", attrs1[0].(slog.Attr).Value.Any())
	assert.Equal(t, "value2", attrs2[0].(slog.Attr).Value.Any())
}

func TestInitLogger(t *testing.T) {
	tests := []struct {
		name        string
		environ     string
		expectError bool
	}{
		{
			name:        "development environment",
			environ:     "dev",
			expectError: false,
		},
		{
			name:        "production environment",
			environ:     "prod",
			expectError: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			logger, err := loggingctx.InitLogger(tc.environ)
			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NotNil(t, logger)
			assert.IsType(t, &slog.Logger{}, logger)
		})
	}
}
