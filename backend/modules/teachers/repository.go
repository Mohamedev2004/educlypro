package teachers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"educlypro/modules/academic"
	"educlypro/modules/auth"

	"gorm.io/gorm"
)

type Repository interface {
	// FindOwnerCenterID resolves the center a center_owner user belongs to.
	// Mirrors staff.Repository.FindOwnerCenterID.
	FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error)

	ListByCenter(ctx context.Context, centerID uint, params ListParams) ([]Teacher, error)
	CountByCenter(ctx context.Context, centerID uint, search string) (int64, error)
	FindByIDInCenter(ctx context.Context, centerID, teacherID uint) (*Teacher, error)
	// FindBySlugInCenter is how a teacher's detail page looks up its
	// subject, mirroring centers.Repository.FindBySlug.
	FindBySlugInCenter(ctx context.Context, centerID uint, slug string) (*Teacher, error)

	// FindSubjectsInCenter resolves the given subject ids to their rows,
	// filtered to those actually belonging to centerID — the caller compares
	// the returned count against len(subjectIDs) to detect ids that don't
	// belong (wrong center, or don't exist at all).
	FindSubjectsInCenter(ctx context.Context, centerID uint, subjectIDs []uint) ([]academic.Subject, error)
	// FindClassesInCenter is FindSubjectsInCenter's counterpart for classes.
	FindClassesInCenter(ctx context.Context, centerID uint, classIDs []uint) ([]academic.Class, error)

	EmailTaken(ctx context.Context, centerID uint, email string, excludeID uint) (bool, error)
	// Create inserts a teacher with no subject/class assignment yet — those
	// are added afterward via ReplaceSubjects/ReplaceClasses, from the
	// teacher detail page.
	Create(ctx context.Context, teacher *Teacher) error
	// UpdateInfo saves the teacher's core identity fields only (not its
	// slug, which is fixed at creation, nor its subject/class assignment).
	UpdateInfo(ctx context.Context, teacher *Teacher) error
	// ReplaceSubjects replaces a teacher's full subject assignment with the
	// given (already-resolved, already-owned) set — never re-saves the
	// subject rows themselves.
	ReplaceSubjects(ctx context.Context, teacherID uint, subjects []academic.Subject) error
	// ReplaceClasses is ReplaceSubjects's counterpart for classes.
	ReplaceClasses(ctx context.Context, teacherID uint, classes []academic.Class) error
	SoftDelete(ctx context.Context, teacherID uint) error
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

func (r *repository) teacherQuery(centerID uint, search string) *gorm.DB {
	query := r.db.Model(&Teacher{}).Where("center_id = ?", centerID)

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(full_name) LIKE ? OR LOWER(email) LIKE ?", like, like)
	}

	return query
}

func (r *repository) ListByCenter(ctx context.Context, centerID uint, params ListParams) ([]Teacher, error) {
	var teachers []Teacher

	orderColumn := "created_at"
	if params.Sort != "" {
		orderColumn = params.Sort
	}

	direction := "DESC"
	if strings.EqualFold(params.Direction, "asc") {
		direction = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage

	err := r.teacherQuery(centerID, params.Search).
		WithContext(ctx).
		Preload("Subjects").
		Preload("Classes").
		Order(fmt.Sprintf("%s %s", orderColumn, direction)).
		Limit(params.PerPage).
		Offset(offset).
		Find(&teachers).Error

	return teachers, err
}

func (r *repository) CountByCenter(ctx context.Context, centerID uint, search string) (int64, error) {
	var count int64
	err := r.teacherQuery(centerID, search).WithContext(ctx).Count(&count).Error
	return count, err
}

func (r *repository) FindByIDInCenter(ctx context.Context, centerID, teacherID uint) (*Teacher, error) {
	var teacher Teacher
	err := r.db.WithContext(ctx).
		Preload("Subjects").
		Preload("Classes").
		Where("id = ? AND center_id = ?", teacherID, centerID).
		First(&teacher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &teacher, nil
}

func (r *repository) FindBySlugInCenter(ctx context.Context, centerID uint, slug string) (*Teacher, error) {
	var teacher Teacher
	err := r.db.WithContext(ctx).
		Preload("Subjects").
		Preload("Classes").
		Where("slug = ? AND center_id = ?", slug, centerID).
		First(&teacher).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &teacher, nil
}

func (r *repository) FindSubjectsInCenter(ctx context.Context, centerID uint, subjectIDs []uint) ([]academic.Subject, error) {
	if len(subjectIDs) == 0 {
		return nil, nil
	}

	var subjects []academic.Subject
	err := r.db.WithContext(ctx).
		Joins("JOIN majors ON majors.id = subjects.major_id").
		Joins("JOIN grades ON grades.id = majors.grade_id").
		Where("subjects.id IN ? AND grades.center_id = ?", subjectIDs, centerID).
		Select("subjects.*").
		Find(&subjects).Error
	return subjects, err
}

func (r *repository) FindClassesInCenter(ctx context.Context, centerID uint, classIDs []uint) ([]academic.Class, error) {
	if len(classIDs) == 0 {
		return nil, nil
	}

	var classes []academic.Class
	err := r.db.WithContext(ctx).
		Joins("JOIN majors ON majors.id = classes.major_id").
		Joins("JOIN grades ON grades.id = majors.grade_id").
		Where("classes.id IN ? AND grades.center_id = ?", classIDs, centerID).
		Select("classes.*").
		Find(&classes).Error
	return classes, err
}

func (r *repository) EmailTaken(ctx context.Context, centerID uint, email string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&Teacher{}).
		Where("center_id = ? AND email = ? AND id <> ?", centerID, email, excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) Create(ctx context.Context, teacher *Teacher) error {
	return r.db.WithContext(ctx).Omit("Subjects", "Classes").Create(teacher).Error
}

func (r *repository) UpdateInfo(ctx context.Context, teacher *Teacher) error {
	return r.db.WithContext(ctx).Model(&Teacher{}).
		Where("id = ?", teacher.ID).
		Updates(map[string]any{
			"full_name": teacher.FullName,
			"email":     teacher.Email,
			"phone":     teacher.Phone,
		}).Error
}

func (r *repository) ReplaceSubjects(ctx context.Context, teacherID uint, subjects []academic.Subject) error {
	return r.db.WithContext(ctx).Model(&Teacher{ID: teacherID}).Association("Subjects").Replace(subjects)
}

func (r *repository) ReplaceClasses(ctx context.Context, teacherID uint, classes []academic.Class) error {
	return r.db.WithContext(ctx).Model(&Teacher{ID: teacherID}).Association("Classes").Replace(classes)
}

func (r *repository) SoftDelete(ctx context.Context, teacherID uint) error {
	return r.db.WithContext(ctx).Delete(&Teacher{}, teacherID).Error
}
