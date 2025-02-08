package cmdhandler

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
)

type UserStateManager interface {
	DropActiveState(userID int64)
	Set(userID int64, state string) (string, error)
}

func HandleCancel(usm UserStateManager) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		slog.InfoContext(update.Context(), "Receive /cancel", loggingctx.GetLogAttrs(update.Context())...)

		userID := update.Message.From.ID
		usm.DropActiveState(userID)

		chatId := tu.ID(update.Message.Chat.ID)
		_, _ = bot.SendMessage(tu.Message(
			chatId,
			"I just interrupted the current operation/s. What can I do for you now?",
		))
	}
}
