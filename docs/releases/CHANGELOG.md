# Changelog

All notable user-facing changes to `biblio-cli` and its generated sites are
recorded here, release by release. See `docs/releases/v<version>/` for the
full user guide of each release.

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
