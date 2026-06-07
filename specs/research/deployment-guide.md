# Guide: Creating, configuring, and deploying a `biblio` site

**Status:** Scaffold — drafted ahead of `biblio-cli`'s implementation, from the decisions confirmed in [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md). All design questions that shaped this guide (tracked in [roadmap.md](roadmap.md) §0) are now resolved; what remains is for the CLI/scaffold to actually be built — this guide will be corrected against real command output and behavior as that lands. See [architecture.md](architecture.md) for diagrams of the system this guide deploys.

## 1. Prerequisites

- A Cloudflare account with **Workers (Pro)** and at least one custom domain already configured
- A **GitHub** repository for the site (CI/CD via GitHub Actions — see [000-tech-stack.md](000-tech-stack.md))
- `biblio-cli` installed — download a prebuilt binary from GitHub Releases (the primary path), or `go install .../biblio-cli@latest` if you have a Go toolchain (see the [confirmed install/distribution decision](001-static-site-platform-with-mcp.md#decision-log))
- The `wrangler` CLI, authenticated against the Cloudflare account

## 2. Choosing a site type

`biblio-cli` scaffolds three site types — pick the one that matches the new site. (See [001 → frontmatter schemas](001-static-site-platform-with-mcp.md#architecture-patterns) for the full content model of each.)

| Site type | Content model (high level) | Status |
| --- | --- | --- |
| Personal profile / CV | Profile, Experience, Skills, Projects | — |
| Technical project log | Dated, tagged, project-linked entries | Priority scaffold (validate first) |
| Project documentation | Ordered/sectioned pages, nested tree | Priority scaffold (validate first) |

## 3. Scaffolding a new site

```sh
biblio-cli new <project-name> --type=<profile|log|docs> --domain=<example.com> --theme=<preset>
```

Flags select the site type, domain, and theme preset directly; omit them to get an interactive Bubble Tea-driven prompt flow instead (the [confirmed `new` command shape](001-static-site-platform-with-mcp.md#decision-log)).

Based on the confirmed architecture, scaffolding a new site is expected to generate:

- A combined **React + Vite + Worker** project, wired per the [selected architecture](architecture.md#1-high-level-architecture-selected)
- A `content/` directory, pre-populated per the [confirmed layout convention](001-static-site-platform-with-mcp.md#decision-log) for the chosen site type — e.g. `profile.md` + `experience|skills|projects/*.md` for a profile/CV site, `entries/YYYY-MM-DD-slug.md` for a project log, `docs/<section>/<page>.md` for documentation — plus `content/site.yaml` for the shared site-level config, and a few example entries
- A generated `tokens.css` from the site type's [theme preset](001-static-site-platform-with-mcp.md#architecture-patterns) (CV/profile: warmer/personal; project log: cooler/technical with monospace accents; docs: neutral/high-contrast)
- A `wrangler.jsonc` with the route/custom-domain binding pre-filled
- A `site.yaml` config file holding `site_name`/`description`/`site_type`, plus `mcp_enabled` (the [confirmed MCP opt-in/out flag](001-static-site-platform-with-mcp.md#decision-log), default `true`)
- A GitHub Actions workflow wired to the [three-step build pipeline](architecture.md#2-build--deployment-pipeline)

## 4. Authoring content

- Content lives under `content/`, laid out per the [confirmed directory convention](001-static-site-platform-with-mcp.md#decision-log) (flat-by-concept directories matching the site's type — `experience/`/`skills/`/`projects/` for a profile site, `entries/` for a log, `docs/<section>/` for documentation — plus `site.yaml` for shared config), as **plain Markdown files with YAML frontmatter** — *not* MDX
- Each file is loaded at build time via Vite's `import.meta.glob`, parsed with `gray-matter`, and rendered with `react-markdown` (+ `remark-gfm`)
- Frontmatter fields are fixed per site type — see the [confirmed schemas](001-static-site-platform-with-mcp.md#architecture-patterns):
  - *Profile/CV:* Profile (`name`/`bio`/`links`), Experience, Skill, Project entries
  - *Project log:* Entry (`title`/`date`/`project`/`tags`/`summary`/`status`) — `slug` derived from the **filename**
  - *Docs:* Page (`title`/`order`/`section`/`summary`) — `slug`/`path` derived from **file location**
- A site-level config (not a per-file frontmatter block) supplies `site_name`/`description`/`site_type`

## 5. Configuring the MCP surface

- The `/mcp` endpoint is **public, with no auth/access control**, by design — Cloudflare's standard protections (rate limiting, WAF, bot management) are the operational safety net (see the [Decision Log](001-static-site-platform-with-mcp.md#decision-log))
- The **`mcp_enabled` flag** in `content/site.yaml` (default `true`) determines whether the content index is generated and the `/mcp` route is mounted at all, so a site that doesn't want an MCP surface — set `mcp_enabled: false` — skips it entirely (see the [confirmed flag shape](001-static-site-platform-with-mcp.md#decision-log))
- The tool surface is fixed per site type — e.g. a project-log site gets `list_entries`/`get_entry`/`search_entries`, a docs site gets `get_toc`/`get_page`/`search_docs`, every site gets `get_site_info` and a generic `search` (full mapping in [001 → tool surface](001-static-site-platform-with-mcp.md#architecture-patterns))

## 6. Local development

```sh
biblio-cli dev
```

`biblio-cli dev` wraps `vite dev` and `wrangler dev` behind one command (the [confirmed local-dev workflow decision](001-static-site-platform-with-mcp.md#decision-log)) — no need to run the two dev servers separately or coordinate them by hand.

## 7. Building and deploying

The [build pipeline](architecture.md#2-build--deployment-pipeline) runs in three steps, normally triggered by a push via the generated GitHub Actions workflow:

1. **`biblio-cli` generates `content-index.json`** from the content directory (the canonical artifact — single source of truth for both the site and the MCP server)
2. **`react-router build`** prerenders every content-driven route to static HTML/CSS/JS, via the framework's native `prerender` config
3. **`wrangler deploy`** bundles the Worker (which `import`s `content-index.json` directly, inlined for in-memory querying) and publishes the static assets + Worker to Cloudflare

```sh
biblio-cli build
biblio-cli deploy
```

`biblio-cli build`/`biblio-cli deploy` are thin wrapper commands over the underlying `vite build`/`wrangler deploy` + CI steps (the [confirmed build/deploy command-surface decision](001-static-site-platform-with-mcp.md#decision-log)) — the CLI stays involved end to end rather than handing you off to raw `package.json` scripts after scaffolding.

## 8. Verifying a deployment

- Visit `https://<domain>/` — the prerendered site should load
- Connect an MCP client to `https://<domain>/mcp` — it should respond to `initialize` and `tools/list`, and `tools/call` should return text-first content (with an optional UI resource for hosts that support it)
- `curl -H "Accept: text/markdown" https://<domain>/<some-page>` — once the [terminal/CLI view](001-static-site-platform-with-mcp.md#architecture-patterns) is built, this should return a Markdown-rendered response with `Vary: Accept, User-Agent` set
- Check **Cloudflare Workers Observability** for errors and request logs

## 9. Escalating storage (only when needed)

This only becomes relevant once a site outgrows the bundled-JSON default — see the [progressive-enhancement storage approach](001-static-site-platform-with-mcp.md#architecture-patterns) and its heuristic thresholds.

- **JSON → D1:** provision a D1 database, add the binding to `wrangler.jsonc`, and switch the content-index generation step to populate D1 (still derived at build time from the same content directory, not hand-maintained)
- **Keyword → Vectorize:** generate embeddings at build time via OpenRouter (see [000-tech-stack.md](000-tech-stack.md)), provision a Vectorize index, and switch `search_docs`/`search_entries` to query it instead of doing in-memory keyword matching

**Open question** (conditional — tracked in [roadmap.md §0](roadmap.md#0-outstanding-design-questions); only worth answering once a real site actually reaches one of the documented thresholds, see [roadmap.md §10](roadmap.md#10-storage-escalation-conditional--not-built-up-front)): concrete step-by-step migration instructions for each escalation path.

## 10. Updating an existing site

```sh
biblio-cli update
```

The *direction* is confirmed — re-scaffolding/upgrading an existing site as shared templates, themes, or schemas evolve is handled through the `biblio-cli` command surface itself (`biblio-cli update`/`biblio-cli upgrade`), not a separate tool or manual process (see the [confirmed decision](001-static-site-platform-with-mcp.md#decision-log)). The exact mechanics — diffing strategy, conflict handling, what gets re-templated vs. preserved — remain **open and deferred** (tracked in [roadmap.md §0](roadmap.md#0-outstanding-design-questions); only answerable once more than one site exists to migrate, see [roadmap.md §9](roadmap.md#9-replicate-to-the-remaining-priority-site-types)).

## References

- [001-static-site-platform-with-mcp.md](001-static-site-platform-with-mcp.md) — the research and decisions this guide assumes throughout
- [architecture.md](architecture.md) — diagrams of the system being created/deployed
- [roadmap.md](roadmap.md) — the implementation task breakdown this guide will be checked against as `biblio-cli` and the scaffold get built (its §0 design-question tracker is now fully resolved)
