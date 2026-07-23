# AstraMind Roadmap

---

# AstraMind Vision

**Intelligent Conversations. Infinite Possibilities.**

AstraMind is evolving from a lightweight AI-powered CLI assistant into a modular AI platform capable of supporting conversational intelligence, knowledge retrieval, automation, enterprise integrations, and autonomous AI agents.

The project emphasizes:

- Clean Architecture
- Extensible AI Provider Framework
- High Performance
- Cross-platform Support
- Testability
- Maintainability
- Enterprise Readiness

---

# AstraMind v1.0.0 Vision

The long-term vision for AstraMind includes:

- Real-time AI conversations
- CLI and Web user interfaces
- Multi-provider AI architecture
- Local and Cloud LLM support
- Persistent AI memory
- Knowledge Base
- RAG (Retrieval-Augmented Generation)
- Semantic Search
- Notes Management
- Task Tracking
- Personal AI Assistant
- Multi-session conversations
- Tool Calling
- Workflow Automation
- Multi-Agent Architecture
- Enterprise Integrations

AstraMind v1.0 provides:

Single-user conversational AI
Real-time prompt-response interaction
CLI-based user experience
API-driven intelligence
Future Roadmap
Phase 1
CLI Chat Assistant
Conversation history
Configuration management
Phase 2
Web-based interface using React
User authentication
Session management
Phase 3
PostgreSQL-backed memory
Persistent chat history
User profiles
Phase 4
Semantic search (✅ shipped in v0.9.0 — embedding-based `/kb ssearch`)
RAG (Retrieval-Augmented Generation) (✅ substantially shipped in v0.9.1 — see Release Roadmap below; deterministic dual-path design, not the free-form-LLM design originally envisioned here)
PDF and document ingestion
Vector database integration
Phase 5
Multi-agent architecture
Autonomous task execution
Tool calling and workflow automation
Enterprise integrations

---

# Release Roadmap

## v0.1.0
### Initial Foundation
- CLI chatbot
- OpenAI integration
- Environment configuration

---

## v0.2.0
### Conversation Management
- Conversation memory
- Conversation history
- Basic chat commands

---

## v0.3.0
### Persistence & Usability
- Persistent chat history
- Session statistics
- Configuration command
- About command
- Improved architecture

---

## v0.4.0
### Multi-Session Support
- Active session tracking
- Session creation
- Session loading
- Session deletion
- Session listing
- Improved error handling

---

## v0.4.1

### Export System
- TXT Export
- Markdown Export
- Export current session

### AI Provider Framework
- AI Provider abstraction
- Provider Factory
- Provider Manager
- Mock AI Provider
- OpenAI Provider
- Automatic Provider Failover

### Testing
- Storage integration tests
- Mock provider tests
- Regression test suite
- Coverage generation

### CI/CD
- GitHub Actions
- go fmt
- go vet
- go build
- go test
- Coverage reports
- Coverage artifacts
- Workflow concurrency
- Workflow timeout
- Least-privilege permissions

### Architecture Improvements
- Modular AI architecture
- Improved storage architecture
- Export subsystem
- Session persistence improvements

---

## v0.5.0 (Latest Stable Release)- Streaming Responses
- OpenAI streaming
- Streaming renderer
- Typing indicator

---

## v0.6.0 - Search System
- /search command
- Search conversations
- Search all sessions
- Search results display

---

## v0.7.0 - Local Models

- Provider selection
- Ollama integration
- Local LLM support
- CPU/GPU model execution

---

## v0.8.0 - Knowledge Base

- /kb import
- /kb list
- /kb search
- /kb remove
- /kb clear
- /kb stats
- Text and Markdown document import
- Automatic document chunking
- Persistent document storage
- Persistent chunk storage
- Keyword search
- Prompt builder for future RAG
- Knowledge Base management
- Repository abstraction

---

## v0.9.0 - Semantic Search

- Embeddings
- `/kb ssearch` — embedding-based semantic search
- Vector database work not yet started (still JSON + linear search)

## v0.9.1 (Current) - Validation Branch, Deterministic RAG

Originally scoped as a pure validation branch (no new features) to close two open loops before v1.0 gets named and scoped for real:

- Validate gemma2:9b on real hardware (answer completeness + usability during generation) — **closed.** Model produces correct output on the target hardware (i5-4210U, 16GB RAM, no GPU); brief stutter only under simultaneous heavy multitasking.
- Niece/lawyer offline demo — capture actual behavior, not just stated feedback — **still open.**

