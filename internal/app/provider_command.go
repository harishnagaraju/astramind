package app

import "fmt"

type providerCommand struct {
	app *App
}

func (c *providerCommand) Name() string {
	return "/provider"
}

func (c *providerCommand) Execute(input string) (bool, error) {

	if input != c.Name() {
		return false, nil
	}

	fmt.Println()
	fmt.Println("Current AI Provider")
	fmt.Println("-------------------")

	fmt.Printf(
		"Provider : %s\n",
		c.app.deps.ProviderManager.ProviderName(),
	)

	fmt.Printf(
		"Model    : %s\n",
		c.app.model,
	)

	fmt.Println()

	return true, nil
}
