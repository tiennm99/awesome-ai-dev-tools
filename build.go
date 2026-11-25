package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// runBuild renders the static dashboard into distDir from committed inputs
// only — data/agents.yml, data/metadata.json and data/history.jsonl. It makes
// no network calls and needs no GITHUB_TOKEN, which is what lets an untrusted
// build environment (Cloudflare Pages) run it.
//
// Splitting this out from the update step also means a pure curation change
// (retagging an entry, editing a note) republishes immediately on push,
// reusing the last fetched star figures instead of waiting for the nightly run.
func runBuild(agentsPath, metadataPath, historyPath, siteDir, distDir string) error {
	agents, err := loadAgents(agentsPath)
	if err != nil {
		return err
	}

	meta, err := readMetadata(metadataPath)
	if err != nil {
		return err
	}

	stats := statsFromMetadata(agents, meta)
	stats = enforceStarFloor(stats)
	if len(stats) == 0 {
		return fmt.Errorf("no entries left to publish after joining %s with %s", agentsPath, metadataPath)
	}
	sortStats(stats)

	history, err := readSnapshots(historyPath)
	if err != nil {
		return err
	}

	// Anchor the delta windows to when the data was fetched, not to wall clock.
	// A rebuild triggered days later (a docs push, a manual redeploy) must show
	// the same "Δ7d" the nightly run computed, not a window silently slid
	// forward past its slack allowance.
	current := Snapshot{
		Date:  meta.FetchedAt.UTC().Format("2006-01-02"),
		Stars: make(map[string]int, len(stats)),
	}
	for _, s := range stats {
		current.Stars[s.CanonicalKey] = s.Stars
	}
	deltas7 := computeDeltaAt(history, current, meta.FetchedAt, 7, 3)
	deltas30 := computeDeltaAt(history, current, meta.FetchedAt, 30, 5)

	if err := copyDir(siteDir, distDir); err != nil {
		return fmt.Errorf("copy %s to %s: %w", siteDir, distDir, err)
	}

	updatedAt := meta.FetchedAt.UTC().Format("2006-01-02 15:04 UTC")
	if err := writeSiteData(filepath.Join(distDir, "data.json"), updatedAt, stats, deltas7, deltas30, history); err != nil {
		return err
	}

	fmt.Printf("built %s: %d entries, data fetched %s\n", distDir, len(stats), updatedAt)
	return nil
}

// copyDir replaces dst with a fresh copy of src. The wipe matters: dist/ is
// build output, and a file deleted from site/ must not survive in a deploy
// because a previous build left it there.
func copyDir(src, dst string) error {
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	return filepath.WalkDir(src, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		target := filepath.Join(dst, rel)
		if d.IsDir() {
			return os.MkdirAll(target, 0o755)
		}
		if !d.Type().IsRegular() {
			return nil // skip symlinks and other irregular entries
		}
		return copyFile(path, target)
	})
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func() { _ = in.Close() }() // read-only; close error is not actionable

	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}
	return atomicWriteFile(dst, func(w io.Writer) error {
		_, err := io.Copy(w, in)
		return err
	})
}
