# Implementation Plan: v0.05 — Prove the platform on documentation and across all site types

**Branch**: `105-release-0-05-docs-site-and-consistency` | **Date**: 2026-06-07 | **Spec**: [spec.md](./spec.md)

**Input**: Release specification from `/specs/releases/0.05-docs-site-and-consistency/spec.md` (mirrors [`specs/releases/0.05-release.md`](../0.05-release.md)), depends on [v0.04](../0.04-project-log-site-type/plan.md), built on feature item [009](../../features/009-replicate-site-types.md)

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Replicate the platform onto its third and final site type — a project-documentation site organized by table-of-contents/section/order hierarchy, deliberately *without* the tag/project filtering that defines the project-log type — and then run the cross-type consistency check that's only meaningful now that all three site types exist side by side: confirm theming presets and the terminal/CLI view both render correctly and consistently across profile, project log, and documentation. Like v0.04, this is a generalization test (assume the platform works unchanged, fix any gap as a defect); unlike v0.04, it closes with a comparative review across all three types rather than a single site's evidence alone.

## Technical Context

**Language/Version**: Go 1.22+ (`biblio-cli` — `new --type=docs` path exercised, no new commands) and TypeScript/React (documentation-specific scaffold template, TOC/hierarchy navigation, MCP tools, terminal rendering)

**Primary Dependencies**: Same confirmed stack as v0.01–v0.04 (Cobra/Viper/Bubble Tea; React + Vite + React Router v7 + Wrangler; `gray-matter`/`react-markdown`/`remark-gfm`; in-Worker MCP server + `md4x`) — this release adds a hierarchical documentation content schema, TOC/navigation UI, and corresponding MCP tools/terminal renderings; the cross-type review itself needs no new dependency, only a structured comparison of the three already-deployed sites

**Storage**: Unchanged — a third `content-index.json` now carries section/order-hierarchical pages; no filtering query surface is added for this type by design (Principle V — no escalation trigger; the "no tag/project filtering" absence is itself a confirmed, deliberate design decision per [roadmap §0.1](../../research/roadmap.md), not a gap)

**Testing**: Go unit tests for the documentation content schema/hierarchy/slug derivation, TypeScript/React tests for TOC/navigation rendering and the *absence* of filter UI/query params (spec Scenario 1.2 — proving a negative is part of this release's contract), MCP tool contract tests for the documentation-specific query surface, and two end-to-end checks: (a) scaffold→author→deploy→navigate for the docs site, and (b) the cross-type side-by-side theming + terminal-rendering review across all three deployed sites — written first per Principle VI

**Target Platform**: A third, independently deployed instance of the same single-Worker topology (Workers Static Assets + `/mcp*` gated by `mcp_enabled`) — plus, for the consistency check, the two already-live sites from v0.02/v0.04 reviewed alongside it

**Project Type**: Web application replication (third site type) plus a cross-cutting verification pass — no new service, host, or project; the comparison runs against three already-deployed Worker instances

**Performance Goals**: Every authored page is reachable and renders correctly through the TOC/section/order navigation alone (FR-002, SC-001); zero cross-type theming or terminal-rendering inconsistencies found in the side-by-side review (FR-003, FR-004, SC-002)

**Constraints**: MUST NOT present any tag- or project-based filtering for this site type — its absence is a confirmed, deliberate design decision being *verified*, not a feature being added (FR-002, spec Assumptions); a tag-filter request against the docs site (a parameter shape that only makes sense for the project-log type) MUST be handled with a documented, sensible default — rejected clearly or cleanly ignored, not silently mishandled (spec Edge Cases); any cross-type theming inconsistency found (e.g. a shared token resolving differently per type) MUST be treated as a defect against the *shared* theming system (§6/v0.01), never patched per-site (spec Edge Cases)

**Scale/Scope**: One real, live documentation site, plus all three site types now live and reviewed together — the first point at which the true cross-type consistency check is even possible (spec Summary); update/upgrade-flow work is deliberately deferred to v0.06 (spec Assumptions, SC-003)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Single-Author, Minimal-Dependency Bias)**: PASS — no new dependencies or config knobs; the deliberate *absence* of filter UI/query params for this type is itself the minimal-surface decision being verified, not something to backfill. The cross-type review reuses existing sites and tooling — it adds no new infrastructure of its own.
- **Principle II (One Pipeline, One Source of Truth)**: PASS — the documentation content schema, hierarchy/slug derivation, and TOC generation extend the same canonical Go-side `internal/content`/`index`/`theme` packages; any cross-type theming-token inconsistency the review finds is explicitly routed back to the *single* shared theming system as a defect (spec Edge Cases) — never forked into a per-type override, which is exactly the divergent-pipeline outcome Principle II forbids.
- **Principle III (CLI-First, Guided Experience)**: PASS — scaffolding/authoring/deploying the docs site uses the identical CLI surface already proven for the other two types (FR-001); no documentation-specific commands.
- **Principle IV (Text-First, Progressive Enhancement)**: PASS — extends the same text-first MCP/terminal contract from v0.03 (proven generalizing in v0.04) to the documentation content shape (FR-004); the cross-type terminal-rendering review explicitly checks that each site returns a rendering "appropriate to its own content shape," not a forced one-size-fits-all template.
- **Principle V (Simplest Topology That Can Grow)**: PASS — third site is a third independent instance of the same simplest-topology deployment; introducing *no* filtering surface for this type (by design) is itself consistent with "escalate only on a documented real trigger" — there is no trigger here, and the spec confirms that absence is a feature, not a gap.
- **Principle VI (Test-First Development, NON-NEGOTIABLE)**: GATE — failing tests for the documentation schema/TOC/tools/rendering, *and* for the absence of filter UI/params (proving a negative is itself testable and required here), MUST be authored and reviewed before implementation; the cross-type consistency review is itself run as a structured, repeatable check (not an ad hoc visual scan) so its findings route back into the shared theming/terminal-view test suites as regressions, per the constitution's Release Process integration-test requirement.
- **Security & Operational Posture**: PASS — `mcp_enabled` remains the sole MCP config surface for the third site too; no new auth/access-control surface introduced anywhere in this release.

