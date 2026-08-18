package subcenters

import (
	"context"
	"errors"
	"strings"
	"time"

	"educlypro/modules/auth"

	"gorm.io/gorm"
)

var (
	ErrSubCenterNotFound = errors.New("sub-center not found")
	ErrSubCenterExists   = errors.New("sub-center already exists")
	ErrSubCenterHasStaff = errors.New("sub-center still has staff assigned")
)

type Service interface {
	// List returns every sub-center in the caller's own center.
	List(ctx context.Context, ownerUserID uint) (*ListResponse, error)
	// ListForCenter is List's "caller already knows and is authorized to act
	// on this center" counterpart — used by the centers module so a
	// super_admin can populate a sub-center picker for an arbitrary center's
	// "add staff" flow. Mirrors staff.Service.CreateForCenter.
	ListForCenter(ctx context.Context, centerID uint) (*ListResponse, error)

	Create(ctx context.Context, ownerUserID uint, req *CreateRequest) (*Response, error)
	Update(ctx context.Context, ownerUserID, subCenterID uint, req *UpdateRequest) (*Response, error)
	Delete(ctx context.Context, ownerUserID, subCenterID uint) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) List(ctx context.Context, ownerUserID uint) (*ListResponse, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	return s.ListForCenter(ctx, centerID)
}

func (s *service) ListForCenter(ctx context.Context, centerID uint) (*ListResponse, error) {
	list, err := s.repo.ListByCenter(ctx, centerID)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(list))
	for _, sc := range list {
		ids = append(ids, sc.ID)
	}

	staffCounts, err := s.repo.StaffCounts(ctx, ids)
	if err != nil {
		return nil, err
	}

	items := make([]Response, 0, len(list))
	for _, sc := range list {
		items = append(items, toResponse(sc, staffCounts[sc.ID]))
	}

	return &ListResponse{Items: items}, nil
}

func (s *service) Create(ctx context.Context, ownerUserID uint, req *CreateRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(req.Name)

	exists, err := s.repo.NameExists(ctx, centerID, name, 0)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, ErrSubCenterExists
	}

	subCenter := &auth.SubCenter{CenterID: centerID, Name: name}
	if err := s.repo.Create(ctx, subCenter); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSubCenterExists
		}
		return nil, err
	}

	resp := toResponse(*subCenter, 0)
	return &resp, nil
}

func (s *service) Update(ctx context.Context, ownerUserID, subCenterID uint, req *UpdateRequest) (*Response, error) {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, subCenterID)
	if err != nil {
		return nil, err
	}
	if existing == nil {
		return nil, ErrSubCenterNotFound
	}

	name := strings.TrimSpace(req.Name)

	if !strings.EqualFold(name, existing.Name) {
		exists, err := s.repo.NameExists(ctx, centerID, name, subCenterID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, ErrSubCenterExists
		}
	}

	updated := &auth.SubCenter{ID: subCenterID, CenterID: centerID, Name: name}
	if err := s.repo.Update(ctx, updated); err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrSubCenterExists
		}
		return nil, err
	}
	updated.CreatedAt = existing.CreatedAt

	staffCounts, err := s.repo.StaffCounts(ctx, []uint{subCenterID})
	if err != nil {
		return nil, err
	}

	resp := toResponse(*updated, staffCounts[subCenterID])
	return &resp, nil
}

func (s *service) Delete(ctx context.Context, ownerUserID, subCenterID uint) error {
	centerID, err := s.repo.FindOwnerCenterID(ctx, ownerUserID)
	if err != nil {
		return err
	}

	existing, err := s.repo.FindByIDInCenter(ctx, centerID, subCenterID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrSubCenterNotFound
	}

	staffCounts, err := s.repo.StaffCounts(ctx, []uint{subCenterID})
	if err != nil {
		return err
	}
	if staffCounts[subCenterID] > 0 {
		return ErrSubCenterHasStaff
	}

	return s.repo.Delete(ctx, subCenterID)
}

func toResponse(sc auth.SubCenter, staffCount int64) Response {
	return Response{
		ID:         sc.ID,
		Name:       sc.Name,
		StaffCount: staffCount,
		CreatedAt:  sc.CreatedAt.Format(time.RFC3339),
	}
}
