package ai

import "github.com/harishnagaraju/astramind/internal/models"

const (
	openAIChatCompletionsEndpoint = "https://api.openai.com/v1/chat/completions"
)

type OpenAIChatRequest struct {
	Model    string           `json:"model"`
	Messages []models.Message `json:"messages"`
	Stream   bool             `json:"stream,omitempty"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message models.Message `json:"message"`
	} `json:"choices"`
}

type OpenAIStreamResponse struct {
	Choices []OpenAIStreamChoice `json:"choices"`
}

type OpenAIStreamChoice struct {
	Delta OpenAIStreamDelta `json:"delta"`
}

type OpenAIStreamDelta struct {
	Content string `json:"content,omitempty"`
	Role    string `json:"role,omitempty"`
}
