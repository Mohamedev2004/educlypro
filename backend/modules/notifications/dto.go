package notifications

const (
	FilterRead   = "read"
	FilterUnread = "unread"
)

var AllowedPerPage = map[int]struct{}{
	10: {},
	20: {},
	30: {},
	40: {},
	50: {},
}

type ListParams struct {
	Page    int
	PerPage int
	Filter  string
}

type Counts struct {
	All    int64 `json:"all"`
	Read   int64 `json:"read"`
	Unread int64 `json:"unread"`
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
	Items      []Notification `json:"items"`
	Filter     string         `json:"filter"`
	Counts     Counts         `json:"counts"`
	Pagination PaginationMeta `json:"pagination"`
}
