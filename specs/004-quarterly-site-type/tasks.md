---

description: "Task list for the quarterly site type"
---

# Tasks: `quarterly` site type

**Input**: Design documents from `specs/004-quarterly-site-type/`

**Prerequisites**: [plan.md](plan.md), [spec.md](spec.md), [research.md](research.md), [data-model.md](data-model.md), [contracts/](contracts/), [quickstart.md](quickstart.md)

**Tests**: **Required.** Constitution Principle VI (test-first, non-negotiable). In every phase, write the test tasks first, run them, and confirm they **fail for the right reason** before starting the implementation tasks that make them pass.

**Organization**: tasks are grouped by user story (US1–US6 from spec.md) so each story can be built and checked on its own.

## Format: `[ID] [P?] [Story] Description`

- **[P]**: can run in parallel (different files, no dependency on an unfinished task)
- **[Story]**: the user story the task belongs to (US1–US6)
- Paths are repo-relative. `Q/` below is shorthand for `templates/site/quarterly/` and appears **only in this legend**; every task spells out the full path.

## Conventions used by every task

- Go code follows the style of its package (doc comments that name the requirement or contract they satisfy, collected errors, table tests).
- Template code follows `templates/site/log/` (same file names, same comment style, Vitest tests next to the code).
- "Reference" links point to the decision that settles a detail; read it before starting the task.
- After each Go task: `go test ./<package>/...`. After each template task: `npm test` inside `templates/site/quarterly/` (run `npm install` there once after T006).

---

## Phase 1: Setup (shared infrastructure)

**Purpose**: dependencies and the empty template tree.

- [ ] T001 Promote `github.com/yuin/goldmark` and `golang.org/x/text` from `// indirect` to direct requirements and add `golang.org/x/image` in `go.mod`/`go.sum` (`go get golang.org/x/image@latest`, then `go mod tidy` once code imports them). Reference: research.md "Dependencies added"
- [ ] T002 Copy the whole `templates/site/log/` tree to `templates/site/quarterly/` (including `.github/workflows/deploy.yml`, but **excluding** `content/posts/`, `content-index.json`, `app/content/filters.ts`, `app/content/filters.test.ts`, `app/mcp/tools/log.ts`, `app/mcp/tools/log.test.ts`, `app/routes/home.tsx`, `app/routes/home.test.tsx`); keep `package.json` `name` as `site`; change the `<title>` in `templates/site/quarterly/app/root.tsx` to `biblio quarterly site` and the server name in `templates/site/quarterly/app/mcp/server.ts` to `biblio-quarterly-site` (temporarily remove the `registerLogTools` import/call there)
- [ ] T003 [P] Add `templates/site/quarterly/.gitignore` containing `node_modules/`, `dist/`, `.wrangler/`, `public/plates/`, `.biblio/`
- [ ] T004 [P] Add bundled fonts to `templates/site/quarterly/public/fonts/`: `Archivo[wdth,wght].woff2` (Latin + Latin-Extended subset), `IBMPlexMono-Regular.woff2`, `IBMPlexMono-Medium.woff2`, `OFL-Archivo.txt`, `OFL-IBMPlexMono.txt`; total under 350 KB (subset with `pyftsubset --flavor=woff2 --unicodes=U+0000-024F,U+2000-206F,U+2190-21FF,U+2600-26FF` if needed). Reference: research.md R14
- [ ] T005 [P] Create a Go test-fixture directory `testdata/quarterly/` (repo root) with: `site.yaml`; `colophon.md`; `issues/01-spring-2026/issue.md` (+ `cover.jpg`) with 2 articles in 2 categories; `issues/02-summer-2026/issue.md` with 3 articles (one `dense-notes.md` with ≥10 footnotes and margin notes incl. a long one, one with `image:` and an inline image under `img/`); `issues/03-autumn-2026/issue.md` with `draft: true` and 1 article; image files `color.jpg`, `alpha.png`, `photo.webp`, `corrupt.jpg` (truncated), `huge-header.png` (header claiming > 40 MP); keep binary files small (≤ 64×64 px except where the test needs size)
- [ ] T006 Run `npm install` in `templates/site/quarterly/` and commit the resulting `templates/site/quarterly/package-lock.json` (dependencies are unchanged from log, so this should be a no-op diff against the copied lock file)

**Checkpoint**: tree exists, Go module resolves.

---

## Phase 2: Foundational (blocking prerequisites)

**Purpose**: type registration, typed settings, the per-kind parser entry point, index options, and the pre-build step for **all** types. Every user story depends on this phase.

**⚠️ No user story work starts until this phase is complete.**

### Tests for Phase 2 (write first, confirm they fail)

