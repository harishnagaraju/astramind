package kb

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Manager coordinates knowledge base operations.
type Manager struct {
	storage  Storage
	embedder Embedder
}

// ErrNoExtractiveMatch is returned by ExtractiveAnswer when no
// extractable item could be found or embedded from the supplied
// results.
var ErrNoExtractiveMatch = fmt.Errorf("no extractive match found")

// NewManager creates a new knowledge base manager.
func NewManager(storage Storage) *Manager {
	return &Manager{
		storage: storage,
	}
}

// SetEmbedder configures the embedder used to generate vector
// embeddings for chunks during import. It is optional - if never
// set, ImportDocument behaves exactly as before (keyword search
// only, no embeddings generated).
func (m *Manager) SetEmbedder(embedder Embedder) {
	m.embedder = embedder
}

func (m *Manager) SaveDocument(doc *Document) error {
	return m.storage.SaveDocument(doc)
}

func (m *Manager) SaveChunks(chunks []Chunk) error {
	return m.storage.SaveChunks(chunks)
}

func (m *Manager) LoadDocument(id string) (*Document, error) {
	return m.storage.LoadDocument(id)
}

func (m *Manager) LoadChunks(documentID string) ([]Chunk, error) {
	return m.storage.LoadChunks(documentID)
}

func (m *Manager) DeleteDocument(id string) error {
	return m.storage.DeleteDocument(id)
}

func (m *Manager) DeleteChunks(documentID string) error {
	return m.storage.DeleteChunks(documentID)
}

// ListKnowledge returns all knowledge base documents.
func (m *Manager) ListKnowledge() ([]Document, error) {
	return m.ListDocuments()
}

// GetKnowledge returns a knowledge base document by ID.
func (m *Manager) GetKnowledge(documentID string) (*Document, error) {
	return m.LoadDocument(documentID)
}

// RemoveKnowledge removes a document and its chunks.
func (m *Manager) RemoveKnowledge(documentID string) error {

	if err := m.DeleteChunks(documentID); err != nil {
		return err
	}

	if err := m.DeleteDocument(documentID); err != nil {
		return err
	}

	return nil
}

// ClearKnowledge removes every knowledge base document.
func (m *Manager) ClearKnowledge() error {

	documents, err := m.ListKnowledge()
	if err != nil {
		return err
	}

	for _, doc := range documents {
		if err := m.RemoveKnowledge(doc.ID); err != nil {
			return err
		}
	}

	return nil
}

func (m *Manager) ListDocuments() ([]Document, error) {
	return m.storage.ListDocuments()
}

// ImportDocument imports a text or markdown file into the knowledge base.
func (m *Manager) ImportDocument(path string) (*Document, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".txt", ".md":
		// supported
	default:
		return nil, ErrInvalidDocument
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	id := generateDocumentID()

	doc := &Document{
		ID:         id,
		Name:       filepath.Base(path),
		Path:       path,
		Content:    string(data),
		CreatedAt:  time.Now(),
		ModifiedAt: info.ModTime(),
	}

	// Split the document into chunks.
	chunks := ChunkDocument(doc, DefaultChunkSize, DefaultOverlap)

	// Record the number of generated chunks.
	doc.ChunkCount = len(chunks)

	// Generate embeddings for each chunk, if an embedder is
	// configured. A failure here does not fail the import - the
	// chunk simply falls back to keyword search, since Embedding
	// stays nil.
	if m.embedder == nil {
		fmt.Println("(no embedder configured - skipping embeddings)")
	} else {

		embedded := 0

		for i := range chunks {

			embedding, err := m.embedder.Embed(chunks[i].Content)
			if err != nil {
				fmt.Printf(
					"(embedding failed for chunk %d: %v)\n",
					i,
					err,
				)
				continue
			}

			chunks[i].Embedding = embedding
			embedded++
		}

		fmt.Printf(
			"(%d/%d chunks embedded)\n",
			embedded,
			len(chunks),
		)
	}

	// Persist the chunks.
	if err := m.SaveChunks(chunks); err != nil {
		return nil, err
	}

	// Persist the document.
	if err := m.SaveDocument(doc); err != nil {
		return nil, err
	}

	return doc, nil
}

func (m *Manager) Search(query string) ([]SearchResult, error) {

	repository := NewRepository(m)

	return repository.Search(query)
}

