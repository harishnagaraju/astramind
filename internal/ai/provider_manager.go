package ai

import "fmt"

type ProviderManager struct {
	provider Provider
	fallback Provider
}

func NewProviderManager(
	p Provider,
) *ProviderManager {

	return &ProviderManager{
		provider: p,
		fallback: &MockProvider{},
	}
}

func (pm *ProviderManager) Provider() Provider {
	return pm.provider
}

func (pm *ProviderManager) ProviderName() string {

	if pm.provider == nil {
		return "None"
	}

	return pm.provider.Name()
}

func (pm *ProviderManager) FallbackProvider() Provider {

	return pm.fallback

}

func (pm *ProviderManager) Chat(
	request ChatRequest,
) (string, error) {

	reply, err := pm.provider.Chat(request)

	if err == nil {
		return reply, nil
	}

	if pm.fallback == nil {
		return "", err
	}

	if pm.provider.Name() == pm.fallback.Name() {
		return "", err
	}

	fmt.Println()
	fmt.Println("⚠ Primary provider failed.")
	fmt.Printf("Switching to %s...\n", pm.fallback.Name())
	fmt.Println()

	pm.provider = pm.fallback

	return pm.provider.Chat(request)
}
