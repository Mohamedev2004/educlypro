package notifications

import (
	"educlypro/modules/auth"

	"gorm.io/gorm"
)

func SeedNotifications(db *gorm.DB, count int) error {
	var users []auth.User
	if err := db.Select("id").Find(&users).Error; err != nil {
		return err
	}

	if len(users) == 0 {
		return nil
	}

	userIDsNeedingSeed := make([]uint, 0, len(users))
	for _, user := range users {
		var existingForUser int64
		if err := db.Model(&Notification{}).
			Where("user_id = ?", user.ID).
			Count(&existingForUser).Error; err != nil {
			return err
		}
		if existingForUser == 0 {
			userIDsNeedingSeed = append(userIDsNeedingSeed, user.ID)
		}
	}

	if len(userIDsNeedingSeed) == 0 {
		return nil
	}

	fakeNotifications := NewFakeNotifications(userIDsNeedingSeed, count)
	if len(fakeNotifications) == 0 {
		return nil
	}

	return db.Create(&fakeNotifications).Error
}
