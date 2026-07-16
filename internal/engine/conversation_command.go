package engine

import (
	"fmt"

	"github.com/harishnagaraju/astramind/internal/infrastructure/config"
	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
	"github.com/harishnagaraju/astramind/internal/infrastructure/storage"
)

type conversationCommand struct{}

func (c *conversationCommand) Execute(
	app *App,
	input string,
) (bool, error) {

	switch input {

	case "/clear":

		app.runtime.Conversation = []models.Message{}

		err := storage.SaveHistory(
			app.activeSession,
			app.runtime.Conversation,
		)

		if err != nil {
			fmt.Println(
				"Error clearing history:",
				err,
			)
		} else {
			fmt.Println(
				"Conversation memory cleared.",
			)
		}

		return true, nil

	case "/history":

		if len(app.runtime.Conversation) == 0 {

			fmt.Println(
				"No conversation history.",
			)

			return true, nil
		}

		fmt.Println("\nConversation History:")

		for i, msg := range app.runtime.Conversation {

			fmt.Printf(
				"%d. [%s] %s\n",
				i+1,
				msg.Role,
				msg.Content,
			)
		}

		return true, nil

	case "/stats":

		userCount := 0
		assistantCount := 0

		for _, msg := range app.runtime.Conversation {

			switch msg.Role {

			case "user":
				userCount++

			case "assistant":
				assistantCount++
			}
		}

		fmt.Println()
		fmt.Println("Session Statistics")
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
			len(app.runtime.Conversation),
		)

		fmt.Printf(
			"Current Model: %s\n",
			app.model,
		)

		fmt.Printf(
			"Memory Limit: %d\n",
			config.MaxMessages,
		)

		return true, nil
	}

	return false, nil
}
