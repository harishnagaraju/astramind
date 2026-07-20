package ai

import (
	"context"
	"hash/fnv"
	"strings"
)

type MockProvider struct{}

func (m *MockProvider) Name() string {
	return "Mock AI"
}

func (m *MockProvider) Chat(
	request ChatRequest,
) (string, error) {

	if len(request.Messages) == 0 {
		return "Mock AI: Empty conversation.", nil
	}

	last := request.Messages[len(request.Messages)-1]

	reply := GetMockResponse(
		strings.TrimSpace(last.Content),
	)

	return reply, nil
}

// Embed returns a deterministic pseudo-embedding derived from the
// input text. It performs no network calls, which makes it safe for
// unit tests that exercise embedding-based logic (e.g. cosine
// similarity ranking) without depending on a real provider.
func (m *MockProvider) Embed(
	request EmbeddingRequest,
) ([]float32, error) {

	const dimensions = 16

	vector := make([]float32, dimensions)

	words := strings.Fields(
		strings.ToLower(request.Text),
	)

	if len(words) == 0 {
		return vector, nil
	}

	for _, word := range words {

		h := fnv.New32a()
		h.Write([]byte(word))
		hash := h.Sum32()

		index := int(hash) % dimensions
		if index < 0 {
			index += dimensions
		}

		vector[index] += 1
	}

	return vector, nil
}

func (m *MockProvider) Stream(
	ctx context.Context,
	request ChatRequest,
) (Stream, error) {

	stream := &openAIStream{
		events: make(chan StreamEvent),
	}

	go func() {
		defer close(stream.events)

		// Respect context cancellation.
		select {
		case <-ctx.Done():
			stream.events <- StreamEvent{
				Type: StreamEventError,
				Err:  ctx.Err(),
			}
			return
		default:
		}

		stream.events <- StreamEvent{
			Type:    StreamEventToken,
			Content: "Hello",
		}

		stream.events <- StreamEvent{
			Type:    StreamEventToken,
			Content: " from",
		}

		stream.events <- StreamEvent{
			Type:    StreamEventToken,
			Content: " Mock",
		}

		stream.events <- StreamEvent{
			Type:    StreamEventToken,
			Content: " AI!",
		}

		stream.events <- StreamEvent{
			Type: StreamEventDone,
		}
	}()

	return stream, nil
}
