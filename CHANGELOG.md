# Changelog

All notable changes to AstraMind are documented in this file.

The project follows [Semantic Versioning](https://semver.org/).

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