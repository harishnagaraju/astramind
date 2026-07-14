package app

import "github.com/harishnagaraju/astramind/internal/models"

// RuntimeContext contains mutable runtime state for
// the current interactive session.
type RuntimeContext struct {
	Conversation []models.Message
}
