package ai

// openAIStream is the concrete implementation of the Stream interface
// used by the OpenAI provider.
type openAIStream struct {
	events chan StreamEvent
}

// Events returns the stream of events emitted by the provider.
func (s *openAIStream) Events() <-chan StreamEvent {
	return s.events
}
