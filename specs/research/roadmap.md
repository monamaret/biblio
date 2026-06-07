# Roadmap: `biblio` implementation plan

**Status:** Draft — tentative, ordered by dependency rather than effort.
**Progress cursor:** → §0 cleared — every blocking design question is resolved (the only items still open — storage-escalation migration steps, site-update mechanics, embedding-model choice — are explicitly conditional/deferred and don't block anything). Next up: begin §1 (`biblio-cli` foundations), the foundation every other deliverable scaffolds from. *(Move this pointer forward — to a section number and short note — as work actually begins/completes; it's the at-a-glance answer to "where are we right now.")*

This roadmap intentionally **makes no assumptions about timelines, developer hours, or any other "real world" delivery metric** — none of those have been discussed or decided anywhere in the research (see the non-goals in [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md)), and `biblio` is a sole-developer, personal-scale project where such estimates would be guesswork. Every task below is tied to a deliverable that is already defined by the confirmed decisions and completed research in [001](001-static-site-platform-with-mcp.md), [000-tech-stack.md](000-tech-stack.md), and [architecture.md](architecture.md) — nothing here is speculative scope.

## 0. Outstanding design questions

This section is now the **single tracked location for every open research/design item** across the `biblio` docs — surfaced during research but deliberately left open rather than guessed at. Resolving them shapes several of the sections below, so they're listed first as the things to settle before (or while) building the dependent pieces. *(Consolidated 2026-06-07: items formerly duplicated as inline TODOs in [deployment-guide.md](deployment-guide.md) and as a stale "Open questions" note in [000-tech-stack.md](000-tech-stack.md) now live here only — those docs point back to this section instead of carrying their own copies.)*

### 0.1 Architecture & data-model

- [x] ~~Define the on-disk content directory layout convention~~ — **Resolved 2026-06-07: flat-by-concept directories under `content/`** (`profile.md`/`experience|skills|projects/*.md` for profile sites, `entries/YYYY-MM-DD-slug.md` for log sites, `docs/<section>/<page>.md` for docs sites, plus `content/site.yaml` for shared config) — not type-prefixed, since each site has exactly one fixed `site_type`. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log) for the full convention and rationale. No longer blocks §1/§2/§3.
- [x] ~~Define the build-time MCP opt-in/out configuration shape~~ — **Resolved 2026-06-07: a single flat boolean, `mcp_enabled` (default `true`), in `content/site.yaml`** alongside `site_name`/`description`/`site_type` — not a nested object, since the toggle is the only MCP-related build-time setting identified. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log) for the full shape and rationale. No longer blocks §1/§4.
- [x] ~~Decide whether the docs frontmatter schema needs `tags`/`project` fields~~ — **Resolved 2026-06-07: no.** The Docs Page schema stays `title`/`order`/`section`/`summary`, with no matching filter params on `search_docs`/`get_toc`. `tags`/`project` exist on the project-log schema specifically because that site type spans multiple projects in one timeline; a docs site is already scoped to the one project it documents, and its `section`/`order` hierarchy is its natural organizing principle. "Configurable by site type" is already satisfied by the existing per-type fixed-schema design itself. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log) for the full rationale. No longer affects §3/§4.

### 0.2 Terminal/CLI-view spikes

