package renderer

import (
	"fmt"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// RenderSearchResults displays search results for the current session.
func RenderSearchResults(results []models.SearchResult) {

	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	fmt.Printf("Found %d match(es):\n\n", len(results))

	for _, result := range results {

		fmt.Printf("[%d] %s\n", result.Index+1, result.Role)
		fmt.Println(result.Content)
		fmt.Println()
	}
}

// RenderSessionSearchResults displays search results across sessions.
func RenderSessionSearchResults(results []models.SessionSearchResult) {

	if len(results) == 0 {
		fmt.Println("No matches found.")
		return
	}

	sessionMap := make(map[string]struct{})

	for _, result := range results {
		sessionMap[result.Session] = struct{}{}
	}

	fmt.Printf(
		"Found %d match(es) across %d session(s).\n\n",
		len(results),
		len(sessionMap),
	)

	currentSession := ""

	for _, result := range results {

		if result.Session != currentSession {

			currentSession = result.Session

			fmt.Println("==================================================")
			fmt.Printf("Session: %s\n", currentSession)
			fmt.Println("==================================================")
			fmt.Println()
		}

		fmt.Printf("[%d] %s\n", result.Index+1, result.Role)
		fmt.Println(result.Content)
		fmt.Println()
	}
}
