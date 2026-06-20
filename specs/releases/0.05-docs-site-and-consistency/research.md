# Phase 0 Research: v0.05 — Prove the platform on documentation and across all site types

**Input**: [plan.md](./plan.md) Technical Context — no `NEEDS CLARIFICATION` markers. As with [v0.04](../0.04-project-log-site-type/research.md), this release's relevant research is *which confirmed platform decisions this replication and review must hold to*, not new technology choices — its own Assumptions state it "proves generalization... assumes the existing platform... works unchanged... surfacing any gap as a defect."

## Decision: Why documentation is replicated last (of the three types)

- **Decision**: Replicate to the documentation site type after the project log, and pair it with the cross-type consistency check.
- **Rationale**: [Roadmap §9](../../research/roadmap.md) sequences documentation last among the three types and explicitly times the cross-type review for "only once all three exist" — the spec Summary states this is "the first point at which a true cross-type consistency check is even possible." Running the review here, rather than splitting it across releases, means it happens exactly once, against a complete and stable set of three real sites — not repeatedly against a moving target.
- **Alternatives considered**: Running partial consistency checks after each replication (v0.02, v0.04) — rejected; would produce noisy, incomplete comparisons against an incomplete set, and risk "fixing" something that only looks inconsistent because the third reference point doesn't exist yet.

## Decision: No tag/project filtering for documentation — confirmed absence, not a gap

- **Decision**: The documentation site type deliberately presents no tag- or project-based filtering UI or query parameters; this release *verifies* that absence rather than adding the feature.
- **Rationale**: Confirmed, deliberate design decision recorded at [roadmap §0.1 (architecture & data model)](../../research/roadmap.md) and restated directly in this release's spec Assumptions ("No tag/project filtering is a feature, not a gap"). A hierarchy-based browsing model (TOC/section/order) is the correct shape for documentation; bolting filter UI onto it would be exactly the kind of speculative, one-size-fits-all generalization Principle I rejects — each site type's content model gets the browsing model that fits it, no more.
- **Alternatives considered**: Adding filter UI/params to the docs type "for consistency" with the project log — explicitly rejected by the confirmed decision; would manufacture a feature with no real authoring or visitor need on this content shape, against Principle I.

## Decision: How a tag-filter request against the docs site is handled

- **Decision**: A terminal/agent client requesting the documentation site with filter parameters that only make sense for the project-log type is handled via "a documented, sensible default" — either rejected clearly or cleanly ignored (spec Edge Cases leaves the exact choice to be settled during design, but requires it be documented and consistent).
- **Rationale**: This mirrors the type-mismatch handling already proven in [v0.04](../0.04-project-log-site-type/research.md) ("Type-mismatch query rejection") — reusing one documented contract shape across types rather than inventing a new one per mismatch direction, consistent with Principle II's "pick one owner" discipline applied to error/edge-case contracts as much as to data ownership.
- **Alternatives considered**: Silently accepting and no-op'ing the filter parameter with no documented behavior — rejected; "silently mishandled" is precisely what the spec Edge Cases calls out as the failure mode to avoid, and would leave agent clients unable to distinguish "no matches" from "parameter ignored."

## Decision: How cross-type theming/rendering inconsistencies are triaged

- **Decision**: Any inconsistency the side-by-side review finds — e.g., a shared design token resolving differently per site type — is treated as a defect against the *shared* theming system (§6, built in v0.01), and fixed there; never patched as a one-off per-site override.
- **Rationale**: Directly stated in the spec's own Edge Cases, and a direct consequence of Principle II ("parsing, schema validation, theming... logic MUST NOT be duplicated... pick one owner"). A per-site patch would be the first crack in exactly the single-pipeline guarantee the constitution treats as load-bearing — this decision exists precisely to prevent that crack from ever being rationalized as "just this once."
- **Alternatives considered**: Allowing site-type-specific theme overrides to "smooth over" inconsistencies — rejected outright; would convert a one-time generated-token bug into a permanently forked theming pipeline, the exact divergence Principle II is written to prevent.

## Open questions resolved

None remain — every decision here either reuses a confirmed prior-release contract (type-mismatch handling, shared-theming ownership) or restates a confirmed, already-recorded design decision (the no-filtering absence, the replication ordering). No new technology or architecture choice is being made in this release.
