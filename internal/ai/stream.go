package ai

import "context"

// StreamEventType identifies the kind of streamed event.
type StreamEventType int

const (
	StreamEventToken StreamEventType = iota
	StreamEventDone
	StreamEventError
)

// StreamEvent represents a single event emitted by a provider.
type StreamEvent struct {
	Type  StreamEventType
	Token string
	Err   error
}

// Stream defines a provider-independent stream of events.
type Stream interface {
	Events() <-chan StreamEvent
}

// StreamRequest contains the information required for a streaming request.
type StreamRequest struct {
	Prompt string
}

// StreamingProvider is implemented by providers that support streaming.
type StreamingProvider interface {
	Stream(ctx context.Context, req StreamRequest) (Stream, error)
}
