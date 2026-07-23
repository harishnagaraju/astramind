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

	// Route enumeration-style questions ("what are all the X", "list
	// every X", "what are the X timings") to a deterministic
	// extraction path instead of the free-form LLM path.
	//
	// This replaces an earlier approach that tried to fix incomplete
	// enumeration by rewording the question and tuning prompt/
	// temperature. That approach was fundamentally limited: once a
	// chunk is retrieved, whether the LLM's free-form generation
	// mentions every item in it is a property of the model's
	// generation process, not something prompt wording can
	// guarantee. Confirmed during testing that even a complete,
	// correct, uncorrupted prompt (chunking bug fixed, retrieval
	// limit not a factor) still produced inconsistent, incomplete
	// answers - the variance was in the LLM step itself, which no
	// amount of prompt engineering can make deterministic.
	//
	// The deterministic path below has no such gap: chunks are
	// already paragraph-structured (see chunker.go), so extracting
	// "all items in the retrieved chunks" is a matter of splitting
	// already-retrieved text back into paragraphs and formatting
	// them - no generation step, so no possibility of an item being
	// silently dropped.
	if kb.IsEnumerationQuery(question) {
		return s.handleKBAskEnumeration(question)
	}

	return s.handleKBAskFreeform(question)
}

// handleKBAskEnumeration answers "list everything matching X" style
// questions deterministically: retrieve relevant chunks, extract
// every item within them, format as a list. No LLM call, no
// possibility of an item being silently omitted by generation.
//
// Deliberately does NOT filter extracted items by keyword overlap
// with the question - retrieval already decided which chunks are
// relevant. Filtering items again afterward is exactly the mechanism
// that caused the original bug (narrow keyword matching silently
// dropping related-but-differently-worded entries). Once a chunk is
// judged relevant, every item in it is included.
func (s *Service) handleKBAskEnumeration(question string) error {
	results, err := s.deps.KnowledgeBase.SemanticSearch(question)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No relevant knowledge found to answer this question.")
		return nil
	}

	items := kb.ExtractItems(results)
	answer := kb.BuildListAnswer(items)

	fmt.Println(answer)

	return nil
}

// handleKBAskFreeform answers single-fact questions. Tries the
// extractive path first (Pattern B: return the matching source
// paragraph verbatim, no LLM paraphrase, so no risk of a generative
// model restating a date or fee incorrectly) and only falls back to
// the LLM RAG path if extraction can't run at all (no embedder
// configured) - the same graceful-degradation pattern ImportDocument
// already uses elsewhere in this package for missing embedders.
//
// A genuine multi-constraint reasoning question ("I broke it after
// 10 days, am I covered?") still needs synthesis across several
// facts, which extractive matching can't do - but that class of
// question is intentionally out of scope for /kb ask's current
// design; it answers questions about what the knowledge base states,
// not compound reasoning over it.
func (s *Service) handleKBAskFreeform(question string) error {
	results, err := s.deps.KnowledgeBase.SemanticSearch(question)
	if err != nil {
		return err
	}

	if len(results) == 0 {
		fmt.Println("No relevant knowledge found to answer this question.")
		return nil
	}

	item, err := s.deps.KnowledgeBase.ExtractiveAnswer(question, results)
	if err == nil {
		fmt.Println(item.Text)
		fmt.Println()
		fmt.Println("Sources:")
		fmt.Printf("  [%s]\n", item.SourceChunkID)
		return nil
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
		return err
	}

	fmt.Println(reply)
	fmt.Println()
	fmt.Println("Sources:")

	seen := make(map[string]bool)

	for _, result := range results {
		if seen[result.DocumentID] {
			continue
		}
		seen[result.DocumentID] = true
		fmt.Printf("  [%s]\n", result.DocumentID)
	}

	return nil
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
