package ai

import (
	"testing"

	"github.com/harishnagaraju/astramind/internal/models"
)

func TestMockProvider(t *testing.T) {

	provider := &MockProvider{}

	reply, err := provider.Chat(
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

func TestProviderName(t *testing.T) {

	p := &MockProvider{}

	if p.Name() != "Mock AI" {
		t.Fatal("Incorrect provider name")
	}
}

func TestInitializeMockCache(t *testing.T) {

	err := InitializeMockCache()

	if err != nil {

		t.Fatalf(
			"InitializeMockCache failed: %v",
			err,
		)

	}

	if len(mockCache) == 0 {

		t.Fatal(
			"Mock cache is empty",
		)

	}
}
