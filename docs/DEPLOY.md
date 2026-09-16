# Deploying to Cloudflare Pages

The site is published by Cloudflare Pages' Git integration. Cloudflare builds
from the repository on every push to `main` and never holds a GitHub token.

## How the split works

Data refresh and site build are two separate steps, in two separate places:

| Step | Command | Runs where | Needs a token? | Writes |
| --- | --- | --- | --- | --- |
| Update | `go run .` | GitHub Actions (`update.yml`, nightly + manual) | Yes — `GITHUB_TOKEN` | `README.md`, `data/history.jsonl`, `data/metadata.json` (all committed) |
| Build | `go run . -build` | Cloudflare Pages | No | `dist/` (never committed) |

`data/metadata.json` is the handoff. The update step records every GitHub-sourced
field there — stars, language, description, push date, archived flag — so the
build step can render the dashboard from committed files alone.

That has two consequences worth knowing:

- Cloudflare's build environment never sees a GitHub token, because it has no
  reason to call the GitHub API.
- A pure curation change (retagging an entry, editing a note in
  `data/agents.yml`) republishes as soon as you push, reusing the last fetched
  star figures. You do not wait for the nightly run.

The nightly Actions run commits refreshed data to `main`; that push fires
Cloudflare's build webhook, which redeploys with the new numbers.

## One-time setup

### 1. Bootstrap `data/metadata.json`

The build fails without it, so generate it before connecting Cloudflare. Either
trigger the **Update rankings** workflow manually (Actions tab →
*Update rankings* → *Run workflow*), or run the updater locally and commit:

```bash
export GITHUB_TOKEN=ghp_your_token_here
go run .
git add data/metadata.json data/history.jsonl README.md
git commit -m "chore: bootstrap metadata snapshot"
git push
```

### 2. Create the Pages project

Cloudflare dashboard → **Workers & Pages** → **Create** → **Pages** →
**Connect to Git** → pick this repository, then set:

| Setting | Value |
| --- | --- |
| Production branch | `main` |
| Framework preset | None |
| Build command | `go run . -build` |
| Build output directory | `dist` |
| Root directory | *(leave blank)* |

### 3. Set the build environment variable

Under **Settings → Environment variables → Production** (and Preview, if you
want PR previews) add:

| Variable | Value |
| --- | --- |
| `GO_VERSION` | `1.23` |

Do **not** add `GITHUB_TOKEN` — the build does not use one, and adding it would
hand a credential to an environment that has no need for it.

Cloudflare's build image ships Go and honours `GO_VERSION`. Even on an older
image, Go's toolchain directive in `go.mod` downloads the matching toolchain
automatically.

### 4. Deploy

Save and deploy. Subsequent pushes to `main` — yours and the nightly bot's —
redeploy automatically.

## Caching

`site/_headers` marks `data.json` as `must-revalidate`, so the edge cannot serve
yesterday's ranking after a refresh. `index.html` and everything else in `site/`
are copied into `dist/` as-is and use Cloudflare's defaults.

## Troubleshooting

**Build fails with "data/metadata.json not found"** — the bootstrap in step 1
has not been committed yet.

**A newly added tool is missing from the site** — expected between merging the
`agents.yml` entry and the next update run. The build logs a warning and omits
entries it has no metadata for, rather than failing the deploy. Trigger the
*Update rankings* workflow manually to fetch it immediately.

**Stars look stale** — check the Actions tab: publishing is healthy, but the
update workflow has not committed recently. The dashboard's "updated" timestamp
reports when the data was fetched, not when the site was built, so a redeploy
never makes stale figures look fresh.
