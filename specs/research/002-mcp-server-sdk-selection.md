# Research: In-Worker MCP server SDK selection

**Status:** Complete
**Date:** 2026-06-19
**Owner:** monamaret

## Problem Statement

`/speckit-analyze` on v0.03 (`specs/releases/0.03-agent-terminal-layer/`) flagged a CRITICAL gap (finding C1): research.md for that release claims "no open questions remain... this release is treated as a fully-specified execution task," but no artifact anywhere — `research.md`, `plan.md`'s Technical Context, or `tasks.md` T001 — names a concrete in-Worker MCP server library. `plan.md` only says "MCP server SDK (in-Worker)"; T001 says "Add the in-Worker MCP server SDK" with no package. This blocks T001 and, transitively, all of v0.03 User Story 1 (the MVP).

- **Goal:** Name one concrete library (with version) that implements the in-Worker `/mcp*` server for v0.03 FR-001, with a recorded Decision Log rationale, so T001 is actually executable.
- **In scope:** The MCP server/transport library choice itself, and whether it requires a Durable Object for session state.
- **Out of scope:** The tool business logic (shared/profile-specific tools — already scoped in v0.03's `tasks.md` T011/T012), the MCP App UI bundle mechanics (already resolved in v0.03's `research.md` "MCP App UI bundle approach"), and the content-negotiation/ANSI-rendering layer (separately resolved, also in v0.03's `research.md`).
- **Success criteria:** A Decision Log entry naming the library + version, confirming it satisfies `contracts/mcp-tools.md`'s tool/endpoint contract (list-tools, call-tool, no auth) and Constitution Principles I (minimal-dependency), II (one content-index source of truth), III (stays behind `biblio-cli build`/`deploy`, no new CLI command), and V (same single Worker, simplest topology that can grow).
- **Non-goals:** Picking the embedding-model/Vectorize question (§0.4 of `roadmap.md` — unrelated, deferred) or revisiting the already-resolved md4x/llms.txt/MCP-App-UI decisions in v0.03's `research.md`.

## Background & Context

- **Relevant prior work:** `AGENTS.md`'s "Hosting/runtime" convention bullet already names a candidate pair in passing — "the in-Worker `/mcp` endpoint (via the `agents`/`workers-mcp` SDK + Durable Objects)" — but this was never carried into v0.03's own `research.md` as a resolved, rationale-backed decision. The constitution's Security & Operational Posture section independently anticipates Durable Objects as a possible MCP-related binding: `biblio-cli` "MUST provide the build-time `mcp_enabled` boolean... deciding in one place whether to... include MCP-related bindings (Durable Objects, etc.) in `wrangler.jsonc`." Both documents point the same direction without ever closing the loop with an actual choice.
- **Domain assumptions:** `/mcp` is intentionally public with no auth (Security & Operational Posture) — this removes session-auth complexity from the transport choice. The deployed site is a single Cloudflare Worker serving both Workers Static Assets and (now) `/mcp*` (Principle V) — whatever is chosen must mount onto that same Worker, not a separate service.
- **Stakeholders:** Single author (monamaret) — this is a `biblio`-internal infra decision, not a chosen-for-others API.

## Libraries & Stack

| Candidate | Purpose | License | Maturity | Maintenance | Fit | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| `agents` (Cloudflare Agents SDK, `McpAgent` class) | Full MCP server framework purpose-built for Workers, with built-in tool registry, Streamable HTTP + legacy SSE transport, and Durable-Object-backed session state | MIT | Actively developed, Cloudflare-maintained, the framework Cloudflare's own docs now point developers to for "build an MCP server on Workers" | High — first-party, frequent releases | Likely best fit — purpose-built for exactly this deployment shape | Requires a Durable Object binding for session state per agent/session; need to confirm a stateless-tool-call-only usage doesn't still force unnecessary DO overhead for `biblio`'s low-traffic, no-auth, single-pilot-site case |
| `workers-mcp` | Lighter, earlier Cloudflare-authored package for exposing Worker RPC methods as MCP tools | MIT | Predates `agents`/`McpAgent`; largely superseded by it | Lower — maintenance has shifted toward `agents` | Plausible but likely the legacy path | Worth a quick check of its current changelog/issues before ruling it out — if it's been formally deprecated in favor of `agents`, that settles this row |
| `@modelcontextprotocol/sdk` (official TypeScript MCP SDK) directly, with a hand-rolled Fetch-based Streamable HTTP transport | The framework-agnostic, spec-canonical SDK; would need biblio to wire its own Workers-compatible transport | MIT | Mature, spec-authoritative | High (maintained by the MCP org) | Most control, most integration work — exactly the kind of "build it ourselves" surface Principle I asks to avoid when a purpose-built option already exists | Only worth it if the spike finds `agents`/`workers-mcp` impose constraints `biblio` can't accept (e.g. mandatory DO costs/latency for a near-zero-traffic single-pilot site) |

- **Selected:** `agents` (Cloudflare Agents SDK), `McpAgent` class — confirmed working end-to-end against the real Workers runtime in Spike 2; see Decision Log.
- **Rejected (with reason):** `workers-mcp` — its own README now redirects users to `agents`/`McpAgent` ("You should start here instead"); also stale on npm (last published 2025-03-19 at `0.0.13`, never reached 1.0) versus `agents` shipping `0.16.2` two days before this spike. `@modelcontextprotocol/sdk` direct (hand-rolled transport) — not needed; `agents` already wraps it and worked cleanly, so the added integration work would have been speculative per Principle I.
- **Open questions:** None remain — see Decision Log for the resolved spike outcomes.

## Architecture Patterns

- **Candidate patterns:**
  1. **Session-based, Durable-Object-backed** (`agents`/`McpAgent` default) — one DO instance per MCP session/connection, holding any connection-scoped state; matches the MCP spec's stateful-session model most directly.
  2. **Stateless, per-request** — treat every `list-tools`/`call-tool` as an independent Fetch request with no session continuity, if the chosen library/transport supports it (Streamable HTTP without SSE can be largely stateless for simple tool-call request/response pairs).
- **Selected pattern(s):** (1) Session-based, Durable-Object-backed — confirmed via Spike 2 to be how `McpAgent` actually works (not optional), and fighting it would itself be the kind of added complexity Principle I warns against.
- **Tradeoffs:** A DO per session is a small, usage-scoped cost (Workers' free tier includes a DO allowance) — likely fine for a single-author, single-pilot-site project, but confirm it doesn't silently introduce a billing surprise or contradict the "no new storage" framing in v0.03's `plan.md` Technical Context (Durable Object *compute/coordination* is distinct from D1 *persisted SQL storage* — the Storage section's "no D1, no Vectorize" wording was about persisted data stores, not DOs as a coordination primitive; this distinction should be made explicit in `plan.md` once resolved, closing finding L1 from `/speckit-analyze`).
- **Integration points:** Same Worker that already serves Workers Static Assets (v0.02) and will host the negotiation layer (v0.03 US3) — `wrangler.jsonc`'s `main` entry (already planned as conditional on `mcp_enabled` per v0.03 tasks T003/T004) is the single mount point for all three.
- **Data model implications:** If a DO is required, `internal/scaffold/generator.go`'s `renderWranglerJSONC` (T004) needs a conditional Durable Object binding/migration block alongside the conditional `main`/`assets.binding` it already plans, gated by the same `mcp_enabled` flag — no new persisted data model, only a new binding.

## Existing Codebase Evaluation

Not yet applicable — no MCP server code exists in this repo (v0.01/v0.02 only cover scaffolding/build/deploy). This spike's output feeds directly into v0.03 `tasks.md` T001 before any implementation begins.

- **Relevant modules:** `internal/scaffold/generator.go` (wrangler.jsonc generation — will need the conditional DO binding noted above if `agents` is selected); `templates/site/profile/workers/entry.ts` (the not-yet-implemented Worker entry script from v0.03 T002/T006, which will host whichever library is chosen).
- **Reuse opportunities:** None yet — this is the first MCP-related code in the project.
- **Refactor candidates:** None.
- **Replacements required:** None.
- **Code health signals:** N/A (pre-implementation).
- **Conventions to follow:** Constitution Principle III — no new `biblio-cli` command; the library must be invokable purely from the generated Worker entry script, deployed via the existing `biblio-cli build`/`deploy` wrappers.

## Security & Compliance

- **Threat surface:** Unchanged from the constitution's existing posture — `/mcp` stays public, no auth, relying on Cloudflare's zone-level protections (rate limiting, WAF, bot management). Whatever library is chosen must not introduce an auth requirement or session-token storage that would conflict with this.
- **Data classification:** No new sensitive data — MCP tools query the existing, low-sensitivity `content-index.json` read-only.
- **Compliance regimes:** None apply (personal-scale, no PII processed).
- **Mitigations:** N/A beyond the existing zone-level protections.

## Performance, Reliability & Operations

- **Expected load:** Near-zero — single pilot site, single author's own usage plus whatever AI assistants/agents query it; no load-testing infrastructure implied.
- **SLOs / SLIs:** None formally set (consistent with Principle I's no-speculative-infra bias) — success is "an MCP client completes connect → list-tools → call-tool on the first attempt" (spec SC-001), not a latency/throughput target.
- **Observability:** Whatever the chosen library's default Workers logging surfaces; no new observability tooling planned.
- **Rollout / migration plan:** Ships as part of the same `biblio-cli build`/`deploy` path already in place since v0.02 — no separate rollout step.
- **Rollback plan:** `mcp_enabled: false` + rebuild/redeploy removes the entire surface (FR-002) — this is the project's existing rollback mechanism and doesn't change based on which library is chosen.

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation | Owner |
| --- | --- | --- | --- | --- |
| `agents` SDK's DO-per-session model turns out to be mandatory and meaningfully more complex/costly than a stateless approach for `biblio`'s near-zero-traffic case | Low | Medium | **Materialized as "mandatory" (confirmed in Spike 2) but not as "costly/complex"** — it's just a `durable_objects` binding + `migrations` entry in `wrangler.jsonc`, no extra code; for a single-pilot-site, near-zero-traffic project this is a non-issue. Resolved, no fallback needed. | monamaret |
| `workers-mcp` is found to be effectively unmaintained/deprecated, narrowing this to a one-candidate decision anyway | Medium | Low | Spike step 1 below checks its repo/changelog directly before further evaluation time is spent on it | monamaret |
| MCP spec/SDK version churn breaks the chosen library before v0.04/v0.05 replicate this to other site types | Low | Low | Pin an exact version in `package.json` (consistent with how v0.02 pinned `react-router@^7.18.0` rather than tracking `latest`); revisit only on a documented trigger | monamaret |

## Open Questions / Spike Plan

This is the concrete plan to resolve C1 — modeled on how v0.03's own §0.2 spikes were resolved (source-level inspection + a real, hands-on minimal deployment, not just documentation comparison):

- [X] **Spike 1 — Confirm `workers-mcp`'s maintenance status.** Check its npm page/GitHub repo for last-publish date and any explicit "superseded by `agents`" notice. Resolves whether it's even still a live candidate. — owner: monamaret — due: before Spike 2 — **Done 2026-06-19, see Decision Log.**
- [X] **Spike 2 — Minimal `agents`/`McpAgent` prototype.** Scaffold a throwaway Worker, define one trivial tool (echo), run it against the real Workers runtime, and confirm a real MCP client can connect → list-tools → call-tool against it. — owner: monamaret — due: before deciding — **Done 2026-06-19 against `wrangler dev` (real `workerd` runtime, not a mock) rather than a live public Cloudflare deploy — see Decision Log and Appendix for why that substitution is sufficient here.**
- [X] **Spike 3 — Confirm Durable Object framing against Principle V/Storage.** — owner: monamaret — due: same time as Spike 2 — **Done 2026-06-19, see Decision Log.**
- [ ] **Spike 4 (conditional) — If Spikes 1-2 disqualify both `agents` and `workers-mcp`**, prototype `@modelcontextprotocol/sdk` directly with a hand-rolled Streamable HTTP Fetch transport. — **Not needed: Spikes 1-2 confirmed `agents`/`McpAgent` works cleanly, so this fallback wasn't triggered.**

Once these resolve, record the outcome in this doc's Decision Log below, then:
1. Update v0.03's `research.md` with a new "Decision: In-Worker MCP server SDK" entry pointing back here.
2. Update v0.03's `plan.md` Technical Context "Primary Dependencies" row and Storage section (per Spike 3).
3. Update v0.03's `tasks.md` T001 (and T004, if a DO binding is needed) to name the concrete library/version and binding.

## Decision Log

- **2026-06-19** — **Selected: `agents` (Cloudflare Agents SDK), `McpAgent` class, with the official `@modelcontextprotocol/sdk` underneath it (the `agents` package's own dependency, not chosen separately).** — Rationale: Spike 1 found `workers-mcp`'s own README now opens with "You should start [here](https://developers.cloudflare.com/agents/guides/remote-mcp-server/) instead — and build a **remote** MCP server," an explicit upstream redirect to `agents`/`McpAgent`; its npm package is also stale (last published 2025-03-19, stuck pre-1.0 at `0.0.13`) versus `agents` actively shipping (`0.16.2`, last published 2026-06-17, two days before this spike). `workers-mcp` is therefore rejected outright — Spike 4's fallback to a hand-rolled `@modelcontextprotocol/sdk` transport was never triggered. — Decided by: monamaret
- **2026-06-19** — **Confirmed: `McpAgent` requires `nodejs_compat`.** A first prototype run failed to build with `Could not resolve "path"` and warnings that `node:async_hooks`/`node:diagnostics_channel`/`node:os` aren't available — `agents`' MCP support (and the official MCP SDK it wraps) uses Node built-ins. Adding `"compatibility_flags": ["nodejs_compat"]` to `wrangler.jsonc` resolved it cleanly. This must be carried into v0.03's `internal/scaffold/generator.go` `renderWranglerJSONC` change (currently planned as T004) alongside the conditional `main`/`durable_objects` wiring. — Decided by: monamaret
- **2026-06-19** — **Confirmed: `McpAgent` is session-based, one Durable Object per session (not stateless-per-request).** Verified directly: a raw `initialize` request returns an `mcp-session-id` response header and an SSE-framed (`text/event-stream`) body, and `wrangler.jsonc` requires a `durable_objects` binding + a `migrations` entry (`new_sqlite_classes`) for the agent class. This matches `agents`' own README ("McpAgent gives you a Durable Object per session with built-in state management... and both SSE and Streamable HTTP transports") and was independently corroborated by Cloudflare's own remote-MCP-server guide. — Decided by: monamaret
- **2026-06-19** — **Resolved (Spike 3): a Durable Object binding does NOT conflict with v0.03 `plan.md`'s "no new storage" framing.** That framing was about *persisted data stores* (D1, Vectorize) — Principle V's own wording already separately names "One Cloudflare Worker serves both the prerendered static site... and the `/mcp` endpoint" without ruling out DOs, and the constitution's Security & Operational Posture section independently anticipated "Durable Objects, etc." as an example MCP-related binding `mcp_enabled` must gate. A Durable Object here is purely a per-session coordination primitive (MCP protocol session state, e.g. the active stream/subscription), not a new persisted-data surface — `content-index.json` remains the sole data source for tool responses (Principle II, unchanged). `plan.md`'s Storage section should be updated to state this explicitly, closing `/speckit-analyze` finding L1.
- **2026-06-19** — **Validated end-to-end:** a real MCP client (the official `@modelcontextprotocol/inspector` CLI) completed a full connect → `tools/list` → `tools/call` exchange against the prototype running on the real Workers runtime (`wrangler dev`, which runs the actual `workerd` engine — not a Node mock), receiving correctly structured JSON-RPC results both times. This directly satisfies the technical risk behind spec SC-001 (v0.03) before any production code is written. — Decided by: monamaret

**Scope note on Spike 2's substitution:** the original plan called for deploying to "a real Cloudflare account." That step was substituted with `wrangler dev` (the real `workerd` runtime, running locally) because no live Cloudflare account/credentials are available in this environment — see the Appendix for the full reasoning. This resolves the *technical* viability question (does `McpAgent` run correctly and satisfy the MCP contract) but does not exercise actual public-internet deployment, DNS/route binding, or production auth-less exposure — those remain to be confirmed once v0.03 is actually deployed against the real pilot site (the same gap already flagged for v0.02's live-deploy/human-verification step).

## References

- [AGENTS.md](../../AGENTS.md) — "Hosting/runtime" convention bullet, the original (unresolved) mention of "`agents`/`workers-mcp` SDK + Durable Objects"
- [.specify/memory/constitution.md](../../.specify/memory/constitution.md) — Security & Operational Posture section, naming "Durable Objects, etc." as an example MCP-related binding gated by `mcp_enabled`
- [specs/releases/0.03-agent-terminal-layer/contracts/mcp-tools.md](../releases/0.03-agent-terminal-layer/contracts/mcp-tools.md) — the tool/endpoint contract this library must satisfy
- [specs/releases/0.03-agent-terminal-layer/research.md](../releases/0.03-agent-terminal-layer/research.md) — the release's other (already-resolved) decisions this one sits alongside
- Cloudflare Agents SDK documentation (`agents` package, `McpAgent` class) — for Spike 2
- `workers-mcp` package repository — for Spike 1

## Appendix

### Why `wrangler dev` instead of a live Cloudflare account deploy

The original Spike 2 plan called for a deploy to "a real Cloudflare account." No live account/credentials are available in this research environment (the same constraint already noted for v0.02's live-deploy step). `wrangler dev` was used instead — this runs the actual `workerd` engine (the same runtime Cloudflare deploys to in production), not a Node.js mock or polyfill layer, so it resolves the real technical question this spike exists to answer: does `McpAgent` actually run correctly in the Workers runtime and satisfy the MCP protocol contract. It does **not** exercise live public-internet routing, DNS, or the zone-level protections the constitution's Security & Operational Posture relies on for the no-auth `/mcp` endpoint — that remains to be confirmed once v0.03 is actually deployed against the real pilot site, alongside the same outstanding live-verification gap already flagged for v0.02.

### Prototype: minimal `McpAgent` echo server

**`wrangler.jsonc`** (the config shape `internal/scaffold/generator.go`'s T004 needs to replicate, conditioned on `mcp_enabled`):

```jsonc
{
  "name": "mcp-spike",
  "main": "src/index.ts",
  "compatibility_date": "2026-01-01",
  "compatibility_flags": ["nodejs_compat"],
  "durable_objects": {
    "bindings": [{ "name": "MyMCP", "class_name": "MyMCP" }]
  },
  "migrations": [{ "tag": "v1", "new_sqlite_classes": ["MyMCP"] }]
}
```

**`src/index.ts`:**

```typescript
import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { McpAgent } from "agents/mcp";
import { z } from "zod";

type Env = {
  MyMCP: DurableObjectNamespace;
};

export class MyMCP extends McpAgent<Env, unknown, {}> {
  server = new McpServer({ name: "biblio-spike", version: "1.0.0" });

  async init() {
    this.server.registerTool(
      "echo",
      {
        description: "Echoes the provided text back as Markdown/plain-text content",
        inputSchema: { text: z.string() },
      },
      async ({ text }: { text: string }) => {
        return { content: [{ type: "text", text: `echo: ${text}` }] };
      },
    );
  }
}

export default MyMCP.serve("/mcp", { binding: "MyMCP" });
```

**Dependencies installed:** `agents@0.16.2`, `@modelcontextprotocol/sdk@1.29.0`, `zod@4.4.3`, `wrangler@4.103.0`.

### First failed build (before `nodejs_compat`)

```
✘ [ERROR] Build failed with 1 error:
✘ [ERROR] Could not resolve "path"
    node_modules/mimetext/node_modules/mime-types/index.js:16:22:
      16 │ var extname = require('path').extname
The package "path" wasn't found on the file system but is built into node.
- Add the "nodejs_compat" compatibility flag to your project.
```

Adding `"compatibility_flags": ["nodejs_compat"]` resolved this cleanly — confirmed in the Decision Log.

### Real MCP client exchange (via `@modelcontextprotocol/inspector` CLI)

`tools/list`:

```json
{
  "tools": [
    {
      "name": "echo",
      "description": "Echoes the provided text back as Markdown/plain-text content",
      "inputSchema": {
        "type": "object",
        "properties": { "text": { "type": "string" } },
        "required": ["text"],
        "$schema": "http://json-schema.org/draft-07/schema#"
      },
      "execution": { "taskSupport": "forbidden" }
    }
  ]
}
```

`tools/call` (`echo`, `text: "hello from the spike"`):

```json
{
  "content": [{ "type": "text", "text": "echo: hello from the spike" }]
}
```

Raw `initialize` request via `curl`, confirming the Streamable HTTP/SSE transport and session header:

```
> POST /mcp HTTP/1.1
> Accept: application/json, text/event-stream
< HTTP/1.1 200 OK
< Content-Type: text/event-stream
< mcp-session-id: 68aa1159c3351ea0ad07c9c1a40e96af4bf9f2219d64203c0635e0a2e51e7e89
event: message
data: {"result":{"protocolVersion":"2025-11-25","capabilities":{"tools":{"listChanged":true}},"serverInfo":{"name":"biblio-spike","version":"1.0.0"}},"jsonrpc":"2.0","id":1}
```
