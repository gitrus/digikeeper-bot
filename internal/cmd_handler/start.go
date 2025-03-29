package cmdhandler

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/gitrus/digikeeper-bot/pkg/loggingctx"
)

var startKeyboard = tu.Keyboard(
	tu.KeyboardRow(
		tu.KeyboardButton("/add"),
		tu.KeyboardButton("/search"),
	),
	tu.KeyboardRow(
		tu.KeyboardButton("More..."),
	),
).WithResizeKeyboard().WithInputFieldPlaceholder("Select action")

func HandleStart(ctx *th.Context, update telego.Update) error {
	chatId := tu.ID(update.Message.Chat.ID)
	_, err := ctx.Bot().SendMessage(ctx, tu.Message(
		chatId,
		"Hello! I'm a digikeeper bot. I can help you to keep your digital life in order.",
	).WithReplyMarkup(startKeyboard))
	if err != nil {
		return err
	}

	slog.InfoContext(ctx.Context(), "Receive /start", loggingctx.GetLogAttrs(ctx.Context())...)

	return nil
}
