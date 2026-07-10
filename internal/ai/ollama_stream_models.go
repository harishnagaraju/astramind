package ai

type OllamaStreamResponse struct {
	Message OllamaStreamMessage `json:"message"`
	Done    bool                `json:"done"`
}

type OllamaStreamMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}
