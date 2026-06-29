package storage

import (
	"os"
	"testing"

	"github.com/harishnagaraju/astramind/internal/testutil"
)

func TestSessionWorkflow(t *testing.T) {

	session := "integration-test"

	sessionFile := "data/sessions/" + session + ".json"
	txtFile := "exports/" + session + ".txt"
	mdFile := "exports/" + session + ".md"

	// Cleanup before starting
	_ = os.Remove(sessionFile)
	_ = os.Remove(txtFile)
	_ = os.Remove(mdFile)

	// Sample conversation
	expected := testutil.LoadConversation(
		t,
		"long",
	)

	// Save session
	err := SaveHistory(session, expected)
	if err != nil {
		t.Fatalf("SaveHistory failed: %v", err)
	}

	testutil.AssertFileExists(t, sessionFile)

	// Load session
	actual, err := LoadHistory(session)
	if err != nil {
		t.Fatalf("LoadHistory failed: %v", err)
	}

	if len(actual) != len(expected) {
		t.Fatalf(
			"Expected %d messages but got %d",
			len(expected),
			len(actual),
		)
	}

	// Export TXT
	err = ExportSession(session, actual)
	if err != nil {
		t.Fatalf("ExportSession failed: %v", err)
	}

	testutil.AssertFileExists(t, txtFile)

	// Export Markdown
	err = ExportMarkdown(session, actual)
	if err != nil {
		t.Fatalf("ExportMarkdown failed: %v", err)
	}

	testutil.AssertFileExists(t, mdFile)

	// Delete session
	err = DeleteSession(session)
	if err != nil {
		t.Fatalf("DeleteSession failed: %v", err)
	}

	// Verify session file is deleted
	_, err = os.Stat(sessionFile)
	if !os.IsNotExist(err) {
		t.Fatal("Session file still exists after deletion")
	}

	// Cleanup exported files
	_ = os.Remove(txtFile)
	_ = os.Remove(mdFile)
}
