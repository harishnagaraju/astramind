package storage

import (
	"os"
	"testing"

	"github.com/harishnagaraju/astramind/internal/testutil"
)

func TestExportMarkdown(t *testing.T) {

	session := "unit-test"

	err := ExportMarkdown(
		session,
		testutil.SampleConversation(),
	)

	if err != nil {
		t.Fatalf(
			"ExportMarkdown failed: %v",
			err,
		)
	}

	testutil.AssertFileExists(
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
