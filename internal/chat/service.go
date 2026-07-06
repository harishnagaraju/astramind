package chat

import (
	"context"
	"io"

	"github.com/harishnagaraju/astramind/internal/ai"
	"github.com/harishnagaraju/astramind/internal/renderer"
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
	ctx context.Context,
	writer io.Writer,
	request ai.ChatRequest,
) (string, bool, error) {

	streamingProvider, ok := s.manager.Provider().(ai.StreamingProvider)

	if ok {

		stream, err := streamingProvider.Stream(
			ctx,
			request,
		)

		if err != nil {
			return "", false, err
		}

		r := renderer.New(writer)

		if err := r.Render(stream); err != nil {
			return "", true, err
		}

		return r.Text(), true, nil
	}

	reply, err := s.manager.Chat(request)

	return reply, false, err
}
