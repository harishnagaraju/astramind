package engine

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// runScript executes commands from a script file, one per line, through
// the same command dispatcher used in interactive mode (a.dispatcher).
// This means script mode supports every registered command - session,
// conversation, search, export, and knowledge base - not just a subset.
func (a *App) runScript() error {

	if len(os.Args) != 3 || os.Args[1] != "--script" {
		return nil
	}

	file, err := os.Open(os.Args[2])
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "#") {
			continue
		}

		if line == "exit" || line == "quit" {
			return nil
		}

		if _, err := a.dispatcher.Execute(line); err != nil {
			return fmt.Errorf("script command %q failed: %w", line, err)
		}
	}

	return scanner.Err()
}
