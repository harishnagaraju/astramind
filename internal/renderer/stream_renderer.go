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

func (r *StreamRenderer) Render(
	stream ai.Stream,
) error {

	for event := range stream.Events() {

		switch event.Type {

		case ai.StreamEventToken:
			// Token rendering will be added next.

		case ai.StreamEventDone:
			return nil

		case ai.StreamEventError:
			return event.Err
		}
	}

	return nil
}
