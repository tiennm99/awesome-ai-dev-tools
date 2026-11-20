# Research Report: audit of every entry in data/agents.yml

**Conducted:** 2026-11-20 (+07) · **Repo:** `tiennm99/awesome-ai-dev-tools` · **Entries audited:** 42

## Executive summary

Of the 42 entries, 21 stay as they are and 21 need action.

- **Rename (2):** `Gitlawb/openclaude` is now `Twigpine/openclaude`, and `Hmbown/Codewhale` is now `codewhale-hq/Codewhale`. Both old keys appear in `data/history.jsonl` (71 snapshots each), so both need a `canonicalKeyMigrations` entry.
- **Remove (3):** `esengine/DeepSeek-Reasonix` was archived on 2026-11-12 with an end-of-maintenance notice. `MoonshotAI/kimi-cli` is archived. `Pythagora-io/gpt-pilot` says it is no longer maintained. Removing the last two reverses the keep decision in commit `98e5b90` (2026-09-29), so the user has to confirm it.
- **Update (16, plus the Codewhale retag):** the changes are tags, descriptions and notes. The ones that change what a reader learns:
  - gemini-cli now serves only enterprise Code Assist users.
  - OpenHands is now "Agent Canvas", an orchestrator.
  - langchain-ai/open-swe is now "Robin".
  - Onlook is now desktop-only.
  - The Codewhale `desktop` tag is wrong: its desktop app has no download.
  - grok-build and oh-my-pi carry `editor-plugin`, but neither ships an editor extension. They reach editors only over ACP.
  - Aider passes six months idle on 2026-11-22, so it needs a decision.
  - Zed's description opens with the marketing adjective "High-performance". claurst and t3code lack an `mcp` tag their docs support.

No entry is below the 1,000-star floor. The lowest is `Kuberwastaken/claurst` at 10,349.

I applied every proposed block and both migrations to a scratch copy of the repo. `go run . -check` passed ("39 agents valid") and `go test ./...` passed. No repo file was modified.

## Method

- **Status:** `gh api repos/OWNER/REPO` for all 42 slugs, run live on 2026-11-20. It returned stars, `pushed_at`, `archived`, topics and the redirect (`full_name`). I cross-checked this against `data/metadata.json` (fetched 2026-11-20 06:15 UTC), and the two agree.
- **Content:** I fetched every README through `gh api repos/.../readme`. I read in full the READMEs of entries that changed or looked likely to, and the head of every other one. For all 42, I grepped for MCP, ACP, headless, desktop, editor, local-model and deprecation terms and read each matching line. Where the README was silent, I read in-repo docs (grok-build user guide, prime-agent `crates/pa-cli/README.md`, MiMo-Code `LICENSE`) or the vendor's docs page (Open Interpreter providers).
- **Staleness:** the warning cutoff is 90 days (`staleWarnAfter` in `github.go`). Removal is at 6 calendar months, which puts the cutoff at a push before 2026-05-20.
- **Description rules:** I checked every description by script: 140 characters at most, no trailing period, no `|`, and no model version names. I read every description for marketing adjectives. Only Zed's "High-performance" fails. All proposed descriptions pass (the longest is 131).
- **Evidence rule:** I added a tag only where there is a quotable primary-source line. Where a short README was silent, I left existing tags alone (codex, claude-code, zed, dyad, deepseek-harness).

## Verdict table

