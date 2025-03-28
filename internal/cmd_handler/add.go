package cmdhandler

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
)

func HandleAdd(usm UserStateManager) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		slog.InfoContext(update.Context(), "Receive /add", loggingctx.GetLogAttrs(update.Context())...)

		userID := update.Message.From.ID
		state, err := usm.Set(userID, "add")
		if err != nil {
			slog.ErrorContext(update.Context(), "Failed to set state", loggingctx.GetLogAttrs(update.Context())...)

			chatId := tu.ID(update.Message.Chat.ID)
			_, err = ctx.Bot().SendMessage(ctx, tu.Message(
				chatId,
				"Another action is in progress. Please finish it first.",
			))
			return err
		}

		logAttrs := append(loggingctx.GetLogAttrs(update.Context()), slog.String("state", state))
		slog.InfoContext(update.Context(), "Set state", logAttrs...)
		return nil
	}
}
