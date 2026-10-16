# Research: `quarterly` site type

**Status:** Draft
**Date:** 2026-10-15
**Owner:** Project owner

## Problem Statement

Add a fourth site type, `quarterly`: a magazine-style blog published in
numbered issues and grouped by category. Its look follows a 1970s academic
book in the early-70s Eastern European style: black-and-white plates, lots of
white space, footnotes and margin notes, one bold accent color, and
typesetting the author can configure. It has no comments.

This is the "fourth site type" item from
[003-post-ga-feature-backlog.md](003-post-ga-feature-backlog.md). That entry
said a research doc should wait for a real use case. This is that use case.
Requirements came from a one-question-at-a-time interview on 2026-10-15
(see the Decision Log).

- **Goal:** `biblio-cli new <name> --type=quarterly` scaffolds a working
  quarterly site. It must have the same guarantees as profile, log and docs:
  prerendered pages, content negotiation (HTML, Markdown and ANSI),
  `llms.txt`, an in-Worker MCP server, the deploy workflow, and
  `update`/`upgrade` support.
- **In scope:**
  - Issue-based content layout.
  - Category, issue and colophon routes.
  - Footnotes and margin notes.
  - A `quarterly-default` theme preset that pages actually use, with
    self-hosted fonts.
  - Typesetting overrides set in `content/site.yaml`.
  - Build-time grayscale conversion of images.
  - Quarterly-specific MCP tools.
- **Out of scope:**
  - Comments.
  - Author pages (`/author/<name>`).
  - A `list_categories` MCP tool.
  - Restyling profile, log or docs. They may later reuse the base-stylesheet
    work below, but that is a separate decision.
  - Paged or print output (PDF). This is covered separately in
    [005-pdf-output-quarterly.md](005-pdf-output-quarterly.md).
- **Success criteria:**
  - A freshly scaffolded quarterly site builds, prerenders every issue,
    article, category and colophon page, and deploys.
  - Changing a `typesetting:` value in `site.yaml` changes the rendered
    site on the next `dev`/`build`.
  - Every image under `content/` is served as grayscale.
  - Footnotes and margin notes render correctly on wide and narrow
    viewports, and they degrade sensibly in the Markdown and ANSI views.
  - `get_article` returns the notes as structured data.
  - The cross-type consistency review (as in v0.05) passes with all four
    types.
- **Non-goals:**
  - A general plugin system for site types.
  - Arbitrary user CSS. Configuration stays inside the typesetting schema.

## Background & Context

- **Relevant prior work:**
  - The three existing site types live in `templates/site/{profile,log,docs}`.
    `log` is the closest model: dated entries, tags, and a filtered MCP list
    tool.
  - `docs` brings ordered sections and table-of-contents generation
    (`app/content/toc.ts`).
  - Theme presets live in `internal/theme`, and `tokens.css` is generated at
    `new` into `public/`.
- **Domain assumptions:**
  - There is one author or editor, matching the other types.
  - An issue is the unit of publication: articles appear when their issue
    folder exists.
  - Categories are free-text labels, not a fixed list.
- **Stakeholders:** the project owner, who is the author and the only user.

### Interview answers (2026-10-15)

| # | Topic | Answer |
|---|---|---|
| 1 | Purpose | A magazine-style blog with categories, minimalist and clean, with no comments. It looks like a 1970s academic book in the early-70s Eastern European style: b&w images, white space, footnotes and margin notes, a bold accent, configurable typesetting |
| 2 | Type ID | `quarterly` |
| 3 | Structure | Issues plus categories: `content/issues/<NN-slug>/issue.md` plus article files, each article with a `category` |
| 4 | Article fields | `author`, `dek`, `image` / `image_alt` / `image_caption`, `order`, `tags` |
| 5 | Notes | Both kinds, with separate syntax: `[^n]` footnotes go to the margin on wide screens and the article foot on narrow ones; `{>> … <<}` makes an unnumbered margin note |
| 6 | Home page | The current issue's contents page |
| 7 | Extra routes | Issue archive, category pages, and a colophon/about page. No author pages |
| 8 | Configurable | Typefaces and scale, accent color, book ornaments, page grid |
| 9 | Config location | Preset defaults plus an optional `typesetting:` block in `site.yaml`, with `tokens.css` regenerated on `dev`/`build` |
| 10 | Default look | Grotesk and red: Inter Tight or Archivo for text, IBM Plex Mono for folios, `#d7261e` accent, `#fbfaf7` background, `#111` text, asymmetric grid, ragged right |
| 11 | MCP tools | `list_issues`, `get_issue`, `list_articles` (filtered), `get_article` (with structured notes) |
| 12 | Images | Build-time grayscale conversion by `biblio-cli` |
| 13 | Next step | This research doc first, then speckit spec, plan and tasks |

