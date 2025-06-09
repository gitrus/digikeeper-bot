package cmdhandler

import (
	"log/slog"
	"strings"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/gitrus/digikeeper-bot/internal/note"
	"github.com/gitrus/digikeeper-bot/pkg/noteparser"
)

func HandleAddNote(svc note.Service) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		slog.InfoContext(ctx.Context(), "Receive /addnote")

		chatID := tu.ID(update.Message.Chat.ID)

		raw := strings.TrimSpace(strings.TrimPrefix(update.Message.Text, "/addnote"))
		createdAt := time.Unix(int64(update.Message.Date), 0)
		parsedNote, _ := noteparser.Parse(createdAt, raw)

		svc.SetPending(update.Message.From.ID, parsedNote)

		msg := "Parsed note:\n" + parsedNote.Payload.Text
		if len(parsedNote.Tags) > 0 {
			msg += "\nAdd tags?"
		}

		keyboard := tu.InlineKeyboard()
		for _, tag := range parsedNote.Tags {
			keyboard.Row(tu.InlineKeyboardButton("+" + tag).WithCallbackData("addtag:" + tag))
		}
		keyboard.Row(tu.InlineKeyboardButton("Save").WithCallbackData("save"))

		_, err := ctx.Bot().SendMessage(ctx, tu.Message(chatID, msg).WithReplyMarkup(keyboard))
		if err != nil {
			slog.ErrorContext(ctx.Context(), "Failed to send message", "error", err)
			return err
		}
		return nil
	}
}
