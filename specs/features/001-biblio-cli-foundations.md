# Feature: `biblio-cli` foundations

**Status:** Not started
**Roadmap reference:** [roadmap.md §1 — `biblio-cli` foundations](../research/roadmap.md#1-biblio-cli-foundations)

## Summary

A working CLI binary, built in Go, that can interactively scaffold a new combined React + Vite + Worker site project end to end — the tool every other `biblio` deliverable is generated or operated through.

## Confirmed decisions this builds on

- [Decision Log — content directory layout convention](../research/001-static-site-platform-with-mcp.md#decision-log): flat-by-concept directories under `content/` (`profile.md`/`experience|skills|projects/`, `entries/YYYY-MM-DD-slug.md`, `docs/<section>/<page>.md`, `site.yaml`) — what the `new`-site command seeds, per chosen site type
- [Decision Log — MCP opt-in/out configuration shape](../research/001-static-site-platform-with-mcp.md#decision-log): a single flat `mcp_enabled` boolean (default `true`) in `content/site.yaml` — exactly what the generated site config must contain
- [Decision Log — CLI implementation stack](../research/001-static-site-platform-with-mcp.md#decision-log): Go + Cobra + Viper, with Charmbracelet's Bubble Tea (and likely Lip Gloss/Glamour) for an interactive TUI — a single static binary for easy distribution
- [Libraries & Stack — confirmed stack table](../research/001-static-site-platform-with-mcp.md#libraries--stack): same stack, listed as **Selected**
- [architecture.md §2 — Build & deployment pipeline](../research/architecture.md#2-build--deployment-pipeline): the three-step pipeline (`biblio-cli` content-index generation → `react-router build` → `wrangler deploy`) this CLI scaffolds the GitHub Actions workflow for
- [Decision Log — build pipeline integration](../research/001-static-site-platform-with-mcp.md#decision-log): `biblio-cli` owns content-index generation as a Go-side generator, not a Vite plugin/npm script
- [Decision Log — `biblio-cli` workflow/UX decisions](../research/001-static-site-platform-with-mcp.md#decision-log): the full confirmed command surface and distribution method —
  - **Install/distribution:** GitHub Releases with prebuilt binaries (`goreleaser`) as the primary path, `go install .../biblio-cli@latest` documented as the alternative — no package-manager tap
  - **`new` command:** `biblio-cli new <project-name> --type=<profile|log|docs> --domain=<domain> --theme=<preset>`, with the Bubble Tea TUI driving an interactive prompt flow whenever flags are omitted
  - **`dev` command:** wraps `vite dev` + `wrangler dev` behind one `biblio-cli dev` command
  - **`build`/`deploy` commands:** `biblio-cli build`/`biblio-cli deploy` wrappers — thin shells over `vite build`/`wrangler deploy` + the generated CI workflow
  - **`update`/`upgrade` command:** site-update/re-scaffold direction folded into this same CLI surface (`biblio-cli update`/`biblio-cli upgrade`); exact mechanics deferred to [009-replicate-site-types](009-replicate-site-types.md)

## Scope

- [ ] Set up the Go module, with Cobra for commands and Viper for config management
- [ ] Set up `goreleaser`-based GitHub Releases distribution (prebuilt binaries) as the primary install path, with `go install .../biblio-cli@latest` as the documented alternative
- [ ] Build the interactive scaffolding flow with Bubble Tea, styled with Lip Gloss
- [ ] Implement the `new`-site command: `biblio-cli new <project-name> --type=<profile|log|docs> --domain=<domain> --theme=<preset>`, with an interactive Bubble Tea prompt flow when flags are omitted — generates a combined React + Vite + Worker project for the chosen site type and domain
- [ ] Implement site-level config generation (`site_name`/`description`/`site_type`, plus the MCP opt-in/out flag)
- [ ] Implement `wrangler.jsonc` generation with the route/custom-domain binding pre-filled
- [ ] Implement `tokens.css` generation from the chosen site type's theme preset
- [ ] Generate the GitHub Actions workflow wired to the three-step build pipeline
- [ ] Implement the `biblio-cli dev` command — wraps `vite dev` + `wrangler dev` behind one command
- [ ] Implement the `biblio-cli build`/`biblio-cli deploy` command wrappers — thin shells over `vite build`/`wrangler deploy` + the generated CI workflow
- [ ] Implement the `biblio-cli update`/`biblio-cli upgrade` command — direction confirmed (CLI-surfaced, not a separate tool/manual process); exact mechanics owned by [009-replicate-site-types](009-replicate-site-types.md) once a second site exists to migrate
- [ ] Evaluate Glamour for in-terminal rendering of the CLI's own output and any generated docs/guides

## Out of scope

- The content-index generator's parsing logic itself — covered by [002-content-pipeline](002-content-pipeline.md), though this CLI is its host binary
- The site project template's internals (page templates, routing, MCP server code) — covered by [003-site-scaffold](003-site-scaffold.md) and [004-in-worker-mcp-server](004-in-worker-mcp-server.md); this feature only scaffolds the project shell
- Theme preset *content* (palettes, fonts) — covered by [006-theming-system](006-theming-system.md); this feature only wires up `tokens.css` generation as a build step

## Open questions

None from [roadmap.md §0](../research/roadmap.md#0-outstanding-design-questions) remain open for this feature — the install/distribution method, `new` command syntax, local dev workflow, and build/deploy/update command surface are all resolved (see Confirmed decisions above). This is now an execution task against a fully-specified command surface.

## Dependencies

None — this is the foundation every other feature's scaffolding output depends on.

## Acceptance criteria

- `biblio-cli new` produces a project that builds and deploys via the three-step pipeline without manual edits
- The generated project includes a working `wrangler.jsonc`, site config, `tokens.css`, and GitHub Actions workflow
- The interactive scaffolding flow runs end to end in a terminal via Bubble Tea
- `biblio-cli dev` runs `vite dev` and `wrangler dev` together behind one command, and `biblio-cli build`/`biblio-cli deploy` correctly wrap the underlying `vite build`/`wrangler deploy` + CI steps
- A user can install the CLI via either a downloaded GitHub Release binary or `go install`, with both paths producing a working `biblio-cli`

## References

- [roadmap.md](../research/roadmap.md) — task breakdown this feature implements
- [001-static-site-platform-with-mcp.md](../research/001-static-site-platform-with-mcp.md) — research and confirmed decisions
- [architecture.md](../research/architecture.md) — reference diagrams
- [deployment-guide.md](../research/deployment-guide.md) — end-to-end usage walkthrough
