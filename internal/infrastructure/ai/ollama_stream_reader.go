package ai

import (
	"bufio"
	"encoding/json"
	"io"
)

func readOllamaStream(
	body io.Reader,
	events chan<- StreamEvent,
) {

	defer close(events)

	scanner := bufio.NewScanner(body)

	for scanner.Scan() {

		var response OllamaStreamResponse

		err := json.Unmarshal(
			scanner.Bytes(),
			&response,
		)

		if err != nil {
			events <- StreamEvent{
				Type: StreamEventError,
				Err:  err,
			}
			return
		}

		if response.Message.Content != "" {
			events <- StreamEvent{
				Type:    StreamEventToken,
				Content: response.Message.Content,
			}
		}

		if response.Done {
			events <- StreamEvent{
				Type: StreamEventDone,
			}
			return
		}
	}

	if err := scanner.Err(); err != nil {
		events <- StreamEvent{
			Type: StreamEventError,
			Err:  err,
		}
	}
}
