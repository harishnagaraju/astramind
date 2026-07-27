package ai

import (
	"net/http"
	"testing"
	"time"

	"github.com/harishnagaraju/astramind/internal/infrastructure/models"
)

// These benchmarks measure real, live round-trip latency to a
// running Ollama instance - genuinely different from a standard Go
// benchmark, which only measures CPU/allocations for pure in-process
// code. Chat/Embed calls are dominated by network + model inference
// time, not Go code execution, so -benchmem's allocation counts
// would be nearly meaningless here; wall-clock time (which "go test
// -bench" reports by default) is the number that actually matters.
//
// Guarded to skip cleanly if Ollama isn't reachable, mirroring the
// same opt-in pattern already used for --web in
// scripts/check_rag_behavior.sh - these must never break a plain
// `go test ./...` run for someone without Ollama running locally.
//
// Run explicitly with:
//   go test -bench=BenchmarkOllama -benchtime=5x ./internal/infrastructure/ai/
//
// -benchtime=Nx (a fixed number of iterations) is deliberately used
// instead of Go's default time-based benchmarking (which runs for a
// fixed duration, executing as many iterations as fit) - letting a
// slow LLM benchmark run for its default ~1 second would call the
// model exactly once or twice and produce a noisy, unrepresentative
// number. A fixed, small iteration count gives a deliberate,
// controlled sample size instead.

const ollamaBenchmarkBaseURL = "http://localhost:11434"

// skipIfOllamaUnreachable does a fast, cheap connectivity check
// before running a real benchmark iteration - failing fast with a
// clear message beats a slow, confusing timeout per iteration.
func skipIfOllamaUnreachable(b *testing.B) {
	b.Helper()

	client := http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(ollamaBenchmarkBaseURL)
	if err != nil {
		b.Skipf("Ollama not reachable at %s - skipping (start Ollama to run this benchmark): %v", ollamaBenchmarkBaseURL, err)
	}
	resp.Body.Close() //nolint:errcheck // read-only handle in a benchmark; a close failure here doesn't lose data
}

// BenchmarkOllamaChatLatency measures real end-to-end latency for a
// single chat completion against whatever model is currently running
// locally. Reports ns/op (nanoseconds per call) - divide by 1e9 for
// seconds, which is the number that actually matters for judging
// "does this feel fast enough" during real usage.
func BenchmarkOllamaChatLatency(b *testing.B) {
	skipIfOllamaUnreachable(b)

	provider := &OllamaProvider{baseURL: ollamaBenchmarkBaseURL, model: "gemma2:9b"}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := provider.Chat(ChatRequest{
			Messages: []models.Message{{Role: "user", Content: "Say OK and nothing else."}},
		})
		if err != nil {
			b.Fatalf("chat request failed: %v", err)
		}
	}
}

// BenchmarkOllamaEmbedLatency measures real end-to-end latency for a
// single embedding call - typically much faster than chat completion
// (no generation, just one forward pass), useful to isolate whether
// retrieval slowness in /kb ask comes from embedding calls or the
// deterministic extraction logic around them.
func BenchmarkOllamaEmbedLatency(b *testing.B) {
	skipIfOllamaUnreachable(b)

	provider := &OllamaProvider{baseURL: ollamaBenchmarkBaseURL}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := provider.Embed(EmbeddingRequest{Text: "what is the zoom meeting id"})
		if err != nil {
			b.Fatalf("embed request failed: %v", err)
		}
	}
}
