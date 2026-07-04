package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

var openAIEndpoint = "https://api.openai.com/v1/chat/completions"

type OpenAIProvider struct{}

func (o *OpenAIProvider) Name() string {
	return "OpenAI"
}

func (o *OpenAIProvider) Chat(
	request ChatRequest,
) (string, error) {

	reqBody := OpenAIChatRequest{
		Model:    request.Model,
		Messages: request.Messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		openAIEndpoint,
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	req.Header.Set(
		"Authorization",
		"Bearer "+request.APIKey,
	)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		var body bytes.Buffer
		body.ReadFrom(resp.Body)

		responseBody := body.String()

		if resp.StatusCode == 429 &&
			strings.Contains(
				responseBody,
				"insufficient_quota",
			) {

			return "", fmt.Errorf(
				"OpenAI quota exceeded.",
			)
		}

		return "", fmt.Errorf(
			"API Error (%d): %s",
			resp.StatusCode,
			responseBody,
		)
	}

	var result OpenAIChatResponse

	err = json.NewDecoder(
		resp.Body,
	).Decode(&result)

	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "No response", nil
	}

	return result.
			Choices[0].
			Message.
			Content,
		nil
}

func (p *OpenAIProvider) Stream(
	ctx context.Context,
	req ChatRequest,
) (Stream, error) {

	stream := &openAIStream{
		events: make(chan StreamEvent),
	}

	go func() {
		defer close(stream.events)

		stream.events <- StreamEvent{
			Type: StreamEventDone,
		}
	}()

	return stream, nil
}
