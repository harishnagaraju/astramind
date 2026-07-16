package renderer

import (
	"fmt"
	"io"
	"strings"

	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

// StreamRenderer renders streaming AI responses.
type StreamRenderer struct {
	writer io.Writer
	buffer strings.Builder
}

func (r *StreamRenderer) Text() string {
	return r.buffer.String()
}

// New creates a new StreamRenderer.
func New(writer io.Writer) *StreamRenderer {
	return &StreamRenderer{
		writer: writer,
	}
}

// Render will consume streaming events.
//

func (r *StreamRenderer) Render(
	stream ai.Stream,
) error {

	for event := range stream.Events() {

		switch event.Type {

		case ai.StreamEventToken:
			r.buffer.WriteString(event.Content)

			_, err := fmt.Fprint(r.writer, event.Content)

			if err != nil {
				return err
			}

		case ai.StreamEventDone:
			if _, err := fmt.Fprintln(r.writer); err != nil {
				return err
			}

			return nil

		case ai.StreamEventError:
			return event.Err
		}
	}

	return nil
}