| # | Entry | ★ (live) | Last push | Verdict | Reason |
|--:|---|--:|---|---|---|
| 1 | Aider-AI/aider | 50,193 | 2026-05-22 | UPDATE (decision) | 182 days idle and over six months on 2026-11-22. Last tag v0.86.2 on 2026-02-12. Recommend keeping it with a historical-significance note, as gpt-engineer has |
| 2 | cline/cline | 71,790 | 2026-11-19 | KEEP | IDE, terminal and desktop, plus a headless CLI for CI/CD; all tags confirmed |
| 3 | anomalyco/opencode | 228,155 | 2026-11-20 | KEEP | Terminal UI and desktop app (beta) confirmed |
| 4 | charmbracelet/crush | 29,147 | 2026-11-20 | KEEP | Switches models mid-session, LSP, MCP; matches |
| 5 | google-gemini/gemini-cli | 108,953 | 2026-11-20 | UPDATE | README now says enterprise Code Assist only; free and Google AI plan users moved to Antigravity CLI on 2026-06-18 |
| 6 | openai/codex | 135,962 | 2026-11-20 | KEEP | Terminal plus IDE install path; tags unchanged |
| 7 | aaif-goose/goose | 56,813 | 2026-11-20 | KEEP (scope flag) | Facts accurate (desktop, CLI, ACP, AAIF). README now describes it as a "general-purpose agent" and leads with the oi-kit library; see unresolved questions |
| 8 | AntonOsika/gpt-engineer | 54,949 | 2025-05-14 | KEEP | Archived; the historical-significance exception is named in CONTRIBUTING |
| 9 | SWE-agent/SWE-agent | 20,837 | 2026-11-16 | KEEP | README still recommends mini-swe-agent; the note is accurate |
| 10 | TabbyML/tabby | 33,930 | 2026-06-30 | KEEP (warn) | 143 days idle; passes six months on 2026-12-30 |
| 11 | zed-industries/zed | 95,387 | 2026-11-20 | UPDATE | Description opens with the marketing adjective "High-performance"; tags unchanged |
| 12 | anthropics/claude-code | 156,403 | 2026-11-20 | KEEP | Unchanged; npm install deprecated, which does not affect tags |
| 13 | avante-corp/avante.nvim | 18,203 | 2026-11-20 | KEEP | Cursor-style Neovim plugin, ACP and MCP confirmed |
| 14 | earendil-works/pi | 123,020 | 2026-11-20 | UPDATE | README reframed as "minimal, extensible agent harness"; documents print/JSON mode, so add `headless` |
| 15 | OpenHands/OpenHands | 93,892 | 2026-11-20 | UPDATE | Repo is now the Agent Canvas control center that runs OpenHands or any ACP agent: add `orchestration`. It ships no terminal agent, so drop `terminal` |
| 16 | warpdotdev/warp | 66,064 | 2026-11-20 | KEEP | AGPL-3.0 confirmed by the license API; built-in agent plus bring-your-own CLI agents. The "Apr 2026" open-sourcing date in the note is unverified |
| 17 | Pythagora-io/gpt-pilot | 33,642 | 2026-06-18 | REMOVE (decision) | README: "no longer actively maintained" and "not being maintained anymore". 155 days idle; over six months on 2026-12-18 |
| 18 | QwenLM/qwen-code | 30,809 | 2026-11-20 | KEEP | Terminal, editor, desktop, browser, local models; matches |
| 19 | Kilo-Org/kilocode | 28,063 | 2026-11-20 | KEEP | VS Code, JetBrains, CLI; matches |
| 20 | onlook-dev/onlook | 27,538 | 2026-11-20 | UPDATE | Now a desktop app that drives your Claude Code or Codex; web version moved to the `v1` branch, so drop `web`; React-only framing outdated |
| 21 | dyad-sh/dyad | 22,332 | 2026-11-19 | KEEP | Local app builder, bring-your-own keys; matches |
| 22 | github/copilot-cli | 11,441 | 2026-11-19 | KEEP | MCP and per-action approval confirmed |
| 23 | deepseek-ai/deepseek-harness | 284,533 | 2026-11-18 | KEEP | Still a developer preview with breaking changes; the note is accurate |
| 24 | openinterpreter/openinterpreter | 68,922 | 2026-11-20 | UPDATE | Add `headless`: it supports Codex `exec` and `interpreter exec`. Ollama and LM Studio are documented, so `local-models` holds |
| 25 | Hmbown/Codewhale | 41,186 | 2026-11-20 | RENAME + UPDATE | Now `codewhale-hq/Codewhale`. Desktop app "in development, no download", so drop `desktop`; add `web`, `headless`, `acp` |
| 26 | esengine/DeepSeek-Reasonix | 36,093 | 2026-11-12 | REMOVE | Archived; "no longer maintained as of 2026-11-12" |
| 27 | Gitlawb/openclaude | 34,440 | 2026-11-20 | RENAME | Now `Twigpine/openclaude`; npm scope moved too; content unchanged |
| 28 | can1357/oh-my-pi | 39,519 | 2026-11-20 | UPDATE | No editor extension; editors reach it only through `omp acp`, so drop `editor-plugin` |
| 29 | xai-org/grok-build | 27,814 | 2026-09-29 | UPDATE | README now names the vendor "SpaceXAI". Editors connect over ACP (Zed, Neovim), with no extension, so drop `editor-plugin` |
| 30 | alibaba/open-code-review | 54,527 | 2026-11-20 | KEEP | Diff-based line-level review CLI; matches |
| 31 | PrimeIntellect-ai/prime-agent | 22,857 | 2026-11-20 | UPDATE | `headless` confirmed in CLI docs; add `mcp` (MCP server config, `/mcp login`, MCP docs) |
| 32 | 1jehuang/jcode | 21,937 | 2026-11-20 | KEEP | Matches |
| 33 | HKUDS/DeepCode | 16,855 | 2026-11-17 | KEEP | TUI, desktop and web over one local service; matches |
| 34 | CodebuffAI/freebuff | 14,249 | 2026-11-20 | UPDATE | Description names desktop and web but the tags do not; add `desktop`, `web`, `app-builder` |
| 35 | MoonshotAI/kimi-cli | 11,404 | 2026-09-22 | REMOVE (decision) | Archived; README says installs "will stop working". Successor `MoonshotAI/kimi-code` has 8,270 ★ |
| 36 | Kuberwastaken/claurst | 10,349 | 2026-10-26 | UPDATE | Beta, multi-provider, ACP, headless all confirmed. Its `docs/mcp.md` documents MCP servers, so add `mcp` |
| 37 | stablyai/orca | 130,472 | 2026-11-20 | UPDATE | README now says "any CLI agent" (Grok, Cursor, Copilot, Antigravity…), not the four named; update description and notes |
| 38 | getpaseo/paseo | 24,510 | 2026-11-20 | UPDATE | Notes miss Antigravity and Muse Code |
| 39 | XiaomiMiMo/MiMo-Code | 13,794 | 2026-11-20 | KEEP | "Fork of opencode" note confirmed by the `LICENSE` copyright line "(c) 2025 opencode" |
| 40 | The-PR-Agent/pr-agent | 16,474 | 2026-11-20 | KEEP | Community-maintained legacy of Qodo; the note is accurate |
| 41 | langchain-ai/open-swe | 11,242 | 2026-11-20 | UPDATE | README renames the product "Robin"; deploying it as-is requires LangSmith Managed Deep Agents |
| 42 | pingdotgg/t3code | 67,084 | 2026-11-20 | UPDATE | Agent list grew (Kiro, Pi, Muse Code, fx, any ACP agent). It drives ACP agents and runs an MCP server for outside agents, so add `acp` and `mcp` |

