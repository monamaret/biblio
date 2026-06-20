# Cross-Type Consistency Record

Closing artifact for v0.05's cross-type review (data-model.md "Cross-Type
Consistency Record"; FR-003, FR-004, SC-002). Produced by running
`go test ./test/integration/cross-type-review/...` against all three site
types after the documentation type's replication (Phase 3) completed.

## sites[]

- profile (v0.02)
- log (v0.04)
- docs (this release, v0.05)

## theming-check

| Site type | Result |
|---|---|
| profile | Renders its own preset (`profile-default`) exactly — `TestCrossTypeTheming_EachSiteRendersItsOwnPreset/profile` |
| log | Renders its own preset (`log-default`) exactly — `TestCrossTypeTheming_EachSiteRendersItsOwnPreset/log` |
| docs | Renders its own preset (`docs-default`) exactly — `TestCrossTypeTheming_EachSiteRendersItsOwnPreset/docs` |

No cross-type bleed-through: each type's distinguishing tokens
(`color-fg`, `color-accent`) are unique per type
(`TestCrossTypeTheming_NoBleedThrough`). `color-bg` is intentionally shared
between profile and docs (two light themes vs. log's dark one), and
`font-heading` is intentionally shared across all three — both confirmed,
deliberate design choices, not bleed-through (see test comments).

## terminal-rendering-check

| Site type | Result |
|---|---|
| profile | Clean Markdown/ANSI rendering, type-agnostic negotiation code (`templates/site/profile/app/negotiation/accept.test.ts`) |
| log | Clean Markdown/ANSI rendering against dated-entry content (`templates/site/log/app/negotiation/accept.test.ts`) |
| docs | Clean Markdown/ANSI rendering against hierarchical-page content (`templates/site/docs/app/negotiation/accept.test.ts`) |

One shared negotiation owner confirmed: `accept.ts`/`ansi.ts`/`wasm.d.ts`
are byte-identical across all three site types
(`TestCrossTypeTerminalRendering_OneSharedNegotiationOwner`) — no per-type
fork exists for any client to render inconsistently against. Each type's
content shape is appropriate to its own model — profile's `role`/`summary`,
log's `date`/`tags`, docs's `section`/`order` — with no cross-contamination
of another type's fields (`TestCrossTypeTerminalRendering_ContentShapeIsAppropriatePerType`).

## inconsistencies-found[]

None. (One over-strict test assertion was found and corrected during
authoring — see Phase 4's T029 note in tasks.md — but it was a defect in
the *review test itself*, not in `internal/theme` or the negotiation code;
no shared-code defect was filed.)
