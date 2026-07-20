package ai

import "github.com/harishnagaraju/astramind/internal/infrastructure/models"

type ChatRequest struct {
	Model    string
	APIKey   string
	Messages []models.Message
}
