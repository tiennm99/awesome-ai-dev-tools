package main

import (
	"cmp"
	"fmt"
	"io"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"text/template"
)

type Row struct {
	Rank          int
	RankLabel     string // medal for the top three, the number otherwise
	NameWithOwner string
	URL           string
	Stars         int
	Delta7d       int
	HasDelta      bool
	Language      string
	PushedAt      string
	Description   string
	Surface       string // the entry's surface tags, as inline code
}

// trendingCount is how many of the biggest 7-day gainers the README lists
// above the ranking.
const trendingCount = 5

// rankMedals labels the top three rows.
var rankMedals = []string{"🥇", "🥈", "🥉"}

func renderReadme(tmplPath, outPath string, stats []Stat, deltas map[string]int) error {
	rows := make([]Row, len(stats))
	for i, s := range stats {
		// Deltas are keyed by CanonicalKey (owner/repo from agents.yml).
		delta, has := deltas[s.CanonicalKey]
		label := strconv.Itoa(i + 1)
		if i < len(rankMedals) {
			label = rankMedals[i]
		}
		rows[i] = Row{
			Rank:          i + 1,
			RankLabel:     label,
			NameWithOwner: s.NameWithOwner,
			URL:           s.URL,
			Stars:         s.Stars,
			Delta7d:       delta,
			HasDelta:      has,
			Language:      s.Language,
			PushedAt:      s.PushedAt.Format("2006-01-02"),
			Description:   sanitizeCell(s.Description),
			Surface:       surfaceCell(s.Tags),
		}
	}

	funcs := template.FuncMap{
		"formatDelta": func(d int, has bool) string {
			if !has {
				return "—"
			}
			if d > 0 {
				return "+" + groupThousands(d)
			}
			if d < 0 {
				return "-" + groupThousands(-d)
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
			"Trending":  trending(rows),
		})
	})
}

// groupThousands formats a non-negative n with comma thousands separators.
func groupThousands(n int) string {
	s := strconv.Itoa(n)
	for i := len(s) - 3; i > 0; i -= 3 {
		s = s[:i] + "," + s[i:]
	}
	return s
}

// trending returns up to trendingCount rows with the biggest positive 7-day
// gain, biggest first. It is empty while history is younger than 7 days.
func trending(rows []Row) []Row {
	var out []Row
	for _, r := range rows {
		if r.HasDelta && r.Delta7d > 0 {
			out = append(out, r)
		}
	}
	slices.SortStableFunc(out, func(a, b Row) int { return cmp.Compare(b.Delta7d, a.Delta7d) })
	if len(out) > trendingCount {
		out = out[:trendingCount]
	}
	return out
}

// surfaceCell renders the surface-facet tags (where a tool runs) as inline
// code for one table cell.
func surfaceCell(tags []string) string {
	var parts []string
	for _, t := range tags {
		if tagFacetIDs[t] == facetSurface {
			parts = append(parts, "`"+t+"`")
		}
	}
	return strings.Join(parts, " ")
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
