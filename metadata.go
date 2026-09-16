package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"
)

// repoMeta is the GitHub-sourced half of a ranking row: every field the
// dashboard needs that cannot be derived from data/agents.yml. It is written
// by the update step (which has a token) and committed, so the build step can
// render the site offline.
type repoMeta struct {
	NameWithOwner string    `json:"nameWithOwner"`
	URL           string    `json:"url"`
	Stars         int       `json:"stars"`
	Language      string    `json:"language"`
	PushedAt      time.Time `json:"pushedAt"`
	Description   string    `json:"description"`
	IsArchived    bool      `json:"archived"`
}

// metadataFile is the on-disk shape of data/metadata.json. Repos is keyed by
// canonical owner/repo (the same key history.jsonl uses), so entries survive
// renames and Go's sorted map marshalling keeps commit diffs minimal.
type metadataFile struct {
	FetchedAt time.Time           `json:"fetchedAt"`
	Repos     map[string]repoMeta `json:"repos"`
}

// writeMetadata persists the fetched GitHub metadata. Unlike the old
// site/data.json this file IS committed: it is the handoff between the update
// step and the build step, and the only reason the build can run without a
// GITHUB_TOKEN.
func writeMetadata(path string, stats []Stat) error {
	repos := make(map[string]repoMeta, len(stats))
	for _, s := range stats {
		repos[s.CanonicalKey] = repoMeta{
			NameWithOwner: s.NameWithOwner,
			URL:           s.URL,
			Stars:         s.Stars,
			Language:      s.Language,
			PushedAt:      s.PushedAt,
			Description:   s.Description,
			IsArchived:    s.IsArchived,
		}
	}

	return atomicWriteFile(path, func(w io.Writer) error {
		enc := json.NewEncoder(w)
		// Indented so the committed file reviews as a readable diff.
		enc.SetIndent("", "  ")
		return enc.Encode(metadataFile{
			FetchedAt: timeNow().UTC(),
			Repos:     repos,
		})
	})
}

func readMetadata(path string) (metadataFile, error) {
	var m metadataFile
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return m, fmt.Errorf("%s not found — run the updater (`go run .` with GITHUB_TOKEN set) to generate it", path)
		}
		return m, err
	}
	if err := json.Unmarshal(data, &m); err != nil {
		return m, fmt.Errorf("parse %s: %w", path, err)
	}
	if len(m.Repos) == 0 {
		return m, fmt.Errorf("%s contains no repos", path)
	}
	return m, nil
}

// statsFromMetadata rejoins the two halves of a row: curation (order-independent
// tags and notes) from agents.yml, live figures from metadata.json.
//
// A repo present in agents.yml but absent from metadata.json is skipped with a
// warning rather than failing the build. That is the normal state between
// "contributor adds an entry" and "the next nightly update fetches it", and a
// hard failure there would block every Cloudflare deploy in the meantime.
func statsFromMetadata(agents []Agent, meta metadataFile) []Stat {
	stats := make([]Stat, 0, len(agents))
	for _, a := range agents {
		key := a.Owner + "/" + a.Repo
		m, ok := meta.Repos[key]
		if !ok {
			fmt.Fprintf(os.Stderr, "warning: %s has no entry in metadata yet — omitted from this build; the next updater run will add it\n", key)
			continue
		}
		stats = append(stats, Stat{
			CanonicalKey:  key,
			Owner:         a.Owner,
			Repo:          a.Repo,
			Tags:          a.Tags,
			Notes:         a.Notes,
			Description:   m.Description,
			Stars:         m.Stars,
			Language:      m.Language,
			PushedAt:      m.PushedAt,
			URL:           m.URL,
			NameWithOwner: m.NameWithOwner,
			IsArchived:    m.IsArchived,
		})
	}
	return stats
}
