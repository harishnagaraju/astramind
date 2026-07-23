# AstraMind
Intelligent Conversations. Infinite Possibilities.

Created and maintained by Harish Nagaraju.

[![Go Build](https://github.com/harishnagaraju/astramind/actions/workflows/go.yml/badge.svg)](https://github.com/harishnagaraju/astramind/actions)
![Go](https://img.shields.io/badge/Go-1.24+-blue)
![License](https://img.shields.io/badge/License-AGPL--3.0-green)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)
Current Release: **v0.9.0**

AstraMind is a modular, AI-powered command-line assistant built in Go that provides a clean, scalable foundation for developing intelligent applications using Large Language Models (LLMs). Designed with a production-ready architecture, it supports multiple AI providers, conversation and session management, persistent chat history, a local Knowledge Base with semantic search and Retrieval-Augmented Generation (RAG), a local web interface, export capabilities, automated testing, and a modern CI/CD pipeline.

Built with simplicity, maintainability, and extensibility in mind, AstraMind demonstrates best practices for integrating AI into real-world applications through clean Go code, modular components, and industry-standard APIs. It runs fully offline via Ollama, or against any OpenAI-compatible cloud endpoint, using the same code path either way.

# Vision

AstraMind aims to become a flexible AI platform that combines conversational intelligence, knowledge retrieval, automation, and decision support into a single extensible ecosystem. The project emphasizes simplicity, performance, and scalability while leveraging the strengths of the Go programming language.

# Current Features

## Knowledge Base & RAG

- Import text and Markdown documents
- Automatic document chunking
- Persistent JSON document and chunk storage
- Keyword search
- **Semantic search** - embedding-based, ranks results by meaning rather than exact wording
- **Retrieval-Augmented Generation (`/kb ask`)** - retrieves relevant content and answers, citing sources with every response. List and single-fact questions are answered deterministically (chunk content returned directly, no LLM call, no risk of a fact being dropped or misstated); free-form LLM synthesis is used only as a fallback when no embedder is configured. See "How `/kb ask` answers a question" below.
- Knowledge Base statistics and management

## Local Web Interface

- `--web` flag launches a local, browser-based UI as an alternative to the CLI
- Import documents and ask questions from a browser, no terminal required
- Single embedded binary - no separate install, works fully offline when using Ollama
- Uses the exact same backend and provider configuration as the CLI

## Multiple AI Providers
- OpenAI (and any OpenAI-compatible endpoint, including OpenRouter - see below)
- Ollama (local)
- Mock AI (for testing/offline development)

## AI Providers
- Multi-provider AI architecture
- OpenAI-compatible API integration
- Native Ollama integration
- Local LLM support
- Embedding support (Ollama, OpenAI-compatible, and Mock)
- Configurable API base URL
- Configurable AI model selection
- Runtime provider selection
- Automatic provider failover

## Streaming Responses
- Real-time token streaming
- Streaming renderer
- OpenAI streaming
- Ollama streaming
- Automatic fallback to non-streaming providers
- Mock streaming provider

## Conversation Management
- Conversation memory
- Persistent conversation history
- Multi-session support
- Active session tracking
- Session statistics
- Automatic history persistence
- Configurable conversation memory limits

### Conversation Search

- Search current conversation
- Search across all saved sessions
- Case-insensitive search
- Session-grouped search results

## Knowledge Base Commands
- `/kb import <file>`   Import a text or markdown document
- `/kb list`            List imported documents
- `/kb search <text>`   Search the knowledge base (keyword)
- `/kb ssearch <text>`  Search the knowledge base (semantic, embedding-based)
- `/kb ask <question>`  Ask a question, answered from your knowledge base (RAG)
- `/kb remove <id>`     Remove a document
- `/kb clear`           Remove all documents
- `/kb stats`           Show knowledge base statistics

## Session Commands
- `/sessions`
- `/new <session>`
- `/load <session>`
- `/delete <session>`
- `/export`
- `/export md`

## Utility Commands
- `/help`
- `/provider`
- `/history`
- `/search`
- `/searchall`
- `/clear`
- `/stats`
- `/config`
- `/about`

## Session Export 
- Plain text export
- Markdown export

## Developer Features
- Modular Go architecture (`engine` / `features` / `infrastructure` layers)
- Dependency injection - every feature service constructed once at startup and shared
- Provider abstraction layer (chat, streaming, and embeddings)
- Configurable request builder
- Structured error handling
- GitHub Actions CI/CD
- Automated testing (unit, integration, and a full manual walkthrough script)
- GitHub Issues & Milestones
- Release management
- Semantic versioning

## Cross-Platform Support
- Windows
- Linux
- macOS

## Quick Start
    
```bash
git clone https://github.com/harishnagaraju/astramind.git

cd astramind

go build ./cmd/astramind

./astramind
```

To use the local web interface instead of the CLI:

```bash
./astramind --web
```

This starts a local server and opens your default browser automatically.

AstraMind supports any OpenAI-compatible API endpoint through the `OPENAI_BASE_URL` configuration.

## Environment Variables

### OpenAI

```env
AI_PROVIDER=openai
OPENAI_API_KEY=your_api_key
OPENAI_MODEL=gpt-4o-mini
OPENAI_BASE_URL=https://api.openai.com/v1
```

### OpenRouter

OpenRouter is OpenAI-compatible, so it's configured the same way as OpenAI - point `OPENAI_BASE_URL` at OpenRouter and use an OpenRouter API key and model ID:

```env
AI_PROVIDER=openai
OPENAI_API_KEY=your_openrouter_api_key
OPENAI_MODEL=openai/gpt-4o-mini
OPENAI_BASE_URL=https://openrouter.ai/api/v1
```

Browse available models (including free-tier options) at https://openrouter.ai/models.

### Ollama

```env
AI_PROVIDER=ollama
OPENAI_MODEL=gemma3:1b
OPENAI_BASE_URL=http://localhost:11434
```

For semantic search and `/kb ask` to work with Ollama, an embedding model is also required:

```bash
ollama pull nomic-embed-text
```

## Using Ollama

Install Ollama

```bash
https://ollama.com
```

Download a model

```bash
ollama pull gemma3:1b
ollama pull nomic-embed-text
```

Run the model

```bash
ollama run gemma3:1b
```

# Testing

AstraMind includes a comprehensive automated testing framework covering unit, integration, and manual end-to-end testing.

## Unit Tests

Run the complete unit test suite:

```bash
go test -v ./...
```

The unit test suite includes:

* Unit tests across every package
* Provider integration tests (chat, streaming, and embeddings)
* HTTP error handling tests
* Invalid JSON handling tests
* Connection failure tests
* Mock provider tests
* Regression tests

## Integration Tests

Run the complete integration test suite.

### Windows

```cmd
tests\integration\run_all.bat
```

### Linux / macOS

```bash
./tests/integration/run_all.sh
```

The integration suite automatically performs:

* Code formatting validation
* Static analysis (`go vet`)
* Project build verification
* Complete unit test execution
* End-to-end Knowledge Base workflow testing

## Knowledge Base Integration Test

Run only the Knowledge Base integration tests.

### Windows

```cmd
tests\integration\run_kb.bat
```

### Linux / macOS

```bash
./tests/integration/run_kb.sh
```

## Manual Testing

A full interactive walkthrough script covers every command, including a live comparison of keyword vs. semantic search, the RAG loop, and a `--web` API smoke test:

```bash
bash tests/integration/manual_testing.sh
```

## Script Execution

AstraMind supports automated command execution using script files, routed through the same command dispatcher as interactive mode.

Example:

### Windows

```cmd
astramind.exe --script tests\integration\commands\kb.txt
```

### Linux / macOS

```bash
./astramind --script tests/integration/commands/kb.txt
```

## Supported AI Providers

| Provider   | Local | Streaming | Embeddings | Status |
|------------|:-----:|:---------:|:----------:|:------:|
| Mock AI    | N/A   | ✅        | ✅         | ✅     |
| OpenAI     | ❌    | ✅        | ✅         | ✅     |
| OpenRouter | ❌    | ✅        | Varies by model | ✅ |
| Ollama     | ✅    | ✅        | ✅         | ✅     |


## Available Commands
   
| Command        	|            Description            |
|-------------------|-----------------------------------|
| /help          	| Show help                         |
| /about         	| About AstraMind                   |
| /history       	| Conversation history              |
| /clear         	| Clear current conversation        |
| /stats         	| Session statistics                |
| /config        	| Configuration                     |
| /sessions      	| List sessions                     |
| /new <name>    	| Create session                    |
| /load <name>   	| Load session                      |
| /delete <name> 	| Delete session                    |
| /export        	| Export current session as TXT     |
| /export md     	| Export current session as Markdown|
| /provider      	| Show AI provider                  |
| /search <text> 	| Search current conversation       |
| /searchall <text>	| Search across all sessions        |

## Knowledge Base

AstraMind includes a built-in Knowledge Base that allows users to import local
text and Markdown documents for later retrieval.

Imported documents are automatically:

- assigned a unique identifier
- split into chunks
- embedded (if an embedding-capable provider is configured)
- stored on disk
- indexed for both keyword and semantic search

### Example

```
/kb import notes.txt

Imported: notes.txt
(1/1 chunks embedded)

/kb list

Knowledge Base Documents
------------------------
8f3a2d...

Name : notes.txt
Chunks : 1

/kb search architecture

Found 3 matching chunks

/kb ssearch how does the system handle errors

Semantic Search Results
------------------------
[8f3a2d...] (similarity: 0.812)
...matching chunk content by meaning, not exact wording...

/kb ask how does the system handle errors

Here is everything found in your knowledge base across 1 relevant section(s):

* ...matching chunk content, returned verbatim - no LLM paraphrase for
  list/single-fact questions like this one...

Sources:
  [8f3a2d...]

/kb stats

Documents : 1
Chunks : 1

/kb remove 8f3a2d...

Removed

/kb clear

Knowledge base cleared.
```

### How `/kb ask` answers a question

As of the v0.9.1 branch, `/kb ask` no longer sends every question to the LLM. It routes based on question shape:

- **"List everything matching X" questions** ("what are all the...", "what are the... timings") are answered deterministically: every relevant chunk found is returned in full, formatted as a list, with no LLM call and no possibility of an item being silently dropped during generation.
- **Single-fact questions** ("is there...", "what time is...", "how much...", "what is the...") are also answered deterministically when an embedding provider is configured: the single most relevant chunk is returned verbatim, so a date, fee, or ID can't be misstated by a model paraphrasing it.
- **The free-form LLM path is retained only as a fallback** - used when no embedding provider is configured. It carries the caveat below.

This replaced an earlier single-path design after validation testing found that free-form LLM enumeration was unreliable in a way prompt wording alone couldn't fix - see [CHANGELOG.md](CHANGELOG.md) (v0.9.1) for the full investigation, including a real measured case where similarity-based sentence windowing was tried and abandoned.

### A note on the LLM fallback path

When the LLM fallback path is used (no embedder configured), answer quality still depends on the capability of the active AI provider. Small local models (e.g. `gemma3:1b`, ~1B parameters) can omit or, in some cases, fabricate details even when the correct source text is present in context - this is a model capability limitation, not a defect in AstraMind's retrieval or prompt construction. For accuracy-critical use on this fallback path, a larger local model or a cloud provider is recommended. See [CHANGELOG.md](CHANGELOG.md) for details on how this was identified and verified.

## Streaming Responses
AstraMind supports real-time streaming responses from compatible AI providers.

Benefits include:

- Lower perceived latency
- Token-by-token rendering
- Improved interactive experience
- Automatic fallback for providers without streaming support

# Project Structure
```
cmd/
    astramind/
internal/
    engine/            - command dispatch, app lifecycle, web server
        webui/         - embedded local web interface
    features/
        chat/          - conversation + knowledge base command handling
        kb/            - knowledge base, chunking, search, RAG
        history/       - conversation history persistence
        session/       - session lifecycle
        search/        - conversation search
        export/        - session export
    infrastructure/
        ai/            - provider abstraction (chat, streaming, embeddings)
        storage/       - file-based persistence
        models/        - shared data types
        renderer/      - terminal output
        config/        - runtime configuration
    testutil/

exports/
data/
```

# Release Management

## v0.9.1 (validation branch, not yet merged to main)
**Deterministic RAG & Real-Hardware Validation**

Originally scoped as validation-only (no new features); real findings during validation required an architectural fix. Documented here rather than silently expanded.

- Closed: `gemma2:9b` validated on real hardware (Intel i5-4210U, 16GB RAM, no GPU) - produces correct output; usable for sequential use, brief stutter only under heavy simultaneous multitasking.
- `/kb ask` now routes enumeration and single-fact questions to deterministic, zero-LLM-call extraction from retrieved chunks, rather than trusting free-form LLM generation to enumerate correctly or restate a fact exactly. The original free-form LLM path is retained as a fallback for single-fact questions only, when no embedder is configured.
- Fixed a real document-chunking bug: byte-offset splitting could corrupt content mid-word, and a CRLF-encoded real-world file bypassed an earlier paragraph-aware fix entirely (a `\r\n\r\n` blank line never matches a `\n\n` split point). Both fixed; both verified against real file bytes, not synthetic test fixtures.
- Added content-fidelity and determinism checks to the manual test walkthrough.
- Two sentence-level windowing approaches for single-fact answers (fixed-size, then embedding-similarity-based) were tried and abandoned after live testing against real embeddings showed no reliable topic-boundary signal at sentence granularity - see [CHANGELOG.md](CHANGELOG.md) for the measured data.
- Still open: real-user (niece/lawyer) offline demo feedback.
- See [CHANGELOG.md](CHANGELOG.md) for the complete list and full investigation detail.

## v0.9.0
**Semantic Search, RAG & Local Web Interface**

- Embedding-based semantic search (`/kb ssearch`), alongside keyword search.
- Retrieval-Augmented Generation (`/kb ask`) - completes the pipeline from import through answer generation, with citations.
- Local web interface (`--web`) - browser-based UI for non-technical users, same backend as the CLI.
- Embedding provider support for Ollama, OpenAI-compatible endpoints, and Mock.
- Architectural cleanup: consistent dependency injection across all feature services, decoupled search/history/session packages, `--script` mode routed through the full command dispatcher.
- Injectable, test-isolated storage for history and sessions.
- Fixed: Ollama RAG answers truncating mid-generation (context window was unset); unbounded semantic search result count; silent embedding failures on import.
- First unit tests for the `history` and `session` packages.
- See [CHANGELOG.md](CHANGELOG.md) for the complete list, including a documented model-capability limitation found through controlled testing.

## v0.8.0
**Knowledge Base**

- Built-in Knowledge Base framework.
- Document import for text and Markdown files.
- Automatic document chunking with configurable chunk size and overlap.
- Persistent document storage using JSON.
- Persistent chunk storage.
- Repository layer for chunk management.
- Keyword-based Knowledge Base search.
- Prompt builder for Retrieval-Augmented Generation (RAG).
- Knowledge Base management API.
- Knowledge Base statistics.
- CLI commands:
  - `/kb import`
  - `/kb list`
  - `/kb search`
  - `/kb remove`
  - `/kb clear`
  - `/kb stats`
- Improved storage architecture for future database backends.
- Comprehensive unit tests for Knowledge Base components.
- End-to-end CLI validation for Knowledge Base workflows.
- Foundation for future RAG integration.

## v0.7.0
**Native Ollama Integration**

- Native Ollama provider.
- Local Large Language Model (LLM) support.
- Runtime provider selection.
- Local model configuration.
- Ollama streaming support.
- Streaming parser for Ollama responses.
- Shared request builder.
- Integration tests for Ollama.
- HTTP error handling tests.
- Invalid JSON handling tests.
- Connection failure tests.
- Streaming integration tests.
- Verified with Gemma3:1b.
- Improved provider architecture

## v0.6.0
**Search System**
- Search current conversation using `/search`.
- Search across all saved sessions using `/searchall`.
- Case-insensitive conversation search.
- Session-grouped search results.
- Dedicated search result renderer.
- Reusable search engine for current and multi-session search.
- Search result data models for extensibility.
- Unit tests for conversation search.
- Integration tests for multi-session search.
- Improved conversation persistence across sessions.
- Improved search result presentation

## v0.5.0
**Streaming Responses & Provider Enhancements**
- Real-time AI streaming responses.
- Streaming renderer for token-by-token output.
- OpenRouter integration.
- Configurable API base URL.
- Configurable AI model selection.
- Shared HTTP request builder.
- Improved HTTP error handling.
- Mock streaming provider.
- End-to-end streaming integration.
- HTTP integration tests using `httptest.Server`.
- Streaming integration tests.
- HTTP error integration tests.
- Enhanced provider architecture.
- Improved OpenAI-compatible API support.

## v0.4.1
**AI Provider Framework & Export**
- Conversation export (TXT & Markdown).
- AI provider abstraction.
- Automatic provider failover.
- Mock AI provider.
- OpenAI provider implementation.
- GitHub Actions CI pipeline.
- Regression test suite.

## v0.4.0
**Multi-Session Support**
- Active session tracking.
- Session listing.
- Session creation.
- Session loading.
- Session deletion.
- Improved API error handling.

## v0.3.0
**Persistence & Usability Improvements**
- Persistent conversation history.
- Session statistics.
- Configuration command.
- About command.
- Improved internal architecture.

## v0.2.x
**Conversation Management**
- Conversation memory.
- Conversation history.
- Interactive chat commands.

## v0.1.x
**Initial Foundation**
- Command-line AI assistant.
- OpenAI integration.
- Environment-based configuration.
- Modular project structure.

## Technology Stack
- Go 1.24+
- OpenAI-compatible APIs (OpenAI, OpenRouter)
- Ollama
- JSON
- REST APIs
- GitHub Actions
- Git
- Markdown

## Architecture
AstraMind follows a modular, layered architecture:

- `engine` - command dispatch, dependency injection, application lifecycle, local web server
- `features` - chat, knowledge base (search/RAG), history, session, search, export
- `infrastructure` - AI provider framework (chat/streaming/embeddings), storage, models, renderer, configuration
- Test utilities

## License
GNU Affero General Public License v3.0 (AGPL-3.0)

This project is licensed under AGPL-3.0. *Note: the licensing terms are under review with legal counsel before any commercial release - the `LICENSE` file in this repository should be reconciled with this statement before v1.0.*

# About
## Why AstraMind?
The name combines:
Astra — derived from the Latin word for "stars," symbolizing exploration, intelligence, and limitless possibilities.
Mind — representing reasoning, learning, and artificial intelligence.
Together, AstraMind represents an AI system designed to explore knowledge, assist users, and continuously evolve toward more advanced forms of intelligence.

## Project Owner
**Harish Nagaraju**

Software Architect | AI Engineer | Founder, RK Consulting

- Website: https://www.rkconsulting.co.in
- GitHub: https://github.com/harishnagaraju

## Maintainers
AstraMind is designed and maintained by Harish Nagaraju.

| Name | Role |
|--------|--------|
| Harish Nagaraju | Creator & Lead Developer |

- harishnagaraju@rkconsulting.co.in
- harishnagaraju@gmail.com

## Project Planning

See [Roadmap](docs/roadmap.md) for upcoming releases and features.