package telegocommandrouter

import (
	"log/slog"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"
)

func HandleUnknownCommand(ctx *th.Context, update telego.Update) error {
	cmd := ""
	if update.Message != nil && update.Message.Text != "" {
		cmd = update.Message.Text
	}
	slog.InfoContext(
		ctx.Context(),
		"Unknown command", append(loggingctx.GetLogAttrs(ctx.Context()), slog.String("command", cmd))...)

	chatId := tu.ID(update.Message.Chat.ID)
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(
		chatId,
		"Unknown command. Type /help to see available commands.",
	))

	return err
}
