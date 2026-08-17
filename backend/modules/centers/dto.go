package centers

// Response is the super-admin-facing shape of a center: identity, its owner
// (if one exists), and a count of the operational staff (scanners +
// receptionists) at that center.
type Response struct {
	ID            uint    `json:"id"`
	Name          string  `json:"name"`
	Slug          string  `json:"slug"`
	OwnerUsername *string `json:"owner_username"`
	OwnerEmail    *string `json:"owner_email"`
	StaffCount    int64   `json:"staff_count"`
	CreatedAt     string  `json:"created_at"`
}

var AllowedPerPage = map[int]struct{}{
	10: {},
	20: {},
	30: {},
	40: {},
	50: {},
}

var AllowedSort = map[string]struct{}{
	"name":       {},
	"slug":       {},
	"created_at": {},
}

type ListParams struct {
	Page      int
	PerPage   int
	Search    string
	Sort      string
	Direction string
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type ListResponse struct {
	Items      []Response     `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

// CreateRequest creates a new center. The slug is always derived from Name
// server-side and de-duplicated automatically — it is never client-supplied.
type CreateRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

// CreateOwnerRequest creates the center_owner account for a center that
// doesn't have one yet.
type CreateOwnerRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=6,max=72"`
}

// AddStaffRequest creates a center_scanner or center_receptionist account
// for a specific center. Unlike staff.CreateRequest, the center is taken
// from the URL (an admin acting on an arbitrary center), not from the
// request body.
type AddStaffRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Role     string `json:"role" binding:"required,oneof=center_scanner center_receptionist"`
}

type OwnerResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	CreatedAt string `json:"created_at"`
}

type StaffMemberResponse struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

// DetailResponse is the full super-admin-facing view of a single center,
// used by the center detail page: identity, owner (if any), and every
// staff member (not just a count, unlike the list Response).
type DetailResponse struct {
	ID        uint                  `json:"id"`
	Name      string                `json:"name"`
	Slug      string                `json:"slug"`
	Owner     *OwnerResponse        `json:"owner"`
	Staff     []StaffMemberResponse `json:"staff"`
	CreatedAt string                `json:"created_at"`
}
