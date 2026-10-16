# Phase 0 Research: `quarterly` site type

**Feature**: [spec.md](spec.md) · **Plan**: [plan.md](plan.md) · **Date**: 2026-10-15

This file settles every open item the plan needed before design: the four
planning questions in [roadmap §0.5](../research/roadmap.md#05-quarterly-site-type)
and the technical unknowns found while reading the code. The upstream
research and its Decision Log are in
[004-quarterly-site-type.md](../research/004-quarterly-site-type.md); this
file does not repeat decisions already recorded there.

Each entry gives **Decision / Rationale / Alternatives considered**.

---

## R1. Default text face (roadmap §0.5)

- **Decision:** **Archivo** for `font_body` and `font_display`; **IBM Plex
  Mono** for `font_mono`. Both are bundled as self-hosted woff2 (for the site)
  with Latin and Latin-Extended subsets, plus the OFL license texts.
- **Rationale:**
  - It is what the research doc's preset sketch already uses.
  - The profile and docs presets both use Inter (`internal/theme/presets.go`).
    Inter Tight would read as a sibling of those. Archivo's grotesk with
    slightly squared curves is closer to the 1970s Central European
    catalogue look the interview asked for.
  - Archivo ships a variable font (weight and width axes), so one file
    covers text, display and small-caps-like lead-ins.
  - Latin-Extended covers Czech and Polish diacritics.
- **Alternatives considered:** Inter Tight (too close to the existing
  presets); system grotesk stacks only (no control over the look, and
  `font_*` overrides would have nothing to fall back from).
- **Follow-up:** the side-by-side sample setting from the research doc is
  not a blocker. Swapping the face later is a one-line preset change plus
  font files, and no requirement depends on it.

## R2. Structured notes: where they are parsed (roadmap §0.5)

- **Decision:** **Go extracts footnotes and margin notes at index time.**
  `internal/content` gains a notes parser built on `goldmark` (already in the
  module graph through glamour; promoted to a direct dependency) with its
  footnote extension, plus a small scanner for `{>> … <<}`. The index entry
  for each article carries `footnotes[]` and `margin_notes[]`
  ([data-model.md](data-model.md#note)).
- **Rationale:**
  - Principle II: `biblio-cli` is the sole producer of the content index,
    and the follow-on PDF feature (005) needs the same parsed notes.
  - `get_article` becomes a pure lookup with no Markdown parsing in the
    Worker.
- **Consequence (accepted):** the HTML renderer in the template still parses
  Markdown with `remark-gfm` plus a small template-local remark plugin for
  margin notes, because rendering always happened in TS (react-markdown).
  That is rendering, not schema parsing, which is the same split the other
  three types already have. To stop the two parsers drifting, a template
  contract test renders every fixture article and asserts that the number,
  order and text of rendered notes equal the index's `footnotes`/`margin_notes`.
- **Reference positions:** `ref_offset` is the 0-based **Unicode code point**
  offset of the note's reference marker in the entry's `body`. Code points,
  not bytes, so TS (`Array.from(body)`) and Go (`utf8`) agree.
- **Alternatives considered:** parsing in TS inside the `get_article` tool
  (rejected: duplicates parsing in a second language and leaves 005 with
  nothing to reuse); Go rewriting the body into pre-rendered HTML (rejected:
  breaks the Markdown and ANSI views, which serve the body).

## R3. Ornament and structural toggles (roadmap §0.5)

- **Decision:** `tokens.css` **stays variables only** (one `:root` block, as
  today). Typesetting values split into two groups, both produced by one Go
  resolver from the same merged settings:
  1. **Visual values CSS can consume through `var()`** become custom
     properties in `tokens.css`: colors, fonts, `base_size`, `scale`
     (pre-computed into `--ts-step-1` … `--ts-step-5`), `leading`,
     `measure`, `align` (`--ts-align`, `--ts-hyphens`), `grid`
     (`--ts-grid-columns`), `rules` (`--ts-rule-width`), `drop_caps`
     (`--ts-dropcap-lines`, `--ts-dropcap-float`), `lead_in`
     (`--ts-lead-in-variant`), `fleurons` (`--ts-fleuron`, a fixed string or
     `none`), and `numerals` (`--ts-numerals: upper-roman | decimal`, used by
     CSS counters).
  2. **Structural switches that change markup** go into the content index
     under a top-level `typesetting` object that components read:
     `margin_notes` (sidenote markup vs. end-of-article notes) and
     `numerals` (also used by text views and MCP output, which have no CSS).
- **Rationale:** CSS cannot branch on a custom property's value without
  container style queries, which are not supported in every target browser.
  Emitting `[data-…]` rule blocks would turn `tokens.css` into a second
  stylesheet with logic in it. Passing the few structural switches through
  the index keeps one owner (Go) and lets React choose markup at prerender
  time, with no client JS.
- **Alternatives considered:** `@container style()` queries (browser support
  gap); generated attribute rules in `tokens.css` (mixes logic into a
  variables file); duplicating settings in a TS config (violates
  Principle II).

## R4. Image output (roadmap §0.5)

- **Decision:**
  - Maximum width **1600 px** (images narrower are not upscaled).
  - Pixel-count safety cap **40 MP**, checked with `image.DecodeConfig`
    before full decode.
  - JPEG stays JPEG (quality 82); PNG stays PNG; WebP is written as **PNG
    if it has an alpha channel, otherwise JPEG** (Go has a WebP decoder in
    `golang.org/x/image/webp`, but no encoder).
  - Grayscale opaque images are encoded as single-channel (`image.Gray`);
    images with alpha keep their alpha (`image.NRGBA` with equal RGB
    channels from Rec. 709 luminance).
  - Output is written to `public/plates/<hash>.<ext>`, where `<hash>` is the
    first 16 hex chars of SHA-256 over the source bytes plus the conversion
    settings (grayscale on/off, max width, encoder version). The hash in the
    name makes the cache check a file-exists test and busts browser caches.
  - Plates no longer referenced by any source are deleted at the end of the
    step.
- **Rationale:** keeps cross-compilation pure Go (decided in 004), keeps
  files small, and the hash-named file is its own cache entry, so no
  separate cache file is needed.
- **Alternatives considered:** always JPEG (loses transparency); keeping
  WebP (no encoder in pure Go); a 2000 px width (larger than any layout
  column needs at 2x on the widest grid).

## R5. How content refers to plates

- **Decision:** Go rewrites image destinations in the indexed `body` to
  plate URLs. The index also carries `plate` (article lead image) and
  `cover_plate` (issue cover). Only destinations that resolve to a converted
  image under `content/` are rewritten; other URLs are left unchanged. The
  destinations come from the goldmark AST, not a regular expression.
- **Rationale:** the HTML, Markdown and ANSI views and `get_article` all
  serve `body`; rewriting once in Go means none of them needs path logic,
  and relative paths like `img/fig1.jpg` never leak into a served view.
- **Alternatives considered:** mirroring plates at their source path so
  relative URLs resolve by themselves (breaks for WebP→JPEG/PNG, and loses
  hash-based cache busting); rewriting in TS (path logic in two languages).

## R6. Markdown and terminal views of notes

- **Decision:** the index carries a second body for quarterly articles,
  `markdown_body`, produced by Go: margin notes are rewritten as footnotes
  with labels `m1`, `m2`, … (renumbered to avoid any label already used in
  the article), and plate URLs are already rewritten. The template's
  Markdown and ANSI negotiation for quarterly articles serves
  `markdown_body` when present, else `body`.
- **Rationale:** the rewrite is a parsing task, so it belongs in Go
  (Principle II). The size cost (two bodies per article) is fine at
  personal scale and only applies to quarterly sites.
- **Alternatives considered:** rewriting `{>> … <<}` in TS at request time
  (parsing in a second language); storing only the rewritten body (the HTML
  renderer would then be unable to tell margin notes from footnotes).

## R7. Notes in HTML without client JS

- **Decision:**
  - **Margin notes:** the Tufte CSS pattern: an inline `<label>` with the
    dagger (†), a hidden checkbox, and an inline `<span class="marginnote">`.
    On wide viewports CSS floats the span into the margin column and hides
    the dagger toggle. On narrow viewports the span is hidden until the
    checkbox is checked. Everything stays phrasing content, so it is valid
    inside `<p>`.
  - **Footnotes:** rendered twice from the same mdast node: as an inline
    sidenote span (shown on wide viewports) and in the `remark-gfm` footnote
    section at the article foot (shown on narrow viewports). The hidden copy
    is `display: none`, which also removes it from the accessibility tree,
    so assistive tech reads each note once.
  - With `margin_notes: false`, the sidenote spans are not rendered, so both
    kinds of note appear at the foot and the text column widens.
  - Breakpoint: the margin column appears at `min-width: 64rem`.
- **Rationale:** the research doc suggested `<details>`, but `<details>` is
  flow content and is invalid inside a paragraph, so browsers would split
  the paragraph. The checkbox pattern is the established CSS-only
  alternative and keeps the prerendered pages free of client JS.
- **Alternatives considered:** `<details>` (invalid nesting); a JS
  positioning pass (adds client JS to prerendered pages); notes always at
  the foot (loses the look the interview asked for).
- **Overlap:** sidenotes are floated with `clear: right` into the margin
  column, so a long note pushes the next one down instead of overlapping
  it. SC-005 checks this at widths from 320 to 1920 px against a
  dense-notes fixture.

## R8. The pre-build step

- **Decision:** a new Go package `internal/prebuild` exposes
  `Run(projectDir string, mode Mode) (Report, error)`, where `Mode` is
  `Build` or `Dev`. It runs, in order:
  1. **Validate settings**, all types: parse `content/site.yaml` (now with
     typed `typesetting:`), check `site_type` is known, and for quarterly
     validate every typesetting key. A `typesetting:` block on a
     non-quarterly site is an error naming the site type, so a setting is
     never silently ignored.
  2. **Tokens**, quarterly only: resolve preset plus overrides and write
     `public/tokens.css`.
  3. **Plates**, quarterly only: convert images (R4).
  4. **Index**, all types: `index.Generate` with `IncludeDrafts: mode == Dev`
     and the plate map, written to `content-index.json`.

  `biblio-cli build` calls it before `vite build`. `biblio-cli dev` calls it
  once before starting the supervisor. Rerunning on every content change is
  deferred (the spec accepts restart-only, FR-024).
  `biblio-cli index` keeps working on its own with publish semantics
  (drafts excluded).
- **Rationale:** one place for the ordering, testable without Node, and
  `build.go`/`dev.go` stay thin callers. Per the clarification, profile, log
  and docs run only steps 1 and 4, so their pages and assets are unchanged.
- **Alternatives considered:** putting the steps inline in `build.go` and
  `dev.go` (duplicated ordering); a file watcher in `dev` now (more moving
  parts than this release needs).

## R9. Settings validation rules

- **Decision:** a typed `theme.Typesetting` struct with pointer fields, so
  "unset" and "set" are distinct. Rules:

  | Kind | Rule |
  |---|---|
  | color (`accent`) | `^#[0-9a-fA-F]{6}$` |
  | font family | `^[A-Za-z0-9][A-Za-z0-9 \-]{0,63}$` (no quotes, commas, semicolons, braces) |
  | length (`base_size`) | `^\d+(\.\d+)?(rem|px)$`, between 0.75rem/12px and 1.5rem/24px |
  | length (`measure`) | `^\d+(\.\d+)?(em|ch)$`, between 20 and 90 |
  | `scale` | number in [1.1, 1.618] |
  | `leading` | number in [1.1, 2.0] |
  | enums | exact match on the listed values |
  | bools | YAML booleans only |

  An unknown key inside `typesetting:` is an error that lists the valid
  keys. Every error message names the key and the accepted values (FR-023).
  Font-file presence: for each configured family, the step looks under
  `public/fonts/` for a file whose name, lowercased with spaces removed,
  starts with the family name in the same form. If none is found, it warns
  (FR-025).
- **Rationale:** no author string reaches CSS without passing a closed
  pattern, which closes the CSS-injection risk named in 004's security
  section.
- **Alternatives considered:** a full CSS value parser (overkill); allowing
  any string for font names (injection risk).

## R10. Content model in Go

- **Decision:**
  - `content.ParseItem` gains a file-kind-aware variant for quarterly:
    `ParseItemKind(relPath, data, siteType)` picks a schema from the path:
    `issues/<dir>/issue.md` → issue, `issues/<dir>/*.md` → article,
    `colophon.md` → colophon, anything else under `content/` → error. The
    existing `ParseItem` signature and per-type behavior for profile, log
    and docs are unchanged (their tests must pass untouched).
  - `index.Generate` gains an options struct: `IncludeDrafts`, `Plates`
    (source path → URL). For quarterly it adds derived fields and the
    `issues`/`categories` site-level arrays (see
    [data-model.md](data-model.md)) and checks cross-file rules (unique
    issue numbers; an article must be inside an issue folder).
  - Category key: lowercase, trim, collapse internal whitespace to `-`, and
    strip characters outside `[a-z0-9-]` after NFKD folding (so `Ukázky`
    → `ukazky`). The display label is the first label seen in contents
    order.
- **Rationale:** keeps the old entry point stable for three types while
  adding the per-file-kind schemas 004 called for.
- **Alternatives considered:** changing `ParseItem`'s signature for every
  caller (riskier refactor, no benefit to other types).

## R11. Template layout and the `.gitignore` gap

- **Decision:** `templates/site/quarterly/` is copied from `templates/site/log`
  and adapted. It ships a `.gitignore` covering `public/plates/` and
  `.biblio/`. `scaffold.Walk` currently skips every dotfile, so it would
  drop that file; narrow the skip to `.gitkeep` only. The other three
  templates contain no other dotfiles apart from `.github/` (a directory,
  already copied), so their scaffold output is unchanged. A test pins that.
- **Rationale:** generated plates should never be committed by the author.
- **Alternatives considered:** shipping the file as `gitignore` and renaming
  it at scaffold time (a special case in the generator for one file).

## R12. Routes and prerendering

- **Decision:** `react-router.config.ts` prerenders `/`, `/issues`,
  `/issues/<issue>` for each published issue, `/category/<key>` for each
  category, each article slug (`/issues/<issue>/<article>`), and `/colophon`
  if present. `routes.ts` declares `index`, `issues`, `issues/:issue`,
  `category/:category`, and the `*` splat for article and colophon detail,
  as in the existing templates.
- **Rationale:** the same index-driven prerender list the other templates
  use; the index's `issues` and `categories` arrays mean the config never
  parses paths.

## R13. MCP tools

- **Decision:** `list_issues`, `get_issue`, `list_articles`, `get_article`
  are registered in `app/mcp/tools/quarterly.ts` next to the shared tools,
  following `log.ts`'s pattern (zod input schemas, `textResult`). Unlike
  log's `get_entry`, which throws for an unknown slug, the quarterly tools
  return a "not found" text result (FR-034). `get_article` returns the text
  answer as `content` and the notes as `structuredContent` (MCP SDK 1.29
  supports it), with the notes also listed in the text answer so
  terminal hosts lose nothing. Contract: [contracts/mcp-tools.md](contracts/mcp-tools.md).
- **Rationale:** text-first (Principle IV), and structured data for hosts
  that can use it.

## R14. Fonts in the template

- **Decision:** `public/fonts/` ships `Archivo[wdth,wght].woff2` (Latin and
  Latin-Extended subsets), `IBMPlexMono-Regular.woff2`,
  `IBMPlexMono-Medium.woff2`, and `OFL-Archivo.txt` and
  `OFL-IBMPlexMono.txt`. `@font-face` rules live in `public/quarterly.css`
  with `font-display: swap`. Target total under 350 KB.
- **Rationale:** self-hosted per FR-027; woff2 only for the web (005 adds
  OTF/TTF later for Typst).
- **Alternatives considered:** a font CDN (rejected in 004 for privacy);
  separate static weights for Archivo (more files, no gain).

## R15. Testing approach (Principle VI)

- **Decision:** write each test before the code it covers, confirm it fails
  for the right reason, then implement.
  - **Go unit:** typesetting validation (one table row per key, valid and
    invalid), tokens generation (golden file), preset resolution, image
    conversion (fixtures: color JPEG, PNG with alpha, WebP, corrupt file,
    over-cap header, symlink escaping `content/`), the plate cache (skip on
    rerun, cleanup of orphans), the notes parser (golden), the
    `markdown_body` rewrite, per-kind schemas, drafts, category keys, issue
    ordering, and duplicate issue numbers.
  - **Go integration** (`test/integration/`): scaffold a quarterly site and
    run the pre-build step end to end (`pilot-site-quarterly`); extend
    `cross-type-review` to four types; extend `update-three-sites` with a
    quarterly upgrade case; a regression test that profile, log and docs
    pre-build output differs only in the automatically written index.
  - **Template (Vitest):** route tests per page, a notes rendering test
    (wide and narrow markup, `margin_notes: false`), the
    notes-agree-with-index contract test (R2), negotiation (Markdown view
    uses `markdown_body`), `llms.txt`, MCP tool contract tests, and a quarterly
    `app/routes/theme.test.tsx` modeled on the profile template's (no inline
    styles), plus a new rule that `quarterly.css` contains no literal colors
    or font family names outside `@font-face`.
- **Rationale:** mirrors the existing split; each requirement maps to at
  least one test before implementation.

## Dependencies added

| Dependency | Where | Why |
|---|---|---|
| `github.com/yuin/goldmark` v1.7.13 (indirect → direct) | Go | Notes extraction, image destinations (R2, R5) |
| `golang.org/x/image` (`draw`, `webp`) | Go | Resize and WebP decode (R4); the only new module |
| `golang.org/x/text` v0.30.0 (indirect → direct) | Go | `unicode/norm` for category key folding (R10) |

No new npm dependencies: `remark-gfm` and `react-markdown` are already in
the log template, and the margin-note remark plugin is template-local
code.
