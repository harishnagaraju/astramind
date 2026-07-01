package testutil

import "testing"

func TestLoadConversation(t *testing.T) {

	messages := LoadConversation(
		t,
		"short",
	)

	if len(messages) != 2 {
		t.Fatalf(
			"Expected 2 messages but got %d",
			len(messages),
		)
	}

	if messages[0].Role != "user" {
		t.Fatal("Unexpected first message role")
	}

	if messages[1].Role != "assistant" {
		t.Fatal("Unexpected second message role")
	}
}
