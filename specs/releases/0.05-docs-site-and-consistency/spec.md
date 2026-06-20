# v0.05 — Prove the platform on documentation and across all site types

**Created**: 2026-06-07
**Status**: Draft
**Roadmap deliverables**: [§9 Replicate to remaining (priority) site types](../research/roadmap.md#9-replicate-to-the-remaining-priority-site-types) (documentation portion + cross-type consistency check)
**Depends on**: [v0.04](0.04-release.md) — the cross-type consistency check is only meaningful once a second site type exists alongside the pilot
**Feature specs**: [009](../features/009-replicate-site-types.md)

## Summary

The author scaffolds, authors, and deploys a project-documentation site — the third and final supported site type — and then confirms that visual theming and terminal/CLI-facing rendering both behave correctly and consistently across all three site types now that all of them exist side by side. This is the first point at which a true cross-type consistency check is even possible: the roadmap explicitly calls for confirming "theming presets and the terminal/CLI view render correctly across all three site types" only once all three exist.

## User Scenarios & Acceptance Criteria

### Scenario 1 — Browse a documentation site by its hierarchy

As a visitor, I navigate the docs site by its table of contents and section/order structure, with no tag/project filters cluttering it.

1. **Given** a freshly scaffolded documentation site with pages organized into sections, **When** a visitor opens it, **Then** the table of contents/navigation correctly reflects the section/order structure and every page is reachable and renders its content correctly.
2. **Given** the deployed documentation site, **When** a visitor or an agent looks for tag- or project-based filtering (as exists on the project-log site type), **Then** no such filtering is present — confirming the documentation type's deliberately narrower, hierarchy-based browsing model.

### Scenario 2 — Compare theming across all three site types

As the author, I view all three deployed sites side by side and see each one's own theme correctly and consistently applied.

1. **Given** all three site types now deployed, **When** each is visited in turn, **Then** each displays its own distinct, correctly-applied visual theme with no cross-type bleed-through or inconsistency.

### Scenario 3 — Compare terminal rendering across all three site types

As a terminal/agent client, I get a correct, type-appropriate Markdown view from any of the three sites.

1. **Given** all three site types now deployed, **When** each is requested by a terminal/Markdown-preferring client, **Then** each returns a clean, correctly rendered terminal-friendly view appropriate to its own content shape.

## Functional Requirements

- **FR-001**: The system MUST let the author scaffold, author content for, and deploy a project-documentation site using the same CLI surface already proven for the other two site types.
- **FR-002**: The system MUST let a visitor browse the documentation site's table of contents/navigation and read its pages organized by their section and order, with no tag- or project-based filtering presented for this site type.
- **FR-003**: The system MUST apply each of the three site types' own distinct visual theme correctly and consistently when all three are deployed and viewed side by side.
- **FR-004**: The system MUST return a correctly rendered, clean terminal/Markdown view, appropriate to each site type's own content shape, when any of the three deployed sites is requested by a terminal- or Markdown-preferring client.

## Edge Cases

- How does the platform behave when a documentation site's authored structure doesn't cleanly map to a single section/order hierarchy (e.g., a page that logically belongs in two sections) — is there a documented, sensible default?
- How does the platform handle a terminal/agent client requesting content from the documentation site type using filter parameters that only make sense for the project-log type (e.g., a tag filter) — is it rejected clearly, or simply ignored?
- What happens when the cross-type theming review finds a subtle inconsistency (e.g., a shared token resolving differently per type) — is it treated as a defect against the shared theming system (§6), not a one-off per-site fix?

## Success Criteria

- **SC-001**: A real project-documentation site is live, and a visitor can reach and correctly read every authored page through the table-of-contents/navigation structure alone.
- **SC-002**: All three supported site types are simultaneously live and real, each one visually distinct per its own theme and each one returning a correct terminal/Markdown rendering on request — with zero cross-type inconsistencies found during a side-by-side review.
- **SC-003**: This release is independently demonstrable purely on this site type's evidence plus the cross-type review — it requires no update/upgrade-flow work (that's v0.06) to show its value.

## Assumptions

- **Third site, same platform, no platform changes assumed**: As with [v0.04](0.04-release.md), this release proves generalization — it assumes the existing platform (CLI, pipeline, theming, scaffold, MCP server, terminal view) works unchanged for the documentation site type, surfacing any gap as a defect to fix.
- **No tag/project filtering is a feature, not a gap**: The absence of tag/project filter UI and query params on the documentation site type is a confirmed, deliberate design decision (see [roadmap §0.1](../research/roadmap.md#01-architecture--data-model)) — this release verifies that absence, it doesn't add the feature.
- **No timeline commitments**: Consistent with the roadmap's explicit choice to avoid effort/timeline estimates, this spec defines what the release contains and how its completion is demonstrated, not when it ships.
