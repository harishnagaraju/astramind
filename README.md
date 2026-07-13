</>

# AstraMind
Intelligent Conversations. Infinite Possibilities.

Created and maintained by Harish Nagaraju.

[![Go Build](https://github.com/harishnagaraju/astramind/actions/workflows/go.yml/badge.svg)](https://github.com/harishnagaraju/astramind/actions)
![Go](https://img.shields.io/badge/Go-1.24+-blue)
![License](https://img.shields.io/badge/License-Apache%202.0-green)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)
Current Stable Release: **v0.8.0**

AstraMind is a modular, AI-powered command-line assistant built in Go that provides a clean, scalable foundation for developing intelligent applications using Large Language Models (LLMs). Designed with a production-ready architecture, it supports multiple AI providers, conversation and session management, persistent chat history, export capabilities, automated testing, and a modern CI/CD pipeline.

Built with simplicity, maintainability, and extensibility in mind, AstraMind demonstrates best practices for integrating AI into real-world applications through clean Go code, modular components, and industry-standard APIs. While it begins as a lightweight CLI assistant, its architecture is designed to evolve into a comprehensive AI platform capable of supporting web applications, local and cloud-based LLMs, Retrieval-Augmented Generation (RAG), streaming responses, AI agents, enterprise copilots, knowledge systems, and domain-specific AI solutions.

# Vision

AstraMind aims to become a flexible AI platform that combines conversational intelligence, knowledge retrieval, automation, and decision support into a single extensible ecosystem. The project emphasizes simplicity, performance, and scalability while leveraging the strengths of the Go programming language.

# Current Features

## Knowledge Base

- Import text and Markdown documents
- Automatic document chunking
- Persistent JSON document storage
- Persistent chunk storage
- Keyword search
- Knowledge Base statistics
- Knowledge Base management

## Multiple AI Providers
- OpenAI
- OpenRouter
- Ollama
- Mock AI

## AI Providers
- Multi-provider AI architecture
- OpenAI-compatible API integration
- OpenRouter integration
- Native Ollama integration
- Local LLM support
- Configurable API base URL
- Configurable AI model selection
- Runtime provider selection
- Mock AI provider
- Automatic provider failov

## Streaming Responses
- Real-time token streaming
- Streaming renderer
- OpenAI streaming
- Ollama streaming
- Automatic fallback to non-streaming providers
- Mock streaming provider
- End-to-end streaming support

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
- `/kb search <text>`   Search the knowledge base
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
- Modular Go architecture
- Provider abstraction layer
- Streaming provider interface
- Configurable request builder
- Structured error handling
- GitHub Actions CI/CD
- Automated testing
- Unit tests
- Integration tests
- Streaming integration tests
- HTTP error integration tests
- Mock provider testing
- Regression test suite
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
AstraMind supports any OpenAI-compatible API endpoint through the OPENAI_BASE_URL configuration.

## Environment Variables

### OpenAI

```env
AI_PROVIDER=openai
OPENAI_API_KEY=your_api_key
OPENAI_MODEL=gpt-4o-mini
OPENAI_BASE_URL=https://api.openai.com/v1
```

### OpenRouter

```env
AI_PROVIDER=openrouter
OPENAI_API_KEY=your_api_key
OPENAI_MODEL=openrouter/auto
OPENAI_BASE_URL=https://openrouter.ai/api/v1
```

### Ollama

```env
AI_PROVIDER=ollama
OPENAI_MODEL=gemma3:1b
OPENAI_BASE_URL=http://localhost:11434
```
## Using Ollama

Install Ollama

```bash
https://ollama.com
```

Download a model

```bash
ollama pull gemma3:1b
```

Run the model

```bash
ollama run gemma3:1b
```

## Testing

Run the complete test suite:

```bash
go test -v ./...
```

The project includes:

- Unit tests
- Provider integration tests
- Streaming integration tests
- HTTP error tests
- Invalid JSON tests
- Connection failure tests
- Mock provider tests
- Regression tests

## Supported AI Providers

| Provider   | Local | Streaming | Status |
|------------|:-----:|:---------:|:------:|
| Mock AI    | N/A   | ✅        | ✅     |
| OpenAI     | ❌    | ✅        | ✅     |
| OpenRouter | ❌    | ✅        | ✅     |
| Ollama     | ✅    | ✅        | ✅     |


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
| /search <text> 	| Search current converstation      |
| /searchall <text>	| Search across all sessions        |

## Knowledge Base

AstraMind includes a built-in Knowledge Base that allows users to import local
text and Markdown documents for later retrieval.

Imported documents are automatically:

- assigned a unique identifier
- split into chunks
- stored on disk
- indexed for keyword search

The Knowledge Base is designed as the foundation for Retrieval-Augmented
Generation (RAG) in future releases.

### Example

/kb import notes.txt

Imported: notes.txt

/kb list

Knowledge Base Documents
------------------------
8f3a2d...

Name : notes.txt
Chunks : 1

/kb search architecture

Found 3 matching chunks

/kb stats

Documents : 1
Chunks : 1

/kb remove 8f3a2d...

Removed

/kb clear

Knowledge base cleared.


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
    ai/
    chat/
    config/
    models/
    renderer/
    storage/
    testutil/

exports/
data/
```

# Release Management

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
- OpenAI APIs
- OpenRouter
- Ollama
- JSON
- REST APIs
- GitHub Actions
- Git
- Markdown

## Architecture
AstraMind follows a modular architecture consisting of:

- AI Provider Framework
- Provider Factory
- Provider Manager
- Streaming Framework
- Storage Layer
- Chat Engine
- Renderer
- Configuration
- Export System
- Test Utilities

## License
Apache License 2.0

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.

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

