# Implementation Plan: v0.03 — Make the site usable by AI agents and terminal clients

**Branch**: `103-release-0-03-agent-terminal-layer` | **Date**: 2026-06-07 | **Spec**: [spec.md](./spec.md)

**Input**: Release specification from `/specs/releases/0.03-agent-terminal-layer/spec.md` (mirrors [`specs/releases/0.03-release.md`](../0.03-release.md)), depends on [v0.02](../0.02-build-deploy-verify/plan.md), built on feature items [004](../../features/004-in-worker-mcp-server.md), [005](../../features/005-mcp-app-ui-layer.md), [007](../../features/007-terminal-cli-view.md)

**Note**: This template is filled in by the `/speckit-plan` command. See `.specify/templates/plan-template.md` for the execution workflow.

## Summary

Mount an in-Worker MCP server and a content-negotiation layer onto the *same* Worker that already serves the v0.02 pilot site — adding structured-query tools for AI assistants, small sandboxed MCP App UI bundles for graphical hosts (always paired with a self-sufficient text/Markdown rendering), and a hand-rolled `Accept`-led negotiation branch that returns clean Markdown/ANSI for terminal and `curl` clients — all gated by the single `mcp_enabled` boolean from `content/site.yaml`, so the same deployment serves human, agent, and terminal consumers alike. Technical approach: build directly on the confirmed, fully-resolved §0.2 research spikes (md4x WASM wiring, Cloudflare's zone-gated Markdown-for-Agents toggle, minimal `llms.txt`) — this release has no open design questions of its own, only execution.

## Technical Context

**Language/Version**: TypeScript (in-Worker MCP server, content-negotiation `fetch` handler, MCP App UI bundles — React + Vite)

**Primary Dependencies**: `agents` (Cloudflare Agents SDK, `McpAgent` class — wraps the official `@modelcontextprotocol/sdk`; confirmed via [research spike 002](../../research/002-mcp-server-sdk-selection.md), requires `"compatibility_flags": ["nodejs_compat"]` and a `durable_objects` binding + `migrations` entry), `md4x` (in-Worker Markdown→ANSI via `renderToAnsi()`, wired through Wrangler's native `.wasm` import per the resolved §0.2 spike — `import wasmModule from "md4x/build/md4x.wasm"` then `init({ wasm: wasmModule })`), small sandboxed React + Vite MCP App UI bundles (postMessage-bridged, importing the same `tokens.css` from v0.01), build-time-generated `llms.txt` derived from `content-index.json`

**Storage**: Unchanged — reads the existing `content-index.json` (read-only) for both MCP tool responses and terminal-rendering source content; no new *persisted data* storage, no D1 (Principle V — keyword search only, still pre-Vectorize). `McpAgent`'s required Durable Object binding (confirmed via [research spike 002](../../research/002-mcp-server-sdk-selection.md)) is a per-session *coordination* primitive, not a persisted-data surface this "no new storage" framing rules out — `content-index.json` remains the sole data source for every tool response (Principle II unchanged).

**Testing**: MCP tool contract tests (each tool's Markdown/plain-text `content` output, route mounting, content-index querying), TypeScript/React tests for negotiation logic and MCP App UI bundles — written first per Principle VI; an end-to-end check of "MCP client connects → lists tools → calls tool → gets correct structured result" against the live pilot site (spec Scenario 1, SC-001) is this release's closing integration test

**Target Platform**: The same single Cloudflare Worker deployed in v0.02 (Workers Static Assets + now also `/mcp*`), mounted on a route — not a subdomain, not a separate service (Principle V topology, completed by this release)

**Project Type**: Web application extension — adds a server-side MCP/negotiation layer and a small UI-bundle sub-project to the existing combined site/Worker project from v0.01/v0.02

**Performance Goals**: An MCP client completes a full connect → list-tools → call-tool exchange on the first attempt (SC-001); terminal/`curl` requests return clean renderings with no leftover browser chrome (SC-002); `Vary: Accept, User-Agent` set on every negotiated response so Cloudflare's cache serves the correct variant per client (per [roadmap §7](../../research/roadmap.md))

**Constraints**: Every MCP tool's primary `content` response MUST be a clean, self-sufficient Markdown/plain-text rendering — richer SVG/HTML MCP App UI is strictly a bonus layered on top, never a requirement (Principle IV, FR-003); when `mcp_enabled` is `false`, no MCP route, tool surface, or generated artifact may be present or reachable at all (FR-002, Scenario 4); the `/mcp` endpoint remains intentionally public with **no auth** — rely on Cloudflare's standard zone-level protections only (Security & Operational Posture); the hand-rolled negotiation/conversion path MUST be built regardless of zone plan tier, since Cloudflare's managed "Markdown for Agents" toggle is unavailable on the Free zone plan and is at most an optional bonus later (resolved §0.2 spike)

**Scale/Scope**: Demonstrated against the same single real pilot site from v0.02 — proving this layer on the other two site types is deferred to v0.04/v0.05 alongside their replication (spec Assumptions); tier-iii full styled terminal "mini-site" is explicitly out of scope (Principle IV)

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-check after Phase 1 design.*

- **Principle I (Single-Author, Minimal-Dependency Bias)**: PASS — `mcp_enabled` (generated in v0.01) remains the sole MCP-related config surface; no nested config objects, no per-site auth, no speculative tool surfaces beyond what each site type's content shape genuinely supports.
- **Principle II (One Pipeline, One Source of Truth)**: PASS — the MCP query layer and the terminal-rendering path both read the same `content-index.json` produced in v0.01; `llms.txt` is build-time-generated *from* that same index (no separate authoring step); MCP App UI bundles import the same `tokens.css` rather than re-deriving theme values.
- **Principle III (CLI-First, Guided Experience)**: PASS — nothing in this release introduces a capability outside the `biblio-cli` lifecycle; the MCP/negotiation layer is generated and deployed through the existing `build`/`deploy` commands, gated purely by the `mcp_enabled` flag set at scaffold time.
- **Principle IV (Text-First, Progressive Enhancement)**: GATE — this is the release that *establishes* this principle in running code: every MCP tool's `content` response MUST be self-sufficient Markdown/plain-text first, with any richer MCP App UI view strictly additive (FR-003); the terminal/CLI view ships only tiers i/ii (raw Markdown passthrough + in-Worker `md4x` ANSI rendering) — tier iii is explicitly and permanently out of this release's scope per the constitution's own wording.
- **Principle V (Simplest Topology That Can Grow)**: PASS — the `/mcp` endpoint mounts onto the *same* Worker already serving the static site (completing, not splitting, the topology named in the constitution) via `McpAgent`'s `durable_objects` binding ([research spike 002](../../research/002-mcp-server-sdk-selection.md) — a per-session coordination primitive, not a new data store); search remains keyword-based; no D1, no Vectorize.
- **Principle VI (Test-First Development, NON-NEGOTIABLE)**: GATE — MCP tool contract tests (content output, route mounting, index querying) and negotiation-logic tests MUST be written and confirmed failing before implementation; the connect→list→call integration check against the live pilot is this release's closing end-to-end test.
- **Security & Operational Posture**: GATE — `/mcp` MUST remain public with no auth/access control, relying solely on Cloudflare's zone-level protections (rate limiting, WAF, bot management); `mcp_enabled` stays the *only* MCP-related configuration surface, deciding in one place whether `wrangler.jsonc`'s `main`, `durable_objects` binding/`migrations`, and `nodejs_compat` flag are present at all — the constitution's own "(Durable Objects, etc.)" example binding, now concretely resolved via [research spike 002](../../research/002-mcp-server-sdk-selection.md) — this release MUST NOT introduce nested MCP config objects or any access-control layer (would require a fresh research pass + Decision Log entry it explicitly has not had).

No violations requiring justification — Complexity Tracking is not needed for this release.

## Project Structure

### Documentation (this feature)

```text
specs/releases/0.03-agent-terminal-layer/
├── plan.md              # This file (/speckit-plan command output)
├── research.md          # Phase 0 output (/speckit-plan command)
├── data-model.md        # Phase 1 output (/speckit-plan command)
├── quickstart.md        # Phase 1 output (/speckit-plan command)
├── contracts/           # Phase 1 output (/speckit-plan command)
└── tasks.md             # Phase 2 output (/speckit-tasks command - NOT created by /speckit-plan)
```

### Source Code (repository root)

```text
# additions to the v0.01/v0.02 combined site/Worker template (TypeScript)
templates/site/<type>/
├── app/
│   └── mcp/
│       ├── server.ts        # in-Worker MCP server: route mounting, tool registry
│       ├── tools/           # per-site-type structured-query tools (shared + type-specific)
│       └── ui/              # small sandboxed React+Vite MCP App UI bundles (postMessage-bridged,
│                            #  import the shared tokens.css; always paired with text/Markdown content)
├── app/
│   └── negotiation/
│       ├── accept.ts        # Accept/User-Agent–based negotiation, sets Vary: Accept, User-Agent
│       ├── markdown.ts      # tier-i raw Markdown passthrough
│       └── ansi.ts          # tier-ii md4x renderToAnsi() via Wrangler .wasm import + init({ wasm })
├── app/
│   └── discovery/
│       └── llms-txt.ts      # build-time llms.txt generation from content-index.json
└── wrangler.jsonc           # extended (not replaced) to mount /mcp* conditionally on mcp_enabled

test/
├── unit/
│   ├── mcp-tools/           # per-tool Markdown/plain-text content contract tests
│   └── negotiation/         # Accept/User-Agent branch + Vary-header tests
└── integration/
    └── pilot-site-mcp/      # connect → list-tools → call-tool, terminal-curl, and
                             #  mcp_enabled=false "zero reachable trace" end-to-end checks
```

**Structure Decision**: Additive-only changes to the existing combined site/Worker template from v0.01/v0.02 — the MCP server, negotiation layer, MCP App UI bundles, and `llms.txt` generation all live inside the same per-site-type template tree and deploy as part of the same single Worker (Principle V). No new top-level project, host, or service is introduced; `wrangler.jsonc` is extended to conditionally mount `/mcp*` based on the `mcp_enabled` flag generated at scaffold time, never unconditionally.

## Complexity Tracking

*No Constitution Check violations — section intentionally left empty.*
