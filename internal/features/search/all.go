package search

import (
	"github.com/harishnagaraju/astramind/internal/features/history"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// SearchAllSessions searches every saved session for the specified query.
// It depends on the History feature, not on storage directly.
func SearchAllSessions(
	historySvc *history.Service,
	query string,
) ([]models.SessionSearchResult, error) {

	sessionNames, err := historySvc.ListSessions()
	if err != nil {
		return nil, err
	}

	results := make([]models.SessionSearchResult, 0)

	for _, session := range sessionNames {

		messages, err := historySvc.Load(session)
		if err != nil {
			continue
		}

		matches := SearchMessages(messages, query)

		for _, match := range matches {
			results = append(results, models.SessionSearchResult{
				Session: session,
				Index:   match.Index,
				Role:    match.Role,
				Content: match.Content,
			})
		}
	}

	return results, nil
}
