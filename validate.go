package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ownerPattern approximates GitHub's username/org rules: alphanumeric runs
// separated by single hyphens — equivalent to
// "^[A-Za-z0-9](?:[A-Za-z0-9]|-(?=[A-Za-z0-9]))*$" (no leading/trailing
// hyphen, no consecutive hyphens) but written without lookahead, which Go's
// RE2-based regexp engine doesn't support.
var ownerPattern = regexp.MustCompile(`^[A-Za-z0-9]+(-[A-Za-z0-9]+)*$`)

// repoPattern approximates GitHub's repo name rules: alphanumeric, dot,
// underscore, hyphen.
var repoPattern = regexp.MustCompile(`^[A-Za-z0-9._-]+$`)

// tagVocabulary is the closed tag list, ordered for display. Tags describe a
// tool along orthogonal facets, so an entry carries several; the grouping is
// what lets the dashboard offer OR-within-facet, AND-across-facet filters.
//
// Deliberately excluded: model names (gpt-4, sonnet, r1 — they churn within
// months) and implementation stacks (rust, nextjs — they say nothing about
// choosing the tool).
//
// This slice is the single source of truth: the lookup map below, the
// validation messages, and the dashboard's filter chips (shipped in
// the generated dist/data.json) are all derived from it.
var tagVocabulary = []tagFacet{
	{ID: facetSurface, Label: "Surface", Tags: []string{"terminal", "editor-plugin", "ide", "desktop", "web", "self-hosted"}},
	{ID: facetModel, Label: "Model access", Tags: []string{"byo-model", "single-vendor", "local-models"}},
	{ID: facetWorkflow, Label: "Workflow", Tags: []string{"interactive", "autonomous", "review", "app-builder", "research", "orchestration"}},
	{ID: facetIntegration, Label: "Integration", Tags: []string{"mcp", "acp", "headless"}},
	{ID: facetOrigin, Label: "Origin", Tags: []string{"vendor", "community"}},
}

// tagFacet is one group of related tags, also the shape the dashboard consumes.
type tagFacet struct {
	ID    string   `json:"id"`
	Label string   `json:"label"`
	Tags  []string `json:"tags"`
}

const (
	facetSurface     = "surface"
	facetModel       = "model"
	facetWorkflow    = "workflow"
	facetIntegration = "integration"
	facetOrigin      = "origin"
)

// tagFacetIDs indexes tag -> facet ID for O(1) validation.
var tagFacetIDs = func() map[string]string {
	m := make(map[string]string)
	for _, f := range tagVocabulary {
		for _, t := range f.Tags {
			m[t] = f.ID
		}
	}
	return m
}()

// facetTags lists one facet's legal tags, for error messages that tell a
// contributor exactly what they may write.
func facetTags(id string) []string {
	for _, f := range tagVocabulary {
		if f.ID == id {
			return f.Tags
		}
	}
	return nil
}

// validateAgents checks data/agents.yml entries offline (no network, no
// token) and returns every violation found — not just the first — so a
// contributor sees the complete list of fixes needed in one pass.
func validateAgents(agents []Agent) []string {
	var violations []string
	seen := make(map[string]int, len(agents)) // lowercase "owner/repo" -> first index seen

	for i, a := range agents {
		ref := fmt.Sprintf("entry %d (owner=%q repo=%q)", i, a.Owner, a.Repo)

		switch {
		case strings.TrimSpace(a.Owner) == "":
			violations = append(violations, fmt.Sprintf("%s: owner is empty", ref))
		case !ownerPattern.MatchString(a.Owner):
			violations = append(violations, fmt.Sprintf("%s: owner %q does not look like a valid GitHub username/org (alphanumeric, single hyphens, no leading/trailing hyphen)", ref, a.Owner))
		}

		switch {
		case strings.TrimSpace(a.Repo) == "":
			violations = append(violations, fmt.Sprintf("%s: repo is empty", ref))
		case !repoPattern.MatchString(a.Repo):
			violations = append(violations, fmt.Sprintf("%s: repo %q contains characters not allowed in a GitHub repo name (allowed: letters, digits, '.', '_', '-')", ref, a.Repo))
		}

		if a.Category != "" {
			violations = append(violations, fmt.Sprintf("%s: category %q is no longer a field — replace it with tags, e.g. tags: [terminal, byo-model, interactive, community]", ref, a.Category))
		}
		violations = append(violations, validateTags(ref, a.Tags)...)

		key := strings.ToLower(a.Owner + "/" + a.Repo)
		if first, dup := seen[key]; dup {
			violations = append(violations, fmt.Sprintf("%s: duplicate of entry %d (case-insensitive owner/repo match)", ref, first))
		} else {
			seen[key] = i
		}
	}

	return violations
}

// runCheck loads data/agents.yml offline and validates it, printing every
// violation to stderr and returning a non-nil error if any are found. It
// never makes a network call, so it's safe to run against fork PRs without
// a GitHub token.
func runCheck(path string) error {
	agents, err := loadAgents(path)
	if err != nil {
		return fmt.Errorf("load %s: %w", path, err)
	}

	violations := validateAgents(agents)
	if len(violations) > 0 {
		for _, v := range violations {
			fmt.Fprintln(os.Stderr, v)
		}
		return fmt.Errorf("%d violation(s) found in %s", len(violations), path)
	}

	fmt.Printf("%s: %d agents valid\n", path, len(agents))
	return nil
}

// validateTags checks one entry's tags against the closed vocabulary: every
// tag known, no repeats, at least one surface tag (the dashboard groups rows
// by where you run them), and at most one origin tag (a tool has one
// publisher).
func validateTags(ref string, tags []string) []string {
	var violations []string
	seen := make(map[string]bool, len(tags))
	facetCount := map[string]int{}

	for _, tag := range tags {
		switch facet, known := tagFacetIDs[tag]; {
		case !known:
			violations = append(violations, fmt.Sprintf("%s: tag %q is not in the vocabulary (see docs/CONTRIBUTING.md)", ref, tag))
		case seen[tag]:
			violations = append(violations, fmt.Sprintf("%s: tag %q is repeated", ref, tag))
		default:
			seen[tag] = true
			facetCount[facet]++
		}
	}

	if facetCount[facetSurface] == 0 {
		violations = append(violations, fmt.Sprintf("%s: needs at least one surface tag (one of: %s)", ref, strings.Join(facetTags(facetSurface), ", ")))
	}
	if facetCount[facetOrigin] > 1 {
		violations = append(violations, fmt.Sprintf("%s: has %d origin tags, expected at most one (%s)", ref, facetCount[facetOrigin], strings.Join(facetTags(facetOrigin), ", ")))
	}

	return violations
}
