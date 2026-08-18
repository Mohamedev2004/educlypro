package academic

import (
	"context"
	"errors"
	"testing"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockRepository struct {
	findOwnerCenterIDFn func(ctx context.Context, ownerUserID uint) (uint, error)
	listByCenterFn      func(ctx context.Context, centerID uint) ([]Grade, error)
	countByCenterFn     func(ctx context.Context, centerID uint) (int64, error)

	gradeNameExistsFn    func(ctx context.Context, centerID uint, name string) (bool, error)
	createGradeFn        func(ctx context.Context, grade *Grade) error
	findGradeByIDInCtrFn func(ctx context.Context, centerID, gradeID uint) (*Grade, error)
	deleteGradeFn        func(ctx context.Context, gradeID uint) error

	majorNameExistsFn      func(ctx context.Context, gradeID uint, name string) (bool, error)
	createMajorWithClassFn func(ctx context.Context, major *Major, class *Class) error
	findMajorInCtrFn       func(ctx context.Context, centerID, majorID uint) (*Major, error)
	deleteMajorFn          func(ctx context.Context, majorID uint) error

	subjectNameExistsFn func(ctx context.Context, majorID uint, name string) (bool, error)
	createSubjectFn     func(ctx context.Context, subject *Subject) error
	findSubjectInCtrFn  func(ctx context.Context, centerID, subjectID uint) (*Subject, error)
	deleteSubjectFn     func(ctx context.Context, subjectID uint) error
}

func (m *mockRepository) FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error) {
	if m.findOwnerCenterIDFn != nil {
		return m.findOwnerCenterIDFn(ctx, ownerUserID)
	}
	return 1, nil
}

func (m *mockRepository) ListByCenter(ctx context.Context, centerID uint) ([]Grade, error) {
	if m.listByCenterFn != nil {
		return m.listByCenterFn(ctx, centerID)
	}
	return nil, nil
}

func (m *mockRepository) CountByCenter(ctx context.Context, centerID uint) (int64, error) {
	if m.countByCenterFn != nil {
		return m.countByCenterFn(ctx, centerID)
	}
	return 0, nil
}

func (m *mockRepository) GradeNameExists(ctx context.Context, centerID uint, name string) (bool, error) {
	if m.gradeNameExistsFn != nil {
		return m.gradeNameExistsFn(ctx, centerID, name)
	}
	return false, nil
}

func (m *mockRepository) CreateGrade(ctx context.Context, grade *Grade) error {
	if m.createGradeFn != nil {
		return m.createGradeFn(ctx, grade)
	}
	grade.ID = 100
	return nil
}

func (m *mockRepository) FindGradeByIDInCenter(ctx context.Context, centerID, gradeID uint) (*Grade, error) {
	if m.findGradeByIDInCtrFn != nil {
		return m.findGradeByIDInCtrFn(ctx, centerID, gradeID)
	}
	return &Grade{ID: gradeID, CenterID: centerID}, nil
}

func (m *mockRepository) DeleteGrade(ctx context.Context, gradeID uint) error {
	if m.deleteGradeFn != nil {
		return m.deleteGradeFn(ctx, gradeID)
	}
	return nil
}

func (m *mockRepository) MajorNameExists(ctx context.Context, gradeID uint, name string) (bool, error) {
	if m.majorNameExistsFn != nil {
		return m.majorNameExistsFn(ctx, gradeID, name)
	}
	return false, nil
}

func (m *mockRepository) CreateMajorWithClass(ctx context.Context, major *Major, class *Class) error {
	if m.createMajorWithClassFn != nil {
		return m.createMajorWithClassFn(ctx, major, class)
	}
	major.ID = 200
	class.ID = 201
	class.MajorID = major.ID
	return nil
}

func (m *mockRepository) FindMajorInCenter(ctx context.Context, centerID, majorID uint) (*Major, error) {
	if m.findMajorInCtrFn != nil {
		return m.findMajorInCtrFn(ctx, centerID, majorID)
	}
	return &Major{ID: majorID}, nil
}

func (m *mockRepository) DeleteMajor(ctx context.Context, majorID uint) error {
	if m.deleteMajorFn != nil {
		return m.deleteMajorFn(ctx, majorID)
	}
	return nil
}

func (m *mockRepository) SubjectNameExists(ctx context.Context, majorID uint, name string) (bool, error) {
	if m.subjectNameExistsFn != nil {
		return m.subjectNameExistsFn(ctx, majorID, name)
	}
	return false, nil
}

func (m *mockRepository) CreateSubject(ctx context.Context, subject *Subject) error {
	if m.createSubjectFn != nil {
		return m.createSubjectFn(ctx, subject)
	}
	subject.ID = 300
	return nil
}

func (m *mockRepository) FindSubjectInCenter(ctx context.Context, centerID, subjectID uint) (*Subject, error) {
	if m.findSubjectInCtrFn != nil {
		return m.findSubjectInCtrFn(ctx, centerID, subjectID)
	}
	return &Subject{ID: subjectID}, nil
}

func (m *mockRepository) DeleteSubject(ctx context.Context, subjectID uint) error {
	if m.deleteSubjectFn != nil {
		return m.deleteSubjectFn(ctx, subjectID)
	}
	return nil
}

// ── AddGrade ──────────────────────────────────────────────────────────────────

