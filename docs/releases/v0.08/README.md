# biblio-cli — v0.08 User Guide (draft: the `quarterly` site type)

**Release**: v0.08 — *draft, not yet released*
**Status**: builds on [v0.07](../../../specs/releases/0.07-release.md) and
the [v0.06 guide](../v0.06/README.md); everything there still applies. This
release adds a fourth site type, `quarterly`, and makes `build`/`dev` index
content automatically for every type.

> This guide is drafted alongside the feature (specs/004-quarterly-site-type,
> tasks.md T097). The release package itself (`specs/releases/0.08-…/`) is
> assembled when v0.08 closes, per the constitution's Release Process.

## What's new

- **`quarterly`** — a magazine published in numbered issues, articles grouped
  by category, in the style of an early-1970s Eastern European academic
  book: black-and-white plates, footnotes and margin notes, one bold accent
  color, and typesetting you configure in `content/site.yaml`.
- **No more manual `biblio-cli index`.** `build` and `dev` now validate
  `content/site.yaml` and regenerate `content-index.json` before Vite runs,
  for every site type. Profile, log and docs sites otherwise build exactly
  as before.
- **Safer `update`/`upgrade`.** `content-index.json`, `public/tokens.css` and
  `public/plates/` are treated as build output, never as template files, so
  they can no longer cause update conflicts.

## The `quarterly` lifecycle

### Create

```sh
biblio-cli new field-notes --type=quarterly --domain=example.com --theme=quarterly-default
cd field-notes && npm install
```

(Omit any flag to be prompted; `quarterly` is the fourth choice.) The new
project has two sample issues, five sample articles, a colophon, self-hosted
Archivo and IBM Plex Mono with their OFL licenses, and a `.gitignore` that
keeps generated plates out of git.

### Configure

**Content layout:**

```text
content/
  site.yaml                 site_type: quarterly, optional typesetting:
  colophon.md               optional; title required
  issues/
    01-spring-2026/         one folder per issue; its name is the issue slug
      issue.md              title, number, date required; season, cover, cover_alt, draft optional
      cover.jpg
      on-margins.md         an article: title, category required;
                            author, dek, image, image_alt, image_caption, order, tags optional
      img/…
```

- The current issue (the home page) is the published issue with the
  highest `number`. Numbers must be unique.
- Articles appear in the contents by `order`, then file name, grouped by
  category in the order their first article appears.
- **Drafts:** add `draft: true` to an `issue.md` to keep that issue and its
  articles out of the built site, `llms.txt` and the MCP tools. `biblio-cli
  dev` still shows it, so you can read it while writing. Remove the line to
  publish.

**Notes:**

```markdown
The margins were wide.[^1] Every page left a third empty,{>> cf. Gerstner, 1964 <<} and …

[^1]: Brno, 1971.
```

`[^n]` is a numbered footnote; `{>> … <<}` is an unnumbered margin note.
On wide screens both sit in the margin beside their line; on narrow screens
footnotes collect at the end and margin notes become † toggles. Markdown and
terminal clients get margin notes as ordinary footnotes. Problems (a missing
definition, an unclosed `{>>`) are warnings in the build output, never
failures.

**Typesetting** — every key optional; omitted keys keep the
`quarterly-default` look:

```yaml
typesetting:
  accent: "#d7261e"          # #rrggbb
  font_body: Archivo         # also font_display, font_mono (family names)
  base_size: 1.0625rem       # 0.75–1.5rem or 12–24px
  scale: 1.25                # 1.1–1.618
  leading: 1.5               # 1.1–2.0
  measure: 34em              # 20–90 em or ch
  align: left                # left | justify (justify turns on hyphenation)
  grid: asymmetric           # asymmetric | centered
  margin_notes: true         # false puts every note at the article's end
  drop_caps: none            # none | drop | raised
  lead_in: small_caps        # none | small_caps
  numerals: roman            # roman | arabic
  rules: heavy               # none | hairline | heavy
  fleurons: false
  grayscale_images: true     # false keeps photos in color (still width-capped)
```

A change takes effect on the next `biblio-cli dev` start or `build`. Any
invalid value stops the build before anything is written, with a message
naming the key and the accepted values. A typeface without a matching file
in `public/fonts/` only warns, and pages fall back to the generic family.

**Images:** put JPEG, PNG or WebP files anywhere under `content/`. Each
build publishes them as grayscale plates no wider than 1600 px in
`public/plates/`, skipping images that haven't changed. Files over
40 megapixels, corrupt files, and symlinks pointing outside `content/` fail
the build with the file named.

### Preview

```sh
biblio-cli dev
```

Runs the pre-build step once (drafts included), then `vite dev` and
`wrangler dev`. Restart it to pick up content or typesetting changes.

### Update

```sh
biblio-cli update
```

As in [v0.06](../v0.06/README.md#update). Your `content/`, `content/site.yaml`
and the generated `content-index.json`, `public/tokens.css` and
`public/plates/` are never touched or reported as conflicts.

### Deploy

```sh
biblio-cli build    # validate → tokens → plates → index → prerender
biblio-cli deploy
```

The generated `.github/workflows/deploy.yml` runs the same two commands on
push.

### Destroy

```sh
npx wrangler delete   # from the project directory
```

As in [v0.06](../v0.06/README.md#destroy): your project directory and
content are untouched.

### What an AI assistant can do

Besides the shared `list_pages`, `get_page` and `search_content`, a
quarterly site's `/mcp` endpoint offers:

| Tool | Answers |
|---|---|
| `list_issues` | every published issue: number, title, date, article count, slug |
| `get_issue` | an issue's contents grouped by category (by slug or number) |
| `list_articles` | articles filtered by issue, category, author, tag and issue-date range |
| `get_article` | an article's full text with every note, plus the footnotes and margin notes as structured data |

Every answer is readable as plain text; drafts are never visible.

## Known scope limits in v0.08

- PDF output for quarterly issues and articles is a separate follow-on
  feature ([research 005](../../../specs/research/005-pdf-output-quarterly.md)).
- `dev` does not watch for content changes; restart it.
