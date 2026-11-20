# Review: agents refresh (uncommitted `data/agents.yml` + `history.go`)

Date: 2026-11-20. Scope: 15 changed entries, 2 renames with history migrations, 1 removal
(`esengine/DeepSeek-Reasonix`), 45 new entries. Checked independently against GitHub API metadata and
each repo's README/docs. No repo files were edited.

## Mechanical checks (all pass)

- `go vet ./...`, `go test ./...` and `go run . -check` (86 agents valid) pass. `go run . -build` (run
  in a scratch copy) succeeds. It omits 47 entries for now, the 45 new ones plus the 2 renamed slugs,
  because they have no metadata yet. That is expected: the next updater run fills it.
- All 86 slugs resolve to themselves on `gh api repos/OWNER/REPO`. None is a GitHub fork.
- Renames are correct. `repos/Gitlawb/openclaude` redirects to `Twigpine/openclaude`, and
  `repos/Hmbown/Codewhale` redirects to `codewhale-hq/Codewhale`. `data/history.jsonl` has 71 snapshots
  under each old key and none under the new keys, so the two `canonicalKeyMigrations` lines are needed
  and enough. No other history key maps to a new entry. The history keys missing from the list
  (Roo-Code, continue, trae-agent, suna, plandex, bolt.new, void) match none of the new entries.
