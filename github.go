package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Stat struct {
	// CanonicalKey is owner/repo from agents.yml — stable across renames.
	CanonicalKey  string
	Owner         string
	Repo          string
	Tags          []string
	Notes         string
	Description   string
	Stars         int
	Language      string
	PushedAt      time.Time
	URL           string
	NameWithOwner string
	IsArchived    bool
}

type repoNode struct {
	StargazerCount  int    `json:"stargazerCount"`
	Description     string `json:"description"`
	PrimaryLanguage *struct {
		Name string `json:"name"`
	} `json:"primaryLanguage"`
	PushedAt      time.Time `json:"pushedAt"`
	URL           string    `json:"url"`
	NameWithOwner string    `json:"nameWithOwner"`
	IsArchived    bool      `json:"isArchived"`
}

type graphQLResponse struct {
	Data   map[string]*repoNode `json:"data"`
	Errors []struct {
		Message string `json:"message"`
		Path    []any  `json:"path"`
	} `json:"errors"`
}

const repoFields = `
		stargazerCount
		description
		primaryLanguage { name }
		pushedAt
		url
		nameWithOwner
		isArchived
	`

// httpClient has a timeout to prevent hung workflow jobs.
var httpClient = &http.Client{Timeout: 30 * time.Second}

// graphqlURL is a seam for tests: production always talks to GitHub, but
// tests point this at an httptest server to exercise the whole fetch path
// (chunking, alias offsets, missing-node, GraphQL-error handling) without
// a network call or token.
var graphqlURL = "https://api.github.com/graphql"

// staleWarnAfter is how long a tracked repo may go without a push before the
// updater flags it for review. It mirrors the review window in the README's
// inclusion criteria; dropping an entry stays a human decision, so this only
// annotates the run.
const staleWarnAfter = 90 * 24 * time.Hour

// minStars is the hard inclusion floor from the README's criteria. Unlike the
// staleness window, this is not a judgement call: an entry below the floor is
// dropped from the published ranking, so the list can never show a repo that
// does not meet it. It lives here rather than in -check because star counts
// need the API and -check runs offline.
const minStars = 1000

// chunkSize is the max aliases per GraphQL request (GitHub node-limit safety margin).
const chunkSize = 50

// maxRetries and retry backoff for transient HTTP/network errors.
const maxRetries = 3

var retryBackoff = []time.Duration{2 * time.Second, 4 * time.Second, 8 * time.Second}

// fetchStats queries GitHub GraphQL in chunks of up to chunkSize repos per
// request. Returns an error if any GraphQL errors are present OR if any
// requested repo is missing from the response — better to fail loud than
// silently publish a shorter README.
func fetchStats(token string, agents []Agent) ([]Stat, error) {
	collected := make(map[string]*repoNode, len(agents))

	for start := 0; start < len(agents); start += chunkSize {
		end := start + chunkSize
		if end > len(agents) {
			end = len(agents)
		}
		chunk := agents[start:end]

		nodes, err := fetchChunk(token, chunk, start)
		if err != nil {
			return nil, err
		}
		for k, v := range nodes {
			collected[k] = v
		}
	}

	stats := make([]Stat, 0, len(agents))
	for i, a := range agents {
		alias := fmt.Sprintf("r%d", i)
		node := collected[alias]
		if node == nil {
			return nil, fmt.Errorf("repo %s/%s (alias %s) missing from GraphQL response — deleted, private, or renamed?", a.Owner, a.Repo, alias)
		}
		lang := ""
		if node.PrimaryLanguage != nil {
			lang = node.PrimaryLanguage.Name
		}

		// Drift detection: surface renames, archival and staleness as GitHub
		// Actions run annotations so a human notices without polling every
		// repo by hand.
		canonicalKey := a.Owner + "/" + a.Repo
		if node.NameWithOwner != "" && !strings.EqualFold(node.NameWithOwner, canonicalKey) {
			fmt.Printf("::warning::repo %s renamed to %s — update data/agents.yml and add a canonicalKeyMigrations entry\n", canonicalKey, node.NameWithOwner)
		}
		switch {
		case node.IsArchived:
			// Archived already implies no further pushes; one warning is enough.
			fmt.Printf("::warning::repo %s is archived — consider removing or annotating\n", canonicalKey)
		case !node.PushedAt.IsZero() && timeNow().Sub(node.PushedAt) > staleWarnAfter:
			days := int(timeNow().Sub(node.PushedAt).Hours() / 24)
			fmt.Printf("::warning::repo %s has no push in %d days — review against the maintenance criterion\n", canonicalKey, days)
		}

		stats = append(stats, Stat{
			CanonicalKey:  canonicalKey,
			Owner:         a.Owner,
			Repo:          a.Repo,
			Tags:          a.Tags,
			Notes:         a.Notes,
			Description:   node.Description,
			Stars:         node.StargazerCount,
			Language:      lang,
			PushedAt:      node.PushedAt,
			URL:           node.URL,
			NameWithOwner: node.NameWithOwner,
			IsArchived:    node.IsArchived,
		})
	}

	// Sort by stars descending. Ties are ordered by CanonicalKey for determinism
	// regardless of map-iteration or agents.yml order.
	sort.Slice(stats, func(i, j int) bool {
		if stats[i].Stars != stats[j].Stars {
			return stats[i].Stars > stats[j].Stars
		}
		return stats[i].CanonicalKey < stats[j].CanonicalKey
	})

	return stats, nil
}

