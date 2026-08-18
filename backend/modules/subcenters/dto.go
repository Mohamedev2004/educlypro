package subcenters

// Response is the shape of a single sub-center: identity and how many staff
// are currently assigned to it.
type Response struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	StaffCount int64  `json:"staff_count"`
	CreatedAt  string `json:"created_at"`
}

// ListResponse is unpaginated — a center realistically has a handful of
// sub-centers, not pages of them (unlike staff or centers).
type ListResponse struct {
	Items []Response `json:"items"`
}

// CreateRequest creates a sub-center under the caller's own center.
type CreateRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}

// UpdateRequest renames an existing sub-center.
type UpdateRequest struct {
	Name string `json:"name" binding:"required,min=2,max=150"`
}