## Renames

Change `owner` in `data/agents.yml`, and add these to `canonicalKeyMigrations` in `history.go` (the formatting is gofmt-aligned):

```go
	"Gitlawb/openclaude":            "Twigpine/openclaude",
	"Hmbown/Codewhale":              "codewhale-hq/Codewhale",
```

Also extend the "Known renames" comment above the map: openclaude moved Gitlawb to Twigpine, and Codewhale moved Hmbown to codewhale-hq.

### Twigpine/openclaude (old key `Gitlawb/openclaude`)

```yaml
  - owner: Twigpine
    repo: openclaude
    description: "Terminal coding-agent CLI that works with cloud and local model providers, with MCP support and a VS Code extension"
    tags: [terminal, editor-plugin, byo-model, local-models, interactive, headless, mcp, community]
```

Evidence: <https://github.com/Twigpine/openclaude>. The API redirect sets `full_name` to `Twigpine/openclaude`. The README says "The package moved to the `@twigpine` scope and the old one is no longer updated".

### codewhale-hq/Codewhale (old key `Hmbown/Codewhale`)

```yaml
  - owner: codewhale-hq
    repo: Codewhale
    description: "Open-source terminal coding agent in Rust that works with hosted or local models and delegates work to sub-agents"
    tags: [terminal, web, editor-plugin, byo-model, local-models, interactive, headless, mcp, acp, community]
```

