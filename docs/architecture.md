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
│   ├── astramind/
│   │   └── main.go
│   └── dev/
│       └── main.go          - task runner: build/test/coverage/junit/regression
│
├── internal/
│   ├── engine/
│   │   └── webui/            - embedded local web interface
│   ├── features/
│   │   ├── chat/
│   │   ├── kb/
│   │   │   ├── chunker.go              - paragraph-aware, CRLF-normalizing chunking
│   │   │   ├── manager.go              - import, semantic search, ExtractiveAnswer
│   │   │   ├── precise_match.go        - findPreciseLine, HasAnyContentWordOverlap,
│   │   │   │                             FilterRelevantChunks
│   │   │   ├── structured_extraction.go - ExtractItems, BuildListAnswer
│   │   │   ├── query_expansion.go      - IsEnumerationQuery (question router)
│   │   │   ├── json_storage.go
│   │   │   └── ...
│   │   ├── history/
│   │   ├── session/
│   │   ├── search/
│   │   └── export/
│   ├── infrastructure/
│   │   ├── ai/                - provider.go, provider_manager.go, factory.go,
│   │   │                        ollama_*.go, openai_*.go, mock_provider.go
│   │   ├── storage/
│   │   ├── models/
│   │   ├── renderer/
│   │   └── config/
│   └── utilityforunittest/    - shared Go test helpers (not related to /tests)
│
├── scripts/                   - every runnable action (build/test/coverage/
│                                 regression/check_knowledge_base/
│                                 check_rag_behavior/manual_walkthrough), .sh + .bat
├── tests/
│   ├── fixtures/               - input data the scripts consume
│   ├── output/                 - generated coverage reports
│   └── docs/                   - test plan reference documents
├── reports/                    - generated junit.xml / regression.xml (gitignored)
│
├── data/
├── exports/
├── docs/
├── .github/workflows/          - CI: Linux + Windows, JUnit reporting
├── README.md
├── CHANGELOG.md
├── go.mod
└── go.sum
```

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

The project includes multiple testing layers. Every build/test/coverage command is implemented exactly once, in `cmd/dev` (a Go program) - `scripts/*.sh` and `scripts/*.bat` are thin wrappers into it, not separate implementations that can silently diverge from each other.

## Unit Tests

- Provider tests
- Chat service tests
- Renderer tests
- Storage tests
- Chunker tests, including a CRLF-specific regression case built from real-world file content, not synthetic LF-only fixtures
- Deterministic extraction tests (`ExtractItems`, `ExtractiveAnswer`), run against real previously-failing document content rather than idealized fixtures
- `FilterRelevantChunks` tests, including the exact real failure case (an unrelated document bundled into an enumeration answer) and a regression guard confirming an earlier, rejected item-level filtering approach is never reintroduced
- Embedding pipeline tests (`ollama_embedding_test.go`, `openai_embedding_test.go`) using `httptest.Server`, not a mocked client field - covering valid responses, empty results, HTTP errors, and malformed JSON for both providers
- Provider fallback behavior tests, documenting that `Stream()`/`Embed()` have no fallback of their own (unlike `Chat()`, which permanently switches to the fallback provider on failure)

## Integration Tests

Using `httptest.Server`:

- Chat API
- Streaming API
- HTTP error handling
- Invalid JSON responses

## Regression Pipeline & Reporting

`scripts/regression.sh`/`.bat` run build, test, coverage, and the KB/RAG behavioral checks in sequence. Each step's real exit code drives its reported status (PASS/FAIL/SKIPPED) - a step is only ever reported as having passed if it actually did, and steps after a failure are marked SKIPPED rather than silently assumed passing.

`go run ./cmd/dev -run=junit` parses `go test -json` output directly into real JUnit XML (`reports/junit.xml`) - no external tool dependency. `go run ./cmd/dev -run=regression-report` writes `reports/regression.xml`, the 4 pipeline steps themselves as JUnit test cases. Both are generated automatically as part of the regression run, and consumed by CI (`.github/workflows/go.yml`) to render results directly on each commit/PR, on both Linux and Windows.

## Manual / Content-Fidelity Testing

`scripts/manual_walkthrough.sh` extends the interactive walkthrough with:

- A content-fidelity scan - imports a fixture with known facts, asks `/kb ask` questions, and greps the transcript for every expected fact, rather than eyeballing output
- A determinism scan - runs the same question multiple times and flags any fact that appears inconsistently across runs
- A web API smoke test (Part 2) - always runs as part of this script; there is no separate flag to enable or disable it

## Runtime Validation

- OpenRouter compatibility
- Ollama compatibility, including a real-hardware validation pass (see v0.9.1 below)
- Streaming validation
- Session persistence
- Export validation
- Cross-platform validation: real runs confirmed on Linux, MINGW64 (git-bash on Windows), and native `cmd.exe`, plus CI on Linux and Windows on every push/PR

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

## v0.9.2 (Current)

Started as four related, targeted changes; grew substantially after real cross-platform manual testing (Linux, MINGW64, native Windows `cmd.exe`) surfaced further correctness and tooling issues.

- `.docx` import, single-fact precision refinement, and a real reliability gap closed (the web UI's `/api/ask` was silently still on the older free-form-LLM path after v0.9.1's redesign).
- Three further `/kb ask` correctness fixes found via real testing, all documented in detail in CHANGELOG.md: plural-question misrouting, a missing relevance threshold for out-of-scope questions, and a short-word substring collision in that threshold's implementation.
- A fourth relevance-filtering gap (enumeration answers bundling unrelated documents in small knowledge bases) found and fixed after release, deliberately at whole-chunk granularity to avoid reintroducing an earlier, rejected item-level filtering approach.
- Project reorganized (`scripts/`, `tests/fixtures|output|docs/`) and a new Go-native task runner (`cmd/dev`) introduced, eliminating a real, recurring class of bug: a `.sh`/`.bat` pair silently diverging from each other. Real instances found and fixed during this release's own testing, on real Linux, MINGW64, and native Windows machines - not just by reading the code:
    - `check_knowledge_base.bat`'s stale pre-reorganization path (root cause of a real `"Could not find astramind.exe"` failure)
    - A `cmd.exe` parser crash in `check_rag_behavior.bat` - literal parentheses inside a parenthesized `if` block, misread as block structure (`"Smoke was unexpected at this time."`), root-caused via a targeted debug line
    - Log files silently missing their own trailer content, in three separate scripts (`check_rag_behavior.sh`/`.bat`, `manual_walkthrough.sh`, `regression.sh`/`.bat`) - fixed via `exec > >(tee ...)` in bash and paired screen+file `echo` in batch, verified by diffing saved logs against real terminal output byte-for-byte
    - `manual_walkthrough.sh` crashing on `--web` - no argument-flag recognition at all, silently misread as a binary path override
    - A stale `.gitignore` silently tracking generated coverage reports into git history
- Test coverage closed: `ollama_embedding_test.go`/`openai_embedding_test.go` (`Embed()` 0% → 90.9% each, `httptest.Server`-based, no mocked client field), `provider_manager_fallback_test.go` (documents `Stream()`/`Embed()` have no fallback of their own, unlike `Chat()`), `json_storage_delete_test.go` (`DeleteDocument`/`DeleteChunks` 50% → 100%, including a documented asymmetry in how each handles a missing file).
- Machine-readable test reporting (`reports/junit.xml`, `reports/regression.xml`, generated by `cmd/dev` directly from `go test -json` - no external tool dependency) and a CI workflow rebuilt to run on both Linux and Windows - previously Linux-only, meaning it would have caught none of the Windows-specific bugs listed above.
- **Real-user (lawyer) offline demo feedback** (#56) — still open, carried forward from v0.9.1.

---

## v0.9.1

Originally scoped to two narrow validation loops only, with no new features intended:

- **Validate gemma2:9b on real hardware** (#55) — **closed.** Runs correctly and produces accurate output on the target hardware (Intel i5-4210U, 16GB RAM, no GPU); brief UI stutter observed only under simultaneous heavy multitasking, not disqualifying for sequential use.
- **Run the offline demo with real users (lawyer)** (#56) — **still open**, carried forward to v0.9.2.

**Deviation from original scope, noted explicitly rather than silently:** validating #55 required testing `/kb ask` output quality, which surfaced a real chunking bug and a hard limit on free-form LLM enumeration reliability. Fixing the second problem required an architectural change (see "RAG Query Routing" above), not a config tweak - so RAG completion, originally item 1 of the v1.0 backlog below, was substantially delivered on this branch as a direct consequence of investigating the hardware question, not as planned scope expansion. The free-form LLM path (item 1's original design) still exists, but only as a fallback.

---

# Version 1.0 Vision

Scope below is a *backlog*, not a commitment. Updated to reflect that RAG completion has substantially landed as of v0.9.1/v0.9.2 (see above) - remaining priority order:

1. ~~**RAG completion**~~ — substantially delivered as a deterministic dual-path design (see "RAG Query Routing" above), not the free-form LLM design originally envisioned. Remaining under this heading: revisit after #56.
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