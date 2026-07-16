package chat

import (
	"bufio"
	"os"
	"strings"
)

// ExecuteScript executes CLI commands from a script file.
func (s *Service) ExecuteScript(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Ignore blank lines.
		if line == "" {
			continue
		}

		// Ignore comments.
		if strings.HasPrefix(line, "#") {
			continue
		}

		// Exit command.
		if line == "exit" || line == "quit" {
			return nil
		}

		// Execute Knowledge Base commands.
		handled, err := s.ExecuteCommand(line)
		if err != nil {
			return err
		}

		if handled {
			continue
		}

		// Future command handlers can be added here.
	}

	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}
