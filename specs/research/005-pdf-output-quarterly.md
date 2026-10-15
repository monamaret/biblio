# Research: PDF output for `quarterly` sites

**Status:** Draft
**Date:** 2026-10-15
**Owner:** Project owner

## Problem Statement

A `quarterly` site ([004-quarterly-site-type.md](004-quarterly-site-type.md))
looks like a 1970s academic book, so it should also produce real book-like
output. Readers should be able to download either of two PDFs:

- **A full issue,** set like a printed number of the journal: cover,
  contents with real page numbers, the articles, and the colophon.
- **A single article,** set as an offprint (*separatum*), the way academic
  presses printed individual articles for authors and readers.

Pages are **A5**. Each PDF also comes in a **print edition**: A5 pages
imposed two-up on landscape **A4**, double-sided, so the sheets fold in half
into booklets. Each article is its own folded booklet (a *signature*). A full
issue is a series of signatures gathered in order, so a reader can print it
on an ordinary home printer, fold each article's sheets, and stack them into
an issue.

- **Goal:**
  - `biblio-cli build` produces, for a quarterly site, one PDF per issue
    and one per article. Each comes in two editions: a reading edition
    (A5 pages in reading order) and a print edition (imposed A4 sheets).
  - The PDFs are served as static assets, linked from the issue and article
    pages, and discoverable by agents.
- **In scope:**
  - Paged layout: A5 pages, margins, recto/verso.
  - Imposition for home printers: two A5 pages per A4 side, booklet order,
    one signature per article, and signature marks for gathering.
  - Running heads and folios.
  - Footnotes at the foot of the page.
  - Margin notes in the outer margin.
  - Embedded self-hosted fonts.
  - Grayscale plates.
  - A contents page with page numbers.
  - PDF metadata and bookmarks.
  - Typesetting driven by the same preset and `site.yaml` `typesetting:`
    values as the site.
  - Delivery through links, content negotiation, `llms.txt` and MCP.
- **Out of scope:**
  - PDF output for profile, log or docs. The design should not block it
    later.
  - Generating PDFs on request in the Worker.
  - EPUB and other e-book formats.
  - Print-shop prepress (bleed, crop marks, CMYK, large-sheet imposition).
    Only home-printer A4 imposition is in scope.
  - Imposition for US Letter paper. Letter users print the A4 edition with
    "fit to page" (about 94%); see Risks.
  - Binding beyond folding and gathering (stapling, sewing, covers that
    wrap the whole issue).
  - Selling or gating downloads.
- **Success criteria:**
  - A fixture issue (cover, 3+ articles with dense footnotes and margin
    notes, plates, colophon) builds to an issue PDF.
  - In that PDF:
    - footnotes land on the page that references them;
    - margin notes sit beside the line that references them;
    - contents page numbers are correct;
    - every article starts a new signature, and so a right-hand page;
    - every signature's page count is a multiple of 4;
    - fonts are embedded;
    - the outline (bookmarks) matches the contents.
  - The fixture issue's print edition, printed double-sided on A4 at 100%
    (flip on short edge), then folded and gathered by its signature marks,
    reads in the right page order with nothing clipped by the printer
    margins.
  - Changing a `typesetting:` value changes both the site and the PDFs on
    the next build.
  - Builds with unchanged content skip PDF generation.
  - The deploy workflow produces the same PDFs in CI as locally.
- **Non-goals:**
  - Pixel parity between the web and PDF layouts. They share the same
    design values, not the same layout.
  - WYSIWYG editing of the PDF.

## Background & Context

- **Relevant prior work:**
  - 004 defines the content model:
    - `content/issues/<NN-slug>/issue.md` plus article files and an
      optional `colophon.md`;
    - footnotes written as `[^n]`;
    - margin notes written as `{>> … <<}`;
    - build-time grayscale plates in `public/plates/`;
    - the `typesetting:` schema;
    - the build step that runs before Vite.
  - The existing templates already serve each page as HTML, Markdown or
    ANSI from the same URL, chosen by the request's `Accept` header
    (`app/negotiation/accept.ts`). PDF can be a fourth format.
- **Domain assumptions:**
  - PDFs are built ahead of time and served statically. Issues rarely
    change once published.
  - Size and timing estimates assume about 5–15 articles per issue and
    about 30–120 pages per issue PDF.
  - Readers read on screen and also print at home on A4 duplex printers.
    The print edition is a first-class output, not an afterthought.
  - A5 is exactly half of A4, so two A5 pages fit a landscape A4 side
    with no scaling.
  - Home printers can't print to the paper edge: there is roughly 3–5 mm
    of unprintable margin. Page margins must stay clear of it.
- **Stakeholders:** the project owner, who is the author and the only user.

### Book conventions to reproduce

