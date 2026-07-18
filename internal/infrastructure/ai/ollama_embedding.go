package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

// Embed generates a vector embedding for the given text using the
// Ollama embeddings API.
func (o *OllamaProvider) Embed(
	request EmbeddingRequest,
) ([]float32, error) {

	model := request.Model
	if model == "" {
		model = "nomic-embed-text"
	}

	req, err := buildOllamaEmbeddingRequest(
		o.baseURL,
		model,
		request.Text,
	)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		var body bytes.Buffer
		body.ReadFrom(resp.Body)

		return nil, handleAPIError(
			resp.StatusCode,
			body.String(),
		)
	}

	var result OllamaEmbeddingResponse

	err = json.NewDecoder(
		resp.Body,
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	if len(result.Embedding) == 0 {
		return nil, fmt.Errorf(
			"ollama returned an empty embedding",
		)
	}

	return result.Embedding, nil
}

// buildOllamaEmbeddingRequest builds an HTTP request for the Ollama
// embeddings API.
func buildOllamaEmbeddingRequest(
	baseURL string,
	model string,
	text string,
) (*http.Request, error) {

	ollamaReq := OllamaEmbeddingRequest{
		Model:  model,
		Prompt: text,
	}

	body, err := json.Marshal(ollamaReq)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal ollama embedding request: %w",
			err,
		)
	}

	req, err := http.NewRequest(
		http.MethodPost,
		baseURL+"/api/embeddings",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf(
			"failed to create ollama embedding request: %w",
			err,
		)
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	return req, nil
}
