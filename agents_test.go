package main

import (
	"os"
	"strings"
	"testing"
)

func TestLoadAgents_ValidationErrors(t *testing.T) {
	tests := []struct {
		name      string
		yaml      string
		expectErr bool
		errSubstr string
		desc      string
	}{
		{
			name: "valid entry",
			yaml: `agents:
  - owner: org
    repo: repo1
    tags: [terminal, community]
`,
			expectErr: false,
			desc:      "should load valid entry",
		},
		{
			name: "missing owner",
			yaml: `agents:
  - repo: repo1
`,
			expectErr: true,
			errSubstr: "entry 0 missing owner or repo",
			desc:      "should fail on missing owner",
		},
		{
			name: "missing repo",
			yaml: `agents:
  - owner: org
`,
			expectErr: true,
			errSubstr: "entry 0 missing owner or repo",
			desc:      "should fail on missing repo",
		},
		{
			name: "empty owner",
			yaml: `agents:
  - owner: ""
    repo: repo1
`,
			expectErr: true,
			errSubstr: "entry 0 missing owner or repo",
			desc:      "should fail on empty owner",
		},
		{
			name: "empty repo",
			yaml: `agents:
  - owner: org
    repo: ""
`,
			expectErr: true,
			errSubstr: "entry 0 missing owner or repo",
			desc:      "should fail on empty repo",
		},
		{
			name: "valid multiple entries",
			yaml: `agents:
  - owner: org1
    repo: repo1
  - owner: org2
    repo: repo2
    tags: [terminal, community]
`,
			expectErr: false,
			desc:      "should load multiple valid entries",
		},
		{
			name: "second entry missing owner",
			yaml: `agents:
  - owner: org1
    repo: repo1
  - repo: repo2
`,
			expectErr: true,
			errSubstr: "entry 1 missing owner or repo",
			desc:      "should report correct index for failing entry",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := tmpDir + "/agents.yml"

			if err := os.WriteFile(tmpFile, []byte(tt.yaml), 0600); err != nil {
				t.Fatalf("WriteFile failed: %v", err)
			}

			agents, err := loadAgents(tmpFile)

			if tt.expectErr {
				if err == nil {
					t.Errorf("%s: expected error, got nil", tt.desc)
				} else if !strings.Contains(err.Error(), tt.errSubstr) {
					t.Errorf("%s: expected error containing %q, got %q", tt.desc, tt.errSubstr, err.Error())
				}
				if agents != nil {
					t.Errorf("%s: expected nil agents on error, got %v", tt.desc, agents)
				}
			} else {
				if err != nil {
					t.Errorf("%s: expected no error, got %v", tt.desc, err)
				}
				if agents == nil {
					t.Errorf("%s: expected non-nil agents", tt.desc)
				}
			}
		})
	}
}

func TestLoadAgents_EmptyList(t *testing.T) {
	// Empty agents list should load successfully (edge case).
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/agents.yml"

	yaml := `agents: []`
	if err := os.WriteFile(tmpFile, []byte(yaml), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	agents, err := loadAgents(tmpFile)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if agents == nil || len(agents) != 0 {
		t.Errorf("expected empty agents slice, got %v", agents)
	}
}

func TestLoadAgents_MalformedYAML(t *testing.T) {
	// Malformed YAML should return a parse error.
	tmpDir := t.TempDir()
	tmpFile := tmpDir + "/agents.yml"

	yaml := `agents:
  - owner: org
    repo: repo1
  - this is not valid yaml: {`

	if err := os.WriteFile(tmpFile, []byte(yaml), 0600); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	agents, err := loadAgents(tmpFile)
	if err == nil {
		t.Errorf("expected error on malformed YAML, got nil")
	}
	if agents != nil {
		t.Errorf("expected nil agents on parse error, got %v", agents)
	}
}
