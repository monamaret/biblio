# Phase 1 Data Model: v0.05 — Prove the platform on documentation and across all site types

**Input**: [spec.md](./spec.md) FR-001..FR-004, [research.md](./research.md) replication/consistency decisions

This release introduces one new authoring-time schema (the documentation page, hierarchical rather than timeline-shaped) plus one new cross-cutting review artifact (the consistency record). Both slot into the existing Content Item / Content Index / Theme Preset / Deployment entities from [v0.01](../0.01-scaffold-and-theme/data-model.md)–[v0.04](../0.04-project-log-site-type/data-model.md) — same canonical owner, same index shape, same MCP/negotiation machinery (Principle II).

## Entity: Documentation Page (extends Authored Content Item)

A single hierarchically-organized Markdown+frontmatter file under `content/`, validated by an extension of the canonical Go-side schema (`internal/content/docs`).

| Field | Type | Notes |
|---|---|---|
| `section` | string | Groups pages into the TOC's top-level structure |
| `order` | integer | Determines a page's position within its section for navigation (FR-002) |
| `slug` / `path` | derived | Same canonical derivation as every Authored Content Item — no docs-specific slug scheme |
| `body` | Markdown | Authored documentation prose |

**Validation rules**: same fail-clearly-not-silently rule as every site type (malformed/missing `section`/`order` surfaced as a clear, per-file error at indexing time); a page that doesn't cleanly map to a single `section`/`order` slot (e.g. logically belongs in two sections) MUST resolve via a documented, sensible default rather than an undefined or duplicated placement (spec Edge Cases) — this default is itself a unit under test-first development before the TOC generator is implemented.

**Relationships**: one Documentation Page = one Authored Content Item, indexed identically into the shared Content Index — there is no separate "docs index," and explicitly **no** filter-relevant fields (`tags[]`/`project`) are part of this schema, by deliberate design (research.md "no tag/project filtering").

## Entity: Table of Contents / Navigation Structure

The hierarchy-based browsing structure derived from `section`/`order` across all Documentation Pages (FR-002, Scenario 1).

| Field | Type | Notes |
|---|---|---|
| `sections[]` | ordered list of sections, each with ordered pages | Derived deterministically from each page's `section`/`order` fields at index/build time |
| filter surface | **absent by design** | No tag/project filter UI or query params exist for this entity — its absence is the thing Scenario 1.2 explicitly verifies, not an omission to fill in |

**Validation rules**: every authored page MUST be reachable and render correctly through this structure alone (SC-001); the structure MUST NOT expose, anywhere in its UI or query surface, the tag/project filtering that exists on the project-log type (Scenario 1.2) — a test asserting this *absence* is as load-bearing as any test asserting a presence (Principle VI).

**Relationships**: derived 1:many from the Documentation Pages in the shared Content Index; is the sole browsing/query surface for this site type — there is no parallel Filter Query entity for documentation (contrast with [v0.04's Filter Query](../0.04-project-log-site-type/data-model.md), which deliberately has no analog here).

## Entity: Cross-Type Consistency Record

The closing review artifact for this release — not a runtime entity, but the structured comparison SC-002 is checked against, produced once all three site types are live.

| Field | Type | Notes |
|---|---|---|
| `sites[]` | exactly the three live deployments | Profile (from v0.02), project log (from v0.04), documentation (this release) |
| `theming-check` | per-site result | Each site displays its own distinct, correctly-applied theme with no cross-type bleed-through (FR-003, Scenario 2) |
| `terminal-rendering-check` | per-site result | Each site returns a clean, correct, type-appropriate terminal/Markdown view on request (FR-004, Scenario 3) |
| `inconsistencies-found[]` | list, each routed to a defect | Any finding (e.g. a shared token resolving differently per type) is filed as a defect against the *shared* theming system — never patched per-site (research.md "How cross-type inconsistencies are triaged") |

**Validation rules**: SC-002 requires this record to show all three types simultaneously live, each visually distinct per its own theme, each returning a correct type-appropriate terminal rendering, with **zero** cross-type inconsistencies remaining open at review close (an inconsistency found-and-then-fixed-as-a-shared-defect is acceptable; one found-and-left-open is not).

**Relationships**: references exactly the three Deployment entities (carried over unchanged in shape from [v0.02's data-model](../0.02-build-deploy-verify/data-model.md)); any `inconsistencies-found[]` entry creates a tracked defect against `internal/theme` or the shared negotiation/rendering code — never a new per-site code path.
