package logs

import (
	"context"
	"strings"
	"time"

	"gorm.io/gorm"
)

type Repository interface {
	Create(ctx context.Context, log *AuditLog) error
	List(ctx context.Context, params ListParams) ([]AuditLog, error)
	Count(ctx context.Context, params ListParams) (int64, error)
	CountByLevel(ctx context.Context, params ListParams) (LevelCounts, error)
	DistinctValues(ctx context.Context) (Facets, error)
	ListAll(ctx context.Context, params ListParams, fn func(AuditLog) error) error
	Chart(ctx context.Context, r ChartRange) ([]ChartPoint, error)
}

type repository struct {
	db *gorm.DB
}

// NewRepository expects the secondary Logs Database connection!
func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, log *AuditLog) error {
	return r.db.WithContext(ctx).Create(log).Error
}

func (r *repository) baseQuery(ctx context.Context, params ListParams) *gorm.DB {
	q := r.db.WithContext(ctx).Model(&AuditLog{})

	if len(params.Levels) > 0 {
		q = q.Where("level IN ?", params.Levels)
	}
	if len(params.Statuses) > 0 {
		q = q.Where("action IN ?", params.Statuses)
	}
	if len(params.StatusCodes) > 0 {
		q = q.Where("status_code IN ?", params.StatusCodes)
	}

	// Duration buckets (ms)
	// If both buckets are selected, do not filter.
	if len(params.Durations) == 1 {
		switch params.Durations[0] {
		case "bigger than 200ms":
			q = q.Where("duration_ms > ?", 200)
		case "less than 200ms":
			q = q.Where("duration_ms < ?", 200)
		}
	}

	if params.From != nil {
		q = q.Where("created_at >= ?", params.From.UTC())
	}
	if params.To != nil {
		q = q.Where("created_at <= ?", params.To.UTC())
	}

	if params.Query != "" {
		like := "%" + strings.ToLower(params.Query) + "%"
		// Note: SQLite uses LIKE case-insensitively by default, Postgres needs LOWER().
		q = q.Where(
			`LOWER(request_id) LIKE ? OR LOWER(entity) LIKE ? OR LOWER(entity_id) LIKE ? OR LOWER(actor_id) LIKE ? OR LOWER(action) LIKE ? OR LOWER(CAST(payload AS TEXT)) LIKE ?`,
			like, like, like, like, like, like,
		)
	}

	return q
}

func (r *repository) List(ctx context.Context, params ListParams) ([]AuditLog, error) {
	var items []AuditLog
	offset := (params.Page - 1) * params.PerPage

	q := r.baseQuery(ctx, params).
		Order("created_at DESC").
		Limit(params.PerPage).
		Offset(offset)

	if err := q.Find(&items).Error; err != nil {
		return nil, err
	}
	return items, nil
}

func (r *repository) Count(ctx context.Context, params ListParams) (int64, error) {
	var total int64
	if err := r.baseQuery(ctx, params).Count(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *repository) Chart(ctx context.Context, rg ChartRange) ([]ChartPoint, error) {
	now := time.Now().UTC()
	steps := 24
	stepDur := time.Hour
	if rg == ChartRange7d {
		steps = 7
		stepDur = 24 * time.Hour
	}

	start := now.Add(-time.Duration(steps) * stepDur)
	if rg == ChartRange7d {
		start = time.Date(start.Year(), start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
	}

	type row struct {
		Bucket time.Time
		Count  int64
	}
	var rows []row

	dialect := r.db.Dialector.Name()
	q := r.db.WithContext(ctx).Model(&AuditLog{}).Where("created_at >= ? AND created_at <= ?", start, now)

	switch dialect {
	case "postgres":
		if rg == ChartRange24h {
			q = q.Select("date_trunc('hour', created_at) AS bucket, COUNT(*) AS count").Group("bucket")
		} else {
			q = q.Select("date_trunc('day', created_at) AS bucket, COUNT(*) AS count").Group("bucket")
		}
		if err := q.Order("bucket ASC").Scan(&rows).Error; err != nil {
			return nil, err
		}
	case "sqlite":
		// SQLite stores timestamps as text; use strftime to bucket in UTC.
		if rg == ChartRange24h {
			q = q.Select("strftime('%Y-%m-%dT%H:00:00Z', created_at) AS bucket, COUNT(*) AS count").Group("bucket")
		} else {
			q = q.Select("strftime('%Y-%m-%dT00:00:00Z', created_at) AS bucket, COUNT(*) AS count").Group("bucket")
		}
		type rowStr struct {
			Bucket string
			Count  int64
		}
		var r2 []rowStr
		if err := q.Order("bucket ASC").Scan(&r2).Error; err != nil {
			return nil, err
		}
		for _, rr := range r2 {
			if t, err := time.Parse(time.RFC3339, rr.Bucket); err == nil {
				rows = append(rows, row{Bucket: t, Count: rr.Count})
			}
		}
	default:
		// Generic fallback: pull timestamps only (bounded by range) and bucket in Go.
		var times []time.Time
		if err := q.Select("created_at").Pluck("created_at", &times).Error; err != nil {
			return nil, err
		}
		m := make(map[time.Time]int64, steps)
		for _, ts := range times {
			t := ts.UTC()
			var b time.Time
			if rg == ChartRange24h {
				b = time.Date(t.Year(), t.Month(), t.Day(), t.Hour(), 0, 0, 0, time.UTC)
			} else {
				b = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC)
			}
			m[b]++
		}
		for b, c := range m {
			rows = append(rows, row{Bucket: b, Count: c})
		}
	}

	countByBucket := make(map[time.Time]int64, len(rows))
	for _, rr := range rows {
		b := rr.Bucket.UTC()
		if rg == ChartRange24h {
			b = time.Date(b.Year(), b.Month(), b.Day(), b.Hour(), 0, 0, 0, time.UTC)
		} else {
			b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)
		}
		countByBucket[b] = rr.Count
	}

	out := make([]ChartPoint, 0, steps)
	for i := steps - 1; i >= 0; i-- {
		b := now.Add(-time.Duration(i) * stepDur)
		if rg == ChartRange24h {
			b = time.Date(b.Year(), b.Month(), b.Day(), b.Hour(), 0, 0, 0, time.UTC)
		} else {
			b = time.Date(b.Year(), b.Month(), b.Day(), 0, 0, 0, 0, time.UTC)
		}
		out = append(out, ChartPoint{
			Date:  b.Format(time.RFC3339),
			Count: countByBucket[b],
		})
	}

	return out, nil
}