## Libraries & Stack

| Candidate | Purpose | License | Maturity | Maintenance | Fit | Notes |
| --- | --- | --- | --- | --- | --- | --- |
| Go stdlib `image`, `image/jpeg`, `image/png` | Decode, convert to grayscale and encode images | BSD-3 | Stable | Go team | High | Pure Go, so goreleaser cross-builds stay simple. Grayscale is a pixel loop into `image.Gray`, or luminance with an alpha-preserving variant for PNG |
| `golang.org/x/image/draw` | Resize to a maximum width | BSD-3 | Stable | Go team | High | `CatmullRom` scaler. Only new Go dependency |
| `golang.org/x/image/webp` | Decode WebP input | BSD-3 | Stable | Go team | Medium | Decode only. Output is JPEG or PNG |
| `libvips` / `bimg` | Image processing | LGPL / MIT | Mature | Active | Low | Needs cgo, which breaks pure-Go cross-compilation. Rejected |
| `remark-gfm` (already a dependency) | `[^n]` footnote parsing | MIT | Mature | Active | High | Already in every template. Footnote nodes are restyled as sidenotes in the renderer |
| Custom remark plugin (in the template) | `{>> … <<}` margin notes | — | New | Ours | High | A small mdast transform over text nodes. The syntax is CriticMarkup's comment form, so editors that know CriticMarkup display it reasonably |
| `goldmark` + footnote extension (already an indirect dependency) | Extract notes at index time | MIT | Mature | Active | Medium | One option for structured notes in `content-index.json` (see Open Questions) |
| Inter Tight / Archivo, IBM Plex Mono (woff2) | Default typefaces | OFL-1.1 | Mature | Active | High | Self-hosted in `public/fonts/`. OFL allows redistribution inside the template |

- **Selected:**
  - Go stdlib plus `x/image` for images.
  - `remark-gfm` for footnotes.
  - A template-local remark plugin for margin notes.
  - Self-hosted OFL fonts.
- **Rejected (with reason):**
  - cgo image libraries, because they break cross-compilation.
  - Font CDNs, because they add a third-party request and a privacy cost.
  - Tailwind or shadcn, already superseded per roadmap §3.
- **Open questions:** Archivo or Inter Tight as the default text face
  (decide by setting a sample article in both).

## Architecture Patterns

- **Candidate patterns:**
  1. Copy `templates/site/log` and adapt it into a fourth independent
     template. This is the existing pattern: each type is its own npm
     project.
  2. Share code across templates. Rejected for now. The templates are
     deliberately independent, and `update`/`upgrade` diffs per type.
- **Selected pattern(s):** Pattern 1, plus three new CLI-side pre-build steps.

### Content layout

```
content/
  site.yaml                 site_type: quarterly, optional typesetting:
  colophon.md               optional; title required
  issues/
    01-autumn-2026/
      issue.md              title, number, date required; cover, cover_alt optional
      cover.jpg
      on-margins.md         title, category required; author, dek,
      grid-systems.md       image, image_alt, image_caption, order, tags optional
      img/…
```

- An article's issue comes from its folder. The index adds derived fields
  `issue` (folder slug), `issue_number` and `issue_date`, so routes and MCP
  tools never parse paths.
- Article order within an issue is `order` ascending, then filename.
  Category order follows the first appearance by `order`.
- The current issue is the one with the highest `number`.

### Routes

| Route | View |
|---|---|
| `/` | Current issue's contents page: masthead, issue number and season, cover plate, articles grouped by category with roman-numeral numbering, dotted leaders and authors |
| `/issues` | Archive of all issues, newest first |
| `/issues/:issue` | That issue's contents page (same component as `/`) |
| `/category/:category` | Articles in a category across issues |
| `*` (splat) | Article or colophon detail, as in the existing templates |

