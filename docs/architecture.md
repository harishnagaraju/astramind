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

# RAG Query Routing (`/kb ask`)

Added during the v0.9.1 validation branch, replacing an earlier single-path design that sent every question through a free-form LLM prompt.

```text
                         +----------------------+
                         |   User Question       |
                         +----------+-----------+
                                    |
                                    v
                         +----------------------+
                         | IsEnumerationQuery    |
                         | (question router)     |
                         +----------+-----------+
                                    |
                 +------------------+------------------+
                 |                                     |
        enumeration-style                      single specific fact
      ("what are all/the X")                  ("is there", "what time
                 |                              is", "how much", ...)
                 v                                     v
      +----------------------+            +----------------------+
      |    ExtractItems      |            |   ExtractiveAnswer    |
      |  + BuildListAnswer   |            | (whole-chunk return)  |
      | (deterministic,      |            | (deterministic,       |
      |  zero LLM calls)     |            |  zero LLM calls)      |
      +----------+-----------+            +----------+-----------+
                                                       |
                                          embedder unavailable, or
                                          no embeddable chunk found
                                                       |
                                                       v
                                          +----------------------+
                                          | BuildSemanticPrompt   |
                                          |   + LLM Chat()        |
                                          | (fallback only)       |
                                          +----------------------+
```

**Why this exists, and why it isn't a single LLM prompt anymore:** the original design sent every `/kb ask` question through `BuildSemanticPrompt` + a chat completion, trusting the model to enumerate every matching fact and restate individual facts correctly. Validation testing (v0.9.1, using `gemma2:9b` on real hardware) found this unreliable in a way that could not be fixed by prompt wording, question rewording, or temperature tuning alone: with a chunking bug fixed, retrieval confirmed complete, and the prompt confirmed correct and unmodified, the same question against the same unchanged content still produced different, incomplete answers across repeated runs. The variance was isolated to the LLM's generation step itself - a known limitation of smaller models on exhaustive-enumeration tasks, not a defect in retrieval or prompt construction (both were independently verified correct).

**The fix is architectural, not a better prompt:** for the two question shapes where the correct answer is knowable in advance from what's already been retrieved (list everything matching X; state one specific fact), skip the LLM's generation step entirely and construct the answer directly from retrieved chunk content:

- **`IsEnumerationQuery`** classifies the question. Default is enumeration (a false positive here just means a single-fact question gets a broader answer than strictly necessary, which is safe; a false negative would mean a "list everything" question is answered by extraction of a single item, which is unsafe). Only patterns that clearly signal one specific fact (`is there`, `what time is`, `how much`, `who is`, `where is`, `when is/does/will`, `what is the`) are excluded from the enumeration default.
- **`ExtractItems` + `BuildListAnswer`** (enumeration path): every chunk `SemanticSearch` retrieved is returned as one item, formatted as a list, with source citations. No filtering of items against the question's wording - retrieval already decided relevance at the chunk level, and filtering again afterward was the exact mechanism that caused the original narrow-keyword-matching bug.
- **`ExtractiveAnswer`** (single-fact path): ranks chunks by embedding similarity to the question, returns the single best-matching chunk's content verbatim. A sentence-level windowing approach (return only the specific matching sentence plus a few neighbors) was tried and measured to fail: live testing against real embeddings showed no reliable topic-boundary signal at sentence granularity - an unrelated sentence scored *higher* similarity to the match than genuinely related content did, apparently from incidental word overlap. Whole-chunk return was kept instead: verbose in cases where a chunk bundles multiple unrelated entries, but never wrong, since chunking already guarantees no entry is corrupted or split mid-content.
- **`BuildSemanticPrompt` + LLM `Chat()`** is retained only as a fallback for the single-fact path, used when no embedder is configured. It is not used for enumeration questions at all, and carries the original model-capability caveat when it is used.

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
        chunker.go        - paragraph-aware, CRLF-normalizing document chunking
        manager.go         - import, semantic search, ExtractiveAnswer
        repository.go
        search.go
        prompt.go          - BuildPrompt (keyword), BuildSemanticPrompt (LLM fallback)
        structured_extraction.go  - ExtractItems, BuildListAnswer (deterministic enumeration)
        query_expansion.go - IsEnumerationQuery (question router)
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

