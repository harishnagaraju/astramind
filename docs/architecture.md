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

---

# Core Components

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

## v0.8.0

Knowledge Base

- Document storage
- Document indexing
- Semantic search
- Question answering
- Retrieval-Augmented Generation (RAG)

---

# Version 1.0 Vision

AstraMind evolves into a complete AI platform featuring:

- Local and cloud AI providers
- Retrieval-Augmented Generation (RAG)
- Knowledge Base
- Long-term memory
- AI Agents
- Tool calling
- Plugin architecture
- Multi-modal support
- Personal productivity assistant

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