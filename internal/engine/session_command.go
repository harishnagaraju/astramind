package engine

import (
	"fmt"
	"strings"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

type sessionCommand struct{}

func (c *sessionCommand) Execute(app *App, input string) (bool, error) {

	// Create once for this command execution.
	svc := app.deps.SessionService

	// /load
	if strings.HasPrefix(input, "/load ") {

		sessionName := strings.TrimSpace(
			strings.TrimPrefix(input, "/load "),
		)

		if sessionName == "" {
			fmt.Println("Usage: /load <session-name>")
			return true, nil
		}

		if !svc.Exists(sessionName) {
			fmt.Printf("Session '%s' does not exist.\n", sessionName)
			return true, nil
		}

		messages, err := svc.Load(sessionName)
		if err != nil {
			return true, err
		}

		app.activeSession = sessionName
		app.runtime.Conversation = messages

		fmt.Printf("Loaded session: %s\n", sessionName)

		return true, nil
	}

	// /new
	if strings.HasPrefix(input, "/new ") {

		sessionName := strings.TrimSpace(
			strings.TrimPrefix(input, "/new "),
		)

		if sessionName == "" {
			fmt.Println("Usage: /new <session-name>")
			return true, nil
		}

		if err := svc.Create(sessionName); err != nil {
			return true, err
		}

		app.activeSession = sessionName
		app.runtime.Conversation = []models.Message{}

		fmt.Printf("Created and switched to session: %s\n", sessionName)

		return true, nil
	}

	// /delete
	if strings.HasPrefix(input, "/delete ") {

		sessionName := strings.TrimSpace(
			strings.TrimPrefix(input, "/delete "),
		)

		if sessionName == "" {
			fmt.Println("Usage: /delete <session-name>")
			return true, nil
		}

		if sessionName == "default" {
			fmt.Println("Default session cannot be deleted.")
			return true, nil
		}

		if sessionName == app.activeSession {
			fmt.Println("Cannot delete active session.")
			return true, nil
		}

		if err := svc.Delete(sessionName); err != nil {
			return true, err
		}

		fmt.Printf("Deleted session: %s\n", sessionName)

		return true, nil
	}

	// /sessions
	if input == "/sessions" {

		sessions, err := svc.List()
		if err != nil {
			return true, err
		}

		fmt.Println()
		fmt.Println("Available Sessions")
		fmt.Println("------------------")

		if len(sessions) == 0 {
			fmt.Println("No sessions found.")
			return true, nil
		}

		for _, session := range sessions {
			fmt.Println(session)
		}

		return true, nil
	}

	return false, nil
}
