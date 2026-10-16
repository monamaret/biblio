# Changelog

All notable user-facing changes to `biblio-cli` and its generated sites are
recorded here, release by release. See `docs/releases/v<version>/` for the
full user guide of each release.

## v0.08 — The `quarterly` site type (unreleased draft)

- Added a fourth site type, `quarterly`: numbered issues, category pages,
  footnotes and margin notes (CSS-only sidenotes on wide screens),
  configurable typesetting in `content/site.yaml`, grayscale image plates,
  issue drafts, and four MCP tools (`list_issues`, `get_issue`,
  `list_articles`, `get_article`). See [docs/releases/v0.08](v0.08/README.md).
- `biblio-cli build` and `dev` now validate `content/site.yaml` and
  regenerate `content-index.json` automatically for every site type; running
  `biblio-cli index` by hand is no longer needed.
- `update`/`upgrade` no longer treat `content-index.json`,
  `public/tokens.css` or `public/plates/` as template files, so regenerated
  build output can't cause update conflicts.
- Every page now has its own `<title>` built from the site name in
  `content/site.yaml` — e.g. `On Margins · 502.io`, or just the site name
  on the home page — instead of a hard-coded `biblio … site` on every page,
  for every site type. `content-index.json` now carries a `site` object
  (name and description) for profile, log, and docs sites too (#215).
- `biblio-cli new` now builds `content-index.json` from the new project's
  own `site.yaml`, so a fresh quarterly site no longer carries the demo
  name "The Quarterly", and the first `build` leaves the index unchanged
  for every site type. The quarterly starter colophon no longer names a
  publication (#216).
- `llms.txt` no longer starts with a blank paragraph when the site
  description is empty, for every site type (#217).
- Every site type now ships a baseline `.gitignore` covering build output,
  `.dev.vars`, `.env`/`.env.*` (keeping `.env.example`), `*.log`, and
  `coverage/` (#219).
- The quarterly deploy workflow now uses Node.js 22, and every site type's
  deploy workflow gains a `workflow_dispatch` trigger for manual redeploys
  (#218).
- Docs sites now ship a `package-lock.json` like the other site types, so
  `npm install` on a fresh docs site no longer fails with an `ERESOLVE`
  conflict on `ai@7` (and `tsc` no longer trips over a duplicated MCP SDK).
  `biblio-cli update` adds the lockfile to existing docs sites; if a site
  already has its own, `update` flags it for review instead (#223).

## v0.07 — Bug-fix release

- Fixed `tokens.css` being silently dropped from production builds: the
  scaffold generator now writes it to `public/` (where `vite build` picks
  it up) instead of the project root (#191).
- Fixed `biblio-cli dev` orphaning `vite`/`wrangler` processes (and
  sometimes hanging itself) on exit by group-killing each wrapped process
  instead of only its top-level `npx` leader (#192).
- Fixed `biblio-cli -v` always printing `dev`: goreleaser now stamps the
  release tag into release binaries via `-ldflags`, and `go install
  .../biblio-cli@latest` now falls back to `runtime/debug.ReadBuildInfo()`
  (#193).
- Fixed `biblio-cli dev` failing on its very first run by creating the
  configured assets directory (e.g. `dist/client`) before launching,
  rather than requiring a prior build to have created it (#194).
- Fixed `wrangler.jsonc` getting an invalid Workers script name from
  project names containing dots/spaces (e.g. a domain-shaped name) — `new`
  now derives a sanitized slug for that field only, leaving the project
  directory and `content/site.yaml`'s `name` untouched (#195).
- Added `.github/workflows/release.yml` — the first workflow to actually
  build and publish `biblio-cli`'s `goreleaser` binaries on a `v*` tag
  push; this release is the first to use it (#196).
- Rewrote `README.md`: clearer framing of what this scaffolds, plus new
  local-dev, troubleshooting, testing, MCP-App-UI-bundle, and
  deployed-site-testing sections (#197).
- Added `docs/agent-integration-guide.md`, documenting how the Spec Kit
  workflow, opencode, and Claude Code plugin config fit together, linked
  from the README (#198).
- Corrected `CLAUDE.md`'s inaccurate description of `.claude/settings.json`
  and added `.gitignore` entries for personal/editor-specific agent config
  (`.claude/settings.local.json`, `.cursor/`, `.windsurf/`, `.clinerules/`)
  so it can't be committed by accident (#199).

## v0.06 — General Availability

- Added `biblio-cli update` and `biblio-cli upgrade <version>`: classify
  every region of an existing project against the scaffold-template it was
  generated from (scaffold-tracked / author-customized / ambiguous),
  advance scaffold-tracked regions, preserve author-customized regions
  byte-for-byte, and surface any ambiguous region as a conflict the author
  resolves explicitly via `--resolve <path>=keep-author|use-scaffold` —
  never silently overwritten.
- Extended `biblio-cli new` to write a `.biblio/scaffold-manifest.json`
  baseline (a content hash per embedded template file) at scaffold time,
  the 3-way-diff baseline `update`/`upgrade` need.
- Confirmed the update/upgrade contract holds, unmodified, across all three
  site types (profile, project log, docs).
- Added `internal/escalation`'s Storage-Escalation Guidance Check:
  documented, checkable entry-count thresholds for the two conditional
  escalation paths (D1, Vectorize), reported against each site's real
  `content-index.json` scale. Recorded real local-only dry-run proofs in
  `docs/storage-escalation-guide.md` that both paths are executable on
  demand, without adopting either as a running dependency.
- Closed the GA Readiness Review (`test/integration/ga-readiness-review/
  GA-READINESS-RECORD.md`): every roadmap deliverable §1–§10 is either
  complete-and-demonstrated or — for the two conditional storage-escalation
  items — proven-executable-on-demand. `ga-determination: true` —
  **the platform is Generally Available.**

## v0.05 — Prove the platform on documentation and across all site types

- Brought `templates/site/docs` up to parity with `templates/site/profile`'s
  React Router v7 / Worker / MCP / negotiation buildout (it had remained the
  thin v0.01-era placeholder through v0.02–v0.04).
- Extended `internal/content/schema.go`'s frontmatter validation with a
  `section` requirement (alongside the existing `order`) and structural
  checks for both, for the `docs` site type.
- Added a TOC/hierarchy-derivation module (`toc()`) grouping pages into
  sections ordered by their lowest-order page, with a deterministic
  slug-based tiebreak default for pages sharing a section/order — the
  single owner shared by the browser TOC view and the MCP tools.
- Added a TOC/navigation home view and page-detail view for the `docs` site
  type, with a dedicated test suite asserting the deliberate *absence* of
  any tag/project filter UI or query surface.
- Added docs-type-specific MCP tools (`list_sections`, `get_doc_page`),
  registered alongside the shared tools; `get_doc_page` is named distinctly
  from the shared `get_page` tool to avoid a duplicate-registration
  collision. A log-shaped filter argument against `list_sections` is
  rejected with a structured error; a profile/log-only tool against a docs
  site returns the MCP SDK's standard "unknown tool" error.
- Closed the cross-type consistency review now that all three site types
  exist: confirmed each type renders its own theme with no cross-type
  bleed-through, confirmed one shared negotiation/rendering code path
  across all three (no per-type fork), and confirmed each type's content
  shape is appropriate to its own model — recorded in
  `test/integration/cross-type-review/CONSISTENCY-RECORD.md` with zero open
  inconsistencies found.

## v0.04 — Prove the platform on the technical project log

- Brought `templates/site/log` up to parity with `templates/site/profile`'s
  React Router v7 / Worker / MCP / negotiation buildout (it had remained
  the thin v0.01-era placeholder through v0.02–v0.03).
- Added an entry-list view (chronological, with date/tags/project) and a
  detail view for the `log` site type, plus date-range/tag/project filter
  controls with intersection semantics over the bundled
  `content-index.json`.
- Extended `internal/content/schema.go`'s frontmatter validation with
  structural checks for `log`-type `date` (parseable YYYY-MM-DD), `tags`
  (list of strings if present), and `project` (string if present) —
  `tags`/`project` remain optional.
- Added log-type-specific MCP tools (`list_entries`, `get_entry`),
  registered alongside the shared tools; a profile-shaped query against a
  log site returns the MCP SDK's standard "unknown tool" error.
- Proved the daily-authoring-rhythm claim with an integration test
  authoring three consecutive dated entries one at a time, confirming
  `wrangler.jsonc`, `content/site.yaml`, and `package.json` never change
  between cycles.

## v0.03 — Make the site usable by AI agents and terminal clients

- Added an in-Worker MCP server (mounted on `/mcp*`, gated entirely by
  `mcp_enabled`) exposing shared (`list_pages`, `get_page`,
  `search_content`) and profile-type-specific (`get_profile`,
  `list_projects`) tools, built on Cloudflare's Agents SDK (`McpAgent`).
- Added `biblio-cli new --mcp-enabled` (default `true`) as the sole
  configuration surface deciding whether `wrangler.jsonc` mounts the
  Worker `main`, `durable_objects` binding, and `nodejs_compat` flag the
  MCP server needs — a build-time decision with no "disabled but
  present" state.
- Added a content-negotiation layer: a page request signaling a
  Markdown/terminal preference (`Accept: text/markdown`, or a recognized
  terminal `User-Agent`) gets clean Markdown or ANSI-rendered output
  (via `md4x`) instead of full browser HTML, with `Vary: Accept,
  User-Agent` set for correct edge caching.
- Added a small, sandboxed MCP App UI bundle (a themed chart view) for
  `search_content` results, always paired with — never replacing — the
  tool's plain-text answer.
- Added build-time-generated `/llms.txt` and a per-page Markdown-discovery
  hint, for agent/crawler discoverability.

## v0.02 — Build, deploy, and verify a real site

- Added `biblio-cli build` to pre-render every content-driven page into
  static output — a thin wrapper over `vite build`, driving React Router
  v7's native `prerender` against `content-index.json`.
- Added `biblio-cli deploy` to publish a built site to its configured
  public address — a thin wrapper over `wrangler deploy`; re-running it
  against an already-live site updates cleanly in place.
- Added React Router v7 framework-mode wiring (`react-router.config.ts`,
  `app/root.tsx`, list/detail route modules) to the `profile` site template,
  reading `content-index.json` directly rather than re-parsing authored
  Markdown on the TypeScript side.
- Added a generated GitHub Actions workflow
  (`.github/workflows/deploy.yml`) invoking `biblio-cli build`/`deploy` on
  push, keeping the whole lifecycle behind the `biblio-cli` surface.

## v0.01 — Scaffold and theme a site

- Added `biblio-cli new <project-name>` to scaffold a complete, runnable,
  themed site project (`profile`, `log`, or `docs`) from one guided command,
  generating `content/`, `content/site.yaml`, `wrangler.jsonc`, and
  `tokens.css`.
- Added `biblio-cli index [path]` to generate the canonical
  `content-index.json` from an authored `content/` directory, validating
  required frontmatter per site type and reporting every failing file on
  error.
- Added `biblio-cli dev` to run `vite dev` and `wrangler dev` together as
  one supervised local-preview process pair.
- Added one ready-to-use theme preset per site type: `profile-default`,
  `log-default`, `docs-default`.
- Added a Bubble Tea guided flow so `new` can be run with no flags at all,
  prompting interactively for site type, domain, and theme.
- Added distribution via prebuilt GitHub Releases binaries and
  `go install github.com/monamaret/biblio/cmd/biblio-cli@latest`.
