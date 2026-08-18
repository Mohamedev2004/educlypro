package auth

import (
	"educlypro/shared/utils"
	"errors"

	"gorm.io/gorm"
)

// SeedUsers provisions the fixed set of roles, one center with two
// sub-centers, and one user per role — including a scanner and a
// receptionist for each sub-center. It is idempotent: re-running it will
// not create duplicates. Returns the seeded center so callers (SeedAll) can
// chain further center-scoped seeding without a second lookup.
func SeedUsers(db *gorm.DB) (Center, error) {
	roleNames := []string{"super_admin", "center_owner", "center_scanner", "center_receptionist"}
	rolesByName := make(map[string]Role, len(roleNames))
	for _, name := range roleNames {
		role := Role{Name: name}
		if err := db.Where("name = ?", name).FirstOrCreate(&role).Error; err != nil {
			return Center{}, err
		}
		rolesByName[name] = role
	}

	ensureUser := func(username, email, plainPassword string, roleID uint, centerID, subCenterID *uint) (User, error) {
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
			Username:    username,
			Email:       email,
			Password:    hashed,
			RoleID:      roleID,
			CenterID:    centerID,
			SubCenterID: subCenterID,
		}
		if err := db.Create(&user).Error; err != nil {
			return User{}, err
		}
		return user, nil
	}

	ensureSubCenter := func(centerID uint, name string) (SubCenter, error) {
		var subCenter SubCenter
		err := db.Where("center_id = ? AND name = ?", centerID, name).
			FirstOrCreate(&subCenter, SubCenter{CenterID: centerID, Name: name}).Error
		return subCenter, err
	}

	// The one center.
	var center Center
	if err := db.Where("slug = ?", "educlypro-center").FirstOrCreate(&center, NewCenter("EduclyPro Center", "educlypro-center")).Error; err != nil {
		return Center{}, err
	}

	// Fixed super admin (not tied to a center).
	if _, err := ensureUser("super-admin", "super-admin@educlypro.com", "password", rolesByName["super_admin"].ID, nil, nil); err != nil {
		return Center{}, err
	}

	// The center's single owner.
	owner, err := ensureUser("center-owner", "owner@educlypro.com", "password", rolesByName["center_owner"].ID, &center.ID, nil)
	if err != nil {
		return Center{}, err
	}

	// Two sub-centers, each staffed with its own scanner and receptionist —
	// every center_scanner/center_receptionist must belong to a sub-center,
	// so seeding "bare" staff with no sub-center would leave a state the
	// app itself can never produce anymore.
	branches := []struct {
		name                 string
		scannerUsername      string
		scannerEmail         string
		receptionistUsername string
		receptionistEmail    string
	}{
		{
			name:                 "Main Branch",
			scannerUsername:      "main-scanner",
			scannerEmail:         "main-scanner@educlypro.com",
			receptionistUsername: "main-receptionist",
			receptionistEmail:    "main-receptionist@educlypro.com",
		},
		{
			name:                 "North Branch",
			scannerUsername:      "north-scanner",
			scannerEmail:         "north-scanner@educlypro.com",
			receptionistUsername: "north-receptionist",
			receptionistEmail:    "north-receptionist@educlypro.com",
		},
	}

	for _, branch := range branches {
		subCenter, err := ensureSubCenter(center.ID, branch.name)
		if err != nil {
			return Center{}, err
		}

		if _, err := ensureUser(branch.scannerUsername, branch.scannerEmail, "password", rolesByName["center_scanner"].ID, &center.ID, &subCenter.ID); err != nil {
			return Center{}, err
		}
		if _, err := ensureUser(branch.receptionistUsername, branch.receptionistEmail, "password", rolesByName["center_receptionist"].ID, &center.ID, &subCenter.ID); err != nil {
			return Center{}, err
		}
	}

	// Link the center back to its owner.
	if err := db.Model(&center).Update("owner_id", owner.ID).Error; err != nil {
		return Center{}, err
	}
	center.OwnerID = &owner.ID

	return center, nil
}
