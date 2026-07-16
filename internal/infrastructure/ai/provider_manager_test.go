package ai

import (
	"fmt"
	"testing"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
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

func TestProviderManagerProviderName(t *testing.T) {

	manager := NewProviderManager(
		&MockProvider{},
	)

	if manager.ProviderName() != "Mock AI" {

		t.Fatalf(
			"Expected Mock AI but got %s",
			manager.ProviderName(),
		)

	}
}

func TestFallbackProvider(t *testing.T) {

	manager := NewProviderManager(
		&OpenAIProvider{},
	)

	if manager.FallbackProvider() == nil {

		t.Fatal(
			"Fallback provider is nil",
		)

	}

	if manager.FallbackProvider().Name() != "Mock AI" {

		t.Fatal(
			"Incorrect fallback provider",
		)

	}

}

type FailingProvider struct{}

func (f *FailingProvider) Name() string {
	return "Broken Provider"
}

func (f *FailingProvider) Chat(
	request ChatRequest,
) (string, error) {

	return "", fmt.Errorf("provider failed")

}

func TestAutomaticFailover(t *testing.T) {

	manager := NewProviderManager(
		&FailingProvider{},
	)

	reply, err := manager.Chat(ChatRequest{
		Messages: []models.Message{
			{
				Role:    "user",
				Content: "hello",
			},
		},
	})

	if err != nil {
		t.Fatal(err)
	}

	if reply != "Hello! Nice to meet you." {

		t.Fatalf(
			"Unexpected reply: %s",
			reply,
		)

	}

	if manager.ProviderName() != "Mock AI" {

		t.Fatal(
			"Provider did not switch to Mock AI",
		)

	}
}
