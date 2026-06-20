---

description: "Task list template for feature implementation"
---

# Tasks: v0.06 — General Availability

**Input**: Design documents from `specs/releases/0.06-general-availability/`

**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/update-and-escalation.md](./contracts/update-and-escalation.md), [quickstart.md](./quickstart.md)

**Tests**: Included and REQUIRED, not optional — Constitution Principle VI (Test-First Development, NON-NEGOTIABLE) and the Release Process both mandate failing tests authored before implementation, plus unit *and* integration coverage before this release can close.

**Organization**: Tasks are grouped by user story, mapped from spec.md's four scenarios in priority order (Scenario 1 → US1/P1 through Scenario 4 → US4/P4). Scenario 4 (GA Readiness Review) is necessarily the closing increment — it aggregates the evidence the other three stories produce (FR-006, SC-002, SC-003).

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: Which user story this task belongs to (US1–US4)
- Paths are relative to the repository root and follow plan.md's Project Structure section

## Path Conventions

Single Go module (`biblio-cli`) at the repository root, per plan.md:
- `internal/migrate/`, `internal/escalation/`, `internal/cli/` — implementation
- `test/unit/`, `test/integration/` — tests
- `docs/storage-escalation-guide.md` — documentation deliverable

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Scaffold the new packages, directories, and doc stub this release adds — no behavior yet.

- [ ] T001 Create `internal/cli/update.go`, `internal/cli/upgrade.go`, `internal/migrate/classify.go`, `internal/migrate/merge.go`, `internal/migrate/conflict.go`, and `internal/escalation/guidance.go` as empty package files with package declarations, per plan.md's Project Structure
- [ ] T002 [P] Create test directory skeletons: `test/unit/migrate/`, `test/unit/escalation/`, `test/integration/update-three-sites/`, `test/integration/escalation-guidance/`, `test/integration/ga-readiness-review/`
- [ ] T003 [P] Create `docs/storage-escalation-guide.md` stub with section headings for Thresholds, Real-Scale Check, D1 Escalation Path, and Vectorize Escalation Path (FR-004/FR-005 deliverable, filled in by US3)

**Checkpoint**: Packages and test layout exist; nothing compiles behavior yet.

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: The classify → merge → conflict machinery in `internal/migrate` that both US1 (single-site update) and US2 (cross-type update) depend on as their shared contract (contracts/update-and-escalation.md Step 1–4). US3 and US4 do not depend on this phase.

**⚠️ CRITICAL**: US1 and US2 cannot begin until this phase is complete.

- [ ] T004 [P] Write failing unit test for Scaffold Region Classification (`scaffold-tracked` / `author-customized` / `ambiguous`, every region classified before any write) in `test/unit/migrate/classify_test.go`
- [ ] T005 Implement Scaffold Region Classification in `internal/migrate/classify.go`, derived from the single canonical scaffold-template definitions (Principle II) — makes T004 pass
- [ ] T006 [P] Write failing unit test for diff/merge logic that advances `scaffold-tracked` regions and preserves `author-customized` regions byte-for-byte in `test/unit/migrate/merge_test.go`
- [ ] T007 Implement merge logic in `internal/migrate/merge.go` (depends on T005) — makes T006 pass
- [ ] T008 [P] Write failing unit test for conflict surfacing — every `ambiguous` classification becomes an Update Conflict that is never auto-resolved, and an interrupted run leaves the site in a consistent, recoverable state — in `test/unit/migrate/conflict_test.go`
- [ ] T009 Implement conflict surfacing and interruption-recovery state in `internal/migrate/conflict.go` (depends on T005, T007) — makes T008 pass

**Checkpoint**: `internal/migrate` classify/merge/conflict logic is implemented and unit-tested. US1 and US2 implementation can begin.

---

## Phase 3: User Story 1 - Update a site to a newer scaffold version without losing content (Priority: P1) 🎯 MVP

**Goal**: `biblio-cli update`, run against one live site, advances scaffold-tracked parts, preserves authored content/customizations untouched, and surfaces any conflicts for the author to resolve (FR-001, FR-002, Scenario 1).

**Independent Test**: Run `biblio-cli update` in `my-profile-site` (quickstart.md Scenario 1); confirm authored content is byte-for-byte intact, scaffold-tracked parts advanced, conflicts (if any) clearly surfaced and not auto-resolved, and that a mid-run `Ctrl+C` leaves the site in a consistent, recoverable state.

