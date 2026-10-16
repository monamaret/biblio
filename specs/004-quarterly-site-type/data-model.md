# Data Model: `quarterly` site type

**Feature**: [spec.md](spec.md) · **Research**: [research.md](research.md)

Three layers: what the author writes (`content/`), what `biblio-cli`
derives (`content-index.json`, `public/tokens.css`, `public/plates/`), and
what the template reads. The Go side owns every derived value
(Principle II); the template only reads the index.

## On-disk layout (authored)

```text
content/
  site.yaml                     site_type: quarterly, name, description, mcp_enabled, typesetting?
  colophon.md                   optional
  issues/
    <NN-slug>/                  issue folder; its name is the issue slug
      issue.md                  required
      <any>.jpg|.png|.webp      cover and plates, any depth
      <article>.md              zero or more articles
```

Any other `.md` file under `content/` (outside `issues/<dir>/` and other than
`colophon.md`) is a validation error naming the file. An `issue.md` nested
deeper than `issues/<dir>/` is also an error.

## Site settings (`content/site.yaml`)

| Field | Type | Required | Notes |
|---|---|---|---|
| `site_type` | `quarterly` | yes | unchanged |
| `name`, `description`, `mcp_enabled` | as today | — | unchanged |
| `typesetting` | map | no | quarterly only; an error on other types |

### Typesetting

