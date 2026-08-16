package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"
)

type Service interface {
	List(ctx context.Context, params ListParams) (*ListResponse, error)
	Export(ctx context.Context, params ListParams, w io.Writer) error
	Chart(ctx context.Context, r ChartRange) ([]ChartPoint, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func normalizeLevel(s string) string {
	upper := strings.ToUpper(strings.TrimSpace(s))
	switch upper {
	case "WARN", "WARNING":
		return "WARN"
	case "INFO":
		return "INFO"
	case "ERROR":
		return "ERROR"
	case "DEBUG":
		return "DEBUG"
	default:
		return upper
	}
}

func levelForUI(dbLevel string) string {
	switch strings.ToUpper(dbLevel) {
	case "WARN", "WARNING":
		return "warning"
	case "ERROR":
		return "error"
	case "DEBUG":
		return "info"
	case "INFO":
		return "info"
	default:
		return strings.ToLower(dbLevel)
	}
}

func durationForUI(ms float64) string {
	if ms <= 0 {
		return "0ms"
	}
	if ms < 1000 {
		return fmt.Sprintf("%dms", int(math.Round(ms)))
	}
	seconds := ms / 1000
	return fmt.Sprintf("%.1fs", seconds)
}

func statusForUI(action string) string {
	switch strings.ToUpper(strings.TrimSpace(action)) {
	case "CREATED":
		return "Created"
	case "UPDATED":
		return "Updated"
	case "FAILED":
		return "Failed"
	case "DELETED":
		return "Deleted"
	default:
		if action == "" {
			return ""
		}
		return strings.ToUpper(action)
	}
}

func messageForUI(l AuditLog) string {
	// Prefer a message field inside payload when available.
	var payload map[string]any
	if len(l.Payload) > 0 {
		_ = json.Unmarshal(l.Payload, &payload)
	}
	if payload != nil {
		if v, ok := payload["message"].(string); ok && strings.TrimSpace(v) != "" {
			return v
		}
		if v, ok := payload["error_message"].(string); ok && strings.TrimSpace(v) != "" {
			return v
		}
	}

	entity := l.Entity
	if strings.TrimSpace(entity) == "" {
		entity = "system"
	}
	action := l.Action
	if strings.TrimSpace(action) == "" {
		action = "EVENT"
	}

	if strings.TrimSpace(l.EntityID) != "" {
		return fmt.Sprintf("%s %s (%s)", action, entity, l.EntityID)
	}
	return fmt.Sprintf("%s %s", action, entity)
}

func (s *service) List(ctx context.Context, params ListParams) (*ListResponse, error) {
	items, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	total, err := s.repo.Count(ctx, params)
	if err != nil {
		return nil, err
	}

	facets, err := s.repo.DistinctValues(ctx)
	if err != nil {
		return nil, err
	}

	totalPages := int(total) / params.PerPage
	if int(total)%params.PerPage != 0 {
		totalPages++
	}
	if totalPages == 0 {
		totalPages = 1
	}

	out := make([]LogItem, 0, len(items))
	for _, l := range items {
		var payload any
		if len(l.Payload) > 0 {
			_ = json.Unmarshal(l.Payload, &payload)
		}

		out = append(out, LogItem{
			ID:         strconv.FormatUint(uint64(l.ID), 10),
			Timestamp:  l.CreatedAt.UTC().Format("2006-01-02T15:04:05Z07:00"),
			Level:      levelForUI(l.Level),
			Service:    l.Entity,
			Message:    messageForUI(l),
			Duration:   durationForUI(l.DurationMs),
			Status:     statusForUI(l.Action),
			StatusCode: l.StatusCode,
			RequestID:  l.RequestID,
			ActorID:    l.ActorID,
			EntityID:   l.EntityID,
			Payload:    payload,
		})
	}

	applied := map[string]any{
		"q":           params.Query,
		"level":       params.Levels,
		"status":      params.Statuses,
		"status_code": params.StatusCodes,
		"duration":    params.Durations,
		"from":        params.From,
		"to":          params.To,
	}

	counts, err := s.repo.CountByLevel(ctx, params)
	if err != nil {
		return nil, err
	}

	return &ListResponse{
		Items: out,
		Facets: Facets{
			Levels:      facets.Levels,
			Statuses:    facets.Statuses,
			Durations:   facets.Durations,
			StatusCodes: facets.StatusCodes,
		},
		Pagination: PaginationMeta{
			Page:       params.Page,
			PerPage:    params.PerPage,
			Total:      total,
			TotalPages: totalPages,
			HasNext:    params.Page < totalPages,
			HasPrev:    params.Page > 1,
		},
		Applied: applied,
		Counts:  counts,
	}, nil
}

func (s *service) Chart(ctx context.Context, r ChartRange) ([]ChartPoint, error) {
	if r != ChartRange24h && r != ChartRange7d {
		r = ChartRange24h
	}
	return s.repo.Chart(ctx, r)
}
