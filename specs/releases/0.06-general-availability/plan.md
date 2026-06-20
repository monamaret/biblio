# Implementation Plan: v0.06 — General Availability

**Branch**: `106-release-0-06-general-availability` | **Date**: 2026-06-07 | **Spec**: [spec.md](./spec.md)

**Input**: Release specification from `/specs/releases/0.06-general-availability/spec.md` (mirrors [`specs/releases/0.06-release.md`](../0.06-release.md)), depends on [v0.04](../0.04-project-log-site-type/plan.md) and [v0.05](../0.05-docs-site-and-consistency/plan.md), built on feature items [009](../../features/009-replicate-site-types.md), [010](../../features/010-storage-escalation.md)

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Design and prove out the one piece of the CLI surface the roadmap deliberately left open until now — a reliable `biblio-cli update`/`upgrade` flow that carries an existing deployed site forward to a newer scaffold version, preserving the author's content and intentional customizations while surfacing genuine conflicts for human resolution — and prove it works across all three real site types now in hand. Pair that with a documented, checkable storage-escalation guide whose recommended D1/Vectorize paths are confirmed *executable on demand* against the three real sites' actual scale. With both in place and every roadmap deliverable (§1–§10) reviewed against this release's bar, the platform reaches **General Availability**. This is correctly sequenced last: per the roadmap, the update/upgrade design is "only answerable once more than one site exists to migrate" (spec Assumptions) — both v0.04 and v0.05 had to land first.

## Technical Context

**Language/Version**: Go 1.22+ (`biblio-cli` — new `update`/`upgrade` commands extending the Cobra tree established in v0.01/v0.02, per the §0.3-confirmed decision to fold this story into the existing build/deploy command surface)

**Primary Dependencies**: Same confirmed stack as every prior release — no new runtime dependency is introduced for the update/upgrade mechanics themselves; they are diff/merge logic over the existing scaffold-template + authored-content separation already established in v0.01. The storage-escalation guidance references the already-named conditional paths (Cloudflare D1, Vectorize) without adopting either as a running dependency in this release (they remain "proven executable on demand," not built up front)

**Storage**: Unchanged for the platform itself — still the generated JSON `content-index.json` across all three real sites (Principle V); this release's storage-relevant work is entirely *documentation and on-demand-execution proof* for the two conditional escalation paths (D1 for relational/volume growth, Vectorize for semantic-search-quality growth), explicitly not a migration of any real site that hasn't crossed a threshold

**Testing**: Go unit tests for the update/upgrade diff/merge/conflict-detection logic (scaffold-tracked vs. author-customized region classification — written first per Principle VI), plus an end-to-end integration suite that runs the update/upgrade flow against all three real sites (profile, project log, documentation) and confirms each completes successfully with content intact (FR-003, Scenario 2); a "storage-escalation guidance" check that runs the documented thresholds against the three sites' real current scale and confirms the recommended paths' steps are concretely executable (FR-004, FR-005, Scenario 3); and a closing GA-readiness review checklist cross-referencing every roadmap deliverable §1–§10 against this release's acceptance bar (Scenario 4, SC-002)

**Target Platform**: The same three already-deployed single-Worker sites from v0.02/v0.04/v0.05 — the update/upgrade flow runs *against* them in place; no new hosting surface

**Project Type**: CLI feature (new `biblio-cli update`/`upgrade` commands) plus a documentation/verification deliverable (storage-escalation guide) plus a closing cross-release review — no new application architecture

**Performance Goals**: All three real sites carried forward to a newer scaffold version "in one sitting," each finishing with authored content fully intact and genuine conflicts clearly identified (SC-001); storage-escalation guidance correctly reports each site's real, current threshold status (Scenario 3, FR-004)

**Constraints**: The update flow MUST preserve authored content and intentional customizations untouched, bring scaffold-tracked parts forward, and clearly surface — never silently resolve — any case where it would otherwise overwrite an intentional customization, leaving the decision to the author (FR-001, FR-002, Scenario 1, spec Edge Cases); it MUST work correctly across all three site types, not only whichever it's first validated against (FR-003, Scenario 2); an interrupted update MUST leave the site in a consistent, recoverable state (spec Edge Cases); GA itself requires every roadmap deliverable to be either complete-and-demonstrated, or — for the two explicitly-conditional storage-escalation items only — proven executable on demand without requiring any real site to have actually crossed a threshold yet (FR-006, spec Assumptions "What GA means here")

