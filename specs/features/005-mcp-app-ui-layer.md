# Feature: MCP App UI layer

**Status:** Not started
**Roadmap reference:** [roadmap.md §5 — MCP App UI layer](../research/roadmap.md#5-mcp-app-ui-layer)

## Summary

The small, sandboxed React + Vite UI-resource bundles that an MCP host can render inline alongside a tool's text-first result — themed automatically via the site's shared `tokens.css`, with hand-rolled SVG/CSS charts instead of a charting-library dependency.

## Confirmed decisions this builds on

- [Decision Log — confirm React for MCP Apps](../research/001-static-site-platform-with-mcp.md#decision-log): React is a required dependency for the MCP server's UI layer, paired with Vite + Tailwind + shadcn/ui
- [Decision Log — no charting-library dependency + text-first rule](../research/001-static-site-platform-with-mcp.md#decision-log): hand-roll minimal SVG/CSS chart components styled via `tokens.css`; every chart-like tool result must produce a self-sufficient Markdown/ASCII rendering as its primary `content`, with the UI resource layered on top only for hosts that support it
- [Decision Log — per-site-type theming system](../research/001-static-site-platform-with-mcp.md#decision-log): one generated `tokens.css` per site, imported by *both* the site build and the MCP App UI bundle — single source of truth for the look
- [Architecture Patterns — MCP App UI rendering](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the rationale for both of the above, including why terminal-based MCP hosts can't render HTML/React UI resources at all
- [architecture.md §4 — MCP tool-call flow](../research/architecture.md#4-mcp-tool-call-flow-eg-list_entries-on-a-project-log-site): the sequence showing the optional `ui://` resource path and the postMessage-bridged render

## Scope

- [ ] Set up the MCP App UI bundle build (small, sandboxed, postMessage-bridged React + Vite bundles)
- [ ] Wire the shared `tokens.css` into the UI bundle build so it inherits the site's theme automatically
- [ ] Hand-roll minimal SVG/CSS chart components (timelines, skill-level bars, simple bar/line charts) styled via `tokens.css` — no charting-library dependency
- [ ] Implement the text-first, progressive-enhancement rule: every chart-like tool result produces a clean, self-sufficient Markdown/ASCII rendering (table, sparkline) as its primary `content`, with the richer UI resource layered on top only for hosts that support it

## Out of scope

- The MCP server/tool implementations that decide *when* to attach a UI resource — covered by [004-in-worker-mcp-server](004-in-worker-mcp-server.md); this feature owns the bundle that gets rendered, not the tool logic that references it
- Generating `tokens.css` itself — covered by [006-theming-system](006-theming-system.md); this feature only consumes it
- The terminal/CLI-optimized rendering of the *site itself* (a related but separate "text-first for non-graphical clients" concern) — covered by [007-terminal-cli-view](007-terminal-cli-view.md)

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) directly block this feature — its shape is fully determined by confirmed decisions above.

## Dependencies

- [006-theming-system](006-theming-system.md) — the bundle build needs a generated `tokens.css` to consume
- [004-in-worker-mcp-server](004-in-worker-mcp-server.md) — needs at least one tool that can reference a `ui://` resource to validate the rendering path end to end

## Acceptance criteria

- A chart-like tool result renders correctly as a standalone Markdown/ASCII table or sparkline with no UI resource present (the terminal-host baseline)
- The same tool result, in a graphical host that supports MCP App UI resources, additionally renders the richer SVG/CSS chart, themed to match the site's `tokens.css`
- No charting-library dependency appears in the bundle's `package.json`

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
