package chat

import (
	"fmt"
	"strings"

	"github.com/harishnagaraju/astramind/internal/features/kb"
	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// HandleKnowledgeCommand processes /kb commands.
func (s *Service) HandleKnowledgeCommand(input string) (bool, error) {
	fields := strings.Fields(input)

	if len(fields) == 0 {
		return false, nil
	}

	if fields[0] != "/kb" {
		return false, nil
	}

	if len(fields) == 1 {
		return true, nil
	}

	// Command recognized.
	switch fields[1] {

	case "import":
		return true, s.handleKBImport(fields)

	case "list":
		return true, s.handleKBList()

	case "search":
		return true, s.handleKBSearch(fields)

	case "ssearch":
		return true, s.handleKBSemanticSearch(fields)

	case "ask":
		return true, s.handleKBAsk(fields)

	case "remove":
		return true, s.handleKBRemove(fields)

	case "clear":
		return true, s.handleKBClear()

	case "stats":
		return true, s.handleKBStats()

	default:
		return true, ErrInvalidCommand
	}
}

func (s *Service) handleKBImport(args []string) error {

	if len(args) < 3 {
		return ErrInvalidCommand
	}

	if s.deps.KnowledgeBase == nil {
		return ErrKnowledgeBaseUnavailable
	}

	doc, err := s.deps.KnowledgeBase.ImportDocument(args[2])
	if err != nil {
		return err
	}

	println("Imported:", doc.Name)

	return nil
}

func (s *Service) handleKBList() error {

	documents, err := s.deps.KnowledgeBase.ListKnowledge()
	if err != nil {
		return err
	}

	if len(documents) == 0 {
		fmt.Println("Knowledge base is empty.")
		return nil
	}

	fmt.Println("Knowledge Base Documents")
	fmt.Println("------------------------")

	for _, doc := range documents {
		fmt.Printf(" %s\n", doc.ID)
		fmt.Printf(" Name   : %s\n", doc.Name)
		fmt.Printf(" Chunks : %d\n\n", doc.ChunkCount)
	}
	return nil
}

func (s *Service) handleKBSearch(args []string) error {

	if len(args) < 3 {
		return ErrInvalidCommand
	}

	query := strings.Join(args[2:], " ")

	results, err := s.deps.KnowledgeBase.Search(query)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No matching knowledge found.")
		return nil
	}

	fmt.Println("Knowledge Search Results")
	fmt.Println("------------------------")

	for _, chunk := range results {
		fmt.Printf(
			"[%s]\n%s\n\n",
			chunk.DocumentID,
			chunk.Content,
		)
	}

	return nil
}

func (s *Service) handleKBSemanticSearch(args []string) error {

	if len(args) < 3 {
		return ErrInvalidCommand
	}

	query := strings.Join(args[2:], " ")

	results, err := s.deps.KnowledgeBase.SemanticSearch(query)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No matching knowledge found.")
		return nil
	}

	fmt.Println("Semantic Search Results")
	fmt.Println("------------------------")

	for _, chunk := range results {
		fmt.Printf(
			"[%s] (similarity: %.3f)\n%s\n\n",
			chunk.DocumentID,
			chunk.Score,
			chunk.Content,
		)
	}

	return nil
}

func (s *Service) handleKBAsk(args []string) error {

	if len(args) < 3 {
		return ErrInvalidCommand
	}

	question := strings.Join(args[2:], " ")

	answer, _, err := s.Ask(question)
	if err != nil {
		return err
	}

	fmt.Println(answer)

	return nil
}