### Tests for User Story 1

- [ ] T010 [P] [US1] Write failing integration test: an Update Run against one real site preserves `author-customized` regions, advances `scaffold-tracked` regions, and reports the preserved/advanced/conflicts triple (data-model.md Update Run) in `test/integration/update-three-sites/single_site_test.go`
- [ ] T011 [P] [US1] Write failing integration test: interrupting an update mid-run (simulated cancel) leaves the site in a consistent, recoverable state with no half-migrated files (spec Edge Cases) in `test/integration/update-three-sites/interruption_test.go`

### Implementation for User Story 1

- [ ] T012 [US1] Implement Update Run orchestration (classify → advance → preserve → surface → commit, contracts/update-and-escalation.md Steps 1–5) in `internal/migrate/update_run.go` (depends on T005, T007, T009) — makes T010, T011 pass
- [ ] T013 [US1] Wire the `biblio-cli update` Cobra command in `internal/cli/update.go`, operating on the current project directory with exit code `0` on clean completion and non-zero with state-recovery guidance on interruption/unresolved conflicts (depends on T012)
- [ ] T014 [US1] Implement the Update Run success/failure report (preserved/advanced/conflicts output) per contracts/update-and-escalation.md "Output (success)" / "Output (failure / interruption)" in `internal/migrate/update_run.go`
- [ ] T015 [US1] Run quickstart.md Scenario 1 against the real `my-profile-site`: confirm content intact after `biblio-cli update`, and that a `Ctrl+C` mid-run recovers cleanly

**Checkpoint**: `biblio-cli update` works end-to-end against one real site. This is the MVP — independently demoable.

---

## Phase 4: User Story 2 - The update flow works across all three site types (Priority: P2)

**Goal**: The exact same update/upgrade contract holds, unmodified, across the profile, project-log, and documentation sites (FR-003, Scenario 2), and `biblio-cli upgrade` (version-targeted variant) is available alongside `update`.

**Independent Test**: Run `biblio-cli update` against `my-profile-site`, `my-project-log`, and `my-docs-site` in turn (quickstart.md Scenario 2); confirm all three complete successfully with content intact, in one sitting (SC-001).

### Tests for User Story 2

- [ ] T016 [P] [US2] Write failing integration test: `biblio-cli update` succeeds against the real project-log site with content intact in `test/integration/update-three-sites/project_log_test.go`
- [ ] T017 [P] [US2] Write failing integration test: `biblio-cli update` succeeds against the real docs site with content intact in `test/integration/update-three-sites/docs_site_test.go`
- [ ] T018 [US2] Write failing integration test: `biblio-cli upgrade <version>` targets a specific scaffold version and succeeds across all three real site types in `test/integration/update-three-sites/upgrade_version_test.go`

### Implementation for User Story 2

- [ ] T019 [US2] Implement the `biblio-cli upgrade` Cobra command in `internal/cli/upgrade.go`, accepting a target-version identifier and sharing the `internal/migrate` classify/merge/conflict machinery with `update` (depends on T012) — makes T018 pass
- [ ] T020 [US2] Run quickstart.md Scenario 2's loop across all three real sites (profile/project-log/docs) in one sitting and confirm SC-001 (depends on T015, T016, T017, T019)

**Checkpoint**: Update/upgrade is proven to generalize across all three real site types, not just the first one it was validated against.

---

## Phase 5: User Story 3 - Storage-escalation guidance is checked against real sites (Priority: P3)

**Goal**: Documented D1/Vectorize storage-escalation guidance correctly reports each real site's threshold status against its actual scale, and each recommended path's steps are proven concretely executable, not just described (FR-004, FR-005, Scenario 3).

**Independent Test**: Check each of the three real sites against `docs/storage-escalation-guide.md`'s documented thresholds (quickstart.md Scenario 3); confirm threshold-status correctly reflects real `content-index.json`-derived scale, and that a dry-run of each recommended path's steps actually completes.

### Tests for User Story 3

- [ ] T021 [P] [US3] Write failing unit test for threshold-status computation (`not-crossed` / `crossed` for D1 volume/relational and Vectorize search-quality triggers) against `content-index.json`-derived scale in `test/unit/escalation/guidance_test.go`
- [ ] T022 [P] [US3] Write failing integration test: the Storage-Escalation Guidance Check, run against all three real sites' actual `content-index.json` scale, reports a status consistent with their real data in `test/integration/escalation-guidance/real_sites_test.go`

