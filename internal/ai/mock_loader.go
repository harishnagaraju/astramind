package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
)

func LoadMockData() ([]MockPair, error) {

	root, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	root = filepath.Join(root, "..", "..")

	file := filepath.Join(
		root,
		"tests",
		"testdata",
		"mock_ai.json",
	)

	data, err := os.ReadFile(file)
	if err != nil {
		return nil, err
	}

	var pairs []MockPair

	err = json.Unmarshal(data, &pairs)
	if err != nil {
		return nil, err
	}

	return pairs, nil
}
