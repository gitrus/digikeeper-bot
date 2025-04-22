package loggingctx_test

import (
	"context"
	"log"
	"log/slog"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
)

func Example() {
	ctx := context.Background()

	ctx = loggingctx.AddLogAttr(ctx, "user_id", 345)
	ctx = loggingctx.AddLogAttr(ctx, "request_id", "12345")
	ctx = loggingctx.AddLogAttrs(
		ctx,
		[]slog.Attr{slog.Int("user_id", 346), slog.String("user_name", "John Doe")},
	)

	logger, err := loggingctx.InitLogger("dev")
	if err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	slog.SetDefault(logger)

	// Later in the request flow, retrieve all attributes for logging
	slog.InfoContext(ctx, "User action completed", loggingctx.GetLogAttrs(ctx)...)
}
