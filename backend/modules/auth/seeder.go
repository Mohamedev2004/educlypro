package auth

import (
	"educlypro/shared/utils"
	"errors"

	"gorm.io/gorm"
)

// SeedUsers provisions the fixed set of roles, one center, and one user per
// role (super_admin, center_owner, center_scanner, center_receptionist).
// It is idempotent: re-running it will not create duplicates.
func SeedUsers(db *gorm.DB) error {
	roleNames := []string{"super_admin", "center_owner", "center_scanner", "center_receptionist"}
	rolesByName := make(map[string]Role, len(roleNames))
	for _, name := range roleNames {
		role := Role{Name: name}
		if err := db.Where("name = ?", name).FirstOrCreate(&role).Error; err != nil {
			return err
		}
		rolesByName[name] = role
	}

	ensureUser := func(username, email, plainPassword string, roleID uint, centerID *uint) (User, error) {
		var user User
		if err := db.Where("email = ?", email).First(&user).Error; err == nil {
			return user, nil
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			return User{}, err
		}

		hashed, err := utils.HashPassword(plainPassword)
		if err != nil {
			return User{}, err
		}

		user = User{
			Username: username,
			Email:    email,
			Password: hashed,
			RoleID:   roleID,
			CenterID: centerID,
		}
		if err := db.Create(&user).Error; err != nil {
			return User{}, err
		}
		return user, nil
	}

	// The one center.
	var center Center
	if err := db.Where("slug = ?", "educlypro-center").FirstOrCreate(&center, NewCenter("EduclyPro Center", "educlypro-center")).Error; err != nil {
		return err
	}

	// Fixed super admin (not tied to a center).
	if _, err := ensureUser("super-admin", "super-admin@educlypro.com", "password", rolesByName["super_admin"].ID, nil); err != nil {
		return err
	}

	// The center's single owner.
	owner, err := ensureUser("center-owner", "owner@educlypro.com", "password", rolesByName["center_owner"].ID, &center.ID)
	if err != nil {
		return err
	}

	// The center's scanner and receptionist.
	if _, err := ensureUser("center-scanner", "scanner@educlypro.com", "password", rolesByName["center_scanner"].ID, &center.ID); err != nil {
		return err
	}
	if _, err := ensureUser("center-receptionist", "receptionist@educlypro.com", "password", rolesByName["center_receptionist"].ID, &center.ID); err != nil {
		return err
	}

	// Link the center back to its owner.
	if err := db.Model(&center).Update("owner_id", owner.ID).Error; err != nil {
		return err
	}

	return nil
}
