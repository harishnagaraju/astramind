package ai

import (
	"testing"

	"github.com/harishnagaraju/astramind/internal/models"
)

func TestProviderManager(t *testing.T) {

	mock := &MockProvider{}

	manager := NewProviderManager(mock)

	if manager.Provider() == nil {
		t.Fatal("Provider manager returned nil provider")
	}

	if manager.Provider().Name() != "Mock AI" {
		t.Fatal("Unexpected provider name")
	}
}

func TestProviderManagerChat(t *testing.T) {

	manager := NewProviderManager(
		&MockProvider{},
	)

	reply, err := manager.Chat(
		ChatRequest{
			Messages: []models.Message{
				{
					Role:    "user",
					Content: "hello",
				},
			},
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	expected := "Hello! Nice to meet you."

	if reply != expected {
		t.Fatalf(
			"Expected '%s' but got '%s'",
			expected,
			reply,
		)
	}
}