### Notes rendering

- **Wide viewports:** a two-column text block (main text and margin), with
  `[^n]` footnotes and `{>> … <<}` margin notes beside the line that
  references them.
- **Narrow viewports:** footnotes collect at the article foot. Margin notes
  become inline toggles marked with a dagger (†). This is CSS-only, using
  `<details>`, so no client JS is needed for prerendered pages.
- **Markdown view (`Accept: text/markdown`):** keep `[^n]` as-is. Rewrite
  margin notes as footnotes (`[^m1]`) so the file reads as plain Markdown.
- **ANSI view:** same as Markdown, then the existing ANSI rendering.

### Theming: the first type whose pages use the tokens

Today `tokens.css` only defines custom properties. No page stylesheet reads
them, and `theme.test.tsx` forbids inline styles. Quarterly needs a real
stylesheet:

- **`public/quarterly.css`** is template-owned and static. It reads only
  `var(--…)` tokens and `[data-*]` toggles, and is linked from `root.tsx`
  after `tokens.css`. `theme.test.tsx`'s "no inline styles" rule stays.
  The test gains an assertion that `quarterly.css` contains no literal
  colors or font names, so it consumes tokens only.
- **`tokens.css`** is still generated by Go from the preset plus the
  `site.yaml` overrides. Continuous values (colors, fonts, scale, measure)
  become custom properties. Toggles (drop caps, numerals, rules, alignment,
  margin notes) become `:root` properties that `quarterly.css` reads, for
  example `--ts-align: justify`. Others become generated attribute rules
  (see Open Questions).
- **Fonts:** `@font-face` rules for the default faces live in
  `quarterly.css`. A font override in `site.yaml` names a family. The
  author either adds matching files under `public/fonts/` or accepts a
  fallback to the generic family in the stack.

### Typesetting schema (`site.yaml`)

| Key | Type / values | Default (`quarterly-default`) |
|---|---|---|
| `accent` | hex color | `#d7261e` |
| `font_body`, `font_display`, `font_mono` | family name | Archivo or Inter Tight, the same, IBM Plex Mono |
| `base_size` | CSS length | `1.0625rem` |
| `scale` | ratio 1.1–1.618 | `1.25` |
| `leading` | number | `1.5` |
| `measure` | `em`/`ch` length | `34em` |
| `align` | `left` \| `justify` | `left` (justify turns on `hyphens: auto`) |
| `grid` | `asymmetric` \| `centered` | `asymmetric` |
| `margin_notes` | bool | `true` |
| `drop_caps` | `none` \| `drop` \| `raised` | `none` |
| `lead_in` | `none` \| `small_caps` | `small_caps` |
| `numerals` | `roman` \| `arabic` | `roman` |
| `rules` | `none` \| `hairline` \| `heavy` | `heavy` |
| `fleurons` | bool | `false` |
| `grayscale_images` | bool | `true` |

Every value is validated in Go before it reaches CSS: enums, a hex regex, a
length regex, numeric ranges, and a font-name character set. An invalid
value fails `index`/`build` with a message naming the key.

### Build pipeline changes

`biblio-cli build` and `dev` currently only wrap `vite build` and
`vite dev`/`wrangler dev`. Generating `content-index.json` is a separate
`biblio-cli index` command. For quarterly, a pre-step runs before Vite:

1. Validate `site.yaml`, including the `typesetting:` block.
2. Regenerate `public/tokens.css` from the preset plus the overrides.
3. Convert images: for each image under `content/`, write a grayscale,
   width-capped copy to `public/plates/<hash>.<ext>`. Unchanged images are
   skipped using a content-hash cache, and the index maps source paths to
   plate URLs.
4. Run `index` as today.

For `dev`, the file watcher reruns steps 1–4 when `content/` changes. If
that's too much for the first cut, rerunning only on restart is acceptable.

- **Tradeoffs:**
  - The pre-step makes `build` more than a thin wrapper, which deviates
    from the roadmap §0.3 wording. In exchange, theming and image handling
    stay in Go, consistent with Principle II.
  - Other types can adopt the same pre-step later.
