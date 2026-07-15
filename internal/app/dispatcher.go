package app

// commandDispatcher routes CLI commands to registered handlers.
type commandDispatcher struct {
	app *App

	commands []Command
}

// newDispatcher creates a new command dispatcher and registers
// all supported commands.
func newDispatcher(app *App) *commandDispatcher {

	d := &commandDispatcher{
		app: app,
	}

	// Register commands here.
	d.commands = []Command{
		&builtinCommand{},
		&sessionCommand{},
		&conversationCommand{},
		&searchCommand{},
	}

	return d
}

// Execute routes a command to the first registered handler
// that can process it.
func (d *commandDispatcher) Execute(input string) (bool, error) {

	for _, cmd := range d.commands {

		handled, err := cmd.Execute(d.app, input)

		if handled || err != nil {
			return handled, err
		}
	}

	// Fall back to existing command handlers.
	return d.app.deps.ChatService.HandleKnowledgeCommand(input)
}
