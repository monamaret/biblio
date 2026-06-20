# Phase 0 Research: v0.02 — Build, deploy, and verify a real site

**Input**: [plan.md](./plan.md) Technical Context — no `NEEDS CLARIFICATION` markers remained; every relevant decision was already settled in [`specs/research/001-static-site-platform-with-mcp.md`](../../research/001-static-site-platform-with-mcp.md) and [roadmap §0.3](../../research/roadmap.md) before this plan was drafted, per the constitution's research-precedes-planning workflow.

## Decision: Pre-rendering approach

- **Decision**: React Router v7's native `prerender`, configured in `react-router.config.ts`, pre-renders every content-driven page (list/detail/navigation per site type) into static output as part of `biblio-cli build`.
- **Rationale**: Confirmed in [roadmap §3](../../research/roadmap.md) ("Set up the React + Vite + React Router v7 project skeleton, with native `prerender` configured in `react-router.config.ts`") — using the framework's native mechanism avoids hand-rolled static-generation logic and keeps the site-build pipeline owned by one well-supported tool.
- **Alternatives considered**: A custom static-site-generation step bolted onto Vite — rejected; would duplicate functionality React Router v7 already provides natively, adding maintenance surface Principle I rules out.

## Decision: Content loading and rendering in the prerendered template

- **Decision**: `import.meta.glob` + `gray-matter` for content loading, rendered via `react-markdown` + `remark-gfm`.
- **Rationale**: Confirmed stack choice in [roadmap §3](../../research/roadmap.md); these are the lightest-weight, most direct tools for "load Markdown+frontmatter, render it" with no MDX (explicitly rejected per Principle I's rationale — "no MDX" is named directly in the constitution).
- **Alternatives considered**: MDX-based rendering — explicitly rejected in the confirmed Decision Log (Principle I rationale cites "no MDX" by name as a confirmed decision this bias is built on).

## Decision: Build/deploy command surface

- **Decision**: `biblio-cli build`/`biblio-cli deploy` are thin wrappers over `vite build`/`wrangler deploy` plus a generated CI workflow; the update/upgrade story (v0.06) is folded into this same command surface (`biblio-cli update`/`upgrade`) rather than split out.
- **Rationale**: Confirmed in [roadmap §0.3](../../research/roadmap.md) — keeps one consistent CLI-first surface end to end rather than splitting "dev/build/deploy is one command, updates are raw `package.json`/CI," which Principle III explicitly forbids.
- **Alternatives considered**: Documenting `vite build`/`wrangler deploy` as commands the author runs directly — rejected; that is exactly the "raw script as the primary surface" Principle III calls out as not acceptable.

## Decision: Deployment topology

- **Decision**: One Cloudflare Worker serves the prerendered static site via Workers Static Assets, mounted on the route/custom-domain binding generated in v0.01's `wrangler.jsonc` — not split across separate hosts/services, and not yet carrying the `/mcp` route (that's v0.03).
- **Rationale**: Confirmed architecture decision (Principle V, "Simplest Topology That Can Grow" — "One Cloudflare Worker serves both the prerendered static site... and the `/mcp` endpoint... never split across separate hosts/services"). This release builds the first half of that topology; the second half (the `/mcp` mount) arrives in v0.03 onto the *same* Worker, not a new one.
- **Alternatives considered**: A separate static-hosting service (e.g. Pages) plus a separate Worker for future MCP — rejected; would pre-split a topology the constitution names as deliberately unified.

## Decision: Validation approach — single real pilot site

- **Decision**: This release is demonstrated and proven exclusively against one real, fully authored personal-profile/CV pilot site — not a synthetic fixture, and not all three site types.
- **Rationale**: Confirmed in [roadmap §8](../../research/roadmap.md) (pilot site validation exists specifically to "validate the scaffold + CLI + architecture before replication") and restated in this release's own Assumptions — proving the other two site types live is deliberately deferred to v0.04/v0.05, so this release's evidence stands on its own (spec SC-003).
- **Alternatives considered**: Demonstrating against a synthetic/seed-only site — rejected; would not prove the real authoring → deploy → live-verification path the release exists to validate (FR-004).

## Open questions resolved

None remain — Technical Context was filled entirely from confirmed decisions in `001-static-site-platform-with-mcp.md` and `roadmap.md` §0.3/§3/§8, with no `NEEDS CLARIFICATION` markers requiring further research.
