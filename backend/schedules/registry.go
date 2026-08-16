package schedules

import (
	"educlypro/modules/auth"
	"educlypro/modules/notifications"
	"time"

	"gorm.io/gorm"
)

// StartAll initializes all background tasks for the application
func StartAll(db *gorm.DB) {
	authRepo := auth.NewRepository(db)
	notifRepo := notifications.NewRepository(db)

	// Register your individual tasks here
	Run("Auth Token Cleanup", 12*time.Hour, authRepo.DeleteExpiredTokens)
	Run("Notification Cleanup", 24*time.Hour, notifRepo.DeleteOldReadNotifications)

	// Future tasks:
	// Run("Daily Analytics", 12*time.Hour, analytics.GenerateReport)
}
