package ai

import (
	"bytes"
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
