package storage

import (
	"os"
	"testing"

	"github.com/harishnagaraju/astramind/internal/testutil"
)

func TestExportTXT(t *testing.T) {

	session := "unit-test"

	err := ExportSession(
		session,
		testutil.LoadConversation(
			t,
			"short",
		),
	)

	if err != nil {
		t.Fatalf(
			"ExportSession failed: %v",
			err,
		)
	}

	testutil.AssertFileExists(
		t,
		"exports/unit-test.txt",
	)

	err = os.Remove(
		"exports/unit-test.txt",
	)

	if err != nil {
		t.Fatalf(
			"Cleanup failed: %v",
			err,
		)
	}
}
