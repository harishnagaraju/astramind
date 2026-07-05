package ai

import (
	"context"
	"testing"
)

func TestMockProviderStream(t *testing.T) {

	provider := &MockProvider{}

	stream, err := provider.Stream(
		context.Background(),
		ChatRequest{},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	var events []StreamEvent

	for event := range stream.Events() {
		events = append(events, event)
	}

	if len(events) != 5 {
		t.Fatalf(
			"expected 5 events, got %d",
			len(events),
		)
	}
}
