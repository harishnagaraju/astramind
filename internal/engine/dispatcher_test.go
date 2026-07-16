package engine

import (
	"fmt"
	"testing"
)

func TestNewDispatcher(t *testing.T) {

	app := &App{}

	dispatcher := newDispatcher(app)

	if dispatcher == nil {
		t.Fatal("Dispatcher should not be nil")
	}

	if dispatcher.app != app {
		t.Fatal("Dispatcher should store the App instance")
	}

	if dispatcher.commands == nil {
		t.Fatal("Dispatcher command list should not be nil")
	}

	if len(dispatcher.commands) != 5 {
		t.Fatalf(
			"Expected 5 registered commands but got %d",
			len(dispatcher.commands),
		)
	}
}

func TestDispatcherRegistrationOrder(t *testing.T) {

	app := &App{}

	dispatcher := newDispatcher(app)

	if len(dispatcher.commands) != 5 {
		t.Fatalf("expected 5 commands, got %d", len(dispatcher.commands))
	}

	expected := []string{
		"*engine.builtinCommand",
		"*engine.sessionCommand",
		"*engine.conversationCommand",
		"*engine.searchCommand",
		"*engine.exportCommand",
	}

	for i, cmd := range dispatcher.commands {
		actual := fmt.Sprintf("%T", cmd)
		if actual != expected[i] {
			t.Fatalf(
				"command %d: expected %s, got %s",
				i,
				expected[i],
				actual,
			)
		}
	}
}
