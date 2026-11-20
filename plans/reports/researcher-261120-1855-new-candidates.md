# New candidate entries, 2026-11-20

Method: about 25 GitHub topic searches (coding-agent, ai-coding, claude-code, agentic-coding, ai-ide, code-review, ai-code-review, acp, agent-client-protocol, mcp, codex, opencode, git-worktree, agent-orchestration and others), re-run in star bands (1k-2k, 2k-4k, 4k-8k) for the large topics; about 20 keyword searches on names and descriptions; a `created:>=2026-08-01` sweep; a scan of new repos in 27 model-lab and big-tech orgs; three web searches for Sep-Nov 2026 launches. Every recommended repo was checked with `gh api repos/OWNER/REPO` (stars, `pushed_at`, `archived`, license) plus the date of the last commit on its default branch, and its README was read for scope, deprecation notices and tag evidence. Data is live as of 2026-11-20.

## Outcome

The list has a large backlog. The 09-29 round used a 10k-star bar and rejected every orchestrator. The current rules have a 1,000-star floor and admit ADEs. Under those rules, I verified 45 repos that qualify cleanly. This is not an exhaustive census: a tail of 1k-5k repos was seen but not reviewed (see Limitations). By shape, the 45 are 20 ADEs, the shape excluded before, 15 terminal or IDE agents (four of them from model labs: Moonshot, Mistral, Z.ai and MiniMax), 7 editor plugins, 2 AI terminals and 1 review tool.

Two "obvious" omissions do not qualify. `continuedev/continue` (36.5k) says in its README that it is "no longer actively maintained and is read-only"; it was already removed in commit 00c90c5. `BloopAI/vibe-kanban` (28.8k) heads its README with "Vibe Kanban is sunsetting".

Strongest picks, if you want a first batch: herdr, cmux, happy, waveterm, Paradev, agent-orchestrator, mini-swe-agent, kimi-code, mistral-vibe, ThinkRail, codecompanion.nvim, gptme.

## Recommended additions

Sorted by stars. "Last commit" is the newest commit on the default branch, which can be older than `pushed_at`.

| Stars | owner/repo | pushed_at | Last commit | Why in scope |
|--:|---|---|---|---|
| 51,992 | herdrdev/herdr | 2026-11-19 | 2026-11-19 | ADE: terminal workspace that runs and monitors existing coding agents |
| 28,841 | manaflow-ai/cmux | 2026-11-20 | 2026-11-20 | ADE: macOS terminal built for running many coding agents |
| 24,722 | slopus/happy | 2026-11-20 | 2026-11-20 | ADE: mobile/web control of Claude Code, Codex and ACP agents |
| 22,723 | wavetermdev/waveterm | 2026-11-19 | 2026-10-05 | AI-assisted terminal (same shape as tracked warp) |
| 22,452 | winfunc/opcode | 2026-11-06 | 2026-11-06 | ADE: desktop session manager for Claude Code |
| 18,849 | gastownhall/gastown | 2026-11-14 | 2026-11-12 | ADE: coordinates many agents with persistent work tracking (rejected 09-11 as orchestration) |
| 16,418 | ParadevAI/Paradev | 2026-11-20 | 2026-11-20 | ADE: parallel CLI agents in worktrees (ex-Superset) |
| 14,960 | OrchestratorInc/agent-orchestrator | 2026-11-20 | 2026-11-20 | ADE: fleet of coding agents from task to merged PR |
| 14,701 | siteboon/claudecodeui | 2026-11-19 | 2026-11-19 | ADE: web/mobile UI for Claude Code, Cursor CLI, Codex |
| 13,092 | openchamber/openchamber | 2026-11-20 | 2026-11-19 | ADE on OpenCode: desktop, web, VS Code, mobile |
| 12,919 | humanlayer/humanlayer | 2026-11-16 | 2026-11-16 | ADE: team control of agents running on your machines |
| 9,885 | SWE-agent/mini-swe-agent | 2026-11-16 | 2026-09-03 | Coding agent; the tracked SWE-agent's notes already point to it |
| 8,839 | smtg-ai/claude-squad | 2026-11-13 | 2026-11-13 | ADE: terminal manager for parallel agent sessions |
| 8,806 | zai-org/ZCode | 2026-11-16 | 2026-11-13 | Vendor coding workspace from Z.ai (new, created 2026-09-20) |
| 8,371 | JetBrains/thinkrail | 2026-11-20 | 2026-11-18 | Agentic IDE (JetBrains incubator) |
| 8,334 | agegr/pi-web | 2026-11-20 | 2026-11-20 | ADE: browser UI for tracked pi sessions |
| 8,270 | MoonshotAI/kimi-code | 2026-11-20 | 2026-11-20 | Vendor terminal agent; successor of tracked, archived kimi-cli |
| 8,202 | intitni/CopilotForXcode | 2026-10-24 | 2026-10-24 | AI-assisted editor plugin for Xcode |
| 8,105 | vastsa/PI-Desktop | 2026-11-20 | 2026-11-20 | Desktop coding agent on the pi harness |
| 7,745 | zhukunpenglinyutong/jetbrains-cc-gui | 2026-11-20 | 2026-11-16 | ADE in the IDE: JetBrains GUI for coding CLIs |
| 7,636 | tailcallhq/forgecode | 2026-11-18 | 2026-10-31 | Terminal coding agent |
| 6,944 | olimorris/codecompanion.nvim | 2026-11-19 | 2026-11-19 | Neovim AI plugin and ACP client |
| 6,833 | stagewise-io/stagewise | 2026-11-16 | 2026-11-04 | Agentic IDE with built-in agent |
| 6,589 | Gaurav-Gosain/tuios | 2026-11-20 | 2026-11-20 | ADE: terminal multiplexer that tracks and fans out to agents |
| 6,047 | PawanOsman/OpenCursor | 2026-10-26 | 2026-10-26 | VS Code coding agent |
| 5,832 | get-bb/bb | 2026-11-20 | 2026-11-20 | ADE: agentic IDE driving existing agent CLIs |
| 5,438 | mistralai/mistral-vibe | 2026-11-18 | 2026-11-17 | Vendor terminal agent (rejected 09-29 only for stars) |
| 5,436 | tiann/hapi | 2026-11-20 | 2026-11-20 | ADE: remote control for local agent sessions |
| 5,208 | Tencent/LoopForge | 2026-11-20 | 2026-11-20 | ADE: self-hosted team workspace running agent CLIs |
| 4,590 | agent-of-empires/agent-of-empires | 2026-11-20 | 2026-11-20 | ADE: session manager with TUI and web dashboard |
| 4,489 | gptme/gptme | 2026-11-20 | 2026-11-20 | Terminal agent (rejected 09-29 only for stars) |
| 4,464 | zgsm-ai/costrict | 2026-09-30 | 2026-09-30 | VS Code/JetBrains coding agent |
| 3,933 | vercel-labs/fx | 2026-11-19 | 2026-11-18 | Terminal coding agent (experimental) |
| 3,703 | twinnydotdev/twinny | 2026-11-18 | 2026-11-18 | VS Code AI assistant for local models |
| 3,645 | CopilotC-Nvim/CopilotChat.nvim | 2026-10-23 | 2026-10-23 | Neovim AI chat plugin with tool calling |
| 3,615 | RunMaestro/Maestro | 2026-11-20 | 2026-11-20 | ADE: desktop agent fleet manager |
| 3,066 | huggingface/tau | 2026-11-16 | 2026-11-13 | Terminal coding agent |
| 2,582 | Nano-Collective/nanocoder | 2026-11-20 | 2026-11-20 | Terminal coding agent |
| 2,269 | MiniMax-AI/minimax-code | 2026-11-17 | 2026-11-04 | Vendor terminal agent |
| 2,250 | microsoft/intelligent-terminal | 2026-11-20 | 2026-11-20 | AI-assisted terminal hosting ACP agents |
| 2,092 | nimbalyst/nimbalyst | 2026-11-18 | 2026-11-18 | ADE: visual workspace for parallel agents |
| 2,081 | coder/xum | 2026-11-20 | 2026-11-20 | Desktop coding agent multiplexer |
| 1,634 | kodustech/kodus-ai | 2026-11-19 | 2026-11-19 | PR review agent (same shape as tracked pr-agent) |
| 1,559 | ggml-org/llama.vscode | 2026-11-20 | 2026-11-20 | VS Code plugin for local completion, chat and agent |
| 1,471 | magiksoft/aider-desk | 2026-11-20 | 2026-11-20 | Desktop AI coding app |

