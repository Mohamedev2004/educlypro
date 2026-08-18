package teachers

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"educlypro/modules/academic"
	"educlypro/shared/utils"

	"gorm.io/gorm"
)

var (
	ErrTeacherNotFound         = errors.New("teacher not found")
	ErrEmailTaken              = errors.New("email already in use")
	ErrInvalidSubjectSelection = errors.New("one or more selected subjects are invalid")
	ErrInvalidClassSelection   = errors.New("one or more selected classes are invalid")
	// ErrSlugGenerationFailed is returned in the practically-impossible case
	// where maxSlugAttempts consecutive candidates all collide.
	ErrSlugGenerationFailed = errors.New("could not generate a unique slug")
)

// maxSlugAttempts bounds the insert-and-retry loop in Create. Each attempt
// is a real, atomic INSERT — collisions are detected via the DB's unique
// index on (center_id, slug), not a separate check-then-insert query, so
// concurrent requests for the same name can never both succeed with the
// same slug. Mirrors centers.Service.Create.
const maxSlugAttempts = 100

type Service interface {
	List(ctx context.Context, ownerUserID uint, params ListParams) (*ListResponse, error)
	// GetBySlug returns a single teacher's full detail (including subjects
	// and classes) — used by the teacher detail page.
	GetBySlug(ctx context.Context, ownerUserID uint, slug string) (*Response, error)
	// Create adds a teacher with just its core identity fields. Subject/
	// class assignment happens afterward via UpdateSubjects/UpdateClasses.
	Create(ctx context.Context, ownerUserID uint, req *CreateRequest) (*Response, error)
	// Update saves a teacher's core identity fields only.
	Update(ctx context.Context, ownerUserID, teacherID uint, req *UpdateRequest) (*Response, error)
	// UpdateSubjects replaces a teacher's full subject assignment.
	UpdateSubjects(ctx context.Context, ownerUserID, teacherID uint, req *UpdateSubjectsRequest) (*Response, error)
	// UpdateClasses replaces a teacher's full class assignment.
	UpdateClasses(ctx context.Context, ownerUserID, teacherID uint, req *UpdateClassesRequest) (*Response, error)
	Delete(ctx context.Context, ownerUserID, teacherID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, ownerUserID uint, params ListParams) (*ListResponse, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	teacherList, err := s.repo.ListByCenter(ctx, centerID, params)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.CountByCenter(ctx, centerID, params.Search)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	items := make([]Response, 0, len(teacherList))
	for _, t := range teacherList {
		items = append(items, toResponse(t))
	}

	return &ListResponse{
		Items: items,
		Pagination: PaginationMeta{
			Page:       params.Page,
			PerPage:    params.PerPage,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
	}, nil
}

func (s *service) GetBySlug(ctx context.Context, ownerUserID uint, slug string) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	teacher, err := s.repo.FindBySlugInCenter(ctx, centerID, slug)
	if err != nil {
		return nil, err
	}
	if teacher == nil {
		return nil, ErrTeacherNotFound
	}

	resp := toResponse(*teacher)
	return &resp, nil
}

func (s *service) Create(ctx context.Context, ownerUserID uint, req *CreateRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	taken, err := s.repo.EmailTaken(ctx, centerID, req.Email, 0)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrEmailTaken
	}

	base := utils.Slugify(req.FullName, "teacher")
	fullName := strings.TrimSpace(req.FullName)
	phone := strings.TrimSpace(req.Phone)

	for attempt := 1; attempt <= maxSlugAttempts; attempt++ {
		slug := base
		if attempt > 1 {
			slug = fmt.Sprintf("%s-%d", base, attempt)
		}

		teacher := &Teacher{
			CenterID: centerID,
			Slug:     slug,
			FullName: fullName,
			Email:    req.Email,
			Phone:    phone,
		}
		err := s.repo.Create(ctx, teacher)
		if err == nil {
			resp := toResponse(*teacher)
			return &resp, nil
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			continue
		}
		return nil, err
	}

	return nil, ErrSlugGenerationFailed
}

