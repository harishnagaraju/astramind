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

	req, err := o.buildRequest(
		context.Background(),
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

	httpReq, err := p.buildRequest(
		ctx,
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

	go p.readStream(
		resp.Body,
		stream,
	)

	return stream, nil
}

func (p *OpenAIProvider) chatCompletionsEndpoint() string {
	return p.baseURL + "/chat/completions"
}
