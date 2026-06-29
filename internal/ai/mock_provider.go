package ai

import (
	"strings"
)

type MockProvider struct{}

func (m *MockProvider) Name() string {
	return "Mock AI"
}

func (m *MockProvider) Chat(
	request ChatRequest,
) (string, error) {

	if len(request.Messages) == 0 {
		return "Mock AI: Empty conversation.", nil
	}

	last := request.Messages[len(request.Messages)-1]

	reply := GetMockResponse(
		strings.TrimSpace(last.Content),
	)

	return reply, nil
}