- [ ] T007 [P] Add a test in `internal/scaffold/generator_test.go` asserting `quarterly` is in `SiteTypes`, that `Generate` with `SiteType: "quarterly", Theme: "quarterly-default"` writes `.gitignore` into the project, and that profile/log/docs scaffold the exact same file list as before (pin the current lists in the test)
- [ ] T008 [P] Add `internal/theme/typesetting_test.go`: a table test with one valid and at least one invalid value per key from contracts/config-and-tokens.md (accent, font_body/display/mono incl. `x;}` and quotes, base_size, scale 1.0/2.0, leading, measure, align `center`, grid, margin_notes non-bool, drop_caps, lead_in, numerals, rules, fleurons, grayscale_images), an unknown key (`colour`), and that each error message contains the key name and the accepted values; plus `Resolve` tests: no overrides = preset defaults, one override changes only that value
- [ ] T009 [P] Add tests in `internal/content/schema_test.go`: `ReadSiteConfig` parses a `typesetting:` block into the typed struct; a `typesetting:` block with `site_type: log` returns the error "typesetting is only supported for site_type quarterly"; existing tests stay unchanged and passing
- [ ] T010 [P] Add `internal/prebuild/prebuild_test.go`: for a scaffolded `log` site, `Run(dir, Build)` writes `content-index.json` equal to what `biblio-cli index` produced before, writes no `public/plates/`, leaves `public/tokens.css` byte-identical; an invalid `site.yaml` returns an error and writes nothing; `Report` lists the steps run
- [ ] T011 [P] Add tests in `internal/cli/build_test.go` and `internal/cli/dev_test.go` (stubbing `BuildCommand` and the dev supervisor as the existing tests do) asserting the pre-build step runs before the wrapped command, and that a validation error stops the command before the wrapped command starts
- [ ] T012 [P] Create `test/integration/prebuild-regression/prebuild_regression_test.go`: for each of profile, log, docs, scaffold, run the pre-build step, then assert every file under the project except `content-index.json` is byte-identical to a scaffold that did not run it, and that `content-index.json` equals `index.Generate` output (FR-036, SC-007)

### Implementation for Phase 2

