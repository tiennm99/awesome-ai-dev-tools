package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"testing"
	"time"
)

func TestComputeDeltas_NoHistory(t *testing.T) {
	// Empty history slice should return empty delta map without panicking.
	current := Snapshot{
		Date: "2026-05-14",
		Stars: map[string]int{
			"org/repo": 100,
		},
	}
	deltas := computeDeltas(nil, current)
	if deltas == nil {
		t.Errorf("expected non-nil empty map, got nil")
	}
	if len(deltas) != 0 {
		t.Errorf("expected empty deltas, got %d entries", len(deltas))
	}
}

func TestComputeDeltas_SevenDayWindow(t *testing.T) {
	// The cutoff window is (cutoff-3d, cutoff] where cutoff = now - 7d.
	// So valid snapshots have date in range (now-10d, now-7d].

	now := time.Now().UTC()
	today := now.Format("2006-01-02")

	tests := []struct {
		name        string
		daysAgo     int
		expectDelta bool
		desc        string
	}{
		{
			name:        "snapshot at exactly 7 days ago",
			daysAgo:     7,
			expectDelta: true,
			desc:        "at cutoff boundary — should be included",
		},
		{
			name:        "snapshot at 9 days ago",
			daysAgo:     9,
			expectDelta: true,
			desc:        "within (cutoff-3d, cutoff] — should be included",
		},
		{
			name:        "snapshot at 11 days ago",
			daysAgo:     11,
			expectDelta: false,
			desc:        "older than cutoff-3d — should be excluded",
		},
		{
			name:        "snapshot at 6 days ago",
			daysAgo:     6,
			expectDelta: false,
			desc:        "newer than cutoff — should be excluded",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			snapshotDate := now.AddDate(0, 0, -tt.daysAgo).Format("2006-01-02")
			history := []Snapshot{
				{
					Date: snapshotDate,
					Stars: map[string]int{
						"org/repo": 100,
					},
				},
			}
			current := Snapshot{
				Date: today,
				Stars: map[string]int{
					"org/repo": 150,
				},
			}
			deltas := computeDeltas(history, current)

			if tt.expectDelta {
				if len(deltas) == 0 {
					t.Errorf("%s: expected delta for 'org/repo', got empty map", tt.desc)
				}
				if delta, ok := deltas["org/repo"]; !ok || delta != 50 {
					t.Errorf("%s: expected delta=50, got %v", tt.desc, delta)
				}
			} else {
				if len(deltas) > 0 {
					t.Errorf("%s: expected no delta, got %v", tt.desc, deltas)
				}
			}
		})
	}
}

func TestReadSnapshots_MalformedLine(t *testing.T) {
	// Temp JSONL with: valid, malformed, empty, valid lines.
	// Should skip malformed/empty and return only the 2 valid snapshots.

	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/history.jsonl"

	content := `{"date":"2026-05-01","stars":{"org/repo":100}}
not json at all { bad
{"date":"2026-05-02","stars":{"org/repo":200}}
{"date":"2026-05-03","stars":{"org/repo":300}}
`

	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	snapshots, err := readSnapshots(tmpFile)
	if err != nil {
		t.Fatalf("readSnapshots failed: %v", err)
	}

	if len(snapshots) != 3 {
		t.Errorf("expected 3 snapshots (skipped 1 malformed), got %d", len(snapshots))
	}

	// Check dates are in order.
	if len(snapshots) >= 3 {
		if snapshots[0].Date != "2026-05-01" {
			t.Errorf("snapshot[0] date: expected 2026-05-01, got %s", snapshots[0].Date)
		}
		if snapshots[1].Date != "2026-05-02" {
			t.Errorf("snapshot[1] date: expected 2026-05-02, got %s", snapshots[1].Date)
		}
		if snapshots[2].Date != "2026-05-03" {
			t.Errorf("snapshot[2] date: expected 2026-05-03, got %s", snapshots[2].Date)
		}
	}
}

func TestReadSnapshots_MissingFile(t *testing.T) {
	// Non-existent file should return nil with no error (not os.IsNotExist error).
	snapshots, err := readSnapshots("/nonexistent/path/history.jsonl")
	if err != nil {
		t.Errorf("expected nil error on missing file, got %v", err)
	}
	if snapshots != nil {
		t.Errorf("expected nil snapshots on missing file, got %v", snapshots)
	}
}

func TestApplyMigrations(t *testing.T) {
	// Test the canonical key migration logic.
	stars := map[string]int{
		"block/goose": 45115,
		"other/repo":  1000,
	}

	result := applyMigrations(stars)

	// The old key should be removed.
	if _, ok := result["block/goose"]; ok {
		t.Errorf("old key 'block/goose' should be removed")
	}

	// The new key should be present with the value.
	if v, ok := result["aaif-goose/goose"]; !ok || v != 45115 {
		t.Errorf("new key 'aaif-goose/goose': expected 45115, got %v", v)
	}

	// Other keys should be unchanged.
	if v, ok := result["other/repo"]; !ok || v != 1000 {
		t.Errorf("'other/repo': expected 1000, got %v", v)
	}
}

