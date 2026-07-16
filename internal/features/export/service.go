package export

import (
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

// Service provides export operations.
type Service struct{}

// NewService creates a new export service.
func NewService() *Service {
	return &Service{}
}

// ExportTXT exports a conversation as plain text.
func (s *Service) ExportTXT(
	session string,
	messages []models.Message,
) error {

	return storage.ExportSession(
		session,
		messages,
	)
}

// ExportMarkdown exports a conversation as Markdown.
func (s *Service) ExportMarkdown(
	session string,
	messages []models.Message,
) error {

	return storage.ExportMarkdown(
		session,
		messages,
	)
}
