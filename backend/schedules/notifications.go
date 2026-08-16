package schedules

import (
	"context"
	"log"
	"time"

	"educlypro/modules/notifications"

	"gorm.io/gorm"
)

func StartNotificationCleanup(db *gorm.DB) {
	Run("notification-cleanup", 24*time.Hour, func() error {
		ctx := context.Background()
		cutoff := time.Now().AddDate(0, -1, 0)

		result := db.WithContext(ctx).
			Where("is_read = ? AND created_at <= ? AND deleted_at IS NULL", true, cutoff).
			Delete(&notifications.Notification{})

		if result.Error != nil {
			return result.Error
		}

		log.Printf("Schedule [notification-cleanup] soft deleted %d read notifications older than %s", result.RowsAffected, cutoff.Format(time.DateOnly))
		return nil
	})
}
