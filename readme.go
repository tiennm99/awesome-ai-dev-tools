package main

import (
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"text/template"
)

type Row struct {
	Rank          int
	NameWithOwner string
	URL           string
	Stars         int
	Delta7d       int
	HasDelta      bool
	Language      string
	PushedAt      string
	Description   string
	Category      string
}

// TopMover is the biggest 7-day gainer, surfaced as a one-line callout above
// the README table. HasMover is false when no repo has a 7-day delta yet
// (e.g. history younger than 7 days) — the template omits the line then.
type TopMover struct {
	NameWithOwner string
	URL           string
	Delta         int
	HasMover      bool
}

func renderReadme(tmplPath, outPath string, stats []Stat, deltas map[string]int) error {
	rows := make([]Row, len(stats))
	var topMover TopMover
	for i, s := range stats {
		// Deltas are keyed by CanonicalKey (owner/repo from agents.yml).
		delta, has := deltas[s.CanonicalKey]
		rows[i] = Row{
			Rank:          i + 1,
			NameWithOwner: s.NameWithOwner,
			URL:           s.URL,
			Stars:         s.Stars,
			Delta7d:       delta,
			HasDelta:      has,
			Language:      s.Language,
			PushedAt:      s.PushedAt.Format("2006-01-02"),
			Description:   sanitizeCell(s.Description),
			Category:      s.Category,
		}
		if has && (!topMover.HasMover || delta > topMover.Delta) {
			topMover = TopMover{NameWithOwner: s.NameWithOwner, URL: s.URL, Delta: delta, HasMover: true}
		}
	}

	funcs := template.FuncMap{
		"formatDelta": func(d int, has bool) string {
			if !has {
				return "—"
			}
			if d > 0 {
				return fmt.Sprintf("+%d", d)
			}
			if d < 0 {
				return fmt.Sprintf("%d", d)
			}
			return "0"
		},
		"formatStars": func(n int) string {
			switch {
			case n < 1000:
				return fmt.Sprintf("%d", n)
			case n < 1_000_000:
				return fmt.Sprintf("%.1fk", float64(n)/1000)
			default:
				return fmt.Sprintf("%.1fM", float64(n)/1_000_000)
			}
		},
	}

	tmpl, err := template.New("").Funcs(funcs).ParseFiles(tmplPath)
	if err != nil {
		return err
	}

	return atomicWriteFile(outPath, func(w io.Writer) error {
		return tmpl.ExecuteTemplate(w, filepath.Base(tmplPath), map[string]any{
			"Rows":      rows,
			"UpdatedAt": timeNow().UTC().Format("2006-01-02 15:04 UTC"),
			"Total":     len(rows),
			"TopMover":  topMover,
		})
	})
}

// sanitizeCell makes a third-party repo description safe to embed in a
// single Markdown table cell. Order matters:
//  1. backslash first, so escaping added by later steps isn't re-escaped
//     (also fixes `\|` rendering as a literal backslash + unescaped pipe,
//     which broke the table row).
//  2. pipe, so the cell can't inject a table column boundary.
//  3. angle brackets → HTML entities, so raw HTML/script tags in a
//     description can't be injected into the rendered README.
//  4. newlines/carriage returns → space, so the cell stays one line.
func sanitizeCell(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "|", "\\|")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.ReplaceAll(s, "\r", " ")
	return strings.TrimSpace(s)
}