func (s *service) Update(ctx context.Context, ownerUserID, teacherID uint, req *UpdateRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, teacherID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrTeacherNotFound
	}

	if !strings.EqualFold(req.Email, existing.Email) {
		taken, err := s.repo.EmailTaken(ctx, centerID, req.Email, teacherID)
		if err != nil {
			return nil, err
		}
		if taken {
			return nil, ErrEmailTaken
		}
	}

	teacher := &Teacher{
		ID:       teacherID,
		CenterID: centerID,
		Slug:     existing.Slug,
		FullName: strings.TrimSpace(req.FullName),
		Email:    req.Email,
		Phone:    strings.TrimSpace(req.Phone),
	}
	if err := s.repo.UpdateInfo(ctx, teacher); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailTaken
		}
		return nil, err
	}
	teacher.CreatedAt = existing.CreatedAt
	teacher.Subjects = existing.Subjects
	teacher.Classes = existing.Classes

	resp := toResponse(*teacher)
	return &resp, nil
}

func (s *service) UpdateSubjects(ctx context.Context, ownerUserID, teacherID uint, req *UpdateSubjectsRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, teacherID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrTeacherNotFound
	}

	subjects, err := s.resolveSubjects(ctx, centerID, req.SubjectIDs)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ReplaceSubjects(ctx, teacherID, subjects); err != nil {
		return nil, err
	}
	existing.Subjects = subjects

	resp := toResponse(*existing)
	return &resp, nil
}

func (s *service) UpdateClasses(ctx context.Context, ownerUserID, teacherID uint, req *UpdateClassesRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, teacherID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrTeacherNotFound
	}

	classes, err := s.resolveClasses(ctx, centerID, req.ClassIDs)
	if err != nil {
		return nil, err
	}

	if err := s.repo.ReplaceClasses(ctx, teacherID, classes); err != nil {
		return nil, err
	}
	existing.Classes = classes

	resp := toResponse(*existing)
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, ownerUserID, teacherID uint) error {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return err
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, teacherID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTeacherNotFound
	}

	return s.repo.SoftDelete(ctx, teacherID)
}

// resolveSubjects de-duplicates the given ids, then validates that every one
// of them belongs to the caller's own center, returning the matching rows
// ready to hand to the repository's ReplaceSubjects.
func (s *service) resolveSubjects(ctx context.Context, centerID uint, subjectIDs []uint) ([]academic.Subject, error) {
	deduped := dedupeIDs(subjectIDs)
	if len(deduped) == 0 {
		return nil, nil
	}

	subjects, err := s.repo.FindSubjectsInCenter(ctx, centerID, deduped)
	if err != nil {
		return nil, err
	}
	if len(subjects) != len(deduped) {
		return nil, ErrInvalidSubjectSelection
	}
	return subjects, nil
}

// resolveClasses is resolveSubjects's counterpart for classes.
func (s *service) resolveClasses(ctx context.Context, centerID uint, classIDs []uint) ([]academic.Class, error) {
	deduped := dedupeIDs(classIDs)
	if len(deduped) == 0 {
		return nil, nil
	}

	classes, err := s.repo.FindClassesInCenter(ctx, centerID, deduped)
	if err != nil {
		return nil, err
	}
	if len(classes) != len(deduped) {
		return nil, ErrInvalidClassSelection
	}
	return classes, nil
}

func dedupeIDs(ids []uint) []uint {
	seen := make(map[uint]struct{}, len(ids))
	deduped := make([]uint, 0, len(ids))
	for _, id := range ids {
		if _, ok := seen[id]; ok {
			continue
		}
		seen[id] = struct{}{}
		deduped = append(deduped, id)
	}
	return deduped
}

func toResponse(t Teacher) Response {
	resp := Response{
		ID:        t.ID,
		Slug:      t.Slug,
		FullName:  t.FullName,
		Email:     t.Email,
		Phone:     t.Phone,
		Subjects:  make([]SubjectRef, 0, len(t.Subjects)),
		Classes:   make([]ClassRef, 0, len(t.Classes)),
		CreatedAt: t.CreatedAt.Format(time.RFC3339),
	}
	for _, sub := range t.Subjects {
		resp.Subjects = append(resp.Subjects, SubjectRef{ID: sub.ID, Name: sub.Name})
	}
	for _, cls := range t.Classes {
		resp.Classes = append(resp.Classes, ClassRef{ID: cls.ID, Name: cls.Name})
	}
	return resp
}
