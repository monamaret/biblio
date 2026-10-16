# Contract: CLI surface

No new commands. Changed behavior only. Existing contracts:
[v0.01 cli-commands](../../releases/0.01-scaffold-and-theme/contracts/cli-commands.md),
[v0.02 cli-commands](../../releases/0.02-build-deploy-verify/contracts/cli-commands.md),
[v0.06 update-and-escalation](../../releases/0.06-general-availability/contracts/update-and-escalation.md).

## `biblio-cli new <name> --type=quarterly [--domain=…] [--theme=quarterly-default]`

- `quarterly` is accepted by `--type` and offered in the guided TUI flow,
  after `docs` (FR-001).
- The only theme for `quarterly` is `quarterly-default`. As for the other
  types, omitting `--theme` opens the guided prompt, which offers that one
  preset.
- Output: the `templates/site/quarterly/` tree, `content/site.yaml`
  (`site_type: quarterly`, no `typesetting:` block, so the defaults apply),
  `wrangler.jsonc`, `public/tokens.css`, the scaffold manifest, and a
  `.gitignore` covering `public/plates/` and `.biblio/`.
- `--type=quarterly --theme=log-default` fails with the existing
  "not a valid preset for site type" error.

## `biblio-cli build` (all site types)

Before `vite build`, runs the pre-build step in Build mode
([research.md R8](../research.md#r8-the-pre-build-step)):

```text
Validating content/site.yaml… ok
Generating public/tokens.css (quarterly-default + 3 overrides)… ok          quarterly only
Converting images: 2 converted, 9 skipped, 1 removed → public/plates/       quarterly only
Indexing: 14 item(s), 1 draft issue excluded → content-index.json
Building (vite build, prerendering content-driven pages)…
```

- Any validation error stops the command before `vite build`, with a
  non-zero exit status. All validation errors are listed together, each
  naming its file or key.
- Warnings (`warning: …`) are printed and the build continues.
- Profile, log and docs print only the validate and index lines.

## `biblio-cli dev` (all site types)

- Runs the pre-build step once, in Dev mode (drafts included), before
  starting `vite dev` and `wrangler dev`. A validation error stops `dev`
  before either server starts.
- Content or settings changes made while `dev` runs take effect on the
  next `dev` start (FR-024). There is no file watcher in this release.

## `biblio-cli index [path]`

- Unchanged for profile, log and docs.
- For quarterly it runs only the validate and index steps, with publish
  semantics (drafts excluded). It does not convert images, so a
  `plate`/`cover_plate` URL is written only when the plate already exists
  in `public/plates/`. `build` and `dev` always run the full step.

## `biblio-cli update` / `upgrade`

- Treat `templates/site/quarterly/` like the other trees: manifest-tracked
  3-way merge, author content and `content/site.yaml` never overwritten
  (FR-003).
- Generated files (`public/tokens.css`, `public/plates/*`,
  `content-index.json`) are not manifest-tracked template files.
