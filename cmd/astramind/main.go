package main
import "github.com/harishnagaraju/astramind/internal/config"
import "github.com/harishnagaraju/astramind/internal/storage"
import "github.com/harishnagaraju/astramind/internal/models"
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

/* func LoadHistory() ([]models.Message, error)
func SaveHistory(messages []models.Message) error */

var conversation []models.Message

/*const MaxMessages = 20*/

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []models.Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message models.Message `json:"message"`
	} `json:"choices"`
}

func main() {
	
	activeSession := "default"

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Could not load .env file", err)
		return
	}

	apiKey := os.Getenv("OPENAI_API_KEY")
        /*
	if apiKey == "your_api_key_here" {
   	    fmt.Println("Please update OPENAI_API_KEY in .env")
	    return
	}
        */

	model := os.Getenv("OPENAI_MODEL")

	if apiKey == "" {
		fmt.Println("OPENAI_API_KEY not found in .env")
		return
	}

	if model == "" {
		model = "gpt-4o-mini"
	}

	reader := bufio.NewReader(os.Stdin)

	conversation, err := storage.LoadHistory(activeSession)
       

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

		fmt.Printf("\n[%s] You: ", activeSession)

		userInput, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Input Error:", err)
			continue
		}

		userInput = strings.TrimSpace(userInput)

		if userInput == "" {
			continue
		}
		if strings.HasPrefix(
			userInput,
			"/load ",
			){

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

			activeSession = sessionName

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
			){

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

			activeSession = sessionName

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
			){

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

			if sessionName == activeSession {

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
		
		switch userInput {

		case "exit", "quit":
			storage.SaveHistory(activeSession, conversation)
			fmt.Println("Goodbye!")
			return

		case "/help":
			fmt.Println("\nAvailable Commands:")
			fmt.Println("/about     - About AstraMind")
			fmt.Println("/help      - Show help")
			fmt.Println("/history   - Show conversation history")
			fmt.Println("/clear     - Clear conversation memory")
			fmt.Println("/stats     - Show session statistics")
			fmt.Println("/config    - Show configuration")
			fmt.Println("/export    - Export session (TXT)")
			fmt.Println("/export md - Export session (Markdown)")
			fmt.Println("/sessions  - List sessions")
			fmt.Println("/new <name> - Create session")
			fmt.Println("/load <name> - Load session")
			fmt.Println("/delete <name> - Delete session")
			fmt.Println("/export    - Export session")
			fmt.Println("exit       - Exit AstraMind")
			fmt.Println("quit       - Exit AstraMind")
			continue

		case "/clear":
			conversation = []models.Message{}

			err := storage.SaveHistory(activeSession, conversation)

			if err != nil {
    				fmt.Println("Error clearing history:", err)
			} else {
    				fmt.Println("Conversation memory cleared.")
			}
			continue
			
		case "/config":

			fmt.Println("\nCurrent Configuration")
			fmt.Println("---------------------")

			fmt.Printf(
				"Model: %s\n",
				model,
			)

			fmt.Printf(
				"Max Messages: %d\n",
				config.MaxMessages,
			)

			fmt.Printf(
				"History Enabled: %t\n",
				true,
			)

			fmt.Printf(
				"History File: %s\n",
				config.HistoryFile,
			)

			continue
		
		case "/about":

			fmt.Println("\nAstraMind")
			fmt.Println("---------")

			fmt.Printf(
				"Version: %s\n",
				config.Version,
			)

			fmt.Println("\nFeatures:")

			fmt.Println("✓ Conversation Memory")
			fmt.Println("✓ Persistent History")
			fmt.Println("✓ Session Statistics")
			fmt.Println("✓ Configuration Display")

			fmt.Printf(
				"\nModel: %s\n",
				model,
			)
			fmt.Println("Author: Harish Nagaraju")
			fmt.Println("Company: RK Consulting")
			
			fmt.Println(
				"Repository: github.com/harishnagaraju/astramind",
			)

			continue
		
		case "/export", "/export txt", "/export md":

				if len(conversation) == 0 {
					fmt.Println("Nothing to export.")
					continue
				}

				switch userInput {

				case "/export", "/export txt":

					err := storage.ExportSession(
						activeSession,
						conversation,
					)

					if err != nil {
						fmt.Println("Export failed:", err)
						continue
					}

					fmt.Printf(
						"Session exported to exports/%s.txt\n",
						activeSession,
					)

				case "/export md":

					err := storage.ExportMarkdown(
						activeSession,
						conversation,
					)

					if err != nil {
						fmt.Println("Export failed:", err)
						continue
					}

					fmt.Printf(
						"Session exported to exports/%s.md\n",
						activeSession,
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
				model,
			)

			continue
		}

		// Create temporary conversation
		// Do NOT save until API succeeds.
		tempConversation := append(conversation, models.Message{
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
		conversation = append(conversation, models.Message{
			Role:    "assistant",
			Content: reply,
		})

		// Keep memory bounded
		if len(conversation) > config.MaxMessages {
			conversation = conversation[len(conversation)-config.MaxMessages:]
		}

		fmt.Println("\nAI:", reply)
	}
}

func askAI(
	apiKey string,
	model string,
	messages []models.Message,
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

		responseBody := body.String()

		if resp.StatusCode == 429 &&
			strings.Contains(responseBody, "insufficient_quota") {

			return "", fmt.Errorf(
				"OpenAI quota exceeded.\n\n" +
					"Please check:\n" +
					"- Billing settings\n" +
					"- Usage limits\n" +
					"- Available credits\n\n" +
					"https://platform.openai.com/usage",
			)
		}

		return "", fmt.Errorf(
			"API Error (%d): %s",
			resp.StatusCode,
			responseBody,
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