Evidence: <https://github.com/codewhale-hq/Codewhale>. The README backs each tag change:

- **`desktop` dropped:** "Desktop app | In development, no download".
- **`web` added:** "`codewhale web` | A local web client on `127.0.0.1`".
- **`headless` added:** "`codewhale exec` | One task without the interface, for scripts and CI".
- **`acp` added:** "`codewhale serve --acp` | The Agent Client Protocol over stdio, for editors such as Zed".
- **`editor-plugin` kept:** the README links a community-maintained VS Code extension that connects from the sidebar.

## Removals

Delete each of these blocks. Their history stays in `data/history.jsonl`.

| Entry | Evidence |
|---|---|
| `esengine/DeepSeek-Reasonix` | <https://github.com/esengine/DeepSeek-Reasonix> (archived: true). README: "Reasonix is no longer maintained as of 2026-11-12. The final release is v2.61.0". Announcement: <https://github.com/esengine/DeepSeek-Reasonix/issues/15109> |
| `MoonshotAI/kimi-cli` | <https://github.com/MoonshotAI/kimi-cli> (archived: true). README: "This project has been archived and is no longer maintained… Existing installations are no longer supported and will stop working" |
| `Pythagora-io/gpt-pilot` | <https://github.com/Pythagora-io/gpt-pilot>. README: "This repository is **no longer actively maintained**… a security cleanup, not a resumption of development" and "This repo is not being maintained anymore" |

The rule in `docs/CONTRIBUTING.md` (since `00c90c5`, 2026-09-16) removes archived or declared-deprecated repos unless they qualify for the historical-significance exception. Neither kimi-cli nor gpt-pilot carries that marker. Commit `98e5b90` annotated and kept both after the rule landed, so these two removals need the user's confirmation. If gpt-pilot is kept instead, its `notes` should say "kept for historical significance", as gpt-engineer's do. It will also go past six months idle on 2026-12-18.

Possible replacement, which this audit did not vet: `MoonshotAI/kimi-code` has 8,270 ★, was pushed 2026-11-20, and is the declared successor. It would need a full scope and tag review before being added.

## Updates

### Aider-AI/aider (decision)

```yaml
  - owner: Aider-AI
    repo: aider
    description: "Terminal pair-programming tool that maps your codebase and edits it with cloud or local LLMs"
    tags: [terminal, byo-model, local-models, interactive, community]
    notes: No push since May 2026; kept for historical significance
```

Evidence: <https://github.com/Aider-AI/aider/commits/main> shows the last commit at 2026-05-22T14:02Z. The newest tags are v0.86.2 and v0.86.3.dev, both on 2026-02-12 (<https://github.com/Aider-AI/aider/tags>). The newest GitHub Release is v0.86.0 from 2025-08-09, so the release page understates activity. The README has no maintenance statement.

The recommended option is to keep it with the note. It has 50k stars, was the early template for terminal pair-programmers, and the gpt-engineer precedent covers this case. The alternative is to delete it once the six months pass on 2026-11-22.

### google-gemini/gemini-cli

```yaml
  - owner: google-gemini
    repo: gemini-cli
    description: "Google's open-source terminal agent for Gemini models, now aimed at enterprise Gemini Code Assist users, with MCP support"
    tags: [terminal, single-vendor, interactive, headless, mcp, vendor]
    notes: Free and Google AI plan access moved to Antigravity CLI in Jun 2026; this repo now serves Gemini Code Assist Standard and Enterprise
```

Evidence: <https://github.com/google-gemini/gemini-cli>. The README notice reads "Enterprise Code Assist only" and "No free tier or Google AI plans: Free individual access and Google AI Pro / Ultra (Google One) moved to Antigravity CLI on June 18, 2026". Transition post: <https://developers.googleblog.com/an-important-update-transitioning-gemini-cli-to-antigravity-cli/>.

The old description listed "built-in search, file and shell tools". I dropped that to make room for the audience change.

### OpenHands/OpenHands