All are unarchived and carry an OSI license (MIT, Apache-2.0, AGPL-3.0, GPL-3.0, LGPL-3.0). Four show `NOASSERTION` or "Other" in the API, and I read their LICENSE files: cmux is a GPL-3.0-or-later client with BUSL-1.1 server directories, nanocoder is plain MIT (the API is confused by a heading), kodus-ai is AGPL-3.0 plus an enterprise license for `ee/` files, and LoopForge is MIT with third-party notices, and none declares itself deprecated.

## YAML

Descriptions were checked mechanically: all are at most 140 characters, contain no `|`, have no trailing period and name no model versions. Origin calls: JetBrains, Vercel, Hugging Face and ggml-org are tagged `community` because they are not model labs. Tencent (Hunyuan) is tagged `vendor`. Microsoft is tagged `vendor` by analogy with the tracked `github/copilot-cli`; see the questions at the end.

```yaml
  - owner: herdrdev
    repo: herdr
    description: "Terminal workspace that keeps Claude Code, Codex, and other agents running in persistent panes and flags the ones that are blocked"
    tags: [terminal, orchestration, community]
    notes: "ADE: owns the terminals of the agents you already run, locally or on SSH machines, without wrapping them."
  - owner: manaflow-ai
    repo: cmux
    description: "Ghostty-based macOS terminal with vertical tabs, agent notifications, and a scriptable in-app browser for running many coding agents"
    tags: [desktop, terminal, orchestration, community]
    notes: "ADE. Client is GPL-3.0-or-later; the server directories are BUSL-1.1."
  - owner: slopus
    repo: happy
    description: "Mobile, web, and CLI client that runs Claude Code, Codex, and ACP agents on your own machine with end-to-end encrypted sync"
    tags: [web, terminal, byo-model, orchestration, acp, community]
    notes: "ADE. Desktop app lives in slopus/happy-desktop."
  - owner: wavetermdev
    repo: waveterm
    description: "Cross-platform terminal with a built-in AI assistant that reads terminal output, plus durable SSH sessions and a graphical editor"
    tags: [terminal, desktop, byo-model, local-models, interactive, community]
  - owner: winfunc
    repo: opcode
    description: "Desktop GUI for Claude Code that manages interactive sessions, custom agents, background runs, and MCP servers"
    tags: [desktop, orchestration, mcp, community]
    notes: "ADE for Claude Code only; not affiliated with Anthropic."
  - owner: ParadevAI
    repo: Paradev
    description: "Desktop app that runs Claude Code, Codex, and other CLI agents in parallel, each in its own git worktree with review and previews"
    tags: [desktop, orchestration, community]
    notes: "ADE. Formerly superset-sh/superset."
  - owner: OrchestratorInc
    repo: agent-orchestrator
    description: "Desktop workspace where an orchestrator splits work across coding agents in separate worktrees and follows it through CI and review"
    tags: [desktop, orchestration, community]
    notes: "ADE with iPhone and Android companion apps."
  - owner: siteboon
    repo: claudecodeui
    description: "Self-hostable web, desktop, and mobile UI for Claude Code, Cursor CLI, and Codex sessions, with files, git, and a shell"
    tags: [web, desktop, self-hosted, orchestration, community]
    notes: "ADE, published as CloudCLI UI."
  - owner: openchamber
    repo: openchamber
    description: "Workspace for running and reviewing OpenCode agent sessions from desktop, web, VS Code, and mobile"
    tags: [desktop, web, editor-plugin, orchestration, community]
    notes: "ADE built on OpenCode; independent of the OpenCode team."
  - owner: SWE-agent
    repo: mini-swe-agent
    description: "Minimal bash-only software engineering agent from the SWE-bench team that runs with any model, locally or in sandboxes"
    tags: [terminal, byo-model, autonomous, research, community]
  - owner: smtg-ai
    repo: claude-squad
    description: "Terminal app that manages several Claude Code, Codex, Gemini, or Aider sessions, each in its own git worktree"
    tags: [terminal, orchestration, community]
    notes: "ADE."
  - owner: zai-org
    repo: ZCode
    description: "Z.ai's AI coding workspace with desktop, browser, and terminal interfaces on one agent runtime"
    tags: [desktop, web, terminal, vendor]
  - owner: JetBrains
    repo: thinkrail
    description: "Agentic IDE built around the pi agent, with git worktree workspaces and a spec graph that keeps project knowledge"
    tags: [ide, desktop, web, byo-model, local-models, interactive, mcp, community]
    notes: JetBrains incubator project
  - owner: MoonshotAI
    repo: kimi-code
    description: "Moonshot's terminal coding agent for Kimi models, configurable for other providers and drivable from editors over ACP"
    tags: [terminal, byo-model, interactive, mcp, acp, vendor]
    notes: Successor to the archived MoonshotAI/kimi-cli
  - owner: intitni
    repo: CopilotForXcode
    description: "Xcode source editor extension that adds GitHub Copilot, Codeium, and ChatGPT suggestions, chat, and code modification"
    tags: [editor-plugin, byo-model, interactive, community]
  - owner: tailcallhq
    repo: forgecode
    description: "Terminal coding agent with an interactive TUI, a one-shot prompt mode, and a zsh plugin, working with many LLM providers"
    tags: [terminal, byo-model, interactive, headless, mcp, community]
  - owner: olimorris
    repo: codecompanion.nvim
    description: "Neovim plugin for chatting with LLMs and driving CLI agents over ACP, with MCP support and many providers"
    tags: [editor-plugin, byo-model, local-models, interactive, acp, mcp, community]
  - owner: stagewise-io
    repo: stagewise
    description: "Agentic IDE with a built-in browser whose coding agent can read the page's console and debugger"
    tags: [ide, byo-model, local-models, interactive, community]
  - owner: PawanOsman
    repo: OpenCursor
    description: "VS Code coding agent that works with your API keys, existing subscriptions, or local Ollama and llama.cpp models"
    tags: [editor-plugin, byo-model, local-models, interactive, mcp, community]
  - owner: get-bb
    repo: bb
    description: "Agentic IDE that drives Claude Code, Codex, Pi, and OpenCode from a desktop app, web app, CLI, or HTTP API"
    tags: [desktop, web, terminal, orchestration, community]
    notes: "ADE; in active development."
  - owner: mistralai
    repo: mistral-vibe
    description: "Mistral's terminal coding agent that asks before risky actions, with a programmatic mode and ACP editor support"
    tags: [terminal, interactive, headless, acp, vendor]
  - owner: tiann
    repo: hapi
    description: "Local-first remote control for Claude Code, Codex, OpenCode, and other agents from mobile apps, web, or Telegram"
    tags: [web, terminal, self-hosted, orchestration, community]
    notes: "ADE."
  - owner: Tencent
    repo: LoopForge
    description: "Self-hosted team workspace that runs Claude Code, Codex, CodeBuddy, and other agent CLIs in Docker sandboxes via chat or SOPs"
    tags: [self-hosted, orchestration, mcp, vendor]
    notes: "ADE."
  - owner: agent-of-empires
    repo: agent-of-empires
    description: "Session manager for coding agents on Linux and macOS, with a TUI, a web dashboard, worktrees, and optional containers"
    tags: [terminal, web, orchestration, community]
    notes: "ADE."
  - owner: gptme
    repo: gptme
    description: "Terminal agent that writes code, runs shell and Python, and browses the web, with any provider or local models"
    tags: [terminal, web, desktop, byo-model, local-models, interactive, headless, mcp, acp, community]
  - owner: zgsm-ai
    repo: costrict
    description: "VS Code and JetBrains coding agent for enterprises, with a strict spec-first mode and private deployment"
    tags: [editor-plugin, byo-model, local-models, interactive, mcp, community]
  - owner: vercel-labs
    repo: fx
    description: "Minimal native coding agent CLI, scriptable with fx ask and drivable from editors over ACP"
    tags: [terminal, byo-model, local-models, interactive, headless, acp, mcp, community]
    notes: Experimental upstream
  - owner: twinnydotdev
    repo: twinny
    description: "VS Code assistant for completion, chat, inline edits, and code review against local or self-hosted model servers"
    tags: [editor-plugin, byo-model, local-models, interactive, community]
  - owner: CopilotC-Nvim
    repo: CopilotChat.nvim
    description: "Neovim plugin for GitHub Copilot Chat with tool calling and custom providers such as Ollama"
    tags: [editor-plugin, byo-model, local-models, interactive, community]
  - owner: RunMaestro
    repo: Maestro
    description: "Desktop app for orchestrating Claude Code, Codex, OpenCode, and other agents across projects, with Auto Run and a headless CLI"
    tags: [desktop, orchestration, headless, community]
    notes: "ADE."
  - owner: huggingface
    repo: tau
    description: "Small, readable Python terminal coding agent ported from Pi, with a TUI and a print mode"
    tags: [terminal, byo-model, local-models, interactive, headless, community]
  - owner: Nano-Collective
    repo: nanocoder
    description: "Community-built terminal coding agent that runs on local Ollama models or any OpenAI-compatible API"
    tags: [terminal, byo-model, local-models, interactive, headless, mcp, community]
  - owner: MiniMax-AI
    repo: minimax-code
    description: "MiniMax's terminal coding agent with a headless exec mode, ACP support, and custom OpenAI- or Anthropic-compatible providers"
    tags: [terminal, byo-model, interactive, headless, mcp, acp, vendor]
    notes: Desktop app is not part of the published source
  - owner: nimbalyst
    repo: nimbalyst
    description: "Desktop workspace with visual editors, a session kanban, and worktrees for running Codex, Claude Code, and other agents"
    tags: [desktop, orchestration, mcp, community]
    notes: "ADE. Successor to stravu/crystal."
  - owner: coder
    repo: xum
    description: "Desktop coding agent multiplexer that runs parallel agent workspaces locally, in git worktrees, or over SSH"
    tags: [desktop, byo-model, local-models, community]
    notes: Renamed from coder/mux
  - owner: kodustech
    repo: kodus-ai
    description: "Self-hostable AI code review for pull requests on GitHub, GitLab, Bitbucket, Azure DevOps, and Forgejo, with any LLM"
    tags: [self-hosted, byo-model, review, headless, community]
    notes: AGPL-3.0, with ee/ files under a separate enterprise license
  - owner: ggml-org
    repo: llama.vscode
    description: "VS Code extension from the llama.cpp project for local code completion, chat, and agentic coding"
    tags: [editor-plugin, local-models, interactive, community]
  - owner: magiksoft
    repo: aider-desk
    description: "Desktop app for AI-assisted development built on Aider and its own agent, with git worktrees and MCP"
    tags: [desktop, byo-model, local-models, interactive, mcp, community]
  - owner: microsoft
    repo: intelligent-terminal
    description: "Windows Terminal fork that runs ACP-compatible coding agents, GitHub Copilot by default, beside your shells"
    tags: [terminal, desktop, interactive, acp, vendor]
  - owner: gastownhall
    repo: gastown
    description: "Workspace manager that coordinates Claude Code, Copilot, Codex, and other agents with git-backed work tracking and handoffs"
    tags: [terminal, orchestration, community]
    notes: "ADE."
  - owner: humanlayer
    repo: humanlayer
    description: "Runs Claude Code, Codex, OpenCode, and Pi on your machines and lets a team steer them from desktop, web, or CLI"
    tags: [desktop, web, terminal, orchestration, community]
    notes: "ADE."
  - owner: agegr
    repo: pi-web
    description: "Local browser UI for the pi coding agent that browses, resumes, and branches sessions on the same config and files"
    tags: [web, orchestration, community]
    notes: "ADE for pi only."
  - owner: vastsa
    repo: PI-Desktop
    description: "Plugin-based desktop coding agent built on the pi harness, with any provider, local models, and MCP servers"
    tags: [desktop, byo-model, local-models, interactive, mcp, community]
  - owner: zhukunpenglinyutong
    repo: jetbrains-cc-gui
    description: "IntelliJ plugin that gives Claude Code, Codex, and other coding CLIs a chat panel inside JetBrains IDEs"
    tags: [editor-plugin, orchestration, mcp, community]
    notes: "Renamed from Claude Code GUI to CC GUI."
  - owner: Gaurav-Gosain
    repo: tuios
    description: "Terminal multiplexer and tiling window manager that tracks coding agent state and can fan one prompt out to several agents"
    tags: [terminal, orchestration, acp, mcp, community]
    notes: "ADE."
```

