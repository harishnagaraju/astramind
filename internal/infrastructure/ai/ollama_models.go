package ai

// OllamaChatRequest represents a chat request sent to the Ollama API.
type OllamaChatRequest struct {
	Model    string              `json:"model"`
	Messages []OllamaChatMessage `json:"messages"`
	Stream   bool                `json:"stream"`
	Options  *OllamaOptions      `json:"options,omitempty"`
}

// OllamaOptions configures generation parameters for the Ollama API.
type OllamaOptions struct {
	// NumCtx sets the context window size, in tokens. Ollama's own
	// default (commonly 2048) is too small for RAG prompts, which
	// combine system instructions, several retrieved chunks, the
	// question, and the answer all within the same window - once
	// that combined total exceeds the window, Ollama truncates the
	// response mid-generation rather than erroring.
	NumCtx int `json:"num_ctx,omitempty"`
}

// OllamaChatMessage represents a single chat message.
type OllamaChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// OllamaChatResponse represents the response returned by the Ollama API.
type OllamaChatResponse struct {
	Model              string            `json:"model"`
	CreatedAt          string            `json:"created_at,omitempty"`
	Message            OllamaChatMessage `json:"message"`
	Done               bool              `json:"done"`
	TotalDuration      int64             `json:"total_duration,omitempty"`
	LoadDuration       int64             `json:"load_duration,omitempty"`
	PromptEvalCount    int               `json:"prompt_eval_count,omitempty"`
	PromptEvalDuration int64             `json:"prompt_eval_duration,omitempty"`
	EvalCount          int               `json:"eval_count,omitempty"`
	EvalDuration       int64             `json:"eval_duration,omitempty"`
}
