package testutil

import "github.com/harishnagaraju/astramind/internal/infrastructure/models"

func SampleConversation() []models.Message {

	return []models.Message{

		{
			Role:    "user",
			Content: "Hello AstraMind",
		},

		{
			Role:    "assistant",
			Content: "Hello Harish! How can I help you today?",
		},

		{
			Role:    "user",
			Content: "Explain Go programming.",
		},

		{
			Role:    "assistant",
			Content: "Go is an open-source programming language developed by Google.",
		},
	}
}
