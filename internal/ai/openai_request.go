package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

func (p *OpenAIProvider) buildRequest(
	ctx context.Context,
	request ChatRequest,
	stream bool,
) (*http.Request, error) {

	reqBody := OpenAIChatRequest{
		Model:    request.Model,
		Messages: request.Messages,
		Stream:   stream,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		openAIChatCompletionsEndpoint,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return nil, err
	}

	httpReq.Header.Set(
		"Content-Type",
		"application/json",
	)

	httpReq.Header.Set(
		"Authorization",
		"Bearer "+request.APIKey,
	)

	if stream {
		httpReq.Header.Set(
			"Accept",
			"text/event-stream",
		)
	}

	return httpReq, nil
}
