# Feature: In-Worker MCP server

**Status:** Not started
**Roadmap reference:** [roadmap.md §4 — In-Worker MCP server](../research/roadmap.md#4-in-worker-mcp-server)

## Summary

The custom-built MCP server, implemented directly inside the same Cloudflare Worker that serves the static site, mounted on `/mcp` with Durable Objects for session state and an in-memory query layer over `content-index.json`.

## Confirmed decisions this builds on

- [Architecture Patterns — selected pattern](../research/001-static-site-platform-with-mcp.md#architecture-patterns): a single Worker serves both the site and `/mcp` in-Worker via the `agents`/`workers-mcp` SDK with Durable Objects, mounted on a route rather than a subdomain
- [Decision Log — standardize on the in-Worker architecture](../research/001-static-site-platform-with-mcp.md#decision-log): Pattern 1 selected; GCP Cloud Run patterns dropped
- [Decision Log — route over subdomain](../research/001-static-site-platform-with-mcp.md#decision-log): `example.com/mcp`, no extra DNS, no CORS concerns (MCP clients connect directly rather than via browser JS)
- [Decision Log — public, no-auth `/mcp` with build-time opt-in/out](../research/001-static-site-platform-with-mcp.md#decision-log): no per-site access control; Cloudflare's standard protections are the operational safety net; a build-time flag gates whether the index/route/bindings are generated at all
- [Decision Log — MCP opt-in/out configuration shape](../research/001-static-site-platform-with-mcp.md#decision-log): the gate is a single flat `mcp_enabled` boolean (default `true`) read from `content/site.yaml` — exactly what the gating logic checks
- [Architecture Patterns — tool surface per site type](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the fixed tool list per site type, plus shared `get_site_info`/`search`
- [architecture.md §3 — Request routing call flow](../research/architecture.md#3-request-routing-call-flow-in-worker-fetch-handler) and [§4 — MCP tool-call flow](../research/architecture.md#4-mcp-tool-call-flow-eg-list_entries-on-a-project-log-site): the in-Worker routing branch and the `tools/call` sequence this server implements

## Scope

- [ ] Set up the `agents`/`workers-mcp` SDK in the Worker, with Durable Objects for session state
- [ ] Implement in-Worker path routing — `/mcp*` to the MCP server, fall through to static assets otherwise
- [ ] Implement the `content-index.json` import and an in-memory query layer (filter/lookup/search)
- [ ] Implement the shared tools: `get_site_info` and the generic `search`
- [ ] Implement the per-site-type tool surface:
  - [ ] Personal profile/CV — `get_profile`, `list_experience`, `list_skills`, `list_projects`, `get_resume`
  - [ ] Technical project log — `list_entries`, `get_entry`, `search_entries`
  - [ ] Project documentation — `get_toc`/`list_pages`, `get_page`, `search_docs`
- [ ] Implement the build-time MCP opt-in/out gating driven by `mcp_enabled` — conditional content-index generation, route mounting, and bindings
- [ ] Confirm the `/mcp` endpoint is reachable with no auth, relying on Cloudflare's standard protections (rate limiting, WAF, bot management)

## Out of scope

- The MCP App UI resources these tools may optionally return — covered by [005-mcp-app-ui-layer](005-mcp-app-ui-layer.md); this feature owns the text-first `content` baseline and the optional `ui://` resource reference, not the bundle that renders it
- The `content-index.json` artifact itself — covered by [002-content-pipeline](002-content-pipeline.md); this feature only queries it
- D1/Vectorize-backed query implementations — covered by [010-storage-escalation](010-storage-escalation.md); this feature's query layer targets the in-memory JSON index only

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) remain open for this feature. (`search_docs`/`get_toc` need no tag/project filter parameters — confirmed the Docs Page schema stays `title`/`order`/`section`/`summary`, with keyword search plus `section`/`order`/path scoping covering the docs site type's needs; see [001 → Decision Log](../research/001-static-site-platform-with-mcp.md#decision-log).)

## Dependencies

- [002-content-pipeline](002-content-pipeline.md) — the query layer is built directly on `content-index.json`
- [001-biblio-cli-foundations](001-biblio-cli-foundations.md) — the opt-in/out gating is wired through CLI-generated config

## Acceptance criteria

- A real MCP client can `initialize`, `tools/list`, and `tools/call` against `/mcp` on a deployed site and get correct results
- Every tool's primary `content` is text-first (Markdown/plain text), per the confirmed text-first rule, independent of whether a UI resource is also returned
- A site with the MCP opt-out flag set produces a Worker with no `/mcp` route, no content-index generation step, and no MCP-related bindings

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