Keys, types and defaults are as in the
[research doc's schema](../research/004-quarterly-site-type.md#typesetting-schema-siteyaml);
validation rules are in [research.md R9](research.md#r9-settings-validation-rules).
Unknown keys are errors.

**Effective typesetting** = `quarterly-default` preset values, then each
set override. Produced once by the Go resolver and written to two places:

- `public/tokens.css`: CSS custom properties ([contracts/config-and-tokens.md](contracts/config-and-tokens.md#tokenscss)).
- `content-index.json` → `typesetting`: only the structural switches
  `margin_notes` and `numerals` ([research.md R3](research.md#r3-ornament-and-structural-toggles-roadmap-05)).

## Entities

### Issue

Source: `issues/<dir>/issue.md`.

| Field | Source | Type | Rules |
|---|---|---|---|
| `title` | frontmatter | string | required, non-empty |
| `number` | frontmatter | integer ≥ 1 | required; unique across all issues, drafts included |
| `date` | frontmatter | `YYYY-MM-DD` | required |
| `season` | frontmatter | string | optional display label, e.g. "Autumn 2026"; if absent, the contents page shows the date |
| `cover` | frontmatter | relative path | optional; must resolve to an image inside the issue folder |
| `cover_alt` | frontmatter | string | required when `cover` is set |
| `draft` | frontmatter | bool | optional, default `false` |
| `kind` | derived | `"issue"` | |
| `slug` | derived | `issues/<dir>` | |
| `issue` | derived | `<dir>` | folder name |
| `cover_plate` | derived | URL | `/plates/<hash>.<ext>` when `cover` is set |
| `article_count` | derived | integer | published articles in the issue |

`season` is not in the original interview list. It is added because the
contents page shows "No. 3 · Autumn 2026" and deriving a season from a date
is locale-dependent. It is optional, so it changes no requirement.

### Article

Source: `issues/<dir>/<name>.md` (any name other than `issue.md`).

| Field | Source | Type | Rules |
|---|---|---|---|
| `title` | frontmatter | string | required |
| `category` | frontmatter | string | required, non-empty after trimming |
| `author` | frontmatter | string | optional |
| `dek` | frontmatter | string | optional |
| `image` | frontmatter | relative path | optional; must resolve to an image under `content/` |
| `image_alt` | frontmatter | string | required when `image` is set |
| `image_caption` | frontmatter | string | optional |
| `order` | frontmatter | integer | optional |
| `tags` | frontmatter | list of strings | optional |
| `kind` | derived | `"article"` | |
| `slug` | derived | `issues/<dir>/<name>` | |
| `issue`, `issue_number`, `issue_date` | derived | from the folder's `issue.md` | |
| `category_key` | derived | string | [research.md R10](research.md#r10-content-model-in-go) |
| `position` | derived | integer, 1-based | position in the issue's contents order; feeds the roman/arabic numeral |
| `plate` | derived | URL | when `image` is set |
| `body` | derived | Markdown | image destinations rewritten to plate URLs |
| `markdown_body` | derived | Markdown | `body` with margin notes rewritten as footnotes `[^mN]` |
| `footnotes` | derived | `Note[]` | in reference order |
| `margin_notes` | derived | `Note[]` | in reference order |

### Note

| Field | Type | Footnote | Margin note |
|---|---|---|---|
| `id` | string | `fn-<label>` | `m1`, `m2`, … in reference order |
| `label` | string | the author's label (`1`, `gerstner`) | — |
| `number` | integer | display number, 1-based in reference order | — |
| `text` | Markdown (inline) | definition text | note text |
| `ref_offset` | integer | code point offset of `[^label]` in `body` | code point offset of `{>>` in `body` |

### Category

Derived, site-level (`content-index.json` → `categories`).

| Field | Type | Rules |
|---|---|---|
| `key` | string | normalized; used in `/category/<key>` |
| `label` | string | first label seen, in issue order (newest first), then contents order |
| `article_count` | integer | published articles |

### Colophon

Source: `colophon.md`. Fields: `title` (required), `kind: "colophon"`,
`slug: "colophon"`, `body`.

### Plate

A converted image. Not an index entry; the mapping lives in the derived
URL fields above.

| Attribute | Rule |
|---|---|
| Source | an image file under `content/`, after symlinks are resolved, still inside `content/` |
| Accepted formats | JPEG, PNG, WebP (by content sniffing, not extension) |
| Limits | ≤ 40 MP by header; otherwise an error naming the file |
| Output | `public/plates/<hash16>.<jpg|png>` ([research.md R4](research.md#r4-image-output-roadmap-05)) |
| Lifecycle | created when missing; deleted when no source maps to it |

## `content-index.json` shape (quarterly)

Existing types keep `{ "entries": [...] }` exactly. Quarterly adds
top-level keys alongside `entries`:

```text
{
  "entries":     [ Issue | Article | Colophon ... ],   flat objects, as today
  "issues":      [ { issue, slug, number, title, date, season?, cover_plate?, cover_alt?, draft, article_count } ]   newest first
  "categories":  [ { key, label, article_count } ]
  "current_issue": "<issue slug>" | null
  "typesetting": { "margin_notes": bool, "numerals": "roman" | "arabic" }
}
```

`entries` order: issues by `number` descending; within an issue, `issue.md`
then articles in contents order; the colophon last.

## Ordering rules

- **Issues**: `number` descending.
- **Current issue**: the highest-numbered issue with `draft: false`, or
  `null` when there is none.
- **Articles in an issue**: `order` ascending, with unordered articles after
  ordered ones; ties broken by file name (byte order).
- **Categories in an issue's contents**: by the contents position of their
  first article.

## Drafts (state)

```text
issue.md draft: false|absent ──▶ published: in every output
issue.md draft: true         ──▶ draft:     build/index → excluded everywhere
                                            dev → included, with draft: true on the issue and its articles
```

There is no other state. Turning a draft into a published issue is
removing the field; the next build picks it up. A draft can never be the
current issue in any mode. Its `number` still counts for uniqueness.

## Validation summary (all errors collected, then reported together)

| Check | Message names |
|---|---|
| Missing or empty required field | file, field |
| Wrong type (`number` not an integer, bad date, `tags` not a list) | file, field, expected type |
| Duplicate issue `number` | both issue folders |
| `.md` file in an unexpected place | file |
| `cover`/`image` set without its `_alt` field | file, field |
| Image missing, corrupt, over cap, or outside `content/` | image path, referencing file |
| Unknown or invalid `typesetting` key | key, accepted values |
| `typesetting` on a non-quarterly site | site type |

Warnings (the build continues): an unreferenced footnote definition, a
reference without a definition, an unclosed `{>>`, and a configured font
family with no matching file in `public/fonts/`.
