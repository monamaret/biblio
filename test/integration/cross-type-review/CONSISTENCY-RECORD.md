# Cross-Type Consistency Record

Closing artifact for v0.05's cross-type review (data-model.md "Cross-Type
Consistency Record"; FR-003, FR-004, SC-002). Produced by running
`go test ./test/integration/cross-type-review/...` against all three site
types after the documentation type's replication (Phase 3) completed.

## sites[]

- profile (v0.02)
- log (v0.04)
- docs (this release, v0.05)

## theming-check

| Site type | Result |
|---|---|
| profile | Renders its own preset (`profile-default`) exactly — `TestCrossTypeTheming_EachSiteRendersItsOwnPreset/profile` |
| log | Renders its own preset (`log-default`) exactly — `TestCrossTypeTheming_EachSiteRendersItsOwnPreset/log` |
| docs | Renders its own preset (`docs-default`) exactly — `TestCrossTypeTheming_EachSiteRendersItsOwnPreset/docs` |

No cross-type bleed-through: each type's distinguishing tokens
(`color-fg`, `color-accent`) are unique per type
(`TestCrossTypeTheming_NoBleedThrough`). `color-bg` is intentionally shared
between profile and docs (two light themes vs. log's dark one), and
`font-heading` is intentionally shared across all three — both confirmed,
deliberate design choices, not bleed-through (see test comments).

## terminal-rendering-check

| Site type | Result |
|---|---|
| profile | Clean Markdown/ANSI rendering, type-agnostic negotiation code (`templates/site/profile/app/negotiation/accept.test.ts`) |
| log | Clean Markdown/ANSI rendering against dated-entry content (`templates/site/log/app/negotiation/accept.test.ts`) |
| docs | Clean Markdown/ANSI rendering against hierarchical-page content (`templates/site/docs/app/negotiation/accept.test.ts`) |

One shared negotiation owner confirmed: `accept.ts`/`ansi.ts`/`wasm.d.ts`
are byte-identical across all three site types
(`TestCrossTypeTerminalRendering_OneSharedNegotiationOwner`) — no per-type
fork exists for any client to render inconsistently against. Each type's
content shape is appropriate to its own model — profile's `role`/`summary`,
log's `date`/`tags`, docs's `section`/`order` — with no cross-contamination
of another type's fields (`TestCrossTypeTerminalRendering_ContentShapeIsAppropriatePerType`).

## inconsistencies-found[]

None. (One over-strict test assertion was found and corrected during
authoring — see Phase 4's T029 note in tasks.md — but it was a defect in
the *review test itself*, not in `internal/theme` or the negotiation code;
no shared-code defect was filed.)

---

## Four-type review — quarterly (specs/004-quarterly-site-type, US6)

Re-run on 2026-10-15 with `go test ./test/integration/cross-type-review/...`
after the review tests were switched from a hard-coded three-type list to
`scaffold.SiteTypes`, so any future type is reviewed automatically.

### sites[]

- profile, log, docs (as above)
- quarterly (feature 011, targeting v0.08)

### theming-check

| Site type | Result |
|---|---|
| quarterly | Renders its own preset (`quarterly-default`) exactly — `TestCrossTypeTheming_EachSiteRendersItsOwnPreset/quarterly` — and also carries its typesetting properties (`--ts-*`) in the same single `:root` block — `TestCrossTypeTheming_QuarterlyCarriesTypesetting` |

No bleed-through: quarterly's `color-fg` (`#111111`) and `color-accent`
(`#d7261e`) are distinct from the other three presets'
(`TestCrossTypeTheming_NoBleedThrough`).

### terminal-rendering-check

| Site type | Result |
|---|---|
| quarterly | Clean Markdown/ANSI rendering of issue contents, archive, category and article pages; articles serve the Go-produced `markdown_body`, so margin notes arrive as footnotes (`templates/site/quarterly/app/negotiation/{accept,markdown}.test.ts`) |

`accept.ts`/`ansi.ts`/`wasm.d.ts` remain byte-identical across all four
types (`TestCrossTypeTerminalRendering_OneSharedNegotiationOwner`).
Quarterly's content shape — `kind`, `issue`, `category`, `footnotes`,
`margin_notes`, `markdown_body` — carries none of the other types' fields
(`TestCrossTypeTerminalRendering_ContentShapeIsAppropriatePerType/quarterly`).

### inconsistencies-found[]

- **Generated `content-index.json` treated as a template region by
  `update`/`upgrade`.** Every template ships a starter `content-index.json`,
  so it was manifest-tracked; once `build`/`dev` began regenerating it
  automatically (all four types), any template release that changed its
  shipped index would have made `update` report a conflict on every site.
  Fixed in `internal/migrate/classify.go` via `scaffold.IsGenerated`:
  `content-index.json`, `public/tokens.css` and `public/plates/` are build
  output and are never classified
  (`test/integration/update-three-sites/quarterly_site_test.go`).
- **Trailing-slash detail URLs after hydration.** Workers Static Assets
  serves prerendered pages at `/slug/`; the detail route's splat then
  carried the slash and the page re-rendered as "Not found" once React
  hydrated. Fixed for quarterly (`templates/site/quarterly/app/routes/detail.tsx`).
  The profile, log and docs detail routes
  (`templates/site/{profile,log,docs}/app/routes/detail.tsx`) read the splat
  the same way and are likely affected too; not changed here, since they
  are outside this feature's scope — flagged for follow-up.
