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
RAG (Retrieval-Augmented Generation) — retrieval → prompt builder → LLM → answer; not yet built, top priority for v1.0
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

## v0.9.0 (In Progress) - Semantic Search

- Embeddings
- `/kb ssearch` — embedding-based semantic search
- Vector database work not yet started (still JSON + linear search)

## v0.9.1 (Current) - Validation Branch

Not a feature release. Scoped to two open loops before v1.0 gets named and scoped for real:

- Validate gemma2:9b on real hardware (answer completeness + usability during generation)
- Niece/lawyer offline demo — capture actual behavior, not just stated feedback

Notes/Tasks/Calendar (previously slated for v0.9.0) has been deprioritized — see v1.0 backlog below. It now ships as a plugin/feature after RAG, not as a v0.9 milestone.

---

## v1.0.0 - Personal AI Assistant Platform

Priority order (backlog, not committed — see "v1.0 roadmap" GitHub issue):

1. RAG completion (`/kb ask`) — blocked on v0.9.1 validation
2. Knowledge Base completion (`/kb info`, `/kb update`, `/kb rebuild`, `/kb export`; PDF/Word/HTML import)
3. Vector store migration (SQLite + vector index)
4. Provider abstraction expansion (Gemini, Claude, OpenRouter, LM Studio)
5. Search improvements (fuzzy matching, ranking, highlighting, filters, regex)
6. Plugin architecture (weather, filesystem, calculator, web)
7. Notes / Tasks / Calendar (as a plugin, not core)

Enterprise items (Web Interface, Authentication, PostgreSQL memory, user profiles, multi-agent workflows, enterprise integrations) are pushed out of v1.0 scope — see Phase 2/3/5 below, now explicitly post-v1.0.

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

- RAG
- PDF ingestion
- Document indexing
- Semantic search (✅ basic version shipped in v0.9.0 - embedding-based /kb ssearch)
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
- ✅ Provider Factory
- ✅ Provider Manager
- ✅ Automatic Provider Failover

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

- ✅ /kb import
- ✅ /kb list
- ✅ /kb search
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
- ✅ Coverage reporting
- ✅ Semantic Versioning
- ✅ GitHub Issues
- ✅ GitHub Milestones
- ✅ Modular Go Architecture

---

# Current Stable Features

## AI Assistant

- OpenAI-compatible API integration
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

---

# Current Release Status

| Version  |     Status            							|
|----------|----------------------------------------------------------------------------|
| v0.4.1 | old Stable Release      							|
| v0.5.0 | old Stable Release      							|
| v0.6.0 | old Stable Release      							|
| v0.7.0 | old Stable Release      							|
| v0.8.0 | ✅ Latest Stable Release							|
| v0.9.0 | 🚧 In Progress (semantic search shipped, RAG not yet) 			|
| v0.9.1 | 🚧 Current — validation branch (hardware test + user demo, no new features) 	|
| v1.0.0 | 📋 Planned (scope pending v0.9.1 results) 					|
| v1.1.0 | 🎯 Long-Term Vision     							|

---

**AstraMind continues to evolve through incremental, test-driven development with a strong focus on architecture, maintainability, and extensibility.**