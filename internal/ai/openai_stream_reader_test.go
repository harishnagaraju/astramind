package ai

import (
	"io"
	"strings"
	"testing"
)

func newTestStream(data string) io.ReadCloser {
	return io.NopCloser(strings.NewReader(data))
}

func TestReadStreamTokenAndDone(t *testing.T) {

	body := newTestStream(
		`data: {"choices":[{"delta":{"content":"Hello"}}]}

		data: {"choices":[{"delta":{"content":" World"}}]}

		data: [DONE]`,
	)

	stream := &openAIStream{
		events: make(chan StreamEvent),
	}

	provider := &OpenAIProvider{}

	go provider.readStream(body, stream)

	var events []StreamEvent

	for event := range stream.Events() {
		events = append(events, event)
	}

	if len(events) != 3 {
		t.Fatalf(
			"expected 3 events, got %d",
			len(events),
		)
	}

	if events[0].Type != StreamEventToken {
		t.Fatalf(
			"expected first event to be %q, got %q",
			StreamEventToken,
			events[0].Type,
		)
	}

	if events[0].Content != "Hello" {
		t.Fatalf(
			"expected first token %q, got %q",
			"Hello",
			events[0].Content,
		)
	}

	if events[1].Type != StreamEventToken {
		t.Fatalf(
			"expected second event to be %q, got %q",
			StreamEventToken,
			events[1].Type,
		)
	}

	if events[1].Content != " World" {
		t.Fatalf(
			"expected second token %q, got %q",
			" World",
			events[1].Content,
		)
	}

	if events[2].Type != StreamEventDone {
		t.Fatalf(
			"expected final event to be %q, got %q",
			StreamEventDone,
			events[2].Type,
		)
	}

}
