package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

type OpenAIProvider struct {
	baseURL string
}

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

	endpoint := o.baseURL + "/chat/completions"

	req, err := http.NewRequest(
		http.MethodPost,
		endpoint,
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
	request ChatRequest,
) (Stream, error) {

	stream := &openAIStream{
		events: make(chan StreamEvent),
	}

	reqBody := OpenAIChatRequest{
		Model:    request.Model,
		Messages: request.Messages,
		Stream:   true,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	endpoint := p.baseURL + "/chat/completions"

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		endpoint,
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

	httpReq.Header.Set(
		"Accept",
		"text/event-stream",
	)

	client := &http.Client{}

	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {

		defer resp.Body.Close()

		var body bytes.Buffer
		body.ReadFrom(resp.Body)

		return nil, fmt.Errorf(
			"API Error (%d): %s",
			resp.StatusCode,
			body.String(),
		)
	}

	go p.readStream(
		resp.Body,
		stream,
	)

	_ = resp

	return stream, nil
}
