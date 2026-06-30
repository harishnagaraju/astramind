</>

# AstraMind
Intelligent Conversations. Infinite Possibilities.

Created and maintained by Harish Nagaraju.

[![Go Build](https://github.com/harishnagaraju/astramind/actions/workflows/go.yml/badge.svg)](https://github.com/harishnagaraju/astramind/actions)
![Go](https://img.shields.io/badge/Go-1.24+-blue)
![License](https://img.shields.io/badge/License-Apache%202.0-green)
![Status](https://img.shields.io/badge/Status-Active-brightgreen)

AstraMind is a modular AI-powered command-line assistant written in Go. It provides a clean architecture for integrating multiple Large Language Model (LLM) providers while supporting session management, conversation persistence, export capabilities, automated testing, and a production-grade CI pipeline. AstraMind is designed to evolve into a complete AI platform supporting local models, search, streaming responses, and Retrieval-Augmented Generation (RAG).

# Vision

AstraMind aims to become a flexible AI platform that combines conversational intelligence, knowledge retrieval, automation, and decision support into a single extensible ecosystem. The project emphasizes simplicity, performance, and scalability while leveraging the strengths of the Go programming language.

# Introduction
AstraMind is a lightweight AI-powered chatbot built in Go that enables natural language conversations through Large Language Models (LLMs). Designed as a simple CLI application, AstraMind serves as a foundation for building advanced AI assistants, knowledge systems, and autonomous agents. AstraMind brings the power of Large Language Models (LLMs) to a simple command-line interface. It demonstrates how to integrate AI capabilities into applications using clean, maintainable Go code and industry-standard APIs.
The project is designed as both a learning platform and a foundation for future AI-powered products. Starting as a lightweight CLI chatbot, AstraMind can evolve into a web-based assistant, Retrieval-Augmented Generation (RAG) system, enterprise copilot, multi-agent platform, or domain-specific AI solution.

## Current Features

### AI
- Multi-provider AI architecture
- OpenAI provider
- Mock AI provider
- Automatic provider failover

### Conversation Management
- Persistent conversation history
- Multi-session support
- Session creation, loading, deletion
- Session statistics

### Export
- TXT export
- Markdown export

### Developer Experience
- GitHub Actions CI
- Automated testing
- Coverage reporting
- Regression test suite
- Semantic versioning

### Cross Platform
- Windows
- Linux
- macOS

# Release Management

 ## v0.4.1 (Released)
    **Export System - Planned:**
      •	Export current session.
      •	TXT export.
      •	Markdown export.
      •	Session export.
  ## v0.4.0
    **Multi-session support**
      •	Active session tracking.
      •	Session listing.
      •	Session creation.
      •	Session loading.
      •	Session deletion.
      •	Improved API error handling.
  ## v0.3.0
    **Persistence and usability improvements**
      •	Persistent history storage.
      •	Statistics command.
      •	Configuration command.
      •	About command.
      •	Improved architecture.
   ## v0.2.x
    **Conversation management**
      •	Conversation memory.
      •	History support.
      •	Chat commands.
  ## v0.1.x
    **Initial AstraMind foundation**
      •	CLI chatbot.
      •	OpenAI integration.
      •	Environment configuration.

  ## Quick Start
    ```bash
    git clone https://github.com/harishnagaraju/astramind.git
    
    cd astramind
    
    go build ./cmd/astramind
    
    ./astramind

 ### Available Commands
 ```markdown
  
         | Command        |            Description     |
         |----------------|----------------------------|
         | /help          | Show help                  |
         | /about         | About AstraMind            |
         | /history       | Conversation history       |
         | /clear         | Clear current conversation |
         | /stats         | Session statistics         |
         | /config        | Configuration              |
         | /sessions      | List sessions              |
         | /new <name>    | Create session             |
         | /load <name>   | Load session               |
         | /delete <name> | Delete session             |
         | /export        | Export TXT                 |
         | /export md     | Export Markdown            |
         | /provider      | Show AI provider           |

### Features Available till now

    **AI Assistant**
        • OpenAI-compatible API integration.
        • Interactive command-line chat.
        • Environment-based configuration.
        • Cross-platform support.
    
    **Conversation Management**
        • Conversation memory.
        • Persistent chat history.
        • Session-aware storage.
        • Active session tracking.
    
    **Session Commands**
        •	/sessions
        •	/new <session>
        •	/load <session>
        •	/delete <session>
    
    **Utility Commands**
        •	/help
        •	/history
        •	/clear
        •	/stats
        •	/config
        •	/about
    
    **Developer Features**
        •	GitHub Actions CI/CD
        •	GitHub Issues & Milestones
        •	Release Management
        •	Semantic Versioning
        •	Modular Go Architecture


## License

Apache License 2.0

This project is licensed under the Apache License 2.0 - see the LICENSE file for details.



# About
Why AstraMind?
The name combines:
Astra — derived from the Latin word for "stars," symbolizing exploration, intelligence, and limitless possibilities.
Mind — representing reasoning, learning, and artificial intelligence.
Together, AstraMind represents an AI system designed to explore knowledge, assist users, and continuously evolve toward more advanced forms of intelligence.




## Project Owner

Harish Nagaraju

Software Architect | AI Engineer | Founder, RK Consulting

- Website: https://www.rkconsulting.co.in
- GitHub: https://github.com/harishnagaraju

## Maintainers
AstraMind is designed and maintained by Harish Nagaraju.

| Name | Role |
|--------|--------|
| Harish Nagaraju | Creator & Lead Developer |

harishnagaraju@rkconsulting.co.in
harishnagaraju@gmail.com

See [Roadmap](docs/roadmap.md) for upcoming releases and features.

