package teachers

import (
	"context"
	"errors"
	"testing"

	"educlypro/modules/academic"

	"gorm.io/gorm"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockRepository struct {
	findOwnerCenterIDFn    func(ctx context.Context, ownerUserID uint) (uint, error)
	listByCenterFn         func(ctx context.Context, centerID uint, params ListParams) ([]Teacher, error)
	countByCenterFn        func(ctx context.Context, centerID uint, search string) (int64, error)
	findByIDInCenterFn     func(ctx context.Context, centerID, teacherID uint) (*Teacher, error)
	findBySlugInCenterFn   func(ctx context.Context, centerID uint, slug string) (*Teacher, error)
	findSubjectsInCenterFn func(ctx context.Context, centerID uint, subjectIDs []uint) ([]academic.Subject, error)
	findClassesInCenterFn  func(ctx context.Context, centerID uint, classIDs []uint) ([]academic.Class, error)
	emailTakenFn           func(ctx context.Context, centerID uint, email string, excludeID uint) (bool, error)
	createFn               func(ctx context.Context, teacher *Teacher) error
	updateInfoFn           func(ctx context.Context, teacher *Teacher) error
	replaceSubjectsFn      func(ctx context.Context, teacherID uint, subjects []academic.Subject) error
	replaceClassesFn       func(ctx context.Context, teacherID uint, classes []academic.Class) error
	softDeleteFn           func(ctx context.Context, teacherID uint) error
}

const ownerCenterID uint = 7

func newRepoForOwner() *mockRepository {
	return &mockRepository{
		findOwnerCenterIDFn: func(ctx context.Context, ownerUserID uint) (uint, error) {
			return ownerCenterID, nil
		},
	}
}

func (m *mockRepository) FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error) {
	return m.findOwnerCenterIDFn(ctx, ownerUserID)
}
func (m *mockRepository) ListByCenter(ctx context.Context, centerID uint, params ListParams) ([]Teacher, error) {
	if m.listByCenterFn != nil {
		return m.listByCenterFn(ctx, centerID, params)
	}
	return nil, nil
}
func (m *mockRepository) CountByCenter(ctx context.Context, centerID uint, search string) (int64, error) {
	if m.countByCenterFn != nil {
		return m.countByCenterFn(ctx, centerID, search)
	}
	return 0, nil
}
func (m *mockRepository) FindByIDInCenter(ctx context.Context, centerID, teacherID uint) (*Teacher, error) {
	if m.findByIDInCenterFn != nil {
		return m.findByIDInCenterFn(ctx, centerID, teacherID)
	}
	return &Teacher{ID: teacherID, CenterID: centerID}, nil
}
func (m *mockRepository) FindBySlugInCenter(ctx context.Context, centerID uint, slug string) (*Teacher, error) {
	if m.findBySlugInCenterFn != nil {
		return m.findBySlugInCenterFn(ctx, centerID, slug)
	}
	return &Teacher{ID: 1, CenterID: centerID, Slug: slug}, nil
}
func (m *mockRepository) FindSubjectsInCenter(ctx context.Context, centerID uint, subjectIDs []uint) ([]academic.Subject, error) {
	if m.findSubjectsInCenterFn != nil {
		return m.findSubjectsInCenterFn(ctx, centerID, subjectIDs)
	}
	subjects := make([]academic.Subject, 0, len(subjectIDs))
	for _, id := range subjectIDs {
		subjects = append(subjects, academic.Subject{ID: id, Name: "Subject"})
	}
	return subjects, nil
}
func (m *mockRepository) FindClassesInCenter(ctx context.Context, centerID uint, classIDs []uint) ([]academic.Class, error) {
	if m.findClassesInCenterFn != nil {
		return m.findClassesInCenterFn(ctx, centerID, classIDs)
	}
	classes := make([]academic.Class, 0, len(classIDs))
	for _, id := range classIDs {
		classes = append(classes, academic.Class{ID: id, Name: "Class"})
	}
	return classes, nil
}
func (m *mockRepository) EmailTaken(ctx context.Context, centerID uint, email string, excludeID uint) (bool, error) {
	if m.emailTakenFn != nil {
		return m.emailTakenFn(ctx, centerID, email, excludeID)
	}
	return false, nil
}
func (m *mockRepository) Create(ctx context.Context, teacher *Teacher) error {
	if m.createFn != nil {
		return m.createFn(ctx, teacher)
	}
	teacher.ID = 42
	return nil
}
func (m *mockRepository) UpdateInfo(ctx context.Context, teacher *Teacher) error {
	if m.updateInfoFn != nil {
		return m.updateInfoFn(ctx, teacher)
	}
	return nil
}
func (m *mockRepository) ReplaceSubjects(ctx context.Context, teacherID uint, subjects []academic.Subject) error {
	if m.replaceSubjectsFn != nil {
		return m.replaceSubjectsFn(ctx, teacherID, subjects)
	}
	return nil
}
func (m *mockRepository) ReplaceClasses(ctx context.Context, teacherID uint, classes []academic.Class) error {
	if m.replaceClassesFn != nil {
		return m.replaceClassesFn(ctx, teacherID, classes)
	}
	return nil
}
func (m *mockRepository) SoftDelete(ctx context.Context, teacherID uint) error {
	if m.softDeleteFn != nil {
		return m.softDeleteFn(ctx, teacherID)
	}
	return nil
}

