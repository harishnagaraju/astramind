package ai

import (
	"encoding/json"
	"os"
	"path/filepath"
	"runtime"
)

func LoadMockData() ([]MockPair, error) {

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return nil, os.ErrNotExist
	}

	root := filepath.Join(
		filepath.Dir(filename),
		"..",
		"..",
		"..",
	)

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
