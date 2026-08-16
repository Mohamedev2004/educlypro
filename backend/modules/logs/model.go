package logs

import (
	"time"

	"gorm.io/datatypes"
)

type AuditLog struct {
	ID         uint           `gorm:"primaryKey"`
	RequestID  string         `gorm:"index;type:varchar(255)"`
	Level      string         `gorm:"index;type:varchar(10)"`
	Entity     string         `gorm:"index;type:varchar(100)"`
	EntityID   string         `gorm:"index;type:varchar(255)"`
	ActorID    string         `gorm:"index;type:varchar(255)"`
	StatusCode int            `gorm:"index"`
	Action     string         `gorm:"type:varchar(100)"`
	Payload    datatypes.JSON `gorm:"type:jsonb"`
	DurationMs float64        `gorm:"index"`
	CreatedAt  time.Time      `gorm:"index"`
}
