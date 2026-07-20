package kb

import "errors"

var (
	ErrDocumentNotFound = errors.New("knowledge base document not found")
	ErrInvalidDocument  = errors.New("invalid knowledge base document")
	ErrInvalidChunk     = errors.New("invalid knowledge base chunk")
)
