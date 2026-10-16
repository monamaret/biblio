# Contract: MCP tools (quarterly)

Registered next to the shared tools `list_pages`, `get_page` and
`search_content`, whose behavior is unchanged. Every tool returns a
self-sufficient text `content` item (Principle IV). Drafts are never
visible. Unknown inputs return a text "not found" result with
`isError: false`, never a thrown error (FR-034).

## `list_issues`

- **Input:** none.
- **Output (text):**

```text
| No. | Title | Date | Articles | Slug |
|---|---|---|---|---|
| 3 | Autumn 2026 | 2026-10-01 | 5 | 03-autumn-2026 |
| 2 | Summer 2026 | 2026-07-01 | 4 | 02-summer-2026 |
```

- Empty: `No issues have been published yet.`

## `get_issue`

- **Input:** `{ issue: string }`: the folder slug (`03-autumn-2026`) or the
  number as a string (`"3"`).
- **Output (text):** heading `# No. 3 · Autumn 2026 — <title>`, then for
  each category in contents order a `## <Category>` heading and lines
  `<numeral>. <title> — <author> (/issues/03-autumn-2026/on-margins)`.
  Numerals follow `typesetting.numerals`.
- Not found: `No issue matches "7". Use list_issues to see published issues.`

## `list_articles`

- **Input** (all optional, combined with AND):

| Field | Type | Match |
|---|---|---|
| `issue` | string | slug or number |
| `category` | string | by normalized key, so "Essays" and "essays" match |
| `author` | string | case-insensitive exact |
| `tag` | string | exact |
| `dateFrom`, `dateTo` | `YYYY-MM-DD` | inclusive, on the issue date |

- **Output (text):** one line per article, newest issue first, then
  contents order:
  `- <title> — <author> · <category> · No. <n> — <dek> (/<slug>)`.
- Empty: `No articles matched the given filter.`
- An invalid date returns a validation message naming the field.

## `get_article`

- **Input:** `{ slug: string }`: the article slug with or without the
  leading `/`; a bare `<issue>/<article>` is also accepted.
- **Output:**
  - `content[0]` (text): `# <title>`, a byline line, then
    `markdown_body` (margin notes already rewritten as footnotes), so a
    terminal host gets every note.
  - `structuredContent`:

```json
{
  "slug": "issues/03-autumn-2026/on-margins",
  "title": "On Margins",
  "author": "A. Novak",
  "category": "Essays",
  "issue": { "slug": "03-autumn-2026", "number": 3, "date": "2026-10-01" },
  "footnotes":   [ { "id": "fn-1", "label": "1", "number": 1, "text": "Brno, 1971.", "ref_offset": 412 } ],
  "marginNotes": [ { "id": "m1", "text": "cf. Gerstner, 1964", "ref_offset": 233 } ]
}
```

  - The tool declares an `outputSchema` matching this shape.
- Not found: `No article matches "<slug>". Use list_articles to find one.`

## Unchanged

The MCP endpoint, the `mcp_enabled` switch, no auth, and the
Durable-Object session are unchanged (constitution, Security & Operational
Posture).
