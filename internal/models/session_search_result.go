package models

// SessionSearchResult represents a search match found in a saved session.
type SessionSearchResult struct {
	// Session is the session name containing the match.
	Session string

	// Index is the zero-based message index.
	Index int

	// Role is the message role.
	Role string

	// Content is the message text.
	Content string
}
