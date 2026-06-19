# Changelog

All notable user-facing changes to `biblio-cli` and its generated sites are
recorded here, release by release. See `docs/releases/v<version>/` for the
full user guide of each release.

## v0.02 — Build, deploy, and verify a real site

- Added `biblio-cli build` to pre-render every content-driven page into
  static output — a thin wrapper over `vite build`, driving React Router
  v7's native `prerender` against `content-index.json`.
- Added `biblio-cli deploy` to publish a built site to its configured
  public address — a thin wrapper over `wrangler deploy`; re-running it
  against an already-live site updates cleanly in place.
- Added React Router v7 framework-mode wiring (`react-router.config.ts`,
  `app/root.tsx`, list/detail route modules) to the `profile` site template,
  reading `content-index.json` directly rather than re-parsing authored
  Markdown on the TypeScript side.
- Added a generated GitHub Actions workflow
  (`.github/workflows/deploy.yml`) invoking `biblio-cli build`/`deploy` on
  push, keeping the whole lifecycle behind the `biblio-cli` surface.

## v0.01 — Scaffold and theme a site

- Added `biblio-cli new <project-name>` to scaffold a complete, runnable,
  themed site project (`profile`, `log`, or `docs`) from one guided command,
  generating `content/`, `content/site.yaml`, `wrangler.jsonc`, and
  `tokens.css`.
- Added `biblio-cli index [path]` to generate the canonical
  `content-index.json` from an authored `content/` directory, validating
  required frontmatter per site type and reporting every failing file on
  error.
- Added `biblio-cli dev` to run `vite dev` and `wrangler dev` together as
  one supervised local-preview process pair.
- Added one ready-to-use theme preset per site type: `profile-default`,
  `log-default`, `docs-default`.
- Added a Bubble Tea guided flow so `new` can be run with no flags at all,
  prompting interactively for site type, domain, and theme.
- Added distribution via prebuilt GitHub Releases binaries and
  `go install github.com/monamaret/biblio/cmd/biblio-cli@latest`.
