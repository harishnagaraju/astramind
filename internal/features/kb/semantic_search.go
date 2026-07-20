package kb

import "sort"

// DefaultSemanticSearchLimit caps how many chunks SemanticSearch
// returns. Without a limit, every embedded chunk in the entire
// knowledge base gets returned regardless of relevance - harmless
// with a handful of test documents, but structurally wrong at real
// scale: a RAG prompt would end up containing the entire knowledge
// base on every single question, diluting relevant content with
// irrelevant chunks and eventually exceeding the model's context
// window outright.
const DefaultSemanticSearchLimit = 5

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
// the given query embedding, highest first, and returns at most
// DefaultSemanticSearchLimit results. Chunks with no stored embedding
// (imported before embeddings were enabled, or where embedding
// failed at import time) are skipped rather than scored as a false
// match.
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

	if len(results) > DefaultSemanticSearchLimit {
		results = results[:DefaultSemanticSearchLimit]
	}

	return results, nil
}
