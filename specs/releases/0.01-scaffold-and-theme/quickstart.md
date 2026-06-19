# Quickstart Validation: v0.01 — Scaffold and theme a site

Validates spec Success Criteria SC-001..SC-003 end to end, purely via local-preview evidence (no live deployment required — that's v0.02).

## Prerequisites

- `biblio-cli` binary installed (GitHub Releases download, or `go install github.com/.../biblio-cli@latest` — FR-006 / [contracts/cli-commands.md](./contracts/cli-commands.md))
- Node + a package manager available locally (the generated project wraps `vite dev` / `wrangler dev`)

## Scenario 1 — Scaffold a new site from one command (spec Scenario 1)

```bash
biblio-cli new my-profile-site --type=profile --domain=example.com --theme=profile-default
cd my-profile-site
ls content/ content/site.yaml wrangler.jsonc tokens.css
```

**Expected outcome**: a complete project directory exists containing a populated `content/` tree, `content/site.yaml` (with `mcp_enabled: true` and the declared site type/name/description), `wrangler.jsonc` (route pre-filled from `--domain`), and `tokens.css` matching the `profile-default` preset (FR-001, FR-002, FR-004; spec Scenario 1 acceptance criteria 1–2).

Repeat for `--type=log` and `--type=docs` to validate SC-002 ("all three site-type scaffolds produce a runnable local preview and a structurally correct content index on the first attempt").

## Scenario 2 — Index authored content (spec Scenario 2)

```bash
biblio-cli index
cat content-index.json
```

**Expected outcome**: `content-index.json` is produced, with one entry per authored page/entry, each carrying a correctly derived `slug`/`path` and its structured frontmatter fields (FR-003; [data-model.md](./data-model.md) Content Index entity). Introduce a file with a malformed/missing required frontmatter field and re-run — expect a clear, actionable per-file error, not a silently incomplete index (spec Edge Cases).

## Scenario 3 — Preview locally (spec Scenario 3)

```bash
biblio-cli dev
```

**Expected outcome**: both the site (`vite dev`) and its supporting Worker (`wrangler dev`) start together behind this one command, with their addresses surfaced; opening the site address in a browser shows the seeded content live, styled per the chosen theme (FR-005; spec Scenario 3).

## Demonstrating SC-001..SC-003

- **SC-001**: Steps above go from "no project" to "a running local preview of a themed, content-indexed site" using only `biblio-cli new` → `biblio-cli index` → `biblio-cli dev` — no manual file assembly or config editing.
- **SC-002**: Repeat Scenario 1 + Scenario 2 for all three `--type` values; each MUST produce a runnable preview and a structurally correct index on the first attempt.
- **SC-003**: All of the above is demonstrable purely from local preview — no deployment (v0.02) or agent/terminal-facing capability (v0.03) is invoked or required to show this release's value.
