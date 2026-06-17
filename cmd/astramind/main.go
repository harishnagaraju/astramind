package main
import "strings"
import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

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

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("===================================")
	fmt.Println("AstraMind v0.1.0")
	fmt.Println("Intelligent Conversations. Infinite Possibilities.")
	fmt.Println("Type 'exit' to quit")
	fmt.Println("===================================")

	for {

		fmt.Print("\nYou: ")

		userInput, _ := reader.ReadString('\n')
		userInput = strings.TrimSpace(userInput)

		if userInput == "exit" || userInput == "quit" {
			fmt.Println("Goodbye!")
			break
		}

		reply, err := askAI(apiKey, userInput)

		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		fmt.Println("\nAI:", reply)
	}
}

func askAI(apiKey, prompt string) (string, error) {

	reqBody := ChatRequest{
		Model: "gpt-4o-mini",
		Messages: []Message{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	jsonData, _ := json.Marshal(reqBody)

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
