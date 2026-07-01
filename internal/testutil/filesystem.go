package testutil

import (
	"os"
	"testing"
)

func CreateTempDirectory(t *testing.T) string {

	t.Helper()

	dir, err := os.MkdirTemp("", "astramind-test-*")
	if err != nil {
		t.Fatalf("Failed to create temporary directory: %v", err)
	}

	return dir
}
