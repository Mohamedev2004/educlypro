package notifications

import (
	"gorm.io/gorm"
)

type seedCount struct {
	Unread int
	Read   int
}

// notificationSeedCountsByRole mirrors the role restriction enforced on the
// /notifications route group (see routes.go) — only these roles have an
// inbox, so only these roles get seeded demo data, each with its own
// unread/read split.
var notificationSeedCountsByRole = map[string]seedCount{
	"super_admin":  {Unread: 30, Read: 20},
	"center_owner": {Unread: 15, Read: 15},
}

// SeedNotifications seeds demo data for the notifications inbox. Each user
// in a notificationSeedCountsByRole role, without existing notifications,
// gets exactly that role's unread + read count of fake notifications.
func SeedNotifications(db *gorm.DB) error {
	roles := make([]string, 0, len(notificationSeedCountsByRole))
	for role := range notificationSeedCountsByRole {
		roles = append(roles, role)
	}

	var rows []struct {
		ID   uint
		Role string
	}
	if err := db.
		Table("users").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("roles.name IN ?", roles).
		Select("users.id AS id, roles.name AS role").
		Find(&rows).Error; err != nil {
		return err
	}

	for _, row := range rows {
		var existing int64
		if err := db.Model(&Notification{}).
			Where("user_id = ?", row.ID).
			Count(&existing).Error; err != nil {
			return err
		}
		if existing > 0 {
			continue
		}

		counts := notificationSeedCountsByRole[row.Role]
		fakeNotifications := NewFakeNotificationsForUser(row.ID, counts.Unread, counts.Read)
		if len(fakeNotifications) == 0 {
			continue
		}
		if err := db.Create(&fakeNotifications).Error; err != nil {
			return err
		}
	}

	return nil
}
