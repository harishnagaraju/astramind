package engine

import "github.com/harishnagaraju/astramind/internal/models"

import (
	"fmt"
	"os"
)

var conversation []models.Message

// Run starts the AstraMind application.
//
// The implementation will be migrated incrementally from
// cmd/astramind/main.go during the Bootstrap milestone.
func (a *App) Run() error {

	if err := a.initialize(); err != nil {
		fmt.Println(err)
		return nil
	}

	// Script execution mode.
	if err := a.runScript(); err != nil {
		return err
	}

	if len(os.Args) == 3 && os.Args[1] == "--script" {
		return nil
	}

	return a.runInteractive()
}
