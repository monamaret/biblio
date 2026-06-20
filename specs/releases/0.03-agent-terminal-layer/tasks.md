---

description: "Task list for v0.03 — Make the site usable by AI agents and terminal clients"
---

# Tasks: v0.03 — Make the site usable by AI agents and terminal clients

**Input**: Design documents from `/specs/releases/0.03-agent-terminal-layer/`
**Prerequisites**: [plan.md](./plan.md), [spec.md](./spec.md), [research.md](./research.md), [data-model.md](./data-model.md), [contracts/mcp-tools.md](./contracts/mcp-tools.md), [quickstart.md](./quickstart.md)

**Tests**: Included and REQUIRED — Constitution Principle VI (Test-First Development, NON-NEGOTIABLE) and plan.md's Constitution Check GATE explicitly require failing MCP tool contract tests, negotiation-logic tests, and Worker route-mounting tests authored before implementation, plus a release-closing connect→list-tools→call-tool integration check against the live pilot site (plan.md Testing; spec Scenario 1, SC-001).

**Organization**: Tasks are grouped by user story, mapped from spec.md Scenarios 1–4 in their natural dependency order — the MCP tool surface (US1) before the rich views layered on top of it (US2); content negotiation (US3) is independent of US1/US2 but ships after them per the spec's own scenario order; the off-switch (US4) is verified last because it needs the US1 surface to exist to prove its absence.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: Can run in parallel (different files, no dependencies)
- **[Story]**: US1 = Scenario 1 (MCP client queries the live site), US2 = Scenario 2 (rich inline view in a graphical host), US3 = Scenario 3 (terminal client gets clean Markdown/ANSI), US4 = Scenario 4 (agent-facing layer turned off)
- File paths are exact, taken from plan.md's Project Structure (Source Code) section, with TS test paths resolved to this repo's actual convention (colocated `*.test.ts(x)` under `templates/site/profile/app/`, not a root-level `test/unit/` — root `test/` is Go-only here, per the `/speckit-analyze` finding on v0.02's tasks.md)

## Path Conventions

Extends the v0.01/v0.02 combined site/Worker template and the existing Go CLI module (plan.md "Structure Decision: additive-only"):
- `templates/site/profile/app/mcp/`, `templates/site/profile/app/negotiation/`, `templates/site/profile/app/discovery/`, `templates/site/profile/workers/` (TS — this release's work stays scoped to the single pilot site type, per spec Assumptions)
- `internal/scaffold/` (Go — `wrangler.jsonc` generation, extended to conditionally mount the Worker `main` entry)
- `test/integration/pilot-site-mcp/` (Go, repo root — matches the `test/integration/pilot-site/` convention from v0.02)

---

## Phase 1: Setup (Shared Infrastructure)

**Purpose**: Add the dependencies and Worker-entry scaffolding both the MCP surface and the negotiation layer need, before any user story work begins

- [ ] T001 Add `agents` (Cloudflare Agents SDK, `McpAgent` class), `@modelcontextprotocol/sdk`, `zod`, and `md4x` to `templates/site/profile/package.json` (research spike 002: `agents` wraps `@modelcontextprotocol/sdk`; both are required direct/transitive deps for the `McpAgent` import path used in T013)
- [ ] T002 [P] Author the Worker entry script skeleton (no routing logic yet — just a `fetch` handler stub) in `templates/site/profile/workers/entry.ts`
- [ ] T003 [P] Write failing test for `internal/scaffold/generator.go`'s `renderWranglerJSONC` conditionally adding `"main": "workers/entry.ts"`, an `assets.binding` of `"ASSETS"`, a `durable_objects` binding + `migrations` entry for the `McpAgent` class, and `"compatibility_flags": ["nodejs_compat"]` when `mcp_enabled` is `true`, and omitting all of the above entirely when `false` (FR-001, FR-002; data-model.md "MCP Server Mount" `mounted`; research spike 002 — `nodejs_compat` and the DO binding/migration are both confirmed-required for `McpAgent`, not optional) in `internal/scaffold/generator_test.go`
- [ ] T004 Implement the `renderWranglerJSONC` change to make T003 pass, in `internal/scaffold/generator.go` (depends on T003)

---

## Phase 2: Foundational (Blocking Prerequisites)

**Purpose**: Shared content-index query access and Worker routing skeleton that every user story depends on

**⚠️ CRITICAL**: No user story work can begin until this phase is complete

- [ ] T005 Implement the content-index-backed query helper shared by every MCP tool and the negotiation layer — reads `content-index.json` read-only (Principle II; data-model.md "MCP Tool" `query-source`) in `templates/site/profile/app/mcp/query.ts`
- [ ] T006 [P] Implement the Worker entry script's routing skeleton: delegate `/mcp*` to the (not-yet-built) MCP server, delegate page requests to the (not-yet-built) negotiation layer, and fall through to `env.ASSETS.fetch(request)` for everything else, in `templates/site/profile/workers/entry.ts` (depends on T002, T005)

**Checkpoint**: Foundation ready — the Worker can route to MCP/negotiation/static-asset branches and query `content-index.json`; user story implementation can begin

---

## Phase 3: User Story 1 — An MCP client queries the live site (Priority: P1) 🎯 MVP

**Goal**: A deployed site with `mcp_enabled: true` exposes `/mcp*` with shared and type-specific tools an MCP client can list and call, each returning a correctly structured, self-sufficient Markdown/plain-text result (FR-001, FR-003 partial; spec Scenario 1)

**Independent Test**: Connect an MCP client to `/mcp`, run `list-tools`, then `call-tool` for a tool appropriate to the profile site type, and confirm a complete connect → list-tools → call-tool exchange succeeds on the first attempt with a clean Markdown/plain-text `content` result — per [quickstart.md](./quickstart.md) Scenario 1 and [contracts/mcp-tools.md](./contracts/mcp-tools.md) Tool contract

### Tests for User Story 1 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T007 [P] [US1] Write failing test for `/mcp*` route mounting when `mcp_enabled` is `true` (list-tools succeeds) in `templates/site/profile/app/mcp/server.test.ts`
- [ ] T008 [P] [US1] Write failing test for the shared (cross-site-type) MCP tools' Markdown/plain-text `content` contract in `templates/site/profile/app/mcp/tools/shared.test.ts`
- [ ] T009 [P] [US1] Write failing test for the profile-type-specific MCP tools' `content` contract in `templates/site/profile/app/mcp/tools/profile.test.ts`
- [ ] T010 [P] [US1] Write failing test for a tool call requesting a capability the profile site type doesn't define returning a clear, structured error — not a silent empty result or crash (spec Edge Cases) in `templates/site/profile/app/mcp/server.test.ts`

### Implementation for User Story 1

- [ ] T011 [P] [US1] Implement the shared MCP tools (cross-site-type queries against `content-index.json`) in `templates/site/profile/app/mcp/tools/shared.ts` (depends on T005)
- [ ] T012 [P] [US1] Implement the profile-type-specific MCP tools in `templates/site/profile/app/mcp/tools/profile.ts` (depends on T005)
- [ ] T013 [US1] Implement the in-Worker MCP server as an `agents/mcp` `McpAgent` subclass (tool registry via `McpServer.registerTool`, `/mcp*` route via `MyMCP.serve("/mcp", { binding: ... })`, structured-error responses for unsupported capabilities) in `templates/site/profile/app/mcp/server.ts` (depends on T006, T011, T012; contracts/mcp-tools.md "Tool contract", "Failure-mode contract"; research spike 002 has a working reference implementation)
- [ ] T014 [US1] Wire the MCP server into the Worker entry script's `/mcp*` delegation in `templates/site/profile/workers/entry.ts` (depends on T006, T013)

**Checkpoint**: User Story 1 is independently functional — an MCP client can connect, list tools, and call them against the deployed pilot site (SC-001)

---

## Phase 4: User Story 2 — A graphical host renders a rich inline view (Priority: P2)

**Goal**: Tool results that include a richer inline view render correctly in graphical MCP hosts, inherit the site's `tokens.css` theme, and are always accompanied by the self-sufficient text/Markdown response from US1 — never a substitute for it (FR-003; spec Scenario 2)

**Independent Test**: Open a tool result that includes a `ui-resource` in a graphical MCP host and confirm it renders themed per `tokens.css` and is accompanied by the same clean `content` response; simulate a sandbox that can't load theme tokens and confirm a usable, unthemed fallback rather than a broken render — per [quickstart.md](./quickstart.md) Scenario 2 and spec Edge Cases

### Tests for User Story 2 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T015 [P] [US2] Write failing test for an MCP App UI bundle rendering correctly and importing `tokens.css` in `templates/site/profile/app/mcp/ui/chart-view.test.tsx`
- [ ] T016 [P] [US2] Write failing test that a tool result carrying a `ui-resource` still produces the unchanged, self-sufficient `content` response alongside it (Principle IV; never a substitute) in `templates/site/profile/app/mcp/tools/shared.test.ts`
- [ ] T017 [P] [US2] Write failing test for the MCP App UI bundle falling back to a usable, unthemed view when its sandbox can't load `tokens.css` (spec Edge Cases) in `templates/site/profile/app/mcp/ui/chart-view.test.tsx`

### Implementation for User Story 2

- [ ] T018 [P] [US2] Implement the small, sandboxed, postMessage-bridged React + Vite MCP App UI bundle for chart-like tool results, importing the shared `tokens.css`, with a try/catch theme-load fallback in `templates/site/profile/app/mcp/ui/chart-view.tsx` (research.md "MCP App UI bundle approach")
- [ ] T019 [US2] Extend the shared/profile tool implementations to attach a `ui-resource` reference (never replacing `content`) for results suited to a richer view (depends on T011, T012, T018; data-model.md "MCP Tool" `ui-resource`)

**Checkpoint**: User Story 2 is independently functional — richer inline views render where supported and degrade gracefully everywhere else, with the text answer always present (Scenario 2)

---

## Phase 5: User Story 3 — A terminal client gets a clean Markdown view (Priority: P3)

**Goal**: A page request signaling a Markdown/terminal preference gets a clean, readable Markdown (tier i) or ANSI (tier ii) rendering instead of full browser HTML, with deterministic negotiation and correct `Vary` caching behavior, plus discoverability artifacts for agents/crawlers (FR-004, FR-005; spec Scenario 3)

**Independent Test**: `curl -H "Accept: text/markdown" <page>` and `curl -A "curl/8.0" <page>` both return clean Markdown/ANSI with no leftover browser chrome and a `Vary: Accept, User-Agent` header; an ambiguous `Accept` header resolves deterministically; `/llms.txt` is present and lists key pages — per [quickstart.md](./quickstart.md) Scenario 3 and [contracts/mcp-tools.md](./contracts/mcp-tools.md) "Negotiated page-rendering contract", "Discovery artifacts"

### Tests for User Story 3 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T020 [P] [US3] Write failing test for `Accept`/`User-Agent`-based negotiation selecting the correct variant (html/markdown/ansi) and setting `Vary: Accept, User-Agent` on every response in `templates/site/profile/app/negotiation/accept.test.ts`
- [ ] T021 [P] [US3] Write failing test for deterministic resolution of an ambiguous/conflicting `Accept` header (spec Edge Cases) in `templates/site/profile/app/negotiation/accept.test.ts`
- [ ] T022 [P] [US3] Write failing test for tier-ii `md4x` `renderToAnsi()` output via the Wrangler `.wasm` import wiring (research.md "In-Worker Markdown→ANSI rendering") in `templates/site/profile/app/negotiation/ansi.test.ts`
- [ ] T023 [P] [US3] Write failing test for the build-time-generated `llms.txt` content (site name/description + flat key-page links derived from `content-index.json`, not a fuller dump) in `templates/site/profile/app/discovery/llms-txt.test.ts`
- [ ] T024 [P] [US3] Write failing test for a negotiated request targeting a page that doesn't exist returning a clear, structured "not found" rather than a silent empty/200 response (spec Edge Cases) in `templates/site/profile/app/negotiation/accept.test.ts`

### Implementation for User Story 3

- [ ] T025 [P] [US3] Implement tier-i raw Markdown passthrough in `templates/site/profile/app/negotiation/markdown.ts` (depends on T005)
- [ ] T026 [P] [US3] Implement tier-ii in-Worker `md4x` `renderToAnsi()` rendering, wiring the WASM module via Wrangler's native `.wasm` import and `init({ wasm: wasmModule })` (research.md) in `templates/site/profile/app/negotiation/ansi.ts` (depends on T001, T025)
- [ ] T027 [US3] Implement the `Accept`/`User-Agent`-led negotiation branch (variant selection, deterministic ambiguous-header resolution, `Vary` header, structured not-found handling) in `templates/site/profile/app/negotiation/accept.ts` (depends on T025, T026; contracts/mcp-tools.md "Negotiated page-rendering contract", "Failure-mode contract")
- [ ] T028 [US3] Wire the negotiation branch into the Worker entry script's page-request delegation in `templates/site/profile/workers/entry.ts` (depends on T006, T027)
- [ ] T029 [P] [US3] Implement the build-time `llms.txt` generator (a Vite plugin writing it into the client build output during `biblio-cli build`, derived from `content-index.json`) in `templates/site/profile/app/discovery/llms-txt.ts` (depends on T005)
- [ ] T030 [US3] Register the `llms.txt` generator plugin in `templates/site/profile/vite.config.ts` (depends on T029)
- [ ] T031 [US3] Add per-page Markdown-friendly discoverability hints (sibling addresses and/or link hints) to the prerendered route output from v0.02 (FR-005) in `templates/site/profile/app/root.tsx`

**Checkpoint**: User Story 3 is independently functional — terminal/`curl` clients get clean renderings, caching is correct, and the site is discoverable by agents/crawlers (SC-002)

---

## Phase 6: User Story 4 — The agent-facing layer can be turned off (Priority: P4)

**Goal**: When `mcp_enabled` is `false` for a site, no `/mcp*` route, tool surface, or related generated artifact is present or reachable at all — a build-time decision, not a runtime toggle (FR-002; spec Scenario 4)

**Independent Test**: Scaffold or reconfigure a project with `mcp_enabled: false`, rebuild, redeploy, then `curl -I` its `/mcp` endpoint and confirm zero reachable trace — per [quickstart.md](./quickstart.md) Scenario 4 and [contracts/mcp-tools.md](./contracts/mcp-tools.md) "Endpoint mounting"

### Tests for User Story 4 ⚠️

> Write these FIRST; confirm they FAIL before implementation (Principle VI)

- [ ] T032 [P] [US4] Write failing test that the Worker entry script never reaches the MCP delegation branch when `mcp_enabled` is `false` (verified via the absent `main` wiring from T003/T004) in `templates/site/profile/workers/entry.test.ts`
- [ ] T033 [P] [US4] Write failing Go test that a scaffolded project with `mcp_enabled: false` produces a `wrangler.jsonc` with no `main`/MCP-related bindings at all (data-model.md "MCP Server Mount" — no "disabled but present" intermediate state) in `internal/scaffold/generator_test.go`

### Implementation for User Story 4

- [ ] T034 [US4] Verify and, if needed, adjust the Worker entry script (T006/T014) and `renderWranglerJSONC` (T004) so the `mcp_enabled: false` path truly emits zero MCP-related route, artifact, or binding (depends on T004, T006, T014, T032, T033)

**Checkpoint**: All user stories independently functional — turning the agent-facing layer off leaves zero reachable trace (SC-003)

---

## Phase 7: Polish & Cross-Cutting Concerns

**Purpose**: The release-closing integration evidence and final consistency pass

- [ ] T035 Author and run the end-to-end connect → list-tools → call-tool integration check against one real, deployed pilot site, plus the terminal-`curl` negotiation check and the `mcp_enabled: false` "zero reachable trace" check, in `test/integration/pilot-site-mcp/` (FR-001..FR-004; SC-001, SC-002, SC-003; constitution Release Process release-closing test)
- [ ] T036 [P] Update `specs/releases/0.03-release.md` status and this release's `AGENTS.md`/CLAUDE.md SPECKIT markers to reflect completion
- [ ] T037 [P] Run `gofmt`/`go vet` (`scripts/check.sh`) and the TS-side lint/typecheck across all files touched in this release

---

## Dependencies & Execution Order

- **Phase 1 (Setup)** → **Phase 2 (Foundational)** → **Phase 3 (US1)** → **Phase 4 (US2)** → **Phase 5 (US3)** → **Phase 6 (US4)** → **Phase 7 (Polish)**
- US2 depends on US1's tool registry existing (a rich view has nothing to attach to otherwise)
- US3 is independent of US1/US2 at the code level (page negotiation, not MCP tools) but is sequenced after them per the spec's own scenario order
- US4 depends on US1 (and the Phase 1 `wrangler.jsonc` gating work) existing — proving "no trace" requires there to be a trace to remove in the first place
- Within each phase, tasks not marked `[P]` touch shared files/state and must run sequentially in listed order; `[P]` tasks touch independent files and may run in parallel

## Parallel Execution Examples

- Phase 1: T002 and T003 in parallel (independent files), both before T004
- Phase 2: T006 depends on both T002 and T005, so it starts only once Setup is fully done
- Phase 3 tests: T007, T008, T009, T010 in parallel before any Phase 3 implementation task
- Phase 3 implementation: T011 and T012 in parallel, both before T013
- Phase 4 tests: T015, T016, T017 in parallel before T018/T019
- Phase 5 tests: T020, T021, T022, T023, T024 in parallel before their implementation tasks
- Phase 5 implementation: T025, T026, T029 in parallel (independent files) before T027/T030
- Phase 6 tests: T032 and T033 in parallel before T034
- Phase 7: T036 and T037 in parallel

## Implementation Strategy

**MVP first**: Phase 1 → Phase 2 → Phase 3 (US1) delivers a working `/mcp*` surface an AI assistant can connect to, list tools on, and call — independently testable and demonstrable before any rich-view, negotiation, or off-switch work exists.

**Incremental delivery**: Add Phase 4 (US2) for richer graphical-host views layered on the same tools, then Phase 5 (US3) for the terminal/`curl`-facing negotiation layer (independent of US1/US2 but sequenced after per spec order), then Phase 6 (US4) to close the loop on the off-switch invariant. Phase 7's integration check is the release-closing evidence for SC-001–SC-003, demonstrable purely on agent/terminal-client evidence per spec SC-004 — never requiring v0.02's browser-facing behavior to be re-proven.
