package config

import "testing"

func TestVersionConstant(t *testing.T) {

	if Version == "" {
		t.Fatal("Version should not be empty")
	}
}

func TestMaxMessagesConstant(t *testing.T) {

	if MaxMessages <= 0 {
		t.Fatal("MaxMessages must be greater than zero")
	}
}

func TestHistoryFileConstant(t *testing.T) {

	if HistoryFile == "" {
		t.Fatal("HistoryFile should not be empty")
	}
}

func TestConfigStruct(t *testing.T) {

	cfg := Config{
		APIKey: "test-key",
		Model:  "gemma3:1b",
	}

	if cfg.APIKey != "test-key" {
		t.Fatal("APIKey not stored correctly")
	}

	if cfg.Model != "gemma3:1b" {
		t.Fatal("Model not stored correctly")
	}
}
