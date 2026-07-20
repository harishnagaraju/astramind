package ai

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestOpenAIProviderStreamIntegration(t *testing.T) {

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
					"text/event-stream",
				)

				flusher, ok := w.(http.Flusher)
				if !ok {
					t.Fatal("streaming unsupported")
				}

				fmt.Fprint(
					w,
					`data: {"choices":[{"delta":{"content":"Hello "}}]}`+"\n\n",
				)
				flusher.Flush()

				fmt.Fprint(
					w,
					`data: {"choices":[{"delta":{"content":"World"}}]}`+"\n\n",
				)
				flusher.Flush()

				fmt.Fprint(
					w,
					"data: [DONE]\n\n",
				)
				flusher.Flush()
			},
		),
	)

	defer server.Close()

	provider := &OpenAIProvider{
		baseURL: server.URL,
	}

	stream, err := provider.Stream(
		context.Background(),
		ChatRequest{
			Model:  "dummy",
			APIKey: "dummy",
		},
	)

	if err != nil {
		t.Fatal(err)
	}

	var text string

	for event := range stream.Events() {

		switch event.Type {

		case StreamEventToken:
			text += event.Content

		case StreamEventDone:

			if text != "Hello World" {
				t.Fatalf(
					"unexpected stream text: %q",
					text,
				)
			}

			return

		case StreamEventError:
			t.Fatal(event.Err)
		}
	}

	t.Fatal("stream ended unexpectedly")
}
