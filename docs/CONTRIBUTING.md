# Contributing

## Adding an Agent

Edit [`data/agents.yml`](../data/agents.yml) and add an entry:

```yaml
agents:
  - owner: github-username-or-org
    repo: repository-name
    category: cli
```

Required fields: `owner`, `repo`, `category`

Optional fields: `notes` (for clarifications or caveats)

## Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `owner` | string | Yes | GitHub user or organization that owns the repo |
| `repo` | string | Yes | Repository name on GitHub |
| `category` | string | Yes | One of the values listed below |
| `notes` | string | No | Additional context or disclaimers |

## Valid Categories

- **cli** — Command-line tools and CLI wrappers
- **ide** — Standalone editors and IDEs
- **extension** — Editor extensions (VS Code, Neovim, etc.)
- **library** — Libraries and SDKs
- **research** — Research papers and proof-of-concept projects
- **web** — Web-based tools and online IDEs

## Handling Duplicates and Changes

**Duplicate repos:** CI rejects them. `go run . -check` fails on a case-insensitive `owner/repo` match, so a duplicate never reaches a daily run.

**Renamed repos:** GitHub redirects the old slug, so the updater keeps working — it prints a `::warning::` naming the new slug. Update `owner`/`repo` in `data/agents.yml` to the new slug and add the old key to `canonicalKeyMigrations` in `history.go`, or the repo's star history detaches and its deltas show `—`.

**Deprecation:** To remove an agent, delete its entry from `data/agents.yml`. The next run drops it from the README; its history stays in `data/history.jsonl`, so re-adding the entry later restores its star chart.

**Staleness:** An entry with no push in **6 months** is dropped. Past **3 months** the updater prints a `::warning::` naming the repo and its days idle, so the daily run surfaces candidates without anyone auditing the list by hand. Removal stays a human decision: a repo can go quiet between releases, and a historically significant one (`gpt-engineer`) is kept with a `notes` marker instead.

## PR Review

- Keep PRs to changes in `data/agents.yml` only (do not edit `README.md` or `data/history.jsonl`)
- The daily GitHub Actions workflow (runs at 00:00 UTC) picks up merged PRs automatically
- No manual review required; the updater regenerates the README after your PR merges

For local testing before opening a PR, see [LOCAL_DEV.md](./LOCAL_DEV.md).
