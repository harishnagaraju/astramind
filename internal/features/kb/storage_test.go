package kb

import (
	"testing"
)

func TestNewJSONStorage(t *testing.T) {
	dir := t.TempDir()

	storage := NewJSONStorage(dir)

	if storage == nil {
		t.Fatal("expected storage instance")
	}
}
