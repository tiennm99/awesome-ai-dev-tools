package main

import (
	"encoding/json"
	"io"
)

// siteRow is one ranked repo in site/data.json, consumed by site/index.html.
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

// writeSiteData emits the JSON payload for the GitHub Pages dashboard.
// The file is generated fresh on every updater run and is not committed;
// the Pages deploy step in the workflow picks it up from the working tree.
func writeSiteData(path string, stats []Stat, deltas7, deltas30 map[string]int, history []Snapshot) error {
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
			UpdatedAt: timeNow().UTC().Format("2006-01-02 15:04 UTC"),
			Rows:      rows,
			History:   history,
			Facets:    tagVocabulary,
		})
	})
}
