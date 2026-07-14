package app

import "fmt"

type helpCommand struct {
	app *App
}

func (c *helpCommand) Name() string {
	return "/help"
}

func (c *helpCommand) Execute(
	input string,
) (bool, error) {

	if input != c.Name() {
		return false, nil
	}

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
}
