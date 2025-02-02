package cmdhandler

import (
	"log/slog"

	ic "github.com/WAY29/icecream-go/icecream"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleCancel(bot *telego.Bot, update telego.Update) {
	ic.Ic("Command /cancel received")

	slog.InfoContext(update.Context(), "Received update", slog.Any("message", update.Message))

	chatId := tu.ID(update.Message.Chat.ID)
	_, _ = bot.SendMessage(tu.Message(
		chatId,
		"I just interrupted the current operation/s. What can I do for you now?",
	))
}