- **Integration points:**
  - `scaffold.SiteTypes`.
  - `content.requiredFields`, which now needs per-file-kind schemas (see
    below).
  - `index.Generate` (derived issue fields, plate paths).
  - `theme` (new preset, an override merge, and a richer generator).
  - `cli/build.go` and `cli/dev.go`.
  - The TUI type list, which updates automatically.
  - `migrate`/manifest, for the new template tree.
- **Data model implications:**
  - Validation is per file kind: `issue.md`, article, and `colophon.md`.
    `content.ParseItem` currently takes only `siteType`, so it needs the
    relative path to choose a schema.
  - New derived index fields: `kind`, `issue`, `issue_number`,
    `issue_date`, `plate`.

### MCP tools (in addition to the shared `list_pages`, `get_page`, `search_content`)

| Tool | Input | Output |
|---|---|---|
| `list_issues` | — | number, title, date, slug, article count; newest first |
| `get_issue` | `issue` (slug or number) | issue metadata plus contents grouped by category |
| `list_articles` | `issue?`, `category?`, `author?`, `tag?`, `dateFrom?`, `dateTo?` | article summaries (title, dek, author, category, issue, slug) |
| `get_article` | `slug` | full body plus `footnotes[]` and `marginNotes[]` as structured objects (`id`, `text`, anchor offset) |

All results are text-first with a plain table or list, following the
existing progressive-enhancement rule.

## Existing Codebase Evaluation

- **Relevant modules:**
  - `internal/scaffold` (`assets.go` SiteTypes, `generator.go`,
    `manifest.go`).
  - `internal/content/schema.go`.
  - `internal/index/generator.go`.
  - `internal/theme`.
  - `internal/cli/{new,build,dev,index}.go`.
  - `internal/tui/new_flow.go`.
  - `templates/site/log/` (the starting copy).
  - `templates/site/docs/app/content/toc.ts` (grouping and ordering logic).
- **Reuse opportunities:**
  - From the log template, copied wholesale: `negotiation/`, `discovery/`,
    `mcp/server.ts`, `mcp/agent.ts`, `mcp/tools/shared.ts`, `workers/`,
    the deploy workflow, and the vite and vitest configs.
  - The filter logic in log's `filters.ts` (date, tag, project) maps
    directly onto `list_articles`.
- **Refactor candidates:**
  - `content.ParseItem`, from type-level to file-kind-level required
    fields.
  - `theme.Preset`, which needs a typed schema for typesetting values
    alongside the free-form `Tokens` map.
  - `theme.GenerateTokensCSS`, which must merge overrides.
- **Replacements required:** none.
- **Code health signals:**
  - Go unit tests exist for theme, scaffold, content and index.
  - Each template has Vitest suites for routes, MCP tools, negotiation and
    `llms.txt`.
  - `scripts/check.sh` is the gate.
- **Conventions to follow:**
  - Templates stay independent npm projects.
  - Theming logic lives only in Go.
  - One styling path through `root.tsx`.
  - Tools are text-first.
  - Each template's `.github/workflows/deploy.yml` is updated alongside.

## Security & Compliance

- **Threat surface:**
  - Author-supplied images are decoded by the CLI. Risks are decompression
    bombs and malformed files.
  - `site.yaml` values are written into CSS, which is an injection risk.
  - Margin-note text is rendered as Markdown.
- **Data classification:** public published content. No personal data
  beyond author bylines the author chooses to publish.
- **Compliance regimes:**
  - None.
  - Font licenses: ship the OFL text alongside the font files.
- **Mitigations:**
  - Images: check dimensions with `image.DecodeConfig` before decoding, and
    reject anything over a pixel cap (e.g. 40 MP). Never follow symlinks
    outside `content/`.
  - CSS: strict typed validation of every typesetting value, and no
    free-form strings in CSS.
  - Markdown: render with react-markdown's default of no raw HTML.

## Performance, Reliability & Operations

- **Expected load:** static prerendered pages served from Workers Static
  Assets. Same as the other types.
- **SLOs / SLIs:** none beyond the existing ones.
- **Build performance:** image conversion is the only new expensive step.
  With hash-based skipping, rebuilds cost close to nothing when images are
  unchanged.
- **Observability:** the pre-step reports how many images it converted and
  skipped, and which tokens it regenerated.
- **Rollout / migration plan:**
  - Ship as a new type in a feature release (v0.08).
  - No impact on existing sites, because profile, log and docs are
    untouched apart from shared Go refactors that keep their behavior.
