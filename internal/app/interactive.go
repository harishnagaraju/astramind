package app

import "github.com/harishnagaraju/astramind/internal/config"
import "github.com/harishnagaraju/astramind/internal/storage"
import "github.com/harishnagaraju/astramind/internal/models"
import "github.com/harishnagaraju/astramind/internal/ai"
import "github.com/harishnagaraju/astramind/internal/renderer"

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
		if strings.HasPrefix(
			userInput,
			"/load ",
		) {

			sessionName := strings.TrimSpace(
				strings.TrimPrefix(
					userInput,
					"/load ",
				),
			)

			if sessionName == "" {

				fmt.Println(
					"Usage: /load <session-name>",
				)

				continue
			}

			if !storage.SessionExists(
				sessionName,
			) {

				fmt.Printf(
					"Session '%s' does not exist.\n",
					sessionName,
				)

				continue
			}

			messages, err :=
				storage.LoadHistory(
					sessionName,
				)

			if err != nil {

				fmt.Println(
					"Error:",
					err,
				)

				continue
			}

			a.activeSession = sessionName

			a.runtime.Conversation = messages

			fmt.Printf(
				"Loaded session: %s\n",
				sessionName,
			)

			continue
		}

		if strings.HasPrefix(
			userInput,
			"/new ",
		) {

			sessionName := strings.TrimSpace(
				strings.TrimPrefix(
					userInput,
					"/new ",
				),
			)

			if sessionName == "" {

				fmt.Println(
					"Usage: /new <session-name>",
				)

				continue
			}

			err := storage.CreateSession(
				sessionName,
			)

			if err != nil {

				fmt.Println(
					"Error:",
					err,
				)

				continue
			}

			a.activeSession = sessionName

			a.runtime.Conversation = []models.Message{}

			fmt.Printf(
				"Created and switched to session: %s\n",
				sessionName,
			)

			continue
		}

		if strings.HasPrefix(
			userInput,
			"/delete ",
		) {

			sessionName := strings.TrimSpace(
				strings.TrimPrefix(
					userInput,
					"/delete ",
				),
			)

			if sessionName == "" {

				fmt.Println(
					"Usage: /delete <session-name>",
				)

				continue
			}

			if sessionName == "default" {

				fmt.Println(
					"Default session cannot be deleted.",
				)

				continue
			}

			if sessionName == a.activeSession {

				fmt.Println(
					"Cannot delete active session.",
				)

				continue
			}

			err := storage.DeleteSession(
				sessionName,
			)

			if err != nil {

				fmt.Println(
					"Error:",
					err,
				)

				continue
			}

			fmt.Printf(
				"Deleted session: %s\n",
				sessionName,
			)

			continue
		}

		if strings.HasPrefix(userInput, "/search ") || userInput == "/search" {

			query := strings.TrimSpace(
				strings.TrimPrefix(userInput, "/search"),
			)

			if query == "" {
				fmt.Println("Usage: /search <text>")
				continue
			}

			results := storage.SearchMessages(a.runtime.Conversation, query)

			if len(results) == 0 {
				fmt.Println("No matches found.")
				continue
			}

			renderer.RenderSearchResults(results)

			continue
		}

		if strings.HasPrefix(userInput, "/searchall ") || userInput == "/searchall" {

			query := strings.TrimSpace(
				strings.TrimPrefix(userInput, "/searchall"),
			)

			if query == "" {
				fmt.Println("Usage: /searchall <text>")
				continue
			}

			results, err := storage.SearchAllSessions(query)
			if err != nil {
				fmt.Println("Search failed:", err)
				continue
			}

			if len(results) == 0 {
				fmt.Println("No matches found.")
				continue
			}

			renderer.RenderSessionSearchResults(results)

			continue
		}

		switch userInput {

		case "exit", "quit":
			storage.SaveHistory(a.activeSession, a.runtime.Conversation)
			fmt.Println("Goodbye!")
			return nil

		case "/export", "/export txt", "/export md":

			if len(a.runtime.Conversation) == 0 {
				fmt.Println("Nothing to export.")
				continue
			}

			switch userInput {

			case "/export", "/export txt":

				err := storage.ExportSession(
					a.activeSession,
					a.runtime.Conversation,
				)

				if err != nil {
					fmt.Println("Export failed:", err)
					continue
				}

				fmt.Printf(
					"Session exported to exports/%s.txt\n",
					a.activeSession,
				)

			case "/export md":

				err := storage.ExportMarkdown(
					a.activeSession,
					a.runtime.Conversation,
				)

				if err != nil {
					fmt.Println("Export failed:", err)
					continue
				}

				fmt.Printf(
					"Session exported to exports/%s.md\n",
					a.activeSession,
				)
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
}
