package kb

import "sort"

// SemanticSearchResult represents an embedding-based search result,
// ranked by cosine similarity rather than keyword count.
type SemanticSearchResult struct {
	DocumentID string
	ChunkID    string
	ChunkIndex int
	Content    string
	Score      float64
}

// SemanticSearch ranks every embedded chunk by cosine similarity to
// the given query embedding, highest first. Chunks with no stored
// embedding (imported before embeddings were enabled, or where
// embedding failed at import time) are skipped rather than scored as
// a false match.
func (r *Repository) SemanticSearch(
	queryEmbedding []float32,
) ([]SemanticSearchResult, error) {

	if len(queryEmbedding) == 0 {
		return []SemanticSearchResult{}, nil
	}

	chunks, err := r.Chunks()
	if err != nil {
		return nil, err
	}

	var results []SemanticSearchResult

	for _, chunk := range chunks {

		if len(chunk.Embedding) == 0 {
			continue
		}

		score := CosineSimilarity(queryEmbedding, chunk.Embedding)

		results = append(results, SemanticSearchResult{
			DocumentID: chunk.DocumentID,
			ChunkID:    chunk.ID,
			ChunkIndex: chunk.Index,
			Content:    chunk.Content,
			Score:      score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results, nil
}
