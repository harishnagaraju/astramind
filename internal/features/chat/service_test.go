package chat

import (
	"bytes"
	"context"
	"testing"

	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
)

func TestServiceStreaming(t *testing.T) {

	provider := &ai.MockProvider{}

	manager := ai.NewProviderManager(provider)

	service := NewService(
		Dependencies{
			ProviderManager: manager,
		},
	)

	var output bytes.Buffer

	reply, streamed, err := service.Chat(
		context.Background(),
		&output,
		ai.ChatRequest{},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if !streamed {
		t.Fatal("expected streaming response")
	}

	expectedReply := "Hello from Mock AI!"

	if reply != expectedReply {
		t.Fatalf(
			"expected reply %q, got %q",
			expectedReply,
			reply,
		)
	}

	expectedOutput := expectedReply + "\n"

	if output.String() != expectedOutput {
		t.Fatalf(
			"expected output %q, got %q",
			expectedOutput,
			output.String(),
		)
	}
}
