# biblio-cli — v0.05 User Guide

**Release**: Prove the platform on documentation and across all site types
**Status**: builds on [v0.04](../v0.04/README.md) — the same CLI, content
pipeline, theming, MCP server, and terminal/negotiation layer now also work
for the `docs` site type. No new `biblio-cli` commands or flags.

## What's new

The third and final supported site type — a hierarchy-based documentation
site, organized by table-of-contents `section`/`order` rather than a
timeline or a flat list — is fully wired up to the same platform v0.01–v0.04
proved on `profile` and `log`. With all three types now live, this release
also closes a cross-type consistency review: confirming theming and
terminal/CLI rendering both work correctly and consistently across all
three side by side.

```sh
biblio-cli new my-docs-site --type=docs --domain=example.com --theme=docs-default
cd my-docs-site
# author pages with `section`/`order` frontmatter under content/
biblio-cli index && biblio-cli build && biblio-cli deploy
```

### Table of contents, hierarchy navigation, and the deliberate absence of filters

The deployed site groups every page under its section heading, ordered
within each section, with every page reachable through the table of
contents alone. Unlike the `log` site type, there is **no** tag/project
filter UI or query surface anywhere — that absence is this site type's
confirmed, deliberate design, not an unfinished feature (verified by a
dedicated test asserting the negative).

### What an AI assistant can do

Connect an MCP client to `https://<your-domain>/mcp` and it can list and
call:

- `list_pages`, `get_page`, `search_content` — shared, work on any site type
- `list_sections`, `get_doc_page` — docs-site-specific; `list_sections`
  lists every section and its pages in hierarchy order

A log-shaped query (e.g. a tag-filter argument) against `list_sections`
is rejected with a clear, structured error — never silently ignored. A
profile/log-only tool (e.g. `get_profile`, `list_entries`) against a docs
site returns the MCP SDK's standard "unknown tool" error — those tools are
simply never registered on this site type.

### What a terminal or `curl` user sees

```sh
curl -H "Accept: text/markdown" https://<your-domain>/getting-started
curl -A "curl/8.0" -H "Accept: text/markdown" https://<your-domain>/getting-started
```

Same content-negotiated behavior as every other site type — the same URL,
no separate address, browsers unaffected.

### The cross-type consistency review

With `profile`, `log`, and `docs` all live, this release closes a
structured side-by-side review (see
`test/integration/cross-type-review/CONSISTENCY-RECORD.md`): each site
renders its own distinct theme with no cross-type bleed-through, and all
three share exactly one negotiation/rendering code path — no per-type
fork — with content appropriate to each one's own shape. Zero open
inconsistencies were found.

## Known scope limits in v0.05

- No update/upgrade-flow work — that's [v0.06](../v0.06/README.md).
- Filtering remains absent by design for the `docs` site type; this is a
  confirmed decision, not a gap to be filled in a later release.
- No authentication on `/mcp`, by design, same as every other site type.
