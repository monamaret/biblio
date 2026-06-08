# Implementation Plan: v0.01 — Scaffold and theme a site

**Branch**: `101-release-0-01-scaffold-and-theme` | **Date**: 2026-06-07 | **Spec**: [spec.md](./spec.md)

**Input**: Release specification from `/specs/releases/0.01-scaffold-and-theme/spec.md` (mirrors [`specs/releases/0.01-release.md`](../0.01-release.md)), built on feature items [001](../../features/001-biblio-cli-foundations.md), [002](../../features/002-content-pipeline.md), [006](../../features/006-theming-system.md)

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Ship `biblio-cli` as a single Go binary that scaffolds a complete, runnable, themed combined site project (content + Wrangler/Vite plumbing) for any of the three site types from one guided command, generates a single canonical `content-index.json` from authored content, and wraps local preview (`vite dev` + `wrangler dev`) behind one `biblio-cli dev` command — the foundation every later release builds on. Technical approach: build the CLI on Cobra/Viper/Bubble Tea per the confirmed stack, with the canonical content-parsing/schema/slug/theming logic owned by the Go side (Principle II) and consumed unchanged by the JS/TS site in later releases.

## Technical Context

**Language/Version**: Go 1.22+ (biblio-cli); TypeScript/Node (generated site project template — React + Vite + React Router v7, consumed by `vite dev`/`vite build` but not authored as part of this release's CLI work)

**Primary Dependencies**: Cobra + Viper (CLI/config), Bubble Tea + Lip Gloss + Glamour (interactive TUI), `goreleaser` (binary distribution); generated-template side: React + Vite + React Router v7, `gray-matter` (frontmatter parsing reference for the template, though canonical parsing logic lives in Go per Principle II), Wrangler/`wrangler.jsonc`

**Storage**: Generated, in-memory-then-bundled `content-index.json` (flat JSON file emitted by the CLI's indexing step) — no database; D1 escalation explicitly out of scope (Principle V)

**Testing**: Go unit/CLI tests (`go test`) for content-index generation, frontmatter schema parsing/validation, slug derivation, and scaffolding/templating logic — written first per Principle VI (Test-First Development, NON-NEGOTIABLE); red-green-refactor for every command and generator added in this release

**Target Platform**: Cross-platform CLI binary (macOS/Linux/Windows via `goreleaser`); generated site projects target Cloudflare Workers (Workers Static Assets) via Wrangler, run locally via `vite dev` + `wrangler dev`

**Project Type**: CLI tool (single Go module) that generates combined web-application projects (site + Worker) as its output artifact

**Performance Goals**: Scaffold-to-runnable-preview in a single guided session with no manual file editing (SC-001); content-indexing completes correctly over an author's full `content/` tree on first run (no perf targets beyond "single-author scale" per Principle I)

**Constraints**: Single-binary distribution only (no package-manager taps, Principle III); canonical parsing/schema/slug/theming logic MUST live in `biblio-cli`, not be duplicated in the generated TS/React template (Principle II); no speculative config knobs beyond the documented `mcp_enabled` boolean and the three site-type/theme presets (Principle I)

**Scale/Scope**: Single author, three site-type scaffolds (profile/CV, project log, documentation), one theme preset per site type — no multi-user, multi-tenant, or productization surface (Principle I)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Single-Author, Minimal-Dependency Bias)**: PASS — scope is exactly the three named site types and one theme preset each; no speculative site types, auth, or config surfaces are introduced. `mcp_enabled` is generated as the sole MCP-related config knob (per Security & Operational Posture), even though the MCP server itself ships in v0.03.
- **Principle II (One Pipeline, One Source of Truth)**: PASS — `content-index.json` generation and `tokens.css` generation are both owned by `biblio-cli` (Go side) as confirmed; the generated React/Vite template consumes their output rather than re-implementing parsing/schema/theming logic. This plan explicitly calls out that frontmatter schema, slug derivation, and theme-token generation logic must not be duplicated on the TS side.
- **Principle III (CLI-First, Guided Experience)**: PASS — every capability (`new`, indexing, `dev`) is a `biblio-cli` command/flag with the Bubble Tea TUI as guided fallback; distribution is single-binary via GitHub Releases (`goreleaser`) + `go install`, matching the confirmed §0.3 resolution.
- **Principle IV (Text-First, Progressive Enhancement)**: N/A for this release — no MCP tools or terminal/CLI rendering surface ships until v0.03; nothing in this release's scope produces chart-like or richer-than-text output.
- **Principle V (Simplest Topology That Can Grow)**: PASS — content index is a generated flat JSON file (no D1); no search capability is introduced yet (keyword search ships with the MCP layer in v0.03, still pre-Vectorize).
- **Principle VI (Test-First Development, NON-NEGOTIABLE)**: GATE — every generator/command in this release (scaffolding templater, frontmatter schema parser/validator, slug deriver, content-indexer, `dev` wrapper) MUST have a failing Go test authored and reviewed before its implementation lands. Tasks generated from this plan must sequence test-authoring before implementation for each unit.
- **Security & Operational Posture**: PASS — scaffolding generates the `mcp_enabled` boolean (default `true`) in `content/site.yaml` as the only MCP-related config surface; no auth/access-control surface is generated (the MCP route itself doesn't exist until v0.03, so there is nothing yet to secure beyond emitting the correct config flag).

No violations requiring justification — Complexity Tracking is not needed for this release.

## Project Structure

### Documentation (this feature)

```text
specs/releases/0.01-scaffold-and-theme/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
# biblio-cli — single Go module, the only thing actually built in this release
cmd/
└── biblio-cli/
    └── main.go

internal/
├── cli/                 # Cobra command tree: new, index, dev (+ root/version)
├── tui/                 # Bubble Tea models + Lip Gloss styles for the guided `new` flow
├── scaffold/            # Site-type templates (profile|log|docs) + project-file generation
├── content/             # Frontmatter schema parsing/validation, slug derivation (canonical, Principle II)
├── index/               # content-index.json generation
├── theme/               # Theme-preset definitions + tokens.css generation (feeds §6/v0.01 jointly)
└── devserver/           # `biblio-cli dev` wrapper around `vite dev` + `wrangler dev`

templates/
└── site/                # Embedded combined React+Vite+Worker project template(s) per site type
    ├── profile/
    ├── log/
    └── docs/

test/
├── unit/                # Per-package Go unit tests (content, index, theme, scaffold)
└── integration/         # End-to-end "scaffold → index → dev" CLI tests per site type

.goreleaser.yaml         # Binary build/release config (GitHub Releases distribution)
```

**Structure Decision**: Single Go module/CLI project (`Option 1` shape, CLI-flavored). All canonical logic (parsing, schema, slugs, theming, scaffolding) lives under `internal/` in Go per Principle II; the generated site projects under `templates/site/` are inert templates populated at scaffold time, not a second implementation of that logic. No `frontend`/`backend` split exists at this layer — the combined React+Vite+Worker shape is the *output* of scaffolding, not this release's own source layout.

## Complexity Tracking

*No Constitution Check violations — section intentionally left empty.*
