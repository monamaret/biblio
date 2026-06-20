# biblio-cli — v0.04 User Guide

**Release**: Prove the platform on the technical project log
**Status**: builds on [v0.03](../v0.03/README.md) — the same CLI, content
pipeline, theming, MCP server, and terminal/negotiation layer now also work
for the `log` site type. No new `biblio-cli` commands or flags.

## What's new

A second, structurally different site type — a running, cross-cutting
timeline of dated, tagged entries — is fully wired up to the same platform
v0.01–v0.03 proved on the `profile` pilot.

```sh
biblio-cli new my-project-log --type=log --domain=example.com --theme=log-default
cd my-project-log
# author dated, tagged entries under content/posts/ (date and title required;
# tags and project are optional)
biblio-cli index && biblio-cli build && biblio-cli deploy
```

### Entry list, detail view, and filtering

The deployed site shows every entry's date, tags, and project, newest
first, with date-range/tag/project filter controls. Filtering reads tag/
project values directly off the entries — a brand-new tag or project is
filterable immediately after the next `biblio-cli index`/`build`, with no
separate registration step. Combining filter dimensions (e.g. a date range
*and* a tag) returns the intersection.

### The daily authoring rhythm

Adding a new day's entry is the same flow every time: author one
Markdown+frontmatter file under `content/posts/`, then
`biblio-cli index && biblio-cli build && biblio-cli deploy` — no
per-entry structural or config changes, however many entries you add.

### What an AI assistant can do

Connect an MCP client to `https://<your-domain>/mcp` and it can list and
call:

- `list_pages`, `get_page`, `search_content` — shared, work on any site type
- `list_entries`, `get_entry` — log-site-specific; `list_entries` accepts
  optional `dateFrom`/`dateTo`/`tag`/`project` arguments and returns the
  same intersection semantics as the browser filter UI

A profile-shaped query (e.g. `get_profile`) against a log site returns the
MCP SDK's standard "unknown tool" error — that tool is simply never
registered on this site type.

### What a terminal or `curl` user sees

```sh
curl -H "Accept: text/markdown" https://<your-domain>/posts/2026-06-01-hello-world
curl -A "curl/8.0" -H "Accept: text/markdown" https://<your-domain>/posts/2026-06-01-hello-world
```

Same content-negotiated behavior as v0.03's `profile` site — the same URL,
no separate address, browsers unaffected.

## Known scope limits in v0.04

- Wired up for the `log` site type only — `docs` and the cross-type
  consistency check are v0.05.
- Filtering is computed client-side over the bundled `content-index.json`
  — no server-side query API beyond the MCP tools above (Principle V: no
  storage escalation triggered by adding a second site type).
- No authentication on `/mcp`, by design, same as every other site type.
