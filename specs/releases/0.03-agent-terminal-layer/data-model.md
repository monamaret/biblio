# Phase 1 Data Model: v0.03 — Make the site usable by AI agents and terminal clients

**Input**: [spec.md](./spec.md) FR-001..FR-005, [research.md](./research.md) MCP/negotiation decisions

This release adds a *query and rendering* layer over the existing, unchanged `content-index.json` (Principle II) — it introduces no new persisted storage. The entities below are the runtime/response shapes this layer must produce and the gating rules that govern their presence.

## Entity: MCP Tool

A structured-query capability exposed by the in-Worker MCP server when `mcp_enabled` is `true` (FR-001).

| Field | Type | Notes |
|---|---|---|
| `name` / `description` | string | Discoverable via MCP `list-tools` (spec Scenario 1) |
| `scope` | enum: `shared` \| `type-specific` | Shared tools work across all site types; type-specific tools match the deployed site's `site-type` (e.g. profile vs. log vs. docs queries) |
| `content` | Markdown/plain-text | MUST be a clean, self-sufficient primary response — Principle IV, FR-003; this is what every tool returns *regardless* of whether a richer view also exists |
| `ui-resource` (optional) | MCP App UI bundle reference | Present only for tool results that include a richer inline view (research.md "MCP App UI bundle approach"); strictly additive — never replaces `content` |
| `query-source` | reference | Reads `content-index.json` read-only — same index produced in v0.01, unchanged (Principle II) |

**Validation rules**: every tool registered MUST produce a valid `content` response on its own, independent of host rendering capability (FR-003); a tool call for a capability the site type doesn't define MUST be rejected clearly (spec Edge Cases) — this is the behavior the failing-test-first MCP tool contract tests must encode (Principle VI).

**Relationships**: registered into exactly one MCP Server Mount (below); queries exactly one Content Index (read-only, from v0.01).

## Entity: MCP Server Mount

The conditional `/mcp*` route surface on the deployed Worker (FR-001, FR-002, Scenario 4).

| Field | Type | Notes |
|---|---|---|
| `mounted` | derived boolean | Equals `content/site.yaml`'s `mcp_enabled` value at build time — the *sole* gating input (Security & Operational Posture; this plan's Constitution Check Security gate) |
| `route` | string | `/mcp*` on the *same* Worker already serving the static site (Principle V topology) — never a subdomain or separate service |
| `tools[]` | array of MCP Tool | Populated only when `mounted` is `true` |

**Validation rules / state transitions**: when `mcp_enabled` is `false`, **no** route, tool surface, or generated artifact related to MCP may be present or reachable — verified by direct request to the endpoint returning no trace whatsoever (FR-002, Scenario 4, SC-003). This is a build-time decision baked into the deployed artifact, not a runtime toggle — there is no "disabled but present" intermediate state.

**Relationships**: one per Deployment (the Deployment entity carried over unchanged from [v0.02's data-model](../0.02-build-deploy-verify/data-model.md)); owns zero-or-many MCP Tools depending on its `mounted` state.

## Entity: Negotiated Page Response

The content-negotiated response a page request produces, depending on the requesting client's signaled preference (FR-004, Scenario 3).

| Field | Type | Notes |
|---|---|---|
| `variant` | enum: `html` \| `markdown` \| `ansi` | Selected by `Accept`/`User-Agent`-based negotiation (research.md "Negotiation and caching correctness") |
| `content` | string | The rendered body for the selected variant — tier-i raw Markdown passthrough or tier-ii in-Worker `md4x` `renderToAnsi()` output for terminal-capable clients; standard HTML for browsers (carried over unchanged from v0.02) |
| `vary-header` | `Vary: Accept, User-Agent` | MUST be set on every negotiated response so Cloudflare's edge cache serves the correct variant to the correct client class |

**Validation rules**: ambiguous/conflicting `Accept` signals MUST resolve via documented, deterministic negotiation logic (spec Edge Cases) — never an arbitrary or inconsistent choice; terminal/`curl` requests MUST receive a clean rendering with no leftover browser-page chrome (SC-002).

**Relationships**: produced per page request against the same Pre-rendered Page Set from v0.02 (no new content source); exactly one variant is returned per request (never multiple, never none).

## Entity: Discovery Artifact (`llms.txt`)

The minimal, build-time-generated machine-readable summary published for agent/crawler consumption (FR-005).

| Field | Type | Notes |
|---|---|---|
| `site_name` / `description` | string | Sourced directly from `content-index.json` / `content/site.yaml` |
| `key-pages[]` | flat list of links | Derived directly from `content-index.json` — no separate authoring step (research.md "`llms.txt` investment level": minimal, not `llms-full.txt`) |

**Validation rules**: generated at build time from the existing index only — introducing any new authored-content shape or per-page summary field for this artifact would violate the confirmed "minimal" decision and Principle I.

**Relationships**: generated 1:1 per Deployment, derived solely from the existing Content Index — adds no new authoring surface.
