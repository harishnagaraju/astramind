package search

import (
	"strings"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// SearchMessages performs a case-insensitive search over a conversation
// already loaded into memory.
func SearchMessages(messages []models.Message, query string) []models.SearchResult {
	query = strings.TrimSpace(query)
	if query == "" {
		return nil
	}

	query = strings.ToLower(query)

	results := make([]models.SearchResult, 0)

	for i, msg := range messages {
		if strings.Contains(strings.ToLower(msg.Content), query) {
			results = append(results, models.SearchResult{
				Index:   i,
				Role:    msg.Role,
				Content: msg.Content,
			})
		}
	}

	return results
}
