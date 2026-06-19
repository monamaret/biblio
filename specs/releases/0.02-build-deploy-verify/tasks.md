---

description: "Task list for v0.02 — Build, deploy, and verify a real site"
---

# Tasks: v0.02 — Build, deploy, and verify a real site

**Input**: Design documents from `/specs/releases/0.02-build-deploy-verify/`
**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/cli-commands.md](./contracts/cli-commands.md), [quickstart.md](./quickstart.md)

**Tests**: Included and REQUIRED — Constitution Principle VI (Test-First Development, NON-NEGOTIABLE) and plan.md's Constitution Check GATE explicitly require failing tests authored before implementation for the `build`/`deploy` CLI wrappers and the prerendered views, plus a release-closing end-to-end integration test (plan.md Testing; spec Edge Cases).

**Organization**: Tasks are grouped by user story, mapped from spec.md Scenarios 1–3 in their natural dependency order — build before deploy before live-visitor verification — to enable independent implementation and testing of each.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 = Scenario 1 (build), US2 = Scenario 2 (deploy), US3 = Scenario 3 (visitor browses live site)
- File paths are exact, taken from plan.md's Project Structure (Source Code) section

## Path Conventions

Extends the existing Go CLI module and v0.01-scaffolded combined-project template (plan.md "Structure Decision"):
- `internal/cli/` (Go, `biblio-cli` Cobra tree)
- `templates/site/<type>/` (TS/React combined-project template; this release's work is on `templates/site/profile/` per the single-pilot-site scope)
- `test/{unit,integration}/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the dependencies and config scaffolding that both the `build`/`deploy` CLI wrappers and the prerendered template need, before any user story work begins

- [ ] T001 Add React Router v7 (`react-router`, `@react-router/dev`) and content-rendering dependencies (`gray-matter`, `react-markdown`, `remark-gfm`) to `templates/site/profile/package.json`
- [ ] T002 [P] Author the base `react-router.config.ts` template (no `prerender` entries yet — wired in US1) in `templates/site/profile/react-router.config.ts`
- [ ] T003 [P] Author the generated CI workflow template invoked by `biblio-cli deploy` in `templates/site/profile/.github/workflows/deploy.yml` (research.md "Build/deploy command surface")

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared content-loading scaffolding that the prerendered routes (US1) and the live-rendered pages (US3) both depend on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T004 Implement the `import.meta.glob` + `gray-matter` content-loading module, reading the unchanged v0.01 `content-index.json` as its single source of truth (Principle II; data-model.md "Pre-rendered Page Set" `source-index`) in `templates/site/profile/app/content/load.ts`
- [ ] T005 [P] Implement the `react-markdown` + `remark-gfm` rendering component used by every detail view in `templates/site/profile/app/content/render.tsx`

**Checkpoint**: Foundation ready — content can be loaded from `content-index.json` and rendered; route/view implementation can begin

---

## Phase 3: User Story 1 — Build a site from authored content (Priority: P1) 🎯 MVP

**Goal**: `biblio-cli build` pre-renders every content-driven page of a scaffolded profile/CV project into static output with no missing or broken pages (FR-001; spec Scenario 1)

**Independent Test**: Run `biblio-cli build` against a scaffolded, authored profile project and confirm one pre-rendered page exists per `content-index.json` entry (profile, experience, skills, projects — list/detail/navigation); introduce a malformed-frontmatter file and confirm the build fails with a clear, per-file message instead of silently skipping the page — per [quickstart.md](./quickstart.md) Scenario 1 and [contracts/cli-commands.md](./contracts/cli-commands.md) `build`

### Tests for User Story 1 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T006 [P] [US1] Write failing test for the profile/experience/skills/projects list and detail route modules rendering authored content correctly in `test/unit/profile_routes_test.go` (or the TS-side equivalent colocated with the route modules, per plan.md Testing)
- [ ] T007 [P] [US1] Write failing test for `biblio-cli build` succeeding with zero missing/broken pages against a fully authored project in `internal/cli/build_test.go`
- [ ] T008 [P] [US1] Write failing test for `biblio-cli build` failing with a clear, per-file message on malformed authored content (spec Edge Cases) in `internal/cli/build_test.go`

### Implementation for User Story 1

- [ ] T009 [P] [US1] Implement the profile list/detail route modules (profile, experience, skills, projects) and shared navigation view in `templates/site/profile/app/routes/`
- [ ] T010 [US1] Wire `prerender` in `templates/site/profile/react-router.config.ts` to enumerate every `content-index.json` entry into a pre-rendered route, satisfying the Pre-rendered Page Set completeness invariant (depends on T002, T004, T009; data-model.md "Pre-rendered Page Set")
- [ ] T011 [US1] Implement `biblio-cli build` as a thin wrapper over `vite build` in `internal/cli/build.go`: invoking the build, surfacing per-file validation failures from `prerender` with a non-zero exit, and the contract's success output (depends on T010; contracts/cli-commands.md `build`)

**Checkpoint**: User Story 1 is independently functional — `biblio-cli build` produces a complete, gap-free static output for the pilot profile site (SC-002 partial)

---

## Phase 4: User Story 2 — Deploy and reach the site live (Priority: P2)

**Goal**: `biblio-cli deploy` publishes a built site to its configured public address from a single command, updating cleanly in place on re-run (FR-002; spec Scenario 2)

**Independent Test**: Run `biblio-cli deploy` against a successful build and confirm the site is reachable at its configured `wrangler.jsonc` address via `curl -I`; re-run deploy against the now-live site and confirm a clean in-place update with no visitor-facing disruption — per [quickstart.md](./quickstart.md) Scenario 2 and [contracts/cli-commands.md](./contracts/cli-commands.md) `deploy`

### Tests for User Story 2 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T012 [P] [US2] Write failing test for `biblio-cli deploy` succeeding against a built project and surfacing the configured public address in `internal/cli/deploy_test.go`
- [ ] T013 [P] [US2] Write failing test for `biblio-cli deploy` surfacing the underlying `wrangler deploy`/CI failure with diagnosable context on error (e.g. missing build artifact, auth failure) in `internal/cli/deploy_test.go`

### Implementation for User Story 2

- [ ] T014 [US2] Implement `biblio-cli deploy` as a thin wrapper over `wrangler deploy` in `internal/cli/deploy.go`: failing early if no build artifact exists, invoking the deploy against the existing `wrangler.jsonc` route binding, and the contract's success/failure output and exit codes (depends on T011; contracts/cli-commands.md `deploy`)
- [ ] T015 [US2] Wire the generated CI workflow (`templates/site/profile/.github/workflows/deploy.yml`, T003) to invoke `biblio-cli build` then `biblio-cli deploy` on push, matching the confirmed one-CLI-surface decision (research.md "Build/deploy command surface")

**Checkpoint**: User Story 2 is independently functional — a built pilot site can be published to and re-deployed at its configured public address (SC-002)

---

## Phase 5: User Story 3 — A visitor browses the live site (Priority: P3)

**Goal**: The deployed pilot site's profile, experience, skills, and projects views display the authored content correctly, styled per the v0.01-generated theme, when an ordinary visitor browses it (FR-003; spec Scenario 3)

**Independent Test**: Open the deployed pilot site's public address in a browser and confirm every view (profile, experience, skills, projects) displays authored content correctly and matches the `tokens.css` theme generated in v0.01 — per [quickstart.md](./quickstart.md) Scenario 3; this is the release-closing human-visitor check that produces the Pilot Site Validation Record (data-model.md)

### Tests for User Story 3 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T016 [P] [US3] Write failing test confirming every route module applies the unchanged `tokens.css` theme (no inline/duplicate styling path) in `test/unit/theme_application_test.go`

### Implementation for User Story 3

- [ ] T017 [US3] Verify and, if needed, adjust the route modules from T009 so each links `tokens.css` and renders without relying on any styling path other than the v0.01-generated theme (depends on T009, T016)
- [ ] T018 [US3] Author and run the end-to-end content → build → deploy → live-verification integration check against one real, fully authored pilot profile/CV site in `test/integration/pilot-site/` (FR-004; data-model.md "Pilot Site Validation Record"; constitution Release Process release-closing test)

**Checkpoint**: All user stories independently functional — the full content → build → deploy → live-verification path is demonstrated on one real pilot site (SC-001, SC-002, SC-003)

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Final consistency pass across the release

- [ ] T019 [P] Update `specs/releases/0.02-release.md` status and this release's `AGENTS.md`/CLAUDE.md SPECKIT markers to reflect completion, per the project's documentation-stays-in-sync convention
- [ ] T020 [P] Run `gofmt`/`go vet` (`scripts/check.sh`) and the TS-side lint/typecheck across all files touched in this release

---

## Dependencies & Execution Order

- **Phase 1 (Setup)** → **Phase 2 (Foundational)** → **Phase 3 (US1)** → **Phase 4 (US2)** → **Phase 5 (US3)** → **Phase 6 (Polish)**
- US2 depends on US1's build output existing (deploy has nothing to publish otherwise)
- US3 depends on US2's live deployment existing (nothing to browse otherwise)
- Within each phase, tasks not marked `[P]` touch shared files/state and must run sequentially in listed order; `[P]` tasks touch independent files and may run in parallel

## Parallel Execution Examples

- Phase 1: T002 and T003 in parallel (independent template files)
- Phase 2: T004 and T005 in parallel (independent modules), both before Phase 3
- Phase 3 tests: T006, T007, T008 in parallel before any Phase 3 implementation task
- Phase 4 tests: T012, T013 in parallel before T014
- Phase 6: T019 and T020 in parallel

## Implementation Strategy

**MVP first**: Phase 1 → Phase 2 → Phase 3 (US1) delivers a working `biblio-cli build` producing a complete pre-rendered static output — independently testable and demonstrable before any deploy capability exists.

**Incremental delivery**: Add Phase 4 (US2) to reach a publicly live site, then Phase 5 (US3) to close the loop with human-visitor verification and the release-closing pilot-site integration test (FR-004, SC-001). Each phase's checkpoint is independently demonstrable per spec SC-003 (no v0.03 agent/terminal capability required at any point).
