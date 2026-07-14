package app

import (
	"os"
)

func (a *App) runScript() error {

	if len(os.Args) != 3 || os.Args[1] != "--script" {
		return nil
	}

	if err := a.service.ExecuteScript(os.Args[2]); err != nil {
		return err
	}

	return nil
}
