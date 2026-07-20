# Changelog

All notable changes to AstraMind are documented in this file.

The project follows [Semantic Versioning](https://semver.org/).
---
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