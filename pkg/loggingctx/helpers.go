// Package loggingctx provides structured logging support through context propagation.
//
// This package includes:
// - adding slog attributes to a context
// - retrieving slog attributes from a context
// - initializing zap logger for slog
//
// Examples:
//
//	// Add attributes to context
//	ctx = loggingctx.AddLogAttr(ctx, "user_id", userId)
//	ctx = loggingctx.AddLogAttr(ctx, "request_id", requestId)
//
//	// Initialize zap logger for slog
//	logger := loggingctx.InitZapLogger()
//	slog.SetDefault(logger)
//
//	// Later in the request flow, retrieve all attributes for logging
//	slog.InfoContext(ctx, "User action completed", loggingctx.GetLogAttrs(ctx)...)
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
	return AddLogAttrs(ctx, []slog.Attr{slog.Any(key, value)})
}

// AddLogAttr adds a logging attributes list to the provided context
func AddLogAttrs(ctx context.Context, attrs []slog.Attr) context.Context {
	logAttrs, ok := ctx.Value(LogAttrsKey).([]slog.Attr)
	if !ok {
		logAttrs = make([]slog.Attr, 0, 9)
	}

	existAttrs := make(map[string]slog.Attr, len(logAttrs))
	for _, attr := range logAttrs {
		existAttrs[attr.Key] = attr
	}

	newAttrs := make([]slog.Attr, 0, len(attrs))
	for _, attr := range attrs {
		if _, ok := existAttrs[attr.Key]; ok {
			existAttrs[attr.Key] = attr
		}
		newAttrs = append(newAttrs, attr)
	}

	result := make([]slog.Attr, 0, len(existAttrs)+len(newAttrs))
	for _, attr := range existAttrs {
		result = append(result, attr)
	}
	if len(newAttrs) > 0 {
		result = append(result, newAttrs...)
	}

	return context.WithValue(ctx, LogAttrsKey, result)
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
	// environ in [dev, <any-else>]
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