### Implementation for User Story 3

- [ ] T023 [US3] Implement the Storage-Escalation Guidance Check (`threshold-status`, `path-executable-on-demand`) in `internal/escalation/guidance.go` — makes T021, T022 pass
- [ ] T024 [US3] Document concrete, checkable D1 (storage volume / relational-query) and Vectorize (search-quality) thresholds and their recommended escalation steps in `docs/storage-escalation-guide.md` (depends on T003)
- [ ] T025 [US3] Dry-run/walk through the D1 escalation path's steps against a copy of real site data and record the proof of executability in `docs/storage-escalation-guide.md` (depends on T024)
- [ ] T026 [US3] Dry-run/walk through the Vectorize escalation path's steps against a copy of real site data and record the proof of executability in `docs/storage-escalation-guide.md` (depends on T024)
- [ ] T027 [US3] Run quickstart.md Scenario 3 against all three real sites and confirm the guidance correctly reflects each site's real, current scale (depends on T023)

**Checkpoint**: Storage-escalation guidance is documented and proven executable on demand, without requiring either path to actually be carried out (Principle V — never pre-built "just in case").

---

## Phase 6: User Story 4 - GA readiness review (Priority: P4)

**Goal**: Cross-reference every roadmap deliverable §1–§10 and every release v0.01–v0.06 against this release's acceptance bar, producing a `ga-determination` record that marks the platform Generally Available (FR-006, Scenario 4, SC-002, SC-003).

**Independent Test**: Run the GA Readiness Review (quickstart.md Scenario 4); confirm every ordinary deliverable is `complete-and-demonstrated` and the two conditional storage-escalation items are `proven-executable-on-demand`, referencing — never re-running — prior releases' closed evidence.

**Depends on**: US1, US2 (update-evidence) and US3 (escalation-evidence) — this story aggregates their results and cannot start before they complete.

### Tests for User Story 4

- [ ] T028 [US4] Write failing integration test: the GA Readiness Review aggregates the Update Run suite (US1/US2) and Storage-Escalation Guidance Checks (US3) by reference, without re-running v0.01–v0.05's already-closed evidence (SC-003) in `test/integration/ga-readiness-review/aggregate_test.go`

### Implementation for User Story 4

- [ ] T029 [US4] Implement the GA Readiness Review artifact (`deliverables[]`, `update-evidence`, `escalation-evidence`, `ga-determination`) per data-model.md in `internal/escalation/ga_review.go` (depends on T020, T027) — makes T028 pass
- [ ] T030 [US4] Cross-reference every roadmap deliverable in `specs/research/roadmap.md` §1–§10, and every release v0.01–v0.06, against this release's acceptance bar, recording each deliverable's status (depends on T029)
- [ ] T031 [US4] Produce the closing `ga-determination: true` record once every deliverable meets its bar (depends on T030)

**Checkpoint**: The platform is marked Generally Available — every roadmap deliverable is complete-and-demonstrated or proven-executable-on-demand.

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: Release-closing housekeeping required by the constitution's Release Process and Decision Recording workflow.

- [ ] T032 [P] Record the real `go build`/`go test`/lint commands now established by this release in `AGENTS.md` (Build / Test / Lint section), keeping `CLAUDE.md` in sync per the constitution's "kept current" requirement
- [ ] T033 [P] Assemble `specs/release_0.06/` packaging directory containing the final feature list (pointing to features 009/010), release notes, and a create/configure/update/deploy/destroy usage guide, per the constitution's Release Process
- [ ] T034 Run the complete quickstart.md validation (Scenarios 1–4) end-to-end and confirm SC-001, SC-002, and SC-003 all hold (depends on T015, T020, T027, T031)
- [ ] T035 [P] Update the progress cursor in `specs/research/roadmap.md` marking §9 (replicate to remaining site types) and §10 (storage escalation) complete/proven-ready
- [ ] T036 Verify the usage guide in `specs/release_0.06/` covers the full `biblio-cli` lifecycle comprehensively — create, configure, update, deploy, and destroy (constitution Release Process item 4) — **and** includes a build-from-source section (cloning the repo, toolchain prerequisites, `go build`/`go install` steps, and running the locally built binary in place of a release download), so GA does not depend on a published binary being available (depends on T032, T033)

---

