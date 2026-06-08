# Phase 0 Research: v0.06 — General Availability

**Input**: [plan.md](./plan.md) Technical Context — no `NEEDS CLARIFICATION` markers. This release's defining research characteristic is the inverse of v0.04/v0.05's: rather than holding existing decisions constant, it *exercises* a decision the roadmap deliberately deferred — "the update/upgrade mechanics are only answerable once more than one site exists to migrate" (spec Assumptions, restating [roadmap §9](../../research/roadmap.md)). This document records why that sequencing is load-bearing and how the deferred design questions resolve now that the precondition is met.

## Decision: Why update/upgrade design had to wait for v0.04 and v0.05

- **Decision**: Treat this release's central design work — the update/upgrade flow — as only now answerable, and therefore correctly sequenced last.
- **Rationale**: The roadmap is explicit, and the spec restates it directly: designing a "carry an existing site forward" flow requires *more than one real site* to migrate against, or the design is pure speculation about what "scaffold-tracked vs. customized" looks like in practice. v0.01 produced one site shape; v0.02 proved it deployable; v0.03 added the agent layer; only v0.04 and v0.05 produced the second and third real, structurally different sites that make "does this generalize across site types" an answerable question rather than a guess (FR-003, Scenario 2).
- **Alternatives considered**: Designing the update/upgrade flow speculatively earlier (e.g., alongside v0.01's scaffold work) — rejected; this is precisely the kind of "designing for a hypothetical future need" Principle I forbids ("never against hypothetical... future needs"). Building it before three real sites existed to validate against would have meant designing the conflict-classification logic against zero evidence of what real customizations look like.

## Decision: Update/upgrade folded into the existing CLI surface

- **Decision**: `biblio-cli update`/`biblio-cli upgrade` land in the same Cobra command tree as `new`/`dev`/`build`/`deploy` — not a separate migration tool, plugin, or raw script.
- **Rationale**: Confirmed directly in [roadmap §0.3](../../research/roadmap.md): "the site-update/re-scaffold story... is folded into this same command surface (e.g. `biblio-cli update`/`biblio-cli upgrade`), keeping one consistent CLI-first surface end to end rather than splitting 'dev/build/deploy is one command, updates are raw `package.json`/CI.'" This decision was made *before* v0.02 shipped — this release is where it's finally implemented, not redesigned.
- **Alternatives considered**: A standalone `biblio-migrate` tool — rejected by the confirmed §0.3 decision and by Principle III's "never as... a separate standalone tool" rule, stated almost verbatim.

## Decision: Conflict handling philosophy — surface, never silently resolve

- **Decision**: Where the update would overwrite an author's intentional customization, the flow MUST surface that conflict clearly for the author to decide — never silently choose a side (FR-002, spec Edge Cases).
- **Rationale**: This is a direct extension of the project's single-author trust model (Principle III's rationale: "a single author trusts end-to-end"; Principle VI's rationale: "a regression here silently corrupts every site it touches"). An update tool that silently overwrites customizations would corrode exactly the trust the CLI-first, single-binary, single-pipeline design exists to build and preserve.
- **Alternatives considered**: A "scaffold always wins" or "customizations always win" automatic policy — both rejected; either would silently destroy real authored work or real upstream improvements respectively, and the spec Edge Cases explicitly frames "does the flow surface the conflict for a human decision, rather than choosing a side automatically" as the question this release must answer "yes" to.

## Decision: What "GA" means and how the storage-escalation items count toward it

- **Decision**: General Availability means every roadmap deliverable (§1–§10) is either complete-and-demonstrated, or — for the two items the roadmap explicitly marks conditional and "not scheduled work" (D1 and Vectorize storage escalation) — proven *executable on demand*, without requiring that any real site has actually crossed a threshold.
- **Rationale**: Stated directly and carefully in the spec's own Assumptions ("What 'GA' means here") specifically to avoid "inventing a GA bar the roadmap doesn't itself set." This is Principle V's escalation discipline made concrete and checkable: "each escalation step is a new, separately-decided change with its own Decision Log entry — not something to pre-build 'just in case.'" Proving the path *walkable* satisfies that discipline; actually walking it without a real trigger would violate it.
- **Alternatives considered**: Requiring an actual D1/Vectorize migration as a GA gate — rejected; would force a speculative escalation with no real trigger, exactly what Principle I and Principle V both forbid by name. Requiring nothing at all from the conditional items — also rejected; would leave "GA" meaning less than the roadmap's own bar, undermining the "complete review against every deliverable" claim SC-002 makes.

## Decision: Reserved version numbers v0.04–v0.06 are intentional headroom, not missing scope

- **Decision**: Treat the gap between v0.03 and v0.04–v0.06 as deliberate headroom for hardening/feedback incorporation — assume no additional deliverables beyond those already named in roadmap §1–§10.
- **Rationale**: Stated directly in the spec's own Assumptions ("Reserved version numbers"). This matters for planning discipline: it means this plan should not invent GA-blocking work that isn't traceable to a named roadmap deliverable — doing so would itself be the kind of scope-invention Principle I and the GA-definition decision above both guard against.
- **Alternatives considered**: Treating the version gap as a signal that undocumented scope exists and should be reverse-engineered into this plan — rejected; would manufacture requirements with no Decision Log basis, violating the constitution's Decision Recording & Research Workflow (every plan traces to confirmed decisions, not invented ones).

## Open questions resolved

None remain — the central "deferred until X exists" question this release exists to answer is resolved by the simple fact that v0.04 and v0.05 have, per this plan's sequencing, already produced the second and third real sites; every other decision here either restates a confirmed prior commitment (§0.3 command-surface folding, Principle V escalation discipline) or is the spec's own explicitly-authored definition of GA.
