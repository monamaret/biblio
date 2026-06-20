---

description: "Task list for v0.04 — Prove the platform on the technical project log"
---

# Tasks: v0.04 — Prove the platform on the technical project log

**Input**: Design documents from `/specs/releases/0.04-project-log-site-type/`
**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/log-site-surface.md](./contracts/log-site-surface.md), [quickstart.md](./quickstart.md)

**Tests**: Included and REQUIRED — Constitution Principle VI and plan.md's Testing section require failing-test-first coverage for the schema/filter/tool/rendering units, plus a release-closing "scaffold → author a run of dated entries → deploy → filter → agent/terminal query" integration check.

**Codebase reality check (read before starting)**: `templates/site/log/` is still the thin v0.01-era Vite/React placeholder (`src/main.tsx`, `index.html`, no `app/` tree) — it never received the React Router v7 / Worker / MCP / negotiation buildout `templates/site/profile/` got in v0.02–v0.03. `internal/scaffold/generator.go`'s `copyTemplateTree` copies each site type's template tree wholesale with no shared base, so `log` must independently reach parity. Phase 1 below does that parity work first; it is **not** optional scaffolding cleanup, it's the actual prerequisite plan.md's "purely additive" framing assumed was already true.

One correction to plan.md's Project Structure sketch: the sketch lists `internal/index/filters.go` (Go) for date/tag/project filtering. Filtering is a derived query over the already-generated `content-index.json`, consumed only by the browser UI and MCP tools — exactly like `searchEntries` in `templates/site/profile/app/mcp/query.ts`, which is TS-only with no Go-side duplicate. A Go-side `filters.go` would create a second filter-logic owner, violating Principle II. Tasks below implement filtering once, in TS, shared by both consumers (T025).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 = Scenario 1 (author/publish a run of dated entries), US2 = Scenario 2 (filter the entry list), US3 = Scenario 3 (daily authoring rhythm), US4 = Scenario 4 (agent/terminal layer for this site type)

## Path Conventions

- `templates/site/log/app/...`, `templates/site/log/workers/entry.ts` — mirrors `templates/site/profile/app/` and `workers/entry.ts`'s structure and naming
- `internal/content/schema.go`, `internal/content/schema_test.go` — log-specific structural validation extends the existing shared schema, no new `internal/content/log/` subpackage (one schema owner, Principle II)
- `test/integration/pilot-site-log/` — matches the `test/integration/pilot-site-mcp/` convention from v0.03

---

## Phase 1: Setup (Bring `templates/site/log` to parity with `templates/site/profile`'s infra)

**Purpose**: Nothing log-specific can be built until the log template has the same React Router v7 / Vite / Worker / test scaffolding profile already has

- [X] T001 [P] Replace `templates/site/log/package.json`'s plain Vite/React deps with the profile set (`react-router`, `@react-router/dev`, `agents`, `@modelcontextprotocol/sdk`, `zod`, `md4x`, `react-markdown`, `remark-gfm`, `isbot`, vitest/testing-library devDeps) per `templates/site/profile/package.json`
- [X] T002 [P] Replace `templates/site/log/vite.config.ts` with the React Router v7 + `llmsTxtPlugin()` config from `templates/site/profile/vite.config.ts` (same rationale: no `@cloudflare/vite-plugin`)
- [X] T003 [P] Add `templates/site/log/react-router.config.ts`, `templates/site/log/tsconfig.json`, `templates/site/log/vitest.config.ts`, `templates/site/log/vitest.setup.ts` mirroring profile's (type-agnostic, no log-specific changes)
- [X] T004 [P] Remove the superseded v0.01 plain-Vite files `templates/site/log/src/main.tsx` and `templates/site/log/index.html`
- [X] T005 [P] Add `templates/site/log/.github/workflows/deploy.yml` (copy of `templates/site/profile/.github/workflows/deploy.yml` — CI/deploy contract is type-agnostic, FR-001)
- [X] T006 [P] Copy `templates/site/profile/app/content/load.ts` and `render.tsx` into `templates/site/log/app/content/` unchanged (Principle II — content loading/rendering is type-agnostic; data-model.md "Project-Log Entry" extends the same Authored Content Item)
- [X] T007 Author `templates/site/log/app/root.tsx` and `templates/site/log/app/routes.ts` (route skeleton only: `index` and a `"*"` detail splat, matching `templates/site/profile/app/routes.ts`; the filter route is added in Phase 4) — `routes.ts` shipped as a valid *empty* `RouteConfig` rather than one importing not-yet-existing `home.tsx`/`detail.tsx`, so it typechecks at this checkpoint; Phase 3 (T019) fills it in

