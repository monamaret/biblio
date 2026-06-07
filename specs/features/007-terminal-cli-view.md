# Feature: Terminal/CLI-optimized view of the site

**Status:** Not started
**Roadmap reference:** [roadmap.md §7 — Terminal/CLI-optimized view of the site](../research/roadmap.md#7-terminalcli-optimized-view-of-the-site)

## Summary

A terminal-optimized rendering of site pages — served via content negotiation from the same Worker `fetch` handler — for `curl`/`lynx`/terminal-coding-agent clients that can't or won't run JS, mirroring the text-first MCP rule applied to the site's own routes.

## Confirmed decisions this builds on

- [Architecture Patterns — Terminal/CLI-optimized view of the site (Researched)](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the full confirmed direction —
  - **Mechanism:** layered `Accept`-led content negotiation (parse q-values, rank `text/markdown` vs. `text/html`, fall back to User-Agent only when `Accept` is uninformative, use `sec-fetch-dest: document` as a browser safeguard, set `Vary: Accept, User-Agent`), complemented by a `.md` URL-sibling convention and a minimal `llms.txt`
  - **Validation:** Cloudflare's "Markdown for Agents" edge feature, Vercel's content-negotiation rewrites (99.37% size reduction reported), the `wttr.in`/`cheat.sh` precedent, and Bun's reported 10× token-usage reduction serving raw Markdown to Claude Code
  - **Rendering spectrum:** raw Markdown passthrough → lightly-formatted → full terminal-styled "mini-site"; **`md4x`** (a ~100 KB gzipped WASM Markdown renderer with `renderToAnsi()`) named as the strongest candidate for collapsing the first two tiers into in-Worker, request-time rendering
- [architecture.md §3 — Request routing call flow](../research/architecture.md#3-request-routing-call-flow-in-worker-fetch-handler): the middle branch of the `fetch` handler's three-way split (`/mcp*` / markdown-negotiated / default HTML) this feature implements
- [Decision Log — confirm `md4x`'s `renderToAnsi()` runs in the Cloudflare Workers runtime](../research/001-static-site-platform-with-mcp.md#decision-log): verified viable via source-level inspection — wire the WASM module manually with `import wasmModule from "md4x/build/md4x.wasm"` (Wrangler's native `.wasm` import) and `init({ wasm: wasmModule })`, since `md4x`'s automatic `init()` doesn't detect the Workers runtime
- [Decision Log — Cloudflare's "Markdown for Agents" reach and plan-tier gating](../research/001-static-site-platform-with-mcp.md#decision-log): the managed feature *does* reach Worker-served responses, but is gated by the **zone's plan tier** (Pro/Business/Enterprise — unavailable on Free, and a different billing construct than the Workers compute plan) plus `Content-Length`/~1–2 MB constraints — so **build the hand-rolled negotiation+conversion path below regardless**; it's the only thing that works on a Free zone plan, and gives strictly better output quality (no nav-chrome/footer/cookie-banner noise) than the generic converter even where Pro+ is available. The managed toggle is an optional zero-cost zone-level bonus to flip on later if the zone ever sits on Pro+, never a dependency
- [Decision Log — scope `llms.txt` minimal](../research/001-static-site-platform-with-mcp.md#decision-log): `site_name`/`description` plus a flat list of key page links derived directly from `content-index.json`, no extra authoring step — not a fuller `llms-full.txt`-style dump, since research found major LLM crawlers don't request it in meaningful volume yet
- [Decision Log — drop the tier-iii "mini-site" view from initial scope](../research/001-static-site-platform-with-mcp.md#decision-log): tiers i/ii already cover the actual need; a full styled terminal view would be speculative investment against unvalidated demand — revisit only if real terminal-client usage data ever shows it's wanted

## Scope

- [ ] Implement layered content negotiation in the Worker's `fetch` handler: parse `Accept` q-values (rank `text/markdown` vs. `text/html` properly, don't substring-match), fall back to User-Agent matching only when `Accept` is uninformative, and use `sec-fetch-dest: document` as a browser safeguard
- [ ] Set `Vary: Accept, User-Agent` on every negotiated response so Cloudflare's cache serves the correct variant to the correct client
- [ ] Add `.md` URL-sibling routes and `Link: rel="alternate"; type="text/markdown"` discoverability headers
- [ ] Generate a minimal `llms.txt` from the content index — confirmed scope: `site_name`/`description` plus a flat list of key page links, no fuller `llms-full.txt`-style dump (see Confirmed decisions above)
- [ ] Implement raw Markdown passthrough rendering (tier i — the cheapest, always-available baseline)
- [ ] Implement in-Worker lightly-formatted Markdown → ANSI rendering via `md4x`'s `renderToAnsi()` (tier ii) — confirmed viable, with no separate build-time pipeline; wire the WASM module manually via `init({ wasm: wasmModule })` (see Confirmed decisions above)

## Out of scope

- The `/mcp*` and default-HTML branches of the `fetch` handler — covered by [004-in-worker-mcp-server](004-in-worker-mcp-server.md) and [003-site-scaffold](003-site-scaffold.md) respectively; this feature owns only the markdown-negotiated branch
- A full terminal-styled "mini-site" view (rendering tier iii) — confirmed dropped from the initial scope entirely (see Confirmed decisions above); tiers i/ii below are the complete deliverable for now
- Any MCP-tool-result rendering (charts, text-first tool `content`) — that's a parallel but distinct concern owned by [005-mcp-app-ui-layer](005-mcp-app-ui-layer.md)

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) remain open for this feature — all four terminal/CLI-view spike items (`md4x` Workers-runtime viability, "Markdown for Agents" reach, `llms.txt` investment level, and the tier-iii mini-site decision) are now resolved (see Confirmed decisions above). This is now an execution task against a fully-specified two-tier (i/ii) scope.

## Dependencies

- [003-site-scaffold](003-site-scaffold.md) — this feature extends the same Worker `fetch` handler and operates on the same prerendered/Markdown content

## Acceptance criteria

- `curl -H "Accept: text/markdown" https://<domain>/<page>` returns a 200 Markdown response with `Vary: Accept, User-Agent` set, while a normal browser request to the same URL still returns the full HTML site
- `.md` URL siblings and `Link: rel="alternate"` headers resolve correctly and are discoverable
- `llms.txt` is generated from the content index and reflects real site content
- Cloudflare's cache correctly serves the HTML variant to browsers and the Markdown variant to negotiated clients for the same URL (no cross-contamination)

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
