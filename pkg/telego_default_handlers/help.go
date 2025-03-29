package cmdhandler

import (
	"log/slog"
	"strings"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleHelpFabric(cmdDescriptions map[string]string) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		logAttrs := loggingctx.GetLogAttrs(update.Context())
		slog.InfoContext(update.Context(), "Receive /help", logAttrs...)

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
		_, err := ctx.Bot().SendMessage(ctx, tu.Message(chatId, helpMessageBuilder.String()))

		return err
	}
}
