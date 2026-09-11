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
	readmeTmpl   = "templates/readme.tmpl"
	readmePath   = "README.md"
	siteDataPath = "site/data.json"
)

func main() {
	check := flag.Bool("check", false, "validate data/agents.yml offline (no network, no token) and exit")
	flag.Parse()

	if *check {
		if err := runCheck(agentsPath); err != nil {
			log.Printf("check failed: %v", err)
			os.Exit(1)
		}
		return
	}

	if err := run(); err != nil {
		log.Fatalf("update failed: %v", err)
	}
}

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

	snapshots, deltas7, deltas30, err := appendHistory(historyPath, stats)
	if err != nil {
		return err
	}

	if err := renderReadme(readmeTmpl, readmePath, stats, deltas7); err != nil {
		return err
	}

	if err := writeSiteData(siteDataPath, stats, deltas7, deltas30, snapshots); err != nil {
		return err
	}

	fmt.Printf("updated %d agents\n", len(stats))
	return nil
}
