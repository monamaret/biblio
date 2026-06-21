# biblio-cli

A single Go binary that scaffolds a **static site that also runs an in-Worker MCP server** — one guided command produces a complete, runnable React + Vite + React Router v7 project that builds to static assets and deploys as a single Cloudflare Worker. That Worker serves the prerendered site *and* an `/mcp` endpoint (via Durable Objects) from the same deployment, plus small MCP App UI resource bundles that richer MCP hosts can render inline alongside a tool's text answer.

This is not a general-purpose site generator: every scaffolded project is one of three fixed site types (`profile`, `log`, `docs`), content is plain Markdown + YAML frontmatter (not MDX, not a CMS), and the whole lifecycle — create, configure, preview, build, deploy, update — is meant to be driven through `biblio-cli`, not by hand-editing the generated Vite/Wrangler config.

## Install

**Download a prebuilt binary**:

Go to [Releases](../../releases) and download the binary for your platform, then move it onto your `PATH`.

**Or install from source** (requires Go 1.24+, see `go.mod`):

```sh
go install github.com/monamaret/biblio/cmd/biblio-cli@latest
```

**Or build from a clone**, with no install step:

```sh
git clone https://github.com/monamaret/biblio.git
cd biblio
go build -o biblio-cli ./cmd/biblio-cli
./biblio-cli --help
```

## Usage

### Scaffold a new site

```sh
biblio-cli new my-site --type=profile --domain=example.com --theme=profile-default
```

Site types: `profile`, `log`, `docs`. Omit any flag to be prompted interactively. See [docs/releases/v0.06/README.md](docs/releases/v0.06/README.md) for the full create → configure → update → deploy lifecycle, including how to prepare a Cloudflare account/domain for deployment.

### Index authored content

```sh
cd my-site
biblio-cli index
```

Reads every `.md` file under `content/`, validates frontmatter, and writes `content-index.json`. Pass an explicit path to target a different content directory.

### Local dev

```sh
biblio-cli dev
```

Starts `vite dev` and `wrangler dev` together as one supervised pair — the former serves the site with hot reload, the latter emulates the Worker (including the `/mcp` Durable Object) at its own local port. Press `Ctrl+C` to stop both; the supervisor kills each process's entire process group, not just its top-level command, so `npx`-wrapped tooling doesn't leave orphaned dev servers behind.

## Troubleshooting

- **`wrangler.jsonc not found`** from `build`/`dev`/`deploy` — you're not in a scaffolded project directory; run `biblio-cli new` first.
- **`wrangler dev` fails on the very first run** — `dist/client` (wrangler's assets directory) doesn't exist yet on a brand-new scaffold or a fresh clone; `biblio-cli dev` creates it automatically before launching, so this should only surface if you're invoking `wrangler dev`/`vite build` directly instead of through `biblio-cli`.
- **A project named after its domain (e.g. `my-site.example.com`) won't deploy** — Cloudflare Workers script names only allow lowercase letters, digits, and hyphens. `biblio-cli new` derives a sanitized name for `wrangler.jsonc`'s `"name"` field automatically and prints it when it differs from the literal project name you typed; the directory name and `content/site.yaml`'s `name` are unaffected.
- **The deployed site has no styling** — make sure you built with `biblio-cli build`, not a bare `vite build` from an old checkout; `tokens.css` is served from `public/` so Vite's build step picks it up.
- **`biblio-cli -v` prints something other than a clean version** — expected for a plain local `go build`/`go test` (you'll see `dev` or a VCS-derived pseudo-version); a real version tag only shows up in binaries from [Releases](../../releases) or `go install .../biblio-cli@latest`.

## Testing

**This repo (the Go CLI itself)**:

```sh
go test ./...          # full unit + integration suite
./scripts/check.sh      # gofmt -l . (fails on unformatted files) + go vet ./...
```

**A scaffolded site** (its own independent npm project — there's no repo-root Node tooling here):

```sh
cd my-site
npm install
npx vitest run          # unit + component tests
npx tsc --noEmit         # typecheck (no separate lint script is configured)
```

## MCP App UI bundle conventions

Tools that return a chart-like result (e.g. `templates/site/*/app/mcp/tools/*.ts`) pair two things in their response: a plain-text/Markdown `content` block (the answer every MCP client can render, including terminal-based ones) and an optional `resource` block carrying a small HTML view — never the other way around. The HTML view (`templates/site/*/app/mcp/ui/chart-view.tsx`) is server-rendered with `react-dom/server`'s `renderToStaticMarkup`, not bundled as a separate build step, and styles itself against the site's `tokens.css` using `var(--token, fallback)` — so it degrades to a plain, usable view if the sandboxed context can't load `tokens.css` at all, with no JS-level error handling needed. When adding a new UI view, follow this same shape: a small, dependency-free React component using CSS-variable fallbacks, rendered to a static HTML string, attached as a `resource` alongside (never instead of) a text-first `content` block.

## Testing the deployed site's views

Every non-`/mcp` request is content-negotiated by `Accept`/`User-Agent`, so the same URL serves three different representations — `text/markdown` only wins when a client actually asks for it; a bare `Accept: */*` (curl's default) always gets HTML:

```sh
curl https://my-site.example.com/                                                       # default Accept: */* → HTML
curl -H "Accept: text/markdown" https://my-site.example.com/                            # curl's own User-Agent → ANSI-rendered terminal output
curl -H "Accept: text/markdown" -A "myscript/1.0" https://my-site.example.com/          # non-terminal User-Agent → plain Markdown
```

(The ANSI tier is selected for any client identifying as `curl`/`wget`/`httpie`/`lynx` *once* `text/markdown` has already won the Accept negotiation — pass a different `User-Agent` to see the plain-Markdown tier instead.)

For the MCP endpoint, use an actual MCP client rather than raw `curl` (it's a stateful, session-based protocol) — the [MCP Inspector](https://github.com/modelcontextprotocol/inspector) is the quickest way to connect and exercise a deployed site's tools interactively:

```sh
npx @modelcontextprotocol/inspector https://my-site.example.com/mcp
```

For agent/crawler discovery, check `https://my-site.example.com/llms.txt` — a build-time-generated summary of the site and its pages.

## Further reading

- [docs/agent-integration-guide.md](docs/agent-integration-guide.md) — how this repo's coding-agent tooling (Spec Kit workflow, opencode, Claude Code plugin config) fits together.
