# v0.02 — Build, deploy, and verify a real site

**Created**: 2026-06-07
**Status**: Draft
**Roadmap deliverables**: [§3 Site scaffold](../research/roadmap.md#3-site-scaffold--react--vite--worker-combined-project-template) · [§8 Pilot site validation](../research/roadmap.md#8-pilot-site-validation)
**Depends on**: [v0.01](0.01-release.md) — needs a scaffolded, indexed, themed project to build and deploy
**Feature specs**: [003](../features/003-site-scaffold.md) · [008](../features/008-pilot-site-validation.md)

## Summary

The author takes the scaffold from v0.01, fills it with their own profile content, and runs one command to build and deploy it — ending with a real, publicly reachable site that proves the whole pipeline (content → build → deploy) works, not just the scaffolding step. This is the first release that produces something real visitors can reach, and the roadmap names pilot validation as existing specifically to "validate the scaffold + CLI + architecture before replication."

## User Scenarios & Acceptance Criteria

### Scenario 1 — Build a site from authored content

As the author, I run the build command and get a fully pre-rendered, deployable site with no missing pages.

1. **Given** a scaffolded profile/CV project with authored content, **When** the author runs the build command, **Then** every content-driven page is correctly pre-rendered into static output with no missing or broken pages.

### Scenario 2 — Deploy and reach the site live

As the author, I run the deploy command and the site becomes reachable at a public address.

1. **Given** a successful build, **When** the author runs the deploy command, **Then** the site becomes reachable at its configured public address and serves the pre-rendered pages correctly.

### Scenario 3 — A visitor browses the live site

As a visitor, I open the deployed site and see the author's content correctly, styled per its theme.

1. **Given** a deployed site, **When** a visitor browses it, **Then** the profile, experience, skills, and projects views all display the authored content correctly styled per the chosen theme.

## Functional Requirements

- **FR-001**: The system MUST pre-render every content-driven page of a scaffolded site (including list/detail and navigation views appropriate to its site type) into static output as part of a single build command.
- **FR-002**: The system MUST let the author publish a built site to a publicly reachable address from a single deploy command.
- **FR-003**: The system MUST serve the deployed site's pre-rendered pages correctly to ordinary browser visitors, styled per the site's chosen theme.
- **FR-004**: The system MUST be demonstrated, by the close of this release, on at least one real, fully authored personal-profile/CV site carried through the entire content → build → deploy → live-verification path.

## Edge Cases

- What happens when the build encounters authored content that fails validation (e.g., a malformed frontmatter field) — does it fail the build with a clear message, or silently skip the page?
- How does the deploy step behave if it's re-run against a site that's already live — does it update cleanly in place with no visitor-facing disruption?

## Success Criteria

- **SC-001**: At least one real, fully authored site is live at a public address, and 100% of its authored pages render correctly and match its chosen visual theme when checked by a human visitor.
- **SC-002**: The full content → build → deploy → live-verification path completes without manual intervention beyond the documented commands.
- **SC-003**: This release is independently demonstrable purely on live-browser evidence — it requires no agent/terminal-facing capability (that's v0.03) to show its value, and showing it never requires temporarily wiring up v0.03's capability.

## Assumptions

- **Single pilot site, single site type**: Per the roadmap's rollout plan, this release proves the pipeline on one personal-profile/CV pilot — proving the other two site types live is explicitly deferred to v0.04/v0.05.
- **No timeline commitments**: Consistent with the roadmap's explicit choice to avoid effort/timeline estimates, this spec defines what the release contains and how its completion is demonstrated, not when it ships.
- **Builds on v0.01's output unchanged**: This release assumes the scaffold, content index, and theme produced by v0.01 are stable inputs — it doesn't redesign any of them, only proves they can be carried through to a live deployment.
