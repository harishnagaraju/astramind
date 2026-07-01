</>

# AstraMind
Intelligent Conversations. Infinite Possibilities.

Created and maintained by Harish Nagaraju.

[![Go Build](https://github.com/harishnagaraju/astramind/actions/workflows/go.yml/badge.svg)](https://github.com/harishnagaraju/astramind/actions)
![Go](https://img.shields.io/badge/Go-1.24+-blue)
![License](https://img.shields.io/badge/License-Apache%202.0-green)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)
Current Stable Release: **v0.4.1**

AstraMind is a modular, AI-powered command-line assistant built in Go that provides a clean, scalable foundation for developing intelligent applications using Large Language Models (LLMs). Designed with a production-ready architecture, it supports multiple AI providers, conversation and session management, persistent chat history, export capabilities, automated testing, and a modern CI/CD pipeline.

Built with simplicity, maintainability, and extensibility in mind, AstraMind demonstrates best practices for integrating AI into real-world applications through clean Go code, modular components, and industry-standard APIs. While it begins as a lightweight CLI assistant, its architecture is designed to evolve into a comprehensive AI platform capable of supporting web applications, local and cloud-based LLMs, Retrieval-Augmented Generation (RAG), streaming responses, AI agents, enterprise copilots, knowledge systems, and domain-specific AI solutions.

# Vision

AstraMind aims to become a flexible AI platform that combines conversational intelligence, knowledge retrieval, automation, and decision support into a single extensible ecosystem. The project emphasizes simplicity, performance, and scalability while leveraging the strengths of the Go programming language.

# Current Features

## AI Assistant
- Multi-provider AI architecture
- OpenAI-compatible API integration
- Mock AI provider
- Automatic provider failover
- Interactive command-line interface (CLI)
- Environment-based configuration

## Conversation Management
- Conversation memory
- Persistent conversation history
- Multi-session support
- Active session tracking
- Session statistics

## Session Commands
- /sessions
- /new <session>
- /load <session>
- /delete <session>
- /export
- /export md

## Utility Commands
- /help
- /provider
- /history
- /clear
- /stats
- /config
- /about

## Developer Features
- Modular Go architecture
- GitHub Actions CI/CD
- Automated testing
- Coverage reporting
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

## Available Commands
   
| Command        |            Description            |
|----------------|-----------------------------------|
| /help          | Show help                         |
| /about         | About AstraMind                   |
| /history       | Conversation history              |
| /clear         | Clear current conversation        |
| /stats         | Session statistics                |
| /config        | Configuration                     |
| /sessions      | List sessions                     |
| /new <name>    | Create session                    |
| /load <name>   | Load session                      |
| /delete <name> | Delete session                    |
| /export        | Export current session as TXT     |
| /export md     | Export current session as Markdown|
| /provider      | Show AI provider                  |


# Release Management

## v0.4.1
### Highlights
- Conversation Export (TXT & Markdown)
- AI Provider Framework
- Automatic Provider Failover
- Mock AI Provider
- OpenAI Provider
- GitHub Actions CI
- Regression Test Suite
## v0.4.0
**Multi-session support**
- Active session tracking.
- Session listing.
- Session creation.
- Session loading.
- Session deletion.
- Improved API error handling.
## v0.3.0
**Persistence and usability improvements**
- Persistent history storage.
- Statistics command.
- Configuration command.
- About command.
- Improved architecture.
## v0.2.x
**Conversation management**
- Conversation memory.
- History support.
- Chat commands.
## v0.1.x
**Initial AstraMind foundation**
- CLI chatbot.
- OpenAI integration.
- Environment configuration.

## Technology Stack

- Go 1.24+
- OpenAI-compatible APIs
- JSON
- REST APIs
- GitHub Actions
- Git
- Markdown

## Architecture

AstraMind follows a modular architecture consisting of:

- AI Provider Framework
- Storage Layer
- Chat Engine
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