func TestApplyMigrations_ChainResolvesFullyRegardlessOfOrder(t *testing.T) {
	// A→B→C chain must always collapse to C, whether the map happens to be
	// iterated A-then-B or B-then-A (Go randomizes map iteration order).
	origMigrations := canonicalKeyMigrations
	defer func() { canonicalKeyMigrations = origMigrations }()
	canonicalKeyMigrations = map[string]string{
		"org/a": "org/b",
		"org/b": "org/c",
	}

	// Run many times: map iteration order is randomized per run, so this
	// gives confidence the result doesn't depend on it.
	for i := 0; i < 50; i++ {
		stars := map[string]int{"org/a": 10, "other/repo": 5}
		result := applyMigrations(stars)

		if _, ok := result["org/a"]; ok {
			t.Fatalf("iteration %d: old key 'org/a' should be removed", i)
		}
		if _, ok := result["org/b"]; ok {
			t.Fatalf("iteration %d: dead intermediate key 'org/b' should not exist", i)
		}
		if v, ok := result["org/c"]; !ok || v != 10 {
			t.Fatalf("iteration %d: terminal key 'org/c': expected 10, got %v (ok=%v)", i, v, ok)
		}
		if v, ok := result["other/repo"]; !ok || v != 5 {
			t.Fatalf("iteration %d: unrelated key 'other/repo' should be unchanged", i)
		}
	}
}

func TestResolveCanonicalKey_CycleDoesNotHang(t *testing.T) {
	origMigrations := canonicalKeyMigrations
	defer func() { canonicalKeyMigrations = origMigrations }()
	canonicalKeyMigrations = map[string]string{
		"org/a": "org/b",
		"org/b": "org/a",
	}

	done := make(chan string, 1)
	go func() { done <- resolveCanonicalKey("org/a") }()

	select {
	case got := <-done:
		if got != "org/a" {
			t.Errorf("expected cycle to resolve back to original key 'org/a', got %q", got)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("resolveCanonicalKey hung on a migration cycle")
	}
}

func TestResolveCanonicalKey_ExceedsHopCapKeepsOriginal(t *testing.T) {
	origMigrations := canonicalKeyMigrations
	defer func() { canonicalKeyMigrations = origMigrations }()
	// A chain longer than maxMigrationHops with no cycle — must still
	// terminate and fall back to the original key.
	migrations := make(map[string]string, maxMigrationHops+5)
	for i := 0; i < maxMigrationHops+5; i++ {
		migrations[fmt.Sprintf("org/k%d", i)] = fmt.Sprintf("org/k%d", i+1)
	}
	canonicalKeyMigrations = migrations

	got := resolveCanonicalKey("org/k0")
	if got != "org/k0" {
		t.Errorf("expected hop-cap fallback to original key 'org/k0', got %q", got)
	}
}

func TestCanonicalKeyMigrations_NoChains(t *testing.T) {
	// Invariant: every migration should point directly at the terminal
	// canonical key. If a future edit introduces A→B, B→C without collapsing
	// it to A→C, this test fails — guarding against the exact bug the
	// transitive resolver in resolveCanonicalKey defends against at runtime.
	for old, canonical := range canonicalKeyMigrations {
		if _, ok := canonicalKeyMigrations[canonical]; ok {
			t.Errorf("canonicalKeyMigrations[%q] = %q, but %q is itself a migration key — collapse the chain to point directly at the terminal canonical key", old, canonical, canonical)
		}
	}
}

func TestComputeDeltaOver_FixedNow_BoundaryEdges(t *testing.T) {
	// Exercise the (cutoff-slackDays, cutoff] window with a fixed "now" so
	// the exact day-count boundaries (7d cutoff, 10d lower bound for the
	// 7-day delta) are deterministic rather than relative-to-test-run-time.
	fixed := time.Date(2026, 8, 9, 12, 0, 0, 0, time.UTC)
	orig := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = orig }()

	mk := func(daysAgo int) Snapshot {
		return Snapshot{
			Date:  fixed.AddDate(0, 0, -daysAgo).Format("2006-01-02"),
			Stars: map[string]int{"org/repo": 50},
		}
	}
	current := Snapshot{Date: fixed.Format("2006-01-02"), Stars: map[string]int{"org/repo": 200}}

	tests := []struct {
		name        string
		daysAgo     int
		expectDelta bool
	}{
		{"exactly at 7d cutoff — included", 7, true},
		{"9d ago — inside (cutoff-3d, cutoff] window", 9, true},
		{"exactly at 10d lower bound — excluded (window is exclusive lower bound)", 10, false},
		{"11d ago — excluded", 11, false},
		{"6d ago — newer than cutoff, excluded", 6, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := []Snapshot{mk(tt.daysAgo)}
			deltas := computeDeltas(history, current)
			_, ok := deltas["org/repo"]
			if ok != tt.expectDelta {
				t.Errorf("daysAgo=%d: expected delta present=%v, got %v", tt.daysAgo, tt.expectDelta, ok)
			}
			if ok && deltas["org/repo"] != 150 {
				t.Errorf("daysAgo=%d: expected delta=150, got %d", tt.daysAgo, deltas["org/repo"])
			}
		})
	}
}

