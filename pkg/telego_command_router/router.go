package telegocommandrouter

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

	return chg
}

func (ch *CommandHandlerGroup) RegisterCommand(command string, handler th.Handler, description string) {
	ch.commands[command] = registeredHandler{h: handler, d: description}

	ch.commands["help"] = registeredHandler{
		h: HandleHelpFabric(ch.getCommandToDescription()),
		d: "Show help message with available commands",
	}
}

func (ch *CommandHandlerGroup) RegisterGroup(bh *th.BotHandler) {
	commands := bh.Group(th.AnyCommand())

	predicates := make([]th.Predicate, 0, len(ch.commands))
	for command, handler := range ch.commands {
		p := th.CommandEqual(command)
		commands.Handle(handler.h, p)
		predicates = append(predicates, p)
	}

	// handle unknown command
	var unknownPredicate th.Predicate = th.None()
	for _, predicate := range predicates {
		unknownPredicate = th.Or(unknownPredicate, th.Not(predicate))
	}
	commands.Handle(HandleUnknownCommand, unknownPredicate)
}

func (ch *CommandHandlerGroup) getCommandToDescription() map[string]string {
	cmdToDesc := make(map[string]string)
	for command, handler := range ch.commands {
		cmdToDesc[command] = handler.d
	}
	return cmdToDesc
}
