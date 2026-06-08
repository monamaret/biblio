# Phase 0 Research: v0.04 — Prove the platform on the technical project log

**Input**: [plan.md](./plan.md) Technical Context — no `NEEDS CLARIFICATION` markers. This release's spec Assumptions state directly that it "proves generalization, not redesign," assuming the v0.01–v0.03 platform "works unchanged for a new site type, surfacing any gap found as a defect to fix rather than a sign the platform needs reshaping" — so the relevant research is *which platform decisions this replication must hold to*, not new technology choices.

## Decision: Why the project log is replicated first

- **Decision**: Replicate to the technical project-log site type before the documentation type.
- **Rationale**: [Roadmap §9](../../research/roadmap.md) names this generalization test — "validate tag/project filtering and the daily-update authoring flow" — as the explicit purpose of replicating to this type *first*, prioritizing it because its content model (a running, cross-cutting timeline) differs the most from the already-proven profile site. Front-loading the largest structural difference surfaces generalization defects earliest, when they're cheapest to fix against only one prior site type.
- **Alternatives considered**: Replicating documentation first — rejected by the roadmap's own ordering; the project log's timeline/filter shape is the harder generalization case and the more valuable one to prove early.

## Decision: Filtering computed over the existing JSON index — no new storage/search

- **Decision**: Date-range, tag, and project filtering (FR-003) is computed directly over the same generated `content-index.json` already produced by the unchanged v0.01 indexing pipeline — no database, no search-index escalation.
- **Rationale**: Principle V (Simplest Topology That Can Grow) requires escalating storage/search only on a documented, real volume/quality trigger. A single author's project-log entry count sits comfortably within "generated, in-memory JSON index" scale; introducing filtering is itself not such a trigger. This keeps the second site type on the exact same storage posture as the first (Principle I/V consistency).
- **Alternatives considered**: Pre-wiring D1 for "queryable" filtering, or a client-side search library — both rejected; Principle I calls out "storage beyond a bundled JSON index... deferred until a documented threshold is actually crossed" by name, and adding either here would be exactly the speculative escalation the constitution forbids.

## Decision: Authoring rhythm — same flow, no special cases

- **Decision**: The "write a new dated entry" workflow MUST be the same `biblio-cli` authoring flow already used for the prior entry — no one-off structural or config adjustments to publish each subsequent entry (FR-004, Scenario 3).
- **Rationale**: This is the load-bearing generalization claim of the whole release: per spec Scenario 3 and SC-002, the author must be able to describe and repeat the same rhythm for "several consecutive published entries with zero one-off adjustments." Anything less would mean the platform doesn't actually generalize — exactly the failure mode this release exists to catch (and, per the Assumptions, fix as a defect rather than redesign around).
- **Alternatives considered**: A log-specific "new entry" CLI subcommand or scaffolding step — rejected; FR-001 requires using "the same CLI surface already proven for the profile site type," and a type-specific authoring command would itself be the kind of special-case the release is meant to prove unnecessary.

## Decision: Type-mismatch query rejection

- **Decision**: An MCP/terminal request whose shape only makes sense for the profile site type (e.g., "get résumé") MUST be rejected clearly by the project-log site, mirroring the v0.03 "unsupported capability" contract rather than introducing a new error shape.
- **Rationale**: [v0.03's contract](../0.03-agent-terminal-layer/contracts/mcp-tools.md) already defines this exact failure mode ("a request for a capability the site type doesn't define MUST be rejected with a clear, structured error"). Reusing it here — rather than inventing a project-log-specific rejection format — is itself part of proving the v0.03 layer generalizes (FR-005), and keeps the contract surface singular per Principle II.
- **Alternatives considered**: A bespoke "type mismatch" error format for the project-log site — rejected; would fragment the MCP error contract across site types for no benefit, against Principle II's "pick one owner" discipline.

## Open questions resolved

None remain — every Technical Context decision in this release is "hold the existing platform decision constant and prove it generalizes," which is itself the confirmed research posture stated in the spec's own Assumptions section. No new technology, storage, or architecture choice is being made here.
