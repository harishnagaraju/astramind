package engine

import "github.com/harishnagaraju/astramind/internal/infrastructure/models"

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
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

	// Local web UI mode.
	if len(os.Args) >= 2 && os.Args[1] == "--web" {

		addr := "localhost:8420"

		go openBrowser("http://" + addr)

		return a.runWeb(addr)
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

// openBrowser opens the given URL in the system's default browser.
// Best-effort - if it fails, the person can still open the URL
// manually, so the error is not treated as fatal.
func openBrowser(url string) {

	var cmd *exec.Cmd

	switch runtime.GOOS {

	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)

	case "darwin":
		cmd = exec.Command("open", url)

	default:
		cmd = exec.Command("xdg-open", url)
	}

	_ = cmd.Start()
}
