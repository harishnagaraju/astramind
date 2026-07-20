package models

// SearchResult represents a single search match within a conversation.
type SearchResult struct {
	// Index is the zero-based index of the matching message
	// within the conversation history.
	Index int

	// Role is the message role ("user", "assistant", etc.).
	Role string

	// Content contains the matching message text.
	Content string
}
