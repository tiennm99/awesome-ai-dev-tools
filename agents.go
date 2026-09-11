package main

import (
	"fmt"
	"os"

	"github.com/goccy/go-yaml"
)

type Agent struct {
	Owner string   `yaml:"owner"`
	Repo  string   `yaml:"repo"`
	Tags  []string `yaml:"tags"`
	Notes string   `yaml:"notes,omitempty"`

	// Category is the retired single-select field that tags replaced. It is
	// still parsed so a stale entry fails validation with a message naming
	// the replacement, rather than being silently ignored.
	Category string `yaml:"category,omitempty"`
}

func loadAgents(path string) ([]Agent, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var cfg struct {
		Agents []Agent `yaml:"agents"`
	}
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	for i, a := range cfg.Agents {
		if a.Owner == "" || a.Repo == "" {
			return nil, fmt.Errorf("entry %d missing owner or repo", i)
		}
	}
	return cfg.Agents, nil
}