| Element | Issue PDF | Article offprint |
|---|---|---|
| Cover | Masthead, number and season, cover plate, accent rule | Header line: "Offprint from *Quarterly* No. 3 · Autumn 2026" |
| Contents | Grouped by category, roman numerals, dotted leaders, page numbers | — |
| Article opening | Starts on a right-hand page; category label in the accent color; title, author, dek; optional drop cap or small-caps lead-in | Same, on page 1 |
| Running heads | Left pages: publication and issue; right pages: article title | Left pages: author; right pages: short title |
| Folios | Mono, outer corner, continuous through the issue | Restart at 1, or keep the issue's page numbers (see Open Questions) |
| Notes | Footnotes at the foot of the page, numbered per article; margin notes in the outer margin | Same |
| Plates | Grayscale, captioned "Fig. n —" | Same |
| Back matter | Colophon (typefaces, masthead, edition) on the last page of its own closing signature | Short colophon line on the last page |
| Signatures | Front matter, then one signature per article (split if long), then the closing signature | One signature (split if long) |
| Signature line | Mono line at the foot of each signature's first page, e.g. `QUARTERLY 3 · 2`, after the norm line of Czech and Polish books of the period | Same, with the offprint's article number |

## Libraries & Stack

The deciding requirement is real paged-media typesetting: footnotes placed at
the foot of the page that references them, margin notes, running heads that
follow the current article, and contents page numbers worked out from the
final layout. Browsers' built-in print CSS can't do most of this.

| Candidate | Purpose | License | Maturity | Maintenance | Fit | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| **Typst** (CLI) | Typesetting engine: Markdown converted to Typst markup, then PDF | Apache-2.0 | Young but fast-moving; stable for production PDFs | Active | **High** | Native footnotes, page-dependent running heads, outline from headings, hyphenation and justification, OpenType features (small caps, oldstyle figures). Single static binary, fast compiles. Margin notes need a layout helper (community packages exist, or a small hand-rolled one) |
| **Paged.js + headless Chromium** (Playwright/Puppeteer) | Polyfill for the CSS paged-media specs (incl. footnotes), printed by Chromium | MIT (Paged.js); BSD/Apache (Chromium) | Mature-ish; maintenance cadence has slowed | Moderate | Rejected | Reuses the site's HTML and `quarterly.css` plus a print stylesheet, so styling has one source. Supports footnotes, margin boxes, running strings and contents page numbers. Cost: a Chromium download (~150 MB+) in local builds and CI, and slower renders |
| **WeasyPrint** | HTML+CSS to PDF with paged-media support | BSD-3 | Mature | Active | Rejected | Good paged-media support including footnotes. Adds a Python runtime and its system libraries to a Go+Node toolchain. Rejected |
| Vivliostyle | CSS typesetting engine | AGPL-3.0 (core) | Mature | Active | Low | Copyleft licensing of the core is a poor fit for something bundled into every scaffolded site |
| Pandoc + LaTeX | Markdown to TeX to PDF | GPL / LPPL | Very mature | Active | Low | Excellent typography, but a multi-GB TeX install and a Haskell binary. Too heavy for a single-author CLI |
| Prince / DocRaptor | Commercial CSS-to-PDF | Proprietary | Very mature | Active | Low | Best paged-media support, but needs a paid license. Rejected |
| Go PDF libraries (gopdf, maroto, archived gofpdf) | Low-level PDF drawing | MIT | Varied | Varied | Low | No line breaking, hyphenation, footnote placement or justification. We would be writing a typesetter. Rejected |
| Browser print (`@media print` only) | Reader prints from the page | — | — | — | Low (as primary) | Zero cost, but Chromium doesn't implement footnote floats or page-margin content. Worth shipping anyway as a print stylesheet for article pages |
| Cloudflare Browser Rendering | Headless Chromium in Workers | Proprietary service | GA | Active | Low | Renders at request time, needs a paid plan, and costs per render. Static pre-built PDFs make it unnecessary |
| `github.com/pdfcpu/pdfcpu` | Imposition: reorder the A5 pages into booklet order and place them two-up on A4 landscape | Apache-2.0 | Mature | Active | High | Pure Go library, used in-process, with no extra binary. Places each source page as a reusable vector form, so text stays vector and fonts are embedded once. Has page-selection, N-up and booklet features; which of them we use is settled in the spike |
| Typst PDF-as-image (second Typst pass) | Imposition alternative: a Typst document that places pages of the compiled A5 PDF two-up | Apache-2.0 | Depends on the pinned version | Active | Medium | Keeps everything in one tool, but depends on the pinned Typst version's support for placing PDF pages as images. To be checked in the spike; pdfcpu is the default |
| `github.com/ulikunitz/xz` | Unpack Typst's `.tar.xz` release archives (Linux/macOS) | BSD-3 | Mature | Active | High | Pure Go. The Go standard library has `archive/tar` and `archive/zip` (Windows archives) but no xz decoder. To be confirmed against the pinned release's archive format |
| `goldmark` (already an indirect dependency) | Parse Markdown in Go | MIT | Mature | Active | High (with Typst) | Reads the article Markdown into a syntax tree, extended with footnote parsing and a margin-note parser. A custom renderer then writes Typst markup |

- **Selected (decided 2026-10-15):**
  - **Typst**, fed by a Go `goldmark`-to-Typst renderer, with Typst page
    templates generated from the theme preset.
  - There is no browser-based fallback.
  - **pdfcpu** for imposing the print edition (provisional; the spike may
    switch it to a second Typst pass).
- **Rejected (with reason):**
  - Paged.js plus Chromium: Chromium's download size, install time and
    render overhead in both local builds and CI. Decided by the project
    owner.
  - LaTeX: install weight.
  - Prince: license cost.
  - Vivliostyle: AGPL.
  - Go PDF libraries: no typesetting.
  - Browser Rendering: runtime cost, and not needed for static output.
  - WeasyPrint: adds a third language runtime (Python and its system
    libraries).
