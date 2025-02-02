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

const logAttrsKey contextKey = "log_attrs"

// AddLogAttr adds a logging attribute to the provided context.
//
// Note: Attributes is a slice, so if there is already an attribute with the same key,
// the new attribute is appended to the existing slice.
func AddLogAttr(ctx context.Context, key string, value any) context.Context {
	attrs, ok := ctx.Value(logAttrsKey).([]slog.Attr)
	if !ok {
		attrs = []slog.Attr{}
	}
	attrs = append(attrs, slog.Any(key, value))
	return context.WithValue(ctx, logAttrsKey, attrs)
}

// GetLogAttrs retrieves logging attributes slice from the context.
func GetLogAttrs(ctx context.Context) []any {
	attrs, ok := ctx.Value(logAttrsKey).([]slog.Attr)
	if !ok {
		return []any{}
	}

	anyattr := make([]any, len(attrs))
	for i, attr := range attrs {
		anyattr[i] = attr
	}
	return anyattr
}

func InitLogger(environ string) error {
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
		return fmt.Errorf("Fail at init zap logger %s", err)
	}

	handler := zapslog.NewHandler(logger.Core())

	slog.SetDefault(slog.New(handler))
	return nil
}
