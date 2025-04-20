package telegomiddleware

import (
	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

// FirstNRunes returns the first n runes of a string
func FirstNRunes(s string, n int) string {
	runes := []rune(s)
	if n > len(runes) {
		n = len(runes)
	}
	return string(runes[:n])
}

func AddSlogAttrs() th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		origCtx := ctx.Context()

		innerCtx := loggingctx.AddLogAttr(origCtx, "update_id", update.UpdateID)

		if update.Message == nil {
			ctx = ctx.WithContext(innerCtx)
			return ctx.Next(update)
		}

		// Add more attributes for messages
		innerCtx = loggingctx.AddLogAttr(innerCtx, "message_id", update.Message.MessageID)
		innerCtx = loggingctx.AddLogAttr(innerCtx, "text_first10", FirstNRunes(update.Message.Text, 10))

		if update.Message.From != nil {
			innerCtx = loggingctx.AddLogAttr(innerCtx, "user_id", update.Message.From.ID)
		}

		innerCtx = loggingctx.AddLogAttr(innerCtx, "chat_id", update.Message.Chat.ID)

		ctx = ctx.WithContext(innerCtx)

		return ctx.Next(update)
	}
}
