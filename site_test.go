package main

import (
	"encoding/json"
	"os"
	"testing"
	"time"
)

func TestWriteSiteData_JSONShapeAndDeltaFields(t *testing.T) {
	fixed := time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC)
	orig := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = orig }()

	stats := []Stat{
		{
			CanonicalKey:  "org/repo1",
			Owner:         "org",
			Repo:          "repo1",
			Tags:          []string{"terminal", "community"},
			Notes:         "some note",
			Description:   "a repo",
			Stars:         100,
			Language:      "Go",
			PushedAt:      fixed,
			URL:           "https://github.com/org/repo1",
			NameWithOwner: "org/repo1",
			IsArchived:    false,
		},
		{
			CanonicalKey:  "org/repo2",
			Owner:         "org",
			Repo:          "repo2",
			Tags:          []string{"web", "community"},
			Description:   "another repo",
			Stars:         50,
			Language:      "TypeScript",
			PushedAt:      fixed,
			URL:           "https://github.com/org/repo2",
			NameWithOwner: "org/repo2",
			IsArchived:    true,
		},
	}
	deltas7 := map[string]int{"org/repo1": 10}
	deltas30 := map[string]int{"org/repo1": 40}
	history := []Snapshot{
		{Date: "2026-08-01", Stars: map[string]int{"org/repo1": 90, "org/repo2": 45}},
	}

	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/data.json"

	if err := writeSiteData(tmpFile, stats, deltas7, deltas30, history); err != nil {
		t.Fatalf("writeSiteData: %v", err)
	}

	raw, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}

	var got siteData
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal output: %v", err)
	}

	if got.UpdatedAt != "2026-08-09 12:00 UTC" {
		t.Errorf("UpdatedAt: expected fixed timeNow value, got %q", got.UpdatedAt)
	}
	if len(got.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(got.Rows))
	}
	if len(got.History) != 1 {
		t.Fatalf("expected history to pass through unchanged, got %d snapshots", len(got.History))
	}

	row1 := got.Rows[0]
	if row1.Key != "org/repo1" || row1.NameWithOwner != "org/repo1" {
		t.Errorf("row1 identity mismatch: %+v", row1)
	}
	if !row1.HasDelta || row1.Delta7d != 10 {
		t.Errorf("row1 delta7d: expected hasDelta=true delta7d=10, got hasDelta=%v delta7d=%d", row1.HasDelta, row1.Delta7d)
	}
	if !row1.HasDelta30 || row1.Delta30d != 40 {
		t.Errorf("row1 delta30d: expected hasDelta30=true delta30d=40, got hasDelta30=%v delta30d=%d", row1.HasDelta30, row1.Delta30d)
	}
	if row1.Archived {
		t.Error("row1: expected archived=false")
	}

	row2 := got.Rows[1]
	if row2.HasDelta || row2.Delta7d != 0 {
		t.Errorf("row2 delta7d: expected no baseline (hasDelta=false, delta7d=0), got hasDelta=%v delta7d=%d", row2.HasDelta, row2.Delta7d)
	}
	if row2.HasDelta30 {
		t.Error("row2: expected no 30d baseline")
	}
	if !row2.Archived {
		t.Error("row2: expected archived=true")
	}

	// Field-name contract with site/index.html: verify the raw JSON keys,
	// since a struct-tag typo wouldn't be caught by round-tripping through
	// the same Go struct above.
	var rawMap map[string]any
	if err := json.Unmarshal(raw, &rawMap); err != nil {
		t.Fatalf("unmarshal raw map: %v", err)
	}
	rows, ok := rawMap["rows"].([]any)
	if !ok || len(rows) == 0 {
		t.Fatalf("expected rows array in raw JSON, got %v", rawMap["rows"])
	}
	firstRow, ok := rows[0].(map[string]any)
	if !ok {
		t.Fatalf("expected row to be an object, got %T", rows[0])
	}
	for _, key := range []string{"key", "nameWithOwner", "url", "stars", "delta7d", "hasDelta", "delta30d", "hasDelta30", "language", "pushedAt", "description", "tags", "archived"} {
		if _, ok := firstRow[key]; !ok {
			t.Errorf("expected JSON field %q in row, got keys: %v", key, mapKeys(firstRow))
		}
	}
}

func mapKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}
