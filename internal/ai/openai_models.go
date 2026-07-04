package ai

import "github.com/harishnagaraju/astramind/internal/models"

const (
	openAIChatCompletionsEndpoint = "https://api.openai.com/v1/chat/completions"
)

type OpenAIChatRequest struct {
	Model    string           `json:"model"`
	Messages []models.Message `json:"messages"`
}

type OpenAIChatResponse struct {
	Choices []struct {
		Message models.Message `json:"message"`
	} `json:"choices"`
}
