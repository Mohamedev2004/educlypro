# Staff Management Documentation

## Overview

The staff module lets a **center owner** manage the two operational roles at
their own center:

- `center_scanner`
- `center_receptionist`

Staff members are regular rows in the `users` table (same `auth.User` model
used everywhere else) — this module does not introduce a separate table. It
only constrains *who* can create/update/delete those users and *which*
center they belong to.

All endpoints live under `/api/v1/staff/*` and require:

- `AuthMiddleware` (valid session)
- `RequireRole("center_owner")` (see `backend/shared/middleware/role.go`)

So only a signed-in `center_owner` can reach this API, and every operation
is scoped to **that owner's own center** — there is no cross-center staff
management and no center picker in the request; `center_id` is always
validated server-side against the owner, never trusted from the client.

## Center scoping

Every service method starts by resolving the owner's center:

- `Repository.FindOwnerCenterID(ctx, ownerUserID)` looks up the caller's
  `center_id` from their own `users` row.

`Create` and `Update` additionally require the request's `center_id` to
match that resolved center; a mismatch returns `ErrCenterMismatch` (mapped
to `403 center_mismatch`) and publishes a `system.events.v1.staff.center_mismatch`
audit event — this is treated as a security-relevant event, not just a
validation error.

## API endpoints

All endpoints are under `/api/v1/staff`.

### `GET /staff`

Lists staff for the caller's center, paginated and filterable.

Query params:

- `page` (default `1`, must be >= 1)
- `per_page` (default `10`, allowed: `10|20|30|40|50`)
- `role` (optional, `center_scanner|center_receptionist`)
- `q` (optional, matches username/email, case-insensitive)
- `sort` (default `created_at`, allowed: `username|email|role|created_at`)
- `direction` (default `desc`, allowed: `asc|desc`)

Response (`ListResponse`):

- `items`: array of `Response` (see Data model)
- `pagination`: `{ page, per_page, total, total_pages, has_next, has_prev }`

### `POST /staff`

Creates a new `center_scanner` or `center_receptionist` user in the owner's
center.

Request body (`CreateRequest`):

- `username` (required, 3–100)
- `email` (required, email, max 150)
- `password` (required, 6–72)
- `role` (required, `center_scanner|center_receptionist`)
- `center_id` (required, must equal the owner's own center)

Errors:

- `center_mismatch` (403) — `center_id` doesn't match the owner's center
- `email_taken` (409) — email already in use
- `staff_create_failed` (500)

`Service.Create` resolves the center from the caller (`FindOwnerCenterID`)
then delegates to a private `createStaffMember` helper. A super admin
managing an arbitrary center's staff — not their own — goes through a
second entrypoint on the same service, `CreateForCenter(ctx, centerID,
actorUserID, req)`, which skips the ownership resolution and trusts the
given `centerID` directly, but funnels into that same `createStaffMember`
helper so the actual validation/creation logic only exists once. This is
what `POST /centers/:id/staff` calls — see `documentation/modules/centers.md`.

### `PUT /staff/:id`

Updates an existing staff member's profile, role, and/or center-scoped
password. Password is **optional** — omit it to leave the current password
unchanged.

Request body (`UpdateRequest`):

- `username` (required, 3–100)
- `email` (required, email, max 150)
- `role` (required, `center_scanner|center_receptionist`)
- `password` (optional, 6–72 when present)
- `center_id` (required, must equal the owner's own center)

Errors:

- `center_mismatch` (403)
- `staff_not_found` (404) — no such staff id in the owner's center
- `email_taken` (409)
- `staff_update_failed` (500)

### `DELETE /staff/:id`

Soft-deletes a staff member (GORM `DeletedAt`) and revokes their active
sessions.

Behavior:

- Verifies the staff id belongs to the owner's center before deleting.
- After soft-delete, best-effort deletes the user's `access` and `refresh`
  token rows (see `documentation/modules/auth.md` — token revocation
  allow-list) so the removed staff member is immediately signed out
  everywhere. A failure to revoke tokens is logged but does not fail the
  request.

Errors:

- `staff_not_found` (404)
- `staff_delete_failed` (500)

## Data model (`staff.Response`)

Shape returned by all endpoints:

- `id`
- `username`
- `email`
- `role`
- `center_id`
- `created_at` (RFC3339)

## Audit events

Like `auth` and `notifications`, this module publishes events module-locally
(no shared publish helper) using the same `types.AuditEvent` struct, always
mirrored to the global sink topic `system.audit_logs`:

- `system.events.v1.staff.created`
- `system.events.v1.staff.updated`
- `system.events.v1.staff.deleted`
- `system.events.v1.staff.center_mismatch` (security-relevant: attempted
  cross-center access)
- `system.events.v1.staff.system_error` (unexpected repository/service
  errors)

`Entity` is always `"Staff"`; `EntityID` is the affected staff user id (or
`"unknown"` for pre-lookup failures); `ActorID` is the owner's user id.

See `documentation/modules/logs.md` for how these events become rows in the
Logs UI.

## Relationship to auth

- Staff accounts are ordinary `auth.User` rows — they log in through the
  normal `POST /auth/login` flow described in `documentation/modules/auth.md`.
- Deleting a staff member reuses the auth module's token store (via the
  `TokenRevoker` interface implemented by `auth.Repository`) to invalidate
  sessions immediately, rather than waiting for token expiry.
