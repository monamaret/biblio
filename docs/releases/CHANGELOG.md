# Changelog

All notable user-facing changes to `biblio-cli` and its generated sites are
recorded here, release by release. See `docs/releases/v<version>/` for the
full user guide of each release.

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
