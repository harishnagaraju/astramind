# AstraMind Platform API v1

The Platform API is the domain-neutral HTTP interface for AstraMind. v1 is intentionally small: it exposes the existing engine capabilities through a stable HTTP boundary without creating a second provider, document, embedding, or search subsystem.

## Endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/api/v1/health` | Liveness check |
| GET | `/api/v1/status` | Runtime provider/model status |
| GET | `/api/v1/version` | AstraMind version |
| GET | `/api/v1/providers` | Active provider discovery |
| GET | `/api/v1/models` | Active model discovery |
| POST | `/api/v1/chat` | Non-streaming chat |
| POST | `/api/v1/chat/stream` | SSE chat streaming |
| GET | `/api/v1/documents` | List knowledge documents |
| POST | `/api/v1/documents` | Import `.txt`, `.md`, or `.docx` |
| GET | `/api/v1/documents/{id}` | Retrieve a document |
| DELETE | `/api/v1/documents/{id}` | Remove a document |
| POST | `/api/v1/documents/{id}/analyze` | Analyze a document with the active model |
| POST | `/api/v1/embeddings` | Generate one text embedding |
| POST | `/api/v1/search` | Keyword or semantic knowledge-base search |

## Chat

Request:

```json
{
  "messages": [
    {"role": "user", "content": "Hello"}
  ]
}
```

Optional fields are `model` and `temperature`. The configured runtime model is used when `model` is omitted.

## Streaming

`POST /api/v1/chat/stream` returns Server-Sent Events:

```
event: token
data: {"content":"Hello"}

event: done
data: {}
```

A provider must support the existing AstraMind `StreamingProvider` capability.

## Documents

Document upload uses multipart form data with field name `file`. Supported formats are `.txt`, `.md`, and `.docx`. Import uses the existing Knowledge Base manager, chunking, storage, and embedding behavior.

## Search

Request:

```json
{"query":"architecture","mode":"keyword"}
```

Use `mode: "semantic"` for the existing embedding-based search. Omitting mode uses keyword search.

## Embeddings

Request:

```json
{"text":"AstraMind platform"}
```

The endpoint delegates to the active provider's existing embedding capability.

## Request IDs and errors

Every Platform API response receives an `X-Request-ID`. A supplied request ID is propagated; otherwise one is generated.

Errors use:

```json
{
  "code": "bad_request",
  "message": "messages is required",
  "request_id": "..."
}
```

## Compatibility and scope

The legacy local web API (`/api/status`, `/api/documents`, `/api/ask`) remains unchanged.

This v1 deliberately does not introduce agents, tool orchestration, API-key management, RBAC, vector databases, reranking, distributed deployment, SDKs, or advanced observability. Those are future upgrades only when an actual application requires them.

**Completion target:** Platform API v1 is considered complete when the endpoints above pass CI and existing AstraMind CLI/web/KB/provider tests remain green. Further development should return to product work, including SkillSifter, rather than expanding the v1 scope.