- The removal is correct: `esengine/DeepSeek-Reasonix` is `archived: true`
  (https://github.com/esengine/DeepSeek-Reasonix).
- All 45 new entries clear the gates:
  - stars: 1,000 or more (lowest is `magiksoft/aider-desk` at 1,471);
  - archived: none;
  - last push: 2026-09-30 or later (oldest is `zgsm-ai/costrict`, 2026-09-30);
  - deprecation: a grep of every README for deprecat/archiv/sunset/"no longer maintained"/"moved to"
    found only false positives. The one real hit is stagewise's account service, covered under Low.
- Descriptions: all are 140 characters or fewer, none ends with a period, none contains `|`, and none
  names a model version. The only adjective hits are the three uses of "Minimal" (see Low).

## Findings

| # | Sev | Entry | Problem | Evidence | Proposed fix |
|---|-----|-------|---------|----------|--------------|
| 1 | Medium | `agegr/pi-web` | Tagged `orchestration`, but it runs pi's agent loop in-process: it depends on `@earendil-works/pi-coding-agent`, `pi-agent-core` and `pi-ai`. It does not drive a separate agent. It is the web counterpart of `vastsa/PI-Desktop` and `JetBrains/thinkrail`, which are correctly tagged `interactive` | https://github.com/agegr/pi-web/blob/main/package.json, https://github.com/agegr/pi-web | `tags: [web, byo-model, interactive, community]` and `notes: Web UI for the pi coding agent; independent of the pi team` |
| 2 | Medium | `zhukunpenglinyutong/jetbrains-cc-gui` | `orchestration` is an over-tag. It is an in-IDE chat panel over one CLI agent at a time, with no parallel sessions, worktrees or remote control. `olimorris/codecompanion.nvim` drives CLI agents the same way and is `interactive`, and so is Zed. `mcp` is also doubtful: "MCP server support" means the wrapped CLI's MCP servers, not MCP spoken by the plugin | https://github.com/zhukunpenglinyutong/jetbrains-cc-gui | `tags: [editor-plugin, interactive, community]` |
| 3 | Medium | `winfunc/opcode` | `mcp` is an over-tag. "MCP Server Management" edits Claude Code's MCP config; opcode does not speak MCP itself. `RunMaestro/Maestro` is a pass-through of the same kind and is correctly left without `mcp` | https://github.com/winfunc/opcode#-mcp-server-management | `tags: [desktop, orchestration, community]` |
| 4 | Medium | `codewhale-hq/Codewhale` (changed) | `editor-plugin` was kept, but the README calls the VS Code extension "community-maintained", published by `HengQuWorld` (10 stars). Editor use is already covered by `acp` (`codewhale serve --acp`). The same diff drops `editor-plugin` from `oh-my-pi` on the same reasoning | https://github.com/codewhale-hq/Codewhale (Usage table), https://github.com/HengQuWorld/CodeWhale-VSCode | `tags: [terminal, web, byo-model, local-models, interactive, headless, mcp, acp, community]` |
| 5 | Medium | `MoonshotAI/kimi-cli` (unchanged, but this diff makes it inconsistent) | Archived upstream, with no historical-significance claim in its notes. The diff removes Reasonix for being archived and adds the successor `MoonshotAI/kimi-code`, so keeping kimi-cli breaks the archived-removal rule in CONTRIBUTING | https://github.com/MoonshotAI/kimi-cli (archived, "no longer maintained") | remove (or, if kept on purpose, change notes to `Archived upstream; kept for historical significance; succeeded by MoonshotAI/kimi-code`) |
| 6 | Low | `Tencent/LoopForge` | `mcp` is weak. The README says roles get "Skills and MCP tools" that are handed to the executor CLIs, so LoopForge configures MCP rather than speaking it. Under-tag rule | https://github.com/Tencent/LoopForge | `tags: [self-hosted, orchestration, vendor]` (see Unresolved Q1 on `vendor`) |
| 7 | Low | `onlook-dev/onlook` (changed) | `app-builder` means prompt-to-app with preview and deploy. The README now describes a design canvas over your existing repo, with export (HTML/PNG/PDF) but no deploy or scaffolding. `byo-model` holds only loosely: it uses your Claude Code or Codex account | https://github.com/onlook-dev/onlook | optional: `tags: [desktop, byo-model, mcp, community]` |
| 8 | Low | `humanlayer/humanlayer` | The description does not say what the tool is. It starts with a verb, against "say what the tool is" | https://github.com/humanlayer/humanlayer | `description: "Platform that runs Claude Code, Codex, OpenCode, and Pi on your machines so a team can steer them from desktop, web, or CLI"` (123 chars) |
| 9 | Low | `earendil-works/pi`, `SWE-agent/mini-swe-agent`, `vercel-labs/fx` | "Minimal" opens all three descriptions. It is accurate (each README self-describes so), not marketing, but it is the only adjective-rule hit | READMEs | optional; acceptable as is |
| 10 | Low | 22 entries with `notes: "ADE."` / `"ADE; ..."` | A bare "ADE." duplicates the `orchestration` tag and tells the reader nothing | n/a | optional: drop bare `notes: "ADE."` lines |
| 11 | Low | `google-gemini/gemini-cli` (changed) | "now aimed at" is time-relative wording in a description. The facts are correct: free and Google AI plans moved to Antigravity CLI on 2026-06-18 | https://github.com/google-gemini/gemini-cli | optional: `description: "Google's open-source terminal agent for Gemini models, maintained for enterprise Gemini Code Assist users, with MCP support"` |

## Scope borderlines (judgement calls; I would keep all, but flag them)

- `gptme/gptme`: the README calls it "A personal AI agent ... general-purpose enough to assist in all
  kinds of knowledge-work". It also says "A great coding agent" and "a development-focused peer" to
  general agents, and positions itself as an alternative to Claude Code and Codex. Kept on the
  development focus. https://github.com/gptme/gptme
- `wavetermdev/waveterm`: an "AI-integrated terminal", where Wave AI is one widget and command
  execution is "Coming Soon". It qualifies only as an AI-assisted terminal (the Warp precedent), and
  the LLM is less central than Warp's. https://github.com/wavetermdev/waveterm
- `Gaurav-Gosain/tuios`: the README headline is a general terminal multiplexer and window manager.
  Agent state, the Inbox, `tuios fan`, and ACP/MCP make up a large Agents section, and the GitHub
  description leads with agents. Kept as an ADE, like herdr and cmux. https://github.com/Gaurav-Gosain/tuios
- `huggingface/tau`: also "a teaching project", but a shippable CLI on PyPI. It does not need
  `research`.

## Edge cases found while scouting

- Staleness: `Aider-AI/aider` (unchanged) was last pushed 2026-05-22, so it
  crosses the 6-month removal line on about 2026-11-22. `TabbyML/tabby` (2026-06-30) and
  `Pythagora-io/gpt-pilot` (2026-06-18) are past the 3-month warning. None of this is in the diff;
  it is a heads-up for the next run.
- `stagewise-io/stagewise`: the stagewise account service is sunset at the end of November 2026. The
  tool keeps working with your own keys or subscription, so it stays in scope and its `byo-model`
  tag stays correct. https://github.com/stagewise-io/stagewise#stagewise-account-sunset
- `coder/xum` (formerly `coder/mux`), `ParadevAI/Paradev` (formerly `superset-sh/superset`) and
  `nimbalyst/nimbalyst` (successor to `stravu/crystal`): none of the old slugs is in the list or in
  history, so they need no migration. Old-slug redirects are confirmed for xum and Paradev.
- No duplicates under other slugs. Pi derivatives (`oh-my-pi`, `tau`, `PI-Desktop`, `pi-web`,
  `thinkrail`), `happy` vs `hapi`, `cmux` vs `xum`, and `costrict` (credits Roo Code, which is not
  listed) are all distinct products.

## Tag spot-checks verified OK

New entries (32 checked against README/docs):

- `herdrdev/herdr`, `manaflow-ai/cmux` (GPL client plus BUSL server note confirmed),
  `slopus/happy` (ACP confirmed in `packages/happy-cli/src/agent/acp/`), `wavetermdev/waveterm`
- `ParadevAI/Paradev`, `OrchestratorInc/agent-orchestrator` (iPhone and Android apps confirmed),
  `siteboon/claudecodeui` (desktop app referenced in README)
- `openchamber/openchamber`, `SWE-agent/mini-swe-agent`, `smtg-ai/claude-squad`, `zai-org/ZCode`
- `JetBrains/thinkrail` (incubator badge confirmed; `community` is correct), `MoonshotAI/kimi-code`,
  `intitni/CopilotForXcode`, `tailcallhq/forgecode`, `olimorris/codecompanion.nvim`
- `stagewise-io/stagewise`, `PawanOsman/OpenCursor`, `get-bb/bb` (Claude Code, Codex and Pi in
  `docs/provider-plugin-api.md`; OpenCode through the ACP provider plugin)
- `mistralai/mistral-vibe` (ACP and programmatic mode in docs), `tiann/hapi`,
  `agent-of-empires/agent-of-empires`, `zgsm-ai/costrict`, `vercel-labs/fx`, `twinnydotdev/twinny`
- `CopilotC-Nvim/CopilotChat.nvim`, `RunMaestro/Maestro`, `huggingface/tau` (`-p` print mode, custom
  OpenAI-compatible and local endpoints), `Nano-Collective/nanocoder`
- `MiniMax-AI/minimax-code` (`mcode exec`, `mcode acp`, OpenAI- and Anthropic-compatible providers),
  `nimbalyst/nimbalyst`, `coder/xum`, `kodustech/kodus-ai`, `ggml-org/llama.vscode`
- `magiksoft/aider-desk` (Ollama listed on aiderdesk.dev docs), `microsoft/intelligent-terminal`,
  `gastownhall/gastown`, `vastsa/PI-Desktop`

Changed entries (all 15 checked):

- `google-gemini/gemini-cli` (notes accurate), `zed-industries/zed`, `earendil-works/pi` (print/JSON,
  RPC and SDK confirm `headless`), `OpenHands/OpenHands` (Agent Canvas runs ACP agents, so
  `orchestration` fits)
- `openinterpreter/openinterpreter` (`interpreter exec`), `can1357/oh-my-pi`, `xai-org/grok-build`
  ("SpaceXAI" matches README)
- `PrimeIntellect-ai/prime-agent` (`mcp` confirmed in `crates/pa-cli/src/mcp_command.rs` and docs),
  `CodebuffAI/freebuff` (desktop, web and app-builder per README product table)
- `Kuberwastaken/claurst` (MCP servers in README), `stablyai/orca` (SSH hosts, remote server, mobile),
  `getpaseo/paseo` (Antigravity and Muse Code confirmed), `langchain-ai/open-swe` (Robin and
  Managed Deep Agents note accurate)
- `pingdotgg/t3code` (ACP agents and "Connect outside agents over MCP" confirmed)
- `Codewhale` is fine apart from finding 4.

## Unresolved questions

1. **`vendor` boundary.** `microsoft/intelligent-terminal` and `Tencent/LoopForge` are `vendor`.
   `JetBrains/thinkrail`, `huggingface/tau` and `vercel-labs/fx` are `community`. Neither Microsoft
   nor Tencent ties its tool to its own models: Intelligent Terminal defaults to Copilot and runs
   any ACP agent, and LoopForge runs ten CLIs. The tags follow existing precedent
   (`github/copilot-cli`, and `alibaba/open-code-review` with `byo-model`), so the working rule seems
   to be "the publisher runs a model lab". Under that rule Microsoft and Tencent are fine, and
   JetBrains, Hugging Face and Vercel are defensibly `community`. Under a strict reading ("first-party
   tool *for the lab's models*"), Microsoft and Tencent should be `community`. Options:
   - (a) keep everything as is;
   - (b) `microsoft/intelligent-terminal` → `tags: [terminal, desktop, interactive, acp, community]`
     and `Tencent/LoopForge` → `tags: [self-hosted, orchestration, community]`.

   Whichever you pick, state the rule in CONTRIBUTING.
2. **Single-agent frontends as `orchestration`.** `opcode` (Claude Code only), `claudecodeui`,
   `openchamber` (OpenCode only) and `t3code` are tagged ADEs because they drive a separate agent's
   sessions. Findings 1 and 2 apply a test: does the tool drive an external agent's sessions, with
   session management, parallelism or remote control (→ `orchestration`), or embed the agent loop
   or act as a chat panel (→ `interactive`)? Confirm that this is the intended rule.
