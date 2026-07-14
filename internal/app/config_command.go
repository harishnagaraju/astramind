package app

import (
	"fmt"

	"github.com/harishnagaraju/astramind/internal/config"
)

type configCommand struct {
	app *App
}

func (c *configCommand) Name() string {
	return "/config"
}

func (c *configCommand) Execute(input string) (bool, error) {

	if input != c.Name() {
		return false, nil
	}

	fmt.Println("\nCurrent Configuration")
	fmt.Println("---------------------")

	fmt.Printf(
		"Model: %s\n",
		c.app.model,
	)

	fmt.Printf(
		"Base URL: %s\n",
		c.app.baseURL,
	)

	fmt.Printf(
		"Max Messages: %d\n",
		config.MaxMessages,
	)

	fmt.Printf(
		"History Enabled: %t\n",
		true,
	)

	fmt.Printf(
		"History File: %s\n",
		config.HistoryFile,
	)

	return true, nil
}
