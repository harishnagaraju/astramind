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

}
