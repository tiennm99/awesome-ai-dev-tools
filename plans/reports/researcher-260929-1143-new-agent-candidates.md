# New agent candidates, 2026-09-29

Method: reused the 2026-09-11 refresh. About 30 `gh search repos` queries and 13 topic searches (stars>9000, sorted by stars), a created>2026-06-01 sweep, GraphQL/REST verification of stars, pushedAt and archived state, README reads for every shortlisted repo, plus 2 web searches (which found nothing new). Cutoff for staleness: pushedAt >= 2026-06-29. Star and push data are live as of 2026-09-29.

## Outcome

Three clear additions and one borderline. The list is nearly saturated: after the 09-11 refresh, the only new qualifying repo published since then is a vendor agent (Xiaomi MiMo Code). The other two are older repos the prior run did not catch (Open SWE, PR-Agent).

| Stars | owner/repo | pushedAt | Verdict |
|--:|---|---|---|
| 13,541 | XiaomiMiMo/MiMo-Code | 2026-09-28 | Add (new vendor agent) |
| 13,179 | The-PR-Agent/pr-agent | 2026-09-28 | Add (review agent, same precedent as alibaba/open-code-review) |
| 10,772 | langchain-ai/open-swe | 2026-09-29 | Add (close to the 10k floor) |
| 10,879 | openai/codex-security | 2026-09-29 | Borderline, your call |

## Ready-to-paste YAML

```yaml
  - owner: XiaomiMiMo
    repo: MiMo-Code
    description: "Terminal coding assistant from Xiaomi with persistent project memory, subagents and any OpenAI-compatible provider"
    tags: [terminal, byo-model, interactive, mcp, headless, vendor]
    notes: Fork of opencode; a separate desktop beta is advertised but not tagged
  - owner: The-PR-Agent
    repo: pr-agent
    description: "Code review agent for pull requests on GitHub, GitLab, Bitbucket and more, run via CLI, GitHub Action or webhook"
    tags: [terminal, self-hosted, byo-model, local-models, review, headless, community]
    notes: Community-maintained legacy project of Qodo, not the Qodo product
  - owner: langchain-ai
    repo: open-swe
    description: "Asynchronous coding agent that investigates issues, opens pull requests and reviews them, built on LangChain Deep Agents"
    tags: [self-hosted, web, byo-model, autonomous, review, mcp, community]
    notes: Under active development; README says no external issues or contributions accepted
```

Borderline (paste only if you accept a security scanner as "reviews code"):

```yaml
  - owner: openai
    repo: codex-security
    description: "OpenAI CLI and SDK that scans code for security vulnerabilities, validates findings and proposes fixes"
    tags: [terminal, single-vendor, review, headless, vendor]
```

## Tag evidence

- **MiMo-Code** (MIT): README says "terminal-native AI coding assistant", "connecting to any mainstream LLM provider API" and Custom Provider for any OpenAI-compatible API (byo-model), MCP server connections, and a `mimo run` headless mode. Xiaomi MiMo is a model lab, so vendor. I did not tag local-models, acp or editor-plugin because the README shows no support. The "ACP" hit is a Windows code page.
- **pr-agent** (Apache-2.0): README lists CLI local usage, GitHub Action, Docker, self-hosted, webhooks, and any model via LiteLLM including Ollama. It is a code review agent, so review. Qodo describes it as community-maintained, hence community.
- **open-swe** (MIT): README describes investigate, implement, open PR, and separate PR review. Triggers are dashboard, GitHub, Slack, and Linear via MCP. Models and sandboxes are configurable. Deployed as a team server. LangChain is not a model lab, so community. Terminal is omitted although a CLI exists, because it needs a deployed server; add it if you prefer.
- **codex-security** (Apache-2.0): README says CLI and SDK for "finding, validating, and fixing security vulnerabilities", CI via API key. It needs Trusted Access for Cyber for some requests.

## Near-misses and rejections

Out of scope (criterion 1, not an agent that writes, edits or reviews code):
- `openai/symphony` 27.5k: orchestrates agent runs.
- `coleam00/Archon` 23.6k: harness builder.
- `github/spec-kit` 139k: spec-driven toolkit.
- `anywhere-labs/dsh-desktop` 29.4k and `dataelement/dsh-desktop` 10.4k: desktop shells around the already tracked deepseek-harness, so wrappers.
- `EKKOLearnAI/ekko-studio` 11.2k: workspace for running other agents.
- `yc-software/qm` 15.3k: Slack/web multiplayer general work agent.
- `OpenBMB/ChatDev` 34.4k: README says it is now a zero-code multi-agent orchestration platform.
- `Fosowl/agenticSeek` 27.4k: general Manus-style assistant (web browsing, planning), Python.
- `FoundationAgents/OpenManus` 58.4k, `lsdefine/GenericAgent` 14.3k, `NousResearch/hermes-agent` 249.9k: general agents.
- `Untrivial-ai/agent-orchestrator` 12.5k, `omnigent-ai/omnigent` 10.3k, `nexu-io/open-design` 98.5k, `cloudflare/security-audit-skill` 22.8k, `langchain-ai/openwiki` 16.8k (docs writer), `DietrichGebert/ponytail` 147.7k (prompt skill): orchestrators, plugins, skills or non-coding.
- Previously rejected and unchanged: claw-code, oh-my-openagent, ruflo, orca, superset, vibe-kanban, herdr, cc-switch, routers and skill packs.

Below 10k: `mistralai/mistral-vibe` 5.0k, `gptme/gptme` 4.4k, `cloudflare/vibesdk` 5.4k, `vercel-labs/open-agents` 5.8k, `dagger/container-use` 4.0k, `sweepai/sweep` 7.7k (also stale).

Stale or archived: `microsoft/vscode-copilot-chat` 9.97k (archived, last push 2026-05-20), `firecrawl/open-lovable` 28.6k (last push 2025-11-19), `stackblitz-labs/bolt.diy` 19.9k (2026-02-07), and the abandoned repos rejected on 09-11.

## Limitations

- Discovery relies on GitHub search ranking. A qualifying repo with no matching keywords or topics and no appearance in the 09-11 probe could be missed.
- Scope calls on review-only and security tools are judgment; the criteria wording ("writes, edits, or reviews") supports pr-agent, and codex-security is the weakest case.
- Web search returned no new vendor agents beyond what GitHub already showed.

## Unresolved questions

1. Include `openai/codex-security` (security scanner and fixer) as a reviewing agent?
2. Should open-swe also carry `terminal`, given the CLI needs a deployed backend?
