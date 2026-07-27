package config

const (
	Version     = "v0.9.2"
	MaxMessages = 20
	HistoryFile = "data/chat_history.json"
)

type Config struct {
	APIKey string
	Model  string
}