- [ ] T013 Add `"quarterly"` to `SiteTypes` and change `Walk` to skip only files named `.gitkeep` (not every dotfile) in `internal/scaffold/assets.go` (makes T007 pass; the TUI in `internal/tui/new_flow.go` picks the type up automatically, so check it renders four choices)
- [ ] T014 Add the `quarterly-default` preset (`SiteType: "quarterly"`, name "Quarterly Default", tokens from research doc 004's appendix: `color-bg #fbfaf7`, `color-fg #111111`, `color-accent #d7261e`, `color-muted #6b6b6b`, `color-rule #111111`, `color-plate-bg #ececec`, `font-body`/`font-heading`/`font-display` `"Archivo", "Helvetica Neue", Arial, sans-serif`, `font-mono` `"IBM Plex Mono", ui-monospace, monospace`) in `internal/theme/presets.go`
- [ ] T015 Implement `Typesetting` (pointer fields per key, YAML tags matching the schema), `Validate() []error`, `Defaults()` for `quarterly-default`, and `Resolve(preset Preset, overrides Typesetting) (Effective, error)` in `internal/theme/typesetting.go` (makes T008 pass). Reference: research.md R9; contracts/config-and-tokens.md "Error format"
- [ ] T016 Add `Typesetting *theme.Typesetting \`yaml:"typesetting"\`` to `SiteConfig` and the non-quarterly rejection in `internal/content/schema.go` (use `yaml` `KnownFields` decoding only for the `typesetting` node so unknown keys inside it fail while the rest of `site.yaml` keeps its current lenient parsing) (makes T009 pass)
- [ ] T017 Add `Options{IncludeDrafts bool; Plates map[string]string}` and `GenerateWithOptions(contentDir string, opts Options) (Index, []error)` to `internal/index/generator.go`; keep `Generate(contentDir)` as `GenerateWithOptions(contentDir, Options{})` so all existing callers and tests are unchanged; give `Index` an optional quarterly section (`Issues`, `Categories`, `CurrentIssue`, `Typesetting`) that is omitted from JSON when empty so profile/log/docs JSON is unchanged
- [ ] T018 Create `internal/prebuild/prebuild.go`: `type Mode int` (`Build`, `Dev`), `type Report struct` (steps run, counts, warnings), `Run(projectDir string, mode Mode) (Report, error)` doing validate → tokens (quarterly only; stub that returns nil until US3) → plates (quarterly only; stub until US4) → index via `GenerateWithOptions` with `IncludeDrafts: mode == Dev`, writing `content-index.json` with the same `json.MarshalIndent(idx, "", "  ")` formatting as `internal/cli/index.go` (makes T010 pass). Reference: research.md R8
- [ ] T019 Call `prebuild.Run(cwd, prebuild.Build)` in `RunBuild` before `runWrapped` in `internal/cli/build.go`, and `prebuild.Run(cwd, prebuild.Dev)` in `runDev` before `devserver.NewSupervisor` in `internal/cli/dev.go`; print the report lines in the format of contracts/cli.md (makes T011 and T012 pass)
- [ ] T020 Run `go test ./...` and `scripts/check.sh` from the repo root; every pre-existing test must pass unchanged

**Checkpoint**: four types registered; `build`/`dev` validate and index every site automatically; profile/log/docs output unchanged.

---

## Phase 3: User Story 1 — Publish a numbered issue (Priority: P1) 🎯 MVP

**Goal**: scaffold a quarterly site, write issues and articles, build, deploy; contents/archive/category/article/colophon pages; drafts.

**Independent Test**: quickstart.md scenarios 1 and 5.

### Tests for User Story 1 (write first, confirm they fail)

- [ ] T021 [P] [US1] Create `internal/content/quarterly_test.go`: `ParseItemKind` picks issue/article/colophon by path; required fields per kind (issue: title, number, date; article: title, category; colophon: title) with errors naming file and field; `number` must be an integer ≥ 1; `date` must be `YYYY-MM-DD`; `cover` without `cover_alt` and `image` without `image_alt` fail; `draft` must be a bool; a `.md` outside `issues/<dir>/` (other than `colophon.md`) fails; `issues/a/b/x.md` fails
- [ ] T022 [P] [US1] Add `CategoryKey` tests to `internal/content/quarterly_test.go`: `"Essays"`, `" essays "`, `"ESSAYS"` → `essays`; `"Field  Notes"` → `field-notes`; `"Ukázky"` → `ukazky`; an all-symbol label returns an error
- [ ] T023 [P] [US1] Create `internal/index/quarterly_test.go` against `testdata/quarterly/`: derived `kind`, `issue`, `issue_number`, `issue_date`, `category_key`, `position`, `article_count`; `entries` order (issues by number desc, `issue.md` then articles by `order` then file name, colophon last); `issues` newest first; `categories` with first-seen labels; `current_issue` = `02-summer-2026` with drafts excluded and with `IncludeDrafts` (draft is never current); draft issue and its articles absent without `IncludeDrafts` and present with `draft: true` with it; duplicate `number` (add a temp copy) errors naming both folders; an issue with no articles has `article_count: 0`; `typesetting` object present with defaults
- [ ] T024 [P] [US1] Create `test/integration/pilot-site-quarterly/pilot_site_quarterly_test.go`: scaffold `--type=quarterly`, run `prebuild.Run(dir, Build)`, assert the template's sample content indexes with no errors, the index has ≥2 issues, ≥2 categories, a colophon, and that `content-index.json` shipped in the template equals the freshly generated one (keeps the shipped file honest)
- [ ] T025 [P] [US1] Create `templates/site/quarterly/app/content/load.test.ts`: accessors `listIssues()`, `getIssue(slugOrNumber)`, `currentIssue()`, `listCategories()`, `articlesInIssue(slug)` (grouped by category in contents order), `articlesInCategory(key)`, `getColophon()` against a fixture index
- [ ] T026 [P] [US1] Create `templates/site/quarterly/app/routes/contents.test.tsx`: renders masthead, "No. N · season" (falls back to the date), cover plate with alt when present, categories in order, each article with numeral (roman by default), title link, leader element and author; the empty-state message when there is no current issue; a link to `/issues`
- [ ] T027 [P] [US1] Create `templates/site/quarterly/app/routes/archive.test.tsx` (every published issue newest first with number, title, season/date, article count, link) and `templates/site/quarterly/app/routes/category.test.tsx` (only that category's articles, newest issue first, each showing its issue; unknown key → 404 response)
- [ ] T028 [P] [US1] Create `templates/site/quarterly/app/routes/detail.test.tsx` (replacing the copied log version): article page shows category label, title, author, dek, lead plate with alt and "Fig. 1 — caption", and body; colophon slug renders the colophon; unknown slug → 404; masthead links to `/colophon` only when it exists
- [ ] T029 [P] [US1] Update `templates/site/quarterly/app/negotiation/accept.test.ts` and add `templates/site/quarterly/app/negotiation/markdown.test.ts`: Markdown view of `/`, `/issues`, `/issues/<issue>` and `/category/<key>` is built from the index as a list; an article's Markdown view is its body (US2 later switches it to `markdown_body`)
- [ ] T030 [P] [US1] Update `templates/site/quarterly/app/discovery/llms-txt.test.ts` to the layout in contracts/site-surface.md "llms.txt" (current issue, issues, current-issue articles, colophon; no drafts)
- [ ] T031 [P] [US1] Create `templates/site/quarterly/react-router.config.test.ts` asserting `prerender()` returns `/`, `/issues`, each published issue, each category, each article and `/colophon`, and nothing for drafts

### Implementation for User Story 1

- [ ] T032 [US1] Implement `ParseItemKind(relPath string, data []byte, siteType string) (Item, error)` and per-kind validation in `internal/content/quarterly.go`, reusing `splitFrontmatter` and `SlugFromPath` from `schema.go`; for non-quarterly site types it delegates to `ParseItem` unchanged (makes T021 pass)
- [ ] T033 [US1] Implement `CategoryKey(label string) (string, error)` in `internal/content/quarterly.go` using `golang.org/x/text/unicode/norm` NFKD and dropping combining marks (makes T022 pass). Reference: research.md R10
- [ ] T034 [US1] Implement the quarterly path in `GenerateWithOptions` in `internal/index/generator.go` (or a new `internal/index/quarterly.go`): two-pass walk (issues first, then articles), derived fields, ordering rules, drafts, `issues`/`categories`/`current_issue`/`typesetting` (defaults until US3 wires overrides), cross-file checks; all errors collected (makes T023 pass). Reference: data-model.md "content-index.json shape", "Ordering rules", "Drafts"
- [ ] T035 [US1] Switch `internal/cli/index.go` to call `GenerateWithOptions(contentDir, index.Options{})` and report `N item(s), M draft issue(s) excluded` for quarterly. Reference: contracts/cli.md "biblio-cli index"
- [ ] T036 [P] [US1] Write the sample content in `templates/site/quarterly/content/`: `issues/01-spring-2026/` (issue.md with `season: Spring 2026`, cover image + alt, 2 articles in "Essays" and "Design") and `issues/02-autumn-2026/` (issue.md, 3 articles across 2 categories, one with `image:`, one `on-margins.md` with dense footnotes and margin notes), plus `colophon.md` naming the typefaces and licenses (FR-002); keep images small grayscale-friendly JPEGs (≤ 200 KB total)
- [ ] T037 [P] [US1] Implement accessors in `templates/site/quarterly/app/content/load.ts` reading the index's `entries`/`issues`/`categories`/`current_issue`/`typesetting` (no path parsing in TS) (makes T025 pass)
- [ ] T038 [US1] Update `templates/site/quarterly/app/routes.ts` to `index("routes/contents.tsx")`, `route("issues", "routes/archive.tsx")`, `route("issues/:issue", "routes/contents.tsx", { id: "issue-contents" })`, `route("category/:category", "routes/category.tsx")`, `route("*", "routes/detail.tsx")`
- [ ] T039 [P] [US1] Implement `templates/site/quarterly/app/routes/contents.tsx` (current issue for `/`, `params.issue` otherwise; numbering via a CSS counter class for the page and `typesetting.numerals` for text) (makes T026 pass)
- [ ] T040 [P] [US1] Implement `templates/site/quarterly/app/routes/archive.tsx` and `templates/site/quarterly/app/routes/category.tsx` (makes T027 pass)
- [ ] T041 [US1] Rewrite `templates/site/quarterly/app/routes/detail.tsx` for article and colophon entries, and add the shared masthead component in `templates/site/quarterly/app/content/masthead.tsx` (makes T028 pass)
- [ ] T042 [US1] Update `templates/site/quarterly/app/negotiation/markdown.ts` (index-built Markdown for list pages) and keep `accept.ts`/`ansi.ts` behavior from log (makes T029 pass)
- [ ] T043 [P] [US1] Update `templates/site/quarterly/app/discovery/llms-txt.ts` (makes T030 pass)
- [ ] T044 [US1] Update `prerender()` in `templates/site/quarterly/react-router.config.ts` to build the list from `issues`, `categories` and `entries` (makes T031 pass). Reference: research.md R12
- [ ] T045 [US1] Generate and commit `templates/site/quarterly/content-index.json` by running `go run ./cmd/biblio-cli index templates/site/quarterly/content` from the repo root (makes T024 pass)
- [ ] T046 [US1] Run `npm test` in `templates/site/quarterly/` and `go test ./...`; then run quickstart.md scenarios 1 and 5 by hand

**Checkpoint**: a quarterly site scaffolds, builds, prerenders every page, handles drafts, and negotiates Markdown/ANSI. MVP.

---

## Phase 4: User Story 2 — Scholarly notes (Priority: P2)

**Goal**: `[^n]` footnotes and `{>> … <<}` margin notes: sidenotes on wide screens, foot list and † toggles on narrow screens, all text kept in Markdown/ANSI views, structured in the index.

**Independent Test**: quickstart.md scenario 2.

### Tests for User Story 2 (write first, confirm they fail)

- [ ] T047 [P] [US2] Create `internal/content/notes_test.go` with golden files under `internal/content/testdata/notes/`: footnotes extracted in reference order with `id`, `label`, `number`, inline-Markdown `text`, and `ref_offset` in **code points** (include a body with Czech diacritics and an emoji before a reference); margin notes `m1…` with text and offset; a margin note containing emphasis and a link; an unclosed `{>>` returns a warning and is left as text; a reference with no definition and an unused definition each return a warning; `{>>` inside a code span or fenced code block is ignored
- [ ] T048 [P] [US2] Add `MarkdownBody` tests to `internal/content/notes_test.go`: margin notes become `[^mN]` references with definitions appended; existing footnotes unchanged; when the author already uses a label `m1`, the generated labels skip to unused ones; no note text is lost (count check)
- [ ] T049 [P] [US2] Extend `internal/index/quarterly_test.go`: articles carry `footnotes`, `margin_notes` and `markdown_body`; issues and colophon do not carry `markdown_body`; note warnings are surfaced in the generator's warnings list (not errors)
- [ ] T050 [P] [US2] Create `templates/site/quarterly/app/content/remark-margin-notes.test.ts`: the plugin turns `{>> text <<}` in text nodes into `marginNote` mdast nodes, handles inline Markdown inside, leaves code untouched, and leaves an unclosed marker as text
- [ ] T051 [P] [US2] Create `templates/site/quarterly/app/content/notes.test.tsx`: with `margin_notes: true` each footnote reference renders an inline sidenote span plus the foot-list item, and each margin note renders label "†" + checkbox + `span.marginnote`; with `margin_notes: false` no sidenote spans render and margin notes appear in the foot list marked †; raw HTML in a note is not rendered
- [ ] T052 [P] [US2] Create `templates/site/quarterly/app/content/notes-agree.test.tsx`: for every article entry in the template's `content-index.json`, render the body and assert the number, order and text of rendered footnotes and margin notes equal the entry's `footnotes`/`margin_notes` (the cross-parser contract from research.md R2)
- [ ] T053 [P] [US2] Extend `templates/site/quarterly/app/negotiation/markdown.test.ts` and `templates/site/quarterly/app/negotiation/ansi.test.ts`: an article's Markdown and ANSI views use `markdown_body` and contain every note's text

### Implementation for User Story 2

- [ ] T054 [US2] Implement `ExtractNotes(body string) (footnotes, marginNotes []Note, warnings []string)` and `MarkdownBody(body string, marginNotes []Note) string` in `internal/content/notes.go` using goldmark with `extension.Footnote` for footnotes and a scanner over text segments (skipping code) for margin notes (makes T047, T048 pass). Reference: research.md R2, R6; data-model.md "Note"
- [ ] T055 [US2] Call `ExtractNotes`/`MarkdownBody` for article entries in the quarterly index path and surface warnings through `Report` in `internal/index/generator.go` (or `internal/index/quarterly.go`) and `internal/prebuild/prebuild.go` (makes T049 pass)
- [ ] T056 [P] [US2] Implement `templates/site/quarterly/app/content/remark-margin-notes.ts` (makes T050 pass)
- [ ] T057 [US2] Implement note rendering in `templates/site/quarterly/app/content/notes.tsx` and wire it plus the remark plugin into `templates/site/quarterly/app/content/render.tsx` via react-markdown `components`, reading `typesetting.margin_notes` from the index (makes T051, T052 pass). Reference: research.md R7
- [ ] T058 [US2] Add the sidenote layout to `templates/site/quarterly/public/quarterly.css` (create the file with only this section for now): margin column at `min-width: 64rem`, sidenotes floated right with `clear: right`, foot list hidden on wide / sidenotes hidden on narrow, † toggle CSS (`:checked + .marginnote`), and the `margin_notes: false` layout; link it from `templates/site/quarterly/app/root.tsx` after `/tokens.css`
- [ ] T059 [US2] Make `templates/site/quarterly/app/negotiation/markdown.ts` serve `markdown_body` for articles when present (makes T053 pass)
- [ ] T060 [US2] Regenerate `templates/site/quarterly/content-index.json` (as in T045), run `npm test` and `go test ./...`, then check quickstart.md scenario 2 at 320, 375, 768, 1024, 1440 and 1920 px with JavaScript disabled (SC-005)

**Checkpoint**: notes render correctly in every view and are structured in the index.

---

## Phase 5: User Story 3 — Configure the typesetting (Priority: P3)

**Goal**: the default look, and every `typesetting:` key changing the site on the next `dev`/`build`, with strict validation.

**Independent Test**: quickstart.md scenario 3.

### Tests for User Story 3 (write first, confirm they fail)

- [ ] T061 [P] [US3] Add golden tests to `internal/theme/theme_test.go` with files under `internal/theme/testdata/`: `GenerateTokensCSS` for `quarterly-default` with no overrides, and with `accent`, `align: justify`, `numerals: arabic`, `grid: centered`, `rules: none`, `drop_caps: raised`, `fleurons: true` — every property in contracts/config-and-tokens.md "tokens.css" present, `--ts-step-1…5` computed from `base_size × scale^n` rounded to 3 decimals, header comment listing overridden keys, output deterministic; profile/log/docs output unchanged byte for byte
- [ ] T062 [P] [US3] Add tests to `internal/prebuild/prebuild_test.go`: for quarterly, `Run` rewrites `public/tokens.css` from `site.yaml` overrides, writes `typesetting.margin_notes`/`numerals` into the index, warns (does not fail) when `font_body` names a family with no matching file in `public/fonts/`, and finds the bundled Archivo file case-insensitively
- [ ] T063 [P] [US3] Create `templates/site/quarterly/app/routes/theme.test.tsx` modeled on `templates/site/profile/app/routes/theme.test.tsx`: no inline `style` props in rendered routes; `public/quarterly.css` contains no color literal (`#hex`, `rgb(`, `hsl(`, CSS named colors) and no font family name outside `@font-face` blocks; `root.tsx` links `/tokens.css` before `/quarterly.css`
- [ ] T064 [P] [US3] Add to `templates/site/quarterly/app/routes/contents.test.tsx`: with `typesetting.numerals: "arabic"` in the fixture index, text output (Markdown view helper) numbers 1, 2, 3

### Implementation for User Story 3

- [ ] T065 [US3] Extend `GenerateTokensCSS` in `internal/theme/tokens.go` to take the resolved `Effective` values for quarterly (keep the existing `Preset`-only behavior for the other presets, e.g. `GenerateTokensCSS(preset)` unchanged plus `GenerateQuarterlyTokensCSS(eff Effective)`), emitting the properties in contracts/config-and-tokens.md (makes T061 pass)
- [ ] T066 [US3] Replace the tokens stub in `internal/prebuild/prebuild.go` with: resolve → write `public/tokens.css`; pass `margin_notes`/`numerals` to the index; implement the font-file warning (makes T062 pass). Reference: research.md R9
- [ ] T067 [US3] Make `internal/scaffold/generator.go` write the quarterly tokens through the same `Resolve` + `GenerateQuarterlyTokensCSS` path at `new` time so a fresh project's `public/tokens.css` equals what the first `build` writes
- [ ] T068 [US3] Complete `templates/site/quarterly/public/quarterly.css`: `@font-face` for the bundled faces (`font-display: swap`), base type (`--ts-base-size`, `--ts-leading`, `--ts-measure`, `--ts-align`, `--ts-hyphens`), the type scale, the asymmetric/centered grid via `--ts-grid-columns`, rules, drop caps, small-caps lead-in, contents numerals via `counter(article, var(--ts-numerals))`, dotted leaders, fleurons, plates on `--color-plate-bg`, and `@media print` basics — values only through `var(--…)` (makes T063 pass)
- [ ] T069 [US3] Use `typesetting.numerals` for the numerals in text outputs (Markdown view of contents pages, `llms.txt` if numbered) in `templates/site/quarterly/app/content/numerals.ts` (shared `toRoman`), and use it from `templates/site/quarterly/app/negotiation/markdown.ts` (makes T064 pass)
- [ ] T070 [US3] Run `npm test` in `templates/site/quarterly/` and `go test ./...`, then quickstart.md scenario 3 including every invalid value

**Checkpoint**: default look in place; every typesetting key works and is validated.

---

## Phase 6: User Story 4 — Black-and-white plates (Priority: P4)

**Goal**: every content image served grayscale, width-capped, cached, and safely decoded.

**Independent Test**: quickstart.md scenario 4.

### Tests for User Story 4 (write first, confirm they fail)

- [ ] T071 [P] [US4] Create `internal/plates/plates_test.go` against `testdata/quarterly/`: JPEG → grayscale JPEG (single-channel); PNG with alpha keeps alpha and has equal RGB channels; WebP with alpha → PNG, without → JPEG; width capped at 1600 with aspect kept, narrower images not upscaled; `grayscale: false` keeps color; `corrupt.jpg` errors naming the file; `huge-header.png` errors before full decode; a symlink pointing outside `content/` is rejected; the output name is `<16-hex>.<ext>` and stable for the same bytes + settings, different when settings change
- [ ] T072 [P] [US4] Add cache tests to `internal/plates/plates_test.go`: a second `Convert` run converts 0 and skips all; changing one source converts 1; a plate whose source was removed is deleted (orphan cleanup); the returned map is `content-relative source path → /plates/<name>`
- [ ] T073 [P] [US4] Create `internal/plates/fuzz_test.go`: `FuzzDecode` feeding arbitrary bytes to the decode path; it must return an error or a bounded image, never panic or allocate beyond the cap
- [ ] T074 [P] [US4] Extend `internal/index/quarterly_test.go`: with a `Plates` map, `plate` and `cover_plate` are set; inline image destinations in `body` and `markdown_body` that resolve to a plate are rewritten to `/plates/…`; external URLs and unresolved paths are left as they are; destinations inside code blocks are not rewritten
- [ ] T075 [P] [US4] Add to `test/integration/pilot-site-quarterly/pilot_site_quarterly_test.go`: after `prebuild.Run(dir, Build)` on the scaffolded site, every image referenced from content has a plate in `public/plates/`, and no `src` in the index points into `content/` (SC-004)

### Implementation for User Story 4

- [ ] T076 [US4] Implement `internal/plates/plates.go`: `type Settings{Grayscale bool; MaxWidth int; MaxPixels int}`, `Convert(contentDir, outDir string, s Settings) (map[string]string, Stats, error)` with content sniffing, `image.DecodeConfig` cap check, symlink resolution confined to `contentDir`, `x/image/draw` CatmullRom resize, grayscale conversion, JPEG q82/PNG encoding, hash naming, skip-if-exists and orphan cleanup (makes T071, T072, T073 pass). Reference: research.md R4; data-model.md "Plate"
- [ ] T077 [US4] Rewrite image destinations from the goldmark AST and set `plate`/`cover_plate` in the quarterly index path in `internal/index/quarterly.go` (makes T074 pass). Reference: research.md R5
- [ ] T078 [US4] Replace the plates stub in `internal/prebuild/prebuild.go`: run `plates.Convert` with `grayscale_images` from the effective typesetting, pass the map to `GenerateWithOptions`, report converted/skipped/removed counts (makes T075 pass)
- [ ] T079 [P] [US4] Use `plate`/`cover_plate` for every `<img>` in `templates/site/quarterly/app/routes/contents.tsx` and `templates/site/quarterly/app/routes/detail.tsx` (inline images already arrive rewritten in `body`); add an assertion to `templates/site/quarterly/app/routes/detail.test.tsx` that no `<img src>` starts with `content/` or is relative
- [ ] T080 [US4] Run `go test ./...` (incl. `go test -fuzz=FuzzDecode -fuzztime=30s ./internal/plates/`) and quickstart.md scenario 4

**Checkpoint**: all content images are plates; reruns are cheap; bad images fail safely.

---

## Phase 7: User Story 5 — Agents read the magazine (Priority: P5)

**Goal**: `list_issues`, `get_issue`, `list_articles`, `get_article` per contracts/mcp-tools.md.

**Independent Test**: quickstart.md scenario 6.

### Tests for User Story 5 (write first, confirm they fail)

- [ ] T081 [P] [US5] Create `templates/site/quarterly/app/mcp/tools/quarterly.test.ts` against a fixture index: each tool's text output matches contracts/mcp-tools.md (table header, ordering, numerals, line formats); `get_issue` by slug and by number; `list_articles` with each filter alone and combined (AND), category matched by key, author case-insensitive, date range inclusive on issue date, invalid date returns a message naming the field; `get_article` text includes every note and `structuredContent` has `footnotes`/`marginNotes` with correct offsets; unknown issue/article returns the not-found text with `isError` false; drafts never visible
- [ ] T082 [P] [US5] Update `templates/site/quarterly/app/mcp/server.test.ts`: `tools/list` returns the three shared tools plus the four quarterly tools, `get_article` advertises an `outputSchema`, and a log-only tool name (`list_entries`) is unknown

### Implementation for User Story 5

- [ ] T083 [P] [US5] Implement article filtering in `templates/site/quarterly/app/content/filters.ts` (issue, category key, author, tag, dateFrom/dateTo) with `templates/site/quarterly/app/content/filters.test.ts` covering the same cases as T081's filter rows (write the test first)
- [ ] T084 [US5] Implement `templates/site/quarterly/app/mcp/tools/quarterly.ts` (`listIssues`, `getIssue`, `listArticles`, `getArticle`, `registerQuarterlyTools`) with zod input schemas and the `outputSchema` for `get_article`, reusing `filters.ts` and `numerals.ts` (makes T081 pass)
- [ ] T085 [US5] Register the tools in `templates/site/quarterly/app/mcp/server.ts` and update `templates/site/quarterly/app/mcp/query.ts` if the shared tools need the new index shape (makes T082 pass)
- [ ] T086 [US5] Run `npm test` in `templates/site/quarterly/`, then quickstart.md scenario 6 with MCP Inspector against `biblio-cli dev`

**Checkpoint**: the agent surface is at parity with the other types.

---

## Phase 8: User Story 6 — Keep it in step with the other types (Priority: P6)

**Goal**: `update`/`upgrade`, cross-type review with four types, no regressions.

**Independent Test**: quickstart.md scenario 7.

### Tests for User Story 6 (write first, confirm they fail)

- [ ] T087 [P] [US6] Extend `test/integration/cross-type-review/theming_review_test.go` and `test/integration/cross-type-review/terminal_rendering_review_test.go` to iterate over all four `scaffold.SiteTypes`, with any quarterly-specific expectations (extra tokens) asserted explicitly
- [ ] T088 [P] [US6] Create `test/integration/update-three-sites/quarterly_site_test.go` following `project_log_test.go`: scaffold quarterly, author an issue and change `site.yaml` typesetting, simulate a template change, run update/upgrade, assert template files advance and `content/` plus `content/site.yaml` are byte-identical; generated files (`public/tokens.css`, `public/plates/`, `content-index.json`) are not reported as conflicts
- [ ] T089 [P] [US6] Add a TUI test in `internal/tui/` (new `new_flow_test.go` if none exists) asserting the type step lists profile, log, docs, quarterly and the theme step for quarterly offers `quarterly-default`

### Implementation for User Story 6

- [ ] T090 [US6] Fix whatever T087–T089 expose in `internal/migrate/` (e.g. excluding `public/plates/` and `.biblio/` from classification in `internal/migrate/classify.go`) and `internal/scaffold/manifest.go`
- [ ] T091 [US6] Update `test/integration/cross-type-review/CONSISTENCY-RECORD.md` with the four-type review results
- [ ] T092 [US6] Run `go test ./...`, `./scripts/check.sh`, `npm test` in all four templates, and quickstart.md scenario 7

**Checkpoint**: four types consistent; nothing regressed.

---

## Phase 9: Polish & cross-cutting concerns

- [ ] T093 [P] Append decisions R1–R4 (text face, notes in Go, toggles split, image output) to the Decision Log in `specs/research/004-quarterly-site-type.md` and tick its matching Open Questions
- [ ] T094 [P] Check off the four §0.5 items and the completed §11 tasks in `specs/research/roadmap.md`, and update `specs/features/011-quarterly-site-type.md` scope checkboxes and status
- [ ] T095 [P] Update `AGENTS.md` and `CLAUDE.md`: four site types, `templates/site/quarterly`, the new `internal/plates` and `internal/prebuild` packages, and that `build`/`dev` now index automatically
- [ ] T096 [P] Update `specs/research/architecture.md` §2 to drop "(planned)" from step 0 and §4's quarterly note once shipped
- [ ] T097 [P] Add a quarterly section to the user guide draft for the release (`docs/releases/v0.08/README.md`): create, configure (typesetting keys, drafts), update, deploy, destroy — the constitution's usage-guide minimum — and an entry in `docs/releases/CHANGELOG.md`
- [ ] T098 Security pass: grep the template for `dangerouslySetInnerHTML` and raw HTML in react-markdown (`rehype-raw` must not be present); confirm every typesetting value in `tokens.css` comes from a validated field; confirm `public/plates/` and `.biblio/` are gitignored in the scaffolded project
- [ ] T099 Run the full quickstart.md (scenarios 1–7) on a freshly scaffolded site and record the results in `test/integration/pilot-site-quarterly/QUICKSTART-RECORD.md`

---

## Dependencies & execution order

### Phase dependencies

- **Setup (Phase 1)**: none.
- **Foundational (Phase 2)**: needs Setup. **Blocks every story.**
- **US1 (Phase 3)**: needs Phase 2. It is the MVP and the base the other stories render into.
- **US2 (Phase 4)**: needs US1's index path (T034) and article page (T041).
- **US3 (Phase 5)**: needs Phase 2 (T015, T018). Its CSS work (T068) extends the file T058 creates, so do US2's T058 first or create the file in whichever comes first.
- **US4 (Phase 6)**: needs US1's index path (T034); independent of US2/US3 apart from T079 touching routes from US1.
- **US5 (Phase 7)**: needs US1's index; `get_article` notes need US2 (T055).
- **US6 (Phase 8)**: needs all of the above to review four complete types.
- **Polish (Phase 9)**: after the stories it documents.

### Story completion order

```text
Setup → Foundational → US1 ─┬─▶ US2 ─┬─▶ US5 ─┐
                            ├─▶ US3  │        ├─▶ US6 → Polish
                            └─▶ US4 ─┘────────┘
```

### Within each story

1. Tests first; run them and confirm they fail for the stated reason.
2. Go side (`content` → `index` → `prebuild`) before the template side that reads its output.
3. Regenerate the template's `content-index.json` whenever the index shape changes.
4. Story checkpoint before moving on.

## Parallel opportunities

- **Phase 1**: T003, T004, T005 together.
- **Phase 2 tests**: T007–T012 all touch different files.
- **US1 tests**: T021–T031 (T021/T022 share a file, so write them together or in sequence).
- **US1 implementation**: T036, T037, T039, T040, T043 in parallel once T034 lands.
- **After US1**: US3 and US4 can be built in parallel; US2 in parallel with US4.
- **US6 tests**: T087, T088, T089 together.
- **Polish**: T093–T097 together.

### Example: US1 tests in parallel

```text
Task: "Create internal/index/quarterly_test.go …"                      (T023)
Task: "Create test/integration/pilot-site-quarterly/… test"            (T024)
Task: "Create templates/site/quarterly/app/content/load.test.ts"       (T025)
Task: "Create templates/site/quarterly/app/routes/contents.test.tsx"   (T026)
Task: "Create templates/site/quarterly/app/routes/archive.test.tsx …"  (T027)
```

### Example: US4 tests in parallel

```text
Task: "Create internal/plates/plates_test.go …"  (T071)
Task: "Create internal/plates/fuzz_test.go"      (T073)
Task: "Extend internal/index/quarterly_test.go"  (T074)
```

## Implementation strategy

### MVP first (US1 only)

1. Phase 1 → Phase 2 → Phase 3.
2. **Stop and validate** with quickstart.md scenarios 1 and 5: a publishable magazine with default styling (no sidenotes yet, color images, tools are the shared three).
3. It can ship to a real domain at this point if wanted.

### Incremental delivery

1. US1 → publish.
2. US2 (notes) → the defining look.
3. US3 (typesetting) and US4 (plates), in either order.
4. US5 (agent tools) → parity.
5. US6 → release gate for v0.08.
6. Polish → release packaging (the `specs/releases/0.08-…/` directory itself is assembled when v0.08 closes, per the constitution).

## Notes

- `[P]` = different files, no unfinished dependency.
- Every story phase ends with a checkpoint that runs its quickstart scenario.
- Commit after each task or logical group; never commit with a failing test that isn't the red step of the current task.
- The shipped `templates/site/quarterly/content-index.json` must always equal freshly generated output (T024 enforces it).
