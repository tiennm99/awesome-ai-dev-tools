package main

import (
	"bufio"
	"cmp"
	"encoding/json"
	"fmt"
	"io"
	"maps"
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
// and the history file contains the old name.
//
// Known renames:
//   - goose moved block → aaif-goose (Agentic AI Foundation); history contains
//     both spellings from earlier flip-flops, so map the non-canonical one.
//   - opencode moved sst → anomalyco
//   - gpt-engineer moved gpt-engineer-org → AntonOsika
//   - avante.nvim moved yetone → avante-corp
//   - openclaude moved Gitlawb → Twigpine
//   - Codewhale moved Hmbown → codewhale-hq
var canonicalKeyMigrations = map[string]string{
	"block/goose":                   "aaif-goose/goose",
	"sst/opencode":                  "anomalyco/opencode",
	"gpt-engineer-org/gpt-engineer": "AntonOsika/gpt-engineer",
	"yetone/avante.nvim":            "avante-corp/avante.nvim",
	"Gitlawb/openclaude":            "Twigpine/openclaude",
	"Hmbown/Codewhale":              "codewhale-hq/Codewhale",
}

// timeNow is a seam for tests: production code always calls time.Now, but
// tests can override this var to exercise fixed-calendar-date scenarios
// (delta window boundaries, migration+delta combinations) deterministically.
var timeNow = time.Now

// historyRetentionDays is how much star history is kept. Older snapshots are
// dropped on each update, so an entry removed from agents.yml keeps its chart
// data until its last snapshot ages past this window, then disappears.
const historyRetentionDays = 365

// appendHistory persists today's snapshot and returns the retained snapshot
// list (oldest first, including today) plus the 7-day and 30-day deltas per
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

	// Drop any pre-existing snapshot for today and any older than the
	// retention window, then append current.
	oldest := timeNow().UTC().AddDate(0, 0, -historyRetentionDays).Format("2006-01-02")
	kept := slices.DeleteFunc(history, func(s Snapshot) bool { return s.Date == today || s.Date < oldest })
	kept = append(kept, current)

	if err := writeSnapshots(path, kept); err != nil {
		return nil, nil, nil, err
	}
	return kept, deltas7, deltas30, nil
}

// historySchema is the version written in the history file's header line.
// Bump it when a change would make an older reader misread the file;
// adding a new per-day field (forks, say) does not need a bump, since
// readers ignore fields they do not know.
const historySchema = 1

// historyHeader is the first line of the history file. Repos lists the
// column order shared by every per-day array that follows.
type historyHeader struct {
	Schema int      `json:"schema"`
	Repos  []string `json:"repos"`
}

// historyRow is one day. Stars[i] belongs to header.Repos[i]; null, or a
// position past the end of a shorter array, means the repo was not tracked
// that day.
type historyRow struct {
	Date  string `json:"date"`
	Stars []*int `json:"stars"`
}

// readSnapshots loads the history file: JSON Lines whose first line is a
// historyHeader and whose other lines are historyRows. A missing or empty
// file is an empty history.
func readSnapshots(path string) ([]Snapshot, error) {
	header, rows, err := readHistoryFile(path)
	if err != nil || header == nil {
		return nil, err
	}

	out := make([]Snapshot, 0, len(rows))
	for _, r := range rows {
		s := Snapshot{Date: r.Date, Stars: map[string]int{}}
		for i, v := range r.Stars {
			if v != nil && i < len(header.Repos) {
				s.Stars[header.Repos[i]] = *v
			}
		}
		// Apply canonical-key migrations so history keyed under old names is
		// transparently remapped to the current agents.yml canonical key.
		s.Stars = applyMigrations(s.Stars)
		out = append(out, s)
	}

	// Delta windows pick the newest snapshot inside the window by scan order,
	// so the slice must be date-ascending. Sort rather than trust the file:
	// a hand edit or a merge of two concurrent runs can interleave lines.
	slices.SortStableFunc(out, func(a, b Snapshot) int { return cmp.Compare(a.Date, b.Date) })

	return out, nil
}

// readHistoryFile parses the raw header and rows. Malformed day lines are
// logged and skipped; a missing, malformed or too-new header is an error,
// because without it no row can be read correctly.
func readHistoryFile(path string) (*historyHeader, []historyRow, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil, nil
		}
		return nil, nil, err
	}
	defer func() { _ = f.Close() }() // read-only; close error is not actionable

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 1024*1024), 16*1024*1024)
	var header *historyHeader
	var rows []historyRow
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		if header == nil {
			var h historyHeader
			if err := json.Unmarshal(line, &h); err != nil || h.Schema == 0 {
				return nil, nil, fmt.Errorf("%s line %d: want a header like {\"schema\":%d,\"repos\":[...]}", path, lineNum, historySchema)
			}
			if h.Schema > historySchema {
				return nil, nil, fmt.Errorf("%s: schema %d is newer than this tool supports (%d)", path, h.Schema, historySchema)
			}
			header = &h
			continue
		}
		var r historyRow
		if err := json.Unmarshal(line, &r); err != nil {
			fmt.Fprintf(os.Stderr, "%s line %d: skipping malformed entry: %v\n", path, lineNum, err)
			continue
		}
		if _, err := time.Parse("2006-01-02", r.Date); err != nil {
			fmt.Fprintf(os.Stderr, "%s line %d: skipping entry with bad date %q\n", path, lineNum, r.Date)
			continue
		}
		rows = append(rows, r)
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, err
	}
	return header, rows, nil
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

// writeSnapshots persists the snapshots in the format readSnapshots reads.
// Each repo is named once, in the header, instead of on every day's line.
// Repos already in the file keep their column, and new ones are appended,
// so adding a repo does not reshuffle every past line in the git diff.
// Repos absent from every snapshot (aged out of retention) are dropped.
func writeSnapshots(path string, snapshots []Snapshot) error {
	present := map[string]bool{}
	for _, s := range snapshots {
		for k := range s.Stars {
			present[k] = true
		}
	}

	var repos []string
	if old, _, err := readHistoryFile(path); err == nil && old != nil {
		for _, k := range old.Repos {
			k = resolveCanonicalKey(k)
			if present[k] {
				repos = append(repos, k)
				delete(present, k)
			}
		}
	}
	repos = append(repos, slices.Sorted(maps.Keys(present))...)

	return atomicWriteFile(path, func(w io.Writer) error {
		enc := json.NewEncoder(w)
		if err := enc.Encode(historyHeader{Schema: historySchema, Repos: repos}); err != nil {
			return err
		}
		for _, s := range snapshots {
			row := historyRow{Date: s.Date, Stars: make([]*int, len(repos))}
			for i, k := range repos {
				if n, ok := s.Stars[k]; ok {
					row.Stars[i] = &n
				}
			}
			// Trailing nulls carry no data: a short array reads the same.
			for len(row.Stars) > 0 && row.Stars[len(row.Stars)-1] == nil {
				row.Stars = row.Stars[:len(row.Stars)-1]
			}
			if err := enc.Encode(row); err != nil {
				return err
			}
		}
		return nil
	})
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
	return computeDeltaAt(history, current, timeNow(), days, slackDays)
}

// computeDeltaAt is computeDeltaOver with an explicit anchor instead of the
// wall clock. The build step passes the time the data was fetched, so
// rebuilding an unchanged dataset later reproduces identical deltas rather
// than sliding the window forward off its baseline.
func computeDeltaAt(history []Snapshot, current Snapshot, anchor time.Time, days, slackDays int) map[string]int {
	deltas := map[string]int{}
	now := anchor.UTC()
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