func (r *repository) DistinctValues(ctx context.Context) (Facets, error) {
	out := Facets{
		Levels:      []string{"Info", "Warning", "Error"},
		Statuses:    []string{"Created", "Updated", "Deleted", "Failed"},
		Durations:   []string{"bigger than 200ms", "less than 200ms"},
		StatusCodes: []int{},
	}

	var codes []int
	if err := r.db.WithContext(ctx).
		Model(&AuditLog{}).
		Where("status_code > 0").
		Distinct().
		Order("status_code ASC").
		Pluck("status_code", &codes).Error; err != nil {
		return Facets{}, err
	}
	out.StatusCodes = codes
	return out, nil
}

func (r *repository) ListAll(ctx context.Context, params ListParams, fn func(AuditLog) error) error {
	query := r.db.WithContext(ctx).Model(&AuditLog{})
	query = applyFilters(query, params) // extract your existing filter logic into this helper

	// Use Find with rows() to stream — never loads all into memory
	rows, err := query.Order("created_at DESC").Rows()
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var log AuditLog
		if err := r.db.ScanRows(rows, &log); err != nil {
			return err
		}
		if err := fn(log); err != nil {
			return err
		}
	}
	return rows.Err()
}

func applyFilters(query *gorm.DB, params ListParams) *gorm.DB {
	if params.Query != "" {
		like := "%" + params.Query + "%"
		query = query.Where(
			"entity ILIKE ? OR entity_id ILIKE ? OR actor_id ILIKE ? OR request_id ILIKE ?",
			like, like, like, like,
		)
	}
	if len(params.Levels) > 0 {
		query = query.Where("level IN ?", params.Levels)
	}
	if len(params.Statuses) > 0 {
		query = query.Where("action IN ?", params.Statuses)
	}
	if params.From != nil {
		query = query.Where("created_at >= ?", params.From)
	}
	if params.To != nil {
		query = query.Where("created_at <= ?", params.To)
	}
	return query
}

// CountByLevel returns the count of logs for each level (Info, Warning, Error) based on the provided filters.
func (r *repository) CountByLevel(ctx context.Context, params ListParams) (LevelCounts, error) {
	type row struct {
		Level string
		Count int64
	}
	var rows []row

	// Use baseQuery but without level filter so all 3 buckets are always returned
	noLevelParams := params
	noLevelParams.Levels = nil

	err := r.baseQuery(ctx, noLevelParams).
		Select("level, COUNT(*) as count").
		Group("level").
		Scan(&rows).Error
	if err != nil {
		return LevelCounts{}, err
	}

	var out LevelCounts
	for _, r := range rows {
		switch strings.ToUpper(r.Level) {
		case "INFO":
			out.Info = r.Count
		case "WARN", "WARNING":
			out.Warning = r.Count
		case "ERROR":
			out.Error = r.Count
		}
	}
	return out, nil
}

// Ensure time import is used in file (kept for future query expansions).
var _ = time.Time{}
