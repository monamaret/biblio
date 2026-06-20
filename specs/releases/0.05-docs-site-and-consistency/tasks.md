---

description: "Task list for v0.05 — Prove the platform on documentation and across all site types"
---

# Tasks: v0.05 — Prove the platform on documentation and across all site types

**Input**: Design documents from `/specs/releases/0.05-docs-site-and-consistency/`
**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/docs-site-and-review.md](./contracts/docs-site-and-review.md), [quickstart.md](./quickstart.md)

**Tests**: Included and REQUIRED — Constitution Principle VI and plan.md's Testing section require failing-test-first coverage for the documentation schema/hierarchy/TOC/tool/rendering units, the *absence* of filter UI/params (spec Scenario 1.2), and the cross-type consistency review itself run as a structured, repeatable check.

**Codebase reality check (read before starting)**: `templates/site/docs/` is still the thin v0.01-era plain Vite/React placeholder (`src/main.tsx`, `index.html`, no `app/` tree, plain `package.json`/`vite.config.ts`) — it never received the React Router v7 / Worker / MCP / negotiation buildout `templates/site/profile/` got in v0.02–v0.03, exactly the situation `templates/site/log/` was in before [v0.04](../0.04-project-log-site-type/tasks.md) Phase 1. Phase 1 below repeats that parity work for `docs`; it is **not** optional cleanup.

**Schema gap (read before starting Phase 2)**: `internal/content/schema.go`'s `requiredFields["docs"]` is currently `{"title", "order"}` — it does not yet require `section`, and there is no `validateDocsFields` analogous to `validateLogFields` (no structural check that `order` is actually numeric or `section` is actually a string). [data-model.md](./data-model.md)'s Documentation Page entity requires both `section` and `order`. Phase 2 closes this gap. The two existing seed pages (`templates/site/docs/content/{index,getting-started}.md`) also lack a `section` field today and must be updated once the schema requires one (Phase 3).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 = Scenario 1 (browse the docs site by hierarchy; filter UI/params confirmed absent), US2 = Scenario 2 (cross-type theming review), US3 = Scenario 3 (cross-type terminal-rendering review)

## Path Conventions

- `templates/site/docs/app/...`, `templates/site/docs/workers/entry.ts` — mirrors `templates/site/profile/app/` and `templates/site/log/app/`'s structure and naming (v0.04 precedent)
- `templates/site/docs/app/routes/{home,detail}.tsx` — flat under `app/routes/`, not nested under an `app/routes/pages/` subdirectory as plan.md's Project Structure sketch suggests; matches the flat layout `templates/site/profile/app/routes/` and `templates/site/log/app/routes/` both already use, so docs doesn't introduce a one-off routing convention
- `internal/content/schema.go`, `internal/content/schema_test.go` — docs-specific structural validation extends the existing shared schema; no new `internal/content/docs/` subpackage (one schema owner, Principle II) — same precedent as `validateLogFields`, despite plan.md's Project Structure sketch suggesting `internal/content/docs/`
- No `internal/scaffold/docs/` package either, despite plan.md's Project Structure sketch suggesting one for "documentation seed-content templates and TOC structure" — seed content lives directly under `templates/site/docs/content/` (T023, the same convention every other site type's seed content already uses) and TOC derivation is TS-only (`templates/site/docs/app/content/toc.ts`, T013), mirroring `templates/site/log/app/content/filters.ts`'s precedent of keeping a single, derived-view owner in TS rather than splitting it across a Go package and a TS module
- `templates/site/docs/app/content/toc.ts` — the TOC/hierarchy-derivation module, mirrors `templates/site/log/app/content/filters.ts`'s pattern: a single TS-only derived-view module shared by the browser UI and MCP tools, not duplicated Go-side (Principle II; data-model.md "Table of Contents / Navigation Structure")
- `test/integration/pilot-site-docs/` — matches the `test/integration/pilot-site-log/` convention from v0.04
- `test/integration/cross-type-review/` — new in this release; the closing artifact for US2/US3 (data-model.md "Cross-Type Consistency Record")

---

## Phase 1: Setup (Bring `templates/site/docs` to parity with `templates/site/profile`'s infra)

