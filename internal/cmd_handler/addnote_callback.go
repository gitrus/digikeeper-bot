package cmdhandler

import (
	"log/slog"
	"strings"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	"github.com/gitrus/digikeeper-bot/internal/note"
)

func HandleAddNoteCallback(svc note.Service) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		if update.CallbackQuery == nil {
			return nil
		}

		data := update.CallbackQuery.Data
		userID := update.CallbackQuery.From.ID
		chatID := tu.ID(update.CallbackQuery.Message.Chat.ID)
		msgID := update.CallbackQuery.Message.MessageID

		if strings.HasPrefix(data, "addtag:") {
			tag := strings.TrimPrefix(data, "addtag:")
			svc.AddTagToPending(userID, tag)
			text := update.CallbackQuery.Message.Text + "\nTag added: " + tag
			_, err := ctx.Bot().EditMessageText(ctx, tu.EditMessageText(chatID, msgID, text))
			if err != nil {
				slog.ErrorContext(ctx.Context(), "Failed to edit message", "error", err)
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(update.CallbackQuery.ID))
			return err
		}

		if data == "save" {
			if err := svc.SavePending(ctx.Context(), userID); err != nil {
				slog.ErrorContext(ctx.Context(), "Failed to save note", "error", err)
				return err
			}
			_, err := ctx.Bot().EditMessageText(ctx, tu.EditMessageText(chatID, msgID, "Note saved"))
			if err != nil {
				slog.ErrorContext(ctx.Context(), "Failed to edit message", "error", err)
				return err
			}
			_ = ctx.Bot().AnswerCallbackQuery(ctx, tu.CallbackQuery(update.CallbackQuery.ID))
			return nil
		}
		return nil
	}
}
