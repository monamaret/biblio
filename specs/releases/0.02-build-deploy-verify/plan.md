# Implementation Plan: v0.02 — Build, deploy, and verify a real site

**Branch**: `102-release-0-02-build-deploy-verify` | **Date**: 2026-06-07 | **Spec**: [spec.md](./spec.md)

**Input**: Release specification from `/specs/releases/0.02-build-deploy-verify/spec.md` (mirrors [`specs/releases/0.02-release.md`](../0.02-release.md)), depends on [v0.01](../0.01-scaffold-and-theme/plan.md), built on feature items [003](../../features/003-site-scaffold.md), [008](../../features/008-pilot-site-validation.md)

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Carry the v0.01 scaffold through to a real, live deployment: pre-render every content-driven page of a scaffolded profile/CV site via React Router v7's native `prerender`, ship `biblio-cli build`/`biblio-cli deploy` as thin wrappers over `vite build`/`wrangler deploy` plus a generated CI workflow, and prove the whole content → build → deploy → live-verification path on one real, fully authored pilot site reachable at a public address. Technical approach: extend the React + Vite + React Router v7 + Wrangler combined-project template (scaffolded in v0.01) with prerender wiring that consumes the unchanged `content-index.json`/`tokens.css` outputs from v0.01, and add the build/deploy command wrappers to `biblio-cli` per the confirmed §0.3 command-surface decision.

## Technical Context

**Language/Version**: TypeScript/React (site project: React 18+, Vite, React Router v7 with native `prerender`); Go 1.22+ (`biblio-cli` — adding `build`/`deploy` command wrappers to the existing Cobra tree)

**Primary Dependencies**: React + Vite + React Router v7 (`prerender` configured in `react-router.config.ts`), `import.meta.glob` + `gray-matter` + `react-markdown`/`remark-gfm` for content loading/rendering, Wrangler (`wrangler deploy`, generated CI workflow) for publishing to Cloudflare Workers (Workers Static Assets)

**Storage**: Unchanged from v0.01 — consumes the generated `content-index.json` as a stable, read-only input to `prerender`; no new storage introduced (Principle V — still pre-D1)

**Testing**: TypeScript/React tests for the site's pre-rendered views (list/detail/navigation per site type) plus Go CLI tests for the new `build`/`deploy` command wrappers — written first per Principle VI; an integration check that the full content → build → deploy path completes and every authored page renders with no missing/broken pages (spec Scenario 1, FR-001) is the release-closing test, in keeping with the constitution's Release Process requirement for both unit and end-to-end integration tests

**Target Platform**: Cloudflare Workers + Workers Static Assets, served at the site's configured public/custom-domain address (route binding generated in v0.01's `wrangler.jsonc`)

**Project Type**: Web application (the generated combined React+Vite+Worker project) plus its CLI driver (`biblio-cli build`/`deploy` wrappers) — this release adds to both halves of the v0.01 scaffold output and the CLI that drives it

**Performance Goals**: Build produces complete static output with zero missing/broken pages on the first run (FR-001, Scenario 1); deploy completes to a publicly reachable address from a single command with no manual intervention (SC-002)

**Constraints**: `build`/`deploy` MUST remain thin wrappers over `vite build`/`wrangler deploy` + generated CI — no parallel hand-rolled build/deploy logic (Principle III); this release MUST NOT redesign the v0.01 scaffold, content index, or theme outputs, only prove they carry through to a live deployment unchanged (spec Assumptions); a malformed-content build failure MUST surface a clear message rather than silently skip pages (spec Edge Cases)

**Scale/Scope**: Single pilot site, single site type (personal-profile/CV) — proving the other two site types live is explicitly deferred to v0.04/v0.05 (spec Assumptions); no multi-site or multi-tenant deploy orchestration (Principle I)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Single-Author, Minimal-Dependency Bias)**: PASS — exactly one pilot site, one site type, proven live; no speculative multi-site deploy tooling, environments, or staging infrastructure introduced.
- **Principle II (One Pipeline, One Source of Truth)**: PASS — `prerender` consumes the unchanged `content-index.json` from v0.01 as its single source of truth; no second content-parsing or theming path is introduced on the build/deploy side. The plan explicitly treats v0.01's outputs as stable, unmodified inputs (spec Assumptions).
- **Principle III (CLI-First, Guided Experience)**: PASS — `biblio-cli build`/`biblio-cli deploy` are thin Cobra-tree wrappers over `vite build`/`wrangler deploy` + generated CI, matching the confirmed §0.3 resolution that folded build/deploy (and later update/upgrade) into one consistent CLI surface — never raw `package.json` scripts or hand-assembled CI steps run directly by the author.
- **Principle IV (Text-First, Progressive Enhancement)**: N/A — no MCP tools or terminal/CLI rendering surface ships in this release (that's v0.03); the deployed site is browser-facing only at this point, exactly as spec SC-003 states.
- **Principle V (Simplest Topology That Can Grow)**: PASS — one Worker serves the prerendered static site via Workers Static Assets on its configured route; no split across hosts/services, no database, no search — nothing here crosses an escalation trigger.
- **Principle VI (Test-First Development, NON-NEGOTIABLE)**: GATE — failing tests for prerendered view rendering (per site type's list/detail/nav views) and for the `build`/`deploy` CLI wrappers MUST be authored and reviewed before their implementations land; the end-to-end content→build→deploy integration check is the release-closing test required by the constitution's Release Process.
- **Security & Operational Posture**: PASS — nothing in this release touches the `/mcp` surface (it doesn't exist yet); the deployed site relies on Cloudflare's standard zone-level protections with no new auth/access-control layer introduced.

No violations requiring justification — Complexity Tracking is not needed for this release.

## Project Structure

### Documentation (this feature)

```text
specs/releases/0.02-build-deploy-verify/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
# biblio-cli additions (Go) — extends the v0.01 Cobra tree, no new module
internal/
└── cli/
    ├── build.go         # `biblio-cli build` — thin wrapper over `vite build` (prerender)
    └── deploy.go        # `biblio-cli deploy` — thin wrapper over `wrangler deploy` + generated CI

# generated combined-project template additions (TS/React) — extends v0.01's templates/site/*
templates/site/<type>/
├── react-router.config.ts   # `prerender` wired against content-index.json
├── app/
│   ├── routes/              # list/detail/navigation route modules per site type
│   └── content/             # `import.meta.glob` + gray-matter loading, react-markdown rendering
└── .github/workflows/       # generated CI workflow invoked by `biblio-cli deploy`

test/
├── unit/                    # prerendered-view + CLI-wrapper unit tests
└── integration/
    └── pilot-site/          # content → build → deploy → live-verification end-to-end check
                             # against the one real, fully authored pilot site (FR-004, SC-001)
```

**Structure Decision**: This release extends two existing trees rather than introducing new top-level structure: the `biblio-cli` Cobra tree gains thin `build`/`deploy` wrappers (Go), and the v0.01-scaffolded combined React+Vite+Worker template gains its `prerender` wiring and rendered views (TypeScript/React). No `frontend`/`backend` split is introduced — the combined-Worker topology from Principle V is preserved untouched, and the pilot site itself is the integration-test fixture, not a new service.

## Complexity Tracking

*No Constitution Check violations — section intentionally left empty.*