**Checkpoint** (verified via a real `biblio-cli new --type=log` smoke scaffold): `templates/site/log` installs, typechecks (`tsc --noEmit`), and runs its (currently single, T012-pulled-forward) vitest suite cleanly. `npm run build`'s prerender step does **not** yet succeed — `react-router.config.ts` unconditionally lists every `content-index.json` entry as a route to prerender, and no route matches until Phase 3 (T019) adds the detail route. This is expected, not a regression; tracked here so it isn't mistaken for a Phase 1 defect.

**Note**: T012 (Phase 2) was pulled forward and completed here — `vite.config.ts` (T002) imports `app/discovery/llms-txt.ts` at config-eval time, so Setup couldn't build at all without it already existing. See Phase 2 below; T012 is marked done, not re-done.

---

## Phase 2: Foundational (Schema validation + MCP/negotiation plumbing every story needs)

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T008 [P] Write failing test for `date` (parseable date string), `tags` (array of strings if present), and `project` (string if present) structural validation for site type `"log"` — malformed values rejected with a clear, per-file error, not silently accepted — in `internal/content/schema_test.go` (data-model.md "Project-Log Entry" validation rules; spec Edge Cases)
- [X] T009 Implement the structural validation in `internal/content/schema.go`'s `ParseItem` to make T008 pass (depends on T008)
- [X] T010 [P] Copy `templates/site/profile/app/mcp/{agent.ts,query.ts,server.ts}` into `templates/site/log/app/mcp/`, renaming the MCP server name to `"biblio-log-site"` and replacing the `registerProfileTools` import/call with a local no-op placeholder `registerLogTools(server: McpServer) {}` in `templates/site/log/app/mcp/tools/log.ts` — so `templates/site/log` builds and typechecks cleanly at every checkpoint through Phase 5; T032 replaces the placeholder body with the real implementation, it does not introduce the function (depends on T006) — also pulled in `tools/shared.ts` and `ui/chart-view.tsx` (production files only, no `.test.ts`/`.test.tsx`), since `server.ts`/`shared.ts` transitively need them to compile; their copied test files were dropped (see note below)
- [X] T011 [P] Copy `templates/site/profile/app/negotiation/{accept.ts,ansi.ts,markdown.ts,wasm.d.ts}` into `templates/site/log/app/negotiation/` unchanged (Principle II — negotiation is type-agnostic, FR-005) (depends on T006)
- [X] T012 [P] Copy `templates/site/profile/app/discovery/llms-txt.ts` (+ test) into `templates/site/log/app/discovery/` unchanged (depends on T006) — completed early in Phase 1, see note above
- [X] T013 Copy `templates/site/profile/workers/entry.ts` into `templates/site/log/workers/entry.ts`, wiring it to the log site's `mcp/agent.ts` and `negotiation/accept.ts` (depends on T010, T011)

**Note**: T010/T011 only call for copying production `.ts` files, not their profile `.test.ts` counterparts — correctly so. `shared.test.ts` and `accept.test.ts` were initially copied too (to keep parity-by-habit) but both assert against profile-specific seed content ("About Me", slug `about`), which doesn't exist on the log site; they were removed rather than adapted, since they were never an actual task. `chart-view.test.tsx` and `ansi.test.ts` are self-contained (their own fixtures, no real content-index.json reads) and were kept.

**Checkpoint** (verified via a real `biblio-cli new --type=log` smoke scaffold, including instantiating `createServer()` directly): Log site installs, typechecks, and passes its full vitest suite (3 files, 10 tests) with the MCP/negotiation/content-loading plumbing in place and log-specific schema validation enforced Go-side. Ready for entry views, filtering, and type-specific tools.

---

## Phase 3: User Story 1 — Author and publish a run of dated entries (Priority: P1) 🎯 MVP

**Goal**: Every authored dated/tagged/project-attributed entry appears correctly in the deployed site's entry list and detail view (FR-001, FR-002; spec Scenario 1)

**Independent Test**: `biblio-cli new --type=log`, author several dated entries (some sharing tags/project), `index && build && deploy`, confirm every entry's date/tags/project display correctly in the list and its own detail view — [quickstart.md](./quickstart.md) Scenario 1

### Tests for User Story 1 ⚠️

- [X] T014 [P] [US1] Write failing test for the entry-list route rendering each entry's date/tags/project in `templates/site/log/app/routes/home.test.tsx`
- [X] T015 [P] [US1] Write failing test for the entry-detail route rendering date/tags/project + body in `templates/site/log/app/routes/detail.test.tsx`
- [X] T016 [P] [US1] Write failing Go test that scaffolding `--type=log`, authoring entries, and running `index` produces `content-index.json` entries carrying correct `date`/`tags`/`project` fields, in `test/integration/scaffold_index_test.go` — extended the existing `TestScaffoldThenIndex_Log` (it already checked `date`) to also assert `tags` is present; `project` is optional per data-model.md and already covered negatively by Phase 2's `TestParseItem_LogProjectNotString`

