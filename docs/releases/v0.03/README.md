# biblio-cli — v0.03 User Guide

**Release**: Make the site usable by AI agents and terminal clients
**Status**: builds on [v0.02](../v0.02/README.md) — the deployed site now
also answers AI assistants over MCP and serves clean Markdown/ANSI to
terminal clients. No new `biblio-cli` commands; everything here is
build-time wiring controlled by one existing flag.

## What's new

A site scaffolded with `biblio-cli new` (agent-facing capability on by
default) now mounts an in-Worker MCP server and a content-negotiation
layer on the **same** Worker that serves the static site — no separate
service, no new deployment step.

### The one configuration surface: `--mcp-enabled`

```sh
biblio-cli new my-profile-site --type=profile --domain=example.com --theme=profile-default
# --mcp-enabled defaults to true; pass --mcp-enabled=false to scaffold
# a site with no agent-facing surface at all.
```

This single flag (written to `content/site.yaml`'s `mcp_enabled`) is a
**build-time** decision, not a runtime toggle:

- `mcp_enabled: true` — `wrangler.jsonc` mounts a Worker `main`
  (`workers/entry.ts`), the `durable_objects` binding the MCP server
  needs, `nodejs_compat`, and `assets.run_worker_first` so the Worker can
  decide per request whether to answer it itself or fall through to the
  static site.
- `mcp_enabled: false` — none of that exists in `wrangler.jsonc` at all.
  There is no "disabled but present" state; the `/mcp` endpoint simply
  isn't part of the deployed artifact.

To change it after scaffolding, re-scaffold with the flag you want, or
edit `content/site.yaml`'s `mcp_enabled` by hand and re-run `biblio-cli
build`/`deploy`.

## What an AI assistant can do

Connect an MCP client to `https://<your-domain>/mcp` and it can list and
call:

- `list_pages`, `get_page`, `search_content` — shared, work on any site type
- `get_profile`, `list_projects` — profile-site-specific

Every tool's answer is plain Markdown/text on its own; `search_content`
also attaches a small themed chart view for graphical MCP hosts that
support it, without ever replacing the text answer. The endpoint is
intentionally public with no authentication, same as the rest of the
deployed site's posture.

## What a terminal or `curl` user sees

```sh
curl -H "Accept: text/markdown" https://<your-domain>/about     # clean Markdown
curl -A "curl/8.0" -H "Accept: text/markdown" https://<your-domain>/about   # ANSI-rendered
```

Every page is reachable this way — the same URL, content-negotiated by
`Accept`/`User-Agent` rather than a separate address. Browsers requesting
the same URL still get the unchanged v0.02 prerendered HTML.

`https://<your-domain>/llms.txt` lists the site's pages for crawler/agent
discovery — generated automatically at build time.

## Known scope limits in v0.03

- Wired up for the `profile` site type only — `log`/`docs` get the same
  treatment in v0.04/v0.05 alongside their own replication.
- The richer MCP App UI view (the chart for `search_content`) is a
  minimal proof of the pattern, not a polished visualization.
- No authentication on `/mcp`, by design — relies on Cloudflare's
  zone-level protections only, same as the rest of the site.
