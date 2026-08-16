package notifications

import (
	"context"
	"log"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, n *Notification) error
	ListForUser(ctx context.Context, userID uint, filter string, limit, offset int) ([]Notification, error)
	CountForUser(ctx context.Context, userID uint, filter string) (int64, error)
	CountSummary(ctx context.Context, userID uint) (Counts, error)
	MarkRead(ctx context.Context, userID uint, id uint) error
	MarkAllRead(ctx context.Context, userID uint) error
	UnreadCount(ctx context.Context, userID uint) (int64, error)
	Delete(ctx context.Context, userID uint, id uint) error
	DeleteAllRead(ctx context.Context, userID uint) error
	DeleteOldReadNotifications() error
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, n *Notification) error {
	return r.db.WithContext(ctx).Create(n).Error
}

func (r *repository) ListForUser(ctx context.Context, userID uint, filter string, limit, offset int) ([]Notification, error) {
	var results []Notification

	query := r.db.WithContext(ctx).Where("user_id = ?", userID)
	if filter == FilterRead {
		query = query.Where("is_read = ?", true)
	}
	if filter == FilterUnread {
		query = query.Where("is_read = ?", false)
	}

	err := query.
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&results).Error
	return results, err
}

func (r *repository) CountForUser(ctx context.Context, userID uint, filter string) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&Notification{}).Where("user_id = ?", userID)
	if filter == FilterRead {
		query = query.Where("is_read = ?", true)
	}
	if filter == FilterUnread {
		query = query.Where("is_read = ?", false)
	}

	err := query.Count(&count).Error
	return count, err
}

func (r *repository) CountSummary(ctx context.Context, userID uint) (Counts, error) {
	var (
		allCount    int64
		unreadCount int64
	)

	if err := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ?", userID).
		Count(&allCount).Error; err != nil {
		return Counts{}, err
	}

	if err := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&unreadCount).Error; err != nil {
		return Counts{}, err
	}

	return Counts{
		All:    allCount,
		Read:   allCount - unreadCount,
		Unread: unreadCount,
	}, nil
}

func (r *repository) MarkRead(ctx context.Context, userID uint, id uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *repository) MarkAllRead(ctx context.Context, userID uint) error {
	now := time.Now()
	return r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Updates(map[string]any{
			"is_read": true,
			"read_at": &now,
		}).Error
}

func (r *repository) UnreadCount(ctx context.Context, userID uint) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error
	return count, err
}

func (r *repository) Delete(ctx context.Context, userID uint, id uint) error {
	return r.db.WithContext(ctx).
		Model(&Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Delete(&Notification{}).Error
}

func (r *repository) DeleteAllRead(ctx context.Context, userID uint) error {
	return r.db.WithContext(ctx).
		Where("user_id = ? AND is_read = ?", userID, true).
		Delete(&Notification{}).Error
}

func (r *repository) DeleteOldReadNotifications() error {
	cutoff := time.Now().AddDate(0, -1, 0)

	result := r.db.
		Where("is_read = ? AND created_at < ?", true, cutoff).
		Delete(&Notification{})

	if result.Error != nil {
		return result.Error
	}

	log.Printf("Schedule [notification-cleanup] soft deleted %d read notifications older than %s",
		result.RowsAffected, cutoff.Format(time.DateOnly))

	return nil
}
