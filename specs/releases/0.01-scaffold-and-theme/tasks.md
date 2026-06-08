---

description: "Task list for v0.01 — Scaffold and theme a site"
---

# Tasks: v0.01 — Scaffold and theme a site

**Input**: Design documents from `/specs/releases/0.01-scaffold-and-theme/`
**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/cli-commands.md](./contracts/cli-commands.md), [quickstart.md](./quickstart.md)

**Tests**: Included and REQUIRED — Constitution Principle VI (Test-First Development, NON-NEGOTIABLE) and plan.md's Constitution Check GATE explicitly require a failing Go test authored before implementation for every generator/command in this release.

**Organization**: Tasks are grouped by user story (mapped from spec.md Scenarios 1–3, in their natural dependency order — scaffold before index before preview) to enable independent implementation and testing of each.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 = Scenario 1 (scaffold), US2 = Scenario 2 (index), US3 = Scenario 3 (preview)
- File paths are exact, taken from plan.md's Project Structure (Source Code) section

## Path Conventions

Single Go module/CLI project (plan.md "Structure Decision: Option 1 shape, CLI-flavored"):
- `cmd/biblio-cli/`, `internal/{cli,tui,scaffold,content,index,theme,devserver}/`, `templates/site/`, `test/{unit,integration}/`

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Initialize the Go module and project skeleton exactly as laid out in plan.md's Project Structure

- [X] T001 Create the directory skeleton (`cmd/biblio-cli/`, `internal/{cli,tui,scaffold,content,index,theme,devserver}/`, `templates/site/{profile,log,docs}/`, `test/{unit,integration}/`) per plan.md Project Structure
- [X] T002 Initialize the Go module (`go.mod`, Go 1.22+) and add Cobra, Viper, Bubble Tea, Lip Gloss, and Glamour dependencies
- [X] T003 [P] Author `.goreleaser.yaml` for cross-platform (macOS/Linux/Windows) binary build/release config (FR-006, research.md "Distribution method")
- [X] T004 [P] Configure Go formatting/vetting (`gofmt`, `go vet`) as a repeatable local check (e.g. `Makefile` or `scripts/check.sh`)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared CLI/TUI scaffolding that `new`, `index`, and `dev` all hang off of

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 Implement the Cobra root command (`biblio-cli`, `--version`) with Viper config loading in `internal/cli/root.go`
- [ ] T006 [P] Implement shared Bubble Tea/Lip Gloss TUI base styles and prompt helpers (used by every guided flow) in `internal/tui/styles.go`
- [ ] T007 [P] Implement the embedded scaffold-template asset loader (Go `embed` over `templates/site/`) in `internal/scaffold/assets.go`

**Checkpoint**: Foundation ready — `biblio-cli --version` runs, shared TUI styles and template-asset loading exist; user story implementation can begin

---

## Phase 3: User Story 1 — Scaffold a new site from one command (Priority: P1) 🎯 MVP

**Goal**: `biblio-cli new <project-name> --type=<profile|log|docs> --domain=<domain> --theme=<preset>` (or its guided TUI equivalent) produces a complete, runnable, themed combined site project (FR-001, FR-002, FR-004; spec Scenario 1)

**Independent Test**: Run `biblio-cli new my-site --type=profile --domain=example.com --theme=profile-default` and confirm `content/`, `content/site.yaml` (with `mcp_enabled: true`), `wrangler.jsonc` (route pre-filled from `--domain`), and `tokens.css` (matching the preset) all exist and are correctly populated — per [quickstart.md](./quickstart.md) Scenario 1 and [contracts/cli-commands.md](./contracts/cli-commands.md) `biblio-cli new`

### Tests for User Story 1 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T008 [P] [US1] Write failing test for Theme Preset definitions + deterministic `tokens.css` generation in `internal/theme/theme_test.go`
- [ ] T009 [P] [US1] Write failing test for the scaffold generator (project-dir creation, seeded `content/`, `site.yaml`, `wrangler.jsonc`, collision detection) in `internal/scaffold/generator_test.go`
- [ ] T010 [P] [US1] Write failing test for `new` command flag validation and error paths (unknown `--type`, `--theme` not valid for `--type`, colliding `project-name`) in `internal/cli/new_test.go`

### Implementation for User Story 1

