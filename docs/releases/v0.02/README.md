# biblio-cli — v0.02 User Guide

**Release**: Build, deploy, and verify a real site
**Status**: builds on [v0.01](../v0.01/README.md) — adds `build` and
`deploy` so a scaffolded project can become a real, publicly reachable
site. No agent/terminal-facing capability (v0.03) yet.

`biblio-cli` now takes a scaffolded, indexed project all the way from
authored content to a live, pre-rendered site, with two more commands:
`build` and `deploy`.

## Commands

### `biblio-cli build`

Pre-renders every content-driven page of the current project into static
output — a thin wrapper over `vite build`, which drives React Router v7's
native `prerender` against the project's `content-index.json`.

```sh
cd my-profile-site
biblio-cli build
```

- Requires `wrangler.jsonc` in the current directory (run `biblio-cli new`
  first if it's missing).
- Produces one pre-rendered HTML page per `content-index.json` entry, plus
  the list/navigation home page, in `dist/client/`.
- **Failure behavior**: a malformed-frontmatter file or any other prerender
  validation failure fails the build with a clear, per-file message —
  never a silently incomplete page set.

### `biblio-cli deploy`

Publishes a built site to its configured public address — a thin wrapper
over `wrangler deploy`.

```sh
biblio-cli deploy
```

- Requires a successful `biblio-cli build` first (checks for `dist/client`)
  and `wrangler.jsonc` in the current directory.
- Publishes to the route/custom-domain binding `biblio-cli new` pre-filled
  from `--domain`, and reports that address on success.
- Re-running `deploy` against an already-live site updates it cleanly in
  place — there is no separate "first deploy" vs. "redeploy" command.
- **Failure behavior**: surfaces the underlying `wrangler deploy` failure
  (auth, route-binding, etc.) with enough context to diagnose — never a
  silent partial publish.

## End-to-end walkthrough

```sh
# 1. Scaffold, author content, and index it (v0.01)
biblio-cli new my-profile-site --type=profile --domain=example.com --theme=profile-default
cd my-profile-site
# ... edit content/about.md, add content/projects/*.md, etc. ...
biblio-cli index

# 2. Build (v0.02)
biblio-cli build

# 3. Deploy (v0.02) — requires Cloudflare/Wrangler authentication
biblio-cli deploy
# open the printed public address in a browser
```

## Known scope limits in v0.02

- One pilot site type (`profile`) has wired-up routes/prerender in this
  release — `log` and `docs` get the same treatment in v0.04/v0.05.
- No agent/terminal-facing (MCP) capability yet — the deployed site is
  browser-facing HTML only; the `/mcp` route mounts in v0.03.
- `biblio-cli deploy` requires you to have already authenticated Wrangler
  with your own Cloudflare account — `biblio-cli` does not manage Cloudflare
  credentials itself.
