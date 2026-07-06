package ai

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIProviderChatIntegration(t *testing.T) {

	server := httptest.NewServer(
		http.HandlerFunc(
			func(
				w http.ResponseWriter,
				r *http.Request,
			) {

				if r.URL.Path != "/chat/completions" {
					t.Fatalf(
						"unexpected path: %s",
						r.URL.Path,
					)
				}

				w.Header().Set(
					"Content-Type",
					"application/json",
				)

				fmt.Fprint(
					w,
					`{
						"choices":[
							{
								"message":{
									"content":"Hello from test server"
								}
							}
						]
					}`,
				)
			},
		),
	)

	defer server.Close()

	provider := &OpenAIProvider{
		baseURL: server.URL,
	}

	reply, err := provider.Chat(
		ChatRequest{
			APIKey: "dummy-key",
			Model:  "dummy-model",
		},
	)

	if err != nil {
		t.Fatalf(
			"unexpected error: %v",
			err,
		)
	}

	if reply != "Hello from test server" {
		t.Fatalf(
			"unexpected reply: %q",
			reply,
		)
	}
}
