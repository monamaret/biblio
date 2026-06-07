# Feature: Replicate to the remaining (priority) site types

**Status:** Not started
**Roadmap reference:** [roadmap.md §9 — Replicate to the remaining (priority) site types](../research/roadmap.md#9-replicate-to-the-remaining-priority-site-types)

## Summary

Build and deploy the two scaffolds the research calls out as priorities — the technical project log and project documentation site types — validating that the platform generalizes beyond the pilot's personal profile/CV case.

## Confirmed decisions this builds on

- [Background & Context — domain assumptions](../research/001-static-site-platform-with-mcp.md#background--context): "the two priority scaffolds are the technical project log (daily updates, high text volume) and project documentation — both need rich text formatting plus project- and tag-based filtering/sorting"
- [Performance, Reliability & Operations — rollout/migration plan](../research/001-static-site-platform-with-mcp.md#performance-reliability--operations): replicate for project-log and documentation site types after the pilot validates
- [Architecture Patterns — frontmatter schema and tool surface per site type](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the confirmed content models and tool surfaces these two site types must implement
- [Decision Log — fold the site-update/re-scaffold story into the `biblio-cli` command surface](../research/001-static-site-platform-with-mcp.md#decision-log): the *direction* is locked — it ships as a `biblio-cli update`/`biblio-cli upgrade` command, not a separate tool or manual process; this feature designs and implements the concrete mechanics (the one part of §0.3 still genuinely open)

## Scope

- [ ] Scaffold and deploy a technical project log site; validate tag/project filtering and the daily-update authoring flow
- [ ] Scaffold and deploy a project documentation site; validate TOC/nav generation and `section`/`order`-based browsing (no tag/project filtering — confirmed as out of scope for this site type, see [001 → Decision Log](../research/001-static-site-platform-with-mcp.md#decision-log))
- [ ] Confirm theming presets and the terminal/CLI view render correctly across all three site types
- [ ] Design and implement the concrete mechanics of the `biblio-cli update`/`biblio-cli upgrade` command (diffing strategy, conflict handling, what gets re-templated vs. preserved) for updating/re-scaffolding existing sites as shared templates evolve — direction already confirmed as a CLI command (see [001-biblio-cli-foundations](001-biblio-cli-foundations.md)); this feature designs the mechanics, first answerable once more than one site exists

## Out of scope

- Building the underlying scaffold/CLI/MCP-server pieces these site types reuse — those are [001](001-biblio-cli-foundations.md)–[007](007-terminal-cli-view.md); this feature only proves they generalize
- Designing new site types beyond the three already defined in the research — out of scope for the whole project per [001 → Problem Statement non-goals](../research/001-static-site-platform-with-mcp.md#problem-statement)

## Open questions

The following [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) item partially shapes this feature:

- Site-update / re-scaffold story — **mechanics only** (§0.3) — the *direction* is already locked (it ships as a `biblio-cli update`/`biblio-cli upgrade` command, see [001 → Decision Log](../research/001-static-site-platform-with-mcp.md#decision-log)); the concrete diffing/conflict-handling/re-templating mechanics first become a real, answerable question once this feature creates a second and third site to migrate

(The docs-schema question is resolved — no `tags`/`project` fields or filtering on the Docs Page type, since a docs site is inherently scoped to one project and `section`/`order` is its natural organizing principle; see [001 → Decision Log](../research/001-static-site-platform-with-mcp.md#decision-log). This feature validates `section`/`order`-based browsing for docs, not filtering parity with the log.)

## Dependencies

- [008-pilot-site-validation](008-pilot-site-validation.md) — the pilot must succeed first; this feature is explicitly the generalization step that follows it

## Acceptance criteria

- A real technical project log site is live, with working tag/project filtering and a validated daily-update authoring flow
- A real project documentation site is live, with working TOC/nav generation and `section`/`order`-based browsing (no tag/project filtering, by design)
- All three deployed sites (pilot + these two) render correctly under their respective theming presets and the terminal/CLI view
- The `biblio-cli update`/`biblio-cli upgrade` command has a concrete, working mechanics design (diffing, conflict handling, re-templating vs. preservation), informed by having actually migrated more than one site

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
