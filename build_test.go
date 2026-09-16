package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// writeBuildFixture lays down a minimal repo layout the build step can read:
// agents.yml (curation), metadata.json (fetched figures), history.jsonl
// (delta baselines) and a site/ source directory.
func writeBuildFixture(t *testing.T, fetchedAt time.Time) (dir string) {
	t.Helper()
	dir = t.TempDir()

	mustWrite := func(rel, content string) {
		t.Helper()
		path := filepath.Join(dir, rel)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			t.Fatalf("mkdir for %s: %v", rel, err)
		}
		if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
			t.Fatalf("write %s: %v", rel, err)
		}
	}

	mustWrite("data/agents.yml", `agents:
  - owner: org
    repo: small
    tags: [terminal, community]
  - owner: org
    repo: big
    tags: [web, vendor]
    notes: a note
`)

	meta := metadataFile{
		FetchedAt: fetchedAt,
		Repos: map[string]repoMeta{
			"org/big": {
				NameWithOwner: "org/big", URL: "https://github.com/org/big",
				Stars: 5000, Language: "Go", PushedAt: fetchedAt,
				Description: "big one",
			},
			"org/small": {
				NameWithOwner: "org/small", URL: "https://github.com/org/small",
				Stars: 2000, Language: "Rust", PushedAt: fetchedAt,
				Description: "small one", IsArchived: true,
			},
		},
	}
	raw, err := json.Marshal(meta)
	if err != nil {
		t.Fatalf("marshal metadata: %v", err)
	}
	mustWrite("data/metadata.json", string(raw))

	// Baseline 7 days before fetchedAt (2026-09-09) so deltas resolve.
	mustWrite("data/history.jsonl", `{"date":"2026-09-09","stars":{"org/big":4900,"org/small":1990}}
{"date":"2026-09-16","stars":{"org/big":5000,"org/small":2000}}
`)

	mustWrite("site/index.html", "<html>dashboard</html>")
	mustWrite("site/_headers", "/data.json\n  Cache-Control: no-cache\n")

	return dir
}

func readBuiltSiteData(t *testing.T, dir string) siteData {
	t.Helper()
	raw, err := os.ReadFile(filepath.Join(dir, "dist", "data.json"))
	if err != nil {
		t.Fatalf("read dist/data.json: %v", err)
	}
	var got siteData
	if err := json.Unmarshal(raw, &got); err != nil {
		t.Fatalf("unmarshal dist/data.json: %v", err)
	}
	return got
}

func runBuildIn(t *testing.T, dir string) error {
	t.Helper()
	return runBuild(
		filepath.Join(dir, "data/agents.yml"),
		filepath.Join(dir, "data/metadata.json"),
		filepath.Join(dir, "data/history.jsonl"),
		filepath.Join(dir, "site"),
		filepath.Join(dir, "dist"),
	)
}

func TestRunBuild_JoinsCurationWithMetadataAndCopiesSite(t *testing.T) {
	fetchedAt := time.Date(2026, 9, 16, 22, 10, 0, 0, time.UTC)
	dir := writeBuildFixture(t, fetchedAt)

	if err := runBuildIn(t, dir); err != nil {
		t.Fatalf("runBuild: %v", err)
	}

	got := readBuiltSiteData(t, dir)

	// updatedAt must label the fetch, not the build.
	if got.UpdatedAt != "2026-09-16 22:10 UTC" {
		t.Errorf("UpdatedAt: expected the metadata fetch time, got %q", got.UpdatedAt)
	}

	if len(got.Rows) != 2 {
		t.Fatalf("expected 2 rows, got %d", len(got.Rows))
	}
	// Ranked by stars descending, independent of agents.yml order.
	if got.Rows[0].Key != "org/big" || got.Rows[1].Key != "org/small" {
		t.Errorf("expected big ranked above small, got %q then %q", got.Rows[0].Key, got.Rows[1].Key)
	}
	// Curation comes from agents.yml...
	if got.Rows[0].Notes != "a note" {
		t.Errorf("expected notes from agents.yml, got %q", got.Rows[0].Notes)
	}
	// ...figures from metadata.json.
	if got.Rows[0].Stars != 5000 || got.Rows[0].Language != "Go" || got.Rows[0].Description != "big one" {
		t.Errorf("row0 metadata not applied: %+v", got.Rows[0])
	}
	if !got.Rows[1].Archived {
		t.Error("expected archived flag to survive the round trip")
	}
	if !got.Rows[0].HasDelta || got.Rows[0].Delta7d != 100 {
		t.Errorf("row0 delta7d: expected +100 from the 2026-09-09 baseline, got hasDelta=%v delta=%d", got.Rows[0].HasDelta, got.Rows[0].Delta7d)
	}
	if len(got.Facets) == 0 {
		t.Error("expected tag facets in the payload")
	}

	// Every static file in site/ must land in dist/ alongside data.json.
	for _, name := range []string{"index.html", "_headers"} {
		if _, err := os.Stat(filepath.Join(dir, "dist", name)); err != nil {
			t.Errorf("expected dist/%s: %v", name, err)
		}
	}
}