**Scale/Scope**: Exactly the three real sites already live from v0.02/v0.04/v0.05 — no synthetic fixtures, no speculative fourth site type, no timeline commitments (spec Assumptions); this release's evidence must stand on its own without re-proving v0.04/v0.05's site-type behaviors, and vice versa (SC-003)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Single-Author, Minimal-Dependency Bias)**: PASS — the update/upgrade mechanics are diff/merge logic over the existing scaffold/content separation, not a new dependency or abstraction layer; the storage-escalation guidance documents *conditional, not-yet-needed* escalation paths exactly as Principle I requires ("deferred until a documented threshold is actually crossed... each one gets its own research pass and Decision Log entry when it happens") — this release proves the path is *ready*, not that it must be taken now.
- **Principle II (One Pipeline, One Source of Truth)**: PASS — the update/upgrade flow's "scaffold-tracked vs. author-customized" classification must itself be derived from the single canonical scaffold-template definitions `biblio-cli` already owns (Principle II) — never a second, divergent notion of "what the scaffold looks like" maintained separately for migration purposes.
- **Principle III (CLI-First, Guided Experience)**: PASS — `update`/`upgrade` land in the *same* `biblio-cli` Cobra tree as `new`/`dev`/`build`/`deploy`, exactly per the confirmed §0.3 resolution that folded this story into one consistent CLI-first surface from the start — never a separate migration tool or raw script.
- **Principle IV (Text-First, Progressive Enhancement)**: N/A — this release adds no new MCP tools or terminal-rendering surface; it operates entirely at the CLI/scaffold-management layer.
- **Principle V (Simplest Topology That Can Grow)**: GATE — this is the release where Principle V's "escalate only on a documented, real trigger" is *exercised as designed*: the two conditional storage-escalation items are proven executable on demand precisely so that, when a real trigger eventually fires, the escalation is "a new, separately-decided change with its own Decision Log entry," not an emergency improvisation. GA explicitly does NOT require either escalation to have actually been carried out (FR-006, spec Assumptions) — only that the path is concretely walkable. Actually executing either escalation here, absent a real crossed threshold, would itself be the speculative-build-up-front Principle V forbids.
- **Principle VI (Test-First Development, NON-NEGOTIABLE)**: GATE — failing tests for the diff/merge/conflict-detection logic, the cross-type update integration suite, and the storage-escalation-guidance check MUST be authored and reviewed before implementation; per the constitution's Release Process, this release closes only once both unit and end-to-end integration tests exist and pass for everything it ships.
- **Security & Operational Posture**: PASS — nothing here touches `/mcp`, auth, or access control; `mcp_enabled` and the no-auth posture are unchanged across all three sites through the update flow.

No violations requiring justification — Complexity Tracking is not needed for this release.

## Project Structure

### Documentation (this feature)

```text
specs/releases/0.06-general-availability/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
# new biblio-cli commands (Go) — extend the existing Cobra tree, no new module/tool
internal/
├── cli/
│   ├── update.go            # `biblio-cli update` — bring scaffold-tracked parts forward
│   └── upgrade.go           # `biblio-cli upgrade` — version-targeted variant of the same flow
├── migrate/
│   ├── classify.go          # scaffold-tracked vs. author-customized region classification
│   │                        #  (derived from the single canonical scaffold-template definitions, Principle II)
│   ├── merge.go             # diff/merge logic preserving authored content + customizations
│   ├── conflict.go          # conflict surfacing — never silent overwrite; recoverable on interruption
│   └── update_run.go        # Update Run orchestration (classify → advance → preserve → surface → commit)
│                            #  and the preserved/advanced/conflicts report, shared by `update` and `upgrade`
└── escalation/
    ├── guidance.go          # storage-escalation threshold checks against real content-index scale,
    │                        #  reporting D1/Vectorize readiness without adopting either as a dependency
    └── ga_review.go         # GA Readiness Review artifact — aggregates Update Run + Guidance Check
                             #  evidence by reference into the closing ga-determination record

docs/
└── storage-escalation-guide.md   # FR-004/FR-005 deliverable: documented thresholds + concretely
                                  #  walkable D1 and Vectorize escalation paths, proven executable on demand

test/
├── unit/
│   └── migrate/             # classify/merge/conflict-detection unit tests (failing-first, Principle VI)
└── integration/
    ├── update-three-sites/  # update/upgrade run against the real profile, log, and docs sites
    │                        #  in turn — each completes successfully, content intact (FR-003, Scenario 2)
    ├── escalation-guidance/ # runs documented thresholds against real site scale; confirms
    │                        #  recommended paths' steps are concretely executable (FR-004, FR-005)
    └── ga-readiness-review/ # closing checklist cross-referencing roadmap §1–§10 against this
                             #  release's acceptance bar (Scenario 4, SC-002)
```

**Structure Decision**: New `internal/migrate` and `internal/escalation` packages plus two new thin Cobra commands (`update`, `upgrade`) extend the existing `biblio-cli` module — no new tool, service, or top-level project. The storage-escalation guide is a documentation deliverable proven against real sites, not new running infrastructure (Principle V — the conditional paths stay conditional). The closing GA-readiness review is itself a structured integration-test artifact, mirroring how v0.05's cross-type consistency review was modeled, so "GA" is demonstrated the same disciplined way every other release-closing claim has been.

## Complexity Tracking

*No Constitution Check violations — section intentionally left empty.*
