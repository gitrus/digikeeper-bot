package cmdhandler

import (
	"log/slog"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
	tu "github.com/mymmrac/telego/telegoutil"

	session "github.com/gitrus/digikeeper-bot/pkg/sessionmanager"
)

func HandleAdd(usm session.UserSessionManager[*session.SimpleUserSession]) th.Handler {
	return func(ctx *th.Context, update telego.Update) error {
		slog.InfoContext(ctx.Context(), "Receive /add")

		userID := update.Message.From.ID
		state, err := usm.Fetch(ctx, userID)
		if err != nil {
			return err
		}

		_, err = usm.Set(
			ctx,
			userID,
			&session.SimpleUserSession{UserID: userID, State: "add", Version: state.Version + 1},
			state.Version,
		)
		if err != nil {
			slog.ErrorContext(ctx.Context(), "Failed to set state")

			chatId := tu.ID(update.Message.Chat.ID)
			_, err = ctx.Bot().SendMessage(ctx, tu.Message(
				chatId,
				"Another action is in progress. Please finish it first.",
			))
			return err
		}

		slog.InfoContext(ctx.Context(), "Set state", slog.String("state", state.State))
		return nil
	}
}