- [x] ~~Confirm whether `md4x`'s `renderToAnsi()` runs in the Cloudflare Workers runtime~~ — **Resolved 2026-06-07: yes, via manual WASM-module wiring.** Source-level inspection of `md4x`'s WASM loader shows its automatic `init()` only auto-detects Node.js (`fs/promises`) or bundler-specific `?url`/`?module` imports — neither of which Wrangler supports — but `init()` also accepts a pre-instantiated `WebAssembly.Module` directly (`init({ wasm })`), which is exactly what Wrangler's native `.wasm` import support produces (`import wasmModule from "md4x/build/md4x.wasm"`). So tier-ii in-Worker Markdown → ANSI rendering via `renderToAnsi()` is viable with no separate build-time pipeline, carrying this one specific integration detail into §7. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log) for the full integration snippet and rationale. No longer blocks §7.
- [x] ~~Confirm whether Cloudflare's "Markdown for Agents" edge feature reaches Worker-served responses~~ — **Resolved 2026-06-07: it does reach them, but the "for free" framing was wrong — it's gated by the zone's (domain's) plan tier, not the Workers compute plan, and is unavailable on the Free zone plan.** Independently confirmed: a real Workers + Static Assets site (architecturally identical to `biblio`'s plan) gets full conversion via the managed toggle — *if* its zone is on Cloudflare Pro/Business/Enterprise (a domain-level plan, separate from the Workers Free/Paid compute plan, and the thing actually gated here). On the **Free zone plan**, the managed toggle is unavailable entirely, and `biblio` must hand-roll the negotiation+conversion branch — which is exactly what §7 already plans to build (build-time-generated Markdown + `Accept`-led negotiation in the Worker's `fetch` handler), independently validated by another author doing the same thing for the same reason. **Net effect: build the hand-rolled §7 path regardless** — it's required on Free, gives full output-quality control either way (no nav-chrome/footer/cookie-banner noise the generic HTML→Markdown converter can't strip), and the managed toggle becomes an optional zone-level bonus to flip on later *if* the zone ever sits on Pro+, not a load-bearing dependency. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log) for the full caveats (response `Content-Length`/size requirements). No longer blocks §7.
- [x] ~~Decide how much to invest in `llms.txt`~~ — **Resolved 2026-06-07: minimal.** A build-time-generated file listing `site_name`/`description` plus a flat list of key page links, derived directly from `content-index.json` with no extra authoring step — not a fuller `llms-full.txt`-style content dump. Rationale: research already found major LLM crawlers don't request it in meaningful volume — it's "expected convention to not be missing" rather than a proven traffic driver, so a fuller dump would be real ongoing investment (size budgeting, staleness risk, format-spec churn) chasing an unmaterialized benefit. The minimal version satisfies "don't be conspicuously missing the convention" at near-zero cost and can expand later if usage data ever justifies it. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log). No longer blocks §7.
- [x] ~~Decide whether a full terminal-styled "mini-site" view (rendering tier iii) is worth building~~ — **Resolved 2026-06-07: no, not for the initial scope.** Tiers i/ii (raw Markdown passthrough + the now-confirmed-viable in-Worker `md4x` → ANSI rendering) already cover the actual need — a terminal/CLI client gets a clean, readable view either way. A full tier-iii view would be a substantial speculative investment (a whole second styled view, its own design/testing surface) for a need that hasn't been validated — exactly the kind of "designing for hypothetical future requirements" the project avoids elsewhere. It becomes a real, answerable question only if real terminal-client usage ever shows demand for more than a clean ANSI render. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log). No longer blocks §7 — the last tier-iii task below is dropped from the initial scope.

### 0.3 CLI & workflow design

*(Moved here from inline TODOs in [deployment-guide.md](deployment-guide.md) — that guide now references this section instead of stating these as its own open items.)*

