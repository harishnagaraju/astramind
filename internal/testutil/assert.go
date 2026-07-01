package testutil

import (
	"os"
	"testing"
)

func AssertFileExists(
	t *testing.T,
	file string,
) {

	t.Helper()

	_, err := os.Stat(file)

	if err != nil {
		t.Fatalf("Expected file does not exist: %s", file)
	}
}
