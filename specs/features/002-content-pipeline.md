# Feature: Content pipeline — `content-index.json` generator

**Status:** Not started
**Roadmap reference:** [roadmap.md §2 — Content pipeline](../research/roadmap.md#2-content-pipeline--content-indexjson-generator)

## Summary

The Go-side generator that walks a site's content directory and produces `content-index.json` — the single canonical, derived artifact consumed by both the prerendered site build and the in-Worker MCP server.

## Confirmed decisions this builds on

- [Decision Log — content directory layout convention](../research/001-static-site-platform-with-mcp.md#decision-log): the exact on-disk shape this walker targets — flat-by-concept directories under `content/` (`profile.md`/`experience|skills|projects/`, `entries/YYYY-MM-DD-slug.md`, `docs/<section>/<page>.md`, `site.yaml`), with slug derivation matching the confirmed rules (filename for log entries, file location for docs pages)
- [Architecture Patterns — progressive-enhancement storage approach](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the bundled JSON content index is the default storage tier, generated at build time, "one source of truth, one pipeline, always in sync"
- [Architecture Patterns — frontmatter schema per site type](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the fixed schemas this generator must parse — Profile/Experience/Skill/Project, Log entry, Docs page, plus the shared site-level config
- [Decision Log — per-site-type frontmatter schemas and storage-escalation thresholds](../research/001-static-site-platform-with-mcp.md#decision-log): the exact field lists locked in for each content type, and the heuristic thresholds for escalating beyond JSON
- [Decision Log — build pipeline integration](../research/001-static-site-platform-with-mcp.md#decision-log): content-index generation is **step 1**, owned by `biblio-cli` as a Go-side generator (mirroring what `gray-matter` would do in JS), run before the site build
- [architecture.md §2 — Build & deployment pipeline](../research/architecture.md#2-build--deployment-pipeline): diagram of where this generator sits in the pipeline and how its output feeds both `react-router build` and the Worker bundle

## Scope

- [ ] Implement the content-directory walker
- [ ] Implement frontmatter + body parsing per site type, matching the confirmed schemas (Profile/Experience/Skill/Project, Log entry, Docs page, shared site-level config)
- [ ] Implement slug/path derivation rules — from filename for log entries, from file location for docs pages
- [ ] Emit `content-index.json` in a shape consumable both by React Router's `prerender` function and the in-Worker MCP query layer
- [ ] Add schema validation for authored frontmatter (Starlight's Zod-validated frontmatter flagged as a useful design-pattern reference)

## Out of scope

- Consuming `content-index.json` in the site build (`prerender`, page templates) — covered by [003-site-scaffold](003-site-scaffold.md)
- Consuming `content-index.json` in the MCP server's query layer — covered by [004-in-worker-mcp-server](004-in-worker-mcp-server.md)
- The D1/Vectorize escalation paths that would eventually supersede or augment this JSON index — covered by [010-storage-escalation](010-storage-escalation.md), and explicitly not built up front

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) remain open for this feature — the layout convention, MCP-flag shape, and docs-schema scope (no `tags`/`project` on Docs Pages, confirmed unchanged) are all resolved. This is now an execution task against fully-specified schemas.

## Dependencies

- [001-biblio-cli-foundations](001-biblio-cli-foundations.md) — this generator is a subcommand/internal package of the CLI binary

## Acceptance criteria

- Running the generator against a populated content directory for each of the three site types produces a valid `content-index.json` matching the confirmed schemas
- Malformed frontmatter is caught by schema validation with a clear error, not a silent bad index
- The same `content-index.json` is consumed without transformation by both the `prerender` step and the Worker's MCP query layer

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
