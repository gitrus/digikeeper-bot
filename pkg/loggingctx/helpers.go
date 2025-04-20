package loggingctx

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/exp/zapslog"
	"go.uber.org/zap/zapcore"
)

type contextKey string

// LogAttrsKey is the context key used for storing logging attributes
const LogAttrsKey contextKey = "LogAttrsKey"

// AddLogAttr adds a logging attribute to the provided context
func AddLogAttr(ctx context.Context, key string, value any) context.Context {
	attrs, ok := ctx.Value(LogAttrsKey).([]slog.Attr)
	if !ok {
		attrs = []slog.Attr{}
	}

	found := false
	for i := range attrs {
		if attrs[i].Key == key {
			attrs[i].Value = slog.AnyValue(value)
			found = true
			break
		}
	}

	if !found {
		attrs = append(attrs, slog.Any(key, value))
	}

	return context.WithValue(ctx, LogAttrsKey, attrs)
}

// GetLogAttrs retrieves logging attributes slice from the context
func GetLogAttrs(ctx context.Context) []any {
	attrs, ok := ctx.Value(LogAttrsKey).([]slog.Attr)
	if !ok {
		return []any{}
	}

	anyattr := make([]any, len(attrs))
	for i, attr := range attrs {
		anyattr[i] = attr
	}
	return anyattr
}

func InitLogger(environ string) (*slog.Logger, error) {
	var logger *zap.Logger
	var err error
	if strings.HasPrefix(environ, "dev") {
		config := zap.Config{
			Encoding:         "json", // Use JSON encoding
			Level:            zap.NewAtomicLevelAt(zap.InfoLevel),
			OutputPaths:      []string{"stdout"},
			ErrorOutputPaths: []string{"stderr"},
			EncoderConfig: zapcore.EncoderConfig{
				TimeKey:      "time",
				LevelKey:     "level",
				MessageKey:   "msg",
				EncodeTime:   zapcore.ISO8601TimeEncoder,
				EncodeLevel:  zapcore.CapitalLevelEncoder,
				CallerKey:    "caller",
				EncodeCaller: zapcore.ShortCallerEncoder,
			},
		}
		logger, err = config.Build()
	} else {
		logger, err = zap.NewProduction()
	}

	if err != nil {
		return nil, fmt.Errorf("Fail at init zap logger %w", err)
	}

	handler := zapslog.NewHandler(logger.Core())
	slogLogger := slog.New(handler)

	return slogLogger, nil
}
