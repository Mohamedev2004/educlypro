package subcenters

import (
	"context"
	"errors"
	"testing"

	"educlypro/modules/auth"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockRepository struct {
	findOwnerCenterIDFn func(ctx context.Context, ownerUserID uint) (uint, error)
	listByCenterFn      func(ctx context.Context, centerID uint) ([]auth.SubCenter, error)
	staffCountsFn       func(ctx context.Context, subCenterIDs []uint) (map[uint]int64, error)
	nameExistsFn        func(ctx context.Context, centerID uint, name string, excludeID uint) (bool, error)
	createFn            func(ctx context.Context, subCenter *auth.SubCenter) error
	findByIDInCenterFn  func(ctx context.Context, centerID, subCenterID uint) (*auth.SubCenter, error)
	updateFn            func(ctx context.Context, subCenter *auth.SubCenter) error
	deleteFn            func(ctx context.Context, subCenterID uint) error
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
func (m *mockRepository) ListByCenter(ctx context.Context, centerID uint) ([]auth.SubCenter, error) {
	if m.listByCenterFn != nil {
		return m.listByCenterFn(ctx, centerID)
	}
	return nil, nil
}
func (m *mockRepository) StaffCounts(ctx context.Context, subCenterIDs []uint) (map[uint]int64, error) {
	if m.staffCountsFn != nil {
		return m.staffCountsFn(ctx, subCenterIDs)
	}
	return map[uint]int64{}, nil
}
func (m *mockRepository) NameExists(ctx context.Context, centerID uint, name string, excludeID uint) (bool, error) {
	if m.nameExistsFn != nil {
		return m.nameExistsFn(ctx, centerID, name, excludeID)
	}
	return false, nil
}
func (m *mockRepository) Create(ctx context.Context, subCenter *auth.SubCenter) error {
	if m.createFn != nil {
		return m.createFn(ctx, subCenter)
	}
	subCenter.ID = 42
	return nil
}
func (m *mockRepository) FindByIDInCenter(ctx context.Context, centerID, subCenterID uint) (*auth.SubCenter, error) {
	if m.findByIDInCenterFn != nil {
		return m.findByIDInCenterFn(ctx, centerID, subCenterID)
	}
	return &auth.SubCenter{ID: subCenterID, CenterID: centerID, Name: "Existing"}, nil
}
func (m *mockRepository) Update(ctx context.Context, subCenter *auth.SubCenter) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, subCenter)
	}
	return nil
}
func (m *mockRepository) Delete(ctx context.Context, subCenterID uint) error {
	if m.deleteFn != nil {
		return m.deleteFn(ctx, subCenterID)
	}
	return nil
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		svc := &service{repo: newRepoForOwner()}

		resp, err := svc.Create(context.Background(), 1, &CreateRequest{Name: "  Main Branch  "})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Main Branch" {
			t.Errorf("expected trimmed name 'Main Branch', got %q", resp.Name)
		}
	})

	t.Run("DuplicateName", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.nameExistsFn = func(ctx context.Context, centerID uint, name string, excludeID uint) (bool, error) {
			return true, nil
		}
		svc := &service{repo: repo}

		_, err := svc.Create(context.Background(), 1, &CreateRequest{Name: "Main Branch"})
		if !errors.Is(err, ErrSubCenterExists) {
			t.Fatalf("expected ErrSubCenterExists, got %v", err)
		}
	})
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestService_Update(t *testing.T) {
	t.Run("NotFoundOutsideOwnCenter", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, subCenterID uint) (*auth.SubCenter, error) {
			return nil, nil
		}
		svc := &service{repo: repo}

		_, err := svc.Update(context.Background(), 1, 999, &UpdateRequest{Name: "New Name"})
		if !errors.Is(err, ErrSubCenterNotFound) {
			t.Fatalf("expected ErrSubCenterNotFound, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		svc := &service{repo: newRepoForOwner()}

		resp, err := svc.Update(context.Background(), 1, 5, &UpdateRequest{Name: "Renamed"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Name != "Renamed" {
			t.Errorf("expected name 'Renamed', got %q", resp.Name)
		}
	})
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestService_Delete(t *testing.T) {
	t.Run("BlockedWhenStaffAssigned", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.staffCountsFn = func(ctx context.Context, subCenterIDs []uint) (map[uint]int64, error) {
			return map[uint]int64{5: 3}, nil
		}
		svc := &service{repo: repo}

		err := svc.Delete(context.Background(), 1, 5)
		if !errors.Is(err, ErrSubCenterHasStaff) {
			t.Fatalf("expected ErrSubCenterHasStaff, got %v", err)
		}
	})

	t.Run("Success", func(t *testing.T) {
		deleted := false
		repo := newRepoForOwner()
		repo.deleteFn = func(ctx context.Context, subCenterID uint) error {
			deleted = true
			return nil
		}
		svc := &service{repo: repo}

		if err := svc.Delete(context.Background(), 1, 5); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if !deleted {
			t.Error("expected Delete to be called")
		}
	})
}
