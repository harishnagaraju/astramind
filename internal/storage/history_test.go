package storage

import (
	"os"
	"testing"

	"github.com/harishnagaraju/astramind/internal/testutil"
)

func TestSaveAndLoadHistory(t *testing.T) {

	session := "unit-test"

	expected := testutil.LoadConversation(
		t,
		"short",
	)

	err := SaveHistory(
		session,
		expected,
	)

	if err != nil {
		t.Fatalf(
			"SaveHistory failed: %v",
			err,
		)
	}

	actual, err := LoadHistory(
		session,
	)

	if err != nil {
		t.Fatalf(
			"LoadHistory failed: %v",
			err,
		)
	}

	if len(actual) != len(expected) {

		t.Fatalf(
			"Expected %d messages but got %d",
			len(expected),
			len(actual),
		)
	}

	testutil.AssertFileExists(
		t,
		"data/sessions/unit-test.json",
	)

	err = os.Remove(
		"data/sessions/unit-test.json",
	)

	if err != nil {
		t.Fatalf(
			"Cleanup failed: %v",
			err,
		)
	}
}

func TestLoadMissingSession(t *testing.T) {

	session := "session-does-not-exist"

	messages, err := LoadHistory(session)

	if err != nil {
		t.Fatalf(
			"Expected no error for missing session, got: %v",
			err,
		)
	}

	if len(messages) != 0 {
		t.Fatalf(
			"Expected empty history but got %d messages",
			len(messages),
		)
	}
}
