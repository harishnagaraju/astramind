package chat

import "strings"

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
