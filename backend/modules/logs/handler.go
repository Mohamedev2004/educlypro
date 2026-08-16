package logs

import (
	"educlypro/shared/utils"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func splitCSV(values []string) []string {
	out := make([]string, 0)
	for _, v := range values {
		for _, part := range strings.Split(v, ",") {
			part = strings.TrimSpace(part)
			if part != "" {
				out = append(out, part)
			}
		}
	}
	return out
}

func parseTimePtr(s string) (*time.Time, error) {
	if strings.TrimSpace(s) == "" {
		return nil, nil
	}
	// Accept RFC3339 timestamps.
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, err
	}
	return &t, nil
}

func queryArrayCompat(c *gin.Context, key string) []string {
	// Support both `key=value1&key=value2` and Axios default `key[]=value1&key[]=value2`.
	values := c.QueryArray(key)
	if len(values) == 0 {
		values = c.QueryArray(key + "[]")
	}
	return values
}

func normalizeLevelFilter(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "info":
		return "INFO"
	case "warning", "warn":
		return "WARN"
	case "error":
		return "ERROR"
	default:
		return strings.ToUpper(strings.TrimSpace(v))
	}
}

func normalizeStatusFilter(v string) string {
	switch strings.ToLower(strings.TrimSpace(v)) {
	case "created":
		return "CREATED"
	case "updated":
		return "UPDATED"
	case "deleted":
		return "DELETED"
	case "failed":
		return "FAILED"
	default:
		return strings.ToUpper(strings.TrimSpace(v))
	}
}

func parseStatusCodes(values []string) ([]int, error) {
	parts := splitCSV(values)
	out := make([]int, 0, len(parts))
	for _, p := range parts {
		n, err := strconv.Atoi(strings.TrimSpace(p))
		if err != nil {
			return nil, err
		}
		if n > 0 {
			out = append(out, n)
		}
	}
	return out, nil
}

// GET /logs
func (h *Handler) List(c *gin.Context) {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil || page < 1 {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_logs_page", "Please provide a valid page number.", map[string]string{"page": "validation.invalid"})
		return
	}

	perPage, err := strconv.Atoi(c.DefaultQuery("per_page", "10"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_logs_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}
	if _, ok := AllowedPerPage[perPage]; !ok {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_logs_per_page", "Please provide a valid page size.", map[string]string{"per_page": "validation.invalid"})
		return
	}

	from, err := parseTimePtr(c.Query("from"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_logs_from", "Please provide a valid from timestamp.", map[string]string{"from": "validation.invalid"})
		return
	}
	to, err := parseTimePtr(c.Query("to"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_logs_to", "Please provide a valid to timestamp.", map[string]string{"to": "validation.invalid"})
		return
	}

	rawLevels := splitCSV(queryArrayCompat(c, "level"))
	levels := make([]string, 0, len(rawLevels))
	for _, v := range rawLevels {
		if vv := normalizeLevelFilter(v); vv != "" {
			levels = append(levels, vv)
		}
	}

	rawStatuses := splitCSV(queryArrayCompat(c, "status"))
	statuses := make([]string, 0, len(rawStatuses))
	for _, v := range rawStatuses {
		if vv := normalizeStatusFilter(v); vv != "" {
			statuses = append(statuses, vv)
		}
	}

	durations := splitCSV(queryArrayCompat(c, "duration"))

	statusCodes, err := parseStatusCodes(queryArrayCompat(c, "status_code"))
	if err != nil {
		utils.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid_logs_status_code",
			"Please provide a valid status code filter.",
			map[string]string{"status_code": "validation.invalid"},
		)
		return
	}

	results, err := h.service.List(c.Request.Context(), ListParams{
		Page:        page,
		PerPage:     perPage,
		Query:       c.Query("q"),
		Levels:      levels,
		Statuses:    statuses,
		StatusCodes: statusCodes,
		Durations:   durations,
		From:        from,
		To:          to,
	})
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "logs_list_failed", "Failed to load logs. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "logs retrieved", results)
}

// GET /logs/chart?range=24h|7d
func (h *Handler) Chart(c *gin.Context) {
	r := ChartRange(strings.TrimSpace(c.DefaultQuery("range", string(ChartRange24h))))
	if r != ChartRange24h && r != ChartRange7d {
		utils.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid_logs_chart_range",
			"Please provide a valid chart range.",
			map[string]string{"range": "validation.invalid"},
		)
		return
	}

	points, err := h.service.Chart(c.Request.Context(), r)
	if err != nil {
		utils.ErrorResponseWithCode(c, http.StatusInternalServerError, "logs_chart_failed", "Failed to load logs chart. Please try again.")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, "logs chart retrieved", points)
}

// GET /logs/export
func (h *Handler) Export(c *gin.Context) {
	// Reuse the exact same param parsing as List — no pagination needed
	from, err := parseTimePtr(c.Query("from"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_logs_from", "Invalid from timestamp.", map[string]string{"from": "validation.invalid"})
		return
	}
	to, err := parseTimePtr(c.Query("to"))
	if err != nil {
		utils.ValidationErrorResponse(c, http.StatusBadRequest, "invalid_logs_to", "Invalid to timestamp.", map[string]string{"to": "validation.invalid"})
		return
	}

	rawLevels := splitCSV(queryArrayCompat(c, "level"))
	levels := make([]string, 0, len(rawLevels))
	for _, v := range rawLevels {
		if vv := normalizeLevelFilter(v); vv != "" {
			levels = append(levels, vv)
		}
	}

	rawStatuses := splitCSV(queryArrayCompat(c, "status"))
	statuses := make([]string, 0, len(rawStatuses))
	for _, v := range rawStatuses {
		if vv := normalizeStatusFilter(v); vv != "" {
			statuses = append(statuses, vv)
		}
	}

	statusCodes, err := parseStatusCodes(queryArrayCompat(c, "status_code"))
	if err != nil {
		utils.ValidationErrorResponse(
			c,
			http.StatusBadRequest,
			"invalid_logs_status_code",
			"Invalid status code filter.",
			map[string]string{"status_code": "validation.invalid"},
		)
		return
	}

	params := ListParams{
		Query:       c.Query("q"),
		Levels:      levels,
		Statuses:    statuses,
		StatusCodes: statusCodes,
		Durations:   splitCSV(queryArrayCompat(c, "duration")),
		From:        from,
		To:          to,
	}

	filename := fmt.Sprintf("audit-logs-%s.xlsx", time.Now().UTC().Format("2006-01-02"))

	c.Header("Content-Type", "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet")
	c.Header("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, filename))
	c.Header("Cache-Control", "no-cache")

	if err := h.service.Export(c.Request.Context(), params, c.Writer); err != nil {
		// Headers already sent — can't send JSON error, just log it
		log.Printf("export error: %v", err)
	}
}
