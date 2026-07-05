package renderer

import (
	"github.com/harishnagaraju/astramind/internal/ai"
)

// StreamRenderer renders streaming AI responses.
type StreamRenderer struct{}

// New creates a new StreamRenderer.
func New() *StreamRenderer {
	return &StreamRenderer{}
}

// Render will consume streaming events.
//
// The implementation will be added in the next step.
func (r *StreamRenderer) Render(stream ai.Stream) error {
	return nil
}