No violations requiring justification — Complexity Tracking is not needed for this release.

## Project Structure

### Documentation (this feature)

```text
specs/releases/0.05-docs-site-and-consistency/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
# extensions to the existing biblio-cli content/scaffold packages (Go) — no new commands
internal/
├── content/
│   └── docs/                # documentation frontmatter schema (section, order) + hierarchy/slug rules
└── scaffold/
    └── docs/                # documentation seed-content templates and TOC structure

# documentation-specific additions to the combined site/Worker template (TypeScript/React)
templates/site/docs/
├── app/
│   ├── routes/
│   │   └── pages/           # hierarchical page views, navigated via TOC/section/order
│   │                        #  (deliberately: no filter routes/components — Scenario 1.2)
│   └── mcp/
│       └── tools/           # documentation-specific structured-query tools
│           #                  (rejects/ignores log-shaped tag-filter params per a documented default)
└── app/negotiation/         # reuses v0.03 negotiation/rendering against docs content shape

test/
├── unit/
│   ├── content-docs/        # schema, hierarchy/slug, TOC-generation unit tests
│   ├── no-filter-surface/   # asserts absence of filter UI/query params for this type (Scenario 1.2)
│   └── mcp-tools-docs/      # documentation tool content-contract tests
└── integration/
    ├── pilot-site-docs/     # scaffold → author hierarchical pages → deploy → navigate (FR-001..FR-002)
    └── cross-type-review/   # structured side-by-side theming + terminal-rendering check
                             #  across the profile (v0.02), log (v0.04), and docs (this release) sites
```

**Structure Decision**: Purely additive, type-scoped extensions following the exact same `internal/content|scaffold`, `templates/site/<type>`, and MCP/negotiation-reuse pattern established in v0.04 — plus one new cross-cutting integration suite (`cross-type-review/`) that is the actual *closing* artifact of this release's consistency-check half. No new top-level project, command, host, or service; any inconsistency the review finds is fixed in the shared `internal/theme`/negotiation code, never forked per type.

## Complexity Tracking

*No Constitution Check violations — section intentionally left empty.*
