package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
)

// OllamaProvider implements the Provider interface
// for locally hosted Ollama models.
type OllamaProvider struct {
	baseURL string
	model   string
}

func (o *OllamaProvider) Name() string {
	return "Ollama"
}

func (o *OllamaProvider) Chat(
	request ChatRequest,
) (string, error) {

	model := o.model
	if model == "" {
		model = "llama3"
	}

	req, err := buildOllamaRequest(
		o.baseURL,
		model,
		request,
		false,
	)
	if err != nil {
		return "", err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		var body bytes.Buffer
		body.ReadFrom(resp.Body)

		return "", handleAPIError(
			resp.StatusCode,
			body.String(),
		)
	}

	var result OllamaChatResponse

	err = json.NewDecoder(
		resp.Body,
	).Decode(&result)

	if err != nil {
		return "", err
	}

	if result.Message.Content == "" {
		return "No response", nil
	}

	return result.Message.Content, nil
}

func (p *OllamaProvider) Stream(
	ctx context.Context,
	request ChatRequest,
) (Stream, error) {

	stream := &ollamaStream{
		events: make(chan StreamEvent),
	}

	model := p.model
	if model == "" {
		model = "llama3"
	}

	httpReq, err := buildOllamaRequest(
		p.baseURL,
		model,
		request,
		true,
	)

	if err != nil {
		return nil, err
	}

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

	go func() {
		defer resp.Body.Close()
		readOllamaStream(
			resp.Body,
			stream.events,
		)
	}()

	return stream, nil
}
