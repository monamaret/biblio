# Quickstart Validation: v0.02 — Build, deploy, and verify a real site

Validates spec Success Criteria SC-001..SC-003 end to end, purely via live-browser evidence (no agent/terminal-facing capability required — that's v0.03). Builds directly on the [v0.01 quickstart](../0.01-scaffold-and-theme/quickstart.md) — start from a scaffolded, indexed, themed profile/CV project.

## Prerequisites

- A scaffolded, content-indexed, themed profile/CV project from [v0.01](../0.01-scaffold-and-theme/quickstart.md) (`biblio-cli new ... --type=profile`, authored with real content for the pilot)
- Cloudflare account + Wrangler authentication configured (required for `biblio-cli deploy`)

## Scenario 1 — Build a site from authored content (spec Scenario 1)

```bash
cd my-profile-site
biblio-cli build
```

**Expected outcome**: every content-driven page (profile, experience, skills, projects — list/detail/navigation as appropriate) is correctly pre-rendered into static output with no missing or broken pages (FR-001; [contracts/cli-commands.md](./contracts/cli-commands.md) `build`). Introduce a malformed-frontmatter file and re-run — expect the build to fail with a clear, per-file message rather than silently skip the page (spec Edge Cases).

## Scenario 2 — Deploy and reach the site live (spec Scenario 2)

```bash
biblio-cli deploy
curl -I https://<configured-domain>
```

**Expected outcome**: the site becomes reachable at its configured public address and serves the pre-rendered pages correctly (FR-002, FR-003). Re-run `biblio-cli deploy` against the now-live site — expect a clean in-place update with no visitor-facing disruption (spec Edge Cases).

## Scenario 3 — A visitor browses the live site (spec Scenario 3)

```text
Open https://<configured-domain> in a browser.
```

**Expected outcome**: the profile, experience, skills, and projects views all display the authored content correctly, styled per the theme generated in v0.01 (FR-003; [data-model.md](./data-model.md) Deployment entity).

## Demonstrating SC-001..SC-003

- **SC-001**: At least one real, fully authored site is live at a public address; a human visitor confirms 100% of its authored pages render correctly and match its chosen theme ([data-model.md](./data-model.md) Pilot Site Validation Record).
- **SC-002**: The full content → build → deploy → live-verification path (Scenarios 1–3 above, run in order) completes with no manual intervention beyond these documented commands.
- **SC-003**: All of the above is demonstrable purely from live-browser evidence — no agent/terminal-facing capability (v0.03) is invoked, and showing it never requires temporarily wiring that capability up.
