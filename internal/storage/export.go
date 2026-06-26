package storage

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/harishnagaraju/astramind/internal/models"
)

func ExportSession(
	session string,
	messages []models.Message,
) error {

	file := fmt.Sprintf(
		"exports/%s.txt",
		session,
	)

	f, err := os.Create(file)
	if err != nil {
		return err
	}

	defer f.Close()

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

	return nil
}