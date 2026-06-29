package ai

type ProviderManager struct {
	provider Provider
}

func NewProviderManager(
	p Provider,
) *ProviderManager {

	return &ProviderManager{
		provider: p,
	}
}

func (pm *ProviderManager) Provider() Provider {
	return pm.provider
}

func (pm *ProviderManager) Chat(
	request ChatRequest,
) (string, error) {

	return pm.provider.Chat(
		request,
	)
}