// SemanticSearch performs an embedding-based search of the knowledge
// base. Unlike ImportDocument's graceful degradation, this fails
// loudly if no embedder is configured - a silent empty result here
// would be indistinguishable from "no matches found".
func (m *Manager) SemanticSearch(query string) ([]SemanticSearchResult, error) {

	if m.embedder == nil {
		return nil, fmt.Errorf(
			"semantic search requires an embedder to be configured",
		)
	}

	queryEmbedding, err := m.embedder.Embed(query)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to embed query: %w",
			err,
		)
	}

	repository := NewRepository(m)

	return repository.SemanticSearch(queryEmbedding)
}

// ExtractiveAnswer answers a single-fact question by finding the
// single chunk whose content is most semantically similar to the
// question, then returning that chunk's full content verbatim - no
// LLM paraphrase step, so no opportunity for a generative model to
// misstate a date, fee, or threshold while rewording it.
//
// This ranks at the fine-grained sentence level internally (to pick
// the single best-matching CHUNK when several were retrieved by
// SemanticSearch), but returns the whole chunk, not a window around
// the matched sentence.
//
// A windowed version of this function was tried and measured against
// real embeddings before landing here. Real similarity data (captured
// via a temporary debug instrumentation, not guessed) showed that
// cosine similarity between short sentences from this embedding model
// does not reliably track topic relevance: an unrelated sentence
// ("Not meeting on 16 February" - a different class's schedule)
// scored HIGHER (0.59) against the anchor ("Meeting ID 795 777
// 3585") than genuinely on-topic sentences did ("Password OMpeace" at
// 0.43, the Zoom URL at 0.49) - apparently because both happen to
// contain the literal word "meeting", despite meaning something
// unrelated. No threshold value could separate on-topic from
// off-topic content using this signal: a threshold strict enough to
// exclude the false-positive "Not meeting..." sentence would also
// exclude the genuinely relevant "Password"/URL sentences. This
// wasn't a calibration problem - the technique itself doesn't carry
// a usable boundary signal at this granularity, with this embedding
// model, on this kind of short-sentence prose.
//
// Given that, whole-chunk return (verbose but always correct, same
// tradeoff already accepted for ExtractItems/BuildListAnswer on
// enumeration queries) is the safer default: chunking (see
// chunker.go) already guarantees no entry is corrupted or split
// mid-content, so returning a whole chunk can never silently include
// wrong information - only extra information, which is a strictly
// safer failure mode than a similarity threshold that was shown to
// sometimes rank irrelevant content above relevant content.
//
// Returns ErrNoExtractiveMatch if results is empty or nothing could
// be embedded.
func (m *Manager) ExtractiveAnswer(
	question string,
	results []SemanticSearchResult,
) (*ExtractedItem, error) {

	if m.embedder == nil {
		return nil, fmt.Errorf(
			"extractive answer requires an embedder to be configured",
		)
	}

	if len(results) == 0 {
		return nil, ErrNoExtractiveMatch
	}

	questionEmbedding, err := m.embedder.Embed(question)
	if err != nil {
		return nil, fmt.Errorf("failed to embed question: %w", err)
	}

	var bestScore float64 = -1
	var bestResult *SemanticSearchResult
	found := false

	for i := range results {
		result := &results[i]

		for _, paragraph := range splitParagraphs(result.Content) {
			embedding, err := m.embedder.Embed(paragraph)
			if err != nil {
				// Skip sentences that fail to embed rather than
				// failing the whole answer - matches the graceful
				// degradation pattern used elsewhere in this
				// package (see ImportDocument).
				continue
			}

			score := CosineSimilarity(questionEmbedding, embedding)
			if score > bestScore {
				bestScore = score
				bestResult = result
				found = true
			}
		}
	}

	if !found {
		return nil, ErrNoExtractiveMatch
	}

	return &ExtractedItem{
		Text:             bestResult.Content,
		SourceDocumentID: bestResult.DocumentID,
		SourceChunkID:    bestResult.ChunkID,
	}, nil
}

// Stats returns knowledge base statistics.
func (m *Manager) Stats() (*Stats, error) {

	documents, err := m.ListKnowledge()
	if err != nil {
		return nil, err
	}

	stats := &Stats{
		DocumentCount: len(documents),
	}

	for _, doc := range documents {
		stats.ChunkCount += doc.ChunkCount
	}

	return stats, nil
}
