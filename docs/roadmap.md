
# AstraMind Roadmap

## Current Status

### Releases

* v0.1.0
* v0.2.1 (Latest)

### Completed Features

* OpenAI-powered CLI chatbot
* Conversation memory
* `/help` command
* `/history` command
* `/clear` command
* Memory size limiting
* Improved error handling
* Environment-driven model selection
* GitHub Actions CI/CD
* GitHub Releases
* Issue tracking workflow
* Feature branch workflow

### Closed Issues

* #1 Conversation Memory
* #2 Add /help command
* #3 Add /history command
* #4 Add /clear command

### Open Issues

* #5 Add streaming responses

---

# Version 0.3.0

## Issue #6: Persistent Chat History

### Goal

Persist conversation history across AstraMind sessions.

### Features

* Save conversations to JSON file
* Load conversations on startup
* Preserve chat context between executions
* Store data under:

```text
data/chat_history.json
```

### Benefits

* Long-term conversation continuity
* Better user experience
* Foundation for future memory features

---

## Issue #7: Session Statistics

### New Command

```text
/stats
```

### Example Output

```text
Messages Sent: 12
Messages Received: 12
Current Model: gpt-4o-mini
Memory Entries: 24
Session Duration: 18 minutes
```

### Benefits

* Better visibility into chatbot usage
* Useful for debugging and monitoring

---

## Issue #8: Configuration Command

### New Command

```text
/config
```

### Example Output

```text
Model: gpt-4o-mini
Max Messages: 20
History Enabled: true
History File: data/chat_history.json
```

### Benefits

* Easier troubleshooting
* Improved transparency

---

# Version 0.4.0

## Issue #5: Streaming Responses

### Goal

Display AI responses incrementally instead of waiting for the entire response.

### Current Behavior

```text
AI: Entire response appears at once
```

### Target Behavior

```text
AI: Hello...
AI: Hello, how...
AI: Hello, how are...
```

### Benefits

* Faster perceived response time
* More natural conversational experience

---

# Version 0.5.0

## Internal Refactoring

### Proposed Structure

```text
internal/
├── ai/
│   └── openai.go
├── chat/
│   ├── commands.go
│   ├── history.go
│   └── memory.go
└── config/
    └── config.go
```

### Goals

* Cleaner architecture
* Easier maintenance
* Better testability

---

# Version 1.0.0

## Personal AI Assistant

### Long-Term Features

* Persistent memory
* Knowledge base
* Notes management
* Task tracking
* Personal assistant capabilities
* Multi-session chat support
* Enhanced CLI experience

### Vision

AstraMind evolves from a simple CLI chatbot into a persistent personal AI assistant with memory, context awareness, and long-term knowledge retention.

---

## Development Workflow

1. Create Issue
2. Create Feature Branch
3. Implement Feature
4. Test Locally
5. Push Branch
6. Create Pull Request
7. Merge into Main
8. Create Release
9. Update Roadmap

---

## Next Immediate Action

Create branch:

```bash
git checkout -b feature/v0.3.0-persistence
```

Implement:

```text
Issue #6 - Persistent Chat History
```

Target Release:

```text
AstraMind v0.3.0
```
