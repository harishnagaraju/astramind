package ai

import (
	"bufio"
	"io"
)

// readStream reads the OpenAI streaming response.
//
// In this phase it only owns the lifecycle of the response body.
// SSE parsing will be implemented in the next phase.
func (p *OpenAIProvider) readStream(
	body io.ReadCloser,
	stream *openAIStream,
) {
	defer body.Close()
	defer close(stream.events)

	scanner := bufio.NewScanner(body)

	for scanner.Scan() {

		line := scanner.Text()

		_ = line
	}
}
