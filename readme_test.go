package main

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestSanitizeCell(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
		desc     string
	}{
		{
			name:     "plain text",
			input:    "hello world",
			expected: "hello world",
			desc:     "unchanged",
		},
		{
			name:     "pipe char",
			input:    "foo | bar",
			expected: "foo \\| bar",
			desc:     "pipe escaped",
		},
		{
			name:     "newline",
			input:    "line1\nline2",
			expected: "line1 line2",
			desc:     "newline becomes space",
		},
		{
			name:     "carriage return",
			input:    "line1\rline2",
			expected: "line1 line2",
			desc:     "carriage return becomes space",
		},
		{
			name:     "pipe and newline",
			input:    "foo | bar\nbaz",
			expected: "foo \\| bar baz",
			desc:     "both handled correctly",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
			desc:     "empty string unchanged",
		},
		{
			name:     "multiple pipes",
			input:    "a | b | c",
			expected: "a \\| b \\| c",
			desc:     "multiple pipes escaped",
		},
		{
			name:     "multiple newlines",
			input:    "line1\n\nline2",
			expected: "line1  line2",
			desc:     "multiple newlines become spaces",
		},
		{
			name:     "whitespace trimming",
			input:    "  text  ",
			expected: "text",
			desc:     "leading/trailing spaces trimmed",
		},
		{
			name:     "complex: whitespace, pipe, newline",
			input:    "  foo | bar\nbaz  ",
			expected: "foo \\| bar baz",
			desc:     "all transformations applied",
		},
		{
			name:     "only newlines and pipes",
			input:    "|\n|",
			expected: "\\| \\|",
			desc:     "pipes and newlines only",
		},
		{
			name:     "mixed line endings",
			input:    "line1\nline2\rline3",
			expected: "line1 line2 line3",
			desc:     "both \\n and \\r converted",
		},
		{
			name:     "backslash then pipe",
			input:    `\|`,
			expected: `\\\|`,
			desc:     "backslash escaped first so it can't neutralize the pipe escape (finding: literal backslash + unescaped pipe previously broke the table row)",
		},
		{
			name:     "lone backslash",
			input:    `a\b`,
			expected: `a\\b`,
			desc:     "backslash doubled so Markdown doesn't interpret it as an escape",
		},
		{
			name:     "raw html img tag",
			input:    "<img src=x onerror=alert(1)>",
			expected: "&lt;img src=x onerror=alert(1)&gt;",
			desc:     "angle brackets entity-escaped so raw HTML can't be injected from a third-party description",
		},
		{
			name:     "raw html script tag",
			input:    "<script>alert(1)</script>",
			expected: "&lt;script&gt;alert(1)&lt;/script&gt;",
			desc:     "script tags neutralized via entity escaping",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeCell(tt.input)
			if result != tt.expected {
				t.Errorf("%s: expected %q, got %q", tt.desc, tt.expected, result)
			}
		})
	}
}

func TestRenderReadme_TableRowsAndCallout(t *testing.T) {
	fixed := time.Date(2026, 9, 11, 3, 39, 0, 0, time.UTC)
	orig := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = orig }()

	stats := []Stat{
		{CanonicalKey: "org/big", NameWithOwner: "org/big", URL: "https://github.com/org/big",
			Stars: 1_500_000, Language: "Go", PushedAt: fixed, Description: "huge | agent", Tags: []string{"terminal", "community"}},
		{CanonicalKey: "org/mid", NameWithOwner: "org/mid", URL: "https://github.com/org/mid",
			Stars: 1234, Language: "Rust", PushedAt: fixed, Description: "mid agent", Tags: []string{"terminal", "community"}},
		{CanonicalKey: "org/small", NameWithOwner: "org/small", URL: "https://github.com/org/small",
			Stars: 999, Language: "", PushedAt: fixed, Description: "small agent", Tags: []string{"terminal", "community"}},
	}
	// org/small has no delta: its row must show the em dash, and it must not
	// win the top-mover callout.
	deltas := map[string]int{"org/big": 12, "org/mid": -3}

	out := t.TempDir() + "/README.md"
	if err := renderReadme("templates/readme.tmpl", out, stats, deltas); err != nil {
		t.Fatalf("renderReadme: %v", err)
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	got := string(raw)

	want := []string{
		"**Last updated:** 2026-09-11 03:39 UTC · **Tracked:** 3 repos",
		"**Top 7-day mover:** [org/big](https://github.com/org/big) (+12 stars)",
		"| 1 | [org/big](https://github.com/org/big) | 1.5M | +12 | Go | 2026-09-11 | huge \\| agent |",
		"| 2 | [org/mid](https://github.com/org/mid) | 1.2k | -3 | Rust | 2026-09-11 | mid agent |",
		"| 3 | [org/small](https://github.com/org/small) | 999 | — |  | 2026-09-11 | small agent |",
	}
	for _, w := range want {
		if !strings.Contains(got, w) {
			t.Errorf("README missing line:\n%s\n--- got ---\n%s", w, got)
		}
	}
}

func TestRenderReadme_NoDeltasOmitsTopMover(t *testing.T) {
	fixed := time.Date(2026, 9, 11, 3, 39, 0, 0, time.UTC)
	orig := timeNow
	timeNow = func() time.Time { return fixed }
	defer func() { timeNow = orig }()

	stats := []Stat{{CanonicalKey: "org/repo", NameWithOwner: "org/repo",
		URL: "https://github.com/org/repo", Stars: 10, PushedAt: fixed, Tags: []string{"terminal", "community"}}}

	out := t.TempDir() + "/README.md"
	if err := renderReadme("templates/readme.tmpl", out, stats, map[string]int{}); err != nil {
		t.Fatalf("renderReadme: %v", err)
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if strings.Contains(string(raw), "Top 7-day mover") {
		t.Error("top-mover callout rendered with no deltas available")
	}
}

func TestRenderReadme_MissingTemplateDoesNotTouchOutput(t *testing.T) {
	out := t.TempDir() + "/README.md"
	if err := os.WriteFile(out, []byte("previous"), 0o644); err != nil {
		t.Fatalf("seed file: %v", err)
	}
	if err := renderReadme("templates/does-not-exist.tmpl", out, nil, nil); err == nil {
		t.Fatal("expected an error for a missing template")
	}
	raw, err := os.ReadFile(out)
	if err != nil {
		t.Fatalf("ReadFile: %v", err)
	}
	if string(raw) != "previous" {
		t.Errorf("output = %q, want the existing README left intact", raw)
	}
}