**Purpose**: Nothing docs-specific can be built until the docs template has the same React Router v7 / Vite / Worker / test scaffolding profile and log already have

- [X] T001 [P] Replace `templates/site/docs/package.json`'s plain Vite/React deps with the profile/log set (`react-router`, `@react-router/dev`, `agents`, `@modelcontextprotocol/sdk`, `zod`, `md4x`, `react-markdown`, `remark-gfm`, `isbot`, vitest/testing-library devDeps) per `templates/site/profile/package.json`
- [X] T002 [P] Replace `templates/site/docs/vite.config.ts` with the React Router v7 + `llmsTxtPlugin()` config from `templates/site/profile/vite.config.ts` (same rationale: no `@cloudflare/vite-plugin`)
- [X] T003 [P] Add `templates/site/docs/react-router.config.ts`, `templates/site/docs/tsconfig.json`, `templates/site/docs/vitest.config.ts`, `templates/site/docs/vitest.setup.ts` mirroring profile's/log's (type-agnostic, no docs-specific changes)
- [X] T004 [P] Remove the superseded v0.01 plain-Vite files `templates/site/docs/src/main.tsx` and `templates/site/docs/index.html`
- [X] T005 [P] Add `templates/site/docs/.github/workflows/deploy.yml` (copy of `templates/site/profile/.github/workflows/deploy.yml` — CI/deploy contract is type-agnostic, FR-001)
- [X] T006 [P] Copy `templates/site/profile/app/content/load.ts` and `render.tsx` into `templates/site/docs/app/content/` unchanged (Principle II — content loading/rendering is type-agnostic; data-model.md "Documentation Page" extends the same Authored Content Item) — also added a `templates/site/docs/content-index.json` fixture (profile/log both ship one; docs didn't yet)
- [X] T007 [P] Copy `templates/site/profile/app/discovery/llms-txt.ts` (+ test) into `templates/site/docs/app/discovery/` unchanged (depends on T006) — pulled forward into Setup since `vite.config.ts` (T002) imports it at config-eval time, matching v0.04's T012 correction
- [X] T008 Author `templates/site/docs/app/root.tsx` and `templates/site/docs/app/routes.ts` (route skeleton only: `index` and a `"*"` detail splat, matching `templates/site/profile/app/routes.ts`; the TOC navigation replaces the plain entry list in Phase 3) — ship `routes.ts` as a valid *empty* `RouteConfig` so it typechecks at this checkpoint; Phase 3 (T020) fills it in

**Checkpoint** (verified via a real `biblio-cli new --type=docs` smoke scaffold): `templates/site/docs` installs, typechecks (`tsc --noEmit`), and runs its (currently single, T007-pulled-forward) vitest suite cleanly. `npm run build`'s prerender step does **not** yet succeed — no route matches until Phase 3 (T020) adds the detail route. Expected, not a regression, per v0.04's identical Phase 1 note.

---

## Phase 2: Foundational (Schema validation + TOC derivation + MCP/negotiation plumbing every story needs)

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [X] T009 [P] Write failing test for `section` (required, non-empty string) and `order` (required, must parse as an integer) structural validation for site type `"docs"` — malformed/missing values rejected with a clear, per-file error, not silently accepted — in `internal/content/schema_test.go` (data-model.md "Documentation Page" validation rules; spec Edge Cases)
- [X] T010 Add `"section"` to `requiredFields["docs"]` and implement `validateDocsFields` (mirroring `validateLogFields`'s shape: `order` must parse as an integer, `section` must be a non-empty string) in `internal/content/schema.go`'s `ParseItem`, to make T009 pass (depends on T009) — also pulled forward adding `section` to the two existing seed pages and the `content-index.json` fixture (originally T023), since the schema change otherwise broke `TestScaffoldThenIndex_Docs` immediately
- [X] T011 [P] Write failing test for `toc()` deterministically grouping entries into `sections[]` (each an ordered list of pages) from their `section`/`order` fields, including the documented default for a page whose `section` doesn't cleanly resolve to one slot (primary-section assignment — spec Edge Cases, data-model.md "Table of Contents / Navigation Structure") in `templates/site/docs/app/content/toc.test.ts`
- [X] T012 [P] Write failing test asserting `toc()`'s output and the routes/components built on it expose no `tags`/`project` fields or filter-query parameters anywhere (spec Scenario 1.2; contracts/docs-site-and-review.md "Filter UI") in `templates/site/docs/app/content/no-filter-surface.test.ts` — not `test/unit/no-filter-surface/`, since that root-level directory has no Node/vitest tooling (no root `package.json`); a standalone `.ts` file there would never actually run. Placed inside the docs template's own vitest suite instead, alongside `toc.test.ts`
- [X] T013 Implement `toc(entries)` — deterministic `section`/`order` → `sections[]` grouping, with the documented primary-section default for ambiguous pages — in `templates/site/docs/app/content/toc.ts` (depends on T006, T011, T012; this is the canonical TOC-computation module shared by the UI (T021) and the MCP tools (Phase 3, T026))
- [X] T014 [P] Copy `templates/site/profile/app/mcp/{agent.ts,query.ts,server.ts}` into `templates/site/docs/app/mcp/`, renaming the MCP server name to `"biblio-docs-site"` and replacing the `registerProfileTools` import/call with a local no-op placeholder `registerDocsTools(server: McpServer) {}` in `templates/site/docs/app/mcp/tools/docs.ts` — so `templates/site/docs` builds and typechecks cleanly at every checkpoint; T027 replaces the placeholder body with the real implementation, it does not introduce the function (depends on T006) — also pull in `tools/shared.ts` and `ui/chart-view.tsx` (production files only, no `.test.ts`/`.test.tsx`, per v0.04's T010 precedent — their copied profile-specific tests assert against profile seed content that doesn't exist here)
- [X] T015 [P] Copy `templates/site/profile/app/negotiation/{accept.ts,ansi.ts,markdown.ts,wasm.d.ts}` into `templates/site/docs/app/negotiation/` unchanged (Principle II — negotiation is type-agnostic, FR-004) (depends on T006)
- [X] T016 Copy `templates/site/profile/workers/entry.ts` into `templates/site/docs/workers/entry.ts`, wiring it to the docs site's `mcp/agent.ts` and `negotiation/accept.ts` (depends on T014, T015)

**Checkpoint** (verified via a real `biblio-cli new --type=docs` smoke scaffold, including instantiating `createServer()` directly): Docs site installs, typechecks, and passes its full vitest suite with the MCP/negotiation/content-loading plumbing in place, `toc()` derivation proven (including the no-filter-surface assertion), and docs-specific schema validation enforced Go-side. Ready for hierarchy navigation and type-specific tools.

---

## Phase 3: User Story 1 — Browse a documentation site by its hierarchy (Priority: P1) 🎯 MVP

**Goal**: Every authored page is reachable and renders correctly through the TOC/section/order structure alone, with no tag/project filtering anywhere in the UI or query surface (FR-001, FR-002; spec Scenario 1)

**Independent Test**: `biblio-cli new --type=docs`, author several pages across two or more sections, `index && build && deploy`, confirm the TOC reflects the section/order structure with every page reachable and rendering correctly, and confirm no tag/project filter UI or query params exist anywhere — [quickstart.md](./quickstart.md) Scenario 1

### Tests for User Story 1 ⚠️

- [X] T017 [P] [US1] Write failing test for the TOC/home route rendering each section's pages in `order`, with no filter controls present, in `templates/site/docs/app/routes/home.test.tsx`
- [X] T018 [P] [US1] Write failing test for the page-detail route rendering its section/order context + body in `templates/site/docs/app/routes/detail.test.tsx`
- [X] T019 [P] [US1] Write failing Go test that scaffolding `--type=docs`, authoring pages across two sections, and running `index` produces `content-index.json` entries carrying correct `section`/`order` fields and never a `tags`/`project` field, in `test/integration/scaffold_index_test.go` — extends the existing `TestScaffoldThenIndex_Docs` (it already checks `order`) to also assert `section` is present and `tags`/`project` are absent

### Implementation for User Story 1

- [X] T020 [US1] Implement `templates/site/docs/app/routes/home.tsx` as the TOC/navigation view, rendering `toc()`'s `sections[]` with no filter controls (depends on T013, T017)
- [X] T021 [US1] Implement `templates/site/docs/app/routes/detail.tsx` as the page-detail view, displaying section/order context + rendered body (depends on T006, T018)
- [X] T022 [US1] Wire `home.tsx`/`detail.tsx` into `templates/site/docs/app/routes.ts` (depends on T008, T020, T021)
- [X] T023 [US1] Add a `section` field to the two existing seed pages `templates/site/docs/content/{index,getting-started}.md` (closing the schema gap noted above) and one more seed page in a second section, so Phase 3's TOC tests and Phase 4/5's cross-type review have a real multi-section hierarchy to exercise (FR-001 seed-content contract) — `section`/`getting-started` were actually pulled forward into T010 (otherwise the schema change broke existing tests immediately); this task added the third page, `templates/site/docs/content/configuration.md` (section "Reference")
- [X] T024 [P] [US1] Write failing test confirming a log-shaped tag-filter parameter sent to the docs-specific `list_sections`/`get_doc_page` MCP tools is rejected with a clear, structured error (the same documented-default contract as v0.04's type-mismatch handling), never silently ignored, in `templates/site/docs/app/mcp/tools/docs.test.ts`
- [X] T025 [P] [US1] Write failing test confirming a profile/log-only tool name is not registered on the docs site's `McpServer` and a call returns the standard unknown-tool error in `templates/site/docs/app/mcp/server.test.ts`
- [X] T026 [US1] Implement `list_sections` and `get_doc_page` tools in `templates/site/docs/app/mcp/tools/docs.ts`, backed by `toc()` (depends on T013, T024); `get_doc_page` is named distinctly from the shared `get_page` tool (`templates/site/profile/app/mcp/tools/shared.ts`, registered on every site type including docs) to avoid a duplicate-tool-name registration collision — same naming precedent as v0.04's `get_entry`; `list_sections` rejects a log-shaped filter argument per T024's documented default
- [X] T027 [US1] Replace T014's no-op `registerDocsTools` placeholder body in `templates/site/docs/app/mcp/tools/docs.ts` with the real tool registration (`list_sections`, `get_doc_page`); `createServer()` in `templates/site/docs/app/mcp/server.ts` already calls it from T014, no further wiring needed (depends on T014, T026)
- [X] T028 [US1] Verify the T007 `llms.txt` copy's output against docs content (section/order-aware key-page links); adjust only if a real gap surfaces (depends on T007) — verified via real build output; all 3 pages listed correctly, no gap, no change needed

**Checkpoint** (verified via a real `biblio-cli new --type=docs` smoke scaffold, full `npm run build` including prerender, plus a live MCP client exchange over `InMemoryTransport`): User Story 1 is independently functional — every authored page prerenders and is reachable via the TOC alone, no filter UI/params exist anywhere (`no-filter-surface` suite green), and a real MCP client lists `list_pages, get_page, search_content, list_sections, get_doc_page` (no name collision — `get_doc_page` is distinct from the shared `get_page`), successfully calls `list_sections`, a log-shaped tag-filter argument is rejected with a structured error, and a `get_entry`/`list_entries` (log-only) call returns `isError: true`. `npx tsc --noEmit` clean, Go suite green (SC-001 complete).

---

## Phase 4: User Story 2 — Compare theming across all three site types (Priority: P2)

**Goal**: With all three site types now live, each displays its own distinct, correctly-applied visual theme with zero cross-type bleed-through (FR-003; spec Scenario 2)

**Independent Test**: Scaffold (or use already-live) profile, log, and docs sites side by side; confirm each renders its own preset's tokens with no cross-type bleed-through — [quickstart.md](./quickstart.md) Scenario 2

### Tests for User Story 2 ⚠️

- [X] T029 [US2] Write a structured comparison test that scaffolds all three site types, generates each one's `tokens.css` via `theme.GenerateTokensCSS`, and asserts (a) each site's tokens.css matches its own preset exactly, and (b) the tokens each preset intends to be genuinely distinguishing per type (`color-fg`, `color-accent`) never collide across types — in `test/integration/cross-type-review/theming_review_test.go` (data-model.md "Cross-Type Consistency Record" theming-check; FR-003, SC-002). Deliberately does *not* assert every token differs across all three types: `color-bg` is intentionally shared between profile and docs (two light themes vs. log's dark one), and `font-heading` is intentionally shared across all three — both legitimate design choices, not bleed-through; an earlier draft of this test asserted `color-bg` must differ and produced a false positive, corrected before this was marked done
- [X] T030 [US2] Run the real side-by-side theming comparison (T029) against all three scaffolded site types; if it surfaces an inconsistency, file and fix it as a defect against the shared `internal/theme` package (never a per-site override — research.md "How cross-type inconsistencies are triaged") (depends on T029) — ran clean; no defect found in `internal/theme` itself, only in the test's own over-strict assertion (fixed in T029)

**Checkpoint**: `theming_review_test.go` passes with zero open inconsistencies across profile/log/docs — User Story 2 is independently demonstrated (SC-002 theming half complete).

---

## Phase 5: User Story 3 — Compare terminal rendering across all three site types (Priority: P3)

**Goal**: With all three site types now live, each returns a clean, correctly rendered, type-appropriate terminal/Markdown view on request (FR-004; spec Scenario 3)

**Independent Test**: Request each of the three deployed/scaffolded sites with `Accept: text/markdown` (and a terminal `User-Agent`); confirm each returns a clean, type-appropriate rendering — [quickstart.md](./quickstart.md) Scenario 3

### Tests for User Story 3 ⚠️

- [X] T031 [US3] Write a structured comparison test in `test/integration/cross-type-review/terminal_rendering_review_test.go` (data-model.md "Cross-Type Consistency Record" terminal-rendering-check; FR-004, SC-002). Reframed from the original plan during authoring: a single Go test cannot directly invoke `negotiatePageResponse`/`negotiateVariant` (those are TS, and the three templates are three independent npm projects with no shared workspace — a cross-template Node script wouldn't have a resolvable `node_modules`). Instead this Go test asserts the two structural guarantees that actually back "each site renders a clean, type-appropriate view": (a) `accept.ts`/`ansi.ts`/`wasm.d.ts` are byte-identical across all three scaffolded site types (one negotiation owner, Principle II), and (b) each type's seed content carries exactly its own content-shape fields with no cross-contamination. The TS-side negotiation behavior itself is exercised per-type by each template's own `app/negotiation/accept.test.ts` (profile/log already had one; this task also added the previously-missing `templates/site/docs/app/negotiation/{accept,ansi}.test.ts`, mirroring log's v0.04 precedent)

### Implementation for User Story 3

- [X] T032 [US3] Run the real side-by-side terminal-rendering comparison (T031) against all three scaffolded site types; if it surfaces an inconsistency, file and fix it as a defect against the shared negotiation/rendering code (never a per-site override) (depends on T031) — ran clean; no defect found
- [X] T033 [US3] Produce the Cross-Type Consistency Record as `test/integration/cross-type-review/CONSISTENCY-RECORD.md`, structured per data-model.md's four fields — `sites[]` (the three live deployments), `theming-check` (per-site result from T030), `terminal-rendering-check` (per-site result from T032), `inconsistencies-found[]` (empty, or each entry linked to the shared-code defect it was routed to and fixed as) — as the closing artifact SC-002 is checked against (depends on T030, T032)

**Checkpoint**: `terminal_rendering_review_test.go` passes with zero open inconsistencies across profile/log/docs — User Story 3 is independently demonstrated, and the Cross-Type Consistency Record shows SC-002 fully satisfied.

---

## Phase 6: Polish & Cross-Cutting Concerns

- [X] T034 [P] Run the full [quickstart.md](./quickstart.md) Scenarios 1–3 against a real scaffolded, deployed docs site plus the already-live profile/log sites; record results as this release's closing evidence (SC-001..SC-003) — **results**: ran against a real `biblio-cli new --type=docs` scaffold (no real Cloudflare deploy available in this environment — same documented limitation as `test/integration/pilot-site`'s own suite). Scenario 1: real `biblio-cli index && build` pre-rendered all 3 pages (`/`, `/index`, `/getting-started`, `/configuration`) with section/order-correct TOC grouping in the output HTML; grepped the rendered HTML for "tag"/"filter"/"project" and found none. Scenario 2: cross-type theming proven via `theming_review_test.go` (Phase 4) — each type renders its own preset exactly, zero bleed-through. Scenario 3: cross-type terminal rendering proven via each template's own `accept.test.ts`/`ansi.test.ts` (mocked WASM, vitest) plus `terminal_rendering_review_test.go`'s structural checks (Phase 5); a direct plain-Node negotiation call hits the same documented `.wasm`-resolution limitation v0.04's T034 already hit (`ansi.ts`'s real `.wasm` import only resolves under Wrangler's bundler or vitest's mocked setup), so it wasn't re-run bare here either
- [X] T035 [P] Update `specs/releases/0.05-release.md` status and this repo's `CLAUDE.md`/`AGENTS.md` SPECKIT markers to reflect completion, following the `**Status**: Shipped (<date>)` convention established in v0.04 (T035) — flipped `0.05-release.md`'s `**Status**` to `Shipped (2026-06-19)`; `CLAUDE.md`'s SPECKIT marker already correctly pointed at `specs/releases/0.05-docs-site-and-consistency/plan.md` (set during this branch's merge with main); `AGENTS.md` has no per-release plan pointer to update. Also added `docs/releases/v0.05/README.md` and a `CHANGELOG.md` entry, mirroring v0.01–v0.04's existing convention
- [X] T036 [P] Run `gofmt`/`go vet` (`scripts/check.sh`) and `npx tsc --noEmit` across every file touched in this release (`templates/site/docs`, `internal/content/schema.go`, `test/integration/cross-type-review/`) — `scripts/check.sh` clean; `npx tsc --noEmit` clean across `templates/site/docs`; full `go test ./...` green

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies - can start immediately
- **Foundational (Phase 2)**: Depends on Setup completion - BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends on Foundational; produces the live docs site US2/US3 review against
- **User Story 2 (Phase 4)** and **User Story 3 (Phase 5)**: Both depend on User Story 1's docs site existing (the third leg of the three-way comparison) — they are independent of each other and may run in parallel
- **Polish (Phase 6)**: Depends on all of Phases 3–5 being complete

### Within Each User Story

- Tests MUST be written and FAIL before implementation
- TOC/schema before routes; routes before MCP tools
- Story complete before moving to next priority

### Parallel Opportunities

- All Setup tasks marked [P] can run in parallel
- All Foundational tasks marked [P] can run in parallel (within Phase 2)
- US2 (Phase 4) and US3 (Phase 5) can run in parallel once US1 (Phase 3) completes
- Tests for a user story marked [P] can run in parallel

---

## Parallel Example: User Story 1

```bash
# Launch all tests for User Story 1 together:
Task: "Failing test for TOC/home route rendering in templates/site/docs/app/routes/home.test.tsx"
Task: "Failing test for page-detail route rendering in templates/site/docs/app/routes/detail.test.tsx"
Task: "Failing Go test for section/order content-index fields in test/integration/scaffold_index_test.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL - blocks all stories)
3. Complete Phase 3: User Story 1 — the docs site itself, live
4. **STOP and VALIDATE**: Test User Story 1 independently (SC-001)
5. Deploy/demo if ready

### Incremental Delivery

1. Complete Setup + Foundational → Foundation ready
2. Add User Story 1 → Test independently → Deploy/Demo (MVP! — SC-001)
3. Add User Story 2 (theming review) → Test independently
4. Add User Story 3 (terminal-rendering review) → Test independently
5. Phases 4 and 5 together close SC-002; Phase 6 closes SC-003

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story for traceability
- Verify tests fail before implementing
- Commit after each task or logical group
- Stop at any checkpoint to validate story independently
- Avoid: vague tasks, same file conflicts, cross-story dependencies that break independence
- Any cross-type inconsistency found in Phases 4–5 is fixed in the shared `internal/theme`/negotiation code, never patched per-site (Principle II; research.md)
