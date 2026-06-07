# Agent Guide

Auto-loaded by opencode as project-level instructions (see
`opencode.json` -> `instructions`). Extend this file as the project
evolves — especially once `biblio-cli` and the site scaffold actually
exist and there are real build/test/lint commands to record.

## Project

`biblio` is a private, single-author platform for standing up React +
Vite static sites (personal profile/CV, technical project log, project
documentation) on Cloudflare Workers, each with a custom-built `/mcp`
endpoint served from the same Worker. A Go-based CLI (`biblio-cli`)
scaffolds, develops, builds, deploys, and updates these sites end to
end. It is not a product for external users — see
[specs/research/001-static-site-platform-with-mcp.md](specs/research/001-static-site-platform-with-mcp.md)
for the full research and confirmed-decision record, and
[specs/research/roadmap.md](specs/research/roadmap.md) for the task
breakdown (all design questions in its §0 are now resolved).

## Conventions

- **Site stack:** React + Vite, React Router v7 (native `prerender` for
  SSG), Tailwind CSS + shadcn/ui (on Radix UI). Content is plain
  Markdown + YAML frontmatter — *not* MDX — loaded via
  `import.meta.glob`, parsed with `gray-matter`, rendered with
  `react-markdown` + `remark-gfm`.
- **Hosting/runtime:** A single Cloudflare Worker serves both the
  static site (Workers Static Assets) and the in-Worker `/mcp` endpoint
  (via the `agents`/`workers-mcp` SDK + Durable Objects), mounted on a
  route (`example.com/mcp`), not a subdomain.
- **CLI stack:** Go + Cobra + Viper, with Charmbracelet's Bubble Tea
  (+ Lip Gloss/Glamour) for the interactive TUI — ships as a single
  static binary.
- **Content authoring:** flat-by-concept directories under `content/`
  (e.g. `profile.md` + `experience|skills|projects/`, `entries/`,
  `docs/<section>/<page>.md`, `site.yaml`) — see the confirmed layout
  convention in the Decision Log of
  [001-static-site-platform-with-mcp.md](specs/research/001-static-site-platform-with-mcp.md#decision-log).
- **No speculative abstraction:** this is a single-private-author
  project — favor the minimal-dependency path, defer
  storage/search escalation until a real site hits a documented
  threshold (see [specs/features/010-storage-escalation.md](specs/features/010-storage-escalation.md)).

## Build / Test / Lint

Not yet established — there is no application source code, build
system, or test suite in this repo yet (see "Project state" in
[CLAUDE.md](CLAUDE.md)). The intended end-to-end pipeline, once
`biblio-cli` exists, is:

```sh
biblio-cli build   # content-index generation -> react-router build (prerender) -> wrangler bundling
biblio-cli deploy  # wrangler deploy
biblio-cli dev     # vite dev + wrangler dev together
```

Update this section with real `install`/`build`/`test`/`lint`/`format`
commands as soon as the CLI and a site scaffold land — see
[specs/research/deployment-guide.md](specs/research/deployment-guide.md)
for the expected end-to-end usage.

## Subagents

None yet — `.opencode/agents/` only contains its own `README.md`
(format reference for defining project-scoped subagents). Add entries
here, one line per agent, as they're created.

## Skills

None yet — there is no `.opencode/skills/` directory. Add entries here,
one line per skill, as project-specific skills are created.

## Notes

- Treat [specs/research/001-static-site-platform-with-mcp.md](specs/research/001-static-site-platform-with-mcp.md)
  (Decision Log) as the source of truth for confirmed architecture/stack
  decisions, and [specs/research/roadmap.md](specs/research/roadmap.md)
  for current task status — both are kept current as decisions land.
- The `/mcp` endpoint is intentionally **public with no auth/access
  control** by design (low-sensitivity, read-oriented content); rely on
  Cloudflare's standard protections (rate limiting, WAF, bot
  management) rather than adding per-site auth.
- Out of scope for the whole project: productizing/selling to other
  users, supporting site generators other than React + Vite, and
  multi-tenant SaaS concerns — don't suggest features that point in
  those directions.