### Implementation for User Story 1

- [X] T017 [US1] Implement `templates/site/log/app/routes/home.tsx` as the chronological entry-list view, displaying date/tags/project per entry (depends on T006, T014)
- [X] T018 [US1] Implement `templates/site/log/app/routes/detail.tsx` as the entry-detail view, displaying date/tags/project + rendered body (depends on T006, T015)
- [X] T019 [US1] Wire `home.tsx`/`detail.tsx` into `templates/site/log/app/routes.ts` (depends on T007, T017, T018)
- [X] T020 [US1] Add two more seed entries to `templates/site/log/content/posts/`, sharing a tag and a project with the existing `2026-06-01-hello-world.md`, to give Phase 4's filter tests real overlapping data (FR-001 seed-content contract) — added `2026-06-05-shipped-the-filter-ui.md` and `2026-06-10-refactored-the-content-index.md`, both tagged `shipping` and attributed to `project: biblio`; also updated the static `content-index.json` fixture to mirror them (it ships pre-baked with the template, not regenerated by `biblio-cli index` until the author's first edit)

**Checkpoint** (verified via a real `biblio-cli new --type=log` smoke scaffold, full `npm run build` including prerender): User Story 1 is independently functional — all 3 entries prerender correctly with date/tags/project visible on both the list and detail HTML output. `npx tsc --noEmit` clean, vitest 5 files/13 tests passing, Go suite green including the extended T016 integration test (SC-001 partial).

---

## Phase 4: User Story 2 — Filter the entry list (Priority: P2)

**Goal**: Filtering by date range, tag, and project (individually or combined, as an intersection) returns only genuinely matching entries, with a brand-new tag/project appearing in filters immediately after the next build (FR-003; spec Scenario 2, Edge Cases)

**Independent Test**: On the deployed site, filter by a date range, a tag, a project, and a date range + tag together; confirm the combined filter returns the intersection with no false positives/omissions; add a new-tag entry, rebuild, confirm it appears with no extra registration step — [quickstart.md](./quickstart.md) Scenario 2

### Tests for User Story 2 ⚠️

- [X] T021 [P] [US2] Write failing test for `filterEntries()` single-dimension filters (date-range only, tag only, project only) in `templates/site/log/app/content/filters.test.ts`
- [X] T022 [P] [US2] Write failing test for `filterEntries()` combined date-range + tag intersection semantics — no false positives, no omissions (spec Edge Cases) in the same file
- [X] T023 [P] [US2] Write failing test confirming a brand-new tag/project value present in `content-index.json` after a rebuild is immediately filterable with no separate registration step (spec Edge Cases) in the same file
- [X] T024 [P] [US2] Write failing test for the filter UI reflecting selected date-range/tag/project state and rendering only matching entries in `templates/site/log/app/routes/home.test.tsx`

### Implementation for User Story 2

- [X] T025 [US2] Implement `filterEntries(entries, query)` — date-range ∩ tag ∩ project intersection over `content-index.json` entries — in `templates/site/log/app/content/filters.ts` (depends on T021, T022, T023; this is the canonical filter-computation module shared by the UI (T026) and the MCP tools (Phase 6, T031))
- [X] T026 [US2] Add date-range/tag/project filter controls to `templates/site/log/app/routes/home.tsx`, wired to `filterEntries()` (depends on T017, T024, T025)

**Note**: T024 surfaced a real test-isolation bug — `vitest.setup.ts` (copied unchanged from profile in Phase 1) never called `@testing-library/react`'s `cleanup()` between tests. With only one `render()` per test file this was latent; `home.test.tsx`'s new multi-test filter-control suite (multiple `render()` calls in one file) made it produce false "multiple elements found" failures. Fixed by adding `afterEach(cleanup)` to `vitest.setup.ts`. Worth carrying back to `templates/site/profile/vitest.setup.ts` at some point, but that's outside this release's scope.

**Checkpoint** (verified via a real `biblio-cli new --type=log` smoke scaffold, full `npm run build`): User Story 2 is independently functional — filtering works correctly across all dimensions and combinations. `npx tsc --noEmit` clean, vitest 6 files/24 tests passing, Go suite green, full build (including the filter UI bundle) prerenders all 3 entries successfully (SC-001 complete).

---

