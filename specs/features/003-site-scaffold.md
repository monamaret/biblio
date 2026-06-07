# Feature: Site scaffold — React + Vite + Worker combined project template

**Status:** Not started
**Roadmap reference:** [roadmap.md §3 — Site scaffold](../research/roadmap.md#3-site-scaffold--react--vite--worker-combined-project-template)

## Summary

The project template `biblio-cli new` generates: a React + Vite site using React Router v7's native prerendering, styled with Tailwind + shadcn/ui, that builds to static assets served by the same Worker that hosts `/mcp` — the static-site half of the selected architecture.

## Confirmed decisions this builds on

- [Decision Log — switch to React Router v7's native prerendering](../research/001-static-site-platform-with-mcp.md#decision-log): `prerender` configured in `react-router.config.ts`, accepting an async function that enumerates routes from the content index — no third-party SSG dependency
- [Decision Log — replace Hugo with React + Vite](../research/001-static-site-platform-with-mcp.md#decision-log): consolidates the stack onto one toolchain shared with the MCP App UI layer
- [Decision Log — author content as plain Markdown + frontmatter](../research/001-static-site-platform-with-mcp.md#decision-log): `import.meta.glob` + `gray-matter` + `react-markdown` (+ `remark-gfm`), no MDX
- [Decision Log — content directory layout convention](../research/001-static-site-platform-with-mcp.md#decision-log): the confirmed on-disk shape (`profile.md`/`experience|skills|projects/`, `entries/YYYY-MM-DD-slug.md`, `docs/<section>/<page>.md`) that `import.meta.glob` patterns and `prerender` route enumeration target directly
- [architecture.md §1 — High-level architecture (selected)](../research/architecture.md#1-high-level-architecture-selected): the Worker serves prerendered static assets via the Workers Static Assets binding alongside `/mcp`
- [Architecture Patterns — frontmatter schema per site type](../research/001-static-site-platform-with-mcp.md#architecture-patterns): the per-site-type page/content models this scaffold's templates and layouts must match

## Scope

- [ ] Set up the React + Vite + React Router v7 project skeleton, with native `prerender` configured in `react-router.config.ts`
- [ ] Wire the `prerender` async function to enumerate content-driven routes from `content-index.json` (log-entry lists, docs TOC pages, etc.)
- [ ] Implement content loading via `import.meta.glob` + `gray-matter`, rendered with `react-markdown` (+ `remark-gfm`)
- [ ] Build the per-site-type page templates/layouts:
  - [ ] Personal profile/CV — profile, experience, skills, and projects views
  - [ ] Technical project log — filterable entry list (by date/tag/project) and entry detail view
  - [ ] Project documentation — TOC/nav generated from file structure, and page view
- [ ] Integrate Tailwind CSS + shadcn/ui, consuming the generated `tokens.css`
- [ ] Wire the Workers Static Assets binding to the prerendered build output

## Out of scope

- The `content-index.json` generator that feeds this scaffold's `prerender` step — covered by [002-content-pipeline](002-content-pipeline.md)
- The MCP server sharing the Worker — covered by [004-in-worker-mcp-server](004-in-worker-mcp-server.md); this feature only owns the static-asset half of the Worker
- Theme preset content (palettes, fonts per site type) — covered by [006-theming-system](006-theming-system.md); this feature only consumes the generated `tokens.css`
- The terminal/CLI-optimized rendering branch in the Worker's `fetch` handler — covered by [007-terminal-cli-view](007-terminal-cli-view.md)

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) remain open for this feature. (The Docs Page templates need no tag/project filter UI — confirmed the schema stays `title`/`order`/`section`/`summary`, with `section`/`order`-driven TOC/nav as the docs site's organizing principle; see [001 → Decision Log](../research/001-static-site-platform-with-mcp.md#decision-log).)

## Dependencies

- [002-content-pipeline](002-content-pipeline.md) — the scaffold's `prerender` step and page templates consume `content-index.json`
- [006-theming-system](006-theming-system.md) — the scaffold imports the generated `tokens.css`

## Acceptance criteria

- `react-router build` prerenders every content-driven route to static HTML/CSS/JS without a runtime SSR server
- Each of the three site types' page templates render correctly from real content-directory input
- The prerendered output binds directly to the Workers Static Assets binding with no transformation step

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
