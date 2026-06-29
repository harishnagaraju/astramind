package ai

import "github.com/harishnagaraju/astramind/internal/models"

type ChatRequest struct {
	Model    string
	APIKey   string
	Messages []models.Message
}