func TestService_AddGrade(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := &service{repo: &mockRepository{}}

		resp, err := svc.AddGrade(context.Background(), 1, &AddGradeRequest{Name: "  Grade 10  "})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Grade 10" {
			t.Errorf("expected trimmed name 'Grade 10', got %q", resp.Name)
		}
		if resp.ID != 100 {
			t.Errorf("expected id 100, got %d", resp.ID)
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		repo := &mockRepository{
			gradeNameExistsFn: func(ctx context.Context, centerID uint, name string) (bool, error) {
				return true, nil
			},
		}
		svc := &service{repo: repo}

		_, err := svc.AddGrade(context.Background(), 1, &AddGradeRequest{Name: "Grade 10"})
		if !errors.Is(err, ErrGradeExists) {
			t.Fatalf("expected ErrGradeExists, got %v", err)
		}
	})
}

// ── RemoveGrade ───────────────────────────────────────────────────────────────

func TestService_RemoveGrade(t *testing.T) {
	t.Run("NotFoundOutsideOwnCenter", func(t *testing.T) {
		repo := &mockRepository{
			findGradeByIDInCtrFn: func(ctx context.Context, centerID, gradeID uint) (*Grade, error) {
				return nil, nil // grade belongs to a different center
			},
		}
		svc := &service{repo: repo}

		err := svc.RemoveGrade(context.Background(), 1, 999)
		if !errors.Is(err, ErrGradeNotFound) {
			t.Fatalf("expected ErrGradeNotFound, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		deleted := false
		repo := &mockRepository{
			deleteGradeFn: func(ctx context.Context, gradeID uint) error {
				deleted = true
				return nil
			},
		}
		svc := &service{repo: repo}

		if err := svc.RemoveGrade(context.Background(), 1, 5); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !deleted {
			t.Error("expected DeleteGrade to be called")
		}
	})
}

// ── AddMajor ──────────────────────────────────────────────────────────────────

func TestService_AddMajor(t *testing.T) {
	t.Run("GradeNotInOwnCenter", func(t *testing.T) {
		repo := &mockRepository{
			findGradeByIDInCtrFn: func(ctx context.Context, centerID, gradeID uint) (*Grade, error) {
				return nil, nil
			},
		}
		svc := &service{repo: repo}

		_, err := svc.AddMajor(context.Background(), 1, 42, &AddMajorRequest{Name: "Science"})
		if !errors.Is(err, ErrGradeNotFound) {
			t.Fatalf("expected ErrGradeNotFound, got %v", err)
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		repo := &mockRepository{
			majorNameExistsFn: func(ctx context.Context, gradeID uint, name string) (bool, error) {
				return true, nil
			},
		}
		svc := &service{repo: repo}

		_, err := svc.AddMajor(context.Background(), 1, 42, &AddMajorRequest{Name: "Science"})
		if !errors.Is(err, ErrMajorExists) {
			t.Fatalf("expected ErrMajorExists, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		svc := &service{repo: &mockRepository{}}

		resp, err := svc.AddMajor(context.Background(), 1, 42, &AddMajorRequest{Name: "Science"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Science" {
			t.Errorf("expected name 'Science', got %q", resp.Name)
		}
	})

	t.Run("AutoCreatesOneClassNamedAfterTheMajor", func(t *testing.T) {
		var capturedClass *Class
		repo := &mockRepository{
			createMajorWithClassFn: func(ctx context.Context, major *Major, class *Class) error {
				major.ID = 200
				class.MajorID = major.ID
				class.ID = 201
				capturedClass = class
				return nil
			},
		}
		svc := &service{repo: repo}

		resp, err := svc.AddMajor(context.Background(), 1, 42, &AddMajorRequest{Name: "Science"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if capturedClass == nil {
			t.Fatal("expected CreateMajorWithClass to be called with a class")
		}
		if capturedClass.Name != "Science class" {
			t.Errorf("expected class name 'Science class', got %q", capturedClass.Name)
		}
		if resp.Class == nil || resp.Class.Name != "Science class" {
			t.Errorf("expected response to include the class, got %+v", resp.Class)
		}
	})
}

// ── AddSubject ────────────────────────────────────────────────────────────────

func TestService_AddSubject(t *testing.T) {
	t.Run("MajorNotInOwnCenter", func(t *testing.T) {
		repo := &mockRepository{
			findMajorInCtrFn: func(ctx context.Context, centerID, majorID uint) (*Major, error) {
				return nil, nil
			},
		}
		svc := &service{repo: repo}

		_, err := svc.AddSubject(context.Background(), 1, 42, &AddSubjectRequest{Name: "Physics"})
		if !errors.Is(err, ErrMajorNotFound) {
			t.Fatalf("expected ErrMajorNotFound, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		svc := &service{repo: &mockRepository{}}

		resp, err := svc.AddSubject(context.Background(), 1, 42, &AddSubjectRequest{Name: "Physics"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Physics" {
			t.Errorf("expected name 'Physics', got %q", resp.Name)
		}
	})
}

// ── Tree ──────────────────────────────────────────────────────────────────────

func TestService_Tree(t *testing.T) {
	repo := &mockRepository{
		listByCenterFn: func(ctx context.Context, centerID uint) ([]Grade, error) {
			return []Grade{
				{
					ID: 1, Name: "Grade 10",
					Majors: []Major{
						{ID: 1, Name: "Science", Subjects: []Subject{{ID: 1, Name: "Physics"}}},
					},
				},
			}, nil
		},
	}
	svc := &service{repo: repo}

	resp, err := svc.Tree(context.Background(), 1)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Grades) != 1 || len(resp.Grades[0].Majors) != 1 || len(resp.Grades[0].Majors[0].Subjects) != 1 {
		t.Fatalf("unexpected tree shape: %+v", resp)
	}
}