## Phase 5: User Story 3 — Keep up a daily authoring rhythm (Priority: P3)

**Goal**: Adding a new day's entry uses the exact same authoring flow as the previous one, with zero one-off structural or config adjustments, repeatable across several consecutive entries (FR-004; spec Scenario 3, SC-002)

**Independent Test**: Author three consecutive dated entries one at a time, running only `biblio-cli index && build && deploy` between each, and confirm no diffs to `wrangler.jsonc`, `content/site.yaml`, or `package.json` are needed between cycles — [quickstart.md](./quickstart.md) Scenario 3

- [X] T027 [US3] Author and run an integration check that authors 3 consecutive dated entries one at a time (separate `index`/`build` cycles per entry) against a scaffolded log site, asserting zero changes to `wrangler.jsonc`/`content/site.yaml`/`package.json` between cycles, in `test/integration/pilot-site-log/` (depends on T017–T020; FR-004, SC-002)

**Checkpoint**: `TestAuthoringRhythm_ThreeConsecutiveEntries` passes — User Story 3 is independently demonstrated. Three separate days each add one dated entry via the real `index.Generate()` → `cli.RunBuild()` path (stubbed `vite build`, matching `pilot-site`'s own convention); `wrangler.jsonc`, `content/site.yaml`, and `package.json` are confirmed byte-for-byte identical across all three cycles — zero one-off adjustments (SC-002). Full Go suite green.

---

## Phase 6: User Story 4 — The agent/terminal layer works for this site type too (Priority: P4)

**Goal**: MCP tools expose timeline/tag/project queries over the log content shape, terminal/`curl` clients get clean Markdown/ANSI renderings, and a profile-shaped query is rejected with the same clear, structured error v0.03 already defines (FR-005; spec Scenario 4, Edge Cases)

**Independent Test**: Connect an MCP client to the deployed log site, `list-tools`/`call-tool` a log-specific tool and confirm a self-sufficient Markdown result; `curl -H "Accept: text/markdown"` and a terminal `User-Agent` both get clean renderings; calling a profile-only tool name (e.g. `get_profile`) returns the MCP SDK's standard unknown-tool error — [quickstart.md](./quickstart.md) Scenario 4

### Tests for User Story 4 ⚠️

- [X] T028 [P] [US4] Write failing test for log-specific MCP tools — `list_entries` (optional date-range/tag/project args, `filterEntries()`-backed) and `get_entry` by slug — content contract in `templates/site/log/app/mcp/tools/log.test.ts`
- [X] T029 [P] [US4] Write failing test confirming a profile-only tool name is not registered on the log site's `McpServer` and a call returns the standard unknown-tool error, not a silent empty result (spec Edge Cases) in `templates/site/log/app/mcp/server.test.ts`
- [X] T030 [P] [US4] Write failing test for `Accept`/`User-Agent`-negotiated Markdown/ANSI rendering of a log entry, re-exercising the T011 negotiation logic against log content shape, in `templates/site/log/app/negotiation/accept.test.ts` — the ansi tier itself is exercised generically by the unchanged, copied `ansi.test.ts` (its own mocking, content-independent); this file's job is `negotiatePageResponse` against real log content plus `negotiateVariant`'s curl→ansi selection, not re-proving WASM rendering already proven elsewhere

### Implementation for User Story 4

- [X] T031 [US4] Implement `list_entries` and `get_entry` tools in `templates/site/log/app/mcp/tools/log.ts` (depends on T025, T028)
- [X] T032 [US4] Replace T010's no-op `registerLogTools` placeholder body in `templates/site/log/app/mcp/tools/log.ts` with the real tool registration (`list_entries`, `get_entry`); `createServer()` in `templates/site/log/app/mcp/server.ts` already calls it from T010, no further wiring needed (depends on T010, T031)
- [X] T033 [US4] Verify the T012 `llms.txt` copy's output against log content (date/tags/project-aware key-page links); adjust only if T030 surfaces a real gap (depends on T012, T030) — verified via real build output; all 3 entries listed correctly, no gap, no change needed

**Checkpoint** (verified via a real `biblio-cli new --type=log` smoke scaffold, including a live MCP client exchange over `InMemoryTransport`): All user stories independently functional. `npx tsc --noEmit` clean, vitest 9 files/37 tests passing, Go suite green. A real MCP client lists `list_pages, get_page, search_content, list_entries, get_entry`, successfully filters via `list_entries({tag: "shipping"})`, and a `get_profile` call (profile-only tool) correctly returns `isError: true` — the agent/terminal layer generalizes to the project-log type (SC-001 complete via agent/terminal evidence).

---

## Phase 7: Polish & Cross-Cutting Concerns

- [X] T034 [P] Run the full [quickstart.md](./quickstart.md) Scenarios 1–4 against a real scaffolded, deployed log site; record results as this release's closing evidence (SC-001..SC-003) — **results**: ran against a real `biblio-cli new --type=log` scaffold (no real Cloudflare deploy available in this environment — same documented limitation as `test/integration/pilot-site`'s own suite). Scenario 1/3: authored 3 consecutive dated entries one at a time via the real binary; `wrangler.jsonc`/`content/site.yaml`/`package.json` md5sums identical across all 3 `biblio-cli index` cycles; real `biblio-cli build` pre-rendered all 6 entries (3 seed + 3 authored) with date/tags/project correctly attributed in the output HTML, including a brand-new tag (`brand-new-tag-here`) appearing with no registration step. Scenario 2: filter-narrowing behavior already proven via real `fireEvent.change` assertions in `home.test.tsx` (Phase 4) — not re-driven through a headless browser here. Scenario 4: a real MCP client over `InMemoryTransport` against the scaffolded project listed `list_pages, get_page, search_content, list_entries, get_entry`, successfully filtered `list_entries({tag: "brand-new-tag-here"})`, and `get_profile` (profile-shaped query) returned `isError: true`; terminal/`curl` Markdown negotiation already proven in `accept.test.ts` (Phase 6) — `negotiation/ansi.ts`'s real `.wasm` import only resolves under Wrangler's bundler or vitest's mocked setup, not plain Node, so it wasn't re-run bare here.
- [X] T035 [P] Update `specs/releases/0.04-release.md` status and this repo's `CLAUDE.md`/`AGENTS.md` SPECKIT markers to reflect completion — flipped `0.04-release.md`'s `**Status**` to `Shipped (2026-06-20)` (a new convention — v0.01–v0.03's release docs all still say `Draft`; documented the going-forward convention in `AGENTS.md`'s Release process section, per explicit user decision not to backfill it onto earlier releases). `CLAUDE.md`'s SPECKIT marker is left pointing at `specs/releases/0.04-project-log-site-type/plan.md` — it's still the correct "current plan" for this branch; v0.05 lives on a separate, not-yet-merged branch (`105-release-0-05-docs-site-and-consistency`) outside this branch's history. Also added `docs/releases/v0.04/README.md` and a `CHANGELOG.md` entry, mirroring v0.01–v0.03's existing user-guide convention (constitution Release Process: release notes + usage guide).
- [X] T036 [P] Run `gofmt`/`go vet` (`scripts/check.sh`) and the TS-side lint/typecheck across every file touched in this release — `scripts/check.sh` clean; no ESLint config exists for the site templates (consistent with Principle I — no added lint infra), so `npx tsc --noEmit` is the TS-side equivalent and is clean across `templates/site/log`.

