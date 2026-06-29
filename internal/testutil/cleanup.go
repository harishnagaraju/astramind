package testutil

import (
	"os"
	"testing"
)

func RemoveDirectory(
	t *testing.T,
	dir string,
) {

	t.Helper()

	err := os.RemoveAll(dir)
	if err != nil {
		t.Fatalf("Failed to remove directory: %v", err)
	}
}
