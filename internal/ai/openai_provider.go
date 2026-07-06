package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
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

	endpoint := o.chatCompletionsEndpoint()

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
		"HTTP-Referer",
		"https://github.com/harishnagaraju/astramind",
	)

	req.Header.Set(
		"X-Title",
		"AstraMind",
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

		return "", handleAPIError(
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

	endpoint := p.chatCompletionsEndpoint()

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
		"HTTP-Referer",
		"https://github.com/harishnagaraju/astramind",
	)

	httpReq.Header.Set(
		"X-Title",
		"AstraMind",
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

		return nil, handleAPIError(
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

func (p *OpenAIProvider) chatCompletionsEndpoint() string {
	return p.baseURL + "/chat/completions"
}
