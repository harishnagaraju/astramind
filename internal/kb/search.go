package kb

import "strings"

// SearchResult represents a keyword search result.
type SearchResult struct {
	DocumentID string
	ChunkID    string
	ChunkIndex int
	Content    string
	Score      int
}

// Search performs a case-insensitive keyword search across all chunks.
func (r *Repository) Search(query string) ([]SearchResult, error) {
	query = strings.TrimSpace(query)
	if query == "" {
		return []SearchResult{}, nil
	}

	chunks, err := r.Chunks()
	if err != nil {
		return nil, err
	}

	query = strings.ToLower(query)

	var results []SearchResult

	for _, chunk := range chunks {
		content := strings.ToLower(chunk.Content)

		score := strings.Count(content, query)
		if score == 0 {
			continue
		}

		results = append(results, SearchResult{
			DocumentID: chunk.DocumentID,
			ChunkID:    chunk.ID,
			ChunkIndex: chunk.Index,
			Content:    chunk.Content,
			Score:      score,
		})
	}

	return results, nil
}