- [ ] T011 [P] [US1] Define the Theme Preset entity and one ready-to-use preset per site type (`profile`, `log`, `docs`) in `internal/theme/presets.go` (data-model.md "Theme Preset"; FR-004)
- [ ] T012 [US1] Implement deterministic `tokens.css` generation from a selected Theme Preset in `internal/theme/tokens.go` (depends on T011; research.md "Theming ownership and artifact" — generated once, owned solely by `biblio-cli`)
- [ ] T013 [P] [US1] Author seed content and scaffold file templates for the three site types in `templates/site/profile/`, `templates/site/log/`, `templates/site/docs/`
- [ ] T014 [US1] Implement the scaffold generator producing the Site Project entity (`content/`, `content/site.yaml` with `mcp_enabled: true`, `wrangler.jsonc` with `--domain`-derived routing, `tokens.css`) in `internal/scaffold/generator.go` (depends on T007, T012, T013; data-model.md "Site Project")
- [ ] T015 [US1] Implement the Bubble Tea guided `new` flow (interactive type/domain/theme prompts, theme choices filtered to the selected type) in `internal/tui/new_flow.go` (depends on T006, T011)
- [ ] T016 [US1] Wire `biblio-cli new <project-name>` in `internal/cli/new.go`: flag parsing, TUI fallback when flags are omitted, generator invocation, and the contract's success/failure output and exit codes (depends on T014, T015; contracts/cli-commands.md `biblio-cli new`)

**Checkpoint**: User Story 1 is independently functional — `biblio-cli new` produces a complete, themed, runnable project for any of the three site types (SC-001, SC-002 partial)

---

## Phase 4: User Story 2 — Index authored content (Priority: P2)

**Goal**: `biblio-cli index [path]` produces the single canonical `content-index.json` from an authored `content/` tree (FR-003; spec Scenario 2)

**Independent Test**: Run `biblio-cli index` against a scaffolded project's `content/` and confirm `content-index.json` has exactly one correctly-sluggged entry per authored item; introduce a file with malformed frontmatter and confirm a clear, actionable per-file error (not a silently incomplete index) — per [quickstart.md](./quickstart.md) Scenario 2

### Tests for User Story 2 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T017 [P] [US2] Write failing test for the canonical frontmatter schema parser/validator across all three site types' field shapes (including malformed/missing-field error surfacing) in `internal/content/schema_test.go`
- [ ] T018 [P] [US2] Write failing test for slug derivation from content path (stable, collision-free within a site) in `internal/content/slug_test.go`
- [ ] T019 [P] [US2] Write failing test for content-index generation (one entry per item, correct `slug`/`path`, per-file error listing on validation failure) in `internal/index/generator_test.go`

### Implementation for User Story 2

- [ ] T020 [P] [US2] Implement the canonical Go-side frontmatter schema parser/validator (the single source of truth per Principle II — never duplicated on the TS side) in `internal/content/schema.go` (data-model.md "Authored Content Item")
- [ ] T021 [P] [US2] Implement slug derivation from content path per the confirmed layout convention in `internal/content/slug.go`
- [ ] T022 [US2] Implement the Content Index generator producing `content-index.json` (`entries[]`, shape consumable by both `prerender` (v0.02) and the MCP query layer (v0.03)) in `internal/index/generator.go` (depends on T020, T021; data-model.md "Content Index"; research.md "Content-index shape and ownership")
- [ ] T023 [US2] Wire `biblio-cli index [path]` in `internal/cli/index.go`: defaulting `path` to the project's `content/`, writing `content-index.json`, and reporting every failing file (not just the first) on validation error (depends on T022; contracts/cli-commands.md `biblio-cli index`)

**Checkpoint**: User Stories 1 AND 2 both work independently — a scaffolded project can now also produce a structurally correct canonical index (SC-002 complete)

---

## Phase 5: User Story 3 — Preview locally (Priority: P3)

**Goal**: `biblio-cli dev` starts the site and its supporting Worker together behind one command (FR-005; spec Scenario 3)

**Independent Test**: Run `biblio-cli dev` in a scaffolded project and confirm both `vite dev` and `wrangler dev` start as one supervised pair, their addresses are surfaced, and the seeded content is visible live in a browser, styled per the chosen theme — per [quickstart.md](./quickstart.md) Scenario 3

### Tests for User Story 3 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T024 [P] [US3] Write failing test for the dev-server supervisor (starts both processes, surfaces both addresses, reports which process failed and why on startup failure, clean shutdown on interrupt) in `internal/devserver/supervisor_test.go`

### Implementation for User Story 3

