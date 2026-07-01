package ai

import "strings"

func GetMockResponse(prompt string) string {

	err := InitializeMockCache()

	if err != nil {
		return "Mock AI data could not be loaded."
	}

	prompt = strings.ToLower(
		strings.TrimSpace(prompt),
	)

	for _, pair := range mockCache {

		if strings.EqualFold(
			pair.Prompt,
			prompt,
		) {

			return pair.Response

		}
	}

	return "Mock AI: I don't have a predefined response for that prompt."
}
