# Research: Post-GA feature backlog

**Status:** Draft
**Date:** 2026-06-20
**Owner:** Project owner

## Purpose

`biblio` reached General Availability at v0.06 — [roadmap.md](roadmap.md) §1–§9 are
complete, §10 (storage escalation) is proven-executable-on-demand but
deliberately not adopted. The roadmap has no section past §10: there is no
tracked scope for what comes next.

This document is **not** a `research-template.md`-style deep dive into one
decided feature. It's a holding pen for candidate post-GA feature ideas —
some surfaced as known gaps during v0.01–v0.06, some new. Each entry should
stay a single paragraph: enough to evaluate whether it's worth a real
research doc. When an entry is picked up, promote it to its own doc under
`specs/research/` (copy `research-template.md`) and a feature spec under
`specs/features/` (copy `TEMPLATE.md`), then mark it Promoted below with a
link.

Nothing here is scheduled, prioritized, or committed. Per the roadmap's own
stance, `biblio` is a sole-developer, personal-scale project — these are
options to draw from when there's reason to, not a queue.

## Candidates carried over from known gaps (already flagged in-repo)

- **Tailwind/shadcn vs. plain-CSS-tokens — fix the `AGENTS.md` mismatch.**
  [roadmap.md §3](roadmap.md) notes that no site template actually depends on
  Tailwind/shadcn/Radix (theming shipped via `internal/theme`-generated
  `tokens.css` instead), but `AGENTS.md`'s Conventions section still describes
  the original Tailwind+shadcn plan. Not a feature — a documentation
  correction — but flagged here since it was explicitly deferred rather than
  fixed in v0.06 Polish.
- **Evaluate Glamour for in-terminal rendering of `biblio-cli`'s own output.**
  [roadmap.md §1](roadmap.md) — `glamour` is only an indirect transitive
  dependency via the Bubble Tea ecosystem today; never invoked. Would improve
  the CLI's own UX (e.g. rendering `--help`, generated guides) but has no
  bearing on generated sites.
- **Embedding model selection for the Vectorize escalation path.**
  [roadmap.md §0.4](roadmap.md) — provider (OpenRouter) is decided; the
  specific model (`text-embedding-3-small` vs. `-large` vs. BGE-M3 vs. Qwen3
  Embedding) is open, and only matters once/if a real site crosses the
  keyword-search-underperforming threshold documented in
  [docs/storage-escalation-guide.md](../../docs/storage-escalation-guide.md).
- **Terminal "mini-site" view (tier iii).**
  [roadmap.md §7](roadmap.md) — explicitly dropped from initial scope for
  lack of validated demand. Revisit only if real terminal-client usage shows
  tiers i/ii (raw Markdown, ANSI rendering) aren't enough.
- **Package-manager distribution tap (Homebrew, etc.) for `biblio-cli`.**
  [001-static-site-platform-with-mcp.md → Decision Log](001-static-site-platform-with-mcp.md#decision-log)
  ruled this out for GA as infrastructure for a multi-user audience the
  project doesn't have. Worth revisiting only if `biblio-cli` ever gets users
  beyond its sole author.

## New candidates (not yet flagged anywhere in-repo)

- **Wire real version info into `biblio-cli` builds, and improve `-h` output.**
  `internal/cli/root.go` hardcodes `var version = "dev"`, and `.goreleaser.yaml`'s
  `builds` block has no `ldflags` to inject the actual tag/commit at build
  time — so `biblio-cli -h`/`--version` reports `dev` even in the prebuilt
  GitHub Releases binaries the [confirmed distribution decision](001-static-site-platform-with-mcp.md#decision-log)
  relies on. Fix is mechanical (goreleaser's standard `-X main.version={{.Version}}`
  pattern), but bundle it with a pass over `-h` help text/examples across all
  commands (`new`, `dev`, `build`, `deploy`, `update`/`upgrade`, `index`)
  while touching this, since both are "polish the CLI's first impression."


- **A fourth site type.** The platform supports exactly three site types
  (profile/CV, project log, docs) with fixed per-type schemas and tool
  surfaces (see [001 → Architecture Patterns](001-static-site-platform-with-mcp.md#architecture-patterns)).
  A fourth (e.g. a link-blog/bookmarks type, or a photo/portfolio gallery
  type) would need its own content schema, MCP tool surface, theme preset,
  and scaffold templates — essentially repeating §3/§4/§6/§9 of the roadmap
  for one more type. Needs a real driving use case before it's worth a
  research doc.
- **Multi-site management from a single `biblio-cli` install.** Today each
  `biblio-cli new` produces one independent project; there's no concept of
  "my sites" across projects. Could matter if the author (or any future user)
  runs more than a couple of sites and wants a single place to check
  update/upgrade status, theme drift, or storage-escalation thresholds across
  all of them, rather than running `biblio-cli` separately in each project
  directory.
- **Analytics/observability for deployed sites.** Nothing in the architecture
  currently surfaces traffic, MCP tool-call volume, or terminal-vs-browser
  client split for a deployed site. Cloudflare Workers Observability is
  already available on the paid plan in use (see
  [000-tech-stack.md](000-tech-stack.md)) and could be wired into a
  `biblio-cli status`/`stats` command without a new paid dependency.
- **Content authoring ergonomics beyond raw Markdown files.** Today authoring
  is "edit a Markdown file in `content/`." A guided `biblio-cli add entry`/
  `add page` command (prompting for required frontmatter per the
  [confirmed schemas](001-static-site-platform-with-mcp.md#decision-log),
  validating before write) could reduce frontmatter mistakes the existing
  build-time schema validation only catches after the fact.
- **Internationalization / multi-language content.** No site type currently
  has any notion of locale. Would touch the content layout convention, the
  MCP tool surface (locale-scoped search), and theming — a cross-cutting
  change, not a small addition. No current driver; listed as a candidate only
  because it's a common "what's next" ask for static-site platforms.
- **Re-evaluate the dropped managed `llms.txt`/Markdown-for-Agents zone
  toggle.** [roadmap.md §0.2](roadmap.md) confirmed the hand-rolled tier-i/ii
  negotiation path is required regardless of zone plan, with the managed
  toggle as "an optional zone-level bonus to flip on later if the zone ever
  sits on Pro+." Worth a one-line check whenever/if a site's zone plan
  changes — not a feature to build, just a flag to revisit.

## How to use this document

- Add a new bullet under the relevant section as gaps or ideas surface —
  keep each entry to a paragraph; don't let this doc become the research
  itself.
- When picking one up: write the real research doc (`research-template.md`)
  and feature spec (`specs/features/TEMPLATE.md`), link them back here, and
  mark the entry `**Promoted →** [NNN-slug.md](NNN-slug.md)`.
- This doc has no roadmap numbering of its own (no §11) until at least one
  candidate is promoted and there's a real deliverable to track — adding an
  empty §11 to [roadmap.md](roadmap.md) ahead of that would be exactly the
  "designing for hypothetical future requirements" this project avoids
  elsewhere.

## References

- [roadmap.md](roadmap.md) — the completed §1–§10 plan this backlog picks up after
- [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) — confirmed architecture/decisions any new feature must fit within
- [000-tech-stack.md](000-tech-stack.md) — the already-paid-for platform list any new feature should draw from first
- [../../docs/releases/v0.06/README.md](../../docs/releases/v0.06/README.md) — current user guide, for what's already shipped
