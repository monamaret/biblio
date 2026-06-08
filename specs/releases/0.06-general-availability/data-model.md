# Phase 1 Data Model: v0.06 — General Availability

**Input**: [spec.md](./spec.md) FR-001..FR-006, [research.md](./research.md) update/upgrade and GA-definition decisions

This release introduces no new authored-content schema (no fourth site type, no new frontmatter shape) — its "data" is entirely about *classifying and migrating existing scaffold output* and *checking real sites against documented thresholds*. The entities below are new to this release; every Site Project / Content Index / Deployment entity from [v0.01](../0.01-scaffold-and-theme/data-model.md)–[v0.05](../0.05-docs-site-and-consistency/data-model.md) carries over unchanged as the *input* to this release's flow.

## Entity: Scaffold Region Classification

The per-file/per-region judgment the update flow makes before touching anything (FR-001, FR-002).

| Field | Type | Notes |
|---|---|---|
| `region` | enum: `scaffold-tracked` \| `author-customized` \| `ambiguous` | Derived from the single canonical scaffold-template definitions `biblio-cli` already owns (Principle II) — never a second, parallel notion of "what the scaffold looks like" |
| `path` | file/region reference | Within the deployed site's project tree |
| classification source | reference | Diffed against the specific newer-scaffold-version the update targets |

**Validation rules**: every region of an existing site MUST be classified before any write occurs; `ambiguous` is a first-class outcome (not an error) and routes directly to the Update Conflict entity below — it is never silently resolved in either direction (FR-002, spec Edge Cases "does the flow surface the conflict for a human decision, rather than choosing a side automatically").

**Relationships**: computed 1:many per Update Run (below) against one Site Project; feeds the Update Conflict entity for every `ambiguous` classification.

## Entity: Update Run

A single execution of `biblio-cli update`/`upgrade` against one deployed site (FR-001, Scenario 1).

| Field | Type | Notes |
|---|---|---|
| `target-site` | reference | One of the three real deployed sites (profile / log / docs) |
| `target-scaffold-version` | version identifier | The newer scaffold version being moved to |
| `preserved[]` | list of regions | Author content + intentional customizations carried through untouched |
| `advanced[]` | list of regions | Scaffold-tracked parts brought up to date |
| `conflicts[]` | list of Update Conflict | Surfaced, not auto-resolved |
| `interruption-state` | invariant | If interrupted mid-run (e.g. author cancels), the site MUST remain in a consistent, recoverable state — never half-migrated with no path forward (spec Edge Cases) |

**Validation rules**: an Update Run completing successfully MUST leave `preserved[]` exactly intact, `advanced[]` correctly brought forward, and `conflicts[]` fully and clearly enumerated for the author — this triple is what FR-001/Scenario 1 defines as success, and what the failing-test-first `internal/migrate` units must encode before implementation.

**Relationships**: one Update Run targets exactly one Site Project at a time; the cross-type integration suite (`test/integration/update-three-sites/`) runs one Update Run per real site and requires all three to succeed (FR-003, Scenario 2) — this is the entity SC-001 ("carry forward all three real, live sites... in one sitting") is checked against.

## Entity: Update Conflict

A surfaced case where advancing the scaffold would otherwise overwrite an author's intentional customization (FR-002, spec Edge Cases).

| Field | Type | Notes |
|---|---|---|
| `region` | reference | The `ambiguous`-classified region in question |
| `scaffold-proposed` | content/diff | What the newer scaffold version would write |
| `author-current` | content/diff | What the author's customization currently contains |
| `resolution` | author decision (required) | The flow MUST wait for and apply the author's explicit choice — it is never auto-resolved (FR-002) |

**Validation rules**: an Update Run with any non-empty `conflicts[]` MUST NOT proceed to overwrite those regions without the corresponding `resolution` being recorded — this is the core trust guarantee the whole release exists to prove (research.md "Conflict handling philosophy").

**Relationships**: produced 1:many by an Update Run; always traces back to exactly one Scaffold Region Classification of `ambiguous`.

## Entity: Storage-Escalation Guidance Check

The documented-threshold check run against each real site's actual current scale (FR-004, FR-005, Scenario 3).

| Field | Type | Notes |
|---|---|---|
| `site` | reference | One of the three real deployed sites |
| `current-scale` | metrics derived from `content-index.json` | Real volume/relational-query and search-quality indicators — not synthetic projections |
| `threshold-status` | per conditional path: `not-crossed` \| `crossed` | Reported correctly against the documented thresholds for D1 (volume/relational) and Vectorize (search quality) escalation (Principle V's named triggers) |
| `path-executable-on-demand` | boolean + concrete steps | For each path the guidance can recommend, its steps MUST be concrete enough to actually carry out — proven, not merely described in principle (FR-005) |

**Validation rules**: the guidance MUST correctly reflect each site's real, current scale (FR-004) — a guide that reports a status inconsistent with the real `content-index.json`/usage data is itself a defect; "proven executable on demand" requires an actual dry-run/walkthrough of each path's steps against (a copy of) real site data, not a paper description (research.md "What GA means").

**Relationships**: one Storage-Escalation Guidance Check per real site; its aggregate result (across all three) is one of the two pillars of the GA Readiness Review below — alongside the Update Run results.

## Entity: GA Readiness Review

The closing, cross-release artifact that marks the platform Generally Available (FR-006, Scenario 4, SC-002).

| Field | Type | Notes |
|---|---|---|
| `deliverables[]` | roadmap §1–§10, each with a status | `complete-and-demonstrated` for ordinary deliverables; for the two conditional storage-escalation items specifically, `proven-executable-on-demand` counts as meeting the bar (research.md "What GA means") |
| `update-evidence` | reference | Aggregate result of the three-site Update Run integration suite (FR-003) |
| `escalation-evidence` | reference | Aggregate result of the Storage-Escalation Guidance Checks (FR-004/FR-005) |
| `ga-determination` | boolean | `true` only when every deliverable meets its bar — this is the literal definition of GA this release establishes |

**Validation rules**: SC-003 requires this review's evidence to stand on its own — it MUST NOT require re-proving v0.04/v0.05's site-type behaviors from scratch (it references their already-closed evidence), and their evidence never depended on this release's capability either.

**Relationships**: aggregates exactly the Update Run suite and the Storage-Escalation Guidance Checks defined above, plus references (not re-runs) every prior release's closing evidence — it is the final node in the dependency chain v0.01 → ... → v0.06 this whole plan sequence describes.