// Ask answers a knowledge-base question and returns the complete,
// display-ready answer text (including any source citations already
// formatted into it) along with a separate, structured list of
// source document IDs for callers that need them independently (e.g.
// a JSON API response field).
//
// This is the single implementation of /kb ask's routing logic - the
// CLI's handleKBAsk and the web UI's /api/ask handler both call this
// rather than each implementing their own version of it. Before this
// refactor, webserver.go had its own separate, hand-written copy of
// the old free-form-LLM-only ask logic, written before this package's
// RAG reliability redesign (deterministic enumeration + extraction,
// see IsEnumerationQuery/ExtractItems/ExtractiveAnswer) and never
// updated when that redesign happened - meaning the web UI kept
// serving the old, less reliable behavior silently, with no test or
// build failure to catch the drift, since the two implementations
// were entirely independent and neither referenced the other. A
// single shared method makes that class of drift structurally
// impossible: there is now only one place this logic can be wrong.
func (s *Service) Ask(question string) (answer string, sources []string, err error) {
	if kb.IsEnumerationQuery(question) {
		return s.askEnumeration(question)
	}

	return s.askFreeform(question)
}

// askEnumeration answers "list everything matching X" style
// questions deterministically: retrieve relevant chunks, extract
// every item within them, format as a list. No LLM call, no
// possibility of an item being silently omitted by generation.
//
// Whole chunks are filtered for relevance (see FilterRelevantChunks)
// before extraction, but individual items within an included chunk
// are never filtered against the question's wording - retrieval
// (and now chunk-level filtering) already decided which chunks are
// relevant. Filtering individual items within a relevant chunk was
// tried and rejected: narrow keyword matching silently dropped
// related-but-differently-worded entries (e.g. "Tuesday Chanting"
// being dropped from a Sanskrit-timings answer because it doesn't
// contain the word "Sanskrit", despite genuinely being part of the
// same schedule). Once a chunk is judged relevant, every item in it
// is included.
func (s *Service) askEnumeration(question string) (string, []string, error) {
	results, err := s.deps.KnowledgeBase.SemanticSearch(question)
	if err != nil {
		return "", nil, err
	}

	if len(results) == 0 {
		return "No relevant knowledge found to answer this question.", nil, nil
	}

	// Relevance gate: SemanticSearch always returns its top-ranked
	// chunks regardless of whether any of them are actually related
	// to the question - "most similar of what exists" isn't the same
	// as "relevant". If the question shares literally zero words with
	// anything retrieved, treat that as no relevant knowledge found
	// rather than confidently presenting unrelated content. See
	// kb.HasAnyContentWordOverlap's doc comment for why this uses
	// word overlap rather than an embedding similarity cutoff.
	if !kb.HasAnyContentWordOverlap(question, results) {
		return "No relevant knowledge found to answer this question.", nil, nil
	}

	// Chunk-level relevance filter: HasAnyContentWordOverlap above
	// only answers "is anything at all relevant" (all-or-nothing for
	// the whole batch) - it does not remove individual irrelevant
	// chunks from an otherwise-valid batch. Confirmed in real testing
	// (real Ollama embeddings): a small knowledge base containing
	// both Sanskrit1.txt and an unrelated service contract answered
	// a Sanskrit-specific question with both documents bundled
	// together, because SemanticSearch returns every chunk whenever
	// the KB's total chunk count is at or below its retrieval limit.
	// FilterRelevantChunks removes only whole chunks with zero
	// content-word overlap - it does NOT filter individual items
	// within a chunk (see the removed comment below for why that was
	// tried and rejected: it silently dropped real entries, like
	// "Tuesday Chanting", from an otherwise-relevant chunk just
	// because the entry itself didn't repeat the word "Sanskrit").
	results = kb.FilterRelevantChunks(question, results)

	items := kb.ExtractItems(results)
	answer := kb.BuildListAnswer(items)

	return answer, uniqueDocumentIDs(results), nil
}

