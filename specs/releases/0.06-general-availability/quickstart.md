# Quickstart Validation: v0.06 — General Availability

Validates spec Success Criteria SC-001..SC-003 — the platform's final release-closing checks. Requires all three real sites from [v0.02](../0.02-build-deploy-verify/quickstart.md)/[v0.04](../0.04-project-log-site-type/quickstart.md)/[v0.05](../0.05-docs-site-and-consistency/quickstart.md) to already be live (this release's central design work was only answerable once they existed — research.md).

## Prerequisites

- The three real, live deployed sites: profile (v0.02), project log (v0.04), documentation (v0.05)
- A newer scaffold-template version to update/upgrade toward (the release's own deliverable)
- `biblio-cli` extended with `update`/`upgrade` commands ([contracts/update-and-escalation.md](./contracts/update-and-escalation.md))

## Scenario 1 — Update a site without losing content (spec Scenario 1)

```bash
cd my-profile-site
biblio-cli update
```

**Expected outcome**: scaffold-tracked parts are brought up to date, authored content and intentional customizations are preserved untouched, and any genuine conflicts are clearly surfaced for the author to resolve — never silently overwritten (FR-001, FR-002; [data-model.md](./data-model.md) Update Run / Update Conflict). Interrupt a run mid-flight (`Ctrl+C`) and confirm the site remains in a consistent, recoverable state (spec Edge Cases).

## Scenario 2 — The update flow works across all three site types (spec Scenario 2)

```bash
for site in my-profile-site my-project-log my-docs-site; do
  (cd "$site" && biblio-cli update)
done
```

**Expected outcome**: the flow completes successfully for every site type — not only whichever it may first have been validated against (FR-003; `test/integration/update-three-sites/`; this is what SC-001's "carry forward all three real, live sites... in one sitting" is checked against).

## Scenario 3 — Storage-escalation guidance is checked against real sites (spec Scenario 3)

```text
For each of the three real sites, check it against
docs/storage-escalation-guide.md's documented thresholds.
```

**Expected outcome**: the guidance correctly reports whether any site has crossed a D1 (volume/relational) or Vectorize (search-quality) threshold, reflecting each site's real, current scale — and for any path it recommends, its steps are concretely walkable, proven via an actual dry-run against real site data, not just described in principle (FR-004, FR-005; [contracts/update-and-escalation.md](./contracts/update-and-escalation.md) Storage-Escalation Guidance).

## Scenario 4 — GA readiness review (spec Scenario 4)

```text
Cross-reference every roadmap deliverable §1–§10, and every release v0.01–v0.06,
against this release's acceptance bar.
```

**Expected outcome**: every deliverable is either complete-and-demonstrated, or — for the two explicitly-conditional storage-escalation items — proven executable on demand, marking the platform as Generally Available ([data-model.md](./data-model.md) GA Readiness Review; research.md "What GA means here").

## Demonstrating SC-001..SC-003

- **SC-001**: Scenarios 1–2 carry forward all three real, live sites to a newer scaffold version in one sitting, each finishing with authored content fully intact and any genuine conflicts clearly identified rather than silently lost.
- **SC-002**: Scenario 4's complete review against every roadmap deliverable §1–§10 shows each one complete-and-demonstrated — or, for the two explicitly-conditional storage-escalation items, shown executable on demand even though no real site has yet required them.
- **SC-003**: This release's update-flow evidence (Scenarios 1–2) never requires re-proving v0.04/v0.05's site-type behaviors from scratch, and vice versa — each release's evidence, including this one, stands on its own.