```yaml
  - owner: OpenHands
    repo: OpenHands
    description: "Agent Canvas, a self-hosted control center that runs OpenHands or other ACP agents locally, in Docker, or in the cloud"
    tags: [web, self-hosted, byo-model, autonomous, orchestration, acp, community]
```

Evidence: <https://github.com/OpenHands/OpenHands>.

- The title is "Agent Canvas: The self-hosted developer control center for coding agents and automations", which supports `orchestration`.
- The repository-boundaries table says this repo owns the "Agent Canvas frontend, user-facing control center, backend selection, and local-stack orchestration". The agent, SDK and server live in `OpenHands/software-agent-sdk`.
- `agent-canvas` only starts the stack, and the UI is at `localhost:8000`. That is why I drop `terminal`.

### earendil-works/pi

```yaml
  - owner: earendil-works
    repo: pi
    description: "Minimal, extensible terminal coding agent harness you customize with extensions, skills, and packages, with an SDK and RPC mode"
    tags: [terminal, byo-model, interactive, headless, community]
```

Evidence: <https://github.com/earendil-works/pi>. The README says "Pi is a minimal, extensible agent harness that you can make your own", and "Use Pi interactively, automate it in print or JSON mode, control it over RPC, or build apps with the Pi TypeScript SDK". The print/JSON mode supports `headless`.

### onlook-dev/onlook

```yaml
  - owner: onlook-dev
    repo: onlook
    description: "Open-source desktop design canvas that runs your real app and has your own Claude Code or Codex agent edit its code"
    tags: [desktop, byo-model, app-builder, mcp, community]
    notes: The earlier browser-based version lives on the repo's v1 branch
```

Evidence: <https://github.com/onlook-dev/onlook>. The README supports each change:

- **Desktop app:** "A desktop app, local-first".
- **Bring-your-own agent:** "We don't ship our own model. Connect the Claude Code or ChatGPT/Codex account you already have… Plug in more tools through MCP connectors".
- **`web` dropped:** "Onlook v1 (web): the browser-based version lives on the `v1` branch".
- **No longer React-only:** "it works with all web frameworks or languages your dev server already supports".
- **`app-builder` kept:** "Start entirely new codebases from scratch".

### openinterpreter/openinterpreter

```yaml
  - owner: openinterpreter
    repo: openinterpreter
    description: "Terminal coding agent written in Rust, tuned for low-cost open-weight models"
    tags: [terminal, byo-model, local-models, interactive, headless, mcp, acp, community]
```

Evidence:

- <https://github.com/openinterpreter/openinterpreter>: "Supports `exec`, MCP, skills, hooks…" and "`interpreter exec --chat-completions`", which support `headless`.
- <https://www.openinterpreter.com/docs/terminal/providers>: built-in `ollama` ("Ollama local OSS") and `lmstudio` providers, which confirm `local-models`.

### can1357/oh-my-pi

```yaml
  - owner: can1357
    repo: oh-my-pi
    description: "Terminal coding agent forked from Pi, with LSP and debugger (DAP) integration built in"
    tags: [terminal, byo-model, local-models, interactive, mcp, acp, headless, community]
```

Evidence: <https://github.com/can1357/oh-my-pi>. The README lists its "Four entry points: interactive, one-shot, RPC, and ACP", and describes running "omp inside Zed" through `omp acp`. It names no VS Code, JetBrains or Neovim extension. The `acp` tag already covers editor use, and `editor-plugin` ("extension for an existing editor") has no evidence.

### xai-org/grok-build

```yaml
  - owner: xai-org
    repo: grok-build
    description: "SpaceXAI's Grok Build terminal coding agent, a full-screen Rust TUI that also runs headless or inside editors over ACP"
    tags: [terminal, single-vendor, interactive, headless, mcp, acp, vendor]
```

Evidence:

