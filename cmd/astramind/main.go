package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const MaxMessages = 20

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load .env file")
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
	model := os.Getenv("OPENAI_MODEL")

	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not found in .env")
		return
	}

	if model == "" {
		model = "gpt-4o-mini"
	}

	reader := bufio.NewReader(os.Stdin)

	var conversation []Message

	fmt.Println("===================================")
	fmt.Println("AstraMind v0.2.1")
	fmt.Println("Intelligent Conversations. Infinite Possibilities.")
	fmt.Println("Type '/help' for commands")
	fmt.Println("===================================")

	for {

		fmt.Print("\nYou: ")

		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input Error:", err)
			continue
		}

		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			continue
		}

		switch userInput {

		case "exit", "quit":
			fmt.Println("Goodbye!")
			return

		case "/help":
			fmt.Println("\nAvailable Commands:")
			fmt.Println("/help      - Show help")
			fmt.Println("/history   - Show conversation history")
			fmt.Println("/clear     - Clear conversation memory")
			fmt.Println("/stats     - Show session statistics")
			fmt.Println("exit       - Exit AstraMind")
			fmt.Println("quit       - Exit AstraMind")
			continue

		case "/clear":
			conversation = nil
			fmt.Println("Conversation memory cleared.")
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
				model,
			)

			continue
		}

		// Create temporary conversation
		// Do NOT save until API succeeds.
		tempConversation := append(conversation, Message{
			Role:    "user",
			Content: userInput,
		})

		reply, err := askAI(
			apiKey,
			model,
			tempConversation,
		)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		// Save user message only after successful API response
		conversation = tempConversation

		// Save assistant response
		conversation = append(conversation, Message{
			Role:    "assistant",
			Content: reply,
		})

		// Keep memory bounded
		if len(conversation) > MaxMessages {
			conversation = conversation[len(conversation)-MaxMessages:]
		}

		fmt.Println("\nAI:", reply)
	}
}

func askAI(
	apiKey string,
	model string,
	messages []Message,
) (string, error) {

	reqBody := ChatRequest{
		Model:    model,
		Messages: messages,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest(
		"POST",
		"https://api.openai.com/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {

		var body bytes.Buffer
		body.ReadFrom(resp.Body)

		return "", fmt.Errorf(
			"API Error (%d): %s",
			resp.StatusCode,
			body.String(),
		)
	}

	var result ChatResponse

	err = json.NewDecoder(resp.Body).Decode(&result)
	if err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "No response", nil
	}

	return result.Choices[0].Message.Content, nil
}