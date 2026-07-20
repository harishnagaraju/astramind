# AstraMind Architecture

## Overview

AstraMind is a modular AI-powered command-line assistant written in Go. It provides a clean architecture for integrating multiple Large Language Model (LLM) providers while supporting streaming responses, persistent conversation history, multi-session management, export capabilities, and comprehensive automated testing.

The application is designed around extensible provider abstractions, making it easy to integrate additional AI providers such as OpenAI, OpenRouter, Ollama, and future local LLMs.

---

# High-Level Architecture

```text
                         +----------------------+
                         |       User           |
                         +----------+-----------+
                                    |
                                    v
                         +----------------------+
                         |   CLI Application    |
                         |      main.go         |
                         +----------+-----------+
                                    |
                                    v
                         +----------------------+
                         |   Command Handler    |
                         | /help /history etc.  |
                         +----------+-----------+
                                    |
                                    v
                         +----------------------+
                         |    Chat Service      |
                         +----------+-----------+
                                    |
                                    v
                         +----------------------+
                         |  Provider Manager    |
                         +----------+-----------+
                                    |
                 +------------------+------------------+
                 |                                     |
                 v                                     v
      +----------------------+            +----------------------+
      |   OpenAI Provider    |            |    Mock Provider     |
      +----------+-----------+            +----------+-----------+
                 |                                     |
                 +------------------+------------------+
                                    |
                                    v
                         OpenAI-Compatible APIs
                  (OpenRouter, OpenAI, Future Providers)
```

---
CLI
 │
 ▼
Command Dispatcher
 │
 ├───────────────┐
 │               │
 ▼               ▼
Chat Service   Knowledge Base
 │               │
 ▼               ▼
Provider     KB Manager
Manager         │
 │              ▼
 ▼         Repository
AI              │
                ▼
          JSON Storage

---

# Streaming Architecture

```text
User
 │
 ▼
Chat Service
 │
 ▼
Provider.Stream()
 │
 ▼
HTTP Streaming (SSE)
 │
 ▼
readStream()
 │
 ▼
Stream Events
 │
 ▼
Renderer
 │
 ▼
Terminal Output
```

---

# Current Project Structure

```text
astramind/
│
├── cmd/
│   └── astramind/
│       └── main.go
│
├── internal/
│   ├── ai/
│   │   ├── provider.go
│   │   ├── provider_manager.go
│   │   ├── factory.go
│   │   ├── openai_provider.go
│   │   ├── mock_provider.go
│   │   ├── stream.go
│   │   ├── renderer.go
│   │   ├── errors.go
│   │   └── ...
│   │
│   ├── chat/
│   │   └── service.go
│   │
│   ├── renderer/
│   │
│   ├── storage/
│   │
│   ├── config/
│   │
│   ├── models/
│   │
│   └── testutil/
│
├── data/
├── exports/
├── docs/
├── .github/
├── README.md
├── CHANGELOG.md
├── go.mod
└── go.sum
```
internal/
    kb/
        chunk.go
        manager.go
        repository.go
        search.go
        prompt.go
        storage.go
        json_storage.go
        stats.go

CLI
 │
 ├── Commands
 │
 ├── Search Renderer
 │
 ├── Chat Service
 │
 └── Storage

Search System

/search

↓

SearchMessages()

↓

Renderer

/searchall

↓

SearchAllSessions()

↓

SearchMessages()

↓

Renderer

---

# Core Components

## Knowledge Base

Responsible for:

- Document import
- Chunk generation
- Repository management
- Keyword search
- Prompt generation
- Knowledge Base statistics

## CLI Layer

Responsible for:

- Interactive command-line interface
- Command processing
- Session selection
- User interaction
- Configuration display

---

## Chat Service

Responsible for:

- Building chat requests
- Managing conversations
- Streaming orchestration
- Provider communication

---

## Provider Manager

Responsible for:

- Active provider selection
- Provider failover
- Provider abstraction

---

## AI Providers

Current providers:

- OpenAI-compatible provider
- Mock provider

Future providers:

- Ollama
- Local LLMs
- Anthropic
- Google Gemini
- Azure OpenAI

---

## Storage Layer

Responsible for:

- Conversation persistence
- Knowledge Base storage
- Chunk storage
- Repository abstraction
- Session management
- Export
- History loading

