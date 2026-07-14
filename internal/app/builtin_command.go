package app

import (
	"fmt"

	"github.com/harishnagaraju/astramind/internal/config"
)

type builtinCommand struct {
	app *App
}

func (c *builtinCommand) Name() string {
	return "builtin"
}

func (c *builtinCommand) Execute(input string) (bool, error) {

	switch input {

	case "/help":

		fmt.Println("\nAvailable Commands:")
		fmt.Println("/about     - About AstraMind")
		fmt.Println("/help      - Show help")
		fmt.Println("/history   - Show conversation history")
		fmt.Println("/clear     - Clear conversation memory")
		fmt.Println("/stats     - Show session statistics")
		fmt.Println("/config    - Show configuration")
		fmt.Println("/export    - Export session (TXT)")
		fmt.Println("/export md - Export session (Markdown)")
		fmt.Println("/sessions  - List sessions")
		fmt.Println("/search <text> - Search current conversation")
		fmt.Println("/searchall <text> - Search all conversation")

		fmt.Println()

		fmt.Println("Knowledge Base")
		fmt.Println("/kb import <file> - Import a text or markdown document")
		fmt.Println("/kb list          - List imported documents")
		fmt.Println("/kb search <text> - Search the knowledge base")
		fmt.Println("/kb remove <id>   - Remove a document")
		fmt.Println("/kb clear         - Remove all documents")
		fmt.Println("/kb stats         - Show knowledge base statistics")

		fmt.Println()

		fmt.Println("/new <name> - Create session")
		fmt.Println("/load <name> - Load session")
		fmt.Println("/delete <name> - Delete session")
		fmt.Println("/export    - Export session")
		fmt.Println("/provider  - Show active AI provider")
		fmt.Println("exit       - Exit AstraMind")
		fmt.Println("quit       - Exit AstraMind")

		return true, nil

	case "/about":

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

	case "/config":

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

	case "/provider":

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

	default:
		return false, nil
	}
}
