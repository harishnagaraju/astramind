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

	reader := bufio.NewReader(os.Stdin)

	conversation, err := storage.LoadHistory(a.activeSession)

	if err != nil {
		fmt.Println("Warning: could not load history:", err)
		conversation = []models.Message{}
	}

	fmt.Printf(
		"Loaded %d messages from history.\n",
		len(conversation),
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

			conversation = messages

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

			conversation = []models.Message{}

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

			results := storage.SearchMessages(conversation, query)

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
			storage.SaveHistory(a.activeSession, conversation)
			fmt.Println("Goodbye!")
			return nil

		case "/clear":
			conversation = []models.Message{}

			err := storage.SaveHistory(a.activeSession, conversation)

			if err != nil {
				fmt.Println("Error clearing history:", err)
			} else {
				fmt.Println("Conversation memory cleared.")
			}
			continue

		case "/provider":

			fmt.Println()

			fmt.Println("Current AI Provider")
			fmt.Println("-------------------")

			fmt.Printf(
				"Provider : %s\n",
				a.deps.ProviderManager.ProviderName(),
			)

			fmt.Printf(
				"Model    : %s\n",
				a.model,
			)

			fmt.Println()

			continue

		case "/export", "/export txt", "/export md":

			if len(conversation) == 0 {
				fmt.Println("Nothing to export.")
				continue
			}

			switch userInput {

			case "/export", "/export txt":

				err := storage.ExportSession(
					a.activeSession,
					conversation,
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
					conversation,
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

			continue

		case "/sessions":

			sessions, err := storage.ListSessions()

			if err != nil {

				fmt.Println(
					"Error loading sessions:",
					err,
				)

				continue
			}

			fmt.Println("\nAvailable Sessions")
			fmt.Println("------------------")

			if len(sessions) == 0 {

				fmt.Println(
					"No sessions found.",
				)

				continue
			}

			for _, session := range sessions {

				fmt.Println(session)
			}

			continue

		case "/history":

			if len(conversation) == 0 {
				fmt.Println("No conversation history.")
				continue
			}

			fmt.Println("\nConversation History:")

			for i, msg := range conversation {
				fmt.Printf(
					"%d. [%s] %s\n",
					i+1,
					msg.Role,
					msg.Content,
				)
			}

			continue

		case "/stats":

			userCount := 0
			assistantCount := 0

			for _, msg := range conversation {

				switch msg.Role {

				case "user":
					userCount++

				case "assistant":
					assistantCount++
				}
			}

			fmt.Println("\nSession Statistics")
			fmt.Println("------------------")

			fmt.Printf(
				"User Messages: %d\n",
				userCount,
			)

			fmt.Printf(
				"Assistant Messages: %d\n",
				assistantCount,
			)

			fmt.Printf(
				"Memory Entries: %d\n",
				len(conversation),
			)

			fmt.Printf(
				"Current Model: %s\n",
				a.model,
			)

			continue
		}

		// Create temporary conversation
		// Do NOT save until API succeeds.
		tempConversation := append(conversation, models.Message{
			Role:    "user",
			Content: userInput,
		})

		reply, streamed, err := a.deps.ChatService.Chat(
			context.Background(),
			os.Stdout,
			ai.ChatRequest{
				Model:    a.model,
				APIKey:   a.apiKey,
				Messages: tempConversation,
			},
		)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Save user message only after successful API response
		conversation = tempConversation

		// Save assistant response
		conversation = append(conversation, models.Message{
			Role:    "assistant",
			Content: reply,
		})

		if err := storage.SaveHistory(a.activeSession, conversation); err != nil {
			fmt.Println("Warning: failed to save conversation:", err)
		}

		// Keep memory bounded
		if len(conversation) > config.MaxMessages {
			conversation = conversation[len(conversation)-config.MaxMessages:]
		}

		if !streamed {
			fmt.Println("\nAI:", reply)
		}
	}
}
