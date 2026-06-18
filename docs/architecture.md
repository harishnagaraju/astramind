# AstraMind Architecture

## Overview

AstraMind is a command-line AI assistant written in Go.

The application provides conversational AI capabilities using OpenAI-compatible APIs and maintains session-level conversation memory.

---

## High-Level Architecture

```text
+-------------+
| User        |
+------+------+
       |
       v
+-------------+
| CLI Layer   |
| main.go     |
+------+------+
       |
       v
+-------------+
| Commands    |
| /help       |
| /history    |
| /clear      |
+------+------+
       |
       v
+-------------+
| Memory      |
| Conversation|
+------+------+
       |
       v
+-------------+
| AI Client   |
| OpenAI API  |
+------+------+
       |
       v
+-------------+
| AI Response |
+-------------+
```

---

## Current Project Structure

```text
astramind/
│
├── cmd/
│   └── astramind/
│       └── main.go
│
├── internal/
│   ├── ai/
│   │   └── openai.go
│   │
│   ├── chat/
│   │   └── chat.go
│   │
│   └── config/
│       └── config.go
│
├── docs/
│   ├── roadmap.md
│   └── architecture.md
│
├── .github/
├── README.md
├── go.mod
└── go.sum
```

---

## Current Capabilities

* OpenAI API integration
* Session conversation memory
* Command processing
* Environment-based configuration
* GitHub Actions CI/CD
* Release management

---

## Planned Architecture Evolution

### Version 0.3.x

Persistent storage:

```text
data/
└── chat_history.json
```

Capabilities:

* Load chat history
* Save chat history
* Restore context

---

### Version 0.4.x

Streaming support:

```text
CLI
  ↓
Streaming API
  ↓
Incremental output
```

---

### Version 0.5.x

Refactoring:

```text
internal/
├── ai/
│   └── openai.go
│
├── chat/
│   ├── commands.go
│   ├── history.go
│   └── memory.go
│
├── storage/
│   └── json.go
│
└── config/
    └── config.go
```

---

## Version 1.0 Vision

AstraMind evolves into a persistent personal AI assistant with:

* Long-term memory
* Knowledge storage
* Personal notes
* Tasks and reminders
* Multiple conversation sessions
* Extensible provider architecture

---

## Design Principles

* Simplicity first
* Minimal dependencies
* Clear separation of concerns
* Git-first workflow
* Incremental feature delivery
* Release-driven development
