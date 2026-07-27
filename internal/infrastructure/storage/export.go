package storage

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

func ExportSession(
	session string,
	messages []models.Message,
) error {

	file := fmt.Sprintf(
		"exports/%s.txt",
		session,
	)

	err := os.MkdirAll(filepath.Dir(file), 0755)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close() //nolint:errcheck // safety-net close on early-return error paths above; the real close+flush is checked explicitly on the success path

	// Export Header
	_, err = fmt.Fprintf(
		f,
		`==================================================
             AstraMind Conversation Export
==================================================

Session      : %s
Exported On  : %s
Messages      : %d

==================================================

`,
		session,
		time.Now().Format("2006-01-02 15:04:05"),
		len(messages),
	)

	if err != nil {
		return err
	}

	// Conversation
	for i, msg := range messages {

		_, err := fmt.Fprintf(
			f,
			`%d. %s
--------------------------------------------------
%s

`,
			i+1,
			strings.ToUpper(msg.Role),
			msg.Content,
		)

		if err != nil {
			return err
		}
	}

	// Footer
	_, err = fmt.Fprintln(
		f,
		`==================================================
End of Conversation
==================================================`,
	)

	if err != nil {
		return err
	}

	// Explicit close (not just the deferred one) so a flush failure -
	// which would mean the exported conversation file was silently
	// truncated on disk - is actually surfaced, not swallowed.
	if err := f.Close(); err != nil {
		return err
	}

	return nil
}