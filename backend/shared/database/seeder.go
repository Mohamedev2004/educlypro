package database

import (
	"educlypro/modules/academic"
	"educlypro/modules/auth"
	"educlypro/modules/notifications"
	"educlypro/modules/teachers"

	"gorm.io/gorm"
)

func SeedAll(db *gorm.DB) {
	// Seed Auth Users (+ the one center + its sub-centers)
	center, err := auth.SeedUsers(db)
	if err != nil {
		panic(err)
	}

	// Seed the Moroccan academic structure (grades -> majors -> subjects,
	// plus one auto-created class per major) for that center.
	seededIDs, err := academic.SeedAcademicStructure(db, center.ID)
	if err != nil {
		panic(err)
	}

	// Seed a couple of demo teachers, assigned to real seeded subjects and
	// classes.
	if err := teachers.SeedTeachers(db, center.ID, seededIDs); err != nil {
		panic(err)
	}

	// Seed Notifications (super_admin only)
	if err := notifications.SeedNotifications(db); err != nil {
		panic(err)
	}
}
