package app

import (
	"fmt"

	"github.com/harishnagaraju/astramind/internal/config"
)

type aboutCommand struct {
	app *App
}

func (c *aboutCommand) Name() string {
	return "/about"
}

func (c *aboutCommand) Execute(input string) (bool, error) {

	if input != c.Name() {
		return false, nil
	}

	fmt.Println("\nAstraMind")
	fmt.Println("---------")

	fmt.Printf(
		"Version: %s\n",
		config.Version,
	)

	fmt.Println("\nFeatures:")

	fmt.Println("✓ Conversation Memory")
	fmt.Println("✓ Persistent History")
	fmt.Println("✓ Session Statistics")
	fmt.Println("✓ Configuration Display")

	fmt.Printf(
		"\nModel: %s\n",
		c.app.model,
	)

	fmt.Println("Author: Harish Nagaraju")
	fmt.Println("Company: RK Consulting")
	fmt.Println("Repository: github.com/harishnagaraju/astramind")

	return true, nil
}