- **Typst provisioning (decided 2026-10-15):** `biblio-cli` downloads a
  pinned, checksum-verified Typst release on first use. See
  [Typst provisioning](#typst-provisioning).
- **Open questions:**
  - Which Typst version should we pin?
  - How should margin notes be implemented in Typst?
  - See Open Questions below.

## Architecture Patterns

- **Candidate patterns:**
  1. **Build-time, CLI-owned rendering.** `biblio-cli build` renders PDFs
     into `public/pdf/` before Vite runs. Vite copies them into
     `dist/client`, and Workers Static Assets serves them.
  2. **Render the HTML, then print it.** After `vite build`, Chromium plus
     Paged.js opens the prerendered issue and article pages from
     `dist/client` and prints each to PDF. Rejected: Chromium overhead.
  3. **Render on request in the Worker.** Rejected: runtime cost, a paid
     plan, and slow responses.
  4. **Print in the reader's browser.** Kept only as a supplementary print
     stylesheet.
- **Selected pattern(s):**
  - Pattern 1 with Typst.
  - Pattern 4 ships as well.

### Pipeline

This extends 004's pre-build step:

```
biblio-cli build (site_type: quarterly, pdf.enabled)
  1. validate site.yaml (+ typesetting:, + pdf:)          ← 004
  2. regenerate public/tokens.css                          ← 004
  3. convert images → public/plates/                       ← 004
  4. index → content-index.json                            ← 004 (adds footnotes/marginNotes)
  5. PDF:
     a. generate theme.typ from preset + typesetting overrides (Go, internal/theme)
     b. for each article: goldmark AST → article.typ (Go, internal/pdf)
     c. assemble issue.typ (front matter, articles, closing signature) and offprint.typ
        per article; every signature is padded with blank pages to a multiple of 4
     d. ensure the pinned Typst binary is present (download + verify on first use)
        typst compile → public/pdf/issues/<issue>.pdf, public/pdf/<article-slug>.pdf   (A5 reading)
     e. impose (Go, pdfcpu) → public/pdf/issues/<issue>-print.pdf,
        public/pdf/<article-slug>-print.pdf                                            (A4 print)
     f. skip any target whose input hash is unchanged (.biblio/pdf-cache.json)
  6. vite build
```

### Page format and imposition

**The A5 page (148 × 210 mm).** One page design is used for both
editions. Suggested starting grid for a right-hand page, measured from the
fold outward, to be tuned in the spike:

| Inner margin (fold) | Text block | Gap | Margin-note column | Outer margin |
|---|---|---|---|---|
| 14 mm | 86 mm | 4 mm | 30 mm | 14 mm |

- Top and bottom margins are at least 12 mm, so running heads, folios and
  signature lines stay clear of printers' unprintable edge (about 3–5 mm).
- Left-hand pages mirror the grid, so margin notes always sit in the outer
  margin.
- With `margin_notes: false`, the text block widens into the note column.
- The print base size is about 9.5 pt, mapped from the `typesetting:`
  `base_size` through the print table in `internal/theme`.

**Signatures.** A signature is a group of A4 sheets folded together in
half; each sheet carries 4 A5 pages.

- The issue is divided into signatures:
  1. **Front matter:** cover, masthead, contents. Padded to 4 pages.
  2. **One signature per article.** It is padded with blank pages to a
     multiple of 4, so the next article always starts a new signature on a
     right-hand page. An article longer than `max_sheets × 4` pages (default
     5 sheets = 20 pages) is split into consecutive signatures. Only the
     last one is padded.
  3. **Closing signature:** the colophon on its last page, so it ends up
     on the back of the gathered issue.
- Padding happens in the Typst layout itself, not during imposition, so
  the reading edition has the same pages and page numbers as the print
  edition, and the contents page is correct for both. Blank pages count
  toward page numbering but show no folio or running head.
- An offprint is the article's signature(s), so the offprint's print
  edition folds into a stand-alone booklet.

**Imposition (booklet order).** For a signature of *n* pages (*n* a
multiple of 4), sheet *s* (0-based) carries:

| Side | Left half | Right half |
|---|---|---|
| Front | page *n − 2s* | page *2s + 1* |
| Back | page *2s + 2* | page *n − 2s − 1* |

Example, an 8-page signature (2 sheets):

```
sheet 1 front: [ 8 | 1 ]    back: [ 2 | 7 ]
sheet 2 front: [ 6 | 3 ]    back: [ 4 | 5 ]
```

Sheet 2 nests inside sheet 1; fold the pair and it reads 1–8.

- Signatures follow one another in the print PDF; they are never nested
  inside each other.
- The print PDF is A4 landscape, two A5 pages per side, at 100% scale, and
  is meant to be printed duplex with **flip on short edge**.
- Imposition only reorders and places the finished A5 pages. It never
  re-typesets anything.

**Gathering aids** (each can be turned off under `pdf.print`):

- **Signature line:** a small mono line at the foot of each signature's
  first page, e.g. `QUARTERLY 3 · 2` (issue 3, signature 2). This follows
  the norm line printed in Czech and Polish books of the period.
- **Collating marks:** a small black block printed on the fold of each
  signature's outer sheet, one step lower per signature. Gathered in
  order, the marks form a diagonal staircase along the spine, which makes
  a misplaced signature obvious.
- **Fold marks:** short ticks at the top and bottom of the centre line,
  inside the printable area.
- **Creep** (inner sheets of a folded signature sticking out at the
  front edge) is ignored. At 5 sheets or fewer of ordinary 80 g/m² paper it
  is well under a millimetre.

**Printing instructions.** A one-page `How to print` sheet (A4 portrait)
is prepended to the print edition. It is excluded from imposition and
covers:

- the printer settings (A4, actual size/100%, double-sided, flip on short
  edge);
- how to fold, nest and gather the signatures;
- a fallback for printers without duplex: print odd sheets, reinsert the
  stack, then print even sheets.

### Commands: `build`, `pdf`, and `dev`

- **`biblio-cli build`** runs the whole pipeline above. It produces every
  PDF the `pdf:` settings enable before Vite runs, so deploys always ship
  current PDFs.
- **`biblio-cli pdf`** is a new command for producing and previewing PDFs
  without a full site build. It runs steps 1, 3, 4 and 5. Step 2
  (`tokens.css`) isn't needed for PDFs, and Vite doesn't run.

  ```
  biblio-cli pdf                       # every enabled issue + article PDF
  biblio-cli pdf --issue 3             # one issue (number or folder slug)
  biblio-cli pdf --article on-margins  # one offprint (article slug)
  biblio-cli pdf --issue 3 --open      # build, then open in the system PDF viewer
  biblio-cli pdf --force               # ignore the PDF cache and rebuild
  ```

  - It always produces both editions (reading and print) of each target.
  - It writes to the same `public/pdf/` paths and uses the same
    `.biblio/pdf-cache.json` as `build`. A later `build` reuses anything it
    already rendered.
  - It reports each PDF's path, page count, size and whether it was built
    or skipped, like `build`.
  - `--issue` and `--article` together are an error. With neither, it
    renders everything enabled.
  - `--open` uses the platform opener (`xdg-open`, `open`, or
    `start`/`rundll32` on Windows). With more than one output, it opens
    only the first.
  - It fails with a clear message on non-quarterly sites, and when
    `pdf.enabled: false`.
- **`biblio-cli dev` never builds PDFs** and never downloads Typst, so the
  preview loop stays fast and works offline.
  - PDF links on pages render only when the index has a `pdf` field and
    the file exists in `public/pdf/`.
  - So during `dev`, links appear for PDFs left by an earlier `build` or
    `pdf` run (possibly stale), and are absent otherwise. They are never
    broken.
  - To check print output while editing, run `biblio-cli pdf --open` in a
    second terminal.

### Typst provisioning

`biblio-cli` downloads Typst itself; the author never installs it. The same
code runs locally and in CI, so both use the same Typst version.

- **The pin.** One file, `internal/pdf/typst_pin.go`, holds:
  - the pinned Typst version;
  - for each supported platform, the official release asset name from the
    `typst/typst` GitHub releases and its SHA-256 checksum.

  The checksums are computed by us and reviewed when the pin changes. A
  script (`scripts/update-typst-pin.sh <version>`) downloads every platform
  asset and rewrites the file. Typst is upgraded only in a `biblio-cli`
  release, together with the golden-file tests.
- **Platforms.** Every platform `biblio-cli` ships for through goreleaser
  must have a matching Typst asset: Linux, macOS and Windows on amd64 and
  arm64. On a platform without a pinned asset, PDF generation fails with an
  error that names the platform and points to the override below.
- **On first use.** This happens on `biblio-cli build` with `pdf.enabled`,
  or on `biblio-cli pdf`.
  1. Look for the binary at
     `<os.UserCacheDir()>/biblio/typst/<version>/typst[.exe]`.
  2. If it's missing, download the asset over HTTPS from
     `https://github.com/typst/typst/releases/download/v<version>/<asset>`,
     printing progress.
  3. Verify the SHA-256 against the pin **before** unpacking. On a
     mismatch, delete the download and fail.
  4. Unpack only the `typst` executable entry from the archive. Reject any
     path-traversal entries and cap the extracted size.
  5. Write to a temporary file in the same directory, `chmod 0755`, and
     rename it into place atomically. A lock file in the version directory
     keeps two concurrent builds from both writing it.
  6. Remove cached versions other than the pinned one.
- **Every later run** uses the cached binary with no network access. The
  binary's checksum is not re-verified on each run; the cache directory is
  trusted like any other user-owned tool cache.
- **The binary on `PATH` is never used implicitly,** so every build gets
  the pinned version.
- **Override (escape hatch).** `BIBLIO_TYPST=/path/to/typst` points at a
  local binary, for offline or air-gapped machines or unsupported
  platforms. `biblio-cli` runs `typst --version`. If the version differs
  from the pin, it fails, unless `BIBLIO_TYPST_ALLOW_MISMATCH=1` is set, in
  which case it prints a warning and continues.
- **Offline first run.** With no cached binary and no network, the build
  fails with a message naming the pinned version, the expected cache path,
  and the `BIBLIO_TYPST` override.
- **CI.** The template's `deploy.yml` doesn't install Typst separately.
  `biblio-cli build` downloads it like it does locally. An `actions/cache`
  step keyed on the `biblio-cli` version (the Typst pin only changes with
  a `biblio-cli` release) restores the cache directory, so most CI runs
  make no download.
- **Licensing.** Typst is Apache-2.0. `biblio-cli` downloads it from
  upstream at run time rather than redistributing it, so no bundled NOTICE
  is needed. The colophon credits Typst anyway.

- **Single source for design values.** `internal/theme` emits both
  `tokens.css` (web) and `theme.typ` (print) from the same merged
  preset-plus-override values. This keeps 004's rule that theming logic
  lives only in Go. CSS-only values such as `measure` in `em` are mapped to
  print equivalents through a small table.
- **Fonts.** Typst is pointed at the same OFL woff2/otf files the template
  ships in `public/fonts/` (using `--font-path`). OFL allows embedding fonts
  in PDFs. Typst can't read woff2, so the template also needs `.otf`/`.ttf`
  copies of each face.
- **Plates.** Reuse the converted grayscale images from step 3. For print,
  cap them at roughly 300 DPI for the A5 text-block width (86 mm, about
  1000 px) rather than the web width.
- **Notes.** `[^n]` maps to Typst `#footnote[...]`. `{>> … <<}` maps to a
  margin-note helper in the generated template. Both parse into the same Go
  syntax tree that step 4 uses for `get_article`'s structured notes. **This
  settles 004's open question on structured notes in favor of extracting
  them in Go at index time.**
- **Contents.** Typst works out where each article starts while laying out
  the pages, so the contents page numbers come straight from the engine.
- **Links.** Cross-references and URLs stay clickable. Each PDF's metadata
  carries the canonical page URL.

### Delivery

| Surface | Change |
|---|---|
| Static assets | Reading editions at `/pdf/issues/<issue>.pdf` and `/pdf/<article-slug>.pdf`; print editions at `/pdf/issues/<issue>-print.pdf` and `/pdf/<article-slug>-print.pdf`. Immutable caching keyed by content hash in the filename or `ETag` |
| Issue and article pages | Two links in the masthead or byline area, set in mono, e.g. `PDF · A5 · 48 pp.` and `Print · A4 folded · 12 sheets` |
| Content negotiation | `Accept: application/pdf` on an article or issue URL returns the reading edition. This extends `app/negotiation/accept.ts`, consistent with the existing HTML/Markdown/ANSI behavior. The print edition is a link only |
| `<link rel="alternate">` | `type="application/pdf"` on article and issue pages |
| `llms.txt` | Lists the issue PDFs |
| MCP | `get_issue` and `get_article` results add a `pdf` URL. Optionally a `resource_link` content item for hosts that support it. The text-first rule is unchanged |
| Print stylesheet | `@media print` rules in `quarterly.css` for readers who print the HTML instead |

### Configuration (`site.yaml`)

```yaml
pdf:
  enabled: true            # default true for quarterly
  issues: true
  articles: true
  folio_style: continuous  # continuous | per-article (offprints)
  print:
    enabled: true          # produce the A4 imposed print editions
    max_sheets: 5          # sheets per signature before a long article is split
    signature_line: true
    collating_marks: true
    fold_marks: true
    instructions: true     # prepend the "How to print" sheet
```

There is no `page_size` setting: the page is always A5 and the print sheet
is always A4. Fixing the format keeps the imposition, the page grid and
the gathering aids simple and testable. There is no `binding` setting
either: one signature per article already makes every article open on a
right-hand page.

The `typesetting:` values (fonts, scale, accent, ornaments, alignment, margin
notes) apply to the PDFs automatically. Settings that only matter in print
live under `pdf:`.

- **Tradeoffs:**
  - Typst means a second layout definition (Typst templates alongside
    `quarterly.css`), so the web and print layouts can drift apart. The
    shared design values from Go reduce, but don't remove, that risk.
    This was accepted over Paged.js's single stylesheet to avoid
    Chromium.
- **Integration points:**
  - `internal/theme`: emit `theme.typ`.
  - New package `internal/pdf`: the Markdown-to-Typst renderer, issue and
    offprint assembly, signature planning and imposition (pdfcpu), the
    PDF cache, and Typst provisioning (the pin,
    download, verification, unpacking and the `BIBLIO_TYPST` override)
    plus the `typst` invocation.
  - `internal/content`/`internal/index`: shared note parsing.
  - `internal/cli/build.go`.
  - New `internal/cli/pdf.go`: the `biblio-cli pdf` command (`--issue`,
    `--article`, `--open`, `--force`).
  - `templates/site/quarterly/`: links, negotiation, `llms.txt`, MCP tool
    results, the print stylesheet, and static `.typ` page templates.
  - The template's `deploy.yml`: an `actions/cache` step for the Typst
    cache directory. No separate Typst install step.
  - `.goreleaser.yaml` (linux, darwin and windows on amd64 and arm64): the release checklist confirms every build
    target has a pinned Typst asset.
- **Data model implications:**
  - Index entries gain `pdf` (URL), `pdf_pages` and `pdf_bytes` for
    articles and issues, plus `pdf_print` (URL) and `pdf_print_sheets` for
    the print edition.
  - A new `.biblio/pdf-cache.json`, gitignored.

## Existing Codebase Evaluation

- **Relevant modules:**
  - `internal/theme` (`tokens.go` is the model for a `typ.go` generator).
  - `internal/cli/build.go`. Currently a thin `npx vite build` wrapper; 004
    already adds a pre-step.
  - `internal/devserver`. No changes: `dev` never builds PDFs.
  - `internal/scaffold/manifest.go`, for `update`/`upgrade` tracking of the
    new `.typ` templates.
  - Template `app/negotiation/accept.ts` and `app/discovery/llms-txt.ts`.
- **Reuse opportunities:**
  - `goldmark` is already in the module graph (via glamour). Its footnote
    extension and custom-renderer API fit the Markdown-to-Typst conversion
    well.
  - `runWrapped` in `internal/cli` already handles running and streaming
    an external process, so `typst compile` can use it.
- **Refactor candidates:**
  - `GenerateTokensCSS` should turn into "resolve the effective design
    values once (preset plus overrides), then write them out as CSS, Typst,
    or more".
- **Replacements required:** none.
- **Code health signals:**
  - Test the Markdown-to-Typst renderer with golden-file Go unit tests (one
    per Markdown construct, plus notes, plates and escaping).
  - Add one integration test that compiles the fixture issue and checks
    the page count, fonts and outline using `pdfinfo`/`pdffonts`, or a Go
    PDF reader in tests only.
- **Conventions to follow:**
  - Theming lives only in Go.
  - `biblio-cli` alone produces generated artifacts.
  - Templates stay independent npm projects.
  - A text-first MCP surface.
  - A research doc comes before the spec.

## Security & Compliance

- **Threat surface:**
  - Author Markdown is translated into Typst markup, so there is an
    injection risk: unescaped `#`, `[`, `]`, `$`, `*`, `_`, `@`, `<`, `` ` ``
    and `\` can run Typst code.
  - `typst compile` reads files from the project.
  - `biblio-cli` downloads and runs a third-party executable, which is a
    supply-chain risk.
- **Data classification:** public published content.
- **Compliance regimes:**
  - None required.
  - Font licenses: ship the OFL texts, and mention them in the colophon.
- **Mitigations:**
  - Escape every text node in the Typst renderer. Never pass raw Markdown
    or HTML through.
  - Fuzz-test the renderer with arbitrary Markdown and assert that the
    compile output has no script evaluation beyond the generated template.
  - Run `typst compile` with `--root` set to the build directory so it
    can't read outside the project.
  - Typst download:
    - HTTPS only, from the fixed `github.com/typst/typst` release URL.
    - SHA-256 pinned in the `biblio-cli` binary itself, so changing the pin
      takes a reviewed `biblio-cli` release.
    - The checksum is verified before unpacking, and a mismatch is fatal.
    - Unpacking takes a single entry, is protected against path traversal,
      and is size-capped.
    - The file is installed with an atomic rename.
  - Vendor every Typst package into the template so compiles never fetch
    from the network.

## Performance, Reliability & Operations

- **Expected load:** build time only. Serving is static.
- **Build cost:**
  - Typst compiles a 100-page document in seconds or less, so a full issue
    plus offprints should take a few seconds. To be confirmed with the
    fixture issue.
  - Hash-based skipping makes unchanged rebuilds close to free.
- **Asset size:**
  - Each issue PDF is roughly 1–10 MB, mostly fonts (subset when embedded)
    and plates.
  - The print edition is about the same size as the reading edition,
    because imposition reuses each page as a vector form rather than
    copying it.
  - Workers Static Assets caps individual files at 25 MiB. The build fails
    clearly if a PDF would exceed it, suggesting a lower plate resolution.
- **Observability:** the build reports each PDF's path, page count, size,
  and whether it was built or skipped.
- **Reliability:**
  - Fail the build on any PDF error by default, because a missing PDF
    behind a "PDF" link is a broken page.
  - An opt-in `pdf.on_error: warn` builds the site without the links.
- **Rollout / migration plan:**
  - Ships with or right after the quarterly site type: v0.08 if it fits,
    otherwise v0.09.
  - Existing quarterly sites get the `.typ` templates and links through
    `biblio-cli upgrade`.
  - Defaults to `pdf.enabled: true`.
- **Rollback plan:** set `pdf.enabled: false`. Removing the step from
  `build` leaves the site intact, because links render only when the index
  has a `pdf` field.

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation | Owner |
| --- | --- | --- | --- | --- |
| Margin notes collide or overflow the page foot in Typst | Medium | Medium | Spike with dense-note fixtures. Fallback setting turns margin notes into footnotes in print | Owner |
| Web and print layouts drift apart | Medium | Low | Shared design values from Go; one visual review checklist covering both | Owner |
| Typst markup or CLI changes between versions break templates | Medium | Medium | Pin the Typst version in `biblio-cli`; golden-file tests; upgrade Typst only in a release | Owner |
| Typst binary missing on the author's machine or in CI | High (first run) | Low | `biblio-cli` downloads the pinned, checksummed binary automatically; CI caches it by version | Owner |
| First-run download fails (offline, GitHub outage, proxy) | Medium | Medium | Clear error with the version, cache path and `BIBLIO_TYPST` override; later runs need no network | Owner |
| Upstream deletes or replaces a release asset | Low | High | Checksum mismatch fails safely. Fix by re-pinning in a `biblio-cli` patch release; `BIBLIO_TYPST` as the stopgap | Owner |
| A goreleaser target has no Typst asset | Low | Medium | Release checklist item; PDF fails with a clear message on that platform, and the site builds with `pdf.enabled: false` | Owner |
| Fonts unavailable in a Typst-readable format | Medium | Medium | Ship `.otf`/`.ttf` alongside woff2 in the template. A build warning names any missing face | Owner |
| PDFs inaccessible to screen readers | Medium | Medium | Require tagged PDF output if the pinned Typst version supports it (verify in the spike); outline, metadata language and image alt text from `image_alt` regardless | Owner |
| Reader prints with the wrong settings (scaled to fit, flip on long edge, single-sided) | High | Medium | Prepended "How to print" sheet; settings repeated in small type on the first sheet's fold margin; the print link's tooltip says "A4, 100%, duplex, flip on short edge" | Owner |
| Printer's unprintable edge clips heads, folios or marks | Medium | Medium | At least 12 mm top and bottom and 14 mm outer margins on the A5 page; marks kept inside the printable area; checked by a real print of the fixture issue | Owner |
| Letter-paper readers can't print at 100% | Medium | Low | Out of scope; the instructions sheet says to use "fit to page" on Letter (about 94%), which still folds correctly | Owner |
| Many blank pages from padding (e.g. a 5-page article becomes 8) | Medium | Low | Accepted as the price of one-article signatures. Blanks carry no folio. The reading edition could drop trailing blanks later if it bothers readers | Owner |
| Long articles produce thick signatures that fold badly | Low | Low | `max_sheets` (default 5) splits long articles into consecutive signatures | Owner |
| Repo or CI size grows from committed PDFs | Low | Low | PDFs are build output only, never committed; `public/pdf/` is gitignored | Owner |

## Open Questions

- [ ] **Typst validation spike.** The engine is decided. Typeset the
  fixture issue in Typst to confirm:
  - footnote and margin-note placement with dense notes;
  - hyphenation and justification quality, including Czech and Polish
    names in the text;
  - build time;
  - tagged/accessible PDF support in the pinned version;
  - padding each signature to a multiple of 4 inside the Typst layout;
  - imposition with pdfcpu (page reordering plus two-up placement)
    against a second Typst pass that places PDF pages as images;
  - a real duplex print of the fixture issue, folded and gathered.

  Owner, before the spec.
- [x] ~~**Typst distribution.**~~ **Resolved 2026-10-15:** `biblio-cli`
  downloads a pinned, checksum-verified copy (see
  [Typst provisioning](#typst-provisioning)). Embedding a WebAssembly
  build and requiring a user install were both declined.
- [x] ~~**Prefetch command.**~~ **Resolved 2026-10-15:** no prefetch
  command. The first PDF-enabled build downloads Typst, and that is
  enough. Offline machines use `BIBLIO_TYPST`.
- [ ] **Margin notes in Typst.** Use a community package or a hand-rolled
  helper? Either way, pin it and vendor it into the template, never fetch
  it at build time. Owner, during the validation spike.
- [x] ~~**Default page size.**~~ **Resolved 2026-10-15:** A5 pages, imposed
  two-up on A4 landscape, with one signature per article (see
  [Page format and imposition](#page-format-and-imposition)).
- [ ] **Offprint numbering.** Should offprints restart at page 1, or keep
  their page numbers from the issue? Since an offprint is now exactly the
  article's signature(s) from the issue, keeping the issue's numbers
  (the traditional behavior) costs almost nothing: the offprint is cut
  from the laid-out issue. Owner, before plan.
- [ ] **Blank pages in the reading edition.** Keep them, so the reading
  and print editions are page-for-page identical, or drop trailing blanks
  from the reading edition? The page numbers would then skip, as they do
  in a printed book. Owner, before plan.
- [x] ~~**`biblio-cli dev`.**~~ **Resolved 2026-10-15:** `dev` never builds
  PDFs. A separate `biblio-cli pdf` command handles it (see
  [Commands](#commands-build-pdf-and-dev)).
- [ ] **Drafts.** If 004 adds `draft: true` issues, drafts should be left
  out of PDFs too. Should `biblio-cli pdf` still be able to preview them?
  Owner, before plan.
- [ ] **Generalizing later.** Keep `internal/pdf` independent of quarterly
  details (the issue and offprint assembly sit behind an interface) so docs
  sites could get "the whole manual as a PDF" later? Owner, before plan.

## Decision Log

- **2026-10-15** — Explore PDF output for individual articles and full
  issues of `quarterly` sites. — It extends the academic-book look of 004
  into an actual book-like artifact. — Project owner
- **2026-10-15** — Build PDFs ahead of time as static assets; no rendering
  on request. — Matches the static, prerendered architecture; no runtime
  cost or paid plan. — Proposed (this doc)
- **2026-10-15** — Typst is the PDF engine, with no Paged.js/Chromium
  fallback. — Avoids Chromium's download, install and render overhead;
  Typst is a single small binary with native paged-media typesetting. —
  Project owner
- **2026-10-15** — `biblio-cli` downloads a pinned, checksum-verified
  Typst release into the user cache on first use. A `BIBLIO_TYPST`
  override covers offline machines and unsupported platforms. — No manual
  install for the author; local and CI builds use the identical version;
  the `biblio-cli` binary stays small. — Project owner
- **2026-10-15** — No separate command to download Typst ahead of time;
  the first PDF-enabled build downloads it. — Keeps the command surface
  small; `BIBLIO_TYPST` already covers offline machines. — Project owner
- **2026-10-15** — `biblio-cli dev` never builds PDFs; a separate
  `biblio-cli pdf` command (`--issue`, `--article`, `--open`, `--force`)
  produces and previews them, sharing `build`'s output paths and cache. —
  Keeps the dev loop fast and offline; PDF checks are deliberate. —
  Project owner
- **2026-10-15** — Pages are A5. Every PDF also gets a print edition:
  pages imposed two-up on landscape A4, double-sided, folding into
  booklets. Each article is its own signature, and an issue is a series of
  signatures. — Readers can print and assemble an issue at home from
  ordinary A4 paper; fits the academic-press look. — Project owner
- **2026-10-15** — Footnote and margin-note parsing moves into Go at index
  time (resolves 004's open question). — Both the PDF renderer and
  `get_article` need the same structured notes. — Proposed (this doc)

## Follow-up

Parked on 2026-10-15 for later.

**Decided by the project owner:**
- Typst is the engine, with no Chromium.
- `biblio-cli` downloads a pinned Typst.
- No separate command to download Typst ahead of time.
- `dev` never builds PDFs; a separate `biblio-cli pdf` command does.
- A5 pages, with an A4 print edition of one signature per article.

**Still proposed:**
- Build PDFs ahead of time as static assets.
- Parse notes in Go at index time.
- Use pdfcpu for imposition.

To pick it up:

1. Pick it up after, or alongside, [004](004-quarterly-site-type.md). This
   feature depends on the quarterly site type existing.
2. Run the Typst validation spike from Open Questions. It also settles
   pdfcpu versus a second Typst pass for imposition, and margin notes in
   Typst.
3. Resolve the remaining Open Questions:
   - offprint numbering;
   - blank pages in the reading edition;
   - drafts;
   - generalizing `internal/pdf`.
4. Confirm or revise the proposed decisions above.
5. Write the feature spec, then run `/speckit-specify`.

## References

- [004-quarterly-site-type.md](004-quarterly-site-type.md) — content model,
  notes syntax, typesetting schema, build pre-step
- [Typst documentation](https://typst.app/docs/) — page setup, footnotes,
  outline, counters
- [Typst Universe](https://typst.app/universe/) — community packages
  (margin-note helpers)
- [Paged.js](https://pagedjs.org/) — considered and rejected (Chromium
  overhead)
- [WeasyPrint](https://weasyprint.org/) — considered and rejected (Python
  runtime)
- [goldmark](https://github.com/yuin/goldmark) — Go Markdown parser and
  renderer API
- [pdfcpu](https://github.com/pdfcpu/pdfcpu) — pure-Go PDF processing:
  page selection, N-up, booklet
- [Cloudflare Workers Static Assets limits](https://developers.cloudflare.com/workers/static-assets/)
  — per-file size cap
- [SIL Open Font License](https://openfontlicense.org/) — embedding terms

## Appendix

### Issue page and signature sequence (A5 pages, A4 sheets)

```
SIGNATURE 1 — front matter (1 sheet, 4 pp.)
  i    cover
  ii   masthead
  iii  contents
  iv   contents (cont.) or blank

SIGNATURE 2 — I. On Margins, 11 pp. → padded to 12 (3 sheets)
  1    opening page (no folio; signature line "QUARTERLY 3 · 2")
  …    running heads: left "QUARTERLY · No. 3"  right "On Margins"
  11   last text page
  12   blank (no folio)

SIGNATURE 3 — II. The Rigid Grid, 7 pp. → padded to 8 (2 sheets)
  13   opening page (signature line "QUARTERLY 3 · 3")
  …
  20   blank

…

SIGNATURE n — closing (1 sheet, 4 pp.)
       blank · blank · blank · colophon
```

### Print edition, one sheet (A4 landscape, signature 3, sheet 1)

```
 front (flip on short edge)          back
┌──────────────┬──────────────┐     ┌──────────────┬──────────────┐
│              │              │     │              │              │
│   page 20    │   page 13    │     │   page 14    │   page 19    │
│   (blank)    │  The Rigid   │     │              │              │
│              │    Grid      │     │              │              │
│              │ QUARTERLY 3·3│     │              │              │
└──────────────┴──────────────┘     └──────────────┴──────────────┘
               ▲ fold                              ▲ fold
               ■ collating mark (step 3)
```

### Article page sketch (right-hand page)

```
┌──────────────────────────────────────────────┐
│ On Margins                              17   │  ← running head + folio (mono)
│                                              │
│  The grid was rigid, and the       │ cf. Ger-│  ← margin note, outer margin
│  margins were wide.¹ Every page    │ stner,  │
│  of the 1971 Brno catalogue …      │ 1964    │
│                                    │         │
│  [ Fig. 2 — plate, grayscale ]     │         │
│                                    │         │
│  ──────                                      │
│  1  Brno, 1971.                              │  ← footnote at page foot
└──────────────────────────────────────────────┘
```
