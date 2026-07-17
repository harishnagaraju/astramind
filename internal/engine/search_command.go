package engine

import (
	"fmt"
	"strings"

	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/features/search"
	"github.com/harishnagaraju/astramind/internal/infrastructure/renderer"
)

type searchCommand struct{}

func (c *searchCommand) Execute(
	app *App,
	input string,
) (bool, error) {

	searchService := search.NewService(history.NewService())

	if strings.HasPrefix(input, "/search ") || input == "/search" {

		query := strings.TrimSpace(
			strings.TrimPrefix(input, "/search"),
		)

		if query == "" {
			fmt.Println("Usage: /search <text>")
			return true, nil
		}

		results := searchService.SearchCurrent(
			app.runtime.Conversation,
			query,
		)

		if len(results) == 0 {
			fmt.Println("No matches found.")
			return true, nil
		}

		renderer.RenderSearchResults(results)

		return true, nil
	}

	if strings.HasPrefix(input, "/searchall ") || input == "/searchall" {

		query := strings.TrimSpace(
			strings.TrimPrefix(input, "/searchall"),
		)

		if query == "" {
			fmt.Println("Usage: /searchall <text>")
			return true, nil
		}

		results, err := searchService.SearchAll(query)
		if err != nil {
			return true, err
		}

		if len(results) == 0 {
			fmt.Println("No matches found.")
			return true, nil
		}

		renderer.RenderSessionSearchResults(results)

		return true, nil
	}

	return false, nil
}
