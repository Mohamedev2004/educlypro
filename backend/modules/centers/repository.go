package centers

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"educlypro/modules/auth"
	"educlypro/modules/staff"

	"gorm.io/gorm"
)

type Repository interface {
	List(ctx context.Context, params ListParams) ([]auth.Center, error)
	Count(ctx context.Context, search string) (int64, error)
	StaffCounts(ctx context.Context, centerIDs []uint) (map[uint]int64, error)

	// FindBySlug is how a single center is looked up everywhere a specific
	// center is addressed via URL — the numeric id is never exposed there.
	FindBySlug(ctx context.Context, slug string) (*auth.Center, error)
	// Create inserts the center as-is. Slug uniqueness is NOT pre-checked
	// here — the caller (Service.Create) is expected to retry with a new
	// candidate slug when this returns a gorm.ErrDuplicatedKey, so the
	// actual DB unique index is the single source of truth (no separate
	// check-then-insert race window).
	Create(ctx context.Context, center *auth.Center) error

	EmailTaken(ctx context.Context, email string) (bool, error)
	FindRoleIDByName(ctx context.Context, name string) (uint, error)
	// AssignOwner creates the given user and sets it as the center's owner
	// atomically, failing if the center already has one.
	AssignOwner(ctx context.Context, centerID uint, user *auth.User) error
	// StaffList returns every center_scanner/center_receptionist user at the
	// given center, newest first.
	StaffList(ctx context.Context, centerID uint) ([]auth.User, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) baseQuery(search string) *gorm.DB {
	query := r.db.Model(&auth.Center{})

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(slug) LIKE ?", like, like)
	}

	return query
}

func (r *repository) List(ctx context.Context, params ListParams) ([]auth.Center, error) {
	var list []auth.Center

	orderColumn := "created_at"
	if params.Sort != "" {
		orderColumn = params.Sort
	}

	direction := "DESC"
	if strings.EqualFold(params.Direction, "asc") {
		direction = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage

	err := r.baseQuery(params.Search).
		WithContext(ctx).
		Preload("Owner").
		Order(fmt.Sprintf("%s %s", orderColumn, direction)).
		Limit(params.PerPage).
		Offset(offset).
		Find(&list).Error

	return list, err
}

func (r *repository) Count(ctx context.Context, search string) (int64, error) {
	var count int64
	err := r.baseQuery(search).WithContext(ctx).Count(&count).Error
	return count, err
}

// StaffCounts returns the number of center_scanner + center_receptionist
// users per center, for the given center ids only.
func (r *repository) StaffCounts(ctx context.Context, centerIDs []uint) (map[uint]int64, error) {
	counts := make(map[uint]int64, len(centerIDs))
	if len(centerIDs) == 0 {
		return counts, nil
	}

	var rows []struct {
		CenterID uint
		Count    int64
	}

	err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Select("center_id, COUNT(*) as count").
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.center_id IN ?", centerIDs).
		Where("roles.name IN ?", []string{staff.RoleScanner, staff.RoleReceptionist}).
		Group("center_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.CenterID] = row.Count
	}
	return counts, nil
}

func (r *repository) FindBySlug(ctx context.Context, slug string) (*auth.Center, error) {
	var center auth.Center
	err := r.db.WithContext(ctx).Preload("Owner").Where("slug = ?", slug).First(&center).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &center, nil
}

func (r *repository) Create(ctx context.Context, center *auth.Center) error {
	return r.db.WithContext(ctx).Create(center).Error
}

func (r *repository) EmailTaken(ctx context.Context, email string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&auth.User{}).Where("email = ?", email).Count(&count).Error
	return count > 0, err
}

func (r *repository) FindRoleIDByName(ctx context.Context, name string) (uint, error) {
	var role auth.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

func (r *repository) AssignOwner(ctx context.Context, centerID uint, user *auth.User) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(user).Error; err != nil {
			return err
		}

		result := tx.Model(&auth.Center{}).
			Where("id = ? AND owner_id IS NULL", centerID).
			Update("owner_id", user.ID)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return ErrOwnerAlreadyAssigned
		}
		return nil
	})
}

func (r *repository) StaffList(ctx context.Context, centerID uint) ([]auth.User, error) {
	var users []auth.User
	err := r.db.WithContext(ctx).
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.center_id = ?", centerID).
		Where("roles.name IN ?", []string{staff.RoleScanner, staff.RoleReceptionist}).
		Preload("Role").
		Select("users.*").
		Order("users.created_at DESC").
		Find(&users).Error
	return users, err
}
