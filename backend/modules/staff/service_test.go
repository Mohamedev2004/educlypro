package staff

import (
	"context"
	"testing"

	"educlypro/modules/auth"

	"github.com/ThreeDotsLabs/watermill/message"
)

// ── Mocks ─────────────────────────────────────────────────────────────────────

type mockPublisher struct {
	publishedTopics []string
}

func (m *mockPublisher) Publish(topic string, messages ...*message.Message) error {
	m.publishedTopics = append(m.publishedTopics, topic)
	return nil
}

func (m *mockPublisher) Close() error { return nil }

type mockTokenRevoker struct {
	deleteTokensByUserFn func(userID uint, tokenType string) error
}

func (m *mockTokenRevoker) DeleteTokensByUserID(userID uint, tokenType string) error {
	if m.deleteTokensByUserFn != nil {
		return m.deleteTokensByUserFn(userID, tokenType)
	}
	return nil
}

type mockRepository struct {
	findOwnerCenterIDFn func(ctx context.Context, ownerUserID uint) (uint, error)
	listByCenterFn      func(ctx context.Context, centerID uint, params ListParams) ([]auth.User, error)
	countByCenterFn     func(ctx context.Context, centerID uint, role, search string) (int64, error)
	findByIDInCenterFn  func(ctx context.Context, centerID, staffID uint) (*auth.User, error)
	findRoleIDByNameFn  func(ctx context.Context, name string) (uint, error)
	createFn            func(ctx context.Context, user *auth.User) error
	updateFn            func(ctx context.Context, user *auth.User) error
	softDeleteFn        func(ctx context.Context, staffID uint) error
	emailTakenFn        func(ctx context.Context, email string, excludeUserID uint) (bool, error)
}

func (m *mockRepository) FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error) {
	return m.findOwnerCenterIDFn(ctx, ownerUserID)
}
func (m *mockRepository) ListByCenter(ctx context.Context, centerID uint, params ListParams) ([]auth.User, error) {
	if m.listByCenterFn != nil {
		return m.listByCenterFn(ctx, centerID, params)
	}
	return nil, nil
}
func (m *mockRepository) CountByCenter(ctx context.Context, centerID uint, role, search string) (int64, error) {
	if m.countByCenterFn != nil {
		return m.countByCenterFn(ctx, centerID, role, search)
	}
	return 0, nil
}
func (m *mockRepository) FindByIDInCenter(ctx context.Context, centerID, staffID uint) (*auth.User, error) {
	if m.findByIDInCenterFn != nil {
		return m.findByIDInCenterFn(ctx, centerID, staffID)
	}
	return nil, nil
}
func (m *mockRepository) FindRoleIDByName(ctx context.Context, name string) (uint, error) {
	if m.findRoleIDByNameFn != nil {
		return m.findRoleIDByNameFn(ctx, name)
	}
	return 2, nil
}
func (m *mockRepository) Create(ctx context.Context, user *auth.User) error {
	if m.createFn != nil {
		return m.createFn(ctx, user)
	}
	user.ID = 42
	return nil
}
func (m *mockRepository) Update(ctx context.Context, user *auth.User) error {
	if m.updateFn != nil {
		return m.updateFn(ctx, user)
	}
	return nil
}
func (m *mockRepository) SoftDelete(ctx context.Context, staffID uint) error {
	if m.softDeleteFn != nil {
		return m.softDeleteFn(ctx, staffID)
	}
	return nil
}
func (m *mockRepository) EmailTaken(ctx context.Context, email string, excludeUserID uint) (bool, error) {
	if m.emailTakenFn != nil {
		return m.emailTakenFn(ctx, email, excludeUserID)
	}
	return false, nil
}

// ── Fixtures ──────────────────────────────────────────────────────────────────

const ownerCenterID uint = 7

func newRepoForOwner() *mockRepository {
	return &mockRepository{
		findOwnerCenterIDFn: func(ctx context.Context, ownerUserID uint) (uint, error) {
			return ownerCenterID, nil
		},
	}
}

// ── Create ────────────────────────────────────────────────────────────────────

func TestService_Create(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		pub := &mockPublisher{}
		repo := newRepoForOwner()
		svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: pub}

		resp, err := svc.Create(context.Background(), 1, &CreateRequest{
			Username: "jane", Email: "jane@example.com", Password: "secret12",
			Role: RoleScanner, CenterID: ownerCenterID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if resp.Role != RoleScanner || resp.CenterID != ownerCenterID {
			t.Errorf("unexpected response: %+v", resp)
		}

		found := false
		for _, topic := range pub.publishedTopics {
			if topic == "system.events.v1.staff.created" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected 'system.events.v1.staff.created' event, got topics: %v", pub.publishedTopics)
		}
	})

	t.Run("CenterMismatch", func(t *testing.T) {
		pub := &mockPublisher{}
		repo := newRepoForOwner()
		svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: pub}

		_, err := svc.Create(context.Background(), 1, &CreateRequest{
			Username: "jane", Email: "jane@example.com", Password: "secret12",
			Role: RoleScanner, CenterID: ownerCenterID + 1, // tampered
		})
		if err != ErrCenterMismatch {
			t.Fatalf("expected ErrCenterMismatch, got %v", err)
		}

		found := false
		for _, topic := range pub.publishedTopics {
			if topic == "system.events.v1.staff.center_mismatch" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected 'system.events.v1.staff.center_mismatch' event, got topics: %v", pub.publishedTopics)
		}
	})

	t.Run("EmailTaken", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.emailTakenFn = func(ctx context.Context, email string, excludeUserID uint) (bool, error) {
			return true, nil
		}
		svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: &mockPublisher{}}

		_, err := svc.Create(context.Background(), 1, &CreateRequest{
			Username: "jane", Email: "jane@example.com", Password: "secret12",
			Role: RoleScanner, CenterID: ownerCenterID,
		})
		if err != ErrEmailTaken {
			t.Fatalf("expected ErrEmailTaken, got %v", err)
		}
	})
}

