# Centers Documentation

## Overview

The centers module gives a **super admin** full visibility and light
provisioning control over every center in the system: a paginated list, a
detail view with the full staff roster, and the ability to create centers,
assign an owner, and add staff — all scoped to an arbitrary center picked by
**slug**, which is what makes this different from the `staff` module (which
only lets a `center_owner` manage *their own* center).

Every endpoint that addresses one specific center takes its **slug** in the
URL, never the numeric id — mirroring the "use the slug for route-model
binding so URLs never expose the id" convention from the campus-connect
reference project's `University` model. The numeric id still exists (it's
the DB primary key and is returned in response bodies for the frontend's
internal bookkeeping), it's just never something a client passes in.

All endpoints live under `/api/v1/centers/*` and require:

- `AuthMiddleware` (valid session)
- `RequireRole("super_admin")` (see `backend/shared/middleware/role.go`)

## Data model

This module does not define its own table — `Center` and `User` are the
same `auth.Center`/`auth.User` models used everywhere else
(`backend/modules/auth/model.go`). It only adds read/write paths over them.

- **Owner**: `auth.Center.Owner`, preloaded from `Center.OwnerID`. A center
  may have no owner (`owner`/`owner_username`/`owner_email` are `null` in
  that case) until `POST /centers/:slug/owner` is called.
- **Staff**: `auth.User` rows at that center whose role is `center_scanner`
  or `center_receptionist` — reusing the role name constants from the
  `staff` module (`staff.RoleScanner`, `staff.RoleReceptionist`) rather than
  duplicating them, so this stays in sync if those role names ever change.
  The center owner is **not** counted as staff — they're surfaced
  separately as `owner`.

## API endpoints

### `GET /centers`

Lists all centers, paginated, searchable, and sortable.

Query params:

- `page` (default `1`, must be >= 1)
- `per_page` (default `10`, allowed: `10|20|30|40|50`)
- `q` (optional, matches center name/slug, case-insensitive)
- `sort` (default `created_at`, allowed: `name|slug|created_at`)
- `direction` (default `desc`, allowed: `asc|desc`)

Response (`ListResponse`):

- `items`: array of `Response` — `{ id, name, slug, owner_username, owner_email, staff_count, created_at }`
- `pagination`: `{ page, per_page, total, total_pages, has_next, has_prev }`

Errors:

- `invalid_centers_page` / `invalid_centers_per_page` / `invalid_centers_sort` / `invalid_centers_direction` (400)
- `centers_list_failed` (500)

### `POST /centers`

Creates a new center. The slug is **never client-supplied** — it's always
derived from `name` server-side (see "Slug generation" below), and the UI
never shows a slug field.

Request body (`CreateRequest`):

- `name` (required, 2–150)

Response: `Response` (same shape as the list — `staff_count: 0`, no owner,
since a center starts empty).

Errors:

- `validation_failed` (400)
- `centers_create_failed` (500)

#### Slug generation

`Create` derives the slug from `name` and guarantees uniqueness under all
inputs, including concurrent requests:

- **Diacritics are stripped** before slugifying (NFD-decompose, drop
  Unicode combining marks, NFC-recompose), so `"Café de Paris"` becomes
  `cafe-de-paris`, not a lossy `caf-de-paris`.
- Remaining non-ASCII-alphanumeric runs (any non-Latin script — Arabic,
  CJK — or symbols/emoji) collapse to a single `-`, trimmed at both ends.
- **Names with no representable ASCII content at all** (pure Arabic,
  pure emoji, symbols-only) would otherwise slugify to an empty string;
  these fall back to the base `"center"` instead.
- The slug is capped at 140 characters, leaving room within the DB's
  150-char column for a `-NNN` uniqueness suffix.
- **Uniqueness is enforced by insert-and-retry, not check-then-insert**:
  `Create` attempts a real `INSERT` with the base slug; if the DB's unique
  index rejects it (surfaced as `gorm.ErrDuplicatedKey`, via
  `gorm.Config{TranslateError: true}` in `shared/database/db.go`), it
  retries with `-2`, `-3`, ... up to 100 attempts before giving up with
  `ErrSlugGenerationFailed`. Because each attempt is an atomic DB
  operation, two concurrent requests for the same name can never both
  succeed with the same slug — there is no separate check-then-insert race
  window to lose.

### `GET /centers/:slug`

Returns the full detail view for one center: identity, owner (if any), and
**every** staff member (not just a count — this is what the list endpoint
doesn't give you).

Response (`DetailResponse`):

- `id`, `name`, `slug`, `created_at`
- `owner`: `{ id, username, email, created_at }` or `null`
- `staff`: array of `{ id, username, email, role, created_at }`

Errors:

- `center_not_found` (404)
- `centers_get_failed` (500)

### `POST /centers/:slug/owner`

Creates the `center_owner` account for a center that doesn't have one yet.
Fails if the center already has an owner — a center can have at most one
owner, and a user can own at most one center (`Center.OwnerID` is a unique
index — see `backend/modules/auth/model.go`).

Request body (`CreateOwnerRequest`):

- `username` (required, 3–100)
- `email` (required, email, max 150)
- `password` (required, 6–72)

Behavior: creates the user and sets `Center.OwnerID` atomically in a single
DB transaction (`Repository.AssignOwner`) — the `UPDATE` is conditioned on
`owner_id IS NULL`, so a race between two concurrent "add owner" requests
can't both succeed.

Response: `OwnerResponse` — `{ id, username, email, created_at }`.

Errors:

- `validation_failed` (400)
- `center_not_found` (404)
- `owner_already_assigned` (409)
- `email_taken` (409)
- `centers_add_owner_failed` (500)

### `POST /centers/:slug/staff`

Creates a `center_scanner` or `center_receptionist` account for the given
center. Unlike `POST /staff` (see `documentation/modules/staff.md`), the
center comes from the URL, not from the caller's own ownership — this
endpoint doesn't delegate to a fresh implementation, it calls
`staff.Service.CreateForCenter`, the same creation logic `POST /staff` uses
internally, just entered from a different, admin-authorized path. See
"Relationship to other modules" below.

Request body (`AddStaffRequest`):

- `username` (required, 3–100)
- `email` (required, email, max 150)
- `password` (required, 6–72)
- `role` (required, `center_scanner|center_receptionist`)

Response: `staff.Response` — same shape `POST /staff` returns.

Errors:

- `validation_failed` (400)
- `center_not_found` (404)
- `email_taken` (409)
- `centers_add_staff_failed` (500)

## Relationship to other modules

- Reuses `auth.Center` / `auth.User` — no new tables, no migrations needed
  beyond what `auth` already owns.
- Reuses the `center_scanner`/`center_receptionist` role constants from the
  `staff` module — see `documentation/modules/staff.md`.
- **`staff.Service.CreateForCenter`**: the `staff` module's `Create` method
  resolves the target center from the caller's own ownership
  (`FindOwnerCenterID`), which doesn't make sense for a super admin acting
  on an arbitrary center. Rather than duplicate the validation/creation
  logic, `staff.Service` exposes a second entrypoint,
  `CreateForCenter(ctx, centerID, actorUserID, req)`, that skips the
  ownership resolution and trusts the given `centerID` directly. Both
  entrypoints funnel into the same private `createStaffMember` helper, so
  email-taken checks, role lookup, password hashing, and audit-event
  publishing only exist in one place.
