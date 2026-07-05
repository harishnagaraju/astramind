package chat

import (
	"bytes"
	"context"
	"testing"

	"github.com/harishnagaraju/astramind/internal/ai"
)

func TestServiceStreaming(t *testing.T) {

	provider := &ai.MockProvider{}

	manager := ai.NewProviderManager(provider)

	service := NewService(manager)

	var output bytes.Buffer

	reply, err := service.Chat(
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

	if reply == "" {
		t.Fatal("expected non-empty reply")
	}
}