// The build must not depend on the clock: it runs whenever Cloudflare happens
// to redeploy, which may be long after the data was fetched.
func TestRunBuild_DeltasAnchoredToFetchTimeNotWallClock(t *testing.T) {
	fetchedAt := time.Date(2026, 9, 16, 22, 10, 0, 0, time.UTC)
	dir := writeBuildFixture(t, fetchedAt)

	orig := timeNow
	// Simulate a redeploy 40 days later; a wall-clock anchor would slide the
	// 7-day window past its slack and drop the delta entirely.
	timeNow = func() time.Time { return fetchedAt.AddDate(0, 0, 40) }
	defer func() { timeNow = orig }()

	if err := runBuildIn(t, dir); err != nil {
		t.Fatalf("runBuild: %v", err)
	}

	got := readBuiltSiteData(t, dir)
	if !got.Rows[0].HasDelta || got.Rows[0].Delta7d != 100 {
		t.Errorf("delta should be unchanged by a late rebuild, got hasDelta=%v delta=%d", got.Rows[0].HasDelta, got.Rows[0].Delta7d)
	}
	if got.UpdatedAt != "2026-09-16 22:10 UTC" {
		t.Errorf("UpdatedAt should still report the fetch time, got %q", got.UpdatedAt)
	}
}

// A newly contributed entry has no metadata until the next update run. The
// build must degrade to omitting it, not fail — otherwise every deploy is
// blocked between the PR merge and the nightly refresh.
func TestRunBuild_SkipsEntriesWithoutMetadataYet(t *testing.T) {
	fetchedAt := time.Date(2026, 9, 16, 22, 10, 0, 0, time.UTC)
	dir := writeBuildFixture(t, fetchedAt)

	agentsPath := filepath.Join(dir, "data/agents.yml")
	existing, err := os.ReadFile(agentsPath)
	if err != nil {
		t.Fatalf("read agents.yml: %v", err)
	}
	added := string(existing) + "  - owner: org\n    repo: brandnew\n    tags: [terminal]\n"
	if err := os.WriteFile(agentsPath, []byte(added), 0o644); err != nil {
		t.Fatalf("write agents.yml: %v", err)
	}

	if err := runBuildIn(t, dir); err != nil {
		t.Fatalf("runBuild should tolerate a missing metadata entry, got: %v", err)
	}

	got := readBuiltSiteData(t, dir)
	if len(got.Rows) != 2 {
		t.Fatalf("expected the un-fetched entry to be omitted, got %d rows", len(got.Rows))
	}
	for _, r := range got.Rows {
		if r.Key == "org/brandnew" {
			t.Error("org/brandnew has no metadata and must not be published")
		}
	}
}

// dist/ is build output: a file removed from site/ must not survive there.
func TestRunBuild_WipesStaleDistFiles(t *testing.T) {
	fetchedAt := time.Date(2026, 9, 16, 22, 10, 0, 0, time.UTC)
	dir := writeBuildFixture(t, fetchedAt)

	stale := filepath.Join(dir, "dist", "old.html")
	if err := os.MkdirAll(filepath.Dir(stale), 0o755); err != nil {
		t.Fatalf("mkdir dist: %v", err)
	}
	if err := os.WriteFile(stale, []byte("stale"), 0o644); err != nil {
		t.Fatalf("write stale file: %v", err)
	}

	if err := runBuildIn(t, dir); err != nil {
		t.Fatalf("runBuild: %v", err)
	}

	if _, err := os.Stat(stale); !os.IsNotExist(err) {
		t.Errorf("expected stale dist file to be removed, stat err = %v", err)
	}
}

func TestRunBuild_MissingMetadataIsActionable(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "data"), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dir, "data/agents.yml"), []byte("agents:\n  - owner: org\n    repo: big\n"), 0o644); err != nil {
		t.Fatalf("write agents.yml: %v", err)
	}

	err := runBuildIn(t, dir)
	if err == nil {
		t.Fatal("expected an error when metadata.json is absent")
	}
	if !strings.Contains(err.Error(), "metadata.json") || !strings.Contains(err.Error(), "GITHUB_TOKEN") {
		t.Errorf("error should name the missing file and how to produce it, got: %v", err)
	}
}
