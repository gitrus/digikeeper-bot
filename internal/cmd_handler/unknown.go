package cmdhandler

import (
	"github.com/mymmrac/telego"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleUnknownCommand(bot *telego.Bot, update telego.Update) {
	chatId := tu.ID(update.Message.Chat.ID)
	_, _ = bot.SendMessage(tu.Message(
		chatId,
		"Unknown command. Type /help to see available commands.",
	))
}
