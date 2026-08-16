package types

import "time"

const NotificationTopic = "notifications.dispatch"

type NotificationChannel string

const (
	ChannelInApp NotificationChannel = "in_app"
	ChannelEmail NotificationChannel = "email"
)

type NotificationEvent struct {
	RequestID    string                `json:"request_id"`
	Topic        string                `json:"topic"`
	RecipientIDs []uint                `json:"recipient_ids"`
	RoleTargets  []string              `json:"role_targets"`
	Title        string                `json:"title"`
	Body         string                `json:"body"`
	Payload      map[string]any        `json:"payload"`
	Channels     []NotificationChannel `json:"channels"`
	Priority     string                `json:"priority"`
	Timestamp    time.Time             `json:"timestamp"`
}
