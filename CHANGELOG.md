# Changelog

All notable changes to AstraMind are documented in this file.

The project follows [Semantic Versioning](https://semver.org/).
---
---
## [v0.9.1] - 2026-07-24

### Highlights

This release closes out the v0.9.1 validation branch: a hardware check on `gemma2:9b` (issue #55) that, in the course of validating it, surfaced a real chunking bug, a real prompt-construction gap, and a hard limit on how far prompt engineering alone can make free-form LLM enumeration reliable. Rather than continuing to tune prompt wording and generation temperature, `/kb ask` was split into two paths: a deterministic, zero-LLM-call extraction path for questions with a discrete, verifiable answer, and the existing free-form RAG path, now used only as a fallback. Every fix below was verified against a real, previously-used document (`Sanskrit1.txt`, CRLF line endings included) through the actual CLI and real Ollama embeddings - not synthetic test fixtures alone.

### Added

#### Deterministic Query Answering

- `kb.IsEnumerationQuery` - classifies a question as enumeration-style ("what are all the X", "what are the X timings") versus a single specific fact, used to route `/kb ask`
- `kb.ExtractItems` / `kb.BuildListAnswer` - deterministic enumeration path: retrieves chunks as usual, returns each retrieved chunk's content as one item, formats as a list. No LLM call, so no possibility of a matching entry being silently dropped during generation
- `kb.ExtractiveAnswer` - deterministic single-fact path: ranks chunks by embedding similarity to the question, returns the single best-matching chunk's content verbatim. No LLM paraphrase step, so no possibility of a generative model misstating a date, fee, or other fact while rewording it
- `ai.ChatRequest.Temperature` (optional, `*float64`) - configurable generation temperature, used at 0.4 for the RAG fallback path

#### Testing

- Content-fidelity and determinism checks in `tests/integration/manual_testing.sh`: imports a fixture with known facts, asks the same question multiple times, and greps the transcript for every expected fact rather than eyeballing the output
- `TestChunkHandlesCRLFLineEndings`, `TestChunkRespectsParagraphBoundaries` - regression tests for the chunking fixes below, using real CRLF content
- `TestExtractItems_RealSanskritDocument`, `TestExtractiveAnswer_SelectsCorrectChunk` - run against the real `Sanskrit1.txt` bytes, not a synthetic fixture

### Fixed

- **Chunker corrupted real documents mid-word.** Byte-offset chunking had no word or paragraph awareness and could split a chunk boundary through the middle of a word (e.g. "Youth group" -> "uth group"). Fixed with paragraph-aware splitting on blank-line boundaries, falling back to the old sliding-window split only for a single paragraph too large to fit in one chunk.
- **The paragraph-aware fix above silently did not apply to real-world CRLF files.** The real test document uses `\r\n` line endings; a CRLF blank line (`\r\n\r\n`) never contains the substring `\n\n`, so the paragraph splitter found zero boundaries in it and fell straight back to the byte-offset splitter - reproducing the exact bug it was meant to fix, invisibly, because every test fixture at the time used LF-only content. Fixed by normalizing `\r\n`/`\r` to `\n` before splitting.
- **`/kb ask` silently omitted valid entries on enumeration questions, even with a complete and correct prompt.** Root-caused through a full pipeline audit (disk read -> chunk -> embed -> prompt -> generation): retrieval, chunking, and prompt content were all confirmed correct and complete in isolation. The remaining variance was in the LLM's free-form generation step - a model reliably following the question's literal keywords over weaker instruction wording, especially at low temperature. Not fixable through prompt wording, question rewording, or temperature tuning alone (several approaches were tried and measured; see Known Limitations). Resolved architecturally: see "Deterministic Query Answering" above.
- `BuildSemanticPrompt` had regressed to its pre-fix form (a single generic trailing instruction instead of explicit per-source enumeration) due to an unrelated file being committed under its name in a prior commit. Restored.

### Known Limitations

- **Single-fact and enumeration answers return whole chunks verbatim, which can be verbose when a chunk bundles multiple unrelated entries.** Two windowing approaches were tried and abandoned before landing on whole-chunk return:
  - A fixed-size window (N sentences before/after the matched sentence) sometimes bled into a genuinely unrelated neighboring entry.
  - A dynamic, embedding-similarity-threshold window (expand while neighboring sentences stay semantically similar to the match) was tried next, on the theory that a real topic boundary would show up as a similarity drop. Live testing against real Ollama embeddings disproved this for short-sentence prose: an unrelated sentence ("Not meeting on 16 February", about a different class) scored *higher* similarity to the anchor ("Meeting ID 795 777 3585") than genuinely related sentences did (the actual Zoom URL and password), apparently because both happened to share the literal word "meeting" despite unrelated meaning. No threshold could separate on-topic from off-topic content using this signal - this was a real, measured finding, not a miscalibration. Whole-chunk return was kept as the safer default: chunking already guarantees no entry is corrupted or split mid-content, so the worst failure mode is extra (correct) text, never wrong or missing text.
- `ExtractiveAnswer` re-embeds every sentence in every retrieved chunk on every single-fact question, with no caching. Fine at the scale tested (single-digit chunk counts); will need embedding caching at import time before this scales to a larger knowledge base.
- The free-form LLM RAG path (`BuildSemanticPrompt` + `Chat`) is now only used as a fallback when no embedder is configured. It is not covered by the completeness/precision guarantees of the deterministic paths above.
- `internal/features/kb/query_expansion.go`'s `ExpandQuery` function (question-rewording via appended instructions) is superseded by the deterministic extraction path and is no longer called anywhere in the codebase. `IsEnumerationQuery` from the same file is still in active use as the query router. Cleanup (removing the now-dead `ExpandQuery` code) is planned but not yet done.

### Tested

- go fmt
- go vet
- go build
- go test -v ./...
- tests/integration/run_all.sh
- tests/integration/manual_testing.sh (extended with content-fidelity and determinism scans)
- Full build and test suite re-verified against a fresh pull of the actual pushed branch from GitHub, independent of local working-directory state

### Verified

- Chunking fix confirmed against the real `Sanskrit1.txt` file's actual bytes (CRLF line endings, real diacritics) - zero corruption, all 9 real entries intact
- Enumeration path (`/kb ask what are the Sanskrit class timings`) confirmed live via the CLI against real Ollama embeddings: all 9 entries present, correctly cited across 3 sources
- Single-fact path (`/kb ask what is the meeting id`) confirmed live via the CLI against real Ollama embeddings: correct chunk selected, single source cited, correct answer present, no fabrication
- Hardware validation (issue #55) closed: `gemma2:9b` produces correct output on the target hardware (Intel i5-4210U, 16GB RAM, no GPU); brief UI stutter observed only under simultaneous heavy multitasking, not disqualifying for sequential use

---

## [v0.9.0] - 2026-07-21

### Highlights

This release turns AstraMind from a keyword-search Knowledge Base into a working Retrieval-Augmented Generation assistant, and closes out a substantial architectural cleanup identified in a full review of the v0.8.0 codebase. It also adds a local, browser-based interface for non-technical users, alongside the existing CLI. A known model-capability limitation was found and documented through controlled testing - see Known Limitations below before treating this as demo-ready for high-stakes use.

### Added

#### Semantic Search & RAG

- `ai.EmbeddingProvider` interface (Ollama, OpenAI, and a deterministic mock for tests), matching the existing `StreamingProvider` pattern
- Embedding generation wired into `/kb import`
- `kb.CosineSimilarity` and `Repository.SemanticSearch`, ranking chunks by embedding similarity rather than keyword count
- `/kb ssearch <text>` - semantic search, kept as a separate command alongside keyword `/kb search`
- `kb.BuildSemanticPrompt` and `/kb ask <question>` - completes the RAG loop (import -> chunk -> embed -> retrieve -> **answer**), citing sources with every response
- Local web UI (`--web` flag): embedded HTTP server and single-file browser interface, exposing import/list/ask over a JSON API, reusing the same backend and provider configuration as the CLI (online or offline)

#### Architecture

- `Dependencies` struct extended with `HistoryService`, `SessionService`, `ExportService`, `SearchService` - every feature service now constructed once at startup and shared, instead of built ad-hoc per command
- `storage.FileHistoryStore` - configurable-directory session storage, mirroring the existing `kb.JSONStorage` pattern
- `history.Store` and `session.Store` interfaces, enabling isolated, test-injectable storage
- `--script` mode routes through the full command dispatcher, matching interactive mode

#### Testing

- First unit tests for `history` and `session` packages (previously untested)
- Semantic search, RAG prompt, and context-window regression test suites
- `tests/integration/manual_testing.sh` - full manual walkthrough including a live keyword-vs-semantic comparison and a `--web` API smoke test

### Improved

- `search` decoupled from `storage`, now depends on `history` (`search -> history -> storage`)
- `session` decoupled from duplicated storage logic, now delegates to `history` for Save/Load/Delete/List
- `README.md` and `docs/roadmap.md` updated to reflect delivered semantic search and RAG

### Fixed

- Silent embedding failure during `/kb import`, caused by two missing provider files - now surfaced with visible import feedback instead of failing invisibly
- Test pollution of the real `data/sessions` folder - storage location is now injectable, and tests use isolated temp directories
- Ollama RAG answers truncating mid-generation: no context window was ever specified, so Ollama fell back to its own small default; requests now set `num_ctx: 8192`
- `SemanticSearch` had no result limit - every embedded chunk in the entire knowledge base was returned on every question with no cap, which is harmless with a few test documents but would stuff the whole knowledge base into every prompt at real scale; capped to the top 5 most relevant chunks
- Broken `.gitignore` line from a prior append with no trailing newline
- Compiled binary (`astramind`) no longer tracked in git

### Removed

- `chat/dispatcher.go` and `chat/script.go` - dead code superseded by dispatcher-based `--script` routing
- `kb.Service` - unused wrapper with zero callers anywhere in the codebase

### Known Limitations

- **`gemma3:1b` (the smallest practical local Ollama model) can fabricate facts and omit valid entries on exhaustive multi-item extraction, even with correct source text directly in its context.** Proven via a controlled experiment: an identical prompt, identical document, and identical question given to `gemma3:1b` versus a ~25B-parameter model (`google/gemma-4-26b-a4b-it` via OpenRouter) - the larger model correctly listed all 5 matching entries with zero errors; the smaller model returned only 2 of 5, and in an earlier run fabricated a duplicate entry with an invented date. This is a model-capability ceiling, not a defect in AstraMind's retrieval, prompt construction, or context handling - all of which were independently verified as correct during this investigation.
- Practical consequence: `/kb ask` and the web UI's answer generation should not be treated as reliable for high-stakes use (e.g. legal research) on `gemma3:1b` or comparably small local models. A real minimum local-model size/hardware floor has not yet been established - testing a mid-size local model (e.g. `gemma2:9b`) against realistic hardware is a recommended next step before any customer-facing demo of the RAG feature.
- `/about` and `/config` still report `v0.8.0` - the version constant has not yet been bumped for this release.

### Tested

- go fmt
- go vet
- go build
- go test -v ./...
- tests/integration/run_all.sh
- tests/integration/manual_testing.sh (including the `--web` API smoke test)

### Verified

- Semantic search and RAG proven end-to-end on real Ollama + `nomic-embed-text`, including a real, previously-unseen document (not a synthetic test fixture)
- No test-created sessions or documents polluting the real `data/` folder after the storage isolation and test-cleanup fixes
- Model-capability limitation isolated via a controlled comparison (see Known Limitations), confirming the issue is not in AstraMind's code

---

## [v0.8.0] - 2026-07-13

### Highlights

This release introduces AstraMind's first Knowledge Base implementation, providing document import, persistent storage, automatic chunking, keyword search, knowledge management commands, and the architectural foundation for Retrieval-Augmented Generation (RAG).

### Added

#### Knowledge Base

- Built-in Knowledge Base framework
- Text document import
- Markdown document import
- Automatic document chunking
- Persistent document storage
- Persistent chunk storage
- Repository abstraction
- Keyword search engine
- Prompt builder for future RAG support
- Knowledge Base statistics
- Knowledge Base management API

#### CLI Commands

- `/kb import <file>`
- `/kb list`
- `/kb search <text>`
- `/kb remove <id>`
- `/kb clear`
- `/kb stats`

#### Testing

- Knowledge Base unit tests
- Repository tests
- Search tests
- Prompt builder tests
- Management tests
- CLI command tests
- End-to-end Knowledge Base workflow validation

### Improved

- Chat service architecture
- Dependency injection for future extensibility
- Storage abstraction
- Repository design
- Internal package organization
- Command dispatch architecture
- Overall extensibility for future RAG integration

### Tested

- go fmt
- go vet
- go build
- go test -v ./...

### Verified

- Document import
- Persistent Knowledge Base storage
- Chunk generation
- Keyword search
- Knowledge Base statistics
- Document removal
- Knowledge Base clearing
- Complete `/kb` command workflow

---

## [v0.7.0] 

### Added

- Native Ollama provider
- Local LLM support through Ollama
- Runtime provider selection
- Streaming support for Ollama
- Streaming integration tests
- HTTP error handling tests
- Invalid JSON handling tests
- Connection failure tests

### Improved

- Provider factory architecture
- Provider manager
- Streaming abstraction
- Request builder reuse
- Test coverage

### Supported Providers

- Mock AI
- OpenAI
- OpenRouter
- Ollama

### Tested

- go fmt
- go vet
- go build
- go test -v ./...

### Verified

- Local Ollama installation
- Gemma3:1b model
- Conversation history
- Session persistence
- Provider switching
- Streaming implementation

---
## v0.6.0 

### Added

- Added `/search` command to search the current conversation.
- Added `/searchall` command to search across all saved sessions.
- Added case-insensitive conversation search.
- Added grouped search results by session.
- Added reusable search renderer.
- Added search models for current and multi-session search.

### Improved

- Improved search result presentation.
- Improved session grouping for multi-session search.

### Fixed

- Fixed conversation persistence after successful AI responses.
- Fixed search across sessions by ensuring conversations are saved immediately.
- Improved session consistency when switching between conversations.

### Tests

- Added unit tests for conversation search.
- Added integration tests for multi-session search.

---
# v0.5.0 

## Highlights

This release introduces real-time streaming AI responses, OpenRouter integration, configurable OpenAI-compatible endpoints, improved provider architecture, and a comprehensive integration testing framework. AstraMind now supports streaming and non-streaming providers through a unified interface while significantly improving maintainability and test coverage.

## Added

### Streaming Responses
- Real-time token-by-token AI streaming.
- Streaming renderer for terminal output.
- Streaming provider interface.
- Mock streaming provider.
- Automatic fallback for non-streaming providers.

### OpenAI-Compatible Providers
- OpenRouter integration.
- Configurable API base URL.
- Configurable AI model selection.
- Shared HTTP request builder.
- Improved provider abstraction.

### Integration Testing
- HTTP integration tests using `httptest.Server`.
- Streaming integration tests.
- HTTP error integration tests.
- Invalid JSON response tests.
- End-to-end provider validation without external services.

### Configuration
- `OPENAI_BASE_URL` environment variable.
- Support for OpenAI-compatible APIs.
- Runtime provider configuration improvements.

## Improved

### AI Provider Framework
- Reduced duplicated request construction.
- Cleaner provider implementation.
- Improved request handling.
- Better error propagation.
- More maintainable streaming architecture.

### Testing
- Expanded automated test suite.
- Increased integration test coverage.
- Improved provider validation.
- Enhanced streaming validation.

### Developer Experience
- Cleaner internal architecture.
- Better separation of provider responsibilities.
- Improved code maintainability.
- Simplified future provider integrations.

## Fixed

- Streaming response handling.
- Provider request construction.
- HTTP error handling consistency.
- OpenAI-compatible endpoint support.
- Runtime streaming stability.

## Testing

Validated successfully with:

- Unit tests
- Integration tests
- Streaming integration tests
- HTTP error integration tests
- Mock provider tests
- Mock streaming provider tests
- Provider manager tests
- Renderer tests
- Storage tests
- Runtime validation against OpenRouter

---

# v0.4.1 

## Highlights

This release introduces a major architectural milestone for AstraMind, including a modular AI provider framework, conversation export capabilities, automated testing, and a production-grade GitHub Actions CI pipeline.

## Added

### Export System
- TXT conversation export
- Markdown conversation export
- Automatic export directory creation

### AI Provider Framework
- AI provider abstraction
- Provider factory
- Provider manager
- Mock AI provider
- OpenAI provider
- Automatic provider failover

### Session Management
- Session creation
- Session loading
- Session deletion
- Session listing
- Active session tracking

### Testing
- Storage integration tests
- Mock provider tests
- Regression test suite
- Coverage generation

### CI/CD
- GitHub Actions workflow
- Automatic formatting (`go fmt`)
- Static analysis (`go vet`)
- Automated builds
- Automated unit tests
- Coverage reporting
- Coverage artifact publishing
- Workflow concurrency
- Workflow timeout
- Least-privilege workflow permissions

## Improved

- Modular project architecture
- Export subsystem
- Storage layer
- Test infrastructure
- Repository documentation
- Project roadmap
- README documentation

## Fixed

- Export failures when the export directory did not exist
- GitHub Actions workflow stability
- Cross-platform CI compatibility

---

# v0.4.0

## Added

- Multi-session support
- Session management
- Improved API error handling

---

# v0.3.0

## Added

- Persistent conversation history
- Statistics command
- Configuration command
- About command

---

# v0.2.0

## Added

- Conversation memory
- History support
- Chat commands

---

# v0.1.0

## Initial Release

- AstraMind CLI chatbot
- OpenAI integration
- Environment configuration