# Research Report: awesome-ai-dev-tools list refresh (additions + stale cleanup)

**Conducted:** 2026-09-11 13:28 (+07) · **Repo:** `tiennm99/awesome-ai-dev-tools` · **Tracked now:** 29

## Executive Summary

List has a **hard gap at the top**: `deepseek-ai/deepseek-harness` (219.6k ★, DeepSeek's official agent harness, shipped Aug 2026) is untracked and would rank **#1**, above `anomalyco/opencode` (206.6k). `xai-org/grok-build` (xAI's official coding agent, 26.7k) is also missing. Both are first-party vendor agents — same tier as claude-code/codex/gemini-cli, which are tracked.

14 repos meet the README's own inclusion criteria (coding agent itself, ≥10k ★, open source, actively maintained) and are absent. 7 tracked repos have no push in >3 months; 3 of them are dead beyond debate (bolt.new 21 mo, plandex 11 mo, trae-agent 7 mo).

Also found one data defect: `data/agents.yml` still lists `yetone/avante.nvim` while GitHub canonical is `avante-corp/avante.nvim`, violating the file's own header rule.

## Methodology

- Sources: live GitHub GraphQL (all 29 tracked + 380 candidate slugs probed), 2 web searches, 2 curated landscape lists (`joylarkin/AI-Coding-Landscape`, `bradAGI/awesome-cli-coding-agents`), 17 candidate READMEs read for scope classification.
- Star/push data: live as of 2026-09-11. Staleness cutoff: `pushedAt < 2026-06-11`.
- Scope filter = README criteria verbatim: agent that **writes, edits, or reviews** code; frameworks/wrappers/observability/workspaces *for* agents excluded.

---

## 1. Additions (meet all criteria, verified)

| ★ | owner/repo | category | pushed | what it is |
|--:|---|---|---|---|
| 219.6k | `deepseek-ai/deepseek-harness` | cli | 09-11 | **DeepSeek's official agent harness (`dsh`)**, plugin architecture. Would be **#1** |
| 68.3k | `openinterpreter/openinterpreter` | cli | 09-09 | Open Interpreter, repositioned as a coding agent for open models |
| 40.9k | `Hmbown/Codewhale` | cli | 09-11 | Rust terminal coding agent (ex-deepseek-tui), 30+ providers |
| 35.5k | `esengine/DeepSeek-Reasonix` | cli | 09-11 | Go binary; CLI/TUI + desktop + VS Code over one local engine |
| 33.1k | `Gitlawb/openclaude` | cli | 09-09 | TS coding-agent CLI, multi-provider, MCP |
| 30.6k | `can1357/oh-my-pi` | cli | 09-11 | Coding agent with LSP/IDE wired in (TS+Rust) |
| 26.7k | `xai-org/grok-build` | cli | 09-09 | **xAI's official coding agent** (TUI, headless, ACP) |
| 22.2k | `alibaba/open-code-review` | cli | 09-11 | Alibaba's AI code-review agent CLI (reviews code ⇒ in scope) |
| 20.5k | `PrimeIntellect-ai/prime-agent` | cli | 09-11 | Self-improving coding/long-horizon agent harness |
| 19.5k | `1jehuang/jcode` | cli | 09-11 | Coding agent TUI, self-updating, YC-launched |
| 16.5k | `HKUDS/DeepCode` | research | 09-09 | Multi-agent code generation; TUI/desktop/web clients |
| 11.9k | `CodebuffAI/freebuff` | cli | 09-11 | Codebuff's free terminal coding agent |
| 11.3k | `MoonshotAI/kimi-cli` | cli | 09-01 | Moonshot's official Kimi Code CLI agent |
| 10.3k | `Kuberwastaken/claurst` | cli | 09-02 | Rust multi-provider terminal agent (at the 10k floor) |

YAML patch (append to `data/agents.yml`):

```yaml
  - owner: deepseek-ai
    repo: deepseek-harness
    category: cli
    notes: DeepSeek official harness (dsh); developer preview, breaking changes expected
  - owner: openinterpreter
    repo: openinterpreter
    category: cli
  - owner: Hmbown
    repo: Codewhale
    category: cli
  - owner: esengine
    repo: DeepSeek-Reasonix
    category: cli
  - owner: Gitlawb
    repo: openclaude
    category: cli
  - owner: can1357
    repo: oh-my-pi
    category: cli
  - owner: xai-org
    repo: grok-build
    category: cli
  - owner: alibaba
    repo: open-code-review
    category: cli
  - owner: PrimeIntellect-ai
    repo: prime-agent
    category: cli
  - owner: 1jehuang
    repo: jcode
    category: cli
  - owner: HKUDS
    repo: DeepCode
    category: research
  - owner: CodebuffAI
    repo: freebuff
    category: cli
  - owner: MoonshotAI
    repo: kimi-cli
    category: cli
  - owner: Kuberwastaken
    repo: claurst
    category: cli
```

## 2. Stale cleanup (no push in >3 months)

