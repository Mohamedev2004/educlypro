package teachers

import (
	"errors"
	"fmt"

	"educlypro/modules/academic"
	"educlypro/shared/utils"

	"gorm.io/gorm"
)

type teacherSeed struct {
	fullName string
	email    string
	phone    string
	// subjectKeys are "Major > Subject" keys matching
	// academic.SeedResult.SubjectIDs.
	subjectKeys []string
	// classKeys are major-name keys matching academic.SeedResult.ClassIDs —
	// a teacher is assigned to a major's class wherever they teach a
	// subject in it.
	classKeys []string
}

// demoTeachers assigns each teacher to the same subject across every major
// that offers it, and to that major's class.
var demoTeachers = []teacherSeed{
	{
		fullName: "Youssef El Amrani",
		email:    "youssef.elamrani@educlypro.com",
		phone:    "0600000001",
		subjectKeys: []string{
			"Sciences Expérimentales > Mathématiques",
			"Sciences Mathématiques A > Mathématiques",
			"Sciences Physiques > Mathématiques",
		},
		classKeys: []string{
			"Sciences Expérimentales",
			"Sciences Mathématiques A",
			"Sciences Physiques",
		},
	},
	{
		fullName: "Fatima Zahra Benali",
		email:    "fatima.benali@educlypro.com",
		phone:    "0600000002",
		subjectKeys: []string{
			"Lettres et Sciences Humaines > Arabe",
			"Sciences Mathématiques A > Philosophie",
			"Sciences Physiques > Philosophie",
		},
		classKeys: []string{
			"Lettres et Sciences Humaines",
			"Sciences Mathématiques A",
			"Sciences Physiques",
		},
	},
}

// SeedTeachers idempotently creates demoTeachers for the given center,
// assigned to subjects/classes resolved from the ids academic.SeedAcademicStructure
// returned.
func SeedTeachers(db *gorm.DB, centerID uint, ids academic.SeedResult) error {
	for _, seed := range demoTeachers {
		var teacher Teacher
		err := db.Where("center_id = ? AND email = ?", centerID, seed.email).First(&teacher).Error
		if err == nil {
			continue // already seeded
		}
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			return err
		}

		slug, err := uniqueTeacherSlug(db, centerID, seed.fullName)
		if err != nil {
			return err
		}

		teacher = Teacher{
			CenterID: centerID,
			Slug:     slug,
			FullName: seed.fullName,
			Email:    seed.email,
			Phone:    seed.phone,
		}
		if err := db.Create(&teacher).Error; err != nil {
			return err
		}

		subjectIDs := make([]uint, 0, len(seed.subjectKeys))
		for _, key := range seed.subjectKeys {
			if id, ok := ids.SubjectIDs[key]; ok {
				subjectIDs = append(subjectIDs, id)
			}
		}
		if len(subjectIDs) > 0 {
			var subjects []academic.Subject
			if err := db.Where("id IN ?", subjectIDs).Find(&subjects).Error; err != nil {
				return err
			}
			if err := db.Model(&teacher).Association("Subjects").Append(subjects); err != nil {
				return err
			}
		}

		classIDs := make([]uint, 0, len(seed.classKeys))
		for _, key := range seed.classKeys {
			if id, ok := ids.ClassIDs[key]; ok {
				classIDs = append(classIDs, id)
			}
		}
		if len(classIDs) > 0 {
			var classes []academic.Class
			if err := db.Where("id IN ?", classIDs).Find(&classes).Error; err != nil {
				return err
			}
			if err := db.Model(&teacher).Association("Classes").Append(classes); err != nil {
				return err
			}
		}
	}

	return nil
}

// uniqueTeacherSlug mirrors the insert-and-retry slug pattern used by
// Service.Create, but as a check-then-pick helper since the seeder's own
// idempotency check (by email, above) already guarantees each demo teacher
// is only ever inserted once.
func uniqueTeacherSlug(db *gorm.DB, centerID uint, fullName string) (string, error) {
	base := utils.Slugify(fullName, "teacher")

	for attempt := 1; attempt <= maxSlugAttempts; attempt++ {
		slug := base
		if attempt > 1 {
			slug = fmt.Sprintf("%s-%d", base, attempt)
		}

		var count int64
		if err := db.Model(&Teacher{}).Where("center_id = ? AND slug = ?", centerID, slug).Count(&count).Error; err != nil {
			return "", err
		}
		if count == 0 {
			return slug, nil
		}
	}

	return "", ErrSlugGenerationFailed
}
