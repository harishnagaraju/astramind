package config

const (
	MaxMessages = 20
	HistoryFile = "data/chat_history.json"
)

type Config struct {
	APIKey string
	Model  string
}