## Tag evidence

Each line gives the README (or doc) the tags come from. "README" means `https://github.com/OWNER/REPO#readme`.

- **herdr**: [README](https://github.com/herdrdev/herdr#readme): "the runtime your coding agents live on", "runs what you already run — claude code, codex, cursor, opencode … herdr doesn't wrap or replace them; it owns their terminals"; topics `agent-orchestration`, `terminal`, `multiplexer`.
- **cmux**: [README](https://github.com/manaflow-ai/cmux#readme): Ghostty-based macOS terminal, notification rings "when coding agents need your attention", scriptable in-app browser; topics `coding-agents`, `parallel-agents`, `macos`. License split per the [LICENSE](https://github.com/manaflow-ai/cmux/blob/main/LICENSE) file.
- **happy**: [README](https://github.com/slopus/happy#readme): "Any model: Claude, Codex, Grok, Gemini", "OpenCode or any Agent Client Protocol agent plugs in as a harness", repo holds "the mobile app (iOS, Android, web), the original CLI, and the relay server". `self-hosted` left off: the relay server exists, but the README gives no self-hosting guide.
- **waveterm**: [README](https://github.com/wavetermdev/waveterm#readme): "AI-integrated terminal for macOS, Linux, and Windows", "Bring your own API keys for OpenAI, Claude, or Gemini, or run local models via Ollama and LM Studio".
- **opcode**: [README](https://github.com/winfunc/opcode#readme): "desktop application that transforms how you interact with Claude Code", "MCP Server Management".
- **Paradev**: [README](https://github.com/ParadevAI/Paradev#readme): "open-source desktop app for running CLI coding agents in parallel … Each task gets its own Git worktree"; topics `ade`, `agent-orchestration`. The old slug `superset-sh/superset` redirects here.
- **agent-orchestrator**: [README](https://github.com/OrchestratorInc/agent-orchestrator#readme): "local desktop workspace … every task gets its own agent, branch, and worktree", supported agents include Claude Code, Codex, Cursor, OpenCode, Cline, Gemini CLI; mobile apps.
- **claudecodeui**: [README](https://github.com/siteboon/claudecodeui#readme): "A desktop and mobile UI for Claude Code, Cursor CLI, and Codex", "Self-Hosted (Open source)", desktop app with "Local CloudCLI".
- **openchamber**: [README](https://github.com/openchamber/openchamber#readme): "workspace for running and reviewing AI coding work on desktop, web, VS Code, and mobile", "OpenChamber uses OpenCode to run coding agents".
- **mini-swe-agent**: [README](https://github.com/SWE-agent/mini-swe-agent#readme): "Supports all models via litellm, openrouter, portkey", "agent that solves GitHub issues or helps you in your command line", built by the SWE-bench team (`research`, matching tracked SWE-agent). `local-models` not tagged: litellm could reach them, but the README does not say so.
- **claude-squad**: [README](https://github.com/smtg-ai/claude-squad#readme): "terminal app that manages multiple Claude Code, Codex, Gemini … and Aider", "git worktrees to isolate codebases".
- **ZCode**: [README.en](https://github.com/zai-org/ZCode/blob/main/README.en.md): "AI coding workspace with desktop, browser, and terminal interfaces", `zcode` (TUI) and `zcode --web`. Model access is not documented in the README, so no model tag.
- **ThinkRail**: [README](https://github.com/JetBrains/thinkrail#readme): "agentic IDE", Monaco editor tabs, "Desktop app or browser", "MCP servers", "connect any provider through pi's own auth", "Any OpenAI-compatible server fits: LiteLLM, sglang, vLLM, Ollama, LM Studio, or llama.cpp".
- **kimi-code**: [README](https://github.com/MoonshotAI/kimi-code#readme): "AI coding agent that runs in your terminal … works out of the box with Moonshot AI's Kimi models and can also be configured to use other compatible providers" (hence `byo-model`), `/mcp-config`, `kimi acp` for Zed and JetBrains.
- **CopilotForXcode**: [README](https://github.com/intitni/CopilotForXcode#readme): "Xcode Source Editor Extension that provides GitHub Copilot, Codeium and ChatGPT support"; features Suggestions, Chat, Modification.
- **forgecode**: [README](https://github.com/tailcallhq/forgecode#readme): Interactive Mode (TUI), "Pass `-p` … to run a single prompt and exit", "MCP Configuration", "Multi-provider support", OpenAI-compatible providers.
- **codecompanion.nvim**: [README](https://github.com/olimorris/codecompanion.nvim#readme): "Neovim AI coding assistant", providers incl. Ollama and OpenRouter, "built-in support for Agent Client Protocol (ACP), Model Context Protocol (MCP)".
- **stagewise**: [README](https://github.com/stagewise-io/stagewise#readme): "open source agentic IDE for developers with a coding agent built right in", "Bring Your Own Key for all popular model providers … custom providers (including local inference!)".
- **OpenCursor**: [README](https://github.com/PawanOsman/OpenCursor#readme): "open-source AI coding agent for VS Code", API providers, custom OpenAI-compatible endpoints, "Ollama and llama.cpp", "MCP and hooks".
- **bb**: [README](https://github.com/get-bb/bb#readme) and [packages/bb-app README](https://github.com/get-bb/bb/blob/main/packages/bb-app/README.md): "agentic IDE … the desktop app, web app, CLI, and HTTP API", providers Claude Code, Codex, Cursor via ACP, Pi, OpenCode ("bb uses the provider CLI you already have authenticated"); topics `ade`, `agent-ide`.
- **mistral-vibe**: [README](https://github.com/mistralai/mistral-vibe#readme): "open-source coding agent for your terminal", sign in "with your Mistral account or an API key". No model tag: the README does not say whether other providers work, and I could not fetch the configuration doc; [docs/README.md](https://github.com/mistralai/mistral-vibe/blob/main/docs/README.md) lists "ACP Setup" and "Programmatic Mode - `vibe --prompt`". `mcp` left off: only code mentions, no doc.
- **hapi**: [README](https://github.com/tiann/hapi#readme): "control them remotely through native iOS / Android apps, Web / PWA, or Telegram", `hapi hub --relay`, "self-hosted options (Cloudflare Tunnel, Tailscale)".
- **LoopForge**: [README](https://github.com/Tencent/LoopForge#readme): "self-hosted workspace where your team ships software with coding agents", "Coding-agent CLIs do the work: CodeBuddy, Codex, Claude Code, OpenCode …", roles get "Skills and MCP tools". `web` not tagged: no web UI described in the README.
- **agent-of-empires**: [README](https://github.com/agent-of-empires/agent-of-empires#readme): "session manager for AI coding agents … monitor them from the TUI or a browser … optional worktree and container isolation".
- **gptme**: [README](https://github.com/gptme/gptme#readme): "runs anywhere a terminal runs", "Works with Anthropic, OpenAI, Google … OpenRouter … or fully local via Ollama and any OpenAI-compatible server", web UI "Self-hostable", desktop app (v0.32.1), "Integrations: MCP & ACP", plain CLI "for scripted and non-interactive use". README calls itself "a development-focused peer to self-hosted agents", which keeps it inside the scope line.
- **costrict**: [README](https://github.com/zgsm-ai/costrict#readme): VS Code extension and JetBrains plugin, "MCP Integration", "support for Anthropic, OpenAI, OpenAI-compatible APIs, and local models".
- **fx**: [README](https://github.com/vercel-labs/fx#readme): "minimal coding agent CLI", "To send one request without the interactive shell, use `fx ask`", "`fx acp`", MCP servers, "OpenAI-compatible Chat Completions endpoint such as Ollama or OpenRouter"; "Status: Experimental".
- **twinny**: [README](https://github.com/twinnydotdev/twinny#readme): "AI coding assistant for Visual Studio Code … Code completion, inline edits, chat, code review", Ollama, LM Studio, llama.cpp and hosted APIs.
- **CopilotChat.nvim**: [README](https://github.com/CopilotC-Nvim/CopilotChat.nvim#readme): "GitHub Copilot Chat capabilities directly into Neovim", "custom providers (Ollama, Mistral.ai)", tool calling.
- **Maestro**: [README](https://github.com/RunMaestro/Maestro#readme): "cross-platform desktop app for orchestrating your fleet of AI agents", supports Claude Code, Codex, OpenCode, Factory Droid; "Full CLI (`maestro-cli`) for headless operation … from cron jobs or CI/CD pipelines".
- **tau**: [README](https://github.com/huggingface/tau#readme): "coding agent that lives in your terminal", "OpenRouter, Hugging Face, and custom OpenAI-compatible endpoints, including local models", "Interactive Textual TUI and non-interactive print mode".
- **nanocoder**: [README](https://github.com/Nano-Collective/nanocoder#readme): "open coding agent for your terminal", "local models via Ollama, or any OpenAI-compatible API", "Non-interactive mode", MCP servers in configuration docs.
- **minimax-code**: [README](https://github.com/MiniMax-AI/minimax-code#readme): Headless `mcode exec`, ACP `mcode acp`, "custom providers with OpenAI- or Anthropic-compatible API formats", MCP; "The published source covers the terminal TUI, headless CLI", not the desktop app.
- **nimbalyst**: [README](https://github.com/nimbalyst/nimbalyst#readme): "visual workspace for building with Codex, Claude Code, OpenCode", "Parallel agent sessions … own git worktree", "MCP client and servers", desktop app. `stravu/crystal` says "Crystal is now Nimbalyst".
- **xum**: [README](https://github.com/coder/xum#readme): "Coding Agent Multiplexer", workspaces local, worktree or SSH, "Ollama supported for local LLMs", "OpenRouter supported". `orchestration` not tagged because it runs its own agent, not other agents' CLIs.
- **kodus-ai**: [README](https://github.com/kodustech/kodus-ai#readme): "Open source AI code review you can self-host, with the LLM you choose", "any OpenAI-compatible endpoint", "CLI + CI/CD Ready".
- **llama.vscode**: [README](https://github.com/ggml-org/llama.vscode#readme): "Local LLM-assisted text completion, chat with AI and agentic coding extension for VS Code", llama.cpp server.
- **aider-desk**: [README](https://github.com/magiksoft/aider-desk#readme): "open-source agentic platform", Git worktrees, "OpenAI, Anthropic, Gemini, DeepSeek, Ollama, and 25+ other providers", "Model Context Protocol (MCP)".
- **gastown**: [README](https://github.com/gastownhall/gastown#readme): "Multi-agent orchestration system for Claude Code, GitHub Copilot, and other AI agents with persistent work tracking", "coordinate multiple AI coding agents (Claude Code, GitHub Copilot, Codex, Gemini, and others)", `gt` CLI. A browser dashboard exists, but `web` is left off because the README treats it as secondary.
- **humanlayer**: [README](https://github.com/humanlayer/humanlayer#readme): "runs coding agents (Claude Code, Codex, OpenCode, Pi) on your machines and lets your team steer them … synced between a desktop app, a web app, and a daemon", `humanlayer` CLI.
- **pi-web**: [README](https://github.com/agegr/pi-web#readme): "Local browser UI for the pi coding agent … browse and resume conversations, run agent turns", Git worktrees.
- **PI-Desktop**: [README](https://github.com/vastsa/PI-Desktop#readme): "modular desktop workspace for AI agents", "Model-agnostic", "OpenAI · Anthropic · OpenAI-Compatible APIs · Custom Gateways · Ollama · LM Studio", "MCP Server: Connect local or remote MCP servers"; GitHub description "Local-first AI coding agent desktop … pi Agent Harness".
- **jetbrains-cc-gui**: [README](https://github.com/zhukunpenglinyutong/jetbrains-cc-gui#readme): "IntelliJ IDEA plugin that provides a visual interface for Claude Code, OpenAI Codex and more AI coding CLIs" (menu lists Grok CLI, Kimi CLI, OpenCode, PI CLI and others), "MCP server support".
- **tuios**: [README](https://github.com/Gaurav-Gosain/tuios#readme): "terminal window manager that knows what your agents are doing", "`tuios start-agent` starts a helper … headless over ACP", "`tuios fan` starts one prompt in several agents", "`tuios mcp` serves the same surface as MCP tools"; agent docs in [docs/AGENT_STATE.md](https://github.com/Gaurav-Gosain/tuios/blob/main/docs/AGENT_STATE.md).
- **intelligent-terminal**: [README](https://github.com/microsoft/intelligent-terminal#readme): "works with any Agent Client Protocol (ACP)-compatible" agent, "GitHub Copilot is the default". Topics `acp`, `terminal`, `windows-terminal`. Model access depends on the hosted agent, so no model tag.

## Borderline (your call, not in the YAML above)

Each fits the letter of the scope but repeats a pattern the list has not ruled on yet.

| Stars | owner/repo | Question |
|--:|---|---|
| 24,041 | deepseek-ai/dsh-superagent-prototype | Vendor research prototype (created 2026-10-31) of a proactive GitHub agent. It ships as a plugin bundle for the tracked deepseek-harness, and its README says "not an official DeepSeek product". Would take `[self-hosted, autonomous, research, vendor]`. |
| 28,108 | openai/symphony | Orchestrates autonomous Codex runs from a Linear board. It is mainly a spec, plus an "experimental reference implementation" labelled an "engineering preview". The last push was 2026-09-15 (66 days ago). |
| 23,963 | coleam00/Archon | YAML workflow engine that runs Claude Code, Codex and Pi. It is a self-hosted tool, not a library, but it calls itself a "harness builder". |
| 14,211 | openai/codex-security | Open since 09-29: is a security scanner and fixer a review tool? |
| 8,453 | vercel-labs/deepsec | Same question: an agent-powered vulnerability scanner. |
| 13,783 | backnotprop/plannotator | Browser review surface for agent plans and diffs. It sends feedback to the agent; the model is not inside it. |
| 10,025 | anthropics/claude-code-action | GitHub Action packaging of tracked Claude Code. Same question for `openai/codex-action` (1.3k) and `github/gh-aw` (5.7k), which compiles agentic workflows into Actions. |
| 6,281 | generalaction/emdash | ADE that qualifies on everything else. The README was renamed to "Emdash Classic" today (2026-11-20), which suggests a successor product. Recheck in a week. |
| 5,480 | cloudflare/vibesdk | Self-deployable prompt-to-app platform (`app-builder`) that also markets itself as a base for building your own platform. |
| 4,092 / 2,822 / 3,162 / 2,205 | copilot.lua, folke/sidekick.nvim, coder/claudecode.nvim, ggml-org/llama.vim | Thin editor bridges to Copilot, Claude Code or llama.cpp. They are in scope as AI editor plugins but have low independent substance. |
| 11,456 | omnigent-ai/omnigent | ADE that qualifies on surface (desktop, web, terminal, phone sessions over Claude Code, Codex, OpenCode and others), but its GitHub description says "AI agent framework and meta-harness" and it carries the `agent-framework` topic. That is the same wording that put Archon here and ruflo in the rejected list. Alpha. Would take `[desktop, web, terminal, byo-model, orchestration, community]`. |
| 6,227 | tw93/Kaku | macOS terminal with agent status, Claude Code notifications and an optional AI chat. The AI is optional rather than central, unlike waveterm. |
| 4,598 | convex/chef | Prompt-to-app builder forked from bolt.diy. It passes the written rule (`pushed_at` 2026-10-28), but its last default-branch commit was 2026-03-13. I would not add it under a stricter last-commit reading. |
| 2,530 | mattzcarey/shippie | PR review agent. It passes the written rule (`pushed_at` 2026-10-21), but its last default-branch commit was 2026-07-02. Same stricter-reading caveat. |
| 1,061 | editor-code-assistant/eca | Qualifies (editor-agnostic agent server with Emacs, VS Code, IntelliJ clients, MCP, Ollama), but it sits 61 stars above the floor. |

## Considered and rejected

| owner/repo | Stars | Rule that excludes it |
|---|--:|---|
| continuedev/continue | 36,527 | Declared unmaintained and read-only in the README. Removed earlier for the same reason. |
| BloopAI/vibe-kanban | 28,754 | Declared sunsetting in the README. |
| CommandCodeAI/command-code | 4,253 | No license and no source: the repo holds only a README and `mods/`. |
| kirodotdev/Kiro | 4,483 | No source: docs and an issue tracker only. |
| sourcery-ai/sourcery | 1,888 | No source: README and a pre-commit hook config. |
| chenhg5/cc-connect | 16,446 | No license file. |
| NanmiCoder/cc-haha | 15,501 | Not open source: "Community Source" license. |
| anus-dev/ANUS | 6,560 | Rebrand of tracked pi ("installed as a dependency: its terminal interface, tools, sessions and providers") with a free-model router. |
| anywhere-labs/dsh-desktop | 36,484 | Desktop shell that "pins and runs a specific upstream version" of tracked deepseek-harness unchanged. |
| ultraworkers/claw-code | 194,758 | README: "Claw Code is not the serious production project here." |
| Orkas-AI/Orkas | 9,474 | General multi-agent "AI assistant that works as a team" (digital-workforce exclusion). Earlier added by PR and then removed. |
| iOfficeAI/AionUi | 34,403 | General "Cowork app". |
| different-ai/openwork | 24,571 | "Alternative to Claude Cowork and Codex … where AI agents do real work on your own files", which makes it a general assistant. |
| yc-software/qm | 15,749 | "Multiplayer agent harness for teams". The 09-29 round judged it a Slack and web general work agent; nothing in it has changed. |
| EKKOLearnAI/ekko-studio | 11,805 | General multi-agent chat and workflow workspace. |
| NousResearch/hermes-agent, paperclipai/paperclip, nanocoai/nanoclaw, makecindy/cindy | 10k-265k | General assistants or agent-management apps, not software-development tools. |
| github/squad | 3,418 | Agent-team configuration installed into other hosts, closer to a prompt/agent collection. Its own shell is deprecated. |
| code-yeongyu/oh-my-openagent, Yeachan-Heo/oh-my-claudecode, alvinunreal/oh-my-opencode-slim | 10k-71k | Plugin bundles for existing hosts. |
| farion1231/cc-switch, router-for-me/CLIProxyAPI, cheaper-inference/OmniRoute, yetone/magpie, Wei-Shaw/sub2api | 39k-156k | Provider switchers and API proxies; not tools where the LLM is central. |
| graphify, claude-mem, context7, serena, codebase-memory-mcp, code-review-graph, agentmemory, rtk, headroom, context-mode | 12k-151k | Context, memory or token infrastructure (MCP servers, plugins) that you plug into an agent. These are building blocks, not tools. |
| github/spec-kit, google/elephant-goldfish-model, DietrichGebert/ponytail, OthmanAdi/planning-with-files | 21k-195k | Toolkits, skills or plugins. |
| ruvnet/ruflo, rocketride-org/rocketride-server, open-multi-agent | 7k-82k | Agent frameworks or harness runtimes meant to be built on. |
| max-sixty/worktrunk, dagger/container-use | 12.7k / 4.1k | Worktree and sandbox utilities for agents; no LLM inside. |
| modem-dev/hunk | 9,855 | Terminal diff viewer; AI annotations are a side feature, not central. |
| steipete/CodexBar, getagentseal/codeburn, graykode/abtop | 3.9k-23k | Usage and observation dashboards. |
| anthropics/oss-scanner | 1,838 | Enrollment repo for a hosted Anthropic service. |
| omnirush-ai/omnirush-gui | 5,072 | Account-only models ("Your own provider keys are not supported"), open-core licence, origin of the free frontier access unclear. Not recommended on quality. |
| getkimchi/kimchi | 2,287 | Built on the pi SDK and bound to the commercial kimchi.dev service. Not recommended on quality; revisit if wanted. |
| vercel-labs/open-agents | 5,906 | "An open source template for building cloud agents", so a building block. |
| plandex-ai/plandex, bytedance/trae-agent, stackblitz/bolt.new, sweepai/sweep | 7.7k-17k | Stale (last push 2025-10 to 2026-02). |
| RooCodeInc/Roo-Code, voideditor/void, smallcloudai/refact, coder/agentapi | 1.5k-29k | Archived. |
| JetBrains/junie, trypear/pearai-app | 570 / 715 | Below the 1,000-star floor. |

## Housekeeping found along the way (not candidates)

- Tracked `Hmbown/Codewhale` now resolves to `codewhale-hq/Codewhale` (41.2k). Per CONTRIBUTING, update the slug and add a `canonicalKeyMigrations` entry.
- `superset-sh/superset` now redirects to `ParadevAI/Paradev`, and `coder/mux` to `coder/xum`. Use the new slugs above.

## Limitations

- **Seen, not reviewed:** these surfaced in the searches with in-scope one-liners, but their READMEs were not checked: zhukunpenglinyutong/desktop-cc-gui (4.6k), spacering-net/codeg (4.5k), zeronsh/zeron (3.7k), fuxicodex/Fuxi (3.2k), chaterm/Chaterm (3.1k), NeuralNomadsAI/CodeNomad (2.7k), happier-dev/happier (2.3k), Dicklesworthstone/pi_agent_rust (2.0k), awslabs/cli-agent-orchestrator (1.6k), HexmosTech/git-lrc (1.5k), johannesjo/parallel-code (1.1k), supabitapp/supacode (2.4k), egoist/waku (1.6k), SeemSeam/claude_codex_bridge (3.6k), tontinton/maki (1.4k), milanglacier/minuet-ai.nvim (1.5k), madox2/vim-ai (1.2k), appleboy/CodeGPT (1.5k), kdlbs/kandev (1.2k), open-mercato/cezar (1.1k) and similar. Expect several more qualifying ADEs and small agents in this band.
- Rows dropped by my first-pass keyword filter (words such as prompt, skill, memory, token, design, sdk) were not skimmed one by one.
- The model-lab org scan read only the 40 newest repos per org, so for large orgs it may not reach back to August 2026.
- Discovery relies on GitHub search ranking, topics and keywords. A repo with no matching topic and no matching words in its name or description could still be missed. Model-lab orgs were scanned only for repos created since 2026-08-01.
- The web searches returned mostly secondary sources, and only `deepseek-ai/dsh-superagent-prototype` was new. I did not use any web source as tag evidence.
- Tags come from README text and linked docs only. The `docs/` trees of most repos were not read, so `headless`, `mcp` and `acp` are probably under-tagged in places.
- Star counts in the 1k-5k band move fast. Re-check any entry near the floor (aider-desk 1.47k, llama.vscode 1.56k, kodus-ai 1.63k) before merging.
- Scope calls on ADEs that front a single agent (opcode for Claude Code, openchamber for OpenCode, pi-web for pi) follow the t3code precedent, but they are still judgment.

## Unresolved questions

1. Do you want all 45 added at once, or a first batch (the twelve named under Outcome)? Twenty are ADEs, against three tracked today (orca, paseo, t3code).
2. Should Microsoft (`intelligent-terminal`) and Tencent (`LoopForge`) carry `vendor`? The vocabulary defines `vendor` as model labs; the tracked `github/copilot-cli` already sets a precedent for GitHub/Microsoft.
3. For the borderline rows: should GitHub Action packagings of tracked agents, security scanners (codex-security, deepsec) and spec-first orchestrators (symphony) count?
4. Should a 1k-star floor still mean adding the long tail of small ADEs (see "Seen, not reviewed"), or do you want a separate quality bar for that band?
5. Should vendor research prototypes that ship as plugins for a tracked agent (dsh-superagent-prototype) be eligible?
