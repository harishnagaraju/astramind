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
	Type    StreamEventType
	Content string
	Err     error
}

// Stream defines a provider-independent stream of events.
type Stream interface {
	Events() <-chan StreamEvent
}

// StreamingProvider is implemented by providers that support streaming.
type StreamingProvider interface {
	Stream(
		ctx context.Context,
		request ChatRequest,
	) (Stream, error)
}
