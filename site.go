package main

import (
	"encoding/json"
	"io"
)

// siteRow is one ranked repo in the generated dist/data.json, consumed by site/index.html.
type siteRow struct {
	Key           string   `json:"key"` // canonical owner/repo, matches history keys
	NameWithOwner string   `json:"nameWithOwner"`
	URL           string   `json:"url"`
	Stars         int      `json:"stars"`
	Delta7d       int      `json:"delta7d"`
	HasDelta      bool     `json:"hasDelta"`
	Delta30d      int      `json:"delta30d"`
	HasDelta30    bool     `json:"hasDelta30"`
	Language      string   `json:"language"`
	PushedAt      string   `json:"pushedAt"`
	Description   string   `json:"description"`
	Tags          []string `json:"tags"`
	Notes         string   `json:"notes,omitempty"`
	Archived      bool     `json:"archived"`
}

type siteData struct {
	UpdatedAt string     `json:"updatedAt"`
	Rows      []siteRow  `json:"rows"`
	History   []Snapshot `json:"history"`
	// Facets drives the dashboard's filter chips, so the tag vocabulary is
	// defined once in Go rather than duplicated in site/index.html.
	Facets []tagFacet `json:"facets"`
}

// writeSiteData emits the JSON payload the dashboard fetches at runtime.
// It is build output (dist/data.json), never committed.
//
// updatedAt is passed in rather than read from the clock because it labels
// data freshness, not build time: a redeploy of unchanged data must not claim
// the figures are newer than the fetch that produced them.
func writeSiteData(path string, updatedAt string, stats []Stat, deltas7, deltas30 map[string]int, history []Snapshot) error {
	rows := make([]siteRow, len(stats))
	for i, s := range stats {
		delta7, has7 := deltas7[s.CanonicalKey]
		delta30, has30 := deltas30[s.CanonicalKey]
		rows[i] = siteRow{
			Key:           s.CanonicalKey,
			NameWithOwner: s.NameWithOwner,
			URL:           s.URL,
			Stars:         s.Stars,
			Delta7d:       delta7,
			HasDelta:      has7,
			Delta30d:      delta30,
			HasDelta30:    has30,
			Language:      s.Language,
			PushedAt:      s.PushedAt.Format("2006-01-02"),
			Description:   s.Description,
			Tags:          s.Tags,
			Notes:         s.Notes,
			Archived:      s.IsArchived,
		}
	}

	return atomicWriteFile(path, func(w io.Writer) error {
		return json.NewEncoder(w).Encode(siteData{
			UpdatedAt: updatedAt,
			Rows:      rows,
			History:   history,
			Facets:    tagVocabulary,
		})
	})
}
