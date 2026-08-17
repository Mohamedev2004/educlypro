package centers

import (
	"context"
	"errors"
	"testing"

	"educlypro/modules/auth"
	"educlypro/modules/staff"

	"gorm.io/gorm"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockRepository struct {
	listFn        func(ctx context.Context, params ListParams) ([]auth.Center, error)
	countFn       func(ctx context.Context, search string) (int64, error)
	staffCountsFn func(ctx context.Context, centerIDs []uint) (map[uint]int64, error)

	findBySlugFn  func(ctx context.Context, slug string) (*auth.Center, error)
	createFn      func(ctx context.Context, center *auth.Center) error
	emailTakenFn  func(ctx context.Context, email string) (bool, error)
	findRoleIDFn  func(ctx context.Context, name string) (uint, error)
	assignOwnerFn func(ctx context.Context, centerID uint, user *auth.User) error
	staffListFn   func(ctx context.Context, centerID uint) ([]auth.User, error)
}

func (m *mockRepository) List(ctx context.Context, params ListParams) ([]auth.Center, error) {
	if m.listFn != nil {
		return m.listFn(ctx, params)
	}
	return nil, nil
}

func (m *mockRepository) Count(ctx context.Context, search string) (int64, error) {
	if m.countFn != nil {
		return m.countFn(ctx, search)
	}
	return 0, nil
}

func (m *mockRepository) StaffCounts(ctx context.Context, centerIDs []uint) (map[uint]int64, error) {
	if m.staffCountsFn != nil {
		return m.staffCountsFn(ctx, centerIDs)
	}
	return map[uint]int64{}, nil
}

func (m *mockRepository) FindBySlug(ctx context.Context, slug string) (*auth.Center, error) {
	if m.findBySlugFn != nil {
		return m.findBySlugFn(ctx, slug)
	}
	return nil, nil
}

func (m *mockRepository) Create(ctx context.Context, center *auth.Center) error {
	if m.createFn != nil {
		return m.createFn(ctx, center)
	}
	center.ID = 42
	return nil
}

func (m *mockRepository) EmailTaken(ctx context.Context, email string) (bool, error) {
	if m.emailTakenFn != nil {
		return m.emailTakenFn(ctx, email)
	}
	return false, nil
}

func (m *mockRepository) FindRoleIDByName(ctx context.Context, name string) (uint, error) {
	if m.findRoleIDFn != nil {
		return m.findRoleIDFn(ctx, name)
	}
	return 2, nil
}

func (m *mockRepository) AssignOwner(ctx context.Context, centerID uint, user *auth.User) error {
	if m.assignOwnerFn != nil {
		return m.assignOwnerFn(ctx, centerID, user)
	}
	user.ID = 7
	return nil
}

func (m *mockRepository) StaffList(ctx context.Context, centerID uint) ([]auth.User, error) {
	if m.staffListFn != nil {
		return m.staffListFn(ctx, centerID)
	}
	return nil, nil
}

type mockStaffService struct {
	createForCenterFn func(ctx context.Context, centerID, actorUserID uint, req *staff.CreateRequest) (*staff.Response, error)
}

func (m *mockStaffService) List(ctx context.Context, ownerUserID uint, params staff.ListParams) (*staff.ListResponse, error) {
	return nil, nil
}

func (m *mockStaffService) Create(ctx context.Context, ownerUserID uint, req *staff.CreateRequest) (*staff.Response, error) {
	return nil, nil
}

func (m *mockStaffService) CreateForCenter(ctx context.Context, centerID, actorUserID uint, req *staff.CreateRequest) (*staff.Response, error) {
	if m.createForCenterFn != nil {
		return m.createForCenterFn(ctx, centerID, actorUserID, req)
	}
	return &staff.Response{ID: 1, Username: req.Username, Email: req.Email, Role: req.Role, CenterID: centerID}, nil
}

func (m *mockStaffService) Update(ctx context.Context, ownerUserID, staffID uint, req *staff.UpdateRequest) (*staff.Response, error) {
	return nil, nil
}

func (m *mockStaffService) Delete(ctx context.Context, ownerUserID, staffID uint) error {
	return nil
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestService_List(t *testing.T) {
	t.Run("Success_WithOwnerAndStaffCount", func(t *testing.T) {
		owner := &auth.User{ID: 2, Username: "jane", Email: "jane@example.com"}
		repo := &mockRepository{
			listFn: func(ctx context.Context, params ListParams) ([]auth.Center, error) {
				return []auth.Center{{ID: 1, Name: "Downtown", Slug: "downtown", Owner: owner}}, nil
			},
			countFn: func(ctx context.Context, search string) (int64, error) {
				return 1, nil
			},
			staffCountsFn: func(ctx context.Context, centerIDs []uint) (map[uint]int64, error) {
				return map[uint]int64{1: 3}, nil
			},
		}
		svc := &service{repo: repo}

		resp, err := svc.List(context.Background(), ListParams{Page: 1, PerPage: 10})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(resp.Items) != 1 {
			t.Fatalf("expected 1 item, got %d", len(resp.Items))
		}

		item := resp.Items[0]
		if item.OwnerUsername == nil || *item.OwnerUsername != "jane" {
			t.Errorf("expected owner username 'jane', got %+v", item.OwnerUsername)
		}
		if item.StaffCount != 3 {
			t.Errorf("expected staff count 3, got %d", item.StaffCount)
		}
		if resp.Pagination.Total != 1 || resp.Pagination.TotalPages != 1 {
			t.Errorf("unexpected pagination: %+v", resp.Pagination)
		}
	})

	t.Run("Success_NoOwner", func(t *testing.T) {
		repo := &mockRepository{
			listFn: func(ctx context.Context, params ListParams) ([]auth.Center, error) {
				return []auth.Center{{ID: 1, Name: "Unclaimed", Slug: "unclaimed", Owner: nil}}, nil
			},
			countFn: func(ctx context.Context, search string) (int64, error) {
				return 1, nil
			},
		}
		svc := &service{repo: repo}

		resp, err := svc.List(context.Background(), ListParams{Page: 1, PerPage: 10})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Items[0].OwnerUsername != nil || resp.Items[0].OwnerEmail != nil {
			t.Errorf("expected nil owner fields, got %+v", resp.Items[0])
		}
	})

	t.Run("PaginationMath", func(t *testing.T) {
		repo := &mockRepository{
			countFn: func(ctx context.Context, search string) (int64, error) {
				return 25, nil
			},
		}
		svc := &service{repo: repo}

		resp, err := svc.List(context.Background(), ListParams{Page: 2, PerPage: 10})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Pagination.TotalPages != 3 {
			t.Errorf("expected 3 total pages, got %d", resp.Pagination.TotalPages)
		}
		if !resp.Pagination.HasNext || !resp.Pagination.HasPrev {
			t.Errorf("expected has_next and has_prev true on page 2 of 3, got %+v", resp.Pagination)
		}
	})
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestService_Create(t *testing.T) {
	t.Run("Success_SlugifiesName", func(t *testing.T) {
		var createdSlug string
		repo := &mockRepository{
			createFn: func(ctx context.Context, center *auth.Center) error {
				createdSlug = center.Slug
				center.ID = 1
				return nil
			},
		}
		svc := &service{repo: repo}

		resp, err := svc.Create(context.Background(), &CreateRequest{Name: "Downtown Campus"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if createdSlug != "downtown-campus" {
			t.Errorf("expected slug 'downtown-campus', got %q", createdSlug)
		}
		if resp.Name != "Downtown Campus" {
			t.Errorf("unexpected response: %+v", resp)
		}
	})

	t.Run("Success_StripsDiacritics", func(t *testing.T) {
		var createdSlug string
		repo := &mockRepository{
			createFn: func(ctx context.Context, center *auth.Center) error {
				createdSlug = center.Slug
				center.ID = 1
				return nil
			},
		}
		svc := &service{repo: repo}

		if _, err := svc.Create(context.Background(), &CreateRequest{Name: "Café de Paris"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if createdSlug != "cafe-de-paris" {
			t.Errorf("expected slug 'cafe-de-paris', got %q", createdSlug)
		}
	})

	t.Run("Success_FallsBackToCenterWhenNameHasNoAsciiContent", func(t *testing.T) {
		var createdSlug string
		repo := &mockRepository{
			createFn: func(ctx context.Context, center *auth.Center) error {
				createdSlug = center.Slug
				center.ID = 1
				return nil
			},
		}
		svc := &service{repo: repo}

		// A name that slugifies to empty (non-Latin script / symbols only)
		// must still produce a valid, non-empty slug.
		if _, err := svc.Create(context.Background(), &CreateRequest{Name: "مركز"}); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if createdSlug != "center" {
			t.Errorf("expected fallback slug 'center', got %q", createdSlug)
		}
	})

	t.Run("Success_RetriesOnDuplicateSlugConflict", func(t *testing.T) {
		// Simulates two concurrent requests racing for the same slug: the
		// first two attempts hit the DB's unique index (gorm.ErrDuplicatedKey),
		// exactly as a real race would surface it — no pre-check involved.
		var attempts []string
		repo := &mockRepository{
			createFn: func(ctx context.Context, center *auth.Center) error {
				attempts = append(attempts, center.Slug)
				if center.Slug == "downtown" || center.Slug == "downtown-2" {
					return gorm.ErrDuplicatedKey
				}
				center.ID = 1
				return nil
			},
		}
		svc := &service{repo: repo}

		resp, err := svc.Create(context.Background(), &CreateRequest{Name: "Downtown"})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Slug != "downtown-3" {
			t.Errorf("expected slug 'downtown-3', got %q", resp.Slug)
		}
		wantAttempts := []string{"downtown", "downtown-2", "downtown-3"}
		if len(attempts) != len(wantAttempts) {
			t.Fatalf("expected attempts %v, got %v", wantAttempts, attempts)
		}
		for i, want := range wantAttempts {
			if attempts[i] != want {
				t.Errorf("attempt %d: expected %q, got %q", i, want, attempts[i])
			}
		}
	})

	t.Run("GivesUpAfterMaxAttempts", func(t *testing.T) {
		repo := &mockRepository{
			createFn: func(ctx context.Context, center *auth.Center) error {
				return gorm.ErrDuplicatedKey
			},
		}
		svc := &service{repo: repo}

		_, err := svc.Create(context.Background(), &CreateRequest{Name: "Downtown"})
		if err != ErrSlugGenerationFailed {
			t.Fatalf("expected ErrSlugGenerationFailed, got %v", err)
		}
	})

	t.Run("PropagatesNonConflictErrors", func(t *testing.T) {
		boom := errors.New("connection reset")
		repo := &mockRepository{
			createFn: func(ctx context.Context, center *auth.Center) error {
				return boom
			},
		}
		svc := &service{repo: repo}

		_, err := svc.Create(context.Background(), &CreateRequest{Name: "Downtown"})
		if err != boom {
			t.Fatalf("expected the underlying error to propagate, got %v", err)
		}
	})
}

// ── GetDetail ─────────────────────────────────────────────────────────────────

func TestService_GetDetail(t *testing.T) {
	t.Run("Success_WithOwnerAndStaff", func(t *testing.T) {
		repo := &mockRepository{
			findBySlugFn: func(ctx context.Context, slug string) (*auth.Center, error) {
				return &auth.Center{
					ID: 1, Name: "Downtown", Slug: "downtown",
					Owner: &auth.User{ID: 2, Username: "jane", Email: "jane@example.com"},
				}, nil
			},
			staffListFn: func(ctx context.Context, centerID uint) ([]auth.User, error) {
				return []auth.User{{ID: 3, Username: "bob", Email: "bob@example.com", Role: auth.Role{Name: staff.RoleScanner}}}, nil
			},
		}
		svc := &service{repo: repo}

		detail, err := svc.GetDetail(context.Background(), "downtown")
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if detail.Owner == nil || detail.Owner.Username != "jane" {
			t.Errorf("expected owner 'jane', got %+v", detail.Owner)
		}
		if len(detail.Staff) != 1 || detail.Staff[0].Username != "bob" {
			t.Errorf("unexpected staff list: %+v", detail.Staff)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := &mockRepository{
			findBySlugFn: func(ctx context.Context, slug string) (*auth.Center, error) {
				return nil, nil
			},
		}
		svc := &service{repo: repo}

		if _, err := svc.GetDetail(context.Background(), "does-not-exist"); err != ErrCenterNotFound {
			t.Fatalf("expected ErrCenterNotFound, got %v", err)
		}
	})
}

// ── AddOwner ──────────────────────────────────────────────────────────────────

func TestService_AddOwner(t *testing.T) {
	baseRepo := func() *mockRepository {
		return &mockRepository{
			findBySlugFn: func(ctx context.Context, slug string) (*auth.Center, error) {
				return &auth.Center{ID: 1, Name: "Downtown", Slug: "downtown"}, nil
			},
		}
	}

	t.Run("Success", func(t *testing.T) {
		repo := baseRepo()
		svc := &service{repo: repo}

		resp, err := svc.AddOwner(context.Background(), "downtown", &CreateOwnerRequest{
			Username: "jane", Email: "jane@example.com", Password: "secret12",
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Username != "jane" {
			t.Errorf("unexpected response: %+v", resp)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := &mockRepository{
			findBySlugFn: func(ctx context.Context, slug string) (*auth.Center, error) {
				return nil, nil
			},
		}
		svc := &service{repo: repo}

		_, err := svc.AddOwner(context.Background(), "does-not-exist", &CreateOwnerRequest{
			Username: "jane", Email: "jane@example.com", Password: "secret12",
		})
		if err != ErrCenterNotFound {
			t.Fatalf("expected ErrCenterNotFound, got %v", err)
		}
	})

	t.Run("AlreadyAssigned", func(t *testing.T) {
		ownerID := uint(5)
		repo := &mockRepository{
			findBySlugFn: func(ctx context.Context, slug string) (*auth.Center, error) {
				return &auth.Center{ID: 1, Name: "Downtown", Slug: "downtown", OwnerID: &ownerID}, nil
			},
		}
		svc := &service{repo: repo}

		_, err := svc.AddOwner(context.Background(), "downtown", &CreateOwnerRequest{
			Username: "jane", Email: "jane@example.com", Password: "secret12",
		})
		if err != ErrOwnerAlreadyAssigned {
			t.Fatalf("expected ErrOwnerAlreadyAssigned, got %v", err)
		}
	})

	t.Run("EmailTaken", func(t *testing.T) {
		repo := baseRepo()
		repo.emailTakenFn = func(ctx context.Context, email string) (bool, error) {
			return true, nil
		}
		svc := &service{repo: repo}

		_, err := svc.AddOwner(context.Background(), "downtown", &CreateOwnerRequest{
			Username: "jane", Email: "jane@example.com", Password: "secret12",
		})
		if err != ErrEmailTaken {
			t.Fatalf("expected ErrEmailTaken, got %v", err)
		}
	})
}

// ── AddStaff ──────────────────────────────────────────────────────────────────

func TestService_AddStaff(t *testing.T) {
	t.Run("Success_DelegatesToStaffService", func(t *testing.T) {
		repo := &mockRepository{
			findBySlugFn: func(ctx context.Context, slug string) (*auth.Center, error) {
				return &auth.Center{ID: 1, Name: "Downtown", Slug: "downtown"}, nil
			},
		}
		var capturedCenterID, capturedActorID uint
		staffSvc := &mockStaffService{
			createForCenterFn: func(ctx context.Context, centerID, actorUserID uint, req *staff.CreateRequest) (*staff.Response, error) {
				capturedCenterID = centerID
				capturedActorID = actorUserID
				return &staff.Response{ID: 9, Username: req.Username, Role: req.Role, CenterID: centerID}, nil
			},
		}
		svc := &service{repo: repo, staffService: staffSvc}

		resp, err := svc.AddStaff(context.Background(), "downtown", 42, &AddStaffRequest{
			Username: "bob", Email: "bob@example.com", Password: "secret12", Role: staff.RoleScanner,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Username != "bob" {
			t.Errorf("unexpected response: %+v", resp)
		}
		if capturedCenterID != 1 || capturedActorID != 42 {
			t.Errorf("expected centerID=1 actorID=42, got centerID=%d actorID=%d", capturedCenterID, capturedActorID)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := &mockRepository{
			findBySlugFn: func(ctx context.Context, slug string) (*auth.Center, error) {
				return nil, nil
			},
		}
		svc := &service{repo: repo, staffService: &mockStaffService{}}

		_, err := svc.AddStaff(context.Background(), "does-not-exist", 42, &AddStaffRequest{
			Username: "bob", Email: "bob@example.com", Password: "secret12", Role: staff.RoleScanner,
		})
		if err != ErrCenterNotFound {
			t.Fatalf("expected ErrCenterNotFound, got %v", err)
		}
	})
}
