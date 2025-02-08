package main

import (
	"log"

	ic "github.com/WAY29/icecream-go/icecream"
	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"

	cmdh "github.com/gitrus/digikeeper-bot/internal/cmd_handler"
	cmdhandler "github.com/gitrus/digikeeper-bot/internal/cmd_handler"
	telegomiddleware "github.com/gitrus/digikeeper-bot/pkg/telego_middleware"
)

func main() {
	config := Configure()

	bot, err := telego.NewBot(config.BotKey, telego.WithDefaultDebugLogger())
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}

	updates, err := bot.UpdatesViaLongPolling(nil)
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	defer bot.StopLongPolling()

	bh, err := th.NewBotHandler(bot, updates)
	if err != nil {
		log.Fatalf("Failed to start bot: %v", err)
	}
	defer bh.Stop()

	// Add global middleware, it will be applied in order of addition
	bh.Use(th.PanicRecovery())
	bh.Use(th.Timeout(config.CommonTimeout))

	bh.Use(telegomiddleware.SlogAddAttrs())

	usm := telegomiddleware.NewMockUserStateManager[string]()
	useStateMiddleware := telegomiddleware.NewUserStateMiddleware(usm)
	bh.Use(useStateMiddleware.Middleware())

	cmdHandlerGroup := cmdh.NewCommandHandlerGroup(usm)
	cmdHandlerGroup.RegisterCommand("start", cmdhandler.HandleStart, "Show start-bot message")
	cmdHandlerGroup.RegisterCommand("cancel", cmdhandler.HandleCancel(usm), "Interrupt any current operation/s")
	cmdHandlerGroup.RegisterCommand("add", cmdhandler.HandleAdd(usm), "Add new note to the list")

	cmdHandlerGroup.RegisterGroup(bh)

	ic.Ic("basegroup", bh.BaseGroup())
	bh.Start()
}