// ── GetBySlug ─────────────────────────────────────────────────────────────────

func TestService_GetBySlug(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findBySlugInCenterFn = func(ctx context.Context, centerID uint, slug string) (*Teacher, error) {
			return &Teacher{ID: 5, CenterID: centerID, Slug: slug, FullName: "Jane"}, nil
		}
		svc := &service{repo: repo}

		resp, err := svc.GetBySlug(context.Background(), 1, "jane")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.FullName != "Jane" {
			t.Errorf("expected full name 'Jane', got %q", resp.FullName)
		}
		if resp.Slug != "jane" {
			t.Errorf("expected slug 'jane', got %q", resp.Slug)
		}
	})

	t.Run("NotFoundOutsideOwnCenter", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findBySlugInCenterFn = func(ctx context.Context, centerID uint, slug string) (*Teacher, error) {
			return nil, nil
		}
		svc := &service{repo: repo}

		_, err := svc.GetBySlug(context.Background(), 1, "missing")
		if !errors.Is(err, ErrTeacherNotFound) {
			t.Fatalf("expected ErrTeacherNotFound, got %v", err)
		}
	})
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := &service{repo: newRepoForOwner()}

		resp, err := svc.Create(context.Background(), 1, &CreateRequest{
			FullName: "Youssef El Amrani", Email: "youssef@example.com", Phone: "0600000000",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Slug != "youssef-el-amrani" {
			t.Errorf("expected slug 'youssef-el-amrani', got %q", resp.Slug)
		}
	})

	t.Run("RetriesSlugOnCollision", func(t *testing.T) {
		attempts := 0
		repo := newRepoForOwner()
		repo.createFn = func(ctx context.Context, teacher *Teacher) error {
			attempts++
			if attempts == 1 {
				return gorm.ErrDuplicatedKey
			}
			teacher.ID = 42
			return nil
		}
		svc := &service{repo: repo}

		resp, err := svc.Create(context.Background(), 1, &CreateRequest{
			FullName: "Jane", Email: "jane@example.com", Phone: "0600000000",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Slug != "jane-2" {
			t.Errorf("expected slug 'jane-2' after collision retry, got %q", resp.Slug)
		}
	})

	t.Run("EmailTaken", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.emailTakenFn = func(ctx context.Context, centerID uint, email string, excludeID uint) (bool, error) {
			return true, nil
		}
		svc := &service{repo: repo}

		_, err := svc.Create(context.Background(), 1, &CreateRequest{
			FullName: "Jane", Email: "jane@example.com", Phone: "0600000000",
		})
		if !errors.Is(err, ErrEmailTaken) {
			t.Fatalf("expected ErrEmailTaken, got %v", err)
		}
	})
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestService_Update(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, teacherID uint) (*Teacher, error) {
			return nil, nil
		}
		svc := &service{repo: repo}

		_, err := svc.Update(context.Background(), 1, 999, &UpdateRequest{
			FullName: "Jane", Email: "jane@example.com", Phone: "0600000000",
		})
		if !errors.Is(err, ErrTeacherNotFound) {
			t.Fatalf("expected ErrTeacherNotFound, got %v", err)
		}
	})

	t.Run("Success_PreservesSlugSubjectsAndClasses", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, teacherID uint) (*Teacher, error) {
			return &Teacher{
				ID: teacherID, CenterID: centerID, Slug: "old-slug", Email: "old@example.com",
				Subjects: []academic.Subject{{ID: 3, Name: "Math"}},
				Classes:  []academic.Class{{ID: 11, Name: "Class A"}},
			}, nil
		}
		svc := &service{repo: repo}

		resp, err := svc.Update(context.Background(), 1, 5, &UpdateRequest{
			FullName: "Jane", Email: "old@example.com", Phone: "0600000000",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.FullName != "Jane" {
			t.Errorf("expected full name 'Jane', got %q", resp.FullName)
		}
		if resp.Slug != "old-slug" {
			t.Errorf("expected slug to stay 'old-slug', got %q", resp.Slug)
		}
		if len(resp.Subjects) != 1 || len(resp.Classes) != 1 {
			t.Errorf("expected existing subjects/classes to be preserved, got %+v / %+v", resp.Subjects, resp.Classes)
		}
	})
}

