package main

import (
	"os"
	"strings"
	"testing"
)

func TestValidateAgents_ValidEntriesNoViolations(t *testing.T) {
	agents := []Agent{
		{Owner: "aider-ai", Repo: "aider", Tags: []string{"terminal", "community"}},
		{Owner: "cline", Repo: "cline", Tags: []string{"terminal", "community"}},
		{Owner: "a", Repo: "b.c-d_e", Tags: []string{"terminal", "community"}},
	}
	violations := validateAgents(agents)
	if len(violations) != 0 {
		t.Errorf("expected no violations, got %v", violations)
	}
}

func TestValidateAgents_EmptyOwner(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "", Repo: "repo1", Tags: []string{"terminal", "community"}}})
	if !anyContains(violations, "owner is empty") {
		t.Errorf("expected 'owner is empty' violation, got %v", violations)
	}
}

func TestValidateAgents_InvalidOwnerChars(t *testing.T) {
	tests := []string{"foo_bar", "-leading", "trailing-", "double--hyphen", "foo owner"}
	for _, owner := range tests {
		t.Run(owner, func(t *testing.T) {
			violations := validateAgents([]Agent{{Owner: owner, Repo: "repo1", Tags: []string{"terminal", "community"}}})
			if !anyContains(violations, "does not look like a valid GitHub username/org") {
				t.Errorf("owner %q: expected invalid-owner violation, got %v", owner, violations)
			}
		})
	}
}

func TestValidateAgents_EmptyRepo(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "", Tags: []string{"terminal", "community"}}})
	if !anyContains(violations, "repo is empty") {
		t.Errorf("expected 'repo is empty' violation, got %v", violations)
	}
}

func TestValidateAgents_InvalidRepoChars(t *testing.T) {
	tests := []string{"foo/bar", "foo bar", "foo@bar", "foo#bar"}
	for _, repo := range tests {
		t.Run(repo, func(t *testing.T) {
			violations := validateAgents([]Agent{{Owner: "org", Repo: repo, Tags: []string{"terminal", "community"}}})
			if !anyContains(violations, "characters not allowed in a GitHub repo name") {
				t.Errorf("repo %q: expected invalid-repo violation, got %v", repo, violations)
			}
		})
	}
}

func TestValidateAgents_LegacyCategoryFieldIsRejected(t *testing.T) {
	// An entry left on the retired schema must fail loudly, not be ignored.
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Category: "cli", Tags: []string{"terminal"}}})
	if !anyContains(violations, "no longer a field") {
		t.Errorf("expected a legacy-category violation, got %v", violations)
	}
}

func TestValidateAgents_MissingSurfaceTag(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Tags: []string{"byo-model", "community"}}})
	if !anyContains(violations, "at least one surface tag") {
		t.Errorf("expected a missing-surface violation, got %v", violations)
	}
}

func TestValidateAgents_NoTagsAtAll(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1"}})
	if !anyContains(violations, "at least one surface tag") {
		t.Errorf("expected a missing-surface violation, got %v", violations)
	}
}

func TestValidateAgents_UnknownTag(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Tags: []string{"terminal", "gpt-4"}}})
	if !anyContains(violations, `tag "gpt-4" is not in the vocabulary`) {
		t.Errorf("expected an unknown-tag violation, got %v", violations)
	}
}

func TestValidateAgents_RepeatedTag(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Tags: []string{"terminal", "terminal"}}})
	if !anyContains(violations, `tag "terminal" is repeated`) {
		t.Errorf("expected a repeated-tag violation, got %v", violations)
	}
}

func TestValidateAgents_TwoOriginTags(t *testing.T) {
	violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Tags: []string{"terminal", "vendor", "community"}}})
	if !anyContains(violations, "origin tags, expected at most one") {
		t.Errorf("expected a multiple-origin violation, got %v", violations)
	}
}

func TestValidateAgents_EveryVocabularyTagAccepted(t *testing.T) {
	// Each legal tag must pass alongside a surface tag, so a typo in
	// tagFacets can't silently make a documented tag unusable.
	for tag := range tagFacetIDs {
		tags := []string{"terminal", tag}
		if tag == "terminal" {
			tags = []string{"terminal"}
		}
		if violations := validateAgents([]Agent{{Owner: "org", Repo: "repo1", Tags: tags}}); len(violations) > 0 {
			t.Errorf("tag %q: expected no violations, got %v", tag, violations)
		}
	}
}

func TestValidateAgents_DuplicateCaseInsensitive(t *testing.T) {
	agents := []Agent{
		{Owner: "Foo", Repo: "Bar", Tags: []string{"terminal", "community"}},
		{Owner: "foo", Repo: "bar", Tags: []string{"terminal", "community"}},
	}
	violations := validateAgents(agents)
	if !anyContains(violations, "duplicate of entry 0") {
		t.Errorf("expected duplicate violation referencing entry 0, got %v", violations)
	}
}

func TestValidateAgents_NoDuplicateForDistinctRepos(t *testing.T) {
	agents := []Agent{
		{Owner: "foo", Repo: "bar", Tags: []string{"terminal", "community"}},
		{Owner: "foo", Repo: "baz", Tags: []string{"terminal", "community"}},
	}
	violations := validateAgents(agents)
	if len(violations) != 0 {
		t.Errorf("expected no violations for distinct repos, got %v", violations)
	}
}

func TestValidateAgents_CollectsAllViolationsNotJustFirst(t *testing.T) {
	agents := []Agent{
		{Owner: "", Repo: "", Tags: nil},
		{Owner: "org", Repo: "repo1", Tags: []string{"terminal", "nonsense", "vendor", "community"}},
	}
	violations := validateAgents(agents)
	// entry 0: owner empty, repo empty, no surface tag = 3 violations.
	// entry 1: unknown tag, two origin tags = 2 violations.
	if len(violations) != 5 {
		t.Errorf("expected 5 violations collected across both entries, got %d: %v", len(violations), violations)
	}
}

func TestRunCheck_CurrentAgentsYML(t *testing.T) {
	// The real data/agents.yml must always pass -check; this is the
	// regression guard for that invariant.
	if err := runCheck("data/agents.yml"); err != nil {
		t.Errorf("expected data/agents.yml to pass validation, got: %v", err)
	}
}

func TestRunCheck_BadYAMLReportsViolationsAndError(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/bad-agents.yml"
	// loadAgents itself already rejects empty owner/repo, so this fixture
	// covers the violation classes only validateAgents catches: bad owner
	// chars, bad repo chars, invalid category, and a case-insensitive dup.
	content := `agents:
  - owner: "bad owner"
    repo: "valid-repo"
    category: cli
  - owner: "org"
    repo: "bad/repo"
    category: not-a-real-category
  - owner: "Dup"
    repo: "Repo"
    category: web
  - owner: "dup"
    repo: "repo"
    category: web
`
	if err := os.WriteFile(tmpFile, []byte(content), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	err := runCheck(tmpFile)
	if err == nil {
		t.Fatal("expected runCheck to fail on invalid agents.yml, got nil error")
	}
	if !strings.Contains(err.Error(), "violation(s) found") {
		t.Errorf("expected error to summarize violation count, got: %v", err)
	}
}

func TestRunCheck_MissingFile(t *testing.T) {
	err := runCheck("/nonexistent/path/agents.yml")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}

func anyContains(list []string, substr string) bool {
	for _, s := range list {
		if strings.Contains(s, substr) {
			return true
		}
	}
	return false
}
