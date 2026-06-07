# Feature: Pilot site validation

**Status:** Not started
**Roadmap reference:** [roadmap.md §8 — Pilot site validation](../research/roadmap.md#8-pilot-site-validation)

## Summary

Build and deploy one real site end to end — validating the scaffold, CLI, and architecture together against an actual production deployment — before replicating the approach to the remaining site types.

## Confirmed decisions this builds on

- [Performance, Reliability & Operations — rollout/migration plan](../research/001-static-site-platform-with-mcp.md#performance-reliability--operations): build the scaffold + CLI, validate with one pilot site (e.g. personal profile/CV), then replicate
- [architecture.md §2 — Build & deployment pipeline](../research/architecture.md#2-build--deployment-pipeline): the exact three-step pipeline this pilot exercises end to end
- [Risks & Mitigations](../research/001-static-site-platform-with-mcp.md#risks--mitigations): "React + Vite + Worker MCP integration proves awkward to scaffold repeatably" — mitigation is explicitly to prototype with one site before generalizing into `biblio-cli`

## Scope

- [ ] Scaffold and deploy a personal profile/CV site as the pilot
- [ ] Verify the full build pipeline end to end: content-index generation → `react-router build` (prerender) → Wrangler deploy
- [ ] Verify the `/mcp` endpoint responds correctly to a real MCP client (`initialize`, `tools/list`, `tools/call`)
- [ ] Verify the MCP App UI renders correctly in at least one graphical host, and degrades correctly (text-first) in a terminal-based host
- [ ] Capture and address any scaffold/pipeline issues found, feeding fixes back into the upstream features

## Out of scope

- Building any of the underlying scaffold/CLI/MCP-server/UI/theming pieces — those are [001](001-biblio-cli-foundations.md)–[006](006-theming-system.md); this feature only validates them together against a real deployment
- Replicating to the project-log and docs site types — covered by [009-replicate-site-types](009-replicate-site-types.md), which depends on this feature succeeding first

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) directly block this feature beyond what already blocks its dependencies — its purpose is in fact to surface *new* issues that the design-stage research couldn't anticipate.

## Dependencies

- [001-biblio-cli-foundations](001-biblio-cli-foundations.md)
- [002-content-pipeline](002-content-pipeline.md)
- [003-site-scaffold](003-site-scaffold.md)
- [004-in-worker-mcp-server](004-in-worker-mcp-server.md)
- [005-mcp-app-ui-layer](005-mcp-app-ui-layer.md)
- [006-theming-system](006-theming-system.md)

(§7 — terminal/CLI view — is not a hard dependency; it can validate after the core pilot succeeds, since it's an additive surface on the same Worker.)

## Acceptance criteria

- A real personal profile/CV site is live on a custom domain, fully built and deployed via the `biblio-cli` → `react-router build` → `wrangler deploy` pipeline with no manual post-scaffold edits
- A real MCP client successfully completes `initialize` → `tools/list` → `tools/call` against the deployed `/mcp` endpoint
- The MCP App UI renders in a graphical host and degrades to text-first content in a terminal-based host, both observed directly (not just unit-tested)
- Every issue found is captured and either fixed upstream or explicitly logged as a new open item

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
