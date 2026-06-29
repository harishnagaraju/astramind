package testutil

import "testing"

func TestCreateTempDirectory(t *testing.T) {

	dir := CreateTempDirectory(t)

	if dir == "" {
		t.Fatal("Temporary directory was not created")
	}

	RemoveDirectory(t, dir)
}