// ── Update ────────────────────────────────────────────────────────────────────

func TestService_Update(t *testing.T) {
	existingStaff := &auth.User{
		ID: 5, Username: "old", Email: "old@example.com", Password: "existing-hash",
		CenterID: func() *uint { c := ownerCenterID; return &c }(),
		Role:     auth.Role{ID: 2, Name: RoleScanner},
	}

	t.Run("Success_PasswordUnchangedWhenBlank", func(t *testing.T) {
		var capturedPassword string
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, staffID uint) (*auth.User, error) {
			return existingStaff, nil
		}
		repo.updateFn = func(ctx context.Context, user *auth.User) error {
			capturedPassword = user.Password
			return nil
		}
		svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: &mockPublisher{}}

		resp, err := svc.Update(context.Background(), 1, 5, &UpdateRequest{
			Username: "new-name", Email: "old@example.com",
			Role: RoleScanner, Password: nil, CenterID: ownerCenterID,
		})
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if capturedPassword != existingStaff.Password {
			t.Errorf("expected password to remain unchanged, got %q", capturedPassword)
		}
		if resp.Username != "new-name" {
			t.Errorf("expected username to update, got %q", resp.Username)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, staffID uint) (*auth.User, error) {
			return nil, nil
		}
		svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: &mockPublisher{}}

		_, err := svc.Update(context.Background(), 1, 999, &UpdateRequest{
			Username: "x", Email: "x@example.com", Role: RoleScanner, CenterID: ownerCenterID,
		})
		if err != ErrStaffNotFound {
			t.Fatalf("expected ErrStaffNotFound, got %v", err)
		}
	})

	t.Run("CenterMismatch", func(t *testing.T) {
		repo := newRepoForOwner()
		svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: &mockPublisher{}}

		_, err := svc.Update(context.Background(), 1, 5, &UpdateRequest{
			Username: "x", Email: "x@example.com", Role: RoleScanner, CenterID: ownerCenterID + 1,
		})
		if err != ErrCenterMismatch {
			t.Fatalf("expected ErrCenterMismatch, got %v", err)
		}
	})
}

// ── Delete ────────────────────────────────────────────────────────────────────

func TestService_Delete(t *testing.T) {
	t.Run("Success_RevokesTokens", func(t *testing.T) {
		revokedTypes := []string{}
		repo := newRepoForOwner()
		repo.findByIDInCenterFn = func(ctx context.Context, centerID, staffID uint) (*auth.User, error) {
			return &auth.User{ID: 5, Username: "jane", Email: "jane@example.com"}, nil
		}
		revoker := &mockTokenRevoker{
			deleteTokensByUserFn: func(userID uint, tokenType string) error {
				revokedTypes = append(revokedTypes, tokenType)
				return nil
			},
		}
		pub := &mockPublisher{}
		svc := &service{repo: repo, tokenRevoker: revoker, publisher: pub}

		if err := svc.Delete(context.Background(), 1, 5); err != nil {
			t.Fatalf("expected no error, got %v", err)
		}
		if len(revokedTypes) != 2 {
			t.Errorf("expected access+refresh tokens revoked, got %v", revokedTypes)
		}

		found := false
		for _, topic := range pub.publishedTopics {
			if topic == "system.events.v1.staff.deleted" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected 'system.events.v1.staff.deleted' event, got topics: %v", pub.publishedTopics)
		}
	})

	t.Run("NotFound", func(t *testing.T) {
		repo := newRepoForOwner()
		svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: &mockPublisher{}}

		if err := svc.Delete(context.Background(), 1, 999); err != ErrStaffNotFound {
			t.Fatalf("expected ErrStaffNotFound, got %v", err)
		}
	})
}

// ── List ──────────────────────────────────────────────────────────────────────

func TestService_List(t *testing.T) {
	repo := newRepoForOwner()
	repo.listByCenterFn = func(ctx context.Context, centerID uint, params ListParams) ([]auth.User, error) {
		return []auth.User{
			{ID: 5, Username: "jane", Email: "jane@example.com", Role: auth.Role{Name: RoleScanner}},
		}, nil
	}
	repo.countByCenterFn = func(ctx context.Context, centerID uint, role, search string) (int64, error) {
		return 1, nil
	}
	svc := &service{repo: repo, tokenRevoker: &mockTokenRevoker{}, publisher: &mockPublisher{}}

	resp, err := svc.List(context.Background(), 1, ListParams{Page: 1, PerPage: 10})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if len(resp.Items) != 1 || resp.Pagination.Total != 1 {
		t.Errorf("unexpected list response: %+v", resp)
	}
}
