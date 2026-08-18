package auth

import (
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:50;not null;uniqueIndex"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type User struct {
	ID        uint           `json:"id" gorm:"primaryKey"`
	Username  string         `json:"username" gorm:"size:100;not null"`
	Email     string         `json:"email" gorm:"size:150;not null;uniqueIndex"`
	Password  string         `json:"-" gorm:"size:255;not null"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
	Tokens    []Token        `json:"-" gorm:"foreignKey:UserID"`

	// Every user has exactly one fixed role.
	RoleID uint `json:"role_id" gorm:"not null;index"`
	Role   Role `json:"-" gorm:"foreignKey:RoleID"`

	// Every center-scoped user (owner/scanner/receptionist) belongs to exactly
	// one center. Null for super_admin.
	CenterID *uint   `json:"center_id" gorm:"index"`
	Center   *Center `json:"-" gorm:"foreignKey:CenterID"`

	// A center_scanner/center_receptionist belongs to exactly one sub-center
	// within their center — which operational location they're assigned to.
	// Null for super_admin and center_owner (and for staff created before
	// this field existed).
	SubCenterID *uint      `json:"sub_center_id" gorm:"index"`
	SubCenter   *SubCenter `json:"-" gorm:"foreignKey:SubCenterID"`
}

type Token struct {
	ID uint `json:"id" gorm:"primaryKey"`
	// This index helps the Middleware
	Token  string `json:"token" gorm:"size:255;not null;index:idx_token_lookup"`
	UserID uint   `json:"user_id" gorm:"not null"`
	Type   string `json:"type" gorm:"size:50;not null;index:idx_token_lookup"`

	// ADD THIS: Standalone index for the DeleteExpiredTokens query
	ExpiresAt time.Time `json:"expires_at" gorm:"not null;index:idx_token_lookup;index:idx_expires_at"`

	CreatedAt time.Time `json:"created_at"`
}

type Center struct {
	ID   uint   `json:"id" gorm:"primaryKey"`
	Name string `json:"name" gorm:"size:150;not null"`
	Slug string `json:"slug" gorm:"size:150;not null;uniqueIndex"`

	// A center has at most one owner, and a user can own at most one center.
	// The FK constraint is intentionally disabled (constraint:-): Center and
	// User reference each other (Center.OwnerID -> User, User.CenterID ->
	// Center), and GORM's AutoMigrate cannot create two tables that both
	// require the other to exist first. The uniqueIndex still gives us the
	// "one owner per center" guarantee at the DB level; the owner_id -> user
	// relationship itself is enforced in application code (the seeder).
	OwnerID *uint `json:"owner_id" gorm:"uniqueIndex"`
	Owner   *User `json:"-" gorm:"foreignKey:OwnerID;constraint:-"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

// SubCenter is an operational location within a center (e.g. a branch or
// annex) — a center has many, a sub-center belongs to exactly one center.
// It has no owner of its own: the center's owner (Center.OwnerID) manages
// every one of its sub-centers.
type SubCenter struct {
	ID       uint    `json:"id" gorm:"primaryKey"`
	CenterID uint    `json:"center_id" gorm:"not null;uniqueIndex:idx_subcenter_center_name"`
	Center   *Center `json:"-" gorm:"foreignKey:CenterID;constraint:OnDelete:CASCADE"`
	Name     string  `json:"name" gorm:"size:150;not null;uniqueIndex:idx_subcenter_center_name"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