- [x] ~~`biblio-cli` install/distribution method~~ — **Resolved 2026-06-07: GitHub Releases with prebuilt binaries (via `goreleaser`) as the primary path, `go install github.com/.../biblio-cli@latest` documented as the alternative for those with a Go toolchain.** No package-manager tap (Homebrew, etc.) — that's distribution infrastructure for a multi-user audience `biblio` doesn't have; a single private author needs at most two install paths. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log). No longer blocks §1.
- [x] ~~Exact `biblio-cli new` command and flags~~ — **Resolved 2026-06-07: `biblio-cli new <project-name> --type=<profile|log|docs> --domain=<domain> --theme=<preset>`, with the Bubble Tea TUI driving an interactive prompt flow whenever flags are omitted.** Flags support quick repeat use/scripting; the guided interactive flow gives the one-time "create a site" moment the polished experience the Bubble Tea investment already signals. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log). No longer blocks §1.
- [x] ~~Local development workflow~~ — **Resolved 2026-06-07: `biblio-cli` wraps `vite dev` + `wrangler dev` behind one `biblio-cli dev` command.** Running two dev servers by hand is friction for a single author with no team to coordinate across; one command matches the smooth, CLI-first DX the Bubble Tea investment already signals. See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log). No longer blocks §1/§3.
- [x] ~~Build/deploy command surface~~ — **Resolved 2026-06-07: `biblio-cli` exposes `biblio-cli build`/`biblio-cli deploy` wrappers** (thin shells over `vite build`/`wrangler deploy` + the generated CI workflow) — **and the site-update/re-scaffold story (item below) is folded into this same command surface** (e.g. `biblio-cli update`/`biblio-cli upgrade`), keeping one consistent CLI-first surface end to end rather than splitting "dev/build/deploy is one command, updates are raw `package.json`/CI." See [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log). No longer blocks §1.
- [ ] **Storage-escalation migration steps** — concrete, step-by-step instructions for the JSON→D1 and keyword→Vectorize paths. *Conditional* — only worth writing once a real site actually crosses one of the documented heuristic thresholds (§10 below, explicitly deferred). Affects: deployment-guide §9.
- [ ] **Site-update / re-scaffold story — exact mechanics** — *direction now locked* (it will be a `biblio-cli` command, per the build/deploy resolution above, e.g. `biblio-cli update`/`biblio-cli upgrade`); the concrete mechanics (diffing strategy, conflict handling, what gets re-templated vs. preserved) remain genuinely *deferred* — only answerable once more than one site exists to migrate (§9 below). Affects: deployment-guide §10.

### 0.4 Deferred / conditional follow-ons

