package ai

import (
	"bufio"
	"io"
	"strings"
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

		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines.
		if line == "" {
			continue
		}

		// Ignore anything that is not an SSE data event.
		if !strings.HasPrefix(line, "data:") {
			continue
		}

		// Remove the SSE prefix.
		data := strings.TrimSpace(
			strings.TrimPrefix(line, "data:"),
		)

		// End of stream.
		if data == "[DONE]" {

			stream.events <- StreamEvent{
				Type: StreamEventDone,
			}

			return
		}

		// JSON parsing will be implemented in the next step.
		_ = data
	}

	if err := scanner.Err(); err != nil {

		stream.events <- StreamEvent{
			Type: StreamEventError,
			Err:  err,
		}
	}
}
