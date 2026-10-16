# Feature Specification: `quarterly` site type

**Feature Branch**: `monamaret/004-quarterly-site-type`

**Created**: 2026-10-15

**Status**: Draft

**Feature item**: [011-quarterly-site-type.md](../features/011-quarterly-site-type.md) · **Roadmap**: [§11](../research/roadmap.md#11-quarterly-site-type)

**Input**: User description: "004" — the `quarterly` site type researched in [specs/research/004-quarterly-site-type.md](../research/004-quarterly-site-type.md). PDF output for quarterly sites ([005](../research/005-pdf-output-quarterly.md)) is a separate follow-on feature and is out of scope here.

## Summary

The author can scaffold, write, and deploy a fourth site type, `quarterly`. It is a magazine published in numbered issues, with articles grouped by category. It looks like an early-1970s Eastern European academic book: black-and-white plates, generous white space, footnotes and margin notes, one bold accent color, and typesetting the author can configure. There are no comments. A quarterly site keeps every guarantee the profile, log, and docs types already give: prerendered pages, HTML/Markdown/terminal views of the same URL, `llms.txt`, an agent (MCP) endpoint, the deploy workflow, and `update`/`upgrade` support.

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Publish a numbered issue (Priority: P1)

As the author, I scaffold a quarterly site, write an issue (an issue file plus a few articles, each with a category), build it, and deploy it. Readers land on the current issue's contents page, can open any article, browse past issues, browse a category across issues, and read the colophon.

**Why this priority**: This is the site type itself. Without it nothing else in this feature has anywhere to live, and it alone is a usable, publishable magazine.

**Independent Test**: Scaffold a new quarterly site, add two issues with three articles across two categories plus a colophon, build and deploy, then visit every route and confirm each page shows the right content in the right order.

**Acceptance Scenarios**:

1. **Given** no existing site, **When** the author runs the create command with the quarterly type, **Then** a working quarterly site is produced that builds and serves out of the box, with a sample issue, sample articles, and a colophon.
2. **Given** a site with issues No. 1 and No. 2, **When** a reader opens the home page, **Then** they see issue No. 2's contents page: masthead, issue number and season, cover plate, and the articles grouped by category with numbering, dotted leaders, and author names.
3. **Given** a site with several issues, **When** a reader opens the issue archive, **Then** every issue is listed newest first, and each links to that issue's contents page.
4. **Given** articles in the category "Essays" across two issues, **When** a reader opens that category's page, **Then** exactly those articles are listed, each showing which issue it belongs to.
5. **Given** an issue whose articles have `order` values, **When** the contents page renders, **Then** articles appear in ascending `order` (ties broken by file name), and categories appear in the order of their first article.
6. **Given** a deployed quarterly site, **When** a client requests an article URL as Markdown or as terminal text, **Then** it receives that representation of the article, just as with the other site types.

---

### User Story 2 - Scholarly notes (Priority: P2)

As the author, I write numbered footnotes and unnumbered margin notes in an article using two distinct syntaxes. Readers on a wide screen see both beside the line that references them. On a narrow screen, footnotes collect at the end of the article and margin notes become small toggles marked with a dagger (†) that open inline.

**Why this priority**: Notes are the defining trait of the academic-book look and the main thing that separates this type from the log type. Without them it is just a blog in issues.

**Independent Test**: Publish one article with several footnotes and margin notes, some long. Check its layout at a wide and a narrow viewport, then fetch it as Markdown and as terminal text.

**Acceptance Scenarios**:

1. **Given** an article with footnotes `[^1]`, `[^2]` and a margin note `{>> … <<}`, **When** it is viewed on a wide screen, **Then** each note appears in the margin next to its reference, footnotes numbered and the margin note unnumbered.
2. **Given** the same article, **When** it is viewed on a narrow screen, **Then** footnotes appear in a list at the end of the article, and each margin note is a † toggle that reveals its text inline. Both work without client-side scripting.
3. **Given** the same article, **When** it is requested as Markdown, **Then** footnotes are unchanged and margin notes have become ordinary footnotes, so the file reads as plain Markdown with no lost text.
4. **Given** the same article, **When** it is requested as terminal text, **Then** the notes appear in the same form as in the Markdown view.
5. **Given** an author who has turned margin notes off in the site settings, **When** the article is viewed on a wide screen, **Then** the text column widens and notes render at the end of the article.

---

### User Story 3 - Configure the typesetting (Priority: P3)

As the author, I get a finished default look without touching anything: a grotesk typeface, a red accent, an off-white page, an asymmetric grid, ragged-right text. I can override any of its typesetting values (typefaces, size and scale, leading, measure, alignment, grid, accent color, and book ornaments such as drop caps, small-caps lead-ins, numeral style, rules, and fleurons) in one block of the site's settings file, and see the change on the next preview or build.

**Why this priority**: Configurable typesetting was an explicit requirement, but the site is already publishable with the defaults, so this layers on top of Stories 1 and 2.

**Independent Test**: Change the accent color, alignment, and numeral style in the settings file, rebuild, and confirm the rendered pages change. Then enter an invalid value and confirm the build stops with a message naming the key.

**Acceptance Scenarios**:

1. **Given** a freshly scaffolded site with no typesetting overrides, **When** it is built, **Then** pages render with the `quarterly-default` look, using typefaces served by the site itself and not from a third-party font service.
2. **Given** an override `accent: "#1e4fd7"`, **When** the author runs the preview or build, **Then** every accent-colored element uses the new color.
3. **Given** `align: justify`, **When** the site is rebuilt, **Then** body text is justified with hyphenation on.
4. **Given** `numerals: arabic`, **When** the contents page renders, **Then** articles are numbered 1, 2, 3 instead of I, II, III.
5. **Given** an invalid value (an unknown enum, a malformed color, a scale outside 1.1–1.618, or a font name with disallowed characters), **When** the author builds, **Then** the build fails before producing any output, with a message naming the offending key and the accepted values.
6. **Given** a typeface override naming a family whose font files are not in the site, **When** the author builds, **Then** the build completes, warns that the family has no matching files, and pages fall back to the generic family.

---

### User Story 4 - Black-and-white plates (Priority: P4)

As the author, I drop ordinary color photos into an issue. The published site shows them as grayscale plates, scaled down to a sensible maximum width, without my editing any image by hand.

**Why this priority**: The b&w plates are part of the look, but the site is still correct with color images, so this is an enhancement on top of Story 1.

**Independent Test**: Add a color JPEG, a PNG with transparency, and a WebP to an issue. Build, then confirm each is served grayscale, width-capped, with transparency kept, and that a rebuild with no image changes skips the conversion.

**Acceptance Scenarios**:

1. **Given** color images in an issue, **When** the site is built, **Then** every image served for the site's content is grayscale and no wider than the configured maximum.
2. **Given** a PNG with transparency, **When** it is converted, **Then** its transparency is preserved.
3. **Given** no images have changed since the last build, **When** the author rebuilds, **Then** no image is reconverted, and the build reports how many images it converted and skipped.
4. **Given** `grayscale_images: false`, **When** the site is built, **Then** images are served in their original color (still width-capped).
5. **Given** an image whose declared dimensions exceed the safety cap, or a corrupt image file, **When** the site is built, **Then** the build fails with a message naming the file, and the image is never fully decoded.

---

### User Story 5 - Agents read the magazine (Priority: P5)

As an MCP client or a coding assistant in a terminal, I can list issues, read an issue's contents, filter articles, and fetch an article with its notes as structured data, each answer readable as plain text.

**Why this priority**: Every site type has an agent surface, so this is required for parity. It builds on Stories 1 and 2 and is tested last.

**Independent Test**: Against a built fixture site, call each quarterly tool and the shared tools, and check that the text answers and structured fields are correct.

**Acceptance Scenarios**:

1. **Given** a site with three issues, **When** an agent calls `list_issues`, **Then** it gets each issue's number, title, date, slug, and article count, newest first, as a readable table.
2. **Given** an issue number or slug, **When** an agent calls `get_issue`, **Then** it gets the issue's metadata and its contents grouped by category, in the same order as the contents page.
3. **Given** filters for issue, category, author, tag, and a date range, alone or combined, **When** an agent calls `list_articles`, **Then** it gets only the articles matching all given filters (title, dek, author, category, issue, slug).
4. **Given** an article with footnotes and margin notes, **When** an agent calls `get_article`, **Then** it gets the full text plus separate lists of footnotes and margin notes, each with an identifier, its text, and where in the body it is referenced.
5. **Given** a quarterly site, **When** an agent calls the shared tools (`list_pages`, `get_page`, `search_content`), **Then** they behave as they do on the other site types.
6. **Given** a request for an unknown issue or article, **When** the tool runs, **Then** it returns a clear "not found" text answer, not an error page or an empty result.

---

### User Story 6 - Keep it in step with the other types (Priority: P6)

As the author, I can `update` and `upgrade` a quarterly site like any other site, my existing profile, log, and docs sites behave exactly as before, and the cross-type consistency review passes with all four types.

**Why this priority**: This is the release gate, not a new capability, but the feature cannot ship without it.

**Independent Test**: Run the full existing test suite with no changes to profile, log, or docs tests. Run `update`/`upgrade` against a quarterly site scaffolded from the previous template revision. Run the cross-type consistency review.

**Acceptance Scenarios**:

1. **Given** existing profile, log, and docs sites, **When** they are built with the new CLI, **Then** their pages, styles, and assets are unchanged, and their content index is regenerated automatically without a separate `index` run.
2. **Given** a quarterly site scaffolded from an older template revision, **When** the author runs `upgrade`, **Then** template files are updated and the author's content and settings are untouched.
3. **Given** the guided (interactive) create flow, **When** the author is asked for a site type, **Then** `quarterly` is offered alongside the other three.

---

### Edge Cases

- An issue folder has no articles: the issue still appears in the archive and its contents page says it has no articles yet.
- Two issues declare the same `number`: the build fails, naming both issue folders.
- An `issue.md` is missing `title`, `number`, or `date`, or an article is missing `title` or `category`: the build fails, naming the file and the missing field.
- An article has no `order`: it sorts after ordered articles, by file name.
- A category label differs only by case or spacing ("Essays" / "essays "): both are treated as one category, with one stable URL.
- A footnote reference has no matching definition, or a definition is never referenced: the build warns, naming the article and the note.
- A margin note contains Markdown (emphasis, a link): it renders as formatted text; raw HTML in a note is not rendered.
- A margin note marker is never closed: the text is left as-is and the build warns.
- A very long margin note on a wide screen: it does not overlap the next note or run off the text block.
- No `colophon.md`: the colophon route returns not-found and no link to it appears.
- No issues at all (a brand-new site with the sample issue removed): the home page shows an empty-state message instead of failing.
- Every issue is a draft: the built site behaves as if there are no issues, while preview shows the drafts.
- The highest-numbered issue is a draft: the home page shows the highest-numbered published issue.
- An issue's `number` duplicates a draft issue's `number`: the build still fails, since numbers must be unique whether or not an issue is a draft.
- An image referenced from content is outside `content/` or reached through a symbolic link that leaves `content/`: it is rejected.

## Requirements *(mandatory)*

### Functional Requirements

**Scaffolding and lifecycle**

- **FR-001**: The CLI MUST offer `quarterly` as a fourth site type in both the flag-driven and guided create flows, producing a site that builds, previews, and deploys with no further setup.
- **FR-002**: A scaffolded quarterly site MUST include a sample issue with at least two articles in two categories, a cover image, and a colophon, and MUST include the licenses of every bundled typeface.
- **FR-003**: Quarterly sites MUST support `update` and `upgrade` with the same guarantees as the other types: template files are updated, and author content and settings are never overwritten.
- **FR-004**: Quarterly sites MUST provide every cross-type guarantee: prerendered pages, HTML/Markdown/terminal content negotiation on every page URL, `llms.txt`, the agent endpoint (switchable by the existing single on/off setting), and the deploy workflow.

**Content model**

- **FR-005**: An issue MUST be a folder under `content/issues/` containing an `issue.md` with required `title`, `number`, and `date`, and optional `season` (a display label such as "Autumn 2026"; when absent, pages show the date), `cover`, `cover_alt` (required when `cover` is set), and `draft` (see FR-012).
- **FR-006**: An article MUST be a Markdown file inside an issue folder with required `title` and `category`, and optional `author`, `dek`, `image`, `image_alt`, `image_caption`, `order`, and `tags`. An article's issue MUST come from its folder, not from a frontmatter field.
- **FR-007**: The site MAY have a `colophon.md` with a required `title`.
- **FR-008**: Required fields MUST be validated per kind of file (issue, article, colophon), and a validation failure MUST name the file and the field.
- **FR-009**: The current issue MUST be the published (non-draft) issue with the highest `number`. Issue numbers MUST be unique within a site.
- **FR-010**: Articles within an issue MUST be ordered by `order` ascending, then by file name; categories within an issue MUST be ordered by the first appearance of an article in that category.
- **FR-011**: Categories MUST be free-text labels with no predefined list; labels that differ only in case or surrounding whitespace MUST be treated as the same category.
- **FR-012**: An issue MAY be marked `draft: true` in `issue.md`. A draft issue and its articles MUST be left out of the built site, the content index, `llms.txt`, category pages, the archive, and agent tools, and MUST NOT count when choosing the current issue. During preview, draft issues MUST still be shown so the author can read them while writing.

**Pages**

- **FR-013**: The site MUST serve these pages: the home page (the current issue's contents), an issue archive (newest first), each issue's contents page, one page per category (articles across all issues), each article, and the colophon when present.
- **FR-014**: A contents page MUST show the masthead, issue number and season/date, the cover plate when present, and articles grouped by category with numbering (roman or arabic per settings), dotted leaders, and author names.
- **FR-015**: The site MUST NOT offer comments or author pages.

**Notes**

- **FR-016**: Authors MUST be able to write numbered footnotes with `[^label]` syntax and unnumbered margin notes with `{>> text <<}` syntax in the same article.
- **FR-017**: On wide viewports, both kinds of note MUST appear in the margin beside the line that references them, without overlapping one another.
- **FR-018**: On narrow viewports, footnotes MUST collect at the article's end and margin notes MUST become † toggles that open inline, working without client-side scripting.
- **FR-019**: In the Markdown and terminal views, footnotes MUST be kept as footnotes and margin notes MUST be rewritten as footnotes, so no note text is lost.
- **FR-020**: Note text MUST support inline Markdown formatting and MUST NOT render raw HTML.

**Typesetting**

- **FR-021**: The type MUST ship a `quarterly-default` look: a self-hosted grotesk text face, a monospaced face for folios and labels, accent `#d7261e`, background `#fbfaf7`, text `#111`, asymmetric grid, ragged-right text, heavy rules, small-caps lead-ins, and roman numerals.
- **FR-022**: The author MUST be able to override, in an optional `typesetting:` block of the site settings file, each of: `accent`, `font_body`, `font_display`, `font_mono`, `base_size`, `scale`, `leading`, `measure`, `align`, `grid`, `margin_notes`, `drop_caps`, `lead_in`, `numerals`, `rules`, `fleurons`, and `grayscale_images`, with the value types and defaults listed in the research doc's typesetting schema.
- **FR-023**: Every typesetting value MUST be validated before any output is produced; an invalid value MUST fail the build with a message naming the key and the accepted values. No free-form author text may reach the page styles unvalidated.
- **FR-024**: A change to any typesetting value MUST take effect on the next preview start or build, with no manual regeneration step.
- **FR-025**: When a typeface override names a family with no matching font files in the site, the build MUST warn and pages MUST fall back to the generic family.
- **FR-026**: The page styles MUST take every color and typeface from the generated design values, never from values hard-coded in the template, so all typesetting stays controlled from one place.
- **FR-027**: Fonts MUST be served by the site itself, never from a third-party font service.

**Images**

- **FR-028**: At build time, every image under `content/` MUST be published as a grayscale copy no wider than a fixed maximum width, unless `grayscale_images` is `false`, in which case it is only width-capped.
- **FR-029**: Image conversion MUST preserve transparency, accept JPEG, PNG, and WebP sources, and skip images unchanged since the last build.
- **FR-030**: The build MUST reject, with a message naming the file, any image that is corrupt, exceeds a pixel-count safety cap (checked before full decoding), or resolves outside `content/`.
- **FR-031**: The build MUST report how many images were converted and how many skipped.

**Agent tools**

- **FR-032**: The agent endpoint MUST offer, in addition to the shared tools, `list_issues`, `get_issue` (by number or slug), `list_articles` (filterable by issue, category, author, tag, and date range, combined as "all must match"), and `get_article` (by slug).
- **FR-033**: `get_article` MUST return footnotes and margin notes as separate structured lists, each note carrying an identifier, its text, and its reference position in the body.
- **FR-034**: Every quarterly tool MUST return a self-sufficient plain-text or Markdown answer (a table or list) as its primary content, and MUST return a clear "not found" answer for unknown issues or articles.

**Build pipeline scope**

- **FR-035**: Before the site is bundled, the build and preview commands MUST run a pre-build step. For every site type it MUST validate the site settings and regenerate the content index, so the author never has to run `index` by hand. For quarterly sites it MUST also regenerate the design values and convert images. The order is: validate settings, regenerate design values (quarterly), convert images (quarterly), regenerate the content index.
- **FR-036**: Existing profile, log, and docs sites MUST produce the same pages, styles, and assets as before this feature. The only change for them is that settings validation and index regeneration now run automatically on build and preview; a settings error that previously surfaced only from `index` now stops `build` or preview with the same message.

### Key Entities *(include if feature involves data)*

- **Issue**: One numbered publication. Has a number (unique), title, date, an optional season label, optional cover plate and its alt text, an optional draft flag, and a slug from its folder name. Contains zero or more articles. The highest-numbered published (non-draft) issue is the current one.
- **Article**: One piece of writing inside exactly one issue. Has title, category, and optionally author, dek, lead image with alt text and caption, order, and tags. Derives its issue slug, number, and date from its folder. Contains zero or more footnotes and margin notes.
- **Category**: A free-text label shared by articles across issues. Has a display label and a normalized form used for its URL and for grouping.
- **Footnote**: A numbered note in an article. Has a label, its text, and the position of its reference in the body.
- **Margin note**: An unnumbered aside in an article. Has a generated identifier, its text, and its position in the body.
- **Colophon**: A single optional page describing the publication (typefaces, masthead, edition).
- **Typesetting settings**: The effective look of the site: the `quarterly-default` preset merged with the author's validated overrides. Feeds the generated design values for every page.
- **Plate**: The published, converted copy of a content image (grayscale unless turned off, width-capped), linked to its source image so content can refer to the original path.

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: An author can go from an empty directory to a deployed quarterly site with one published issue of three articles in under 30 minutes, using only the CLI and a text editor.
- **SC-002**: A freshly scaffolded quarterly site builds and prerenders 100% of its issue, article, category, archive, and colophon pages with zero errors.
- **SC-003**: Changing any one typesetting value is visible in the rendered site on the very next preview or build, for 100% of the settings in the typesetting schema.
- **SC-004**: 100% of content images in a built site are served grayscale (with the default settings), and a rebuild with no image changes reconverts zero images.
- **SC-005**: On a fixture article with at least ten notes, at viewport widths from 320 px to 1920 px, no note overlaps another note or the body text, and no note text is missing from the HTML, Markdown, or terminal view.
- **SC-006**: Every quarterly agent tool returns a correct, readable text answer for each of the fixture site's issues and articles, and `get_article` returns every note in the fixture article with its correct reference position.
- **SC-007**: 100% of the existing profile, log, and docs tests pass unchanged, and the cross-type consistency review passes with all four types.
- **SC-008**: Every invalid typesetting value in a test matrix covering each key is rejected before any output is written, with a message that names the key.

## Assumptions

- **Scope:** PDF output (issue PDFs, offprints, the A4 print edition, the `pdf` command) is the separate follow-on feature in [005](../research/005-pdf-output-quarterly.md). This feature must not block it: notes are parsed once into structured data that a later PDF renderer can reuse, and design values are produced from one place.
- **Default text face:** Archivo, as in the research doc's preset sketch. The research doc's open question (Archivo or Inter Tight) can still swap it during planning after a side-by-side sample; this does not change any requirement.
- **Image limits:** the maximum published width defaults to 1600 px and the pixel-count safety cap to 40 megapixels. Output keeps the source format (JPEG stays JPEG, PNG stays PNG); WebP sources are written as PNG if they have transparency, otherwise JPEG. Planning may tune these numbers.
- **Preview loop:** if rerunning the pre-build steps on every content change during preview is too costly for the first cut, rerunning them when the preview starts is acceptable (per the research doc).
- **Clarified 2026-10-15 (project owner):** issues can be drafts, hidden from everything published but shown in preview (FR-012). The automatic pre-build step covers settings validation and indexing for all four types; design values and image conversion stay quarterly-only (FR-035, FR-036). Both answers should be appended to the 004 research doc's Decision Log, and 005's "Drafts" question now follows from FR-012.
- **Single author:** one author/editor per site, as for the other types. No accounts, roles, or comments.
- **Security posture:** the agent endpoint stays public with no auth, per the constitution.
- **Release:** ships as a new site type in a feature release (v0.08 per the research doc), with no migration needed for existing sites.
- **Process pipeline:** feature item [011-quarterly-site-type.md](../features/011-quarterly-site-type.md), [roadmap §11](../research/roadmap.md#11-quarterly-site-type) (with open planning items in §0.5), and the architecture sync ([architecture.md §2](../research/architecture.md#2-build--deployment-pipeline), pre-build step) were added on 2026-10-15, as the constitution's research-to-Spec-Kit pipeline requires.
