package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

// Paths the updater reads and writes, relative to the repository root.
const (
	agentsPath   = "data/agents.yml"
	historyPath  = "data/history.jsonl"
	metadataPath = "data/metadata.json"
	readmeTmpl   = "templates/readme.tmpl"
	readmePath   = "README.md"
	siteDir      = "site"
	distDir      = "dist"
)

// The tool has three modes, deliberately separated so the only step that needs
// network access and a GITHUB_TOKEN is the one that runs in CI:
//
//	(default)  update — fetch GitHub, refresh history/metadata/README
//	-build            — render dist/ from committed data, offline
//	-check            — validate data/agents.yml, offline
//
// The split is what lets an untrusted build environment (Cloudflare Pages)
// publish the site without ever holding a token.
func main() {
	check := flag.Bool("check", false, "validate data/agents.yml offline (no network, no token) and exit")
	build := flag.Bool("build", false, "render the static site into dist/ from committed data (no network, no token) and exit")
	flag.Parse()

	if *check && *build {
		log.Print("-check and -build are mutually exclusive")
		os.Exit(1)
	}

	switch {
	case *check:
		if err := runCheck(agentsPath); err != nil {
			log.Printf("check failed: %v", err)
			os.Exit(1)
		}
	case *build:
		if err := runBuild(agentsPath, metadataPath, historyPath, siteDir, distDir); err != nil {
			log.Printf("build failed: %v", err)
			os.Exit(1)
		}
	default:
		if err := run(); err != nil {
			log.Fatalf("update failed: %v", err)
		}
	}
}

// run is the update step: it is the only mode that talks to GitHub. It writes
// data/metadata.json, which the build step later consumes offline.
func run() error {
	agents, err := loadAgents(agentsPath)
	if err != nil {
		return err
	}
	if len(agents) == 0 {
		return fmt.Errorf("no agents in %s", agentsPath)
	}

	token := os.Getenv("GITHUB_TOKEN")
	if token == "" {
		return fmt.Errorf("GITHUB_TOKEN env var required")
	}

	stats, err := fetchStats(token, agents)
	if err != nil {
		return err
	}

	stats = enforceStarFloor(stats)
	if len(stats) == 0 {
		return fmt.Errorf("no agents in %s meet the %d-star minimum", agentsPath, minStars)
	}

	_, deltas7, _, err := appendHistory(historyPath, stats)
	if err != nil {
		return err
	}

	if err := renderReadme(readmeTmpl, readmePath, stats, deltas7); err != nil {
		return err
	}

	if err := writeMetadata(metadataPath, stats); err != nil {
		return err
	}

	fmt.Printf("updated %d agents\n", len(stats))
	return nil
}
