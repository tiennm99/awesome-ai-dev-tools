# Contributing

## Adding an Agent

Edit [`data/agents.yml`](../data/agents.yml) and add an entry:

```yaml
agents:
  - owner: github-username-or-org
    repo: repository-name
    tags: [terminal, byo-model, interactive, community]
```

Required fields: `owner`, `repo`, `tags`

Optional fields: `notes` (for clarifications or caveats)

## Field Reference

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `owner` | string | Yes | GitHub user or organization that owns the repo |
| `repo` | string | Yes | Repository name on GitHub |
| `tags` | list | Yes | Tags from the vocabulary below; at least one surface tag, at most one origin tag |
| `notes` | string | No | Additional context or disclaimers |

## Scope

Two kinds of project belong here:

1. **Coding agents** — they write, edit, or review code themselves.
2. **Agent development environments (ADEs)** — their primary purpose is running and
   coordinating those agents: parallel worktrees, session management, remote or mobile
   control. Tag these `orchestration`.

An ADE earns a row because it is the surface a developer actually works in, the same way
a coding agent is. What stays out is anything that only *assists* agents without being a
place you run them: libraries and SDKs, prompt or skill collections, dashboards and
observability-only layers, and single-purpose wrappers around one agent's config.

The other two criteria in the [README](../README.md#contributing) — roughly 10,000+ stars
and active maintenance — apply to both kinds equally.

## Tag Vocabulary

Tags replaced the old single-select `category` field, because one slot cannot
describe a tool that ships as a CLI, an editor plugin and a desktop app at the
same time — which most of them now do. An entry carries several tags across
five facets. `data/agents.yml` is the source of truth; the vocabulary itself
lives in `tagVocabulary` (`validate.go`) and reaches the dashboard through
`site/data.json`, so it is defined exactly once.

**Surface** — where you run it. At least one required.

- **terminal** — a CLI or TUI you run in a shell
- **editor-plugin** — extension for an existing editor (VS Code, JetBrains, Neovim)
- **ide** — a standalone editor or IDE
- **desktop** — a native or Electron/Tauri desktop app
- **web** — runs in a browser, hosted or local
- **self-hosted** — a server you deploy, with clients or editor plugins on top

**Model access** — which models it can drive.

- **byo-model** — bring your own: multiple providers, OpenAI-compatible endpoints, or OpenRouter
- **single-vendor** — built for one lab's models
- **local-models** — runs against local inference (Ollama, llama.cpp, vLLM, SGLang)

**Workflow** — how you work with it.

- **interactive** — conversational pair programming, you stay in the loop
- **autonomous** — takes a goal or issue and runs long stretches unattended
- **review** — reviews diffs or existing code rather than writing it
- **app-builder** — prompt-to-app scaffolding, with preview and deploy
- **research** — published as a research artifact or proof of concept
- **orchestration** — runs and coordinates *other* coding agents as its primary purpose (an ADE), rather than editing code itself

**Integration** — what it plugs into.

- **mcp** — speaks Model Context Protocol
- **acp** — speaks Agent Client Protocol
- **headless** — a non-interactive mode for scripting or CI

**Origin** — who publishes it. At most one.

- **vendor** — first-party tool from a model lab (Anthropic, OpenAI, Google, xAI, DeepSeek, Alibaba, Moonshot, …)
- **community** — everyone else

### The evidence rule

Apply a tag only when the repo's own README, docs, or GitHub topics support it.
Do not tag from reputation or from a blog post. Under-tagging is better than a
wrong tag: a missing tag hides a row from one filter, a wrong one sends someone
to a tool that cannot do what they need.

Deliberately **not** tags: model names (`gpt-4`, `sonnet`, `r1`) because they
churn within months, and implementation stacks (`rust`, `nextjs`) because they
say nothing about choosing the tool. GitHub topics are a drafting aid only —
12 of the tracked repos have no topics at all, including several in the top ten.

## Handling Duplicates and Changes

**Duplicate repos:** CI rejects them. `go run . -check` fails on a case-insensitive `owner/repo` match, so a duplicate never reaches a daily run. It also rejects unknown tags, a repeated tag, an entry with no surface tag, two origin tags, and any leftover `category:` key.

**Renamed repos:** GitHub redirects the old slug, so the updater keeps working — it prints a `::warning::` naming the new slug. Update `owner`/`repo` in `data/agents.yml` to the new slug and add the old key to `canonicalKeyMigrations` in `history.go`, or the repo's star history detaches and its deltas show `—`.

**Deprecation:** To remove an agent, delete its entry from `data/agents.yml`. The next run drops it from the README; its history stays in `data/history.jsonl`, so re-adding the entry later restores its star chart.

**Staleness:** An entry with no push in **6 months** is dropped. Past **3 months** the updater prints a `::warning::` naming the repo and its days idle, so the daily run surfaces candidates without anyone auditing the list by hand. Removal stays a human decision: a repo can go quiet between releases, and a historically significant one (`gpt-engineer`) is kept with a `notes` marker instead.

## PR Review

- Keep PRs to changes in `data/agents.yml` only (do not edit `README.md` or `data/history.jsonl`)
- The daily GitHub Actions workflow (runs at 00:00 UTC) picks up merged PRs automatically
- No manual review required; the updater regenerates the README after your PR merges

For local testing before opening a PR, see [LOCAL_DEV.md](./LOCAL_DEV.md).
