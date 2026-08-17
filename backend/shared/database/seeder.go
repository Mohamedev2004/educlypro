package database

import (
	"educlypro/modules/auth"
	"educlypro/modules/notifications"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	// Seed Auth Users
	if err := auth.SeedUsers(db); err != nil {
		panic(err)
	}

	// Seed Notifications (super_admin only)
	if err := notifications.SeedNotifications(db); err != nil {
		panic(err)
	}
}
