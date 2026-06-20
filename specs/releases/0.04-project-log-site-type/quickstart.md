# Quickstart Validation: v0.04 — Prove the platform on the technical project log

Validates spec Success Criteria SC-001..SC-003, demonstrable purely on this site type's own evidence — no documentation-site replication or cross-type consistency check required (those are [v0.05](../0.05-release.md)). Reuses the v0.01–v0.03 CLI/quickstart flows verbatim, exercised against `--type=log`.

## Prerequisites

- `biblio-cli` from v0.01–v0.03 (no new commands), Cloudflare/Wrangler auth configured (as in [v0.02](../0.02-build-deploy-verify/quickstart.md))

## Scenario 1 — Author and publish a run of dated entries (spec Scenario 1)

```bash
biblio-cli new my-project-log --type=log --domain=log.example.com --theme=log-default
cd my-project-log
# author several dated, tagged entries (some sharing tags or a project) under content/
biblio-cli index && biblio-cli build && biblio-cli deploy
```

**Expected outcome**: every entry appears correctly in the deployed site's entry list with its date, tags, and project correctly attributed (FR-001, FR-002; [contracts/log-site-surface.md](./contracts/log-site-surface.md)).

## Scenario 2 — Filter the entry list (spec Scenario 2)

```text
On the deployed site, filter the entry timeline by:
  - a date range
  - a tag
  - a project
  - a date range AND a tag together
```

**Expected outcome**: only genuinely matching entries are shown for each filter, with the combined filter returning the intersection — no false positives or omissions (FR-003; [data-model.md](./data-model.md) Filter Query / Filtered Entry Set; spec Edge Cases). Add an entry with a brand-new tag, rebuild/redeploy, and confirm it appears in filters immediately with no extra registration step.

## Scenario 3 — Keep up a daily authoring rhythm (spec Scenario 3)

```bash
# Repeat, for several consecutive days:
#   author one new dated entry the same way as the last → index → build → deploy
```

**Expected outcome**: each new entry appears in the deployed site after the normal build/deploy step, with zero one-off adjustments to structure or configuration (FR-004; [contracts/log-site-surface.md](./contracts/log-site-surface.md) authoring-rhythm contract; SC-002).

## Scenario 4 — The agent/terminal layer works for this site type too (spec Scenario 4)

```bash
# Repeat the v0.03 quickstart's Scenario 1 (MCP connect/list/call) and Scenario 3
# (terminal/curl Markdown rendering) against this deployed project-log site.
# Then issue a profile-shaped query (e.g. "get résumé") against it.
```

**Expected outcome**: the same agent- and terminal-facing capabilities validated in [v0.03](../0.03-agent-terminal-layer/quickstart.md) behave correctly for this site type's own structured-query and terminal-rendering surface (FR-005); the profile-shaped query is rejected clearly, reusing the v0.03 "unsupported capability" contract (spec Edge Cases).

## Demonstrating SC-001..SC-003

- **SC-001**: A real technical project-log site is live, and 100% of its authored entries are correctly discoverable through every filter dimension (date, tag, project) by a human visitor (Scenario 1–2).
- **SC-002**: The author follows the same day-to-day authoring rhythm for at least several consecutive published entries with zero one-off adjustments (Scenario 3).
- **SC-003**: All of the above is demonstrable purely on this one site type's evidence — no documentation-site replication or cross-type consistency check (both [v0.05](../0.05-release.md)) is required to show this release's value.
