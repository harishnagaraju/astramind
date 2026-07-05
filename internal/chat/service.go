package chat

import (
	"github.com/harishnagaraju/astramind/internal/ai"
)

type Service struct {
	manager *ai.ProviderManager
}

func NewService(
	manager *ai.ProviderManager,
) *Service {

	return &Service{
		manager: manager,
	}
}

func (s *Service) Chat(
	request ai.ChatRequest,
) (string, error) {

	return s.manager.Chat(request)
}
