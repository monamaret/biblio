# Agent Guide

Auto-loaded by opencode as project-level instructions (see
`opencode.json` -> `instructions`). Extend this file as the project
evolves — especially once `biblio-cli` and the site scaffold actually
exist and there are real build/test/lint commands to record.

## Project

`biblio` is a private, single-author platform for standing up React +
Vite static sites (personal profile/CV, technical project log, project
documentation, and an issue-based `quarterly` magazine) on Cloudflare
Workers, each with a custom-built `/mcp` endpoint served from the same
Worker. A Go-based CLI (`biblio-cli`) scaffolds, develops, builds,
deploys, and updates these sites end to end. It is not a product for
external users — see
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
- **Test-driven development (NON-NEGOTIABLE):** write a failing test before
  writing the implementation that makes it pass — red-green-refactor, for
  every feature, fix, and refactor, in both `biblio-cli` (Go) and the
  site/MCP server (TypeScript/React + MCP tool contract tests). No change
  lands without a test that demonstrates the gap it closes.
- **No speculative abstraction:** this is a single-private-author
  project — favor the minimal-dependency path, defer
  storage/search escalation until a real site hits a documented
  threshold (see [specs/features/010-storage-escalation.md](specs/features/010-storage-escalation.md)).

## Build / Test / Lint

`biblio-cli` (the Go module at the repo root) is real, working source
code as of v0.01, extended through v0.06:

```sh
go build ./...        # build everything, including cmd/biblio-cli
go run ./cmd/biblio-cli --help
go test ./...          # full unit + integration suite
./scripts/check.sh     # gofmt -l . (fails on unformatted files) + go vet ./...
```

`go install github.com/monamaret/biblio/cmd/biblio-cli@latest` installs
the CLI directly; prebuilt binaries are also published via `goreleaser`
on GitHub Releases (Principle III) — see
[docs/releases/v0.06/README.md](docs/releases/v0.06/README.md) for a
build-from-source walkthrough.

Each scaffolded site project (`templates/site/{profile,log,docs,quarterly}`) is its
own independent npm project (React + Vite + React Router v7 + Wrangler);
once scaffolded with `biblio-cli new`, the generated project's own
`package.json` carries its `build`/`test`/`lint`/`dev` scripts (e.g.
`npx vitest`, `npx tsc --noEmit`) — there is no repo-root Node tooling,
since this repo itself ships only the Go CLI and the embedded template
sources.

**Never leave `node_modules/` inside `templates/site/<type>/`.** The
templates are embedded with `go:embed all:site`, which embeds whatever is
on disk: an installed `node_modules/` (~550 MB) gets compiled into every
Go test binary and can fill `/tmp` mid-link. To run a template's own
tests, mirror it elsewhere (e.g. `rsync -a --exclude node_modules
templates/site/quarterly/ <scratch>/` and `npm install` there).
`scaffold.Walk` also skips `node_modules/` and `.react-router/`, so a
stray install never reaches a scaffolded project.

`biblio-cli build` and `biblio-cli dev` run a pre-build step
(`internal/prebuild`) before Vite for every site type: validate
`content/site.yaml`, then regenerate `content-index.json` — so
`biblio-cli index` no longer needs to be run by hand. Quarterly sites
also regenerate `public/tokens.css` from their `typesetting:` block and
convert content images to grayscale plates (`internal/plates`) in
`public/plates/`. `content-index.json`, `public/tokens.css` and
`public/plates/` are build output (`scaffold.IsGenerated`): `update`/
`upgrade` never treat them as template regions.

Tests are written and confirmed failing *before* the corresponding
implementation, on both sides — see the test-driven development
convention above. No change lands without a failing test demonstrating
the gap it closes (Constitution Principle VI, NON-NEGOTIABLE).

## Subagents

None yet — `.opencode/agents/` only contains its own `README.md`
(format reference for defining project-scoped subagents). Add entries
here, one line per agent, as they're created.

## Skills

None yet — there is no `.opencode/skills/` directory. Add entries here,
one line per skill, as project-specific skills are created.

## Feature development process

The fixed pipeline from idea to shipped feature — steps are not skipped,
reordered, or collapsed, even for small features:

1. **Numbered research discussion document** — open or extend a numbered doc
   under `specs/research/` (copy `research-template.md` for a new one) and
   record confirmed decisions in its Decision Log (date, decision, rationale,
   owner — see `001-static-site-platform-with-mcp.md` for the pattern;
   superseded entries are kept, marked **Superseded**, not deleted).
2. **Feature item creation** — create a feature item under `specs/features/`
   from `specs/features/TEMPLATE.md`, pointing back to the now-confirmed
   decisions rather than re-deciding anything inline.
3. **Roadmap inclusion** — add the feature item to
   [specs/research/roadmap.md](specs/research/roadmap.md) in its
   dependency-ordered place, keeping the roadmap the single current view of
   "what's next."
4. **Architecture sync** — update
   [specs/research/architecture.md](specs/research/architecture.md) with any
   resulting diagram/call-flow changes *before* implementation starts.
5. **Spec Kit workflow** — only then run the Spec Kit (speckit) workflow
   (`/speckit.specify` → `/speckit.clarify` → `/speckit.plan` →
   `/speckit.tasks` → `/speckit.implement`, etc.) to produce spec/plan/tasks
   artifacts and drive implementation, following the documented process at
   https://github.com/github/spec-kit/blob/main/spec-driven.md rather than
   improvising a different artifact shape.

## Release process

Each release is packaged as a directory `specs/releases/X.XX-name/` (e.g.
`specs/releases/0.06-general-availability/`) once its features have actually
shipped — not drafted ahead of the work it documents. It MUST contain:

- **Final feature list** — the definitive shipped-feature list, each entry
  linking back to its `specs/features/` item (and from there to the research
  doc/Decision Log entries that justified it).
- **Release notes** — a plain, owner-facing summary of what changed and why
  (no external-user/marketing framing — this is a single-author project).
- **Unit and integration tests** — both layers, covering the release's
  features end to end, per the test-driven-development convention above; a
  release isn't "final" until both exist and pass.
- **Usage guide** — covering, at minimum, how to **create**, **configure**,
  **update**, **deploy**, and **destroy** a `biblio` site with that release's
  `biblio-cli` — the full site lifecycle, end to end.

Each `specs/releases/0.0X-release.md` spec carries a `**Status**` field,
starting as `Draft`. Going forward, flip it to `Shipped` once that release's
`tasks.md` is fully implemented (every task `[X]`) and its closing
quickstart/integration evidence has run successfully — do this as part of
the Polish phase's doc-update task. (v0.01–v0.03's release docs still say
`Draft` despite having shipped; this convention starts with v0.04 rather
than being backfilled onto those.)

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
