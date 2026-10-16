# Feature: `quarterly` site type

**Status:** Complete — implemented and tested on branch `monamaret/004-quarterly-site-type`; ships with v0.08
**Roadmap reference:** [roadmap.md §11 — Quarterly site type](../research/roadmap.md#11-quarterly-site-type)

## Summary

A fourth site type, `quarterly`: an issue-based magazine with articles grouped by category, set in the style of an early-1970s Eastern European academic book (b&w plates, footnotes and margin notes, one accent color, configurable typesetting, no comments). It keeps every guarantee the profile, log and docs types give. Shape and decisions come from [004-quarterly-site-type.md](../research/004-quarterly-site-type.md); the Spec Kit spec is [specs/004-quarterly-site-type/spec.md](../004-quarterly-site-type/spec.md).

## Confirmed decisions this builds on

- [003 → A fourth site type](../research/003-post-ga-feature-backlog.md): the backlog item said a fourth type needs a real driving use case; `quarterly` is that use case, promoted to research on 2026-10-15
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: add `quarterly` as an issue-based magazine
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: content layout is a folder per issue plus a `category` field on each article
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: typesetting is preset defaults plus an optional `typesetting:` block in `site.yaml`, with `tokens.css` regenerated on `dev`/`build`
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: images are converted to grayscale at build time in pure Go
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: footnotes `[^n]` and margin notes `{>> … <<}` are separate syntaxes
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: issues can be drafts (`draft: true`), hidden from everything published but shown in `dev`
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: the automatic pre-build step validates settings and regenerates the index for all four types; tokens and plates stay quarterly-only
- [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log), 2026-10-15: PDF output is a separate follow-on feature ([005](../research/005-pdf-output-quarterly.md))
- [004 → Architecture Patterns](../research/004-quarterly-site-type.md#architecture-patterns): routes, notes rendering, typesetting schema and MCP tool surface

## Scope

- [x] Scaffold a `quarterly` template under `templates/site/quarterly/`, copied from `log` and adapted, registered in `SiteTypes` and the TUI type list
- [x] Per-file-kind content schemas (issue, article, colophon) and derived index fields (`kind`, `issue`, `issue_number`, `issue_date`, `plate`), including `draft` handling
- [x] Routes: current-issue contents at `/`, `/issues`, `/issues/:issue`, `/category/:category`, article and colophon detail
- [x] Footnotes and margin notes: wide-screen sidenotes, narrow-screen foot notes and † toggles (CSS only), Markdown/ANSI rewriting, structured notes in the index
- [x] `quarterly-default` preset, the `typesetting:` schema with Go validation, and a template stylesheet that reads tokens only; self-hosted OFL fonts with their licenses
- [x] Build-time grayscale, width-capped plates with a content-hash cache and decode safety checks
- [x] Automatic pre-build step in `build`/`dev`: settings validation and `index` for all four types; tokens and plates for quarterly
- [x] MCP tools `list_issues`, `get_issue`, `list_articles`, `get_article` (structured notes), text-first
- [x] `update`/`upgrade` support for the new template tree, and the cross-type consistency review extended to four types

## Out of scope

- PDF output (issue PDFs, offprints, the A4 print edition, `biblio-cli pdf`) — follow-on feature from [005-pdf-output-quarterly.md](../research/005-pdf-output-quarterly.md); this feature must only avoid blocking it
- Comments, author pages, and a `list_categories` tool — excluded in [004 → Problem Statement](../research/004-quarterly-site-type.md#problem-statement)
- Restyling profile, log or docs to read their tokens — a separate later decision
- A plugin system for site types, or arbitrary author CSS — non-goals of 004

## Open questions

None open. The four planning items from [roadmap.md §0.5](../research/roadmap.md#05-quarterly-site-type) were resolved on 2026-10-15 and logged in [004 → Decision Log](../research/004-quarterly-site-type.md#decision-log):

- Default text face: Archivo
- Structured notes: parsed in Go at index time
- Ornament toggles: `:root` custom properties only; structural switches in the index
- Image output: keep JPEG/PNG (WebP → PNG/JPEG), 1600 px max, 40 MP cap

## Dependencies

- [009-replicate-site-types](009-replicate-site-types.md) — the three existing types, the `update`/`upgrade` mechanics and the cross-type consistency review this feature extends
- [006-theming-system](006-theming-system.md) — the token generator this feature extends with a typed schema and overrides

## Acceptance criteria

- A freshly scaffolded quarterly site builds, prerenders every issue, article, category, archive and colophon page, and deploys
- Changing a `typesetting:` value changes the rendered site on the next `dev`/`build`; an invalid value fails the build naming the key
- Every content image is served grayscale by default, and unchanged images are not reconverted
- Footnotes and margin notes render correctly on wide and narrow viewports and keep all their text in the Markdown and ANSI views
- `get_article` returns footnotes and margin notes as structured data
- Draft issues are absent from the built site, index, `llms.txt` and tools, and visible in `dev`
- Profile, log and docs sites produce unchanged pages and assets, and no longer need a manual `index` run
- The cross-type consistency review passes with all four types

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [004-quarterly-site-type.md](../research/004-quarterly-site-type.md) — research and confirmed decisions
- [005-pdf-output-quarterly.md](../research/005-pdf-output-quarterly.md) — follow-on PDF feature this must not block
- [architecture.md](../research/architecture.md) — reference diagrams
- [specs/004-quarterly-site-type/spec.md](../004-quarterly-site-type/spec.md) — Spec Kit specification
