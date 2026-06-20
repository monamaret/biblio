# Phase 0 Research: v0.03 — Make the site usable by AI agents and terminal clients

**Input**: [plan.md](./plan.md) Technical Context — the spec's own Assumptions state plainly that "the roadmap records all four blocking research spikes for the terminal/CLI view (§0.2) as resolved — this release is treated as a fully-specified execution task, not one with open design questions of its own." This document consolidates those resolved spikes into the form this plan depends on.

## Decision: In-Worker MCP server SDK

- **Decision**: `agents` (Cloudflare Agents SDK), `McpAgent` class, with the official `@modelcontextprotocol/sdk` underneath it (the `agents` package's own dependency, not chosen separately) — session-based, one Durable Object per MCP session, mounted on `/mcp` via `MyMCP.serve("/mcp", { binding: "MyMCP" })`. Requires `"compatibility_flags": ["nodejs_compat"]` in `wrangler.jsonc` alongside the `durable_objects` binding and `migrations` entry.
- **Rationale**: This decision was missing from the original Technical Context (only "MCP server SDK (in-Worker)" was named, with no concrete package) — flagged CRITICAL by `/speckit-analyze` (finding C1) and resolved via a hands-on research spike, [specs/research/002-mcp-server-sdk-selection.md](../../research/002-mcp-server-sdk-selection.md). The spike confirmed: `workers-mcp`'s own README now redirects users to `agents`/`McpAgent` and is stale on npm (pre-1.0, last published 2025-03-19) versus `agents` actively shipping; a real minimal `McpAgent` prototype run against the actual Workers runtime (`wrangler dev`) completed a genuine connect → list-tools → call-tool exchange via the real `@modelcontextprotocol/inspector` CLI; and the mandatory Durable-Object-per-session model was confirmed to be a pure coordination primitive, not new persisted data storage, so it doesn't conflict with this plan's "no new storage" Storage section below (Principle II's `content-index.json` remains the sole data source for tool responses).
- **Alternatives considered**: `workers-mcp` — rejected; superseded by `agents` per its own README. `@modelcontextprotocol/sdk` direct with a hand-rolled Streamable HTTP Fetch transport — rejected as unnecessary; `agents` already wraps it and worked cleanly in the spike, so hand-rolling the transport would have been speculative extra integration work Principle I rules out.

## Decision: In-Worker Markdown→ANSI rendering (tier ii)

- **Decision**: Use `md4x`'s `renderToAnsi()` in-Worker, wired by importing the WASM module via Wrangler's native `.wasm` import support (`import wasmModule from "md4x/build/md4x.wasm"`) and passing it explicitly to `init({ wasm: wasmModule })`.
- **Rationale**: Resolved 2026-06-07 ([roadmap §0.2](../../research/roadmap.md)) — source-level inspection showed `md4x`'s automatic `init()` only auto-detects Node.js or bundler-specific import styles, neither supported by Wrangler, but `init()` accepts a pre-instantiated `WebAssembly.Module` directly, which is exactly what Wrangler's native `.wasm` import produces. No separate build-time pipeline is required.
- **Alternatives considered**: A build-time Markdown→ANSI conversion step outside the Worker — rejected; the in-Worker path works with no extra pipeline once the manual WASM-wiring detail is known, and keeps rendering co-located with the negotiation logic that selects it.

## Decision: Cloudflare's managed "Markdown for Agents" toggle

- **Decision**: Build the hand-rolled negotiation + conversion branch (build-time-generated Markdown + `Accept`-led negotiation in the Worker's `fetch` handler) regardless of the zone's plan tier; treat Cloudflare's managed toggle as an optional zone-level bonus to flip on later *if* the zone ever sits on Pro+, never as a load-bearing dependency.
- **Rationale**: Resolved 2026-06-07 — the managed toggle does reach Worker-served responses, but is gated by the zone's (domain's) plan tier, not the Workers compute plan, and is entirely unavailable on the Free zone plan. The hand-rolled path is therefore required regardless, gives full output-quality control (no nav-chrome/footer/cookie-banner noise a generic converter can't strip), and was independently validated by another author solving the same problem the same way.
- **Alternatives considered**: Relying solely on the managed toggle — rejected outright; would leave the feature entirely broken on the Free zone plan, which `biblio` (a single-author, cost-conscious project per Principle I) cannot assume it will be above.

## Decision: `llms.txt` investment level

- **Decision**: A minimal, build-time-generated file listing `site_name`/`description` plus a flat list of key page links, derived directly from `content-index.json` — not a fuller `llms-full.txt`-style content dump.
- **Rationale**: Resolved 2026-06-07 — research found major LLM crawlers don't request `llms.txt` in meaningful volume; it's "expected convention to not be missing" rather than a proven traffic driver. A fuller dump would be ongoing investment (size budgeting, staleness risk, format-spec churn) chasing an unmaterialized benefit — exactly the kind of speculative generalization Principle I rejects. The minimal version satisfies the convention at near-zero cost and can expand later if usage data ever justifies it.
- **Alternatives considered**: A fuller `llms-full.txt` content dump — rejected per Principle I; would be ongoing maintenance investment against an unproven benefit.

## Decision: Negotiation and caching correctness

- **Decision**: Set `Vary: Accept, User-Agent` on every negotiated response so Cloudflare's edge cache serves the correct variant (HTML vs. Markdown/ANSI) to the correct client.
- **Rationale**: Named directly in [roadmap §7](../../research/roadmap.md) as a required step of the terminal/CLI view build-out — without it, a cached response negotiated for one client class could be incorrectly served to another, breaking both Scenario 3 (terminal clients get Markdown) and the browser-facing behavior proven in v0.02.
- **Alternatives considered**: Disabling caching for negotiated routes entirely — rejected; would needlessly sacrifice Cloudflare's edge-caching benefits (cost/perf, relevant under Principle I's minimal-footprint bias) when correct `Vary` handling achieves the same correctness for free.

## Decision: MCP App UI bundle approach

- **Decision**: Small, sandboxed, postMessage-bridged React + Vite bundles that an MCP host renders inline, importing the same `tokens.css` generated once in v0.01 — always paired with, never replacing, the tool's primary Markdown/plain-text `content` response.
- **Rationale**: Confirmed rendering decisions in [`001-static-site-platform-with-mcp.md`](../../research/001-static-site-platform-with-mcp.md) and named directly in [roadmap §5](../../research/roadmap.md); Principle IV makes the text-first requirement non-negotiable — "any richer SVG/HTML MCP App UI resource is layered on top strictly as a bonus, never a requirement."
- **Alternatives considered**: A heavier, standalone UI-rendering service or full embedded app per tool — rejected; conflicts with both Principle I (minimal-dependency bias) and Principle V (simplest topology), and would risk the rich view becoming load-bearing rather than additive.

## Open questions resolved

None remain — this release explicitly inherits its research posture from the already-resolved §0.2 spikes and the confirmed architecture/rendering decisions in `001-static-site-platform-with-mcp.md`; no `NEEDS CLARIFICATION` marker was needed in the Technical Context.
