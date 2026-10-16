# Quickstart: validating the `quarterly` site type

Run these scenarios end to end once implementation is done. Each one names
the user story and success criterion it proves. Details live in
[contracts/](contracts/) and [data-model.md](data-model.md).

## Prerequisites

- Go 1.24+, Node 20+, npm.
- From the repo root: `go build -o /tmp/biblio-cli ./cmd/biblio-cli`
  (or `go install ./cmd/biblio-cli`).
- The automated gate passes: `go test ./...` and `./scripts/check.sh`.
- A scratch directory outside the repo, e.g. `cd "$(mktemp -d)"`.

## 1. Scaffold, build and browse (US1 · SC-001, SC-002)

```sh
biblio-cli new field-notes --type=quarterly --domain=example.com --theme=quarterly-default
cd field-notes && npm install
biblio-cli build
```

**Expect:**
- The build output shows the validate, tokens, images and index lines
  ([contracts/cli.md](contracts/cli.md#biblio-cli-build-all-site-types)),
  then the prerender completes with no errors.
- `dist/client/` contains `index.html`, `issues/index.html`, one
  `issues/<issue>/index.html` per sample issue, one page per category,
  one per article, and `colophon/index.html`.

Then `biblio-cli dev` and open the local URL:
- `/` shows the newest sample issue's contents: masthead,
  "No. N · <season>", cover plate, categories in order, roman numerals,
  dotted leaders, authors.
- `/issues` lists every issue, newest first. A category link lists only
  that category's articles, each with its issue.
- `curl -H 'Accept: text/markdown' <local>/issues/<issue>/<article>`
  returns Markdown.

## 2. Notes (US2 · SC-005)

Open the sample article with dense notes (the template ships one).

- At a 1440 px viewport: footnotes (numbered) and margin notes
  (unnumbered) sit in the margin next to their references, with no
  overlap.
- At 375 px: footnotes are listed at the article end, and each margin note
  is a † that expands inline. Disable JavaScript and repeat: same result.
- `curl -H 'Accept: text/markdown' …`: margin notes appear as `[^m1]`…
  footnotes; count every note and compare with the source file.
- Set `typesetting: { margin_notes: false }`, restart `dev`: no margin
  column, all notes at the foot.

## 3. Typesetting (US3 · SC-003, SC-008)

Add to `content/site.yaml`:

```yaml
typesetting:
  accent: "#1e4fd7"
  align: justify
  numerals: arabic
```

Restart `biblio-cli dev`. **Expect:** blue accent, justified hyphenated
text, contents numbered 1, 2, 3. `public/tokens.css` lists the three
overrides in its header.

Then try each invalid value (one at a time) and run `biblio-cli build`:
`scale: 2`, `align: center`, `accent: red`, `font_body: "x;}"`,
`colour: "#000"`. **Expect:** non-zero exit before Vite runs, with a message
naming the key ([contracts/config-and-tokens.md](contracts/config-and-tokens.md#error-format)).
Set `font_body: "Söhne"` with no font file: the build succeeds and warns.

## 4. Plates (US4 · SC-004)

Copy a color JPEG, a PNG with transparency, and a WebP into the sample
issue folder, and reference them from an article (`image:` and an inline
`![…](img/…)`).

```sh
biblio-cli build     # expect "3 converted"
biblio-cli build     # expect "0 converted, … skipped"
```

**Expect:** every `<img>` in `dist/client` points at `/plates/…`; the
files are grayscale and at most 1600 px wide; the PNG keeps its
transparency. Replace one image with a truncated file: the build fails
naming it. With `grayscale_images: false`, the plates are in color.

## 5. Drafts (FR-012)

Add `draft: true` to the newest issue's `issue.md`.

- `biblio-cli build`: that issue is absent from `/`, `/issues`, category
  pages, `llms.txt` and `content-index.json`; `/` shows the previous
  issue.
- `biblio-cli dev`: the draft issue is visible.

## 6. Agent tools (US5 · SC-006)

With `biblio-cli dev` running, point an MCP client (e.g. MCP Inspector) at
`<local>/mcp`:

- `list_issues` → a table, newest first.
- `get_issue {"issue":"1"}` → contents grouped by category.
- `list_articles {"category":"essays"}` → only Essays; add `"author"` to
  narrow it.
- `get_article {"slug":"<article with notes>"}` → text with every note,
  plus `structuredContent.footnotes`/`marginNotes` with offsets.
- `get_issue {"issue":"99"}` → a "not found" text answer.

## 7. Other types unchanged (US6 · SC-007)

```sh
go test ./test/integration/...      # cross-type review (four types), update/upgrade, regressions
```

Manually: scaffold a `log` site, run `biblio-cli build` without ever
running `biblio-cli index`. It builds, and its `dist/` matches a build
from the previous release apart from the regenerated index.
