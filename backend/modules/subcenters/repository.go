package subcenters

import (
	"context"
	"errors"
	"strings"

	"educlypro/modules/auth"

	"gorm.io/gorm"
)

type Repository interface {
	// FindOwnerCenterID resolves the center a center_owner user belongs to.
	// Mirrors staff.Repository.FindOwnerCenterID.
	FindOwnerCenterID(ctx context.Context, ownerUserID uint) (uint, error)

	ListByCenter(ctx context.Context, centerID uint) ([]auth.SubCenter, error)
	// StaffCounts returns how many staff (center_scanner + center_receptionist)
	// are assigned to each of the given sub-centers.
	StaffCounts(ctx context.Context, subCenterIDs []uint) (map[uint]int64, error)

	NameExists(ctx context.Context, centerID uint, name string, excludeID uint) (bool, error)
	Create(ctx context.Context, subCenter *auth.SubCenter) error
	FindByIDInCenter(ctx context.Context, centerID, subCenterID uint) (*auth.SubCenter, error)
	Update(ctx context.Context, subCenter *auth.SubCenter) error
	Delete(ctx context.Context, subCenterID uint) error
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

func (r *repository) ListByCenter(ctx context.Context, centerID uint) ([]auth.SubCenter, error) {
	var list []auth.SubCenter
	err := r.db.WithContext(ctx).Where("center_id = ?", centerID).Order("id ASC").Find(&list).Error
	return list, err
}

func (r *repository) StaffCounts(ctx context.Context, subCenterIDs []uint) (map[uint]int64, error) {
	counts := make(map[uint]int64, len(subCenterIDs))
	if len(subCenterIDs) == 0 {
		return counts, nil
	}

	var rows []struct {
		SubCenterID uint
		Count       int64
	}

	err := r.db.WithContext(ctx).
		Model(&auth.User{}).
		Select("sub_center_id, COUNT(*) as count").
		Where("sub_center_id IN ?", subCenterIDs).
		Group("sub_center_id").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	for _, row := range rows {
		counts[row.SubCenterID] = row.Count
	}
	return counts, nil
}

func (r *repository) NameExists(ctx context.Context, centerID uint, name string, excludeID uint) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&auth.SubCenter{}).
		Where("center_id = ? AND LOWER(name) = LOWER(?) AND id <> ?", centerID, strings.TrimSpace(name), excludeID).
		Count(&count).Error
	return count > 0, err
}

func (r *repository) Create(ctx context.Context, subCenter *auth.SubCenter) error {
	return r.db.WithContext(ctx).Create(subCenter).Error
}

func (r *repository) FindByIDInCenter(ctx context.Context, centerID, subCenterID uint) (*auth.SubCenter, error) {
	var sc auth.SubCenter
	err := r.db.WithContext(ctx).
		Where("id = ? AND center_id = ?", subCenterID, centerID).
		First(&sc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &sc, nil
}

func (r *repository) Update(ctx context.Context, subCenter *auth.SubCenter) error {
	return r.db.WithContext(ctx).
		Model(&auth.SubCenter{}).
		Where("id = ?", subCenter.ID).
		Update("name", subCenter.Name).Error
}

func (r *repository) Delete(ctx context.Context, subCenterID uint) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Soft-deleted staff rows keep their sub_center_id — soft delete
		// only sets deleted_at, it doesn't clear foreign keys — so without
		// this they'd still trip the FK constraint below even though the
		// service's staff-count check (which correctly excludes
		// soft-deleted staff) already confirmed no *active* staff remain.
		// Unscoped so soft-deleted rows are included.
		if err := tx.Unscoped().Model(&auth.User{}).
			Where("sub_center_id = ?", subCenterID).
			Update("sub_center_id", nil).Error; err != nil {
			return err
		}

		return tx.Delete(&auth.SubCenter{}, subCenterID).Error
	})
}
