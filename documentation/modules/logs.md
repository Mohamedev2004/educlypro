# Logs (Audit Logging) Documentation

## Overview

This codebase implements **audit logging via events**:

- Application code publishes `types.AuditEvent` messages to Watermill.
- A dedicated subscriber persists them as `logs.AuditLog` rows.
- Audit logs are written to a **separate database connection** (`LogsDB`).

The backend also exposes an authenticated **Logs HTTP API** under `/api/v1/logs` for listing, exporting, and serving chart aggregates. The super admin UI at `frontend/src/pages/super-admin/logs.tsx` consumes these endpoints.

## Data flow

### 1) Request correlation

- The frontend attaches an `X-Request-ID` header to every API request.
- Backend middleware ensures an id exists and injects it into `context.Context` so service code can read it.

### 2) Event publishing

Services publish audit events to:

- **`system.audit_logs`** (global sink topic)
- and, optionally, a more specific topic (e.g. `system.events.v1.auth.logged_in`)

The auth service is a primary publisher of these events.

**Important**: The codebase currently does **not** have a shared `types.PublishAuditEvent` helper. Publishing is implemented module-locally (see `backend/modules/auth/service.go`), but all publishers use the same `types.AuditEvent` struct and the same sink topic (`system.audit_logs`).

### 3) Event persistence

At runtime, `backend/cmd/main.go` wires:

- Watermill pub/sub (`gochannel`)
- a subscriber: `logs.Subscriber.ProcessLogEvent`
- a handler registered on the topic **`system.audit_logs`**

That subscriber:

- unmarshals the event payload
- marshals the dynamic `Payload` into JSON
- persists a `logs.AuditLog` row in `LogsDB`

Malformed events are dropped (no retry). Database failures return an error so Watermill can retry.

## Event schema (`types.AuditEvent`)

Fields:

- `request_id`: string
- `level`: `DEBUG | INFO | WARN | ERROR`
- `entity`: string (e.g. `"User"`, `"UserSession"`, `"Token"`)
- `entity_id`: string
- `actor_id`: string
- `action`: `CREATED | UPDATED | DELETED | FAILED`
- `payload`: arbitrary JSON object (sanitized when needed)
- `duration_ms`: float64 (derived from request start time)
- `timestamp`: time

Notes:

- `duration_ms` is computed from a `startTime` value injected by `RequestTimingMiddleware`. If missing, duration defaults to `0`.
- The system intentionally avoids logging plaintext secrets (e.g. passwords). Auth registration explicitly publishes a sanitized payload.

## Storage schema (`logs.AuditLog`)

Persisted columns:

- `request_id` (indexed)
- `level` (indexed)
- `entity` (indexed)
- `entity_id` (indexed)
- `actor_id` (indexed)
- `action`
- `payload` (`jsonb`)
- `duration_ms` (indexed)
- `created_at` (indexed)

## Logs HTTP API

All endpoints are under `/api/v1` and require `AuthMiddleware` + `RequireRole("super_admin")` — audit logs span every center and every user's actions, so this is enforced at the route group, not just hidden from other roles' UI. See:

- `backend/modules/logs/routes.go`
- `backend/modules/logs/handler.go`

### `GET /logs`

Lists logs with pagination and server-side filtering.

Query params:

- `page` (default `1`, must be >= 1)
- `per_page` (default `10`, allowed: `10|20|30|40|50|100`)
- `q` (optional, full-text-ish search across request id/entity/entity_id/actor_id/action/payload)
- `level` (optional array; supports both `level=value&level=value` and Axios `level[]=value`)
- `status` (optional array; maps to `action` in storage; supports array formats like `level`)
- `duration` (optional array; UI buckets: `"bigger than 200ms"` / `"less than 200ms"`)
- `from` / `to` (optional RFC3339 timestamps; filters by `created_at`)

Response (`ListResponse`):

- `items`: array of `LogItem` formatted for the UI
- `facets`: { `levels`, `statuses`, `durations` } (UI-friendly values)
- `counts`: { `info`, `warning`, `error` } counts by level (computed against the same filters, **except** the level filter is ignored so the three buckets are always present)
- `pagination`: { `page`, `per_page`, `total`, `total_pages`, `has_next`, `has_prev` }
- `applied`: echoes the applied filter values

### `GET /logs/chart?range=24h|7d`

Returns server-side aggregated buckets for the logs trend chart.

- `range=24h`: returns **24 points** bucketed by hour
- `range=7d`: returns **7 points** bucketed by day

Response:

- array of `{ date: RFC3339 string, count: number }`

Design notes:

- This endpoint is intended to be **independent of UI filters** (it shows the global trend for the selected range).
- Bucketing is implemented in the repository with DB-specific queries (Postgres/SQLite) and a safe fallback.

### `GET /logs/export`

Exports logs as an `.xlsx` file.

Query params match the filter set of `GET /logs` (`q`, `level`, `status`, `duration`, `from`, `to`).

Response:

- binary Excel file (sets `Content-Disposition` and the appropriate content type)

## How to add new audit logs

Preferred pattern:

- Publish a `types.AuditEvent` from the service layer where the business action is decided.
- Include:
  - a stable `entity` + `action`
  - a safe `payload` (never include plaintext secrets)
  - `actor_id` when known (user id as string), otherwise `"system"` / `"unknown"`

## Troubleshooting

### Logs not appearing in DB

Check:

- the event router is running (Watermill router is started in `main.go`)
- events are published to `system.audit_logs`
- `LogsDB` is connected and migrated

### Missing request ids / duration

Check middleware order in `backend/routes/routes.go`:

- `RequestIDMiddleware` must run before service code that reads the request id from context
- `RequestTimingMiddleware` must run before service code that reads `startTime`

If you publish audit events from middleware that runs **before** those two middlewares, `request_id` / `duration_ms` may be empty/0 for blocked requests.

