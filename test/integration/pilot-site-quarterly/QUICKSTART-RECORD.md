# Quickstart record — `quarterly` site type

Results of running [specs/004-quarterly-site-type/quickstart.md](../../../specs/004-quarterly-site-type/quickstart.md)
scenarios 1–7 against real scaffolded sites (tasks.md T099), 2026-10-15.
Each scenario was run when its user story landed, against a site
scaffolded and built with the `biblio-cli` binary of that commit;
scenario 1 was repeated at the end on a fresh scaffold with the final
binary. Node 22.15.1 (the template's CI uses Node 20), Go 1.24 module
built with go1.27, Chromium headless for the visual checks.

| # | Scenario | Result | Evidence |
|---|---|---|---|
| 1 | Scaffold, build and browse (US1 · SC-001, SC-002) | Pass | Fresh `new --type=quarterly` + `build`: validate, tokens, plates (2 converted), index (8 items); 12 pages prerendered (`/`, `/issues`, 2 issues, 2 categories, 5 articles, `/colophon`) plus the SPA fallback; `llms.txt` written. Contents page reviewed at 1440 px: masthead, `No. 2 · Autumn 2026`, categories in order, roman numerals, dotted leaders, small-caps authors |
| 2 | Notes (US2 · SC-005) | Pass | `on-margins`: 3 sidenotes and 3 margin notes in the HTML. Screenshots at 320, 375, 1024, 1440 and 1920 px: wide screens put both kinds in the margin column without overlap; narrow screens show the foot list and † toggles; also checked with JavaScript disabled. Markdown view lists every note as a footnote (`[^m1]`…). One defect found and fixed: detail pages served at `/slug/` turned into "Not found" after hydration (trailing slash in the splat) |
| 3 | Typesetting (US3 · SC-003, SC-008) | Pass | `accent`, `align: justify`, `numerals: arabic` → `tokens.css` header lists the three overrides; contents numbered 1, 2, 3. Each of `scale: 2`, `align: center`, `accent: red`, `font_body: "x;}"`, `colour: "#000"` exits 1 before Vite runs, naming the key. `font_body: "Söhne"` builds with a warning |
| 4 | Plates (US4 · SC-004) | Pass | Both sample images published as single-channel grayscale JPEGs in `public/plates/<hash16>.jpg`; every `<img>` points at `/plates/…`; a rerun reports `0 converted, 2 skipped`; a truncated JPEG stops the build naming the file. Colour mode, PNG alpha, WebP and the 40 MP header cap are covered by `internal/plates` tests (plus a 30 s fuzz run) |
| 5 | Drafts (FR-012) | Pass | `draft: true` on issue No. 2: `/` shows No. 1, `dist/client/issues/` holds only `01-spring-2026`, `llms.txt` and `content-index.json` omit it, build reports `1 draft issue(s) excluded`. Dev-mode inclusion covered by `internal/index` tests |
| 6 | Agent tools (US5 · SC-006) | Pass | Over streamable HTTP against `biblio-cli dev` (wrangler dev): `tools/list` shows the 3 shared + 4 quarterly tools; `list_issues` table newest first; `get_issue "1"` grouped contents; `list_articles {category: essays, author: A. Novak}` → 2 articles; `get_article` text with every note plus `structuredContent` (3 footnotes, 3 margin notes); `get_issue "99"` → not-found text. Markdown and ANSI negotiation and plate serving checked on the same server |
| 7 | Other types unchanged (US6 · SC-007) | Pass | `go test ./...` green, incl. `prebuild-regression` (profile/log/docs output identical apart from `content-index.json`), the four-type cross-type review, and the quarterly update/upgrade test. All four templates' Vitest suites (41/37/30/95) and `tsc` pass. A fresh `log` site builds with no manual `index` run, `tokens.css` unchanged, no plates |
