package centers

import (
	"context"
	"errors"
	"fmt"
	"time"

	"educlypro/modules/auth"
	"educlypro/modules/staff"
	"educlypro/modules/subcenters"
	"educlypro/shared/utils"

	"gorm.io/gorm"
)

var (
	ErrCenterNotFound       = errors.New("center not found")
	ErrOwnerAlreadyAssigned = errors.New("center already has an owner")
	ErrEmailTaken           = errors.New("email already in use")
	// ErrSlugGenerationFailed is returned in the practically-impossible case
	// where maxSlugAttempts consecutive candidates all collide.
	ErrSlugGenerationFailed = errors.New("could not generate a unique slug")
)

const ownerRoleName = "center_owner"

type Service interface {
	List(ctx context.Context, params ListParams) (*ListResponse, error)
	Create(ctx context.Context, req *CreateRequest) (*Response, error)
	GetDetail(ctx context.Context, slug string) (*DetailResponse, error)
	AddOwner(ctx context.Context, slug string, req *CreateOwnerRequest) (*OwnerResponse, error)
	AddStaff(ctx context.Context, slug string, actorUserID uint, req *AddStaffRequest) (*staff.Response, error)
	// ListSubCenters returns every sub-center of the given center — used to
	// populate the sub-center picker in the "add staff" flow.
	ListSubCenters(ctx context.Context, slug string) (*subcenters.ListResponse, error)
}

type service struct {
	repo              Repository
	staffService      staff.Service
	subCentersService subcenters.Service
}

func NewService(repo Repository, staffService staff.Service, subCentersService subcenters.Service) Service {
	return &service{repo: repo, staffService: staffService, subCentersService: subCentersService}
}

func (s *service) List(ctx context.Context, params ListParams) (*ListResponse, error) {
	list, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx, params.Search)
	if err != nil {
		return nil, err
	}

	ids := make([]uint, 0, len(list))
	for _, c := range list {
		ids = append(ids, c.ID)
	}

	staffCounts, err := s.repo.StaffCounts(ctx, ids)
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

	items := make([]Response, 0, len(list))
	for _, c := range list {
		items = append(items, toResponse(c, staffCounts[c.ID]))
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

// maxSlugAttempts bounds the insert-and-retry loop in Create. Each attempt
// is a real, atomic INSERT — collisions are detected via the DB's unique
// index on slug, not a separate check-then-insert query, so concurrent
// requests for the same name can never both succeed with the same slug.
const maxSlugAttempts = 100

func (s *service) Create(ctx context.Context, req *CreateRequest) (*Response, error) {
	base := utils.Slugify(req.Name, "center")

	for attempt := 1; attempt <= maxSlugAttempts; attempt++ {
		slug := base
		if attempt > 1 {
			slug = fmt.Sprintf("%s-%d", base, attempt)
		}

		center := &auth.Center{Name: req.Name, Slug: slug}
		err := s.repo.Create(ctx, center)
		if err == nil {
			resp := toResponse(*center, 0)
			return &resp, nil
		}
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			continue
		}
		return nil, err
	}

	return nil, ErrSlugGenerationFailed
}

func (s *service) GetDetail(ctx context.Context, slug string) (*DetailResponse, error) {
	center, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if center == nil {
		return nil, ErrCenterNotFound
	}

	staffUsers, err := s.repo.StaffList(ctx, center.ID)
	if err != nil {
		return nil, err
	}

	detail := &DetailResponse{
		ID:        center.ID,
		Name:      center.Name,
		Slug:      center.Slug,
		CreatedAt: center.CreatedAt.Format(time.RFC3339),
		Staff:     make([]StaffMemberResponse, 0, len(staffUsers)),
	}

	if center.Owner != nil {
		detail.Owner = &OwnerResponse{
			ID:        center.Owner.ID,
			Username:  center.Owner.Username,
			Email:     center.Owner.Email,
			CreatedAt: center.Owner.CreatedAt.Format(time.RFC3339),
		}
	}

	for _, u := range staffUsers {
		detail.Staff = append(detail.Staff, StaffMemberResponse{
			ID:        u.ID,
			Username:  u.Username,
			Email:     u.Email,
			Role:      u.Role.Name,
			CreatedAt: u.CreatedAt.Format(time.RFC3339),
		})
	}

	return detail, nil
}

func (s *service) AddOwner(ctx context.Context, slug string, req *CreateOwnerRequest) (*OwnerResponse, error) {
	center, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if center == nil {
		return nil, ErrCenterNotFound
	}
	if center.OwnerID != nil {
		return nil, ErrOwnerAlreadyAssigned
	}

	taken, err := s.repo.EmailTaken(ctx, req.Email)
	if err != nil {
		return nil, err
	}
	if taken {
		return nil, ErrEmailTaken
	}

	roleID, err := s.repo.FindRoleIDByName(ctx, ownerRoleName)
	if err != nil {
		return nil, err
	}

	hashed, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	centerID := center.ID
	user := &auth.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashed,
		RoleID:   roleID,
		CenterID: &centerID,
	}
	if err := s.repo.AssignOwner(ctx, centerID, user); err != nil {
		return nil, err
	}

	return &OwnerResponse{
		ID:        user.ID,
		Username:  user.Username,
		Email:     user.Email,
		CreatedAt: user.CreatedAt.Format(time.RFC3339),
	}, nil
}

func (s *service) AddStaff(ctx context.Context, slug string, actorUserID uint, req *AddStaffRequest) (*staff.Response, error) {
	center, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if center == nil {
		return nil, ErrCenterNotFound
	}

	return s.staffService.CreateForCenter(ctx, center.ID, actorUserID, &staff.CreateRequest{
		Username:    req.Username,
		Email:       req.Email,
		Password:    req.Password,
		Role:        req.Role,
		CenterID:    center.ID,
		SubCenterID: req.SubCenterID,
	})
}

func (s *service) ListSubCenters(ctx context.Context, slug string) (*subcenters.ListResponse, error) {
	center, err := s.repo.FindBySlug(ctx, slug)
	if err != nil {
		return nil, err
	}
	if center == nil {
		return nil, ErrCenterNotFound
	}

	return s.subCentersService.ListForCenter(ctx, center.ID)
}

func toResponse(c auth.Center, staffCount int64) Response {
	resp := Response{
		ID:         c.ID,
		Name:       c.Name,
		Slug:       c.Slug,
		StaffCount: staffCount,
		CreatedAt:  c.CreatedAt.Format(time.RFC3339),
	}

	if c.Owner != nil {
		resp.OwnerUsername = &c.Owner.Username
		resp.OwnerEmail = &c.Owner.Email
	}

	return resp
}
