# AstraMind Platform API v1

The Platform API is the domain-neutral HTTP interface for AstraMind.

## Phase 1 endpoints

| Method | Endpoint | Purpose |
|---|---|---|
| GET | `/api/v1/health` | Liveness/health check |
| GET | `/api/v1/status` | Runtime provider and model status |
| GET | `/api/v1/version` | AstraMind version |

## Request IDs

Every Platform API response includes `X-Request-ID`.

If the client sends an `X-Request-ID` request header, AstraMind propagates that value. Otherwise, the API generates a request ID.

## Error response

Errors use the following JSON shape:

```json
{
  "code": "method_not_allowed",
  "message": "method not allowed",
  "request_id": "..."
}
```

## Compatibility

The versioned Platform API is mounted alongside the existing legacy local web API. Existing routes such as `/api/status`, `/api/documents`, and `/api/ask` are not replaced by Phase 1.

## Scope

Phase 1 establishes the API foundation only. Chat, streaming, documents, embeddings, and search endpoints are intentionally deferred until this foundation passes validation.
