package cmdhandler

import (
	ic "github.com/WAY29/icecream-go/icecream"
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleStart(bot *telego.Bot, update telego.Update) {
	ic.Ic("Command /start received")

	chatId := tu.ID(update.Message.Chat.ID)
	_, _ = bot.SendMessage(tu.Message(
		chatId,
		"Hello! I'm a digikeeper bot. I can help you to keep your digital life in order.",
	))
}
