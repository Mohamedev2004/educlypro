# Auth System Documentation

## Overview

The auth system is **cookie-first JWT auth** with **refresh rotation** and **DB-backed token revocation**.

- **Backend**: Gin + GORM. Endpoints live under `/api/v1/auth/*`.
- **Frontend**: Axios uses `withCredentials: true` and automatically calls `/auth/refresh` on most `401`s.
- **Security model**: JWT signature validation **and** a DB lookup of the hashed token to support logout/revocation.

## Session cookies

On register/login/refresh, the backend sets:

- **`auth_token`**: access JWT
- **`refresh_token`**: refresh JWT
- **`session_exists`**: `"true"` (non-HttpOnly, used by the frontend to decide whether to call `/auth/me`)

Cookie behavior is controlled by:

- `COOKIE_SECURE`, `COOKIE_HTTP_ONLY`, `COOKIE_DOMAIN`, `COOKIE_SAMESITE`

## Tokens & claims

### JWT contents

Both access and refresh tokens include:

- `user_id`
- `role`
- `token_type` (`"access"` or `"refresh"`)
- standard registered claims (iat/exp/jti)

Signing and expiry config:

- `JWT_SECRET`
- `JWT_ACCESS_EXPIRY_MINUTES` (default 15)
- `JWT_REFRESH_EXPIRY_DAYS` (default 7)

### Token storage (revocation allow-list)

The DB stores **SHA-256 hashes** of token strings in the `tokens` table:

- `type`: `access | refresh | reset_password`
- `token`: SHA-256 hash
- `expires_at`: enforced/queried for validity + background cleanup

## Protected routes: `AuthMiddleware`

Token extraction order:

1. Cookie `auth_token`
2. Header `Authorization: Bearer <token>`

Validation rules:

- JWT must be valid and not expired
- JWT `token_type` must be `"access"`
- Token must exist in DB (`tokens` row with matching hashed token, `type="access"`, and `expires_at > now()`)

Context values set:

- `userID` (uint)
- `role` (string)

## Frontend behavior

### When `/auth/me` runs

The frontend only enables the `me` query when `document.cookie` contains `session_exists=true`.

### Auto refresh on 401

- Most `401` responses trigger a single retry:
  - call `POST /auth/refresh`
  - then retry the original request
- Refresh calls are deduplicated through a shared in-flight promise.
- If refresh fails:
  - the frontend clears `session_exists`
  - dispatches `auth:logout`
  - auth context clears cached state and navigates to `/login`

## API endpoints

All endpoints below are under `/api/v1`.

### `POST /auth/register`

Creates a new user and issues tokens.

Request body:

- `username` (required, 3–100)
- `email` (required, email, max 150)
- `password` (required, 6–72)
- `role` (required, `admin|user`)

Behavior:

- Creates user (password hashed)
- Issues access + refresh JWTs
- Stores **hashed** token rows in DB
- Sets cookies: `auth_token`, `refresh_token`, `session_exists`
- Emits audit events and notification events (admins + welcome)

Errors:

- `email_taken` (409) when email already exists
- `invalid_role` (400) when role is not allowed

### `POST /auth/login`

Authenticates a user and rotates tokens.

Request body:

- `email` (required, email, max 150)
- `password` (required, 6–72)

Behavior:

- Validates credentials
- Deletes existing `access` and `refresh` token rows for the user
- Issues new tokens, stores hashed rows, sets cookies

Errors:

- `invalid_credentials` (401)

### `POST /auth/refresh`

Refreshes the session using the `refresh_token` cookie.

Rules:

- Requires cookie `refresh_token`
- Must be a valid JWT with `token_type == "refresh"`
- Must exist in DB as `type="refresh"` and not expired

Rotation behavior:

- Deletes the used refresh token row
- Deletes all access tokens for the user
- Issues a new access token and a new refresh token
- Resets cookies

Errors:

- `refresh_token_required` (401) when missing cookie (cookies are cleared)
- `invalid_refresh_token` (401) when invalid/revoked/expired (cookies are cleared)

### `POST /auth/logout` (protected)

Deletes all access + refresh token rows for the user and clears cookies.

### `GET /auth/me` (protected)

Returns the authenticated user profile.

### `PUT /auth/profile` (protected)

Updates `username` and `email`.

Errors:

- `email_taken` (409)

### `PUT /auth/password` (protected)

Updates the user password.

Request body:

- `currentPassword` (required)
- `newPassword` (required, 6–72)

Errors:

- `invalid_current_password` (400)

### `POST /auth/forgot-password`

Creates a **reset token** stored in DB as `type="reset_password"` (SHA-256 hashed), expiring in **30 minutes**.

Security behavior:

- If the email does not exist, the endpoint still returns success (silent fail).

### `POST /auth/reset-password`

Resets password using a reset token.

Rules:

- Token is SHA-256 hashed and looked up in DB
- Must be `type="reset_password"` and not expired
- On success:
  - updates password
  - deletes all reset_password tokens for the user

Errors:

- `invalid_or_expired_token` (400)

## Background maintenance

The backend runs a scheduled cleanup:

- deletes expired rows in `tokens` every **12 hours**

## Request tracing & events

- Frontend sends `X-Request-ID` on every request.
- Backend ensures a request id exists and injects it into the Go context for downstream services.
- Backend measures request start time and stores it in context; auth events include duration.

The auth service publishes:

- domain audit topics (e.g. `system.events.v1.auth.logged_in`)
- a global audit sink topic: `system.audit_logs`
- notification events to `notifications.dispatch`

See `documentation/modules/logs.md` and `documentation/modules/notifications.md`.