| tracked repo | ★ | last push | age | archived | verdict |
|---|--:|---|--:|---|---|
| `stackblitz/bolt.new` | 16.5k | 2024-12-17 | 21 mo | no | **remove** — abandoned, `bolt.diy` fork also stale (2026-02-07) |
| `plandex-ai/plandex` | 15.6k | 2025-10-03 | 11 mo | no | **remove** |
| `bytedance/trae-agent` | 12.1k | 2026-02-05 | 7 mo | no | **remove** |
| `AntonOsika/gpt-engineer` | 55.1k | 2025-05-14 | 16 mo | yes | keep only if the "historical significance" carve-out stands; else remove |
| `RooCodeInc/Roo-Code` | 24.3k | 2026-05-15 | 3.9 mo | yes | **user call** — archived; historical value as Cline fork lineage |
| `voideditor/void` | 28.8k | 2026-06-02 | 3.3 mo | yes | **user call** — archived, empty description |
| `Aider-AI/aider` | 48.9k | 2026-05-22 | 3.7 mo | no | **user call** — hardest one; 2026 roundups still call it a top OSS agent, but `pushedAt` shows no commits in 3.7 mo |

Watch list (inside 3 mo, trending stale): `TabbyML/tabby` (2026-06-30), `Pythagora-io/gpt-pilot` (2026-06-18).

## 3. Rejected candidates (why)

Out per criterion 1 — orchestrators, runtimes, wrappers, skill packs, general assistants:

`NousResearch/hermes-agent` 244.3k (general personal agent, not coding-specific) · `ultraworkers/claw-code` 195k (Rust `claw` harness, but README self-describes as an "agent-managed museum exhibit"; star count disproportionate to substance — **flagged, user call**) · `bytedance/deer-flow` 82.3k (general SuperAgent) · `ruvnet/ruflo` 72k (meta-harness/swarm) · `code-yeongyu/oh-my-openagent` 68.9k (plugin bundle for existing hosts) · `stablyai/orca` 66.2k, `superset-sh/superset` 14.1k, `gastownhall/gastown` 18k, `BloopAI/vibe-kanban` 28.1k, `different-ai/openwork` 23.5k, `getpaseo/paseo`, `iOfficeAI/AionUi` (fleet/orchestration/workspace) · `herdrdev/herdr` 37.6k (runtime agents run on) · `sipeed/picoclaw` 30k (personal assistant on $10 hardware) · `musistudio/claude-code-router` 37.2k, `farion1231/cc-switch` 132k, `router-for-me/CLIProxyAPI` 51.3k (routers) · `trycua/cua` 22.5k (computer-use) · `FoundationAgents/MetaGPT` (framework) · skills/prompt repos (`obra/superpowers`, `anthropics/skills`, `addyosmani/agent-skills`, …).

Below the 10k floor: `aws/amazon-q-developer-cli` 2.0k · `smallcloudai/refact` 3.5k (archived) · `superagent-ai/grok-cli` 3.5k · `just-every/code` 4.0k · `codestoryai/aide` 2.2k (archived).

Dead + below relevance: `smol-ai/developer` 12.2k (2024-04-07) · `Doriandarko/claude-engineer` 11.2k (2024-12-12) · `stitionai/devika` 19.6k (2025-09-25) — all >10k ★ but long abandoned; adding them would contradict the cleanup in §2.

No repo to track (closed source): Cursor, Antigravity (the closed Gemini CLI successor already noted in `agents.yml`), Kiro, Devin, Amp, Factory Droid, Augment.

## 4. Data defect

`data/agents.yml` lists `owner: yetone` / `repo: avante.nvim`; canonical slug is now **`avante-corp/avante.nvim`** (API redirects, README already renders the new name). The file header mandates the current canonical slug + a `history.go` migration. Fix:

1. `data/agents.yml` → `owner: avante-corp`
2. `history.go` `canonicalKeyMigrations` → add `"yetone/avante.nvim": "avante-corp/avante.nvim"`

Without step 2, avante's star history detaches and Δ7d shows `—`.

## Next steps

1. Apply §1 YAML patch (14 adds) + §4 fix; run `go run . -check` then the updater.
2. Decide the 4 user-call rows in §2 (aider, Roo-Code, void, gpt-engineer) and `claw-code` in §3.
3. Consider adding to CONTRIBUTING: an explicit **staleness rule** (e.g. no push in 6 months ⇒ drop unless historically significant), since "actively maintained" is currently unquantified — that is what let bolt.new sit for 21 months.

## Unresolved questions

1. `Aider-AI/aider` — keep despite 3.7 mo of no pushes (name recognition) or drop?
2. Archived-but-significant (`void`, `Roo-Code`, `gpt-engineer`) — keep marked, or is the carve-out only for gpt-engineer?
3. `ultraworkers/claw-code` (195k ★, functional Rust harness, joke-flavored README) — include, exclude, or exclude with a documented reason?
4. Staleness threshold to codify: 3 months (your ask) or 6 months (less churn)?
5. `alibaba/open-code-review` is review-only, not an editing agent — in scope per criteria wording; confirm intent?

## Sources

- [Best Open Source CLI Coding Agents in 2026 — Pinggy](https://pinggy.io/blog/best_open_source_cli_coding_agents/)
- [Best Terminal AI Coding Agents in 2026 — amux](https://amux.io/blog/best-terminal-ai-coding-agents-2026/)
- [Best CLI AI Tools in 2026 — kilo.ai](https://kilo.ai/articles/best-cli-coding-agents)
- [AI-Coding-Landscape](https://github.com/joylarkin/AI-Coding-Landscape) · [awesome-cli-coding-agents](https://github.com/bradAGI/awesome-cli-coding-agents)
- Live GitHub GraphQL API (stars, pushedAt, isArchived) — 2026-09-11
