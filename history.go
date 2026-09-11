package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"slices"
	"time"
)

type Snapshot struct {
	Date  string         `json:"date"`
	Stars map[string]int `json:"stars"`
}

// canonicalKeyMigrations remaps old history keys to the current canonical
// owner/repo from agents.yml. Add entries here when a tracked repo is renamed
// and history.jsonl contains the old name.
//
// Known renames:
//   - goose moved block → aaif-goose (Agentic AI Foundation); history contains
//     both spellings from earlier flip-flops, so map the non-canonical one.
//   - opencode moved sst → anomalyco
//   - gpt-engineer moved gpt-engineer-org → AntonOsika
//   - avante.nvim moved yetone → avante-corp
var canonicalKeyMigrations = map[string]string{
	"block/goose":                   "aaif-goose/goose",
	"sst/opencode":                  "anomalyco/opencode",
	"gpt-engineer-org/gpt-engineer": "AntonOsika/gpt-engineer",
	"yetone/avante.nvim":            "avante-corp/avante.nvim",
}

// timeNow is a seam for tests: production code always calls time.Now, but
// tests can override this var to exercise fixed-calendar-date scenarios
// (delta window boundaries, migration+delta combinations) deterministically.
var timeNow = time.Now

// appendHistory persists today's snapshot and returns the full snapshot list
// (oldest first, including today) plus the 7-day and 30-day deltas per
// canonical key.
func appendHistory(path string, stats []Stat) (snapshots []Snapshot, deltas7 map[string]int, deltas30 map[string]int, err error) {
	today := timeNow().UTC().Format("2006-01-02")

	// Key snapshots by canonical owner/repo from agents.yml, not by the
	// API-returned nameWithOwner, so renames don't orphan historical data.
	current := Snapshot{Date: today, Stars: map[string]int{}}
	for _, s := range stats {
		current.Stars[s.CanonicalKey] = s.Stars
	}

	history, err := readSnapshots(path)
	if err != nil {
		return nil, nil, nil, err
	}

	deltas7 = computeDeltas(history, current)
	deltas30 = computeDeltaOver(history, current, 30, 5)

	// Drop any pre-existing snapshot for today, then append current.
	kept := slices.DeleteFunc(history, func(s Snapshot) bool { return s.Date == today })
	kept = append(kept, current)

	if err := writeSnapshots(path, kept); err != nil {
		return nil, nil, nil, err
	}
	return kept, deltas7, deltas30, nil
}

func readSnapshots(path string) ([]Snapshot, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	defer func() { _ = f.Close() }() // read-only; close error is not actionable

	var out []Snapshot
	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var s Snapshot
		if err := json.Unmarshal(line, &s); err != nil {
			// Log corruption so it's observable; don't silently drop lines.
			fmt.Fprintf(os.Stderr, "history.jsonl line %d: skipping malformed entry: %v\n", lineNum, err)
			continue
		}
		// Apply canonical-key migrations so history keyed under old names is
		// transparently remapped to the current agents.yml canonical key.
		s.Stars = applyMigrations(s.Stars)
		out = append(out, s)
	}
	return out, scanner.Err()
}

// applyMigrations rewrites any deprecated history keys to their current
// canonical form. Old keys are removed; new keys accumulate stars additively
// (in practice the old key had no concurrent new entry, so max is the same).
//
// Each old key is resolved transitively via resolveCanonicalKey, so a chain
// of renames (A→B, B→C) always lands on the terminal key C regardless of the
// order canonicalKeyMigrations happens to be iterated in (map order is
// randomized by Go at runtime).
func applyMigrations(stars map[string]int) map[string]int {
	for old := range canonicalKeyMigrations {
		v, ok := stars[old]
		if !ok {
			continue
		}
		canonical := resolveCanonicalKey(old)
		if canonical == old {
			continue // cycle or hop-cap hit; resolveCanonicalKey already logged
		}
		if stars[canonical] < v {
			stars[canonical] = v
		}
		delete(stars, old)
	}
	return stars
}

// maxMigrationHops caps chain resolution so a misconfigured cycle in
// canonicalKeyMigrations can't hang the updater.
const maxMigrationHops = 10

// resolveCanonicalKey follows canonicalKeyMigrations transitively from key
// until it reaches a terminal key (one with no further migration entry).
// If following the chain would exceed maxMigrationHops, or a cycle is
// detected, the cycle/misconfiguration is logged to stderr and the original
// key is returned unchanged — this degrades to "no migration applied"
// rather than dropping data or looping forever.
func resolveCanonicalKey(key string) string {
	seen := map[string]bool{key: true}
	current := key
	for hop := 0; hop < maxMigrationHops; hop++ {
		next, ok := canonicalKeyMigrations[current]
		if !ok {
			return current
		}
		if seen[next] {
			fmt.Fprintf(os.Stderr, "canonicalKeyMigrations: cycle detected resolving %q (hit %q again); keeping original key unchanged\n", key, next)
			return key
		}
		seen[next] = true
		current = next
	}
	fmt.Fprintf(os.Stderr, "canonicalKeyMigrations: exceeded %d hops resolving %q; keeping original key unchanged\n", maxMigrationHops, key)
	return key
}

// writeSnapshots writes to a temp file then renames atomically so a crash
// mid-write never leaves history.jsonl truncated or partially written.
func writeSnapshots(path string, snapshots []Snapshot) error {
	tmp := path + ".tmp"
	f, err := os.Create(tmp)
	if err != nil {
		return err
	}

	enc := json.NewEncoder(f)
	writeErr := error(nil)
	for _, s := range snapshots {
		if err := enc.Encode(s); err != nil {
			writeErr = err
			break
		}
	}

	if syncErr := f.Sync(); syncErr != nil && writeErr == nil {
		writeErr = syncErr
	}
	// A failed close on a write path can hide lost data — surface it.
	if closeErr := f.Close(); closeErr != nil && writeErr == nil {
		writeErr = closeErr
	}

	if writeErr != nil {
		_ = os.Remove(tmp) // best-effort cleanup; the write error is what matters
		return writeErr
	}

	return os.Rename(tmp, path)
}

// computeDeltas returns stars-now minus stars-at-or-before-cutoff for each
// repo, cutoff = 7 days ago UTC, with a 3-day slack window. See
// computeDeltaOver for the general form (used for the 30-day/"momentum"
// delta too).
func computeDeltas(history []Snapshot, current Snapshot) map[string]int {
	return computeDeltaOver(history, current, 7, 3)
}

// computeDeltaOver returns stars-now minus stars-at-or-before-cutoff for
// each repo, where cutoff = days ago (UTC). The chosen prior snapshot must
// fall within a slackDays window of the cutoff — i.e. in
// (cutoff - slackDays, cutoff]; if cron was skipped for longer than that,
// the delta would be misleadingly labeled "Δ<days>d", so we return no delta
// for that repo instead.
func computeDeltaOver(history []Snapshot, current Snapshot, days, slackDays int) map[string]int {
	deltas := map[string]int{}
	now := timeNow().UTC()
	cutoff := now.AddDate(0, 0, -days).Format("2006-01-02")
	lowerBound := now.AddDate(0, 0, -(days + slackDays)).Format("2006-01-02")

	var base *Snapshot
	for i := range history {
		d := history[i].Date
		if d > lowerBound && d <= cutoff {
			base = &history[i]
		}
	}
	if base == nil {
		return deltas
	}
	for repo, cur := range current.Stars {
		if prev, ok := base.Stars[repo]; ok {
			deltas[repo] = cur - prev
		}
	}
	return deltas
}