- **Rollback plan:** remove `quarterly` from `SiteTypes`. Existing quarterly
  sites keep working because their templates are vendored.

## Risks & Mitigations

| Risk | Likelihood | Impact | Mitigation | Owner |
| --- | --- | --- | --- | --- |
| The `build` pre-step drifts from "thin wrapper" and diverges between types | Medium | Medium | Gate the pre-step on `site_type`; record a Decision Log entry; consider extending it to all types later | Owner |
| Sidenote layout fails at awkward widths or with long notes | Medium | Medium | CSS-only fallback to foot notes below a breakpoint; visual check of fixture articles with dense notes | Owner |
| Overriding `font_*` without the files silently falls back | High | Low | `build` warns when a configured family has no matching `public/fonts/` file | Owner |
| Refactoring `content.ParseItem` breaks existing types | Low | High | Keep the per-type defaults; existing type tests must pass unchanged | Owner |
| Grayscale conversion damages images the author wanted in color | Low | Low | `grayscale_images: false` override | Owner |
| Template size grows from bundled fonts | Medium | Low | woff2 only, Latin and Latin-Extended subsets, which also cover Czech and Polish diacritics | Owner |

## Open Questions

- [x] ~~Default text face: Archivo or Inter Tight?~~ **Resolved
  2026-10-15:** Archivo (plan research R1). See the Decision Log.
- [x] ~~Structured notes: extract them in Go at index time, or parse them
  in TS when `get_article` runs?~~ **Resolved 2026-10-15:** in Go, at
  index time (plan research R2). See the Decision Log.
- [x] ~~Ornament toggles: pure `:root` custom properties, or also
  generated `[data-…]` rule blocks in `tokens.css`?~~ **Resolved
  2026-10-15:** `:root` properties only; structural switches go in the
  index (plan research R3). See the Decision Log.
- [x] ~~Should the pre-step (tokens, images, index) apply to all four types?~~
  **Resolved 2026-10-15:** partly. Settings validation and `index` run
  automatically for all four types; tokens and images stay quarterly-only.
  See the Decision Log.
- [x] ~~Image output: keep the source format, or always write JPEG for
  photos and PNG for images with alpha? What maximum width?~~ **Resolved
  2026-10-15:** keep JPEG/PNG; WebP becomes PNG if it has transparency,
  otherwise JPEG; 1600 px max, 40 MP cap (plan research R4). See the
  Decision Log.
- [x] ~~Can an issue be marked as a draft (an `issue.md` `draft: true` field
  excluded from build and index)?~~ **Resolved 2026-10-15:** yes; drafts
  are shown in `dev`. See the Decision Log.

## Decision Log

- **2026-10-15** — Add a fourth site type, `quarterly`, an issue-based
  magazine. — It meets backlog 003's bar of a real driving use case. —
  Project owner (interview)
- **2026-10-15** — Content layout is issues plus categories (folder per
  issue, `category` frontmatter). — Matches the periodical model and gives
  the contents page its structure. — Project owner
- **2026-10-15** — Typesetting is configured as preset defaults plus
  `site.yaml` overrides, with `tokens.css` regenerated on `dev`/`build`. —
  Per-site control without hand-edited CSS; theming logic stays in Go. —
  Project owner
- **2026-10-15** — Images are converted to grayscale at build time in Go
  (pure Go, no cgo). — True b&w and smaller files; keeps cross-compilation.
  — Project owner
- **2026-10-15** — Footnotes `[^n]` and margin notes `{>> … <<}` are
  separate syntaxes. — Numbered scholarly notes and unnumbered asides are
  different things in a 1970s academic book. — Project owner
- **2026-10-15** — Write a research doc before the speckit spec. — Follows
  CLAUDE.md's process for non-trivial features. — Project owner
- **2026-10-15** — PDF output (A5 reading edition plus an A4 imposed print
  edition, Typst engine) is a separate follow-on feature, researched in
  [005-pdf-output-quarterly.md](005-pdf-output-quarterly.md). — Keeps the
  site type shippable on its own. — Project owner
- **2026-10-15** — An issue can be marked `draft: true` in `issue.md`. A
  draft issue and its articles are left out of the built site, the index,
  `llms.txt`, category pages, the archive and the MCP tools, and never
  count as the current issue; `dev` still shows them. — Lets the author
  write the next issue in the same repo without publishing it. — Project
  owner (spec clarification)
