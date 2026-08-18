package teachers

// SubjectRef is the minimal shape of a subject a teacher is assigned to.
type SubjectRef struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// ClassRef is the minimal shape of a class a teacher is assigned to.
type ClassRef struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

type Response struct {
	ID        uint         `json:"id"`
	Slug      string       `json:"slug"`
	FullName  string       `json:"full_name"`
	Email     string       `json:"email"`
	Phone     string       `json:"phone"`
	Subjects  []SubjectRef `json:"subjects"`
	Classes   []ClassRef   `json:"classes"`
	CreatedAt string       `json:"created_at"`
}

var AllowedPerPage = map[int]struct{}{
	10: {},
	20: {},
	30: {},
	40: {},
	50: {},
}

var AllowedSort = map[string]struct{}{
	"full_name":  {},
	"email":      {},
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

// CreateRequest creates a teacher in the caller's own center with just its
// core identity fields. Subject/class assignment happens afterward, from
// the teacher's detail page (see UpdateSubjectsRequest/UpdateClassesRequest).
type CreateRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=150"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Phone    string `json:"phone" binding:"required,min=6,max=30"`
}

// UpdateRequest replaces an existing teacher's core identity fields. Slug is
// intentionally not re-derived here — it's fixed at creation time.
type UpdateRequest struct {
	FullName string `json:"full_name" binding:"required,min=2,max=150"`
	Email    string `json:"email" binding:"required,email,max=150"`
	Phone    string `json:"phone" binding:"required,min=6,max=30"`
}

// UpdateSubjectsRequest replaces a teacher's full subject assignment
// (SubjectIDs is the new complete set, not a delta). Used by the teacher
// detail page's subjects editor.
type UpdateSubjectsRequest struct {
	SubjectIDs []uint `json:"subject_ids"`
}

// UpdateClassesRequest is UpdateSubjectsRequest's counterpart for classes.
type UpdateClassesRequest struct {
	ClassIDs []uint `json:"class_ids"`
}
