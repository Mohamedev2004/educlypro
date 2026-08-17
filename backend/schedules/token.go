package schedules

import (
	"log"
	"time"

	"educlypro/modules/auth"

	"gorm.io/gorm"
)

// StartTokenCleanup periodically deletes auth tokens whose expiry has passed.
func StartTokenCleanup(db *gorm.DB) {
	Run("Auth Token Cleanup", 12*time.Hour, func() error {
		result := db.Where("expires_at < ?", time.Now()).Delete(&auth.Token{})
		if result.Error != nil {
			return result.Error
		}

		log.Printf("Schedule [Auth Token Cleanup] deleted %d expired tokens", result.RowsAffected)
		return nil
	})
}
