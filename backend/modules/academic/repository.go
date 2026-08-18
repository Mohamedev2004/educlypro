package academic

import (
	"context"
	"errors"
	"strings"

	"educlypro/modules/auth"

	"gorm.io/gorm"
)

type Repository interface {
	// FindOwnerCenterID resolves the center a center_owner user belongs to.
	// Mirrors staff.Repository.FindOwnerCenterID.
	FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error)

	// ListByCenter returns every grade for the center, with majors and
	// subjects preloaded, ordered oldest-first (creation order).
	ListByCenter(ctx context.Context, centerID uint) ([]Grade, error)
	// CountByCenter returns how many grades the center has — the source of
	// truth for whether a center_owner has completed academic setup.
	CountByCenter(ctx context.Context, centerID uint) (int64, error)

	GradeNameExists(ctx context.Context, centerID uint, name string) (bool, error)
	CreateGrade(ctx context.Context, grade *Grade) error
	FindGradeByIDInCenter(ctx context.Context, centerID, gradeID uint) (*Grade, error)
	DeleteGrade(ctx context.Context, gradeID uint) error

	MajorNameExists(ctx context.Context, gradeID uint, name string) (bool, error)
	// CreateMajorWithClass creates the major and its one auto-generated
	// class (major.ID is set on class.MajorID) atomically, so a failure on
	// either side never leaves a major without its class.
	CreateMajorWithClass(ctx context.Context, major *Major, class *Class) error
	// FindMajorInCenter looks up a major and verifies (via its grade) that it
	// belongs to the given center — the authorization check for major-scoped
	// actions.
	FindMajorInCenter(ctx context.Context, centerID, majorID uint) (*Major, error)
	DeleteMajor(ctx context.Context, majorID uint) error

	SubjectNameExists(ctx context.Context, majorID uint, name string) (bool, error)
	CreateSubject(ctx context.Context, subject *Subject) error
	// FindSubjectInCenter looks up a subject and verifies (via its major and
	// grade) that it belongs to the given center.
	FindSubjectInCenter(ctx context.Context, centerID, subjectID uint) (*Subject, error)
	DeleteSubject(ctx context.Context, subjectID uint) error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).Select("center_id").First(&user, ownerUserID).Error; err != nil {
		return 0, err
	}
	if user.CenterID == nil {
		return 0, errors.New("owner has no center")
	}
	return *user.CenterID, nil
}

func (r *repository) ListByCenter(ctx context.Context, centerID uint) ([]Grade, error) {
	var grades []Grade
	err := r.db.WithContext(ctx).
		Where("center_id = ?", centerID).
		Preload("Majors", func(db *gorm.DB) *gorm.DB {
			return db.Order("majors.id ASC")
		}).
		Preload("Majors.Subjects", func(db *gorm.DB) *gorm.DB {
			return db.Order("subjects.id ASC")
		}).
		Preload("Majors.Classes", func(db *gorm.DB) *gorm.DB {
			return db.Order("classes.id ASC")
		}).
		Order("grades.id ASC").
		Find(&grades).Error
	return grades, err
}

func (r *repository) CountByCenter(ctx context.Context, centerID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Grade{}).Where("center_id = ?", centerID).Count(&count).Error
	return count, err
}

func (r *repository) GradeNameExists(ctx context.Context, centerID uint, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Grade{}).
		Where("center_id = ? AND LOWER(name) = LOWER(?)", centerID, strings.TrimSpace(name)).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) CreateGrade(ctx context.Context, grade *Grade) error {
	return r.db.WithContext(ctx).Create(grade).Error
}

func (r *repository) FindGradeByIDInCenter(ctx context.Context, centerID, gradeID uint) (*Grade, error) {
	var grade Grade
	err := r.db.WithContext(ctx).
		Where("id = ? AND center_id = ?", gradeID, centerID).
		First(&grade).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &grade, nil
}

func (r *repository) DeleteGrade(ctx context.Context, gradeID uint) error {
	return r.db.WithContext(ctx).Delete(&Grade{}, gradeID).Error
}

func (r *repository) MajorNameExists(ctx context.Context, gradeID uint, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Major{}).
		Where("grade_id = ? AND LOWER(name) = LOWER(?)", gradeID, strings.TrimSpace(name)).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) CreateMajorWithClass(ctx context.Context, major *Major, class *Class) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(major).Error; err != nil {
			return err
		}
		class.MajorID = major.ID
		return tx.Create(class).Error
	})
}

func (r *repository) FindMajorInCenter(ctx context.Context, centerID, majorID uint) (*Major, error) {
	var major Major
	err := r.db.WithContext(ctx).
		Joins("JOIN grades ON grades.id = majors.grade_id").
		Where("majors.id = ? AND grades.center_id = ?", majorID, centerID).
		Select("majors.*").
		First(&major).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &major, nil
}

func (r *repository) DeleteMajor(ctx context.Context, majorID uint) error {
	return r.db.WithContext(ctx).Delete(&Major{}, majorID).Error
}

func (r *repository) SubjectNameExists(ctx context.Context, majorID uint, name string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Subject{}).
		Where("major_id = ? AND LOWER(name) = LOWER(?)", majorID, strings.TrimSpace(name)).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) CreateSubject(ctx context.Context, subject *Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

func (r *repository) FindSubjectInCenter(ctx context.Context, centerID, subjectID uint) (*Subject, error) {
	var subject Subject
	err := r.db.WithContext(ctx).
		Joins("JOIN majors ON majors.id = subjects.major_id").
		Joins("JOIN grades ON grades.id = majors.grade_id").
		Where("subjects.id = ? AND grades.center_id = ?", subjectID, centerID).
		Select("subjects.*").
		First(&subject).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &subject, nil
}

func (r *repository) DeleteSubject(ctx context.Context, subjectID uint) error {
	return r.db.WithContext(ctx).Delete(&Subject{}, subjectID).Error
}
