package ai

// ollamaStream is the concrete implementation of the Stream interface
// used by the Ollama provider.
type ollamaStream struct {
	events chan StreamEvent
}

// Events returns the stream of events emitted by the provider.
func (s *ollamaStream) Events() <-chan StreamEvent {
	return s.events
}
