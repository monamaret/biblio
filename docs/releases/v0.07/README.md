# biblio-cli — v0.07 User Guide (Bug-fix release)

**Release**: Bug-fix release
**Status**: builds on [v0.06](../v0.06/README.md) (General Availability) — the
same CLI, content pipeline, theming, MCP server, and terminal/negotiation
layer; this release fixes nine bugs filed against v0.06's General
Availability surface, with no new commands or roadmap deliverables. Unlike
v0.01–v0.06, this release wasn't planned against a roadmap item — it's the
result of triaging real issues found using the GA build.

## What's new

Nine issues filed against v0.06, all fixed:

- **`tokens.css` was missing from production builds** (#191) — it lived at
  the project root, which `vite build` never copies (only `publicDir`
  does); `vite dev` masked this because its dev server serves arbitrary
  root-relative files. The scaffold generator now writes it to `public/`
  instead, so deployed sites are no longer unstyled. The served URL is
  unchanged (`/tokens.css`).
- **`biblio-cli dev` left orphaned `vite`/`wrangler` processes on exit**
  (#192) — `npx vite dev`/`npx wrangler dev` spawn the real dev server as a
  *child* process rather than exec'ing into it; killing only the `npx`
  leader orphaned that child (and could hang the supervisor itself, since
  `Wait()` blocks on the inherited output pipe until every holder closes
  it). Each process now starts in its own process group and the whole
  group is killed on every shutdown path.
- **`biblio-cli -v` always printed `dev`** (#193) — release binaries had no
  `-ldflags` set, and `go install .../biblio-cli@latest` (how CI installs
  it) doesn't apply ldflags either way. Goreleaser now stamps the release
  tag in; the CLI also falls back to `runtime/debug.ReadBuildInfo()` for
  the `go install` path, which Go populates automatically.
- **`biblio-cli dev` failed on the very first run** (#194) — `wrangler dev`
  requires its assets directory to exist at startup, before any build has
  ever run, but nothing created `dist/client` until a build did. `dev` now
  creates the configured directory itself, right before launching.
- **Project names with dots/spaces broke `wrangler.jsonc`** (#195) —
  Cloudflare Workers script names only allow lowercase letters, digits,
  and hyphens, but the literal project name (e.g. a domain-shaped name
  like `notebook.example.com`) was written straight into `wrangler.jsonc`'s
  `"name"` field. `new` now derives a sanitized, Workers-safe slug for that
  field only — the project directory and `content/site.yaml`'s `name` are
  unaffected — and prints the substitution when it applies.
- **No workflow ever published a release binary** (#196) — `.goreleaser.yaml`
  existed, but there was no `.github/workflows/release.yml` to run it.
  Added one, triggered on `v*` tags, gating on `go test ./...` first. This
  release is the first to actually exercise it.
- **The README didn't explain what this project does or how to operate it
  day to day** (#197) — rewritten with a clear framing (a static-site
  scaffold that also runs an in-Worker MCP server, not a general site
  generator), plus new local-dev, troubleshooting, testing, MCP-App-UI-bundle,
  and deployed-site-testing sections.
- **No documentation of how the coding-agent tooling fits together** (#198)
  — added [docs/agent-integration-guide.md](../../agent-integration-guide.md),
  covering the Spec Kit workflow, opencode wiring, and Claude Code plugin
  config, linked from the README.
- **Local agent configuration policy was undocumented and unenforced**
  (#199) — audit found nothing actually committed in violation, but
  `CLAUDE.md` inaccurately described `.claude/settings.json`'s contents;
  corrected it, stated the policy explicitly, and added `.gitignore`
  entries so personal/editor-specific agent config can't be committed by
  accident later.

## The complete `biblio-cli` lifecycle

Unchanged from v0.06 — see the [v0.06 guide](../v0.06/README.md#the-complete-biblio-cli-lifecycle)
for create/configure/update/deploy/destroy. The one addition: a real,
binary release now exists, so "download a prebuilt binary" in the
top-level [README.md](../../../README.md#install) actually has something
to download as of this release.

### What an AI assistant can do

Unchanged from v0.05/v0.06 — connect an MCP client to
`https://<your-domain>/mcp`; see the
[v0.05 guide](../v0.05/README.md#what-an-ai-assistant-can-do).

## Known scope limits in v0.07

- This release fixes bugs found in v0.06's GA surface; it doesn't add any
  new CLI command, site type, or roadmap deliverable.
- The Windows-specific orphaned-process fix from #192 keeps the pre-fix
  single-process kill behavior (no POSIX process groups there) — CI only
  exercises `ubuntu-latest` today, so this is unverified on a real Windows
  machine.
