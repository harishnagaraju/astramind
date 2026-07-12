package chat

import "strings"
import "fmt"

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
		fmt.Printf(
			"%s (%d chunks)\n",
			doc.Name,
			doc.ChunkCount,
		)
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
