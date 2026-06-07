# Feature: Theming system

**Status:** Not started
**Roadmap reference:** [roadmap.md §6 — Theming system](../research/roadmap.md#6-theming-system)

## Summary

The shared, token-driven theming system: one shadcn/ui CSS-variable contract, three per-site-type presets, and a `biblio-cli`-generated `tokens.css` that both the site build and the MCP App UI bundle import as their single source of truth for look-and-feel.

## Confirmed decisions this builds on

- [Decision Log — lock in the per-site-type theming system](../research/001-static-site-platform-with-mcp.md#decision-log): the shared CSS-variable contract (colors + `-foreground` pairs, `border`/`input`/`ring`, `radius`, plus Tailwind font/spacing extensions), the three presets (CV/profile: warmer/personal; project log: cooler/technical with monospace accents; docs: neutral/high-contrast), and the single generated `tokens.css` as shared source of truth
- [Architecture Patterns — per-site-type theming for MCP App UIs](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the full rationale, plus the *Retro/terminal CSS research* findings — `terminal.css`/WebTUI/`system.css`/`tui.css` surveyed and **not** adopted as parallel frameworks (collision risk with Tailwind/shadcn); instead their palettes (Gruvbox/Catppuccin/Nord/Dracula), monospace type choices, and box-drawing/CRT effects should be mined directly into the project-log and docs presets; WebTUI named as the strongest integration candidate if a parallel framework is ever wanted after all

## Scope

- [ ] Define the shared shadcn/ui CSS-variable token contract (colors + `-foreground` pairs, `border`/`input`/`ring`, `radius`, plus Tailwind font/spacing extensions)
- [ ] Build the three per-site-type presets — CV/profile (warmer/personal), project log (cooler/technical, monospace accents), docs (neutral/high-contrast)
- [ ] Mine `terminal.css` / WebTUI / `system.css` / `tui.css` for palette values (Gruvbox, Catppuccin, Nord, Dracula), monospace type choices, and box-drawing/CRT-style effects to bake into the project-log and docs presets
- [ ] Implement `biblio-cli`'s `tokens.css` generation from a chosen preset — one file, imported by both the site build and the MCP App UI bundle

## Out of scope

- Adopting any retro/terminal CSS framework wholesale (e.g. WebTUI) as a parallel styling system — explicitly rejected; only their palette/type/effect *values* are mined in
- Consuming `tokens.css` in the site build or the MCP App UI bundle — covered by [003-site-scaffold](003-site-scaffold.md) and [005-mcp-app-ui-layer](005-mcp-app-ui-layer.md) respectively; this feature only defines and generates the file
- Terminal/CLI-view-specific styling — [007-terminal-cli-view](007-terminal-cli-view.md)'s rendering tiers (i/ii) don't consume this feature's retro/terminal token research; the full styled "mini-site" view (tier iii) that would have reused it was confirmed dropped from the initial scope (see [001 → Decision Log](../research/001-static-site-platform-with-mcp.md#decision-log))

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) directly block this feature — the contract, presets, and generation approach are all confirmed decisions; this is largely an execution task.

## Dependencies

- [001-biblio-cli-foundations](001-biblio-cli-foundations.md) — `tokens.css` generation is a CLI scaffolding step

## Acceptance criteria

- Each of the three presets compiles to a valid `tokens.css` satisfying the shared CSS-variable contract
- The same generated `tokens.css`, imported into both the site build and an MCP App UI bundle, produces visually consistent theming in both
- No retro/terminal framework appears as a runtime dependency; only mined palette/type/effect values appear in the preset definitions

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
