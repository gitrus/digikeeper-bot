package telegomiddleware

import (
	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func firstNRunes(s string, n int) string {
	runes := []rune(s)
	if n > len(runes) {
		n = len(runes)
	}
	return string(runes[:n])
}

func SlogAddAttrs() th.Middleware {
	return func(bot *telego.Bot, update telego.Update, next th.Handler) {
		ctx := update.Context()
		ctx = loggingctx.AddLogAttr(ctx, "update_id", update.UpdateID)
		ctx = loggingctx.AddLogAttr(ctx, "user_id", update.Message.From.ID)
		ctx = loggingctx.AddLogAttr(ctx, "chat_id", update.Message.Chat.ID)
		ctx = loggingctx.AddLogAttr(ctx, "message_id", update.Message.MessageID)
		ctx = loggingctx.AddLogAttr(ctx, "text_first10", firstNRunes(update.Message.Text, 10))

		next(bot, update.WithContext(ctx))
	}
}
