package chat

import "errors"

var (
	ErrInvalidCommand = errors.New("invalid command")

	ErrKnowledgeBaseUnavailable = errors.New(
		"knowledge base is not configured",
	)
)