// fetchChunk sends one GraphQL request for a slice of agents. aliasOffset
// ensures alias names (r0, r1, …) are globally unique across chunks.
func fetchChunk(token string, agents []Agent, aliasOffset int) (map[string]*repoNode, error) {
	var b strings.Builder
	b.WriteString("query {\n")
	for i, a := range agents {
		fmt.Fprintf(&b, "  r%d: repository(owner: %q, name: %q) {%s}\n", aliasOffset+i, a.Owner, a.Repo, repoFields)
	}
	b.WriteString("}\n")

	body, err := json.Marshal(map[string]string{"query": b.String()})
	if err != nil {
		return nil, err
	}

	raw, statusCode, err := doWithRetry(token, body)
	if err != nil {
		return nil, err
	}
	if statusCode != http.StatusOK {
		return nil, fmt.Errorf("graphql HTTP %d: %s", statusCode, raw)
	}

	var out graphQLResponse
	if err := json.Unmarshal(raw, &out); err != nil {
		return nil, fmt.Errorf("decode response: %w (body=%s)", err, raw)
	}

	// Treat any GraphQL-level error as fatal — partial data silently shrinks
	// the README and corrupts future delta math.
	if len(out.Errors) > 0 {
		msgs := make([]string, len(out.Errors))
		for i, e := range out.Errors {
			msgs[i] = fmt.Sprintf("%s (%s)", e.Message, describeErrorPath(e.Path, agents, aliasOffset))
		}
		return nil, fmt.Errorf("graphql errors: %s", strings.Join(msgs, "; "))
	}

	return out.Data, nil
}

// describeErrorPath translates a GraphQL error path such as
// ["r17", "stargazerCount"] into a human-readable repo reference, e.g.
// "repo foo/bar (alias r17)", using the alias→Agent mapping for this chunk
// (aliasOffset + local index). Falls back to the raw path if the alias
// can't be resolved (unexpected path shape or out-of-range index).
func describeErrorPath(path []any, agents []Agent, aliasOffset int) string {
	if len(path) == 0 {
		return "path=[]"
	}
	aliasStr, ok := path[0].(string)
	if !ok || !strings.HasPrefix(aliasStr, "r") {
		return fmt.Sprintf("path=%v", path)
	}
	idx, err := strconv.Atoi(aliasStr[1:])
	if err != nil {
		return fmt.Sprintf("path=%v", path)
	}
	local := idx - aliasOffset
	if local < 0 || local >= len(agents) {
		return fmt.Sprintf("path=%v (alias %s)", path, aliasStr)
	}
	a := agents[local]
	return fmt.Sprintf("repo %s/%s (alias %s)", a.Owner, a.Repo, aliasStr)
}

// doWithRetry executes the GraphQL POST with exponential backoff on transient
// errors (network failures, HTTP 5xx, HTTP 429). 4xx other than 429 are not
// retried.
func doWithRetry(token string, body []byte) ([]byte, int, error) {
	var lastErr error
	for attempt := 0; attempt < maxRetries; attempt++ {
		if attempt > 0 {
			log.Printf("retry attempt %d after %v: %v", attempt, retryBackoff[attempt-1], lastErr)
			time.Sleep(retryBackoff[attempt-1])
		}

		req, err := http.NewRequest("POST", graphqlURL, bytes.NewReader(body))
		if err != nil {
			return nil, 0, err
		}
		req.Header.Set("Authorization", "Bearer "+token)
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("User-Agent", "awesome-coding-agents-updater")

		resp, err := httpClient.Do(req)
		if err != nil {
			lastErr = err
			continue // network error — retry
		}

		raw, err := io.ReadAll(resp.Body)
		_ = resp.Body.Close() // body fully read (or read failed); close error is not actionable
		if err != nil {
			lastErr = fmt.Errorf("read body: %w", err)
			continue
		}

		sc := resp.StatusCode
		if sc == http.StatusOK {
			return raw, sc, nil
		}
		if sc == http.StatusTooManyRequests || sc >= 500 {
			lastErr = fmt.Errorf("HTTP %d: %s", sc, raw)
			continue // retryable
		}
		// 4xx (except 429) — not retryable
		return raw, sc, nil
	}
	return nil, 0, fmt.Errorf("all %d attempts failed; last error: %w", maxRetries, lastErr)
}

// enforceStarFloor drops entries below minStars and annotates each one as a run
// error. Dropping rather than failing the run keeps one below-floor entry from
// blocking the refresh of every other repo, while the ::error:: annotation
// makes the removal impossible to miss in the Actions log.
func enforceStarFloor(stats []Stat) []Stat {
	kept := make([]Stat, 0, len(stats))
	for _, s := range stats {
		if s.Stars < minStars {
			fmt.Printf("::error::repo %s has %d stars, below the %d minimum — dropped from the ranking; remove its entry from data/agents.yml\n", s.CanonicalKey, s.Stars, minStars)
			continue
		}
		kept = append(kept, s)
	}
	return kept
}
