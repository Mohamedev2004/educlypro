# Notifications Documentation

## Overview

Notifications are implemented as an **event-driven dispatch system** plus an **in-app inbox API**.

- **Dispatching**: services publish `types.NotificationEvent` messages to Watermill.
- **Delivery**: a subscriber receives events and fans out to one or more delivery channels.
- **Inbox**: the backend exposes `/api/v1/notifications/*` endpoints, restricted to `super_admin` and `center_owner`, for listing and managing in-app notifications persisted in the main DB. No other role has a notifications inbox.

Runtime wiring is split into:

- **Event subscriber** (Watermill): in `backend/cmd/main.go`
- **HTTP API wiring** (Gin): in `backend/routes/routes.go`

## Concepts

### Channels

Defined channels:

- `in_app`
- `email`

Delivery channels are selected per event via `event.channels`.

### Recipients

Recipients are resolved using a union of:

- **Direct recipients**: `recipient_ids` (explicit user ids)
- **Role targets**: `role_targets` (resolved into user ids by querying users with matching roles)

Duplicates are removed during fan-out.

## Event dispatch architecture

### Publishing

The canonical publish topic is:

- `types.NotificationTopic` = **`notifications.dispatch`**

Publishers should set:

- `topic` (business topic like `user.registered`, `user.welcome`)
- `title`, `body`
- `payload` (JSON object)
- `channels` (one or more of `in_app`, `email`)
- either `recipient_ids` and/or `role_targets`
- `timestamp`

The publisher may omit `request_id`; the auth service injects `request_id` from context at publish time.

### Subscriber + service dispatch

At runtime (`backend/cmd/main.go`):

- a Watermill handler is registered on `notifications.dispatch`
- `notifications.Subscriber.ProcessNotificationEvent` unmarshals the event and calls `notifications.Service.Dispatch`

Dispatch behavior:

1. Build a recipient set from `recipient_ids`
2. If `role_targets` is provided, resolve all users with those roles and add them to the set
3. For each recipient and each channel:
   - resolve email address (needed for email channel)
   - call the matching dispatcher implementation

Failures:

- Malformed events are dropped (no retry).
- Dispatch errors return an error so Watermill retries the message.
- Per-recipient/per-channel send failures are logged and do **not** stop other sends.

## Delivery implementations

### In-app (`in_app`)

`delivery.inAppDispatcher` persists a row into the `notifications` table with:

- `request_id`
- `user_id` (the recipient)
- `topic`, `title`, `body`
- `payload` (JSON)
- `channel = "in_app"`
- `is_read = false`
- `created_at = event.timestamp`

### Email (`email`)

`delivery.emailDispatcher` sends email using `shared/utils/email.go`.

Special-cased topics:

- `user.welcome`: sends a welcome email
- `user.registered`: sends an admin “new registration” email

All other topics fall back to a generic email send using `title` and `body`.

Templates used by the mailer live under `backend/templates/emails/*`.

## In-app inbox API

All endpoints are under `/api/v1` and require `AuthMiddleware` + `RequireRole("super_admin", "center_owner")` — the whole `/notifications` route group is gated at once (`backend/modules/notifications/routes.go`), not per-route. Each user only ever sees their own notifications regardless of role (every handler scopes by the caller's own `userID`) — the role gate controls who has an inbox at all, not which rows within it.

HTTP wiring lives in `backend/routes/routes.go`, which builds:

- `notifications.Repository` (MainDB)
- `auth.UserResolver` (MainDB)
- a dispatcher list: `delivery.NewInAppDispatcher(MainDB)` and `delivery.NewEmailDispatcher()`
- `notifications.Service` and `notifications.Handler`

Then registers routes under `/api/v1/notifications`.

### `GET /notifications`

Query params:

- `page` (default `1`, must be >= 1)
- `per_page` (default `10`, allowed: `10|20|30|40|50`)
- `filter` (default `unread`, allowed: `read|unread`)

Response shape (`ListResponse`):

- `items`: array of notifications
- `filter`: the filter that was applied
- `counts`: `{ all, read, unread }`
- `pagination`: `{ page, per_page, total, total_pages, has_next, has_prev }`

### `GET /notifications/unread-count`

Returns:

- `{ count: number }`

### `PATCH /notifications/:id/read`

Marks a single notification as read (sets `is_read=true` and `read_at=now`).

### `PATCH /notifications/read-all`

Marks all unread notifications as read for the authenticated user.

### `DELETE /notifications/:id`

Deletes a notification owned by the authenticated user (soft delete via GORM `DeletedAt`).

## Data model (`notifications.Notification`)

Persisted columns include:

- `request_id` (indexed)
- `user_id` (indexed)
- `topic`, `title`, `body`
- `payload` (JSON)
- `channel`
- `is_read`, `read_at`
- `created_at`
- `deleted_at` (soft delete)

## Seeding

`notifications.SeedNotifications` (`backend/modules/notifications/seeder.go`)
seeds demo inbox data for **`super_admin` and `center_owner` users only**
(`notificationSeedCountsByRole`), matching the inbox's role restriction
above. For each qualifying user that doesn't already have notifications, it
creates a role-specific, deterministic split of fake notifications via
`NewFakeNotificationsForUser` (not a random per-notification coin flip):

- `super_admin`: 30 unread + 20 read (50 total)
- `center_owner`: 15 unread + 15 read (30 total)

Called from `database.SeedAll`.

## Current publishers

No module currently publishes to `notifications.dispatch` — the dispatch
architecture above (event → subscriber → in-app/email delivery) is wired up
and ready, but nothing in `auth`, `staff`, or `centers` calls it today. The
in-app inbox is populated by the seeder for demo purposes instead. If you
wire up a real publisher, `role_targets` should use this app's actual role
names (`super_admin`, `center_owner`, `center_scanner`,
`center_receptionist`), and — since only `super_admin`/`center_owner` have
an inbox — `role_targets` values of `center_scanner`/`center_receptionist`,
or `recipient_ids` for users in those roles, would dispatch successfully but
never be visible to anyone through the HTTP API.

## Watermill wiring (quick reference)

At runtime (`backend/cmd/main.go`), the subscriber is registered like:

- topic: `types.NotificationTopic` (`notifications.dispatch`)
- handler: `notifications.Subscriber.ProcessNotificationEvent`

This is separate from the HTTP inbox API: the inbox reads/writes the `notifications` table directly (via `notifications.Repository`) and does not require Watermill to be running.
