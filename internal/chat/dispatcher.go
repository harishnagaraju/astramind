package chat

// ExecuteCommand executes a single CLI command.
//
// Returns:
//
//	handled = true  -> the command was recognized.
//	handled = false -> not a command.
//
// error is returned if the command execution failed.
func (s *Service) ExecuteCommand(input string) (bool, error) {

	// Knowledge Base commands.
	handled, err := s.HandleKnowledgeCommand(input)
	if handled || err != nil {
		return handled, err
	}

	// Future command handlers.
	//
	// handled, err = s.HandleSessionCommand(input)
	// if handled || err != nil {
	//     return handled, err
	// }

	// handled, err = s.HandleProviderCommand(input)
	// ...

	return false, nil
}
