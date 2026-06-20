<!--
Sync Impact Report
Version change: 1.3.0 → 1.3.1
Modified principles: n/a (no principle redefined)
Corrected sections:
  - Release Process — fixed the packaging directory pattern to match the
    convention every release from v0.01 onward actually uses
    (`specs/releases/X.XX-name/`, not the originally-written
    `specs/release_X.XX/`); found via /speckit-analyze on v0.06's artifacts.
Templates requiring updates: none — this is a wording correction, not a
  rule change.
Follow-up TODOs: none.

Sync Impact Report (prior)
Version change: 1.2.0 → 1.3.0
Modified principles: n/a (no existing principle redefined)
Added sections:
  - Release Process — defines the specs/release_X.XX/ packaging convention:
    final feature list, release notes, unit + integration tests, and a
    create/configure/update/deploy/destroy usage guide per release
Removed sections: none
Templates requiring updates:
  - .specify/templates/plan-template.md ⚠ pending (Constitution Check gate is generic;
    consider adding an explicit "tests written and failing first?" sub-gate)
  - .specify/templates/spec-template.md ⚠ pending (no test-first/release references to add yet)
  - .specify/templates/tasks-template.md ⚠ pending (ensure generated task ordering puts
    test-authoring tasks before implementation tasks, and surfaces release-packaging
    tasks for release-closing work, once real work is planned)
  - .specify/templates/commands/*.md ⚠ pending (verify no outdated agent-specific references)
  - AGENTS.md ✅ updated (new "Release process" section mirrors specs/release_X.XX/ contents)
  - CLAUDE.md ✅ no change needed (defers to AGENTS.md for conventions)
Follow-up TODOs: none — concrete release-numbering start point (e.g. 1.00 vs 0.01)
  remains a call for the project owner to make when the first release is actually cut.

Sync Impact Report (prior)
Version change: 1.1.0 → 1.2.0
Modified principles: n/a (no existing principle redefined)
Expanded sections:
  - Decision Recording & Research Workflow — replaced the prior loose
    description with the explicit, ordered feature pipeline: numbered research
    discussion doc -> feature item creation -> roadmap inclusion -> architecture
    sync -> Spec Kit workflow (linking spec-driven.md as the documented process
    to follow for the final step)

Sync Impact Report (prior)
Version change: 1.0.0 → 1.1.0
Modified principles: n/a (no existing principle redefined)
Added sections:
  - Core Principle VI. Test-First Development (NON-NEGOTIABLE) — TDD/red-green-
    refactor mandate across biblio-cli (Go) and the site/MCP server (TS/React)

Sync Impact Report (prior)
Version change: [TEMPLATE] → 1.0.0
Modified principles: n/a (initial ratification — template placeholders replaced)
Added sections:
  - Core Principles I–V (Single-Author Minimal-Dependency Bias; One Pipeline One
    Source of Truth; CLI-First Guided Experience; Text-First Progressive
    Enhancement for MCP Tools; Simplest Topology That Can Grow)
  - Security & Operational Posture
  - Decision Recording & Research Workflow
  - Governance
-->

# Biblio Constitution

## Core Principles

### I. Single-Author, Minimal-Dependency Bias (No Speculative Abstraction)
Every dependency, abstraction, configuration knob, or generalization MUST justify
itself against `biblio`'s actual scope: a private, single-author platform at
personal/portfolio scale — never against hypothetical multi-user, team, or
future needs. Escalations (storage beyond a bundled JSON index, search beyond
keyword matching, charting libraries, nested config shapes, per-site auth,
alternate site generators, etc.) are deferred until a documented threshold is
actually crossed, and each one gets its own research pass and Decision Log
entry when it happens. Features, refactors, or designs that point toward
productizing, multi-tenancy, or supporting external users or other site
generators MUST be rejected outright — they are explicitly out of scope.
**Rationale**: This bias is the throughline of every confirmed decision in
`specs/research/001-static-site-platform-with-mcp.md` (no MDX, no charting
library, flat `mcp_enabled` boolean, JSON-before-D1, keyword-before-Vectorize).
Building for a hypothetical future user or scale is waste for a project that
will only ever have one.

### II. One Pipeline, One Source of Truth
Content authored as Markdown + YAML frontmatter under `content/` MUST drive
*both* the static site build and the MCP server's data, through a single
content-index generation step that `biblio-cli` owns on the Go side
(`content-index.json`, step 1 of the build pipeline, ahead of `react-router
build` and the Wrangler bundle/deploy steps). Per-site design tokens
(`tokens.css`) MUST be generated once and imported by both the site build and
the MCP App UI bundle. Parsing, schema validation, theming, and slug-derivation
logic MUST NOT be duplicated between the Go CLI and the JS/TS site — pick one
owner and have the other consume its output.
**Rationale**: Confirmed explicitly in the Decision Log — the canonical
parsing/schema logic lives in `biblio-cli` (reusable for scaffolding and
validation too) precisely so the project never ends up bridging two divergent
content or theming pipelines across languages.

### III. CLI-First, Guided Experience
`biblio-cli` (Go + Cobra + Viper, with a Bubble Tea/Lip Gloss/Glamour TUI) MUST
remain the single, consistent entry point for the entire site lifecycle:
create → dev → build → deploy → update. Every new capability is exposed as a
`biblio-cli` command and/or flag, with the interactive TUI as the guided
fallback when flags are omitted — never as a raw `package.json` script, a
hand-assembled CI step, or a separate standalone tool. Distribution stays
single-binary (GitHub Releases via `goreleaser`, plus `go install` for those
with a Go toolchain) — no package-manager taps or multi-user packaging
infrastructure.
**Rationale**: Locked in across four workflow/UX Decision Log entries — a
consistent, guided, one-tool-does-everything surface is the natural fit for a
single author and matches the investment already committed to the TUI
experience.

### IV. Text-First, Progressive Enhancement for MCP and Terminal Surfaces
Every MCP tool MUST produce a clean, self-sufficient Markdown/plain-text
rendering (e.g. a table or ASCII sparkline for anything chart-like) as its
primary `content` response, because terminal-based MCP clients can only see
that — any richer SVG/HTML MCP App UI resource is layered on top strictly as a
bonus, never a requirement for the tool to be useful. The same rule governs the
site's terminal/CLI view: raw Markdown passthrough and in-Worker `md4x` →
ANSI rendering (tiers i/ii) are the supported baseline; a full styled
terminal "mini-site" (tier iii) is explicitly out of initial scope and is
revisited only if real usage data shows the need.
**Rationale**: `biblio`'s MCP server is expected to be used by terminal-based
coding assistants that cannot render HTML/React UI resources — textual output
must stand alone on its own merits, exactly as the Decision Log's
no-charting-library and dropped-tier-iii decisions establish.

### V. Simplest Topology That Can Grow
Default to the simplest viable architecture and escalate only on a documented,
real trigger — never speculatively:
- One Cloudflare Worker serves both the prerendered static site (Workers Static
  Assets) and the `/mcp` endpoint, mounted on a route (`example.com/mcp`), not
  a subdomain — never split across separate hosts/services.
- Default to a generated, in-memory JSON content index; escalate to D1 only
  when volume/relational-query needs genuinely outgrow it.
- Default to keyword search; escalate to Vectorize-backed semantic search only
  when search *quality* (not volume) demonstrably fails on conceptual queries.
Each escalation step is a new, separately-decided change with its own Decision
Log entry — not something to pre-build "just in case."
**Rationale**: This is the explicitly named "simplest topology that can grow"
philosophy from the architecture and storage-escalation decisions — it
collapses what would otherwise be an open-ended hosting/storage design space
into one replicable, low-maintenance pattern for a single owner.

### VI. Test-First Development (NON-NEGOTIABLE)
Tests MUST be written before implementation, MUST be reviewed and confirmed to
fail for the right reason, and only then is implementation written to make them
pass — red-green-refactor, strictly. This applies across the stack: Go unit/CLI
tests for `biblio-cli` (content-index generation, frontmatter schema parsing
and validation, slug derivation, scaffolding/templating logic), and
TypeScript/React tests plus MCP tool contract tests for the site and in-Worker
`/mcp` server (each tool's Markdown/plain-text `content` output, route
mounting, content-index querying). No feature, fix, or refactor lands without
a failing test demonstrating the gap it closes, followed by the change that
closes it. Once real build/test/lint commands exist they MUST be recorded in
`AGENTS.md` (see "Build / Test / Lint") and kept current.
**Rationale**: `biblio-cli` is the single source of truth for content parsing,
schemas, and scaffolding (Principle II) and ships as the one tool a single
author trusts end-to-end (Principle III) — a regression here silently corrupts
every site it touches. Writing the test first is what keeps that trust
warranted as the codebase grows from today's "no source code yet" baseline.

## Security & Operational Posture

The `/mcp` endpoint is public with intentionally **no auth/access control** —
this is a deliberate decision for `biblio`'s low-sensitivity, read-oriented
content and tool surface, not an oversight to "fix" later. Rely on
Cloudflare's standard zone-level protections (rate limiting, WAF, bot
management) rather than adding per-site authentication or access-control
layers. `biblio-cli` MUST provide the build-time `mcp_enabled` boolean
(default `true`) in `content/site.yaml` as the *only* MCP-related
configuration surface, deciding in one place whether to generate the content
index, mount the `/mcp*` route, and include MCP-related bindings
(Durable Objects, etc.) in `wrangler.jsonc`. Do not introduce per-site auth,
nested MCP config objects, or other access-control surface without a fresh
research pass and Decision Log entry justifying the change in threat model.

## Decision Recording & Research Workflow

The standard path from idea to shipped feature is a fixed pipeline — steps are
not skipped, reordered, or collapsed, even for small features:

1. **Numbered research discussion document** — open (or extend) a numbered doc
   under `specs/research/` (e.g. `001-static-site-platform-with-mcp.md`,
   `000-tech-stack.md`), copied from `research-template.md` when starting a new
   one, with every section (Problem Statement, Background, Libraries & Stack,
   Architecture Patterns, Existing Codebase Evaluation, Security & Compliance,
   Performance/Reliability/Ops, Risks, Open Questions, Decision Log,
   References) filled in. Confirmed decisions — including supersessions of
   earlier ones — MUST be appended to that document's Decision Log with date,
   decision, rationale, and owner, exactly as
   `001-static-site-platform-with-mcp.md` already does; superseded entries are
   kept (marked **Superseded**) rather than deleted, preserving the history of
   *why* the project changed direction.
2. **Feature item creation** — once the research settles a buildable shape,
   create a feature item under `specs/features/`, copied from
   `specs/features/TEMPLATE.md`, that points back to the confirmed decisions
   in the research doc rather than restating or re-deciding them (per the
   template's "Confirmed decisions this builds on" section — every line there
   must already be a link to something decided, not a new assumption).
3. **Roadmap inclusion** — add the new feature item to
   `specs/research/roadmap.md` in its dependency-ordered place, so the roadmap
   stays the single, current view of "what's next" and "where we are" (its
   progress cursor is updated as work begins/completes).
4. **Architecture sync** — update `specs/research/architecture.md` with any
   resulting changes to the reference diagrams/call-flows *before*
   implementation starts, so it never drifts out of sync with confirmed
   decisions and roadmap reality.
5. **Spec Kit workflow** — only then start the Spec Kit (speckit) workflow
   (`/speckit.specify` → `/speckit.clarify` → `/speckit.plan` → `/speckit.tasks`
   → `/speckit.implement`, etc.) to turn the now-settled feature item into
   spec/plan/tasks artifacts and drive implementation. See
   https://github.com/github/spec-kit/blob/main/spec-driven.md for the
   spec-driven-development process this workflow implements — follow it as
   documented rather than improvising a different shape for these artifacts.

`AGENTS.md` is the canonical, opencode-loaded project-instructions file and
the living source of truth for conventions and build/test/lint commands;
`CLAUDE.md` MUST be kept in sync with it rather than allowed to drift. Both
files — and this pipeline description itself — are updated alongside the code
that makes their current "not yet established" sections concrete.

## Release Process

Each release MUST be packaged as a directory `specs/releases/X.XX-name/`
(e.g. `specs/releases/0.01-scaffold-and-theme/`,
`specs/releases/0.06-general-availability/`), containing exactly:

1. **Final feature list** — the definitive list of features shipped in that
   release, each entry pointing back to its `specs/features/` item (and, through
   it, to the research doc and Decision Log entries that justified it) — this
   is the closing record of "what actually shipped," not a re-statement of
   open-ended roadmap intent.
2. **Release notes** — a human-readable summary of what changed and why,
   written for the project owner as the sole consumer (no marketing framing,
   no external-user assumptions — consistent with Principle I).
3. **Unit and integration tests** — the release MUST include both unit tests
   (per-package/per-component, in `biblio-cli` and the site/MCP server) *and*
   integration tests that exercise the release's features end to end, in
   keeping with Principle VI (Test-First Development); a release is not
   "final" until both layers exist and pass.
4. **Usage guide** — a guide covering, at minimum, how to **create**,
   **configure**, **update**, **deploy**, and **destroy** a `biblio` site using
   that release's `biblio-cli` — the full lifecycle a single author needs to
   operate a site end to end, mirroring the lifecycle named in the Project
   description and Principle III (CLI-First, Guided Experience).

A release directory is the closing artifact of the feature pipeline above —
it is assembled once the features it bundles have shipped, not drafted
speculatively ahead of the work it documents.

## Governance

This constitution supersedes any ad hoc convention or prior practice where the
two conflict. Amendments require: (1) a proposed change to this file with a
clear rationale, (2) a version bump per the policy below, (3) propagation of
any changed rule into `AGENTS.md`, `CLAUDE.md`, and the affected
`.specify/templates/*` files in the same change set, and (4) a Decision Log
entry in the relevant research document if the amendment reflects a new or
changed architectural/technical decision (per the workflow above).

Versioning policy (semantic versioning for governance):
- **MAJOR** — a principle is removed or redefined in a backward-incompatible
  way (e.g., reversing the no-auth posture, abandoning the CLI-first surface).
- **MINOR** — a new principle or section is added, or existing guidance is
  materially expanded (e.g., adding a new escalation rule or workflow gate).
- **PATCH** — wording, clarification, typo, or other non-semantic refinement.

Every plan, spec, task breakdown, or research document that touches an area
governed by a principle above MUST state how it complies, or document and
justify any deviation as a deliberate, scoped exception (never as silent
drift). Reviewers (i.e., the project owner, since this is a single-author
project) check new work against these principles before merging — complexity
that cannot be justified against Principle I is the default reason to push
back. Use `AGENTS.md` for day-to-day runtime development guidance; this
constitution governs the *non-negotiable* rules beneath it.

**Version**: 1.3.1 | **Ratified**: 2026-06-07 | **Last Amended**: 2026-06-20