- [ ] T025 [US3] Implement the supervisor wrapping and forwarding output for `vite dev` + `wrangler dev` as one supervised process pair in `internal/devserver/supervisor.go`
- [ ] T026 [US3] Wire `biblio-cli dev` in `internal/cli/dev.go`: launching the supervisor, surfacing addresses, and the contract's failure reporting (port collision, missing `wrangler.jsonc`) and exit-on-interrupt behavior (depends on T025; contracts/cli-commands.md `biblio-cli dev`)

**Checkpoint**: All three user stories are independently functional — the full `new` → `index` → `dev` path now runs end to end (SC-001, SC-003)

---

## Phase 6: Polish & Cross-Cutting Concerns

**Purpose**: Cross-story wiring, the binary entrypoint, distribution, and final validation against the spec's Success Criteria

- [ ] T027 [P] Write end-to-end "scaffold → index → dev" integration tests for each of the three site types in `test/integration/`
- [ ] T028 [P] Wire `cmd/biblio-cli/main.go` as the binary entrypoint invoking the root command from `internal/cli`
- [ ] T029 [P] Document install/usage (GitHub Releases binary download and `go install`) per FR-006 and research.md "Distribution method"
- [ ] T030 Run [quickstart.md](./quickstart.md) end to end for all three site types and confirm SC-001..SC-003

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS all user stories
- **User Story 1 (Phase 3)**: Depends only on Foundational — this is the MVP
- **User Story 2 (Phase 4)**: Depends only on Foundational; in practice exercised against a project produced by US1, but its own implementation (`internal/content`, `internal/index`, `internal/cli/index.go`) shares no files with US1 and is independently testable against any `content/` tree, scaffolded or hand-built
- **User Story 3 (Phase 5)**: Depends only on Foundational; independently testable against any scaffolded project directory (its own files — `internal/devserver/*`, `internal/cli/dev.go` — don't overlap US1/US2)
- **Polish (Phase 6)**: Depends on US1, US2, and US3 all being complete (T027/T030 exercise the full chain)

### Within Each User Story

- Tests are written and confirmed FAILING before implementation (Principle VI, NON-NEGOTIABLE)
- Entities/presets before generators/services that consume them
- Generators/services before the Cobra command that wires them
- Story checkpoint reached only once its CLI command is fully wired and contract-passing

### Parallel Opportunities

- T003, T004 (Setup) in parallel
- T006, T007 (Foundational) in parallel — different files, T005 can also proceed independently
- All three US1 test tasks (T008, T009, T010) in parallel; T011 and T013 in parallel
- All three US2 test tasks (T017, T018, T019) in parallel; T020 and T021 in parallel
- Once Foundational is done, US1, US2, and US3 implementation can all proceed in parallel — they touch disjoint package trees (`internal/{scaffold,theme,tui}` vs `internal/{content,index}` vs `internal/devserver`) and disjoint CLI command files
- T027, T028, T029 (Polish) in parallel

---

## Parallel Example: User Story 1

```bash
# Tests first, together:
Task: "Write failing test for Theme Preset + tokens.css generation in internal/theme/theme_test.go"
Task: "Write failing test for scaffold generator in internal/scaffold/generator_test.go"
Task: "Write failing test for `new` command flag validation in internal/cli/new_test.go"

# Then entity/template work, together:
Task: "Define Theme Preset entity + per-type presets in internal/theme/presets.go"
Task: "Author seed content + scaffold templates in templates/site/{profile,log,docs}/"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (blocks everything else)
3. Complete Phase 3: User Story 1 — `biblio-cli new` produces a complete, themed, runnable project
4. **STOP and VALIDATE**: run [quickstart.md](./quickstart.md) Scenario 1 for all three `--type` values (SC-002 partial)

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. Add User Story 1 → validate independently → MVP demonstrable (scaffold + theme)
3. Add User Story 2 → validate independently → canonical index demonstrable
4. Add User Story 3 → validate independently → full local-preview loop demonstrable
5. Polish → run full [quickstart.md](./quickstart.md), confirming SC-001..SC-003 end to end — independent of v0.02/v0.03 (spec SC-003)

---

## Notes

- [P] tasks touch different files with no incomplete dependency between them
- [Story] label maps each task to its spec.md Scenario for traceability
- Canonical parsing/schema/slug/theming logic lives only in Go (`internal/content`, `internal/theme`) — Principle II forbids a second implementation in the generated TS/React template
- Confirm every test in T008–T010, T017–T019, and T024 fails before writing the corresponding implementation
- Commit after each task or logical group
- Stop at each phase checkpoint to validate that story independently before moving on
