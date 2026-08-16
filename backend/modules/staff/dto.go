package staff

const (
	RoleScanner      = "center_scanner"
	RoleReceptionist = "center_receptionist"
)

// CreateRequest creates a center_scanner or center_receptionist user.
// CenterID is required and is validated server-side against the requesting
// owner's own center — it is never trusted for authorization on its own.
type CreateRequest struct {
	Username string `json:"username" binding:"required,min=3,max=100"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Password string `json:"password" binding:"required,min=6,max=72"`
	Role     string `json:"role" binding:"required,oneof=center_scanner center_receptionist"`
	CenterID uint   `json:"center_id" binding:"required"`
}

// UpdateRequest updates an existing staff member. Password is optional — a
// nil/omitted value leaves the current password unchanged.
type UpdateRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=100"`
	Email    string  `json:"email" binding:"required,email,max=150"`
	Role     string  `json:"role" binding:"required,oneof=center_scanner center_receptionist"`
	Password *string `json:"password" binding:"omitempty,min=6,max=72"`
	CenterID uint    `json:"center_id" binding:"required"`
}

type Response struct {
	ID        uint   `json:"id"`
	Username  string `json:"username"`
	Email     string `json:"email"`
	Role      string `json:"role"`
	CenterID  uint   `json:"center_id"`
	CreatedAt string `json:"created_at"`
}

var AllowedPerPage = map[int]struct{}{
	10: {},
	20: {},
	30: {},
	40: {},
	50: {},
}

var AllowedSort = map[string]struct{}{
	"username":   {},
	"email":      {},
	"role":       {},
	"created_at": {},
}

type ListParams struct {
	Page      int
	PerPage   int
	Role      string
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
