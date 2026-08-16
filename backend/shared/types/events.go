package types

import "time"

type EventAction string
type LogLevel string // NEW: Define the LogLevel type

const (
	ActionCreated EventAction = "CREATED"
	ActionUpdated EventAction = "UPDATED"
	ActionDeleted EventAction = "DELETED"
	ActionFailed  EventAction = "FAILED"
)

// NEW: Define your standard levels
const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

type AuditEvent struct {
	RequestID  string      `json:"request_id"`
	Level      LogLevel    `json:"level"`
	Entity     string      `json:"entity"`
	EntityID   string      `json:"entity_id"`
	ActorID    string      `json:"actor_id"`
	StatusCode int         `json:"status_code"`
	Action     EventAction `json:"action"`
	Payload    any         `json:"payload"`
	DurationMs float64     `json:"duration_ms"`
	Timestamp  time.Time   `json:"timestamp"`
}
