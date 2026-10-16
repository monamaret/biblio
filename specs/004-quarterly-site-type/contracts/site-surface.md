# Contract: site surface (routes, views, discovery)

## Routes

| Path | Page | 404 when |
|---|---|---|
| `/` | Contents page of the current issue; an empty-state message when there is no published issue | never |
| `/issues` | Archive: every published issue, newest first: number, title, season or date, article count, link | never |
| `/issues/<issue>` | That issue's contents page (same component as `/`) | unknown or draft issue (in build) |
| `/category/<key>` | Articles in the category across issues, newest issue first, then contents order; each shows its issue | unknown key |
| `/issues/<issue>/<article>` | Article | unknown, or in a draft issue (in build) |
| `/colophon` | Colophon | no `colophon.md` |

All pages are prerendered. No route exists for authors or comments
(FR-015). The masthead links to `/issues`, and to `/colophon` only when it
exists.

## Contents page

Top to bottom: masthead (site name), "No. <number> · <season or date>", the
cover plate with `cover_alt` when present, then "Contents". Under it, for
each category in contents order: the category label (accent color), then
its articles, each showing the numeral (by `--ts-numerals`; text views use
`typesetting.numerals`), the title linking to the article, a dotted leader,
and the author when present. Under the list: a link to the previous issues
(`/issues`).

## Article page

Category label (accent), title, author, dek, the lead plate with
`image_alt` and "Fig. 1 — <image_caption>", then the body. Notes follow
[research.md R7](../research.md#r7-notes-in-html-without-client-js):

- Wide viewports (`min-width: 64rem`) with `margin_notes: true`: footnotes
  and margin notes sit in the margin column next to their reference;
  footnotes are numbered and margin notes are not. The foot list is hidden.
- Narrow viewports: footnotes are in a list at the article foot; each
  margin note is a "†" toggle that reveals its text inline. No client JS.
- `margin_notes: false`: no margin column; both kinds of note are listed at
  the foot (footnotes numbered, margin notes marked †).
- Note text renders inline Markdown; raw HTML is never rendered.

## Content negotiation (unchanged mechanism)

| Request | Article response | Other pages |
|---|---|---|
| `Accept: text/html` | prerendered HTML | prerendered HTML |
| `Accept: text/markdown` or `.md` sibling | `markdown_body` (margin notes as `[^mN]` footnotes; plate URLs) | Markdown built from the index (contents list, archive, category list) |
| terminal user agents / ANSI | ANSI rendering of the same Markdown | same |

`Vary: Accept, User-Agent` and `Link: rel="alternate"` as in the existing
templates.

## `llms.txt`

Site name and description, then:

```text
## Current issue
- No. 3 · Autumn 2026: /issues/03-autumn-2026

## Issues
- No. 3 · Autumn 2026: /issues/03-autumn-2026
- No. 2 · Summer 2026: /issues/02-summer-2026

## Articles (current issue)
- On Margins — A. Novak: /issues/03-autumn-2026/on-margins

## About
- Colophon: /colophon
```

Drafts never appear in it.

## Images

Every `<img>` for content uses a `/plates/…` URL. No `src` points into
`content/`.
