# Changelog

All notable changes to AstraMind are documented in this file.

The project follows [Semantic Versioning](https://semver.org/).
---
## [v0.7.0] - 2026-07-10

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
## v0.6.0 - 2026-07-07

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
# v0.5.0 - 2026-07-06

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

# v0.4.1 - 2026-07-01

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