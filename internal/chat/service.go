package chat

import (
	"context"
	"io"

	"github.com/harishnagaraju/astramind/internal/infrastructure/ai"
	"github.com/harishnagaraju/astramind/internal/kb"
	"github.com/harishnagaraju/astramind/internal/renderer"
)

// Dependencies contains all services used by the chat package.
// New subsystem dependencies should be added here instead of
// directly expanding the Service struct.
type Dependencies struct {
	ProviderManager *ai.ProviderManager
	KnowledgeBase   *kb.Manager
}

type Service struct {
	deps Dependencies
}

func NewService(
	deps Dependencies,
) *Service {

	return &Service{
		deps: deps,
	}
}

func (s *Service) Chat(
	ctx context.Context,
	writer io.Writer,
	request ai.ChatRequest,
) (string, bool, error) {

	streamingProvider, ok := s.deps.ProviderManager.Provider().(ai.StreamingProvider)

	/* if ok {
		println("STREAMING ENABLED")
	} else {
		println("STREAMING DISABLED")
	} */

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

	reply, err := s.deps.ProviderManager.Chat(request)

	return reply, false, err
}