- <https://github.com/xai-org/grok-build>: "Grok Build is SpaceXAI's terminal-based AI coding agent… interactively, headlessly for scripting/CI, or embedded in editors via the Agent Client Protocol (ACP)".
- The agent-mode guide in the repo user guide (<https://github.com/xai-org/grok-build/tree/main/crates/codegen/xai-grok-pager/docs/user-guide>) lists Zed and Neovim (CodeCompanion, avante.nvim) as ACP clients, with JetBrains "Coming soon". It ships no extension of its own, so I drop `editor-plugin`.
- The pushed date of 2026-09-29 is still inside the 90-day window. The README says the repo is "synced periodically from the SpaceXAI monorepo", so long gaps between pushes are expected.

### PrimeIntellect-ai/prime-agent

```yaml
  - owner: PrimeIntellect-ai
    repo: prime-agent
    description: "Terminal coding and research agent driven from a persistent Python REPL, with subagents and a self-refining harness"
    tags: [terminal, byo-model, autonomous, headless, mcp, vendor]
```

Evidence:

- <https://github.com/PrimeIntellect-ai/prime-agent/blob/main/crates/pa-cli/README.md> covers mode selection ("interactive/headless/json/daemon subcommands"), "The headless print/json modes", MCP servers in settings `mcpServers`, and `/mcp login`.
- The repo also has `docs/google-workspace-mcp.md` and `docs/slack-mcp.md`.
- The README says "run `/login` to choose a subscription or API-key provider", which keeps `byo-model`.

### CodebuffAI/freebuff

```yaml
  - owner: CodebuffAI
    repo: freebuff
    description: "Free, ad-supported coding agent available as a CLI, desktop app, web app, and cloud agent for GitHub repos"
    tags: [terminal, desktop, web, interactive, app-builder, community]
```

Evidence: <https://github.com/CodebuffAI/freebuff>. The README's product table lists "Freebuff Desktop: Run parallel agents locally (macOS, Windows, Linux)" and "Freebuff Web: Build and ship full-stack apps". It adds that "Web and Cloud provide sandboxes, previews, terminals, and deployment workflows".

The model catalog is curated, not bring-your-own, so there is still no model-access tag.

### stablyai/orca

```yaml
  - owner: stablyai
    repo: orca
    description: "Desktop agent development environment that runs Claude Code, Codex, or any CLI agent in parallel, each in its own git worktree"
    tags: [desktop, terminal, byo-model, orchestration, headless, community]
    notes: "ADE: runs any terminal coding agent side by side, each in its own git worktree, with remote servers and a mobile companion app."
```

Evidence: <https://github.com/stablyai/orca>. The README says "Run Claude Code, Codex, and any other coding agent in parallel, each in its own worktree", and under "Supported Agents", "Works with **any CLI agent**". The list includes Grok, Cursor, GitHub Copilot, OpenCode and Antigravity. It also covers "Pair a Remote Orca Server" and the mobile apps for iOS and Android.

### getpaseo/paseo

```yaml
  - owner: getpaseo
    repo: paseo
    description: "Self-hosted daemon for orchestrating Claude Code, Codex, OpenCode, and other agents, driven from desktop, mobile, web, or CLI"
    tags: [desktop, web, terminal, editor-plugin, self-hosted, byo-model, orchestration, headless, mcp, community]
    notes: "ADE: self-hosted daemon that runs Claude Code, Codex, Copilot, OpenCode, Pi, Antigravity and Muse Code on your own machines, driven from desktop, mobile, web or CLI."
```

Evidence: <https://github.com/getpaseo/paseo>. The README says "One interface for Claude Code, Codex, Copilot, OpenCode, Pi, Antigravity, and Muse Code agents".

### langchain-ai/open-swe

```yaml
  - owner: langchain-ai
    repo: open-swe
    description: "Robin, LangChain's deployable coding agent that turns issues into pull requests, reviews PRs, and monitors CI, built on Deep Agents"
    tags: [self-hosted, web, byo-model, autonomous, review, mcp, community]
    notes: Renamed Robin in its README; deploying it as-is requires LangSmith Managed Deep Agents
```

Evidence: <https://github.com/langchain-ai/open-swe>. The README's heading is "# Robin", and it says "investigate a codebase, implement changes, validate them, and deliver a pull request. It also reviews pull requests… monitors CI". It also says "Deploying it as is requires Managed Deep Agents", and "We're not accepting issues or external contributions at this time".

### pingdotgg/t3code

```yaml
  - owner: pingdotgg
    repo: t3code
    description: "Control surface for coding agents already on your machine, such as Claude Code, Codex, and OpenCode, from desktop, web, or mobile"
    tags: [desktop, web, self-hosted, byo-model, orchestration, mcp, acp, community]
    notes: "ADE: local server that drives your existing Claude Code, Codex, Cursor, OpenCode, Grok Build, Antigravity, Pi and any ACP agent, with remote access from mobile."
```

Evidence: <https://github.com/pingdotgg/t3code>. The README says "It works with Claude Code, Codex, Cursor, OpenCode, Grok Build, Antigravity, Kiro, Pi, Muse Code, fx, and any ACP agent". The ACP provider guide in the repo is `docs/user/providers-acp.md`. For `mcp`, the repo's `docs/user/outside-agents.md` says "Agents that T3 Code did not start can work with an environment through its MCP server".

### zed-industries/zed

```yaml
  - owner: zed-industries
    repo: zed
    description: "Multiplayer code editor with a built-in agent panel that also hosts external agents over ACP"
    tags: [ide, byo-model, local-models, interactive, acp, community]
```

Evidence: CONTRIBUTING excludes marketing adjectives. The README's own phrase is "a high-performance, multiplayer code editor" (<https://github.com/zed-industries/zed>), so "High-performance" is copied marketing. This is a wording fix only.

### Kuberwastaken/claurst

```yaml
  - owner: Kuberwastaken
    repo: claurst
    description: "Multi-provider terminal coding agent in Rust, begun as a clean-room reimplementation of Claude Code; beta"
    tags: [terminal, byo-model, interactive, headless, mcp, acp, community]
```

Evidence: <https://github.com/Kuberwastaken/claurst/blob/main/docs/mcp.md> describes "connecting Claurst to external tool, resource, and prompt servers". The README documents `claurst -p` headless runs and ACP over stdio.

## Trade-offs and risk

- **Dropping `editor-plugin` from grok-build and oh-my-pi.** The criterion I used is that an editor extension made for this tool exists, whoever publishes it. Codewhale and paseo keep the tag because each README links a VS Code extension built for that tool, even though a community member publishes it. grok-build and oh-my-pi reach editors only through generic ACP clients, which are not extensions for them. The cost is that both drop out of the editor-plugin filter, but `acp` still surfaces them to Zed and Neovim users. The alternative is to redefine `editor-plugin` so it includes "usable from an editor over ACP". That would also pull in claurst, goose and openinterpreter, which do not carry it today. The strict reading is consistent with how those three are tagged now.
- **Keeping `single-vendor` on grok-build.** Its user guide documents custom providers and self-hosted endpoints. I kept the tag for consistency with codex and copilot-cli, which are tagged the same way despite similar escape hatches. Under-tagging is safer than changing a facet convention in an audit.

## Limitations

- I did not open the docs sites of entries marked KEEP beyond their READMEs and keyword hits, so a capability added only in external docs (for example a new `headless` mode) may be missing.
- I did not test scope eligibility again for any entry except goose.
- I did not vet `MoonshotAI/kimi-code` or `SWE-agent/mini-swe-agent` (9,885 ★) as additions.

## Unresolved questions

1. Removing kimi-cli and gpt-pilot reverses the keep in `98e5b90`. Confirm removal, or keep them with a historical-significance note.
2. Aider passes six months idle on 2026-11-22. Should it be kept with the historical note (recommended) or removed?
3. Goose's README now calls it "a general-purpose agent for coding, research, writing, automation, and data analysis" and leads with the oi-kit library. Does it still pass the "about software development" test? I recommend keeping it, because coding is its first-listed use and it ships a complete CLI and desktop agent, but the user should decide.
4. Should `editor-plugin` cover ACP-only editor use? The answer decides the grok-build and oh-my-pi changes.

Status: DONE_WITH_CONCERNS
Summary: 42 entries audited: 21 KEEP, 16 UPDATE, 2 RENAME (openclaude; Codewhale, also retagged; both with migration keys), 3 REMOVE (Reasonix, kimi-cli, gpt-pilot). The proposed YAML passed `go run . -check` and `go test` on a scratch copy.
Concerns/Blockers: Two removals and the Aider keep-or-remove call are user decisions that reverse or extend earlier choices.