---

## Dependencies & Execution Order

- **Phase 1 (Setup)** → **Phase 2 (Foundational)** → **Phase 3 (US1)** → **Phase 4 (US2)** → **Phase 5 (US3)** → **Phase 6 (US4)** → **Phase 7 (Polish)**
- US2 depends on US1's entry-list view existing (filters narrow what US1 already renders)
- US3 depends on US1's seed-content/build flow existing (proves repeatability of something that already works once)
- US4 depends on US2's `filterEntries()` (T025) for its `list_entries` tool, and on Phase 2's MCP/negotiation plumbing (T010–T013)
- Within each phase, tasks not marked `[P]` touch shared files/state and run sequentially in listed order; `[P]` tasks touch independent files and may run in parallel

## Parallel Execution Examples

- Phase 1: T001–T006 in parallel (independent files), T007 after T006
- Phase 2: T008 before T009; T010, T011, T012 in parallel (independent files) before T013
- Phase 3 tests: T014, T015, T016 in parallel before any Phase 3 implementation task
- Phase 4 tests: T021, T022, T023, T024 in parallel before T025/T026
- Phase 6 tests: T028, T029, T030 in parallel before T031–T033
- Phase 7: T034, T035, T036 in parallel

## Implementation Strategy

**MVP first**: Phase 1 → Phase 2 → Phase 3 (US1) delivers a real, deployed log site whose entries display correctly with date/tags/project — independently testable and demonstrable before any filtering, rhythm-repetition, or agent/terminal work exists.

**Incremental delivery**: Add Phase 4 (US2) for filtering, Phase 5 (US3) to prove the authoring rhythm generalizes with zero one-off adjustments, then Phase 6 (US4) to extend the v0.03 agent/terminal layer to this site type. Phase 7's quickstart run is the release-closing evidence for SC-001–SC-003.
