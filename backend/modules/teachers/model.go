package teachers

import (
	"time"

	"educlypro/modules/academic"

	"gorm.io/gorm"
)

// Teacher is a plain data record, not an auth.User — teachers never sign in.
// It is scoped to a single center (like Grade/SubCenter) and is many-to-many
// with both the subjects it teaches and the classes it's assigned to.
type Teacher struct {
	ID uint `json:"id" gorm:"primaryKey"`
	// CenterID also anchors the composite unique index on Slug below —
	// GORM keeps it usable for plain center_id lookups too, since it's the
	// index's leading column.
	CenterID uint   `json:"center_id" gorm:"not null;uniqueIndex:idx_teacher_center_slug"`
	FullName string `json:"full_name" gorm:"size:150;not null"`
	// Slug is derived from FullName at creation time and never changes
	// afterward, so links to a teacher's detail page stay stable even if
	// they're later renamed. Unique per center (see CenterID comment), not
	// globally — see Email below for why per-center is correct here.
	Slug string `json:"slug" gorm:"size:170;not null;uniqueIndex:idx_teacher_center_slug"`
	// Unique per center, not globally — teachers aren't auth principals, so
	// two different centers can each have their own teacher at the same
	// email without conflict.
	Email string `json:"email" gorm:"size:150;not null;uniqueIndex:idx_teacher_center_email"`
	Phone string `json:"phone" gorm:"size:30;not null"`

	// OnDelete:CASCADE on both sides of each join table: deleting a teacher
	// drops their subject/class links, and deleting a subject or class (a
	// normal action, e.g. from the onboarding wizard) drops it from any
	// teacher's list instead of being blocked by the join table's FK — the
	// same class of bug already fixed for sub-centers, avoided here from
	// the start.
	Subjects []academic.Subject `json:"subjects,omitempty" gorm:"many2many:teacher_subjects;constraint:OnDelete:CASCADE"`
	Classes  []academic.Class   `json:"classes,omitempty" gorm:"many2many:teacher_classes;constraint:OnDelete:CASCADE"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}