- Document import, with paragraph-aware, CRLF-safe chunking (a byte-offset splitter with no word-boundary awareness previously corrupted real-world documents mid-word; the real-world failure case was a CRLF-encoded file, since a `\r\n\r\n` blank line never matches a `\n\n` split point - both are normalized before splitting now)
- Repository management
- Keyword search
- Semantic (embedding-based) search
- Deterministic enumeration and single-fact query answering (see "RAG Query Routing" above)
- LLM-based prompt generation, retained as a fallback only
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
- Ollama (local)
- Mock provider

Future providers:

- Local LLMs beyond Ollama
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
- Ollama (local)

---

# Testing Architecture

The project includes multiple testing layers.

## Unit Tests

- Provider tests
- Chat service tests
- Renderer tests
- Storage tests
- Chunker tests, including a CRLF-specific regression case built from real-world file content, not synthetic LF-only fixtures
- Deterministic extraction tests (`ExtractItems`, `ExtractiveAnswer`), run against real previously-failing document content rather than idealized fixtures

## Integration Tests

Using `httptest.Server`:

- Chat API
- Streaming API
- HTTP error handling
- Invalid JSON responses

## Manual / Content-Fidelity Testing

`tests/integration/manual_testing.sh` extends the interactive walkthrough with:

- A content-fidelity scan - imports a fixture with known facts, asks `/kb ask` questions, and greps the transcript for every expected fact, rather than eyeballing output
- A determinism scan - runs the same question multiple times and flags any fact that appears inconsistently across runs

## Runtime Validation

- OpenRouter compatibility
- Ollama compatibility, including a real-hardware validation pass (see v0.9.1 below)
- Streaming validation
- Session persistence
- Export validation

---

# Current Capabilities

## Knowledge Base

- Document import (paragraph-aware, CRLF-safe chunking)
- Persistent storage
- Chunk repository
- Keyword search
- Semantic search
- Deterministic RAG answering (`/kb ask`) for enumeration and single-fact questions
- LLM-based RAG answering, as a fallback
- Knowledge Base management

## AI

- Multi-provider architecture
- OpenAI-compatible APIs
- OpenRouter integration
- Native Ollama integration, including embeddings
- Configurable API endpoint
- Configurable generation temperature per request
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

## v0.8.0

Knowledge Base

- Document import
- Persistent storage
- Chunk repository
- Keyword search
- Prompt builder
- Knowledge Base management

---

## v0.9.0

Semantic Search

- Embeddings
- `/kb ssearch` — embedding-based semantic search over the KB
- Vector database work has **not** started (still JSON + linear search under the hood)

---

## v0.9.1 (Current — validation branch)

Originally scoped to two narrow validation loops only, with no new features intended:

- **Validate gemma2:9b on real hardware** (#55) — **closed.** Runs correctly and produces accurate output on the target hardware (Intel i5-4210U, 16GB RAM, no GPU); brief UI stutter observed only under simultaneous heavy multitasking, not disqualifying for sequential use.
- **Run the offline demo with real users (lawyer)** (#56) — **still open.**

**Deviation from original scope, noted explicitly rather than silently:** validating #55 required testing `/kb ask` output quality, which surfaced a real chunking bug and a hard limit on free-form LLM enumeration reliability. Fixing the second problem required an architectural change (see "RAG Query Routing" above), not a config tweak - so RAG completion, originally item 1 of the v1.0 backlog below, was substantially delivered on this branch as a direct consequence of investigating the hardware question, not as planned scope expansion. The free-form LLM path (item 1's original design) still exists, but only as a fallback.

---

# Version 1.0 Vision

Scope below is a *backlog*, not a commitment. Updated to reflect that RAG completion has substantially landed on the v0.9.1 branch (see above) - remaining priority order:

1. ~~**RAG completion**~~ — substantially delivered on v0.9.1 as a deterministic dual-path design (see "RAG Query Routing" above), not the free-form LLM design originally envisioned. Remaining under this heading: revisit after #56.
2. **Knowledge Base completion** — `/kb info`, `/kb update`, `/kb rebuild`, `/kb export`; PDF/Word/HTML import. Now the top active priority.
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
- Prefer deterministic, verifiable code paths over LLM generation wherever the correct answer is already knowable from retrieved data - added as an explicit principle following the v0.9.1 RAG reliability investigation