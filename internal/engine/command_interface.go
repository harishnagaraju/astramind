package engine

// Command represents a CLI command handler.
type Command interface {
	Execute(app *App, input string) (bool, error)
}
