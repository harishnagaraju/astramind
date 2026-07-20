package engine

import (
	"fmt"
)

type exportCommand struct{}

func (c *exportCommand) Execute(
	app *App,
	input string,
) (bool, error) {

	exportService := app.deps.ExportService

	switch input {

	case "/export", "/export txt":

		if len(app.runtime.Conversation) == 0 {
			fmt.Println("Nothing to export.")
			return true, nil
		}

		err := exportService.ExportTXT(
			app.activeSession,
			app.runtime.Conversation,
		)

		if err != nil {
			return true, err
		}

		fmt.Printf(
			"Session exported to exports/%s.txt\n",
			app.activeSession,
		)

		return true, nil

	case "/export md":

		if len(app.runtime.Conversation) == 0 {
			fmt.Println("Nothing to export.")
			return true, nil
		}

		err := exportService.ExportMarkdown(
			app.activeSession,
			app.runtime.Conversation,
		)

		if err != nil {
			return true, err
		}

		fmt.Printf(
			"Session exported to exports/%s.md\n",
			app.activeSession,
		)

		return true, nil
	}

	return false, nil
}
