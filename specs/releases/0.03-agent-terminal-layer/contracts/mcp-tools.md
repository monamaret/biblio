# Contracts: In-Worker MCP server (v0.03)

**Type of contract**: MCP tool/endpoint contract — the structured-query interface this release exposes to AI assistants and MCP hosts (FR-001). Mounted on `/mcp*` of the same Worker deployed in v0.02, gated entirely by `mcp_enabled`.

## Endpoint mounting

| Condition | Behavior |
|---|---|
| `content/site.yaml` → `mcp_enabled: true` | `/mcp*` is mounted on the deployed Worker; tool registry populated per site type |
| `content/site.yaml` → `mcp_enabled: false` | **No** `/mcp*` route, tool surface, or related generated artifact is present or reachable — verified by direct request returning no trace (FR-002, Scenario 4, SC-003) |

This is a **build-time** decision baked into the deployed artifact (Security & Operational Posture) — there is no runtime feature flag and no intermediate "present but disabled" state.

## Tool contract (every tool, shared and type-specific)

| Aspect | Contract |
|---|---|
| Discovery | Listed via standard MCP `list-tools`; client can list and successfully call tools appropriate to the deployed site's type (Scenario 1) |
| Primary response (`content`) | MUST be a clean, self-sufficient Markdown/plain-text rendering — a complete answer on its own, with no dependency on richer rendering (Principle IV, FR-003) |
| Optional richer view (`ui-resource`) | When present, renders inline in graphical MCP hosts, inherits the site's `tokens.css` theme, and is **always** accompanied by the same self-sufficient `content` (Scenario 2) — never a substitute for it |
| Unsupported capability | A request for a capability the site type doesn't define MUST be rejected with a clear, structured error — not a silent empty result or a crash (spec Edge Cases) |
| Auth | None — `/mcp` is intentionally public, relying on Cloudflare's zone-level protections only (Security & Operational Posture) |

## Negotiated page-rendering contract

| Client signal | Response |
|---|---|
| Browser-preferring `Accept`/`User-Agent` | Standard HTML (carried over unchanged from v0.02) |
| Terminal/Markdown-preferring signal, ANSI-capable | Tier-ii: in-Worker `md4x` `renderToAnsi()` rendering (FR-004) |
| Terminal/Markdown-preferring signal, plain | Tier-i: raw Markdown passthrough |
| Ambiguous/conflicting `Accept` | Resolved via documented, deterministic negotiation logic — never an arbitrary choice (spec Edge Cases) |
| Every negotiated response | `Vary: Accept, User-Agent` set, so Cloudflare's edge cache serves the correct variant to the correct client (research.md) |

## Discovery artifacts

| Artifact | Contract |
|---|---|
| Per-page Markdown-friendly form | Every page is discoverable in its Markdown-friendly form via sibling addresses and/or discoverability hints (FR-005) |
| `llms.txt` | Minimal, build-time-generated from `content-index.json` — `site_name`/`description` + flat key-page links only (research.md "`llms.txt` investment level"); not a fuller content dump |

## Failure-mode contract

| Scenario | Contract |
|---|---|
| Page/tool target doesn't exist | Clear, structured "not found" — not a silent empty/200 response (spec Edge Cases) |
| Graphical host can't load theme tokens for a `ui-resource` | Falls back to a usable, if unthemed, view — never a broken/blank render (spec Edge Cases) |