// askFreeform answers single-fact questions. Tries the extractive
// path first (Pattern B: return the matching source paragraph
// verbatim, no LLM paraphrase, so no risk of a generative model
// restating a date or fee incorrectly) and only falls back to the
// LLM RAG path if extraction can't run at all (no embedder
// configured) - the same graceful-degradation pattern ImportDocument
// already uses elsewhere in this package for missing embedders.
//
// A genuine multi-constraint reasoning question ("I broke it after
// 10 days, am I covered?") still needs synthesis across several
// facts, which extractive matching can't do - but that class of
// question is intentionally out of scope for /kb ask's current
// design; it answers questions about what the knowledge base states,
// not compound reasoning over it.
func (s *Service) askFreeform(question string) (string, []string, error) {
	results, err := s.deps.KnowledgeBase.SemanticSearch(question)
	if err != nil {
		return "", nil, err
	}

	if len(results) == 0 {
		return "No relevant knowledge found to answer this question.", nil, nil
	}

	// Same relevance gate as askEnumeration - see that function's
	// comment for why this exists and why it uses word overlap.
	if !kb.HasAnyContentWordOverlap(question, results) {
		return "No relevant knowledge found to answer this question.", nil, nil
	}

	item, err := s.deps.KnowledgeBase.ExtractiveAnswer(question, results)
	if err == nil {
		sources := []string{item.SourceChunkID}
		return formatAnswerWithSources(item.Text, sources), sources, nil
	}
	// Any extractive failure (no embedder configured, no embeddable
	// item found) falls through to the LLM path below rather than
	// failing the command outright.

	prompt := kb.BuildSemanticPrompt(question, results)

	// Moderate-low temperature: reduces unnecessary variance in
	// wording without the over-correction seen at very low values
	// (0.2 was tested and made the model MORE narrowly literal, the
	// opposite of what was needed for enumeration - enumeration no
	// longer goes through this path at all, but kept moderate rather
	// than default since this is still an extraction task, not open
	// conversation).
	ragTemperature := 0.4

	reply, err := s.deps.ProviderManager.Chat(ai.ChatRequest{
		Messages: []models.Message{
			{Role: "user", Content: prompt},
		},
		Temperature: &ragTemperature,
	})
	if err != nil {
		return "", nil, err
	}

	sources := uniqueDocumentIDs(results)
	return formatAnswerWithSources(reply, sources), sources, nil
}

// formatAnswerWithSources appends a "Sources:" section to answer
// text, matching the format kb.BuildListAnswer already uses for the
// enumeration path - so callers (CLI, web UI) can treat every Ask()
// result uniformly, regardless of which internal path produced it.
func formatAnswerWithSources(answer string, sources []string) string {
	if len(sources) == 0 {
		return answer
	}

	var b strings.Builder
	b.WriteString(answer)
	b.WriteString("\n\nSources:\n")
	for _, src := range sources {
		b.WriteString("  [")
		b.WriteString(src)
		b.WriteString("]\n")
	}
	return b.String()
}

func uniqueDocumentIDs(results []kb.SemanticSearchResult) []string {
	seen := make(map[string]bool)
	var out []string

	for _, r := range results {
		if seen[r.DocumentID] {
			continue
		}
		seen[r.DocumentID] = true
		out = append(out, r.DocumentID)
	}

	return out
}

func (s *Service) handleKBRemove(args []string) error {

	if len(args) != 3 {
		return ErrInvalidCommand
	}

	documentID := args[2]

	if err := s.deps.KnowledgeBase.RemoveKnowledge(documentID); err != nil {
		return err
	}

	fmt.Println("Removed:", documentID)

	return nil
}

func (s *Service) handleKBClear() error {

	if err := s.deps.KnowledgeBase.ClearKnowledge(); err != nil {
		return err
	}

	fmt.Println("Knowledge base cleared.")

	return nil
}

func (s *Service) handleKBStats() error {

	stats, err := s.deps.KnowledgeBase.Stats()
	if err != nil {
		return err
	}

	fmt.Println("Knowledge Base Statistics")
	fmt.Println("-------------------------")
	fmt.Printf("Documents : %d\n", stats.DocumentCount)
	fmt.Printf("Chunks    : %d\n", stats.ChunkCount)

	return nil
}