## Dependencies & Execution Order

### Phase Dependencies

- **Setup (Phase 1)**: No dependencies — can start immediately
- **Foundational (Phase 2)**: Depends on Setup — BLOCKS User Story 1 and User Story 2
- **User Story 1 (Phase 3)**: Depends on Foundational completion
- **User Story 2 (Phase 4)**: Depends on Foundational completion; T019/T020 additionally depend on User Story 1's `update_run.go` (T012) and command wiring (T013)
- **User Story 3 (Phase 5)**: Depends on Setup only (T003) — independent of Foundational/US1/US2, can run in parallel with Phase 2–4
- **User Story 4 (Phase 6)**: Depends on User Story 1, User Story 2, AND User Story 3 completing — it aggregates their evidence and cannot start before
- **Polish (Phase 7)**: Depends on all four user stories being complete

### User Story Dependencies

- **User Story 1 (P1)**: No dependencies on other stories — true MVP
- **User Story 2 (P2)**: Builds on US1's update machinery (shared `internal/migrate` + `update_run.go`) but is independently testable once US1 lands
- **User Story 3 (P3)**: No dependencies on US1/US2 — can be developed in parallel with them
- **User Story 4 (P4)**: Depends on US1 + US2 + US3 — the only story that is not independent, by design (it is the closing aggregation)

### Parallel Opportunities

- All Setup tasks marked [P] (T002, T003) can run in parallel
- All Foundational test-authoring tasks marked [P] (T004, T006, T008) can run in parallel; their implementation tasks (T005, T007, T009) are sequential due to shared-file dependencies
- Once Foundational completes, User Story 1 and User Story 3 can proceed in parallel (different packages, no shared files)
- User Story 2 can start once User Story 1's T012/T013 land, in parallel with continued User Story 3 work
- User Story 4 must wait for US1, US2, and US3 to all complete

---

## Parallel Example: Foundational Phase

```bash
# Launch all Foundational unit tests together (different files, no dependencies):
Task: "Write failing unit test for Scaffold Region Classification in test/unit/migrate/classify_test.go"
Task: "Write failing unit test for diff/merge logic in test/unit/migrate/merge_test.go"
Task: "Write failing unit test for conflict surfacing in test/unit/migrate/conflict_test.go"
```

## Parallel Example: User Story 1 and User Story 3

```bash
# Once Phase 2 (Foundational) completes, these two stories can proceed in parallel:
Task: "User Story 1 — implement Update Run orchestration in internal/migrate/update_run.go"
Task: "User Story 3 — implement Storage-Escalation Guidance Check in internal/escalation/guidance.go"
```

---

## Implementation Strategy

### MVP First (User Story 1 Only)

1. Complete Phase 1: Setup
2. Complete Phase 2: Foundational (CRITICAL — blocks US1 and US2)
3. Complete Phase 3: User Story 1
4. **STOP and VALIDATE**: Run quickstart.md Scenario 1 against the real profile site independently
5. This is the MVP: `biblio-cli update` works, content is preserved, conflicts are surfaced

### Incremental Delivery

1. Setup + Foundational → foundation ready
2. Add User Story 1 → validate against `my-profile-site` → MVP
3. Add User Story 2 → validate `update`/`upgrade` across all three real sites → SC-001 demonstrated
4. Add User Story 3 (can run in parallel with 2–3) → validate guidance against real site scale
5. Add User Story 4 (only once 1–3 are done) → GA Readiness Review → SC-002/SC-003 demonstrated → **General Availability**

### Single-Author Sequencing Note

Per Principle I/V and research.md, this release's central design work was deliberately deferred until v0.04/v0.05 produced the second and third real sites — there is no "team" to parallelize across in this single-author project; the parallel opportunities above describe what is *independently sequenceable*, not concurrent staffing.

---

## Notes

- [P] tasks = different files, no dependencies
- [Story] label maps task to specific user story (US1–US4) for traceability
- Tests are authored first and confirmed to fail for the right reason before implementation, per Constitution Principle VI (NON-NEGOTIABLE) — this is not optional for this release
- Every Update Run task touches real, live sites — no synthetic fixtures (spec Assumptions, plan.md Scale/Scope)
- Commit after each task or logical group
- Stop at any checkpoint to validate a story independently
- US4 (GA Readiness Review) is intentionally the one story that is *not* independent — it is the closing aggregation FR-006 requires