---

## Renderer

Responsible for:

- Streaming token rendering
- Console output
- Error display

---

# Configuration

Environment variables:

```text
OPENAI_API_KEY
OPENAI_MODEL
OPENAI_BASE_URL
```

Supported endpoints include:

- OpenAI
- OpenRouter
- Any OpenAI-compatible API

---

# Testing Architecture

The project includes multiple testing layers.

## Unit Tests

- Provider tests
- Chat service tests
- Renderer tests
- Storage tests

## Integration Tests

Using `httptest.Server`:

- Chat API
- Streaming API
- HTTP error handling
- Invalid JSON responses

## Runtime Validation

- OpenRouter compatibility
- Streaming validation
- Session persistence
- Export validation

---

# Current Capabilities

## Knowledge Base

- Document import
- Automatic chunking
- Persistent storage
- Chunk repository
- Keyword search
- Prompt builder
- Knowledge Base management

## AI

- Multi-provider architecture
- OpenAI-compatible APIs
- OpenRouter integration
- Configurable API endpoint
- Automatic provider failover

## Streaming

- Token streaming
- Streaming renderer
- Mock streaming provider
- End-to-end streaming support

## Sessions

- Persistent history
- Multi-session support
- Session export
- Statistics
- Configuration display

## Developer Experience

- Modular architecture
- Automated testing
- GitHub Actions CI
- Integration tests
- Semantic versioning
- Release management

---

# Architecture Roadmap

## v0.6.0

Search System

- Conversation search
- Session search
- Search highlighting

---

## v0.7.0

Local Models

- Ollama integration
- Provider selection
- Local LLM execution

---

## v0.8.0 (Latest Stable Release)

Knowledge Base

- Document import
- Persistent storage
- Chunk repository
- Keyword search
- Prompt builder
- Knowledge Base management

---

## v0.9.0 (In Progress)

Semantic Search

- Embeddings
- `/kb ssearch` — embedding-based semantic search over the KB
- Vector database work has **not** started (still JSON + linear search under the hood)
- Note: semantic search returning ranked chunks is not the same as RAG. Retrieval still stops before a prompt is built or sent to an LLM.

---

## v0.9.1 (Current — validation branch, not a feature branch)

Scoped narrowly to closing two open loops before any further architecture is committed:

- Validate gemma2:9b on real hardware (answer completeness + system usability during generation)
- Run the offline demo with real users (niece/lawyer) and capture actual behavior, not just stated reactions

No new features land in this branch. It exists to inform what v1.0 should actually be, rather than deciding that in advance.

---

# Version 1.0 Vision

Scope below is a *backlog*, not a commitment — it is intentionally re-ordered from earlier drafts of this document to prioritize finishing what's already started (RAG, KB, retrieval) over adding new surface area (notes/tasks/plugins). Order and inclusion are subject to change based on v0.9.1 validation results. See the "v1.0 roadmap" GitHub issue for the working priority list.

1. **RAG completion** — wire semantic search into a prompt builder → LLM → answer flow (`/kb ask`), blocked on v0.9.1 hardware validation
2. **Knowledge Base completion** — `/kb info`, `/kb update`, `/kb rebuild`, `/kb export`; PDF/Word/HTML import
3. **Vector store migration** — JSON/linear search → SQLite + vector index (FAISS not required yet)
4. **Provider abstraction expansion** — Gemini, Claude, OpenRouter, LM Studio
5. **Search improvements** — fuzzy matching, ranking, highlighting, filters, regex
6. **Plugin architecture** — weather, filesystem, calculator, web — sequenced after RAG, not before
7. **Notes / Tasks / Calendar** — lowest priority; ships as a plugin/feature, not core architecture

Identity check for v1.0: someone downloading it should think "this is my private AI assistant," not "this is a chat program with notes." The differentiator is local LLMs, multiple providers, streaming, Knowledge Base, semantic search, RAG, search, sessions, and script mode — everything else is post-v1.0.

---

# Design Principles

- Clean Architecture
- Modular provider abstraction
- Streaming-first design
- Minimal external dependencies
- Extensible provider framework
- Comprehensive automated testing
- Git-first workflow
- Incremental feature delivery
- Release-driven development
- Production-ready engineering practices