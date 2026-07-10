package ai

import "testing"

func TestFactoryReturnsMockProvider(t *testing.T) {

	p := NewProvider(
		ProviderConfig{},
	)

	if p.Name() != "Mock AI" {
		t.Fatalf(
			"Expected Mock AI but got %s",
			p.Name(),
		)
	}
}

func TestFactoryReturnsOpenAIProvider(t *testing.T) {

	p := NewProvider(
		ProviderConfig{
			APIKey: "dummy-key",
		},
	)

	if p.Name() != "OpenAI" {
		t.Fatalf(
			"Expected OpenAI but got %s",
			p.Name(),
		)
	}
}

func TestExplicitMockProvider(t *testing.T) {

	p := NewProvider(
		ProviderConfig{
			Provider: "mock",
			APIKey:   "dummy-key",
		},
	)

	if p.Name() != "Mock AI" {
		t.Fatalf(
			"Expected Mock AI but got %s",
			p.Name(),
		)
	}
}

func TestExplicitOpenAIProvider(t *testing.T) {

	p := NewProvider(
		ProviderConfig{
			Provider: "openai",
		},
	)

	if p.Name() != "OpenAI" {
		t.Fatalf(
			"Expected OpenAI but got %s",
			p.Name(),
		)
	}
}

func TestExplicitOllamaProvider(t *testing.T) {

	p := NewProvider(
		ProviderConfig{
			Provider: "ollama",
		},
	)

	if p == nil {
		t.Fatal("Expected provider but got nil")
	}

	if p.Name() != "Ollama" {
		t.Fatalf(
			"Expected Ollama but got %s",
			p.Name(),
		)
	}
}
