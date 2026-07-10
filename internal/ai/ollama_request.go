package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// buildOllamaRequest builds an HTTP request for the Ollama Chat API.
func buildOllamaRequest(
	baseURL string,
	model string,
	request ChatRequest,
) (*http.Request, error) {

	if model == "" {
		model = "llama3"
	}

	ollamaReq := OllamaChatRequest{
		Model:  model,
		Stream: false,
	}

	for _, msg := range request.Messages {
		ollamaReq.Messages = append(
			ollamaReq.Messages,
			OllamaChatMessage{
				Role:    msg.Role,
				Content: msg.Content,
			},
		)
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal ollama request: %w",
			err,
		)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/api/chat",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create ollama request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	return req, nil
}
