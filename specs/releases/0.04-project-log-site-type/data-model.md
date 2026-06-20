# Phase 1 Data Model: v0.04 — Prove the platform on the technical project log

**Input**: [spec.md](./spec.md) FR-001..FR-005, [research.md](./research.md) replication decisions

This release introduces exactly one new authoring-time schema — the project-log entry — slotting into the *existing* Content Item / Content Index / Theme Preset entities defined in [v0.01's data-model](../0.01-scaffold-and-theme/data-model.md) (Principle II: same canonical owner, same index shape, no parallel pipeline). The entities below describe only what's new or filter-relevant; everything else (Site Project, Content Index shape contract, Deployment, MCP Tool/Server Mount, Negotiated Page Response) carries over unchanged from v0.01–v0.03.

## Entity: Project-Log Entry (extends Authored Content Item)

A single dated, tagged, project-attributed Markdown+frontmatter file under `content/`, validated by an extension of the canonical Go-side schema (Principle II — lives in `internal/content/log`, not duplicated in TS).

| Field | Type | Notes |
|---|---|---|
| `date` | structured date | Drives chronological ordering and date-range filtering (FR-002, FR-003) |
| `tags[]` | array of strings | Zero or more; drives tag filtering; a brand-new tag MUST appear correctly in filters immediately after the next build with no extra registration step (spec Edge Cases) |
| `project` | string (optional) | Drives project filtering; same "appears immediately after next build" guarantee as `tags[]` |
| `slug` / `path` | derived | Same canonical derivation rules as every Authored Content Item ([v0.01 data-model](../0.01-scaffold-and-theme/data-model.md)) — no project-log-specific slug scheme |
| `body` | Markdown | Authored update text |

**Validation rules**: `date`, `tags[]`, and `project` are validated by the same fail-clearly-not-silently rule established in v0.01 (malformed/missing structured fields surfaced as a clear, per-file error at indexing time); every entry MUST appear in the deployed site's list with its date/tags/project correctly attributed (Scenario 1, FR-002).

**Relationships**: one Project-Log Entry = one Authored Content Item, indexed identically into the shared Content Index — there is no separate "log index."

## Entity: Filter Query

A visitor- or agent-issued request to narrow the entry list by one or more dimensions (FR-003, Scenario 2).

| Field | Type | Notes |
|---|---|---|
| `date-range` (optional) | start/end date | |
| `tag` (optional) | string | |
| `project` (optional) | string | |
| combination semantics | intersection | Combining dimensions (e.g. a date range *and* a tag) MUST return the intersection — exactly what a visitor would expect (spec Edge Cases) |

**Validation rules**: a Filter Query with no matching entries returns an empty result set (no false positives), and every genuinely matching entry MUST appear (no omissions) — both directions checked by the failing-test-first filter-computation unit (`internal/index/filters.go`) before the UI/MCP surfaces consume it.

**Relationships**: evaluated against the existing Content Index's `entries[]` (read-only, same index as every other site type); produces a Filtered Entry Set (below). Also exposed as a structured-query shape through the project-log-specific MCP tools (extending [v0.03's MCP Tool entity](../0.03-agent-terminal-layer/data-model.md)).

## Entity: Filtered Entry Set

The visible result of applying a Filter Query to the project log's entries — what a visitor sees in the UI and what an MCP/terminal client receives as a tool result.

| Field | Type | Notes |
|---|---|---|
| `entries[]` | subset of Project-Log Entry | The intersection result for the applied Filter Query |
| `rendering` | per-surface | Browser: filtered entry-list view; MCP/terminal: the same self-sufficient Markdown/plain-text `content` contract from v0.03, extended to this site type's shape (FR-005) |

**Validation rules**: an MCP/terminal request shaped for the *profile* type's query surface (e.g. "get résumé") MUST be rejected clearly by this site type, reusing the v0.03 "unsupported capability" contract rather than a bespoke error shape (research.md "Type-mismatch query rejection").

**Relationships**: derived from exactly one Filter Query applied to the Content Index; rendered through the same Negotiated Page Response / MCP Tool machinery proven in v0.03, parameterized by this site type's content shape.
