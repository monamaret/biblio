# Implementation Plan: `quarterly` site type

**Branch**: `monamaret/004-quarterly-site-type` | **Date**: 2026-10-15 | **Spec**: [spec.md](spec.md)

**Input**: Feature specification from `specs/004-quarterly-site-type/spec.md`

**Feature item**: [011-quarterly-site-type.md](../features/011-quarterly-site-type.md) · **Roadmap**: [§11](../research/roadmap.md#11-quarterly-site-type) · **Target release**: v0.08

## Summary

Add `quarterly`, a fourth site type: an issue-based magazine with
categories, footnotes and margin notes, configurable typesetting and
grayscale plates. The work has three parts:

1. **Go (`biblio-cli`)**:
   - per-file-kind content schemas and a quarterly-aware index (derived
     issue fields, categories, drafts, structured notes, a Markdown-view
     body);
   - a typed typesetting schema with strict validation and a resolver that
     writes `tokens.css` and the index's structural switches;
   - pure-Go grayscale plate conversion with hash-named outputs;
   - a new pre-build step that `build` and `dev` run for every site type
     (validate + index), with tokens and plates for quarterly only.
2. **Template**: `templates/site/quarterly/`, copied from `log` and
   adapted, with contents/archive/category/article/colophon routes,
   CSS-only sidenotes, a tokens-only stylesheet, self-hosted fonts, and
   four MCP tools.
3. **Parity**: register the type, extend the cross-type review and
   `update`/`upgrade` tests, and leave profile, log and docs output
   unchanged.

All design decisions are in [research.md](research.md) (R1–R15).

## Technical Context

**Language/Version**: Go 1.24 (`biblio-cli`); TypeScript 5.6 with React 19 and React Router 7.18 (template)

**Primary Dependencies**:
- Go, existing: Cobra, Viper, Bubble Tea, `go.yaml.in/yaml/v3`.
- Go, new direct: `github.com/yuin/goldmark` v1.7.13 (was indirect), `golang.org/x/text` (was indirect), `golang.org/x/image` (only new module).
- Template, unchanged from `log`: react-markdown, remark-gfm, `agents`, `@modelcontextprotocol/sdk`, zod, md4x, Vite 6, Wrangler 4.

**Storage**: files only: `content/` (authored), `content-index.json`, `public/tokens.css`, `public/plates/` (generated). No database (Principle V).

**Testing**: `go test ./...` (unit and `test/integration/`), `./scripts/check.sh`, and Vitest (`npm test`) in the template.

**Target Platform**:
- CLI: Linux, macOS and Windows on amd64 and arm64 (goreleaser, pure Go, no cgo).
- Sites: Cloudflare Workers with Static Assets.

**Project Type**: CLI plus embedded site templates (each template is its own npm project)

**Performance Goals**:
- A pre-build step with no changed images takes well under a second on a 15-article issue (only file-exists checks).
- The first conversion of a 12 MP photo takes about a second or less.
- No client JS is added to prerendered quarterly pages.

**Constraints**:
- No cgo. No font CDN.
- `tokens.css` holds variables only.
- No author string reaches CSS without a closed validation pattern.
- Profile, log and docs pages and assets are byte-identical before and after.
- Template fonts total under 350 KB.

**Scale/Scope**:
- One author per site; about 5–15 articles per issue; tens of issues over the site's life.
- New Go code: 1 package (`internal/prebuild`), plus additions to `content`, `index`, `theme`, `scaffold` and `cli`.
- Template: about 10 routes, components and tools, plus their tests.

No `NEEDS CLARIFICATION` remains: the four planning items from roadmap §0.5 are settled in research.md R1–R4.

## Constitution Check

*GATE: Must pass before Phase 0 research. Re-checked after Phase 1 design.*

| Principle | Status | How this plan complies |
|---|---|---|
| **I. Single-author, minimal dependencies** | Pass | One new Go module (`x/image`); goldmark and x/text were already in the graph. No new npm packages. No plugin system, no arbitrary CSS, no author pages. Every setting maps to an interview answer. |
| **II. One pipeline, one source of truth** | Pass, with one noted split | Go is the only producer of the index, `tokens.css`, plates, structured notes and `markdown_body`; the template reads them. The template still parses Markdown for **rendering** (react-markdown plus a margin-note remark plugin), as all three existing types already do. Drift between the two is caught by the notes-agree-with-index contract test (R2). Recorded under Complexity Tracking. |
| **III. CLI-first** | Pass | No new command, no npm script the author must run. `build`/`dev` gain the pre-build step; the deploy workflow still calls only `biblio-cli build`/`deploy`; `quarterly` appears in the TUI from `scaffold.SiteTypes`. |
| **IV. Text-first MCP and terminal** | Pass | All four tools return self-sufficient text; `get_article`'s `structuredContent` is extra, and its text answer includes every note. Markdown and ANSI views keep every note (margin notes as footnotes). |
| **V. Simplest topology** | Pass | Same single Worker and in-memory JSON index. No D1 or Vectorize. Plates are static assets. |
| **VI. Test-first** | Pass (planned) | research.md R15 lists the failing tests that come before each part. `/speckit-tasks` must put test tasks before their implementation tasks. |
| **Security & operational posture** | Pass | `/mcp` stays public and unauthenticated; `mcp_enabled` stays the only MCP setting. New risks are handled: CSS injection (R9 closed patterns), image decompression bombs and symlink escapes (R4, data-model Plate), raw HTML in notes (not rendered). |
| **Decision recording and research workflow** | Pass | Research doc 004 with its Decision Log, feature item 011, roadmap §11 and §0.5, and the architecture.md §2 sync were all in place before this plan. New decisions R1–R4 must be added to 004's Decision Log and §0.5 checked off (a tasks item). |
| **Release process** | Deferred to release | The `specs/releases/0.08-…/` package (feature list, notes, tests, usage guide with create/configure/update/deploy/destroy) is put together when v0.08 closes, not now. |

**Gate result (pre-research):** pass. **Gate result (post-design):** pass.
Phase 1 added nothing that changes a row. `markdown_body` and the index's
`typesetting` object both keep Go as the single producer.

## Project Structure

### Documentation (this feature)

```text
specs/004-quarterly-site-type/
├── spec.md
├── plan.md              # this file
├── research.md          # Phase 0: R1–R15
├── data-model.md        # Phase 1
├── quickstart.md        # Phase 1: validation scenarios
├── contracts/
│   ├── cli.md                 # new/build/dev/index/update changes
│   ├── config-and-tokens.md   # typesetting schema, errors, tokens.css, quarterly.css rules
│   ├── site-surface.md        # routes, pages, negotiation, llms.txt
│   └── mcp-tools.md           # list_issues, get_issue, list_articles, get_article
├── checklists/requirements.md
└── tasks.md             # Phase 2 (/speckit-tasks, not created here)
```

### Source Code (repository root)

```text
internal/
├── content/
│   ├── schema.go            # + ParseItemKind (issue/article/colophon), SiteConfig.Typesetting
│   ├── quarterly.go         # new: per-kind field validation, category keys
│   └── notes.go             # new: goldmark footnotes + {>> <<} scanner → Note[], markdown_body rewrite
├── index/
│   └── generator.go         # + Options{IncludeDrafts, Plates}; quarterly derived fields,
│                            #   issues/categories/current_issue/typesetting; cross-file checks
├── theme/
│   ├── presets.go           # + quarterly-default
│   ├── typesetting.go       # new: Typesetting struct, validation, Resolve(preset, overrides)
│   └── tokens.go            # GenerateTokensCSS → emits resolved quarterly properties too
├── plates/                  # new: grayscale/resize/encode, hash naming, safety caps, orphan cleanup
├── prebuild/                # new: Run(projectDir, Mode) → validate, tokens, plates, index; Report
├── scaffold/
│   ├── assets.go            # SiteTypes += "quarterly"; Walk skips only .gitkeep
│   └── generator.go         # unchanged logic; quarterly flows through
└── cli/
    ├── build.go             # call prebuild.Run(Build) before vite build
    ├── dev.go               # call prebuild.Run(Dev) before the supervisor
    └── index.go             # unchanged for 3 types; quarterly uses index.Options

templates/site/quarterly/    # new, copied from templates/site/log, then adapted
├── .gitignore               # public/plates/, .biblio/
├── .github/workflows/deploy.yml
├── content/                 # sample: 2 issues, 5 articles (one with dense notes), colophon, images
├── public/
│   ├── quarterly.css        # tokens-only stylesheet, @font-face, sidenotes, print basics
│   └── fonts/               # Archivo variable woff2, IBM Plex Mono woff2, OFL texts
├── app/
│   ├── root.tsx             # links /tokens.css then /quarterly.css
│   ├── routes.ts            # index, issues, issues/:issue, category/:category, * splat
│   ├── routes/              # contents.tsx, archive.tsx, category.tsx, detail.tsx (+ tests, theme.test.tsx)
│   ├── content/             # load.ts (issues/categories accessors), render.tsx, notes.tsx,
│   │                        #   remark-margin-notes.ts, filters.ts (articles)
│   ├── negotiation/         # markdown.ts prefers markdown_body; ansi.ts, accept.ts as in log
│   ├── discovery/llms-txt.ts
│   └── mcp/                 # server.ts, agent.ts, query.ts, tools/{shared,quarterly}.ts (+ tests)
├── react-router.config.ts   # prerender list from index issues/categories/entries
└── package.json, vite/vitest/tsconfig, workers/entry.ts   # as in log

test/integration/
├── pilot-site-quarterly/    # new: scaffold → prebuild → index/tokens/plates assertions
├── cross-type-review/       # extend theming and terminal reviews to four types
├── update-three-sites/      # + quarterly upgrade case
└── prebuild-regression/     # new: profile/log/docs output unchanged apart from the index
```

**Structure decision**: follow the existing layout. Go logic goes into the
packages that already own each concern. Two new packages: `plates`, so
image handling is isolated and fuzzable, and `prebuild`, so ordering lives
in one place used by both `build` and `dev`. The template is a fourth
independent tree, as the research doc chose (copying, not sharing code
between templates).

## Implementation phases (input to `/speckit-tasks`)

Order follows the spec's priorities; within each phase, tests come first.

1. **Foundation (blocks everything):**
   - `SiteTypes` and Walk change, `quarterly-default` preset, `theme.Typesetting` with validation and resolver, `ParseItemKind`, index options.
   - `prebuild` with validate and index wired into `build`/`dev`.
   - The regression test proving profile, log and docs are unchanged.
2. **US1, publish an issue (P1):**
   - Template copy; issue/article/colophon schemas; derived fields; ordering; categories; drafts.
   - Routes, prerender list, negotiation and `llms.txt`.
   - Sample content.
3. **US2, notes (P2):** Go notes parser and `markdown_body`; remark margin-note plugin; sidenote markup and CSS; contract test that rendered notes agree with the index.
4. **US3, typesetting (P3):** tokens generation for quarterly; `quarterly.css`; fonts; the font-file warning; `theme.test.tsx`.
5. **US4, plates (P4):** the `plates` package (with fuzz tests on decode); body and cover rewriting; cache and orphan cleanup; reporting.
6. **US5, agent tools (P5):** `quarterly.ts` tools and their contract tests.
7. **US6, parity (P6):** cross-type review with four types; `update`/`upgrade` integration; TUI check; docs:
   - 004 Decision Log (R1–R4), roadmap §0.5 checked off;
   - AGENTS.md mentions four types;
   - user guide additions for the release.

## Complexity Tracking

| Deviation | Why needed | Simpler alternative rejected because |
|---|---|---|
| `biblio-cli build`/`dev` stop being thin wrappers (the roadmap §0.3 wording) | Validation, tokens, plates and indexing must run before Vite, in Go (Principle II); the clarification extended validate and index to all types | Running them by hand or as npm scripts breaks Principle III and is the manual `index` gap this removes |
| Markdown is parsed in Go (index, notes) and in TS (rendering) | Rendering has always been react-markdown in every template; notes must be structured in Go for `get_article` and 005 | Rendering HTML in Go would break the template model; parsing notes in TS would duplicate schema logic. A contract test holds the two parsers together |
| Two bodies per quarterly article (`body`, `markdown_body`) | The HTML renderer needs the margin-note syntax; text views need plain footnotes | Rewriting in TS at request time puts parsing in a second language |
