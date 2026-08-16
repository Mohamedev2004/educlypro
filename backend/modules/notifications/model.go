package notifications

import (
	"time"

	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Notification struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	RequestID string         `gorm:"index" json:"request_id"`
	UserID    uint           `gorm:"index" json:"user_id"`
	Topic     string         `json:"topic"`
	Title     string         `json:"title"`
	Body      string         `json:"body"`
	Payload   datatypes.JSON `json:"payload"`
	Channel   string         `json:"channel"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	ReadAt    *time.Time     `json:"read_at"`
	CreatedAt time.Time      `json:"created_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
