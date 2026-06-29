package config

const (
	Version     = "v0.4.1-dev"
	MaxMessages = 20
	HistoryFile = "data/chat_history.json"
)

type Config struct {
	APIKey string
	Model  string
}
