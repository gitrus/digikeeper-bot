package cmdhandler

import (
	"strings"

	ic "github.com/WAY29/icecream-go/icecream"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleHelpFabric(cmdDescriptions map[string]string) th.Handler {
	return func(bot *telego.Bot, update telego.Update) {
		ic.Ic("Command /help received")

		helpMessageBuilder := strings.Builder{}
		helpMessageBuilder.WriteString("/help    Show this help message\n")
		for command, description := range cmdDescriptions {
			if command == "help" {
				continue
			}
			helpMessageBuilder.WriteString("/" + command)
			helpMessageBuilder.WriteString("   " + description)
			helpMessageBuilder.WriteString("\n")
		}

		chatId := tu.ID(update.Message.Chat.ID)
		_, _ = bot.SendMessage(tu.Message(chatId, helpMessageBuilder.String()))
	}
}
