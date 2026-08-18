package academic

import (
	"context"
	"errors"
	"strings"

	"gorm.io/gorm"
)

var (
	ErrGradeNotFound   = errors.New("grade not found")
	ErrGradeExists     = errors.New("grade already exists")
	ErrMajorNotFound   = errors.New("major not found")
	ErrMajorExists     = errors.New("major already exists")
	ErrSubjectNotFound = errors.New("subject not found")
	ErrSubjectExists   = errors.New("subject already exists")
)

type Service interface {
	// Tree returns the calling center_owner's full grade -> major -> subject
	// structure, used to render the onboarding wizard.
	Tree(ctx context.Context, ownerUserID uint) (*TreeResponse, error)

	AddGrade(ctx context.Context, ownerUserID uint, req *AddGradeRequest) (*GradeResponse, error)
	RemoveGrade(ctx context.Context, ownerUserID, gradeID uint) error

	AddMajor(ctx context.Context, ownerUserID, gradeID uint, req *AddMajorRequest) (*MajorResponse, error)
	RemoveMajor(ctx context.Context, ownerUserID, majorID uint) error

	AddSubject(ctx context.Context, ownerUserID, majorID uint, req *AddSubjectRequest) (*SubjectResponse, error)
	RemoveSubject(ctx context.Context, ownerUserID, subjectID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Tree(ctx context.Context, ownerUserID uint) (*TreeResponse, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	grades, err := s.repo.ListByCenter(ctx, centerID)
	if err != nil {
		return nil, err
	}

	resp := &TreeResponse{Grades: make([]GradeResponse, 0, len(grades))}
	for _, g := range grades {
		resp.Grades = append(resp.Grades, toGradeResponse(g))
	}
	return resp, nil
}

func (s *service) AddGrade(ctx context.Context, ownerUserID uint, req *AddGradeRequest) (*GradeResponse, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)

	exists, err := s.repo.GradeNameExists(ctx, centerID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrGradeExists
	}

	grade := &Grade{CenterID: centerID, Name: name}
	if err := s.repo.CreateGrade(ctx, grade); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrGradeExists
		}
		return nil, err
	}

	resp := toGradeResponse(*grade)
	return &resp, nil
}

func (s *service) RemoveGrade(ctx context.Context, ownerUserID, gradeID uint) error {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return err
	}

	grade, err := s.repo.FindGradeByIDInCenter(ctx, centerID, gradeID)
	if err != nil {
		return err
	}
	if grade == nil {
		return ErrGradeNotFound
	}

	return s.repo.DeleteGrade(ctx, gradeID)
}

func (s *service) AddMajor(ctx context.Context, ownerUserID, gradeID uint, req *AddMajorRequest) (*MajorResponse, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	grade, err := s.repo.FindGradeByIDInCenter(ctx, centerID, gradeID)
	if err != nil {
		return nil, err
	}
	if grade == nil {
		return nil, ErrGradeNotFound
	}

	name := strings.TrimSpace(req.Name)

	exists, err := s.repo.MajorNameExists(ctx, gradeID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrMajorExists
	}

	major := &Major{GradeID: gradeID, Name: name}
	class := &Class{Name: ClassNameForMajor(name)}
	if err := s.repo.CreateMajorWithClass(ctx, major, class); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrMajorExists
		}
		return nil, err
	}
	major.Classes = []Class{*class}

	resp := toMajorResponse(*major)
	return &resp, nil
}

// ClassNameForMajor is the fixed naming rule for a major's one
// auto-generated class — exported so the seeder can produce the exact same
// name a real AddMajor call would.
func ClassNameForMajor(majorName string) string {
	return majorName + " class"
}

func (s *service) RemoveMajor(ctx context.Context, ownerUserID, majorID uint) error {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return err
	}

	major, err := s.repo.FindMajorInCenter(ctx, centerID, majorID)
	if err != nil {
		return err
	}
	if major == nil {
		return ErrMajorNotFound
	}

	return s.repo.DeleteMajor(ctx, majorID)
}

func (s *service) AddSubject(ctx context.Context, ownerUserID, majorID uint, req *AddSubjectRequest) (*SubjectResponse, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	major, err := s.repo.FindMajorInCenter(ctx, centerID, majorID)
	if err != nil {
		return nil, err
	}
	if major == nil {
		return nil, ErrMajorNotFound
	}

	name := strings.TrimSpace(req.Name)

	exists, err := s.repo.SubjectNameExists(ctx, majorID, name)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSubjectExists
	}

	subject := &Subject{MajorID: majorID, Name: name}
	if err := s.repo.CreateSubject(ctx, subject); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSubjectExists
		}
		return nil, err
	}

	resp := toSubjectResponse(*subject)
	return &resp, nil
}

func (s *service) RemoveSubject(ctx context.Context, ownerUserID, subjectID uint) error {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return err
	}

	subject, err := s.repo.FindSubjectInCenter(ctx, centerID, subjectID)
	if err != nil {
		return err
	}
	if subject == nil {
		return ErrSubjectNotFound
	}

	return s.repo.DeleteSubject(ctx, subjectID)
}

func toSubjectResponse(s Subject) SubjectResponse {
	return SubjectResponse{ID: s.ID, Name: s.Name}
}

func toMajorResponse(m Major) MajorResponse {
	resp := MajorResponse{ID: m.ID, Name: m.Name, Subjects: make([]SubjectResponse, 0, len(m.Subjects))}
	for _, sub := range m.Subjects {
		resp.Subjects = append(resp.Subjects, toSubjectResponse(sub))
	}
	if len(m.Classes) > 0 {
		resp.Class = &ClassResponse{ID: m.Classes[0].ID, Name: m.Classes[0].Name}
	}
	return resp
}

func toGradeResponse(g Grade) GradeResponse {
	resp := GradeResponse{ID: g.ID, Name: g.Name, Majors: make([]MajorResponse, 0, len(g.Majors))}
	for _, m := range g.Majors {
		resp.Majors = append(resp.Majors, toMajorResponse(m))
	}
	return resp
}
