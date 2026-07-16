package engine

import "github.com/harishnagaraju/astramind/internal/infrastructure/models"

// RuntimeContext contains mutable runtime state for
// the current interactive session.
type RuntimeContext struct {
	Conversation []models.Message
}
