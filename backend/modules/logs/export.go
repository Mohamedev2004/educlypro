package logs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/xuri/excelize/v2"
)

var exportHeaders = []string{
	"ID", "Timestamp", "Level", "Service", "Entity ID",
	"Actor ID", "Action", "Message", "Duration", "Request ID", "Payload",
}

var exportCols = []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J", "K"}

func rowValues(l AuditLog) []any {
	payloadStr := ""
	if len(l.Payload) > 0 {
		var p any
		if json.Unmarshal(l.Payload, &p) == nil {
			b, _ := json.Marshal(p)
			payloadStr = string(b)
		}
	}
	return []any{
		l.ID,
		l.CreatedAt.UTC().Format(time.RFC3339),
		levelForUI(l.Level),
		l.Entity,
		l.EntityID,
		l.ActorID,
		statusForUI(l.Action),
		messageForUI(l),
		durationForUI(l.DurationMs),
		l.RequestID,
		payloadStr,
	}
}

func (s *service) Export(ctx context.Context, params ListParams, w io.Writer) error {
	// ── Pass 1: measure max content width per column ───────────────────────
	maxWidths := make([]float64, len(exportCols))
	for i, h := range exportHeaders {
		maxWidths[i] = float64(len(h))
	}

	if err := s.repo.ListAll(ctx, params, func(l AuditLog) error {
		for i, val := range rowValues(l) {
			if w := float64(len(fmt.Sprintf("%v", val))); w > maxWidths[i] {
				maxWidths[i] = w
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("export measure pass error: %w", err)
	}

	// ── Build workbook ─────────────────────────────────────────────────────
	f := excelize.NewFile()
	defer f.Close()

	sheet := "Audit Logs"
	f.SetSheetName("Sheet1", sheet)

	// ── Header style ───────────────────────────────────────────────────────
	headerStyle, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Bold: true, Color: "FFFFFF", Size: 11, Family: "Arial"},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"1F3864"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center", Vertical: "center"},
		Border: []excelize.Border{
			{Type: "bottom", Color: "FFFFFF", Style: 2},
		},
	})

	// ── Apply auto-sized column widths ─────────────────────────────────────
	for i, col := range exportCols {
		width := maxWidths[i]*1.1 + 2 // 10% padding + breathing room
		if width > 60 {
			width = 60 // cap payload / message columns
		}
		if width < 8 {
			width = 8 // minimum readable width
		}
		f.SetColWidth(sheet, col, col, width)
	}

	// ── Write header row ───────────────────────────────────────────────────
	for i, h := range exportHeaders {
		cell := exportCols[i] + "1"
		f.SetCellValue(sheet, cell, h)
		f.SetCellStyle(sheet, cell, cell, headerStyle)
	}
	f.SetRowHeight(sheet, 1, 20)

	// ── Level badge styles ─────────────────────────────────────────────────
	styleInfo, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "155724", Family: "Arial", Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"D4EDDA"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	styleWarn, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "856404", Family: "Arial", Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"FFF3CD"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	styleError, _ := f.NewStyle(&excelize.Style{
		Font:      &excelize.Font{Color: "721C24", Family: "Arial", Size: 10},
		Fill:      excelize.Fill{Type: "pattern", Color: []string{"F8D7DA"}, Pattern: 1},
		Alignment: &excelize.Alignment{Horizontal: "center"},
	})
	styleDefault, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "Arial", Size: 10},
	})
	styleAlt, _ := f.NewStyle(&excelize.Style{
		Font: &excelize.Font{Family: "Arial", Size: 10},
		Fill: excelize.Fill{Type: "pattern", Color: []string{"F8F9FA"}, Pattern: 1},
	})

	levelStyle := map[string]int{
		"INFO":  styleInfo,
		"WARN":  styleWarn,
		"ERROR": styleError,
	}

	// ── Pass 2: write rows ─────────────────────────────────────────────────
	rowNum := 2
	if err := s.repo.ListAll(ctx, params, func(l AuditLog) error {
		r := rowNum

		rowStyle := styleDefault
		if r%2 == 0 {
			rowStyle = styleAlt
		}

		for i, val := range rowValues(l) {
			cell := fmt.Sprintf("%s%d", exportCols[i], r)
			f.SetCellValue(sheet, cell, val)
			f.SetCellStyle(sheet, cell, cell, rowStyle)
		}

		// Override level cell with badge style
		if ls, ok := levelStyle[l.Level]; ok {
			levelCell := fmt.Sprintf("C%d", r)
			f.SetCellStyle(sheet, levelCell, levelCell, ls)
		}

		rowNum++
		return nil
	}); err != nil {
		return fmt.Errorf("export write pass error: %w", err)
	}

	// ── Freeze header row ──────────────────────────────────────────────────
	if err := f.SetPanes(sheet, &excelize.Panes{
		Freeze:      true,
		Split:       false,
		TopLeftCell: "A2",
	}); err != nil {
		return fmt.Errorf("set panes error: %w", err)
	}

	// ── Auto-filter on header ──────────────────────────────────────────────
	f.AutoFilter(sheet, fmt.Sprintf("A1:K%d", rowNum-1), nil)

	// ── Stream directly to response writer — no temp file ─────────────────
	_, err := f.WriteTo(w)
	return err
}
