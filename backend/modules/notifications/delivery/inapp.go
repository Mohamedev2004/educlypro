package delivery

import (
	"educlypro/shared/types"
	"context"
	"encoding/json"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type inAppDispatcher struct {
	db *gorm.DB
}

func NewInAppDispatcher(db *gorm.DB) Dispatcher {
	return &inAppDispatcher{db: db}
}

func (d *inAppDispatcher) Channel() types.NotificationChannel {
	return types.ChannelInApp
}

func (d *inAppDispatcher) Send(ctx context.Context, userID uint, email string, event *types.NotificationEvent) error {
	payloadBytes, _ := json.Marshal(event.Payload)

	// Persisting directly via gorm.DB to avoid circular imports
	// using Table("notifications") as instructed
	return d.db.WithContext(ctx).Table("notifications").Create(map[string]any{
		"request_id": event.RequestID,
		"user_id":    userID,
		"topic":      event.Topic,
		"title":      event.Title,
		"body":       event.Body,
		"payload":    datatypes.JSON(payloadBytes),
		"channel":    string(types.ChannelInApp),
		"is_read":    false,
		"created_at": event.Timestamp,
	}).Error
}
