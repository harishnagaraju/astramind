package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// Embed generates a vector embedding for the given text using the
// OpenAI-compatible embeddings API.
func (o *OpenAIProvider) Embed(
	request EmbeddingRequest,
) ([]float32, error) {

	model := request.Model
	if model == "" {
		model = "text-embedding-3-small"
	}

	req, err := o.buildEmbeddingRequest(
		context.Background(),
		model,
		request,
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

	var result OpenAIEmbeddingResponse

	err = json.NewDecoder(
		resp.Body,
	).Decode(&result)

	if err != nil {
		return nil, err
	}

	if len(result.Data) == 0 {
		return nil, fmt.Errorf(
			"openai returned no embedding data",
		)
	}

	return result.Data[0].Embedding, nil
}

func (o *OpenAIProvider) embeddingsEndpoint() string {
	return o.baseURL + "/embeddings"
}

func (o *OpenAIProvider) buildEmbeddingRequest(
	ctx context.Context,
	model string,
	request EmbeddingRequest,
) (*http.Request, error) {

	reqBody := OpenAIEmbeddingRequest{
		Model: model,
		Input: request.Text,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		o.embeddingsEndpoint(),
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

	return httpReq, nil
}
