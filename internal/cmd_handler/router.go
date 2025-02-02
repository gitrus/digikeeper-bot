package cmdhandler

import (
	th "github.com/mymmrac/telego/telegohandler"
)

type registeredHandler struct {
	h th.Handler
	d string
}

type CommandHandlerGroup struct {
	commands map[string]registeredHandler
}

func NewCommandHandlerGroup() *CommandHandlerGroup {
	chg := &CommandHandlerGroup{
		commands: make(map[string]registeredHandler),
	}

	chg.commands["start"] = registeredHandler{h: HandleStart, d: "Show start-bot message"}
	chg.commands["cancel"] = registeredHandler{h: HandleCancel, d: "Interrupt any current operation/s"}

	chg.commands["help"] = registeredHandler{
		h: HandleHelpFabric(chg.getCommandToDescription()),
		d: "Show help message with available commands",
	}

	return chg
}

func (ch *CommandHandlerGroup) getCommandToDescription() map[string]string {
	cmdToDesc := make(map[string]string)
	for command, handler := range ch.commands {
		cmdToDesc[command] = handler.d
	}
	return cmdToDesc
}

func (ch *CommandHandlerGroup) RegisterGroup(bh *th.BotHandler) {
	commands := bh.Group(th.AnyCommand())

	predicats := make([]th.Predicate, 0, len(ch.commands))
	for command, handler := range ch.commands {
		p := th.CommandEqual(command)
		commands.Handle(handler.h, p)
		predicats = append(predicats, p)
	}

	// handle unknown command
	var unknownPredicat th.Predicate = th.None()
	for _, predicat := range predicats {
		unknownPredicat = th.Or(unknownPredicat, th.Not(predicat))
	}
	commands.Handle(HandleUnknownCommand, unknownPredicat)
}
