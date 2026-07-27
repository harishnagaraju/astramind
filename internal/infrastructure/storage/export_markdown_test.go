package storage

import (
	"os"
	"testing"

	"github.com/harishnagaraju/astramind/internal/utilityforunittest"
)

func TestExportMarkdown(t *testing.T) {

	session := "unit-test"

	err := ExportMarkdown(
		session,
		utilityforunittest.LoadConversation(
			t,
			"short",
		),
	)

	if err != nil {
		t.Fatalf(
			"ExportMarkdown failed: %v",
			err,
		)
	}

	utilityforunittest.AssertFileExists(
		t,
		"exports/unit-test.md",
	)

	err = os.Remove(
		"exports/unit-test.md",
	)

	if err != nil {
		t.Fatalf(
			"Cleanup failed: %v",
			err,
		)
	}
}