func TestComputeDeltaOver_30Day_FixedNow(t *testing.T) {
	// Same boundary shape as the 7-day delta, but for the 30d/5d-slack
	// "momentum" window used for Delta30d.
	fixed := time.Date(2026, 8, 9, 0, 0, 0, 0, time.UTC)
	orig := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = orig }()

	mk := func(daysAgo, stars int) Snapshot {
		return Snapshot{
			Date:  fixed.AddDate(0, 0, -daysAgo).Format("2006-01-02"),
			Stars: map[string]int{"org/repo": stars},
		}
	}
	current := Snapshot{Date: fixed.Format("2006-01-02"), Stars: map[string]int{"org/repo": 1000}}

	tests := []struct {
		name        string
		daysAgo     int
		expectDelta bool
	}{
		{"exactly 30d cutoff — included", 30, true},
		{"33d ago — inside (cutoff-5d, cutoff] window", 33, true},
		{"exactly 35d lower bound — excluded", 35, false},
		{"36d ago — excluded", 36, false},
		{"29d ago — newer than cutoff, excluded", 29, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			history := []Snapshot{mk(tt.daysAgo, 400)}
			deltas := computeDeltaOver(history, current, 30, 5)
			_, ok := deltas["org/repo"]
			if ok != tt.expectDelta {
				t.Errorf("daysAgo=%d: expected delta present=%v, got %v", tt.daysAgo, tt.expectDelta, ok)
			}
			if ok && deltas["org/repo"] != 600 {
				t.Errorf("daysAgo=%d: expected delta=600, got %d", tt.daysAgo, deltas["org/repo"])
			}
		})
	}
}

func TestWriteSnapshots_AtomicWrite(t *testing.T) {
	// Verify that writeSnapshots uses atomic rename (writes to .tmp first).
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/history.jsonl"

	snapshots := []Snapshot{
		{Date: "2026-05-01", Stars: map[string]int{"org/repo": 100}},
		{Date: "2026-05-02", Stars: map[string]int{"org/repo": 200}},
	}

	err := writeSnapshots(tmpFile, snapshots)
	if err != nil {
		t.Fatalf("writeSnapshots failed: %v", err)
	}

	// Verify file exists and has correct content.
	if _, err := os.Stat(tmpFile); err != nil {
		t.Fatalf("output file missing: %v", err)
	}

	content, err := os.ReadFile(tmpFile)
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Verify JSON is valid and correct.
	var readBack []Snapshot
	decoder := json.NewDecoder(bytes.NewReader(content))
	for decoder.More() {
		var s Snapshot
		if err := decoder.Decode(&s); err != nil {
			t.Fatalf("decode failed: %v", err)
		}
		readBack = append(readBack, s)
	}

	if len(readBack) != 2 {
		t.Errorf("expected 2 snapshots, got %d", len(readBack))
	}
	if readBack[0].Date != "2026-05-01" || readBack[0].Stars["org/repo"] != 100 {
		t.Errorf("snapshot 0 mismatch: %+v", readBack[0])
	}
	if readBack[1].Date != "2026-05-02" || readBack[1].Stars["org/repo"] != 200 {
		t.Errorf("snapshot 1 mismatch: %+v", readBack[1])
	}

	// .tmp file should not exist (cleaned up after rename).
	if _, err := os.Stat(tmpFile + ".tmp"); !os.IsNotExist(err) {
		t.Errorf("temp file should not exist after successful write")
	}
}

func TestReadSnapshots_SortsOutOfOrderDates(t *testing.T) {
	// A hand edit or a merge of two concurrent updater runs can leave lines
	// out of date order; delta windows depend on ascending order.
	path := t.TempDir() + "/history.jsonl"
	lines := `{"date":"2026-09-03","stars":{"org/repo":300}}
{"date":"2026-08-27","stars":{"org/repo":100}}
{"date":"2026-09-10","stars":{"org/repo":400}}
{"date":"2026-09-01","stars":{"org/repo":200}}
`
	if err := os.WriteFile(path, []byte(lines), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}

	got, err := readSnapshots(path)
	if err != nil {
		t.Fatalf("readSnapshots: %v", err)
	}
	want := []string{"2026-08-27", "2026-09-01", "2026-09-03", "2026-09-10"}
	if len(got) != len(want) {
		t.Fatalf("got %d snapshots, want %d", len(got), len(want))
	}
	for i, d := range want {
		if got[i].Date != d {
			t.Errorf("snapshot %d date = %s, want %s", i, got[i].Date, d)
		}
	}

	// With the dates ordered, the 7-day base is 2026-09-03 (not the
	// last line in the file), so the delta is 400-300.
	fixed := time.Date(2026, 9, 10, 0, 0, 0, 0, time.UTC)
	orig := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = orig }()

	deltas := computeDeltas(got, Snapshot{Date: "2026-09-10", Stars: map[string]int{"org/repo": 400}})
	if deltas["org/repo"] != 100 {
		t.Errorf("delta = %d, want 100", deltas["org/repo"])
	}
}
