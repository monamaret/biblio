# Quickstart Validation: v0.05 — Prove the platform on documentation and across all site types

Validates spec Success Criteria SC-001..SC-003, demonstrable on this site type's evidence plus the cross-type review — no update/upgrade-flow work required (that's [v0.06](../0.06-general-availability/quickstart.md)). Reuses the v0.01–v0.04 CLI/quickstart flows verbatim, exercised against `--type=docs`, plus a new comparison pass across all three now-live sites.

## Prerequisites

- `biblio-cli` from v0.01–v0.04 (no new commands)
- The two already-live sites from [v0.02](../0.02-build-deploy-verify/quickstart.md) (profile) and [v0.04](../0.04-project-log-site-type/quickstart.md) (project log)

## Scenario 1 — Browse a documentation site by its hierarchy (spec Scenario 1)

```bash
biblio-cli new my-docs-site --type=docs --domain=docs.example.com --theme=docs-default
cd my-docs-site
# author pages with `section`/`order` frontmatter, organized into a TOC hierarchy
biblio-cli index && biblio-cli build && biblio-cli deploy
```

```text
On the deployed site:
  1. Navigate via the table of contents — confirm it reflects the section/order
     structure and every page is reachable and renders correctly (FR-002, Scenario 1.1).
  2. Look for any tag- or project-based filtering (as exists on the project-log site)
     — confirm none is present, anywhere in the UI or query surface (Scenario 1.2).
```

**Expected outcome**: full hierarchical navigation works end to end, and the deliberate absence of filter UI/params is confirmed — not found missing-by-accident, but verified present-by-design ([contracts/docs-site-and-review.md](./contracts/docs-site-and-review.md); FR-002).

## Scenario 2 — Compare theming across all three site types (spec Scenario 2)

```text
Visit the profile site (v0.02), the project-log site (v0.04), and this
documentation site in turn.
```

**Expected outcome**: each displays its own distinct, correctly-applied visual theme with no cross-type bleed-through or inconsistency (FR-003; [data-model.md](./data-model.md) Cross-Type Consistency Record — theming check). Any inconsistency found is filed as a defect against the shared theming system, not patched per-site.

## Scenario 3 — Compare terminal rendering across all three site types (spec Scenario 3)

```bash
for site in <profile-domain> <log-domain> <docs-domain>; do
  curl -H "Accept: text/markdown" "https://$site/"
done
```

**Expected outcome**: each of the three deployed sites returns a clean, correctly rendered, type-appropriate terminal-friendly view (FR-004; [data-model.md](./data-model.md) Cross-Type Consistency Record — terminal-rendering check). Also send a project-log-style tag-filter parameter to the documentation site and confirm it's handled via the documented, sensible default — clearly rejected or cleanly ignored, never silently mishandled (spec Edge Cases).

## Demonstrating SC-001..SC-003

- **SC-001**: A real project-documentation site is live, and a visitor reaches and correctly reads every authored page through the TOC/navigation structure alone (Scenario 1).
- **SC-002**: All three supported site types are simultaneously live and real, each visually distinct per its own theme and each returning a correct terminal/Markdown rendering — with **zero** cross-type inconsistencies remaining open at the close of the side-by-side review (Scenarios 2–3; [data-model.md](./data-model.md) Cross-Type Consistency Record).
- **SC-003**: All of the above is demonstrable purely on this site type's evidence plus the cross-type review — no update/upgrade-flow work ([v0.06](../0.06-general-availability/quickstart.md)) is required to show this release's value.
