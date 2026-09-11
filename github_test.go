package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

// aliasRe extracts the numeric alias index from a GraphQL query fragment
// like "  r17: repository(owner: \"foo\", name: \"bar\") {...}".
var aliasRe = regexp.MustCompile(`r(\d+):\s*repository`)

// withGraphQLURL points graphqlURL at srv for the duration of the test and
// restores the original value on cleanup.
func withGraphQLURL(t *testing.T, srv *httptest.Server) {
	t.Helper()
	orig := graphqlURL
	graphqlURL = srv.URL
	t.Cleanup(func() { graphqlURL = orig })
}

func readGraphQLQuery(t *testing.T, r *http.Request) string {
	t.Helper()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		t.Fatalf("read request body: %v", err)
	}
	var req struct {
		Query string `json:"query"`
	}
	if err := json.Unmarshal(body, &req); err != nil {
		t.Fatalf("unmarshal request body: %v (body=%s)", err, body)
	}
	return req.Query
}

func TestFetchStats_ChunkingHappyPath(t *testing.T) {
	// 55 agents forces 2 chunks (50 + 5) — chunk-boundary code has zero
	// production coverage today (only 29 agents exist), so this is the
	// first exercise of the alias-offset math across a chunk boundary.
	const n = 55
	agents := make([]Agent, n)
	for i := range agents {
		agents[i] = Agent{Owner: "org", Repo: fmt.Sprintf("repo%02d", i), Category: "cli"}
	}

	var requestCount int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestCount++
		query := readGraphQLQuery(t, r)

		data := map[string]any{}
		for _, m := range aliasRe.FindAllStringSubmatch(query, -1) {
			idx, err := strconv.Atoi(m[1])
			if err != nil {
				t.Fatalf("parse alias index from %q: %v", m[0], err)
			}
			data["r"+m[1]] = map[string]any{
				"stargazerCount":  100 + idx,
				"description":     "desc " + m[1],
				"primaryLanguage": map[string]string{"name": "Go"},
				"pushedAt":        "2026-08-01T00:00:00Z",
				"url":             fmt.Sprintf("https://github.com/org/repo%02d", idx),
				"nameWithOwner":   fmt.Sprintf("org/repo%02d", idx),
				"isArchived":      false,
			}
		}
		out, err := json.Marshal(map[string]any{"data": data})
		if err != nil {
			t.Fatalf("marshal response: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		if _, err := w.Write(out); err != nil {
			t.Fatalf("write response: %v", err)
		}
	}))
	defer srv.Close()
	withGraphQLURL(t, srv)

	stats, err := fetchStats("test-token", agents)
	if err != nil {
		t.Fatalf("fetchStats: %v", err)
	}
	if requestCount != 2 {
		t.Errorf("expected 2 chunk requests (50 + 5), got %d", requestCount)
	}
	if len(stats) != n {
		t.Fatalf("expected %d stats, got %d", n, len(stats))
	}
	// Sorted descending by stars: repo54 (154 stars) must be first.
	if stats[0].NameWithOwner != "org/repo54" || stats[0].Stars != 154 {
		t.Errorf("expected top stat org/repo54 with 154 stars, got %s with %d", stats[0].NameWithOwner, stats[0].Stars)
	}
	// Spot-check a repo from the second chunk (alias offset 50) resolved
	// correctly — this is exactly the code path that has never run in
	// production (only 1 chunk is used with 29 agents).
	var found bool
	for _, s := range stats {
		if s.NameWithOwner == "org/repo52" {
			found = true
			if s.Stars != 152 {
				t.Errorf("org/repo52: expected 152 stars, got %d", s.Stars)
			}
		}
	}
	if !found {
		t.Error("expected org/repo52 (second chunk) in results")
	}
}

func TestFetchStats_MissingNodeNamesTheRepo(t *testing.T) {
	agents := []Agent{
		{Owner: "foo", Repo: "bar0", Category: "cli"},
		{Owner: "foo", Repo: "bar1", Category: "cli"},
		{Owner: "foo", Repo: "bar2", Category: "cli"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := readGraphQLQuery(t, r)
		data := map[string]any{}
		for _, m := range aliasRe.FindAllStringSubmatch(query, -1) {
			if m[1] == "2" {
				continue // simulate repo r2 (foo/bar2) deleted/private/renamed away
			}
			data["r"+m[1]] = map[string]any{
				"stargazerCount": 10,
				"nameWithOwner":  "foo/bar" + m[1],
			}
		}
		out, _ := json.Marshal(map[string]any{"data": data}) // static map of strings/ints/bools never fails to marshal
		_, _ = w.Write(out)                                  // httptest ResponseRecorder write error is not actionable in a test fake
	}))
	defer srv.Close()
	withGraphQLURL(t, srv)

	_, err := fetchStats("test-token", agents)
	if err == nil {
		t.Fatal("expected error for missing node, got nil")
	}
	if !strings.Contains(err.Error(), "foo/bar2") {
		t.Errorf("expected error to name the missing repo foo/bar2, got: %v", err)
	}
	if !strings.Contains(err.Error(), "alias r2") {
		t.Errorf("expected error to name the alias r2, got: %v", err)
	}
}

