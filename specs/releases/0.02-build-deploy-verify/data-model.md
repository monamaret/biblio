# Phase 1 Data Model: v0.02 — Build, deploy, and verify a real site

**Input**: [spec.md](./spec.md) FR-001..FR-004, [research.md](./research.md) prerendering/deployment decisions

This release adds no new persistent data shapes — it carries the v0.01 `content-index.json`/`tokens.css` outputs through to a deployed artifact unchanged (spec Assumptions: "Builds on v0.01's output unchanged"). The entities below describe the *build and deployment artifacts* this release produces and the contract each must satisfy, not a new schema.

## Entity: Pre-rendered Page Set (build output)

The static output of `biblio-cli build`, produced by React Router v7's `prerender` from the unchanged `content-index.json`.

| Field | Type | Notes |
|---|---|---|
| `pages[]` | array of pre-rendered routes | One per content-driven page, including list/detail and navigation views appropriate to the site type (FR-001) |
| `source-index` | reference | The exact `content-index.json` produced by v0.01's indexing step — read-only input, not regenerated or reshaped here (Principle II, spec Assumptions) |
| completeness | invariant | Every entry in `source-index` MUST correspond to exactly one rendered page; zero missing or broken pages (spec Scenario 1, FR-001) |

**Validation rules**: a build that encounters content failing validation (e.g. malformed frontmatter) MUST fail with a clear message rather than silently skip the page (spec Edge Cases) — this is the behavior the failing-test-first build-wrapper unit must encode (Principle VI).

**Relationships**: derived 1:1 from the Content Index entity defined in [v0.01's data-model](../0.01-scaffold-and-theme/data-model.md); consumed by the Deployment entity below.

## Entity: Deployment

The published, publicly reachable instance of a Pre-rendered Page Set, produced by `biblio-cli deploy` (FR-002).

| Field | Type | Notes |
|---|---|---|
| `public-address` | URL (custom domain / route) | The configured address from `wrangler.jsonc` (generated in v0.01, route pre-filled from `--domain`) |
| `worker-binding` | reference | The single Cloudflare Worker (Workers Static Assets) serving this deployment — same Worker that will later also mount `/mcp` in v0.03 (Principle V topology) |
| `theme` | reference | The `tokens.css` generated in v0.01, served and applied unchanged (FR-003) |
| re-deploy semantics | invariant | Re-running deploy against an already-live site MUST update cleanly in place with no visitor-facing disruption (spec Edge Cases) |

**Validation rules**: after deploy, the site MUST be reachable at `public-address` and serve every page from the Pre-rendered Page Set correctly, styled per `theme` (FR-003, Scenario 3).

**Relationships**: produced 1:1 from a Pre-rendered Page Set; one Deployment exists per scaffolded project in this release's scope (single pilot site — spec Assumptions).

## Entity: Pilot Site Validation Record

The closing evidence artifact for this release — not a runtime entity, but the thing SC-001/SC-002 are checked against.

| Field | Type | Notes |
|---|---|---|
| `site` | reference | The one real, fully authored personal-profile/CV site this release is demonstrated against (FR-004, research.md "single real pilot site") |
| `path-completed` | checklist | content authoring → `biblio-cli build` → `biblio-cli deploy` → live human-visitor verification, completed without manual intervention beyond documented commands (SC-002) |
| `pages-verified` | percentage | MUST be 100% of authored pages rendering correctly and matching the chosen theme on human-visitor check (SC-001) |

**Relationships**: references exactly one Deployment and one Pre-rendered Page Set; this record is what the release-closing integration test in `test/integration/pilot-site/` exercises and what the constitution's Release Process "integration tests that exercise the release's features end to end" requirement is satisfied by.
