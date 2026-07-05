package renderer

import (
	"bytes"
	"errors"
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

	expectedOutput := "Hello World\n"

	if output.String() != expectedOutput {
		t.Fatalf(
			"expected output %q, got %q",
			expectedOutput,
			output.String(),
		)
	}

	expectedText := "Hello World"

	if renderer.Text() != expectedText {
		t.Fatalf(
			"expected text %q, got %q",
			expectedText,
			renderer.Text(),
		)
	}

}

func TestRenderError(t *testing.T) {

	expectedErr := errors.New("stream failure")

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
			Type: ai.StreamEventError,
			Err:  expectedErr,
		}

		close(stream.events)
	}()

	err := renderer.Render(stream)

	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if err != expectedErr {
		t.Fatalf(
			"expected %v, got %v",
			expectedErr,
			err,
		)
	}
}
