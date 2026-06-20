# Phase 1 Data Model: v0.01 — Scaffold and theme a site

**Input**: [spec.md](./spec.md) Functional Requirements FR-001..FR-006, [research.md](./research.md) decisions on content-index/theming ownership

This release's "data" is entirely file-shaped — generated project files, generated config, and one generated JSON index — owned canonically by `biblio-cli` (Go side, Principle II). There is no database and no runtime entity model; the entities below are the *generated artifacts* and the *authoring-time schema* `biblio-cli` parses and validates.

## Entity: Site Project (scaffold output)

The combined, runnable project directory produced by `biblio-cli new`.

| Field | Type | Notes |
|---|---|---|
| `project-name` | string (dir name) | From the `new` command's positional argument |
| `site-type` | enum: `profile` \| `log` \| `docs` | From `--type`; selects scaffold template, theme preset, and seed content shape |
| `domain` | string | From `--domain`; pre-fills the route/custom-domain binding in `wrangler.jsonc` |
| `theme` | preset identifier | From `--theme`; selects which theme preset feeds `tokens.css` generation |
| `content/` | directory tree | Seeded Markdown + YAML frontmatter content matching `site-type` |
| `content/site.yaml` | config file | Declares site type, name/description, and `mcp_enabled` (boolean, default `true` — the sole MCP config surface per Security & Operational Posture) |
| `wrangler.jsonc` | config file | Deployment-routing config; route/custom-domain binding pre-filled from `domain` |
| `tokens.css` | generated CSS file | Output of theme-preset + site-type selection (Decision: theming ownership) |

**Validation rules**: `site-type` MUST be one of the three supported values; `domain` and `project-name` MUST produce a valid, non-colliding project directory and route binding; `theme` MUST resolve to one of the presets defined for the chosen `site-type` (FR-004 — one distinct preset per site type).

**Relationships**: One Site Project has exactly one `content/site.yaml`, one `wrangler.jsonc`, one `tokens.css`, and one authored `content/` tree, which in turn produces exactly one Content Index (below).

## Entity: Authored Content Item (frontmatter schema)

A single Markdown + YAML-frontmatter file under `content/`, as parsed and validated by the canonical Go-side schema (Principle II — this schema MUST NOT be re-implemented on the TS side).

| Field | Type | Notes |
|---|---|---|
| `path` | derived (filesystem path under `content/`) | Source of slug derivation |
| `slug` | derived string | Canonical, structurally-derived from `path` per the confirmed layout convention; MUST be stable and collision-free within a site |
| structured frontmatter fields | per `site-type` schema | e.g. profile entries vs. dated/tagged log entries vs. hierarchical doc pages — exact per-type field sets are defined by the seed-content templates this release ships, and re-validated, not re-designed, in v0.04/v0.05 when the other site types are replicated |
| `body` | Markdown | Authored prose body |

**Validation rules**: malformed or missing structured fields MUST be surfaced as a clear, actionable error at indexing time — never silently produce an incomplete index entry (spec Edge Cases). This is the behavior the failing-test-first units in `internal/content` must encode before the parser/validator is implemented (Principle VI).

**State/lifecycle**: authored → parsed → validated → indexed. There is no further runtime state in this release (no editing UI, no publish workflow beyond re-running the indexing step).

## Entity: Content Index (`content-index.json`)

The single canonical index generated from an author's `content/` tree (FR-003).

| Field | Type | Notes |
|---|---|---|
| `entries[]` | array of indexed items | One per Authored Content Item; carries its derived `slug`/`path` and structured frontmatter fields |
| shape contract | — | MUST be directly consumable by both React Router's `prerender` (site build, v0.02) and the in-Worker MCP query layer (v0.03) — designed now so neither later consumer requires a divergent second index (Decision: content-index shape and ownership) |

**Validation rules**: every authored page/entry MUST appear exactly once, with a correctly derived `slug`/`path` (spec Scenario 2); the index generation step is itself a unit under test-first development (Principle VI).

**Relationships**: generated 1:1 from a Site Project's `content/` tree; consumed (read-only) by the site build and, from v0.03 onward, the MCP query layer — `biblio-cli` is its sole producer (Principle II: "pick one owner and have the other consume its output").

## Entity: Theme Preset

A named, ready-to-use visual-theme definition selectable at scaffold time (FR-004).

| Field | Type | Notes |
|---|---|---|
| `preset-id` | string identifier | Selected via `--theme`; one distinct preset per `site-type` |
| `site-type` | enum | The preset's matching site type |
| token definitions | design-token values | Source values from which `tokens.css` is generated |

**Validation rules**: each `site-type` MUST have at least one ready-to-use preset (FR-004); preset → `tokens.css` generation MUST be deterministic and owned solely by `biblio-cli` (Decision: theming ownership — generated once, imported by both site build and future MCP App UI bundle).

**Relationships**: selected by exactly one Site Project at scaffold time; produces exactly one `tokens.css` per Site Project.