func TestFetchStats_GraphQLErrorNamesTheRepo(t *testing.T) {
	agents := []Agent{
		{Owner: "foo", Repo: "bar0", Category: "cli"},
		{Owner: "foo", Repo: "bar1", Category: "cli"},
		{Owner: "foo", Repo: "bar2", Category: "cli"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]any{
			"data": nil,
			"errors": []map[string]any{
				{"message": "Could not resolve to a Repository.", "path": []any{"r2"}},
			},
		}
		out, _ := json.Marshal(resp) // static map of strings/ints/bools never fails to marshal
		_, _ = w.Write(out)          // httptest ResponseRecorder write error is not actionable in a test fake
	}))
	defer srv.Close()
	withGraphQLURL(t, srv)

	_, err := fetchStats("test-token", agents)
	if err == nil {
		t.Fatal("expected error for GraphQL-level error, got nil")
	}
	if !strings.Contains(err.Error(), "repo foo/bar2") {
		t.Errorf("expected error to name repo foo/bar2, got: %v", err)
	}
	if !strings.Contains(err.Error(), "alias r2") {
		t.Errorf("expected error to name alias r2, got: %v", err)
	}
	if !strings.Contains(err.Error(), "Could not resolve to a Repository.") {
		t.Errorf("expected original GraphQL message preserved, got: %v", err)
	}
}

func TestFetchStats_DriftWarnings(t *testing.T) {
	// Renamed and archived repos should not error the run — fetchStats
	// still succeeds and reports the drift as a warning (checked via
	// captured stdout) rather than failing the whole daily update.
	agents := []Agent{
		{Owner: "old-owner", Repo: "renamed-repo", Category: "cli"},
		{Owner: "org", Repo: "archived-repo", Category: "cli"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		data := map[string]any{
			"r0": map[string]any{
				"stargazerCount": 5,
				"nameWithOwner":  "new-owner/renamed-repo",
				"isArchived":     false,
			},
			"r1": map[string]any{
				"stargazerCount": 5,
				"nameWithOwner":  "org/archived-repo",
				"isArchived":     true,
			},
		}
		out, _ := json.Marshal(map[string]any{"data": data}) // static map of strings/ints/bools never fails to marshal
		_, _ = w.Write(out)                                  // httptest ResponseRecorder write error is not actionable in a test fake
	}))
	defer srv.Close()
	withGraphQLURL(t, srv)

	stats, err := fetchStats("test-token", agents)
	if err != nil {
		t.Fatalf("fetchStats: %v (drift should warn, not fail)", err)
	}
	if len(stats) != 2 {
		t.Fatalf("expected 2 stats, got %d", len(stats))
	}
	for _, s := range stats {
		if s.CanonicalKey == "org/archived-repo" && !s.IsArchived {
			t.Error("expected org/archived-repo Stat.IsArchived=true")
		}
	}
}

// captureStdout redirects os.Stdout for the duration of fn and returns what
// was written, so the ::warning:: annotations fetchStats emits are assertable.
func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	orig := os.Stdout
	os.Stdout = w
	defer func() { os.Stdout = orig }()

	fn()

	if err := w.Close(); err != nil {
		t.Fatalf("close pipe writer: %v", err)
	}
	out, err := io.ReadAll(r)
	if err != nil {
		t.Fatalf("read pipe: %v", err)
	}
	return string(out)
}

func TestFetchStats_StaleWarning(t *testing.T) {
	fixed := time.Date(2026, 9, 11, 0, 0, 0, 0, time.UTC)
	origNow := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = origNow }()

	agents := []Agent{
		{Owner: "org", Repo: "fresh", Category: "cli"},
		{Owner: "org", Repo: "stale", Category: "cli"},
		{Owner: "org", Repo: "stale-and-archived", Category: "cli"},
	}

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		node := func(name, pushedAt string, archived bool) map[string]any {
			return map[string]any{
				"stargazerCount": 5,
				"nameWithOwner":  name,
				"pushedAt":       pushedAt,
				"isArchived":     archived,
			}
		}
		out, _ := json.Marshal(map[string]any{"data": map[string]any{
			// 2 days idle: inside the window, no warning.
			"r0": node("org/fresh", "2026-09-09T00:00:00Z", false),
			// 120 days idle: past staleWarnAfter.
			"r1": node("org/stale", "2026-05-14T00:00:00Z", false),
			// Also stale, but archived — the archived warning covers it.
			"r2": node("org/stale-and-archived", "2025-01-01T00:00:00Z", true),
		}})
		_, _ = w.Write(out) // httptest write error is not actionable in a test fake
	}))
	defer srv.Close()
	withGraphQLURL(t, srv)

	var stats []Stat
	var fetchErr error
	logs := captureStdout(t, func() { stats, fetchErr = fetchStats("test-token", agents) })

	if fetchErr != nil {
		t.Fatalf("fetchStats: %v (staleness should warn, not fail)", fetchErr)
	}
	if len(stats) != 3 {
		t.Fatalf("expected 3 stats, got %d", len(stats))
	}

	if want := "::warning::repo org/stale has no push in 120 days"; !strings.Contains(logs, want) {
		t.Errorf("missing stale warning %q in:\n%s", want, logs)
	}
	if strings.Contains(logs, "org/fresh has no push") {
		t.Errorf("warned about a repo pushed 2 days ago:\n%s", logs)
	}
	if strings.Contains(logs, "org/stale-and-archived has no push") {
		t.Errorf("stale warning duplicates the archived warning:\n%s", logs)
	}
	if want := "::warning::repo org/stale-and-archived is archived"; !strings.Contains(logs, want) {
		t.Errorf("missing archived warning %q in:\n%s", want, logs)
	}
}
