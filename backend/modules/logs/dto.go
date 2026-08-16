package logs

import "time"

var AllowedPerPage = map[int]struct{}{
	10:  {},
	20:  {},
	30:  {},
	40:  {},
	50:  {},
	100: {},
}

type Facets struct {
	Levels    []string `json:"levels"`
	Statuses  []string `json:"statuses"`
	Durations []string `json:"durations"`
	StatusCodes []int  `json:"status_codes"`
}

type ListParams struct {
	Page    int
	PerPage int

	// Filters
	Query     string
	Levels    []string
	Statuses  []string
	StatusCodes []int
	Durations []string
	From      *time.Time
	To        *time.Time
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasNext    bool  `json:"has_next"`
	HasPrev    bool  `json:"has_prev"`
}

type LogItem struct {
	ID        string `json:"id"`
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Service   string `json:"service"`
	Message   string `json:"message"`
	Duration  string `json:"duration"`
	Status    string `json:"status"`
	StatusCode int   `json:"status_code"`
	// Details for the expanded view (optional usage by UI)
	RequestID string      `json:"request_id,omitempty"`
	ActorID   string      `json:"actor_id,omitempty"`
	EntityID  string      `json:"entity_id,omitempty"`
	Payload   interface{} `json:"payload,omitempty"`
}

type LevelCounts struct {
	Info    int64 `json:"info"`
	Warning int64 `json:"warning"`
	Error   int64 `json:"error"`
}

type ListResponse struct {
	Items      []LogItem      `json:"items"`
	Facets     Facets         `json:"facets"`
	Pagination PaginationMeta `json:"pagination"`
	Applied    map[string]any `json:"applied"`
	Counts     LevelCounts    `json:"counts"`
}

type ChartRange string

const (
	ChartRange24h ChartRange = "24h"
	ChartRange7d  ChartRange = "7d"
)

type ChartPoint struct {
	Date  string `json:"date"`
	Count int64  `json:"count"`
}
