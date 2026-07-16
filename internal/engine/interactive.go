package engine

import "github.com/harishnagaraju/astramind/internal/infrastructure/config"
import "github.com/harishnagaraju/astramind/internal/infrastructure/storage"
import "github.com/harishnagaraju/astramind/internal/infrastructure/models"
import "github.com/harishnagaraju/astramind/internal/infrastructure/ai"

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
)

func (a *App) runInteractive() error {

	var err error
	reader := bufio.NewReader(os.Stdin)

	a.runtime.Conversation, err = storage.LoadHistory(
		a.activeSession,
	)

	if err != nil {
		fmt.Println("Warning: could not load history:", err)
		a.runtime.Conversation = []models.Message{}
	}

	fmt.Printf(
		"Loaded %d messages from history.\n",
		len(a.runtime.Conversation),
	)

	fmt.Println("===================================")
	fmt.Printf("AstraMind %s\n", config.Version)
	fmt.Println("Intelligent Conversations. Infinite Possibilities.")
	fmt.Println("Type '/help' for commands")
	fmt.Println("===================================")

	for {

		fmt.Printf("\n[%s] You: ", a.activeSession)

		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input Error:", err)
			continue
		}

		userInput = strings.TrimSpace(userInput)

		handled, err := a.dispatcher.Execute(userInput)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if handled {
			continue
		}

		if userInput == "" {
			continue
		}

		switch userInput {

		case "exit", "quit":
			storage.SaveHistory(a.activeSession, a.runtime.Conversation)
			fmt.Println("Goodbye!")
			return nil

		}
		// Create temporary conversation
		// Do NOT save until API succeeds.
		updatedConversation := append(a.runtime.Conversation, models.Message{
			Role:    "user",
			Content: userInput,
		})

		reply, streamed, err := a.deps.ChatService.Chat(
			context.Background(),
			os.Stdout,
			ai.ChatRequest{
				Model:    a.model,
				APIKey:   a.apiKey,
				Messages: updatedConversation,
			},
		)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Save user message only after successful API response
		a.runtime.Conversation = updatedConversation

		// Save assistant response
		a.runtime.Conversation = append(a.runtime.Conversation, models.Message{
			Role:    "assistant",
			Content: reply,
		})

		if err := storage.SaveHistory(a.activeSession, a.runtime.Conversation); err != nil {
			fmt.Println("Warning: failed to save conversation:", err)
		}

		// Keep memory bounded
		if len(a.runtime.Conversation) > config.MaxMessages {
			a.runtime.Conversation = a.runtime.Conversation[len(a.runtime.Conversation)-config.MaxMessages:]
		}

		if !streamed {
			fmt.Println("\nAI:", reply)
		}
	}
}
