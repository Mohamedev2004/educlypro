package schedules

import (
	"gorm.io/gorm"
)

// StartAll initializes all background tasks for the application
func StartAll(db *gorm.DB) {
	StartTokenCleanup(db)
	StartNotificationCleanup(db)
}