*(Moved here from a stale "Open questions" note in [000-tech-stack.md](000-tech-stack.md) — that doc's embedding-provider question is resolved; only this narrower follow-on remains, and only matters if the Vectorize escalation in §10 is ever actually built.)*

- [ ] **Which embedding model to use within OpenRouter** for build-time embedding generation (e.g. `text-embedding-3-small` vs. `text-embedding-3-large` vs. BGE-M3 vs. Qwen3 Embedding) — the *provider* choice (OpenRouter over DeepSeek) is already decided (see [000-tech-stack.md → Decision Log](000-tech-stack.md#decision-log)); only the specific model is open, and only matters if/when §10 (Vectorize escalation) is actually built.

## 1. `biblio-cli` foundations

Deliverable: a working CLI binary (Go + Cobra + Viper + Bubble Tea, per the [confirmed stack](001-static-site-platform-with-mcp.md#libraries--stack)) that can scaffold a new site project end to end.
**Feature spec:** [specs/features/001-biblio-cli-foundations.md](../features/001-biblio-cli-foundations.md)

- [ ] Set up the Go module, with Cobra for commands and Viper for config management
- [ ] Set up `goreleaser`-based GitHub Releases distribution (prebuilt binaries) as the primary install path, with `go install github.com/.../biblio-cli@latest` documented as the alternative — confirmed §0.3 resolution; no package-manager tap (see [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Build the interactive scaffolding flow with Bubble Tea, styled with Lip Gloss
- [ ] Implement the `new`-site command — confirmed syntax `biblio-cli new <project-name> --type=<profile|log|docs> --domain=<domain> --theme=<preset>`, with the Bubble Tea TUI driving an interactive prompt flow whenever flags are omitted (see [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log)): generates a combined React + Vite + Worker project for the chosen site type and domain, seeding `content/` per the [confirmed layout convention](001-static-site-platform-with-mcp.md#decision-log)
- [ ] Implement `site.yaml` generation (`site_name`/`description`/`site_type`, plus `mcp_enabled` per the [confirmed flag shape](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Implement `wrangler.jsonc` generation with the route/custom-domain binding pre-filled
- [ ] Implement `tokens.css` generation from the chosen site type's theme preset (depends on §6)
- [ ] Generate the GitHub Actions workflow wired to the [three-step build pipeline](architecture.md#2-build--deployment-pipeline)
- [ ] Implement the `biblio-cli dev` command — wraps `vite dev` + `wrangler dev` behind one command (confirmed §0.3 resolution; see [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Implement the `biblio-cli build`/`biblio-cli deploy` command wrappers — thin shells over `vite build`/`wrangler deploy` + the generated CI workflow (confirmed §0.3 resolution; see [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Implement the `biblio-cli update`/`biblio-cli upgrade` command — direction confirmed (folded into the CLI command surface per [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log)); exact mechanics remain deferred until §9 produces a second site to migrate
- [ ] Evaluate Glamour for in-terminal rendering of the CLI's own output and any generated docs/guides

## 2. Content pipeline — `content-index.json` generator

Deliverable: the Go-side generator that produces the single canonical artifact consumed by both the site build and the MCP server, per the [progressive-enhancement storage approach](001-static-site-platform-with-mcp.md#architecture-patterns).
**Feature spec:** [specs/features/002-content-pipeline.md](../features/002-content-pipeline.md)

- [ ] Implement the content-directory walker, targeting the [confirmed layout convention](001-static-site-platform-with-mcp.md#decision-log) (`profile.md`/`experience|skills|projects/`, `entries/`, `docs/<section>/`, `site.yaml`)
- [ ] Implement frontmatter + body parsing per site type, matching the [confirmed schemas](001-static-site-platform-with-mcp.md#architecture-patterns) — Profile/Experience/Skill/Project, Log entry, Docs page, and the shared site-level config
- [ ] Implement slug/path derivation rules — from filename for log entries, from file location for docs pages
- [ ] Emit `content-index.json` in a shape consumable both by React Router's `prerender` function (§3) and the in-Worker MCP query layer (§4)
- [ ] Add schema validation for authored frontmatter (the Starlight survey flagged Zod-validated frontmatter as a good design-pattern reference for this)

## 3. Site scaffold — React + Vite + Worker combined project template

Deliverable: the project template that the `new`-site command (§1) generates — the static-site half of the [selected architecture](architecture.md#1-high-level-architecture-selected).
**Feature spec:** [specs/features/003-site-scaffold.md](../features/003-site-scaffold.md)

- [ ] Set up the React + Vite + React Router v7 project skeleton, with native `prerender` configured in `react-router.config.ts`
- [ ] Wire the `prerender` async function to enumerate content-driven routes from `content-index.json` (log-entry lists, docs TOC pages, etc.)
- [ ] Implement content loading via `import.meta.glob` + `gray-matter`, rendered with `react-markdown` (+ `remark-gfm`)
- [ ] Build the per-site-type page templates/layouts:
  - [ ] Personal profile/CV — profile, experience, skills, and projects views
  - [ ] Technical project log — filterable entry list (by date/tag/project) and entry detail view
  - [ ] Project documentation — TOC/nav generated from file structure, and page view (no tag/project filter UI — see [confirmed schema decision](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Integrate Tailwind CSS + shadcn/ui, consuming the generated `tokens.css` (depends on §6)
- [ ] Wire the Workers Static Assets binding to the prerendered build output

## 4. In-Worker MCP server

Deliverable: the custom MCP server implementation sharing the Worker with the static site, per the [selected pattern](architecture.md#1-high-level-architecture-selected).
**Feature spec:** [specs/features/004-in-worker-mcp-server.md](../features/004-in-worker-mcp-server.md)

- [ ] Set up the `agents`/`workers-mcp` SDK in the Worker, with Durable Objects for session state
- [ ] Implement in-Worker path routing — `/mcp*` to the MCP server, fall through to static assets otherwise
- [ ] Implement the `content-index.json` import and an in-memory query layer (filter/lookup/search)
- [ ] Implement the shared tools: `get_site_info` and the generic `search`
- [ ] Implement the per-site-type tool surface (full mapping in [001 → tool surface](001-static-site-platform-with-mcp.md#architecture-patterns)):
  - [ ] Personal profile/CV — `get_profile`, `list_experience`, `list_skills`, `list_projects`, `get_resume`
  - [ ] Technical project log — `list_entries`, `get_entry`, `search_entries`
  - [ ] Project documentation — `get_toc`/`list_pages`, `get_page`, `search_docs` (keyword search + section/path scoping only — no tag/project filter params, see [confirmed schema decision](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Implement the build-time MCP opt-in/out gating driven by `mcp_enabled` — conditional content-index generation, route mounting, and bindings, per the [confirmed flag shape](001-static-site-platform-with-mcp.md#decision-log)
- [ ] Confirm the `/mcp` endpoint is reachable with no auth, relying on Cloudflare's standard protections (rate limiting, WAF, bot management)

## 5. MCP App UI layer

Deliverable: the small React + Vite UI-resource bundles an MCP host renders inline, per the [confirmed rendering decisions](001-static-site-platform-with-mcp.md#architecture-patterns).
**Feature spec:** [specs/features/005-mcp-app-ui-layer.md](../features/005-mcp-app-ui-layer.md)

- [ ] Set up the MCP App UI bundle build (small, sandboxed, postMessage-bridged React + Vite bundles)
- [ ] Wire the shared `tokens.css` into the UI bundle build so it inherits the site's theme automatically
- [ ] Hand-roll minimal SVG/CSS chart components (timelines, skill-level bars, simple bar/line charts) styled via `tokens.css` — no charting-library dependency
- [ ] Implement the text-first, progressive-enhancement rule: every chart-like tool result produces a clean, self-sufficient Markdown/ASCII rendering (table, sparkline) as its primary `content`, with the richer UI resource layered on top only for hosts that support it

## 6. Theming system

Deliverable: the per-site-type, token-driven theming shared between the site and its MCP App UI, per the [confirmed theming decision](001-static-site-platform-with-mcp.md#architecture-patterns).
**Feature spec:** [specs/features/006-theming-system.md](../features/006-theming-system.md)

- [ ] Define the shared shadcn/ui CSS-variable token contract (colors + `-foreground` pairs, `border`/`input`/`ring`, `radius`, plus Tailwind font/spacing extensions)
- [ ] Build the three per-site-type presets — CV/profile (warmer/personal), project log (cooler/technical, monospace accents), docs (neutral/high-contrast)
- [ ] Mine `terminal.css` / WebTUI / `system.css` / `tui.css` for palette values (Gruvbox, Catppuccin, Nord, Dracula), monospace type choices, and box-drawing/CRT-style effects to bake into the project-log and docs presets
- [ ] Implement `biblio-cli`'s `tokens.css` generation from a chosen preset — one file, imported by both the site build and the MCP App UI bundle

## 7. Terminal/CLI-optimized view of the site

Deliverable: a terminal-optimized rendering of site pages served to `curl`/`lynx`/terminal-coding-agent clients, per the [confirmed research findings](001-static-site-platform-with-mcp.md#architecture-patterns). All four §0.2 terminal/CLI-view spikes are now resolved — this is a fully-specified two-tier (i/ii) execution task.
**Feature spec:** [specs/features/007-terminal-cli-view.md](../features/007-terminal-cli-view.md)

- [ ] Implement layered content negotiation in the Worker's `fetch` handler: parse `Accept` q-values (rank `text/markdown` vs. `text/html` properly, don't substring-match), fall back to User-Agent matching only when `Accept` is uninformative, and use `sec-fetch-dest: document` as a browser safeguard
- [ ] Set `Vary: Accept, User-Agent` on every negotiated response so Cloudflare's cache serves the correct variant to the correct client
- [ ] Add `.md` URL-sibling routes and `Link: rel="alternate"; type="text/markdown"` discoverability headers
- [ ] Generate a minimal `llms.txt` from the content index — confirmed scope: `site_name`/`description` plus a flat list of key page links derived directly from `content-index.json`, no fuller `llms-full.txt`-style dump (see [confirmed §0.2 resolution](#02-terminalcli-view-spikes) / [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Implement raw Markdown passthrough rendering (tier i — the cheapest, always-available baseline)
- [ ] Implement in-Worker lightly-formatted Markdown → ANSI rendering via `md4x`'s `renderToAnsi()` (tier ii) — confirmed viable, with no separate build-time pipeline; wire the WASM module via Wrangler's native `.wasm` import (`import wasmModule from "md4x/build/md4x.wasm"`) and pass it to `init({ wasm: wasmModule })`, since `md4x`'s automatic init doesn't detect the Workers runtime (see [confirmed §0.2 resolution](#02-terminalcli-view-spikes) / [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log))
- [x] ~~Design and build a full terminal-styled "mini-site" view (tier iii), reusing the retro/terminal token research from §6~~ — **Dropped from initial scope 2026-06-07: tiers i/ii cover the need on their own; tier iii would be speculative investment with no validated demand.** See [confirmed §0.2 resolution](#02-terminalcli-view-spikes) / [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log). Revisit only if real terminal-client usage shows demand for more than a clean ANSI render.

## 8. Pilot site validation

Deliverable: one real site, built and deployed end to end, validating the scaffold + CLI + architecture before replication — per the [rollout plan](001-static-site-platform-with-mcp.md#performance-reliability--operations).
**Feature spec:** [specs/features/008-pilot-site-validation.md](../features/008-pilot-site-validation.md)

- [ ] Scaffold and deploy a personal profile/CV site as the pilot
- [ ] Verify the full build pipeline end to end: content-index generation → `react-router build` (prerender) → Wrangler deploy
- [ ] Verify the `/mcp` endpoint responds correctly to a real MCP client (`initialize`, `tools/list`, `tools/call`)
- [ ] Verify the MCP App UI renders correctly in at least one graphical host, and degrades correctly (text-first) in a terminal-based host
- [ ] Capture and address any scaffold/pipeline issues found, feeding fixes back into §1–§6

## 9. Replicate to the remaining (priority) site types

Deliverable: the two scaffolds called out as priorities in the research — project log and documentation — built and validated, per [001 → Background](001-static-site-platform-with-mcp.md#background--context).
**Feature spec:** [specs/features/009-replicate-site-types.md](../features/009-replicate-site-types.md)

- [ ] Scaffold and deploy a technical project log site; validate tag/project filtering and the daily-update authoring flow
- [ ] Scaffold and deploy a project documentation site; validate TOC/nav generation and `section`/`order`-based browsing (no tag/project filtering — confirmed as out of scope for this site type, see [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log))
- [ ] Confirm theming presets and the terminal/CLI view render correctly across all three site types
- [ ] Design and implement the concrete mechanics of the `biblio-cli update`/`biblio-cli upgrade` command (diffing strategy, conflict handling, what gets re-templated vs. preserved) — *direction* already locked as a CLI command (see [§0.3 resolution](#03-cli--workflow-design) / [001 → Decision Log](001-static-site-platform-with-mcp.md#decision-log)); the mechanics are the one part still genuinely open, first answerable once this section creates a second and third site to migrate

## 10. Storage escalation (conditional — not built up front)

Deliverable: the D1 and Vectorize escalation paths described in the [progressive-enhancement storage approach](001-static-site-platform-with-mcp.md#architecture-patterns). These are **explicitly not scheduled work** — they're only triggered if/when a real site crosses one of the documented heuristic thresholds.
**Feature spec:** [specs/features/010-storage-escalation.md](../features/010-storage-escalation.md)

- [ ] *(Conditional — bundled JSON index outgrows comfortable in-memory handling)* Provision a D1 database, add the binding, and switch content-index population to populate D1 from the same content directory at build time
- [ ] *(Conditional — real-usage signals show keyword search underperforming on conceptual queries, most likely on a docs site)* Wire OpenRouter build-time embedding generation (per [000-tech-stack.md](000-tech-stack.md)), provision a Vectorize index, and switch `search_docs`/`search_entries` to query it

## References

- [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) — the research and confirmed decisions every task above is derived from
- [architecture.md](architecture.md) — diagrams of the system these tasks build out
- [deployment-guide.md](deployment-guide.md) — the usage guide that documents the result; its open questions are tracked in §0 here rather than duplicated there
- [000-tech-stack.md](000-tech-stack.md) — the platform/service responsibility mapping this plan builds on
- [../features/](../features/) — feature item specs (one per §1–§10 deliverable above, from [../features/TEMPLATE.md](../features/TEMPLATE.md)), referenced when development on each begins
