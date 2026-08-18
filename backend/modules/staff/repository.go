package staff

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"educlypro/modules/auth"

	"gorm.io/gorm"
)

type Repository interface {
	FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error)
	ListByCenter(ctx context.Context, centerID uint, params ListParams) ([]auth.User, error)
	CountByCenter(ctx context.Context, centerID uint, role, search string) (int64, error)
	FindByIDInCenter(ctx context.Context, centerID, staffID uint) (*auth.User, error)
	FindRoleIDByName(ctx context.Context, name string) (uint, error)
	// FindSubCenter looks up a sub-center by id, regardless of which center
	// it belongs to — the caller is responsible for checking CenterID
	// matches the center the staff member is being assigned into.
	FindSubCenter(ctx context.Context, subCenterID uint) (*auth.SubCenter, error)
	Create(ctx context.Context, user *auth.User) error
	Update(ctx context.Context, user *auth.User) error
	SoftDelete(ctx context.Context, staffID uint) error
	EmailTaken(ctx context.Context, email string, excludeUserID uint) (bool, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error) {
	var user auth.User
	if err := r.db.WithContext(ctx).Select("center_id").First(&user, ownerUserID).Error; err != nil {
		return 0, err
	}
	if user.CenterID == nil {
		return 0, errors.New("owner has no center")
	}
	return *user.CenterID, nil
}

func (r *repository) staffQuery(centerID uint, role, search string) *gorm.DB {
	query := r.db.
		Model(&auth.User{}).
		Joins("JOIN roles ON roles.id = users.role_id").
		Where("users.center_id = ?", centerID).
		Where("roles.name IN ?", []string{RoleScanner, RoleReceptionist})

	if role != "" {
		query = query.Where("roles.name = ?", role)
	}

	if search != "" {
		like := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(users.username) LIKE ? OR LOWER(users.email) LIKE ?", like, like)
	}

	return query
}

func (r *repository) ListByCenter(ctx context.Context, centerID uint, params ListParams) ([]auth.User, error) {
	var users []auth.User

	orderColumn := "users.created_at"
	if params.Sort == "role" {
		orderColumn = "roles.name"
	} else if params.Sort != "" {
		orderColumn = "users." + params.Sort
	}

	direction := "DESC"
	if strings.EqualFold(params.Direction, "asc") {
		direction = "ASC"
	}

	offset := (params.Page - 1) * params.PerPage

	err := r.staffQuery(centerID, params.Role, params.Search).
		WithContext(ctx).
		Preload("Role").
		Preload("SubCenter").
		Select("users.*").
		Order(fmt.Sprintf("%s %s", orderColumn, direction)).
		Limit(params.PerPage).
		Offset(offset).
		Find(&users).Error

	return users, err
}

func (r *repository) CountByCenter(ctx context.Context, centerID uint, role, search string) (int64, error) {
	var count int64
	err := r.staffQuery(centerID, role, search).WithContext(ctx).Count(&count).Error
	return count, err
}

func (r *repository) FindByIDInCenter(ctx context.Context, centerID, staffID uint) (*auth.User, error) {
	var user auth.User
	err := r.db.WithContext(ctx).
		Preload("Role").
		Preload("SubCenter").
		Where("id = ? AND center_id = ?", staffID, centerID).
		First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *repository) FindRoleIDByName(ctx context.Context, name string) (uint, error) {
	var role auth.Role
	if err := r.db.WithContext(ctx).Where("name = ?", name).First(&role).Error; err != nil {
		return 0, err
	}
	return role.ID, nil
}

func (r *repository) FindSubCenter(ctx context.Context, subCenterID uint) (*auth.SubCenter, error) {
	var sc auth.SubCenter
	err := r.db.WithContext(ctx).First(&sc, subCenterID).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sc, nil
}

func (r *repository) Create(ctx context.Context, user *auth.User) error {
	return r.db.WithContext(ctx).Create(user).Error
}

func (r *repository) Update(ctx context.Context, user *auth.User) error {
	return r.db.WithContext(ctx).Model(&auth.User{}).
		Where("id = ?", user.ID).
		Updates(map[string]any{
			"username":      user.Username,
			"email":         user.Email,
			"role_id":       user.RoleID,
			"password":      user.Password,
			"sub_center_id": user.SubCenterID,
		}).Error
}

func (r *repository) SoftDelete(ctx context.Context, staffID uint) error {
	return r.db.WithContext(ctx).Delete(&auth.User{}, staffID).Error
}

func (r *repository) EmailTaken(ctx context.Context, email string, excludeUserID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Where("email = ? AND id <> ?", email, excludeUserID).
		Count(&count).Error
	return count > 0, err
}
