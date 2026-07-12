package chat

import "testing"

func TestHandleKnowledgeCommand(t *testing.T) {
	service := &Service{}

	handled, err := service.HandleKnowledgeCommand("/kb")
	if err != nil {
		t.Fatal(err)
	}

	if !handled {
		t.Fatal("expected command to be handled")
	}
}

func TestHandleKnowledgeCommandNonKB(t *testing.T) {
	service := &Service{}

	handled, err := service.HandleKnowledgeCommand("/search")
	if err != nil {
		t.Fatal(err)
	}

	if handled {
		t.Fatal("expected command not to be handled")
	}
}

func TestHandleKBImportInvalidArguments(t *testing.T) {
	service := &Service{}

	err := service.handleKBImport([]string{
		"/kb",
		"import",
	})

	if err == nil {
		t.Fatal("expected error")
	}
}

func TestHandleKnowledgeCommandRoot(t *testing.T) {
	service := &Service{}

	handled, err := service.HandleKnowledgeCommand("/kb")
	if err != nil {
		t.Fatal(err)
	}

	if !handled {
		t.Fatal("expected command to be handled")
	}
}

func TestHandleKnowledgeCommandInvalidSubcommand(t *testing.T) {

	service := &Service{}

	handled, err := service.HandleKnowledgeCommand("/kb foo")

	if !handled {
		t.Fatal("expected command to be handled")
	}

	if err == nil {
		t.Fatal("expected invalid command error")
	}
}
