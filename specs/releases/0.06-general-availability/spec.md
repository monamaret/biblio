# v0.06 — General Availability

**Created**: 2026-06-07
**Status**: Draft
**Roadmap deliverables**: [§9 Replicate to remaining (priority) site types](../research/roadmap.md#9-replicate-to-the-remaining-priority-site-types) (`update`/`upgrade` mechanics) · [§10 Storage escalation (conditional)](../research/roadmap.md#10-storage-escalation-conditional--not-built-up-front)
**Depends on**: [v0.04](0.04-release.md) and [v0.05](0.05-release.md) — the update/upgrade mechanics are, per the roadmap, "only answerable once more than one site exists to migrate"
**Feature specs**: [009](../features/009-replicate-site-types.md) · [010](../features/010-storage-escalation.md)

## Summary

With three real, deployed sites of three different types now in hand, the author can finally design and prove out the one piece of the CLI surface the roadmap deliberately left open until this point: a reliable way to carry an existing site forward as the underlying scaffold evolves, without losing authored content or unintentionally overwriting customizations. With that in place — and with the conditional storage-escalation paths confirmed ready to engage on demand — every roadmap deliverable is complete or proven ready, and the platform reaches **General Availability**.

## User Scenarios & Acceptance Criteria

### Scenario 1 — Update a site to a newer scaffold version without losing content

As the author, I run the update flow against a live site and my content survives intact while the scaffold-tracked parts move forward.

1. **Given** an existing deployed site and a newer version of the underlying site scaffold, **When** the author runs the update/upgrade flow against that site, **Then** the parts of the site meant to track the scaffold are brought up to date, the author's own authored content and intentional customizations are preserved untouched, and any unavoidable conflicts are clearly surfaced for the author to resolve rather than silently overwritten.

### Scenario 2 — The update flow works across all three site types

As the author, I run the same update flow against my profile, project-log, and documentation sites and it succeeds for each.

1. **Given** all three real sites (profile, project log, documentation), **When** the update/upgrade flow is run against each in turn, **Then** it completes successfully for every site type, not just the one it may have first been designed against.

### Scenario 3 — Storage-escalation guidance is checked against real sites

As the author, I check my real sites against the documented storage-escalation guidance and trust what it reports.

1. **Given** the platform's documented storage-escalation guidance and the real scale of the three deployed sites, **When** the author checks each site against the documented thresholds, **Then** the guidance correctly reports whether any site has crossed a threshold, and — for any path it recommends — the recommended steps can actually be carried out as written.

### Scenario 4 — GA readiness review

As the author, I review every roadmap deliverable against this release's bar and confirm the platform is ready for General Availability.

1. **Given** every roadmap deliverable (§1–§10) and every release in the v0.01–v0.06 plan, **When** the author reviews them against this release's acceptance bar, **Then** all are complete, or — for the conditional storage-escalation items specifically — proven ready to engage on demand, marking the platform as Generally Available.

## Functional Requirements

- **FR-001**: The system MUST provide a way for the author to bring an existing deployed site's scaffold-tracking parts up to date with a newer version of the underlying scaffold, while preserving the author's own authored content and intentional customizations.
- **FR-002**: The system MUST clearly surface, rather than silently resolve, any case where an update would otherwise overwrite an author's intentional customization — leaving the decision to the author.
- **FR-003**: The update/upgrade capability MUST work correctly across all three supported site types, not only the one it may first have been validated against.
- **FR-004**: The system MUST provide documented guidance the author can check any deployed site against to determine whether it has crossed a storage-scale or search-quality threshold that calls for an escalation, and that guidance MUST correctly reflect each site's real, current scale.
- **FR-005**: For each storage-escalation path the guidance can recommend, the system MUST make that path's steps concrete enough to actually carry out, not merely described in principle — so that if a real site ever needs it, the path is provably executable on demand rather than only theoretically planned.
- **FR-006**: By the close of this release, every roadmap deliverable (§1–§10) MUST be either complete and demonstrated, or — for the conditional storage-escalation items specifically — proven ready to engage the moment real usage calls for them; this state constitutes General Availability.

## Edge Cases

- What happens when the update/upgrade flow encounters a site whose author has heavily customized parts of the scaffold that a newer version also wants to change — does the flow surface the conflict for a human decision, rather than choosing a side automatically?
- What happens if, by GA, real usage on any of the three sites has in fact crossed a storage-escalation threshold — does GA then require that escalation to have actually been carried out, or only that the path be proven executable on demand?
- How does the update flow behave if it's interrupted partway through (e.g., the author cancels mid-run) — does the site remain in a consistent, recoverable state?

## Success Criteria

- **SC-001**: The author can carry forward all three real, live sites to a newer scaffold version in one sitting, finishing with each site's authored content fully intact and any genuine conflicts clearly identified rather than silently lost.
- **SC-002**: A complete review against every deliverable in [roadmap.md §1–§10](../research/roadmap.md) shows each one complete and demonstrated — or, for the two explicitly conditional storage-escalation items, shown to be executable on demand even though no real site has yet required them.
- **SC-003**: This release's update-flow evidence never requires re-proving v0.04/v0.05's site-type behaviors from scratch, and vice versa — each release's evidence stands on its own.

## Assumptions

- **What "GA" means here**: General Availability is defined, per the roadmap's own remaining open items, as the point where every deliverable is either complete-and-demonstrated or — for the two items the roadmap explicitly marks as conditional and "not scheduled work" (the D1 and Vectorize storage-escalation paths) — proven *executable on demand* without requiring that any real site has actually crossed a threshold yet. This avoids inventing a GA bar the roadmap doesn't itself set, while still giving "GA" a concrete, testable meaning.
- **Sequencing is load-bearing, not arbitrary**: The roadmap is explicit that the update/upgrade mechanics are "only answerable once more than one site exists to migrate" — meaning this release's central piece of design work is *only possible* once [v0.04](0.04-release.md) and [v0.05](0.05-release.md) have produced the second and third real sites. It is therefore correctly sequenced last.
- **No timeline commitments**: Consistent with the roadmap's explicit choice to avoid effort/timeline estimates, this spec defines what GA contains and how its readiness is demonstrated, not when it ships.
- **Reserved version numbers**: The version gap between v0.03 and v0.04–v0.06 in the overall plan is assumed to be intentional headroom for hardening and feedback incorporation, not missing scope — no additional deliverables are assumed beyond those already named in [roadmap.md §1–§10](../research/roadmap.md).