**Scope note:** validating the hardware question required testing `/kb ask` answer quality, which surfaced a real document-chunking bug and a hard reliability limit in free-form LLM enumeration - neither fixable by configuration alone. Fixing the second required an architectural change: `/kb ask` now routes enumeration and single-fact questions to deterministic, zero-LLM-call extraction, with the original free-form LLM path kept only as a fallback (no embedder configured). This is a real feature addition on a branch that was scoped to have none - flagged here explicitly rather than left implicit, since it changes what "RAG completion" (previously the #1 v1.0 backlog item) means going forward.

Notes/Tasks/Calendar (previously slated for v0.9.0) remains deprioritized — ships as a plugin/feature after Knowledge Base completion, not as a near-term milestone.

---

## v1.0.0 - Personal AI Assistant Platform

Priority order (backlog, not committed — see "v1.0 roadmap" GitHub issue):

1. ~~RAG completion (`/kb ask`)~~ — substantially delivered in v0.9.1 as a deterministic dual-path design. Free-form LLM synthesis remains available only as a fallback path, not the primary mechanism.
2. **Knowledge Base completion** (`/kb info`, `/kb update`, `/kb rebuild`, `/kb export`; PDF/Word/HTML import) — now the top active priority
3. Vector store migration (SQLite + vector index)
4. Provider abstraction expansion (Gemini, Claude, OpenRouter, LM Studio)
5. Search improvements (fuzzy matching, ranking, highlighting, filters, regex)
6. Plugin architecture (weather, filesystem, calculator, web)
7. Notes / Tasks / Calendar (as a plugin, not core)

Enterprise items (Web Interface, Authentication, PostgreSQL memory, user profiles, multi-agent workflows, enterprise integrations) are pushed out of v1.0 scope — see Phase 2/3/5 above, now explicitly post-v1.0.

---

# Long-Term Technical Roadmap

## Phase 1
### Core AI Assistant

- CLI chatbot
- Conversation management
- Configuration management
- Session persistence

---

## Phase 2
### Modern User Experience

- React frontend
- Authentication
- User profiles
- Session synchronization

---

## Phase 3
### Persistent Intelligence

- PostgreSQL storage
- Long-term memory
- User preferences
- Conversation analytics

---

## Phase 4
### Knowledge Platform

- RAG (✅ substantially shipped in v0.9.1 - deterministic dual-path design)
- PDF ingestion
- Document indexing
- Semantic search (✅ shipped in v0.9.0 - embedding-based /kb ssearch)
- Vector database integration

---

## Phase 5
### Autonomous AI

- Multi-agent architecture
- Tool calling
- Workflow automation
- Enterprise integrations

---

# Current Project Status

## AI Provider Framework

- ✅ OpenAI Provider
- ✅ Mock Provider
- ✅ Ollama Provider (including embeddings)
- ✅ Provider Factory
- ✅ Provider Manager
- ✅ Automatic Provider Failover
- ✅ Configurable per-request generation temperature

---

## Conversation Management

- ✅ Conversation memory
- ✅ Persistent history
- ✅ Session-aware storage
- ✅ Multi-session support

---

## Export System

- ✅ TXT export
- ✅ Markdown export
- ✅ Export current session

---

## Session Commands

- ✅ /sessions
- ✅ /new <session>
- ✅ /load <session>
- ✅ /delete <session>

---

## Utility Commands

- ✅ /help
- ✅ /history
- ✅ /clear
- ✅ /stats
- ✅ /config
- ✅ /about
- ✅ /provider
- ✅ /export
- ✅ /export md

---
## Knowledge Base Commands

- ✅ /kb import (paragraph-aware, CRLF-safe chunking)
- ✅ /kb list
- ✅ /kb search
- ✅ /kb ssearch
- ✅ /kb ask (deterministic enumeration + single-fact extraction, LLM fallback)
- ✅ /kb remove
- ✅ /kb clear
- ✅ /kb stats
- ✅ Document chunking
- ✅ Repository
- ✅ Prompt builder

---
## Developer Experience

- ✅ GitHub Actions CI
- ✅ Automated testing
- ✅ Regression suite
- ✅ Content-fidelity and determinism test scans
- ✅ Coverage reporting
- ✅ Semantic Versioning
- ✅ GitHub Issues
- ✅ GitHub Milestones
- ✅ Modular Go Architecture

---

# Current Stable Features

## AI Assistant

- OpenAI-compatible API integration
- Native Ollama integration (chat, streaming, embeddings)
- Mock AI for offline development
- Interactive command-line chat
- Cross-platform support
- Environment-based configuration

---

## Storage

- Persistent conversations
- Session storage
- Conversation export
- JSON-based storage

---

## Quality

- Unit testing
- Integration testing
- Regression testing
- Content-fidelity and determinism testing (RAG-specific)
- Automated CI
- Code formatting
- Static analysis
- Coverage reports

---

# Development Principles

AstraMind follows the following engineering principles:

- Clean Architecture
- SOLID Principles
- Modular Design
- Incremental Development
- Test-Driven Improvements
- Continuous Integration
- Continuous Refactoring
- Semantic Versioning
- GitHub Issue-Driven Development
- Small, Reviewable Commits
- Prefer deterministic, verifiable code paths over LLM generation wherever the correct answer is already knowable from retrieved data

---

# Current Release Status

| Version  |     Status            							|
|----------|----------------------------------------------------------------------------|
| v0.4.1 | old Stable Release      							|
| v0.5.0 | old Stable Release      							|
| v0.6.0 | old Stable Release      							|
| v0.7.0 | old Stable Release      							|
| v0.8.0 | old Stable Release      							|
| v0.9.0 | ✅ Latest Merged Release (semantic search; RAG completed on v0.9.1 branch, not yet merged) |
| v0.9.1 | 🚧 Current — validation branch; #55 closed, RAG dual-path shipped, #56 (user demo) still open |
| v1.0.0 | 📋 Planned (scope pending #56) 					|
| v1.1.0 | 🎯 Long-Term Vision     							|

---

**AstraMind continues to evolve through incremental, test-driven development with a strong focus on architecture, maintainability, and extensibility.**