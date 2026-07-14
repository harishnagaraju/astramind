package app

// Command represents a CLI command.
type Command interface {
	Name() string
	Execute(input string) (bool, error)
}
