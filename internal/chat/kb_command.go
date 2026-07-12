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

	// Command recognized.
	return true, nil
}
