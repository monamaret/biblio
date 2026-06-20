# v0.01 — Scaffold and theme a site

**Created**: 2026-06-07
**Status**: Draft
**Roadmap deliverables**: [§1 `biblio-cli` foundations](../research/roadmap.md#1-biblio-cli-foundations) · [§2 Content pipeline](../research/roadmap.md#2-content-pipeline--content-indexjson-generator) · [§6 Theming system](../research/roadmap.md#6-theming-system)
**Depends on**: nothing — this is the foundation release every later one builds on
**Feature specs**: [001](../features/001-biblio-cli-foundations.md) · [002](../features/002-content-pipeline.md) · [006](../features/006-theming-system.md)

## Summary

The sole `biblio` author can run one CLI command and get a fully scaffolded, themed site project — content directory, config, build/deploy plumbing, and a generated content index — ready to fill with their own writing, without hand-assembling any of the underlying tooling. This is the foundation layer the roadmap itself names as "next up" and as the thing "every other deliverable scaffolds from."

## User Scenarios & Acceptance Criteria

### Scenario 1 — Scaffold a new site from one command

As the author, I run a single guided command, choose a site type/domain/theme, and get back a complete, runnable project.

1. **Given** no existing project, **When** the author runs the scaffolding command with a project name, site type, domain, and theme, **Then** a complete combined site project is generated locally with seeded content, generated config files, and a matching theme applied.
2. **Given** a scaffolded project, **When** the author inspects it, **Then** it contains a populated `content/` tree, a `site.yaml`/`wrangler.jsonc` pair, and a `tokens.css` matching the chosen theme preset.

### Scenario 2 — Index authored content

As the author, I run the indexing step over my content directory and get back one canonical, structurally correct index.

1. **Given** a scaffolded project with authored content, **When** the author runs the content-indexing step, **Then** a single canonical content index is produced that correctly reflects every authored page/entry and its derived slug/path.

### Scenario 3 — Preview locally

As the author, I start a local preview with one command and see my content live.

1. **Given** a scaffolded project, **When** the author runs the local development command, **Then** the site and its supporting services start together behind one command and the author can preview their content live.

## Functional Requirements

- **FR-001**: The system MUST let the author scaffold a complete, runnable combined site project (content + build/deploy plumbing) for any of the three supported site types from a single guided command.
- **FR-002**: The system MUST generate, as part of scaffolding, the site's configuration files (including its declared site type, name/description, and agent-facing on/off setting), its deployment-routing configuration, and its generated visual-theme file.
- **FR-003**: The system MUST produce a single canonical content index from an author's content directory, correctly reflecting every authored page/entry's structured fields and derived slug/path, in a shape usable by both the rendered site and its agent-facing layer (built in v0.03).
- **FR-004**: The system MUST provide a distinct, ready-to-use visual-theme preset for each of the three site types, selectable at scaffold time.
- **FR-005**: The system MUST let the author run a local preview of their scaffolded site (content + supporting services together) from a single command.
- **FR-006**: The system MUST let the author distribute/install the CLI tool itself via at least one straightforward method (prebuilt download or toolchain-based install).

## Edge Cases

- What happens when the author scaffolds a site type (project log or documentation) that won't be carried through to a real deployment until v0.04/v0.05 — is the scaffold still fully usable on its own, just not yet proven live?
- How does the indexing step handle authored content with malformed or missing structured fields — is the problem surfaced clearly rather than silently producing an incomplete index?

## Success Criteria

- **SC-001**: A person can go from "no project" to "a running local preview of a themed, content-indexed site" using only the documented CLI commands, with no manual file assembly or configuration editing required.
- **SC-002**: All three site-type scaffolds (profile, project log, documentation) produce a runnable local preview and a structurally correct content index on the first attempt.
- **SC-003**: This release is independently demonstrable purely on local-preview evidence — it requires no live deployment (that's v0.02) and no agent/terminal-facing capability (that's v0.03) to show its value.

## Assumptions

- **Single-author scope**: Per the roadmap's "sole-developer, personal-scale project" framing, this release is demonstrated by one author scaffolding their own projects, not a multi-user rollout.
- **No timeline commitments**: Consistent with the roadmap's explicit choice to avoid effort/timeline estimates, this spec defines what the release contains and how its completion is demonstrated, not when it ships.
- **Theming depends on its own preset design**: `tokens.css` generation (§1) depends on the preset definitions produced by §6; both are scoped into this release together so the dependency is satisfied within one release rather than split across two.
