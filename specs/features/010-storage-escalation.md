# Feature: Storage escalation (conditional)

**Status:** Not started — **conditional, not scheduled**
**Roadmap reference:** [roadmap.md §10 — Storage escalation (conditional — not built up front)](../research/roadmap.md#10-storage-escalation-conditional--not-built-up-front)

## Summary

The D1 and Vectorize escalation paths for the MCP server's content index — explicitly **not** built up front. This feature item exists so the work is *defined* and ready to pick up the moment a real site crosses one of the documented heuristic thresholds, not so it gets scheduled now.

## Confirmed decisions this builds on

- [Architecture Patterns — progressive-enhancement storage approach](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the three-tier plan (bundled JSON → D1 → Vectorize), each tier remaining a **derived artifact** generated at build time from the same content directory — never hand-maintained
- [Decision Log — per-site-type frontmatter schemas and storage-escalation thresholds](../research/001-static-site-platform-with-mcp.md#decision-log): the heuristic thresholds — JSON → D1 around the low-thousands-of-entries/megabyte range or when genuinely relational queries are needed; keyword → Vectorize when real-usage signals show conceptual queries are poorly served by exact-keyword matching, most likely first on the docs site type
- [000-tech-stack.md — Decision Log](../research/000-tech-stack.md#decision-log): OpenRouter selected (over DeepSeek, which has no embeddings endpoint) as the build-time embedding-generation provider, contingent on this escalation actually happening

## Scope

- [ ] *(Conditional — bundled JSON index outgrows comfortable in-memory handling)* Provision a D1 database, add the binding, and switch content-index population to populate D1 from the same content directory at build time
- [ ] *(Conditional — real-usage signals show keyword search underperforming on conceptual queries, most likely on a docs site)* Wire OpenRouter build-time embedding generation, provision a Vectorize index, and switch `search_docs`/`search_entries` to query it

## Out of scope

- Building either path now, speculatively — explicitly against the confirmed progressive-enhancement philosophy; this feature is a *definition*, not a schedule commitment
- Choosing the specific OpenRouter embedding model — tracked separately as a deferred follow-on in [roadmap.md §0.4](../research/roadmap.md#04-deferred--conditional-follow-ons), to be resolved only once this feature is actually triggered

## Open questions

- Storage-escalation migration steps (§0.3 in [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions)) — concrete step-by-step instructions are explicitly deferred until a real site crosses a threshold; writing them now would be speculative
- Which embedding model to use within OpenRouter (§0.4) — only matters if the Vectorize path of this feature is actually triggered

## Dependencies

- [002-content-pipeline](002-content-pipeline.md) — both escalation paths are alternate population strategies for the same canonical content-index pipeline, not replacements for it
- [004-in-worker-mcp-server](004-in-worker-mcp-server.md) — the query layer this feature would extend/replace already needs to exist and be observed in real use first

## Acceptance criteria

*(To be defined concretely once triggered — premature to specify now.)* At a minimum, before this feature is picked up at all:

- A real, deployed site has been observed crossing one of the documented heuristic thresholds (entry count/index size for D1; search-quality signal for Vectorize) — not a projection, an observation
- The trigger condition and the chosen migration approach are recorded here and in the roadmap before implementation begins

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
