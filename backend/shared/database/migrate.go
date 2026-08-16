package database

import (
	"educlypro/modules/auth"
	"educlypro/modules/logs" // Your new logs module
	"educlypro/modules/notifications"
	"log"
)

func Migrate() {
	// 1. Migrate core tables to the Main Database
	if err := MainDB.AutoMigrate(
		&auth.Role{},
		&auth.User{},
		&auth.Center{},
		&auth.Token{},
		&notifications.Notification{},
	); err != nil {
		log.Fatal("Main DB Migration failed:", err)
	}

	if MainDB.Migrator().HasTable("notification_preferences") {
		if err := MainDB.Migrator().DropTable("notification_preferences"); err != nil {
			log.Fatal("Failed to drop notification_preferences table:", err)
		}
	}

	// Drop the old many-to-many join tables — users now have a single
	// RoleID/CenterID instead of a user_roles/center_users membership table.
	if MainDB.Migrator().HasTable("user_roles") {
		if err := MainDB.Migrator().DropTable("user_roles"); err != nil {
			log.Fatal("Failed to drop user_roles table:", err)
		}
	}
	if MainDB.Migrator().HasTable("center_users") {
		if err := MainDB.Migrator().DropTable("center_users"); err != nil {
			log.Fatal("Failed to drop center_users table:", err)
		}
	}

	// Drop the old dual center_admin_id/responsible_id columns — replaced by
	// the single owner_id column.
	if MainDB.Migrator().HasColumn(&auth.Center{}, "center_admin_id") {
		if err := MainDB.Migrator().DropColumn(&auth.Center{}, "center_admin_id"); err != nil {
			log.Fatal("Failed to drop centers.center_admin_id column:", err)
		}
	}
	if MainDB.Migrator().HasColumn(&auth.Center{}, "responsible_id") {
		if err := MainDB.Migrator().DropColumn(&auth.Center{}, "responsible_id"); err != nil {
			log.Fatal("Failed to drop centers.responsible_id column:", err)
		}
	}

	// 2. Migrate the audit log table to the separate Logs Database
	if err := LogsDB.AutoMigrate(
		&logs.AuditLog{},
	); err != nil {
		log.Fatal("Logs DB Migration failed:", err)
	}
}