- **2026-10-15** — The automatic pre-build step in `build`/`dev` runs
  settings validation and `index` for all four site types; regenerating
  `tokens.css` and converting images stay quarterly-only. Profile, log
  and docs pages, styles and assets are unchanged. — Removes the manual
  `index` step everywhere without restyling existing sites; making their
  pages read the tokens remains a separate decision. — Project owner
  (spec clarification)
- **2026-10-15** — Default faces: Archivo (body and display) and IBM Plex
  Mono, self-hosted woff2 with OFL texts. — Distinct from the Inter used by
  the profile and docs presets; Latin-Extended covers Czech and Polish. —
  Project owner (plan, [research R1](../004-quarterly-site-type/research.md#r1-default-text-face-roadmap-05))
- **2026-10-15** — Footnotes and margin notes are extracted in Go at index
  time; the index also carries a `markdown_body` with margin notes as
  footnotes. A template contract test keeps the TS renderer in agreement.
  (Confirms 005's proposal.) — Project owner (plan, R2, R6)
- **2026-10-15** — `tokens.css` stays variables only; structural switches
  (`margin_notes`, `numerals`) go into the content index for components to
  read. — No CSS style queries needed; one producer (Go). — Project owner
  (plan, R3)
- **2026-10-15** — Plates: keep JPEG/PNG format; WebP → PNG with alpha,
  else JPEG; max 1600 px wide; 40 MP safety cap; hash-named files in
  `public/plates/`; image links rewritten in Go. — Pure Go, cache-busting
  names, no path logic in TS. — Project owner (plan, R4, R5)

## Follow-up

**Picked up 2026-10-15.** The spec is written
([specs/004-quarterly-site-type/spec.md](../004-quarterly-site-type/spec.md)),
along with the feature item
([011-quarterly-site-type.md](../features/011-quarterly-site-type.md)) and
[roadmap §11](roadmap.md#11-quarterly-site-type). The remaining open
questions are tracked in [roadmap §0.5](roadmap.md#05-quarterly-site-type)
and settle during planning.

Original parking note (2026-10-15):

1. Resolve the Open Questions above. The ones blocking the spec are the
   default text face and whether the build pre-step applies to all types.
2. Review [005](005-pdf-output-quarterly.md) alongside this doc. It
   depends on this doc's content model, notes syntax, typesetting schema
   and pre-step.
3. Write the feature spec (`specs/features/`, from `TEMPLATE.md`), then run
   `/speckit-specify` for the target release.

## References

- [003-post-ga-feature-backlog.md](003-post-ga-feature-backlog.md) — the
  "fourth site type" backlog item this resolves
- [roadmap.md §3, §6](roadmap.md) — template and theming history, including
  the superseded Tailwind plan
- [specs/features/006-theming-system.md](../features/006-theming-system.md) —
  token contract intent
- `templates/site/log/`, `templates/site/docs/app/content/toc.ts` — starting
  points
- [CriticMarkup](https://github.com/CriticMarkup/CriticMarkup-toolkit) —
  origin of the `{>> … <<}` comment syntax
- [Tufte CSS](https://edwardtufte.github.io/tufte-css/) — prior art for
  sidenotes that collapse responsively
- [SIL Open Font License](https://openfontlicense.org/) — redistribution
  terms for bundled fonts

## Appendix

### Default preset sketch (`quarterly-default`)

```
color-bg      #fbfaf7      color-fg     #111111
color-accent  #d7261e      color-muted  #6b6b6b
color-rule    #111111      color-plate-bg #ececec
font-body     "Archivo", "Helvetica Neue", Arial, sans-serif
font-display  "Archivo", "Helvetica Neue", Arial, sans-serif
font-mono     "IBM Plex Mono", ui-monospace, monospace
```

### Contents page sketch

```
QUARTERLY                         No. 3 · Autumn 2026
──────────────────────────────────────────────────────
[ cover plate, b&w ]

CONTENTS

ESSAYS
  I.   On Margins ............ A. Novak
  II.  The Rigid Grid ........ J. Kral

DESIGN
  III. Type in Brno .......... M. Hora

                        Previous issues →
```
