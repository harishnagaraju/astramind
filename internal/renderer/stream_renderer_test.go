package renderer

import (
	"bytes"
	"testing"

	"github.com/harishnagaraju/astramind/internal/ai"
)

type mockStream struct {
	events chan ai.StreamEvent
}

func (m *mockStream) Events() <-chan ai.StreamEvent {
	return m.events
}

func TestRenderTokens(t *testing.T) {

	var output bytes.Buffer

	renderer := New(&output)

	stream := &mockStream{
		events: make(chan ai.StreamEvent),
	}

	go func() {

		stream.events <- ai.StreamEvent{
			Type:    ai.StreamEventToken,
			Content: "Hello",
		}

		stream.events <- ai.StreamEvent{
			Type:    ai.StreamEventToken,
			Content: " World",
		}

		stream.events <- ai.StreamEvent{
			Type: ai.StreamEventDone,
		}

		close(stream.events)
	}()

	err := renderer.Render(stream)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}
}