// ── UpdateSubjects ────────────────────────────────────────────────────────────

func TestService_UpdateSubjects(t *testing.T) {
	t.Run("NotFound", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, teacherID uint) (*Teacher, error) {
			return nil, nil
		}
		svc := &service{repo: repo}

		_, err := svc.UpdateSubjects(context.Background(), 1, 999, &UpdateSubjectsRequest{SubjectIDs: []uint{1}})
		if !errors.Is(err, ErrTeacherNotFound) {
			t.Fatalf("expected ErrTeacherNotFound, got %v", err)
		}
	})

	t.Run("Success_ReplacesSubjects", func(t *testing.T) {
		var replaced []academic.Subject
		repo := newRepoForOwner()
		repo.replaceSubjectsFn = func(ctx context.Context, teacherID uint, subjects []academic.Subject) error {
			replaced = subjects
			return nil
		}
		svc := &service{repo: repo}

		resp, err := svc.UpdateSubjects(context.Background(), 1, 5, &UpdateSubjectsRequest{SubjectIDs: []uint{3, 3}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(replaced) != 1 || replaced[0].ID != 3 {
			t.Errorf("expected duplicate subject id deduped to [3], got %+v", replaced)
		}
		if len(resp.Subjects) != 1 {
			t.Errorf("expected 1 subject in response, got %d", len(resp.Subjects))
		}
	})

	t.Run("InvalidSubjectSelection", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findSubjectsInCenterFn = func(ctx context.Context, centerID uint, subjectIDs []uint) ([]academic.Subject, error) {
			return []academic.Subject{{ID: 1, Name: "Math"}}, nil // only 1 of 2 requested resolved
		}
		svc := &service{repo: repo}

		_, err := svc.UpdateSubjects(context.Background(), 1, 5, &UpdateSubjectsRequest{SubjectIDs: []uint{1, 999}})
		if !errors.Is(err, ErrInvalidSubjectSelection) {
			t.Fatalf("expected ErrInvalidSubjectSelection, got %v", err)
		}
	})
}

// ── UpdateClasses ─────────────────────────────────────────────────────────────

func TestService_UpdateClasses(t *testing.T) {
	t.Run("Success_ReplacesClasses", func(t *testing.T) {
		var replaced []academic.Class
		repo := newRepoForOwner()
		repo.replaceClassesFn = func(ctx context.Context, teacherID uint, classes []academic.Class) error {
			replaced = classes
			return nil
		}
		svc := &service{repo: repo}

		resp, err := svc.UpdateClasses(context.Background(), 1, 5, &UpdateClassesRequest{ClassIDs: []uint{11}})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(replaced) != 1 || replaced[0].ID != 11 {
			t.Errorf("expected classes replaced with [11], got %+v", replaced)
		}
		if len(resp.Classes) != 1 {
			t.Errorf("expected 1 class in response, got %d", len(resp.Classes))
		}
	})

	t.Run("InvalidClassSelection", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findClassesInCenterFn = func(ctx context.Context, centerID uint, classIDs []uint) ([]academic.Class, error) {
			return nil, nil // none of the requested classes resolved
		}
		svc := &service{repo: repo}

		_, err := svc.UpdateClasses(context.Background(), 1, 5, &UpdateClassesRequest{ClassIDs: []uint{999}})
		if !errors.Is(err, ErrInvalidClassSelection) {
			t.Fatalf("expected ErrInvalidClassSelection, got %v", err)
		}
	})
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestService_Delete(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		deleted := false
		repo := newRepoForOwner()
		repo.softDeleteFn = func(ctx context.Context, teacherID uint) error {
			deleted = true
			return nil
		}
		svc := &service{repo: repo}

		if err := svc.Delete(context.Background(), 1, 5); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !deleted {
			t.Error("expected SoftDelete to be called")
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, teacherID uint) (*Teacher, error) {
			return nil, nil
		}
		svc := &service{repo: repo}

		if err := svc.Delete(context.Background(), 1, 999); !errors.Is(err, ErrTeacherNotFound) {
			t.Fatalf("expected ErrTeacherNotFound, got %v", err)
		}
	})
}
