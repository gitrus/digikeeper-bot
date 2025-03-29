package main

import (
	"context"
	"log"
	"log/slog"

	th "github.com/mymmrac/telego/telegohandler"

	cmdh "github.com/gitrus/digikeeper-bot/internal/cmd_handler"
	cmdhandler "github.com/gitrus/digikeeper-bot/internal/cmd_handler"
	telegomiddleware "github.com/gitrus/digikeeper-bot/pkg/telego_middleware"
)

func main() {
	config := configure()
	logger := slog.Default()

	ctx := context.Background()

	bot, updates, err := initBot(ctx, config)

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	defer bh.Stop()

	// Add global middleware, it will be applied in order of addition
	bh.Use(th.PanicRecovery())
	bh.Use(th.Timeout(config.CommonTimeout))

	bh.Use(telegomiddleware.AddSlogAttrs())

	usm := telegomiddleware.NewMockUserStateManager[string]()
	useStateMiddleware := telegomiddleware.NewUserStateMiddleware(usm)
	bh.Use(useStateMiddleware.Middleware())

	cmdHandlerGroup := cmdh.NewCommandHandlerGroup(usm)
	cmdHandlerGroup.RegisterCommand("start", cmdhandler.HandleStart, "Show start-bot message")
	cmdHandlerGroup.RegisterCommand("cancel", cmdhandler.HandleCancel(usm), "Interrupt any current operation/s")
	cmdHandlerGroup.RegisterCommand("add", cmdhandler.HandleAdd(usm), "Add new note to the list")

	cmdHandlerGroup.RegisterGroup(bh)

	logger.Info("CmdHandlerGroup", "group", cmdHandlerGroup)

	logger.Info("Starting bot ...")
	bh.Start()
}
