package kb

import "testing"

func TestNewManager(t *testing.T) {
	storage := NewJSONStorage(t.TempDir())

	manager := NewManager(storage)

	if manager == nil {
		t.Fatal("expected manager")
	}

	if manager.storage == nil {
		t.Fatal("expected storage")
	}
}
