# biblio-cli — v0.06 User Guide (General Availability)

**Release**: General Availability
**Status**: builds on [v0.05](../v0.05/README.md) — the same CLI, content
pipeline, theming, MCP server, and terminal/negotiation layer; this release
adds the `update`/`upgrade` commands and closes the platform's last
roadmap-deferred item, marking **General Availability**.

## What's new

The one piece of the CLI surface the roadmap deliberately deferred until
more than one real site existed to migrate — a reliable way to carry an
existing site forward as the scaffold evolves, without losing authored
content — is now built and proven across all three site types. Paired with
a documented, dry-run-proven storage-escalation guide for the two
conditional paths (D1, Vectorize), every roadmap deliverable (§1–§10) is
now either complete-and-demonstrated or proven-executable-on-demand. See
[GA-READINESS-RECORD.md](../../../test/integration/ga-readiness-review/GA-READINESS-RECORD.md)
for the full cross-reference.

## The complete `biblio-cli` lifecycle

### Build from source

A prebuilt binary isn't required — `biblio-cli` builds from source with
nothing but a Go toolchain:

```sh
git clone https://github.com/monamaret/biblio.git
cd biblio
go build -o biblio-cli ./cmd/biblio-cli   # requires Go 1.24+ (see go.mod)
./biblio-cli --help
```

Or run it directly without producing a binary on disk:

```sh
go run ./cmd/biblio-cli --help
```

Or install it onto your `$GOPATH/bin` for use from anywhere:

```sh
go install github.com/monamaret/biblio/cmd/biblio-cli@latest
```

(Note the `/cmd/biblio-cli` suffix — the module root itself has no Go package. If `biblio-cli` isn't found afterward, `$(go env GOPATH)/bin` likely isn't on your `PATH` yet — see the root [README.md](../../../README.md#troubleshooting).)

Prebuilt binaries are also published via `goreleaser` on GitHub Releases
(Principle III) — use whichever path is more convenient; all three produce
the same CLI.

### Create

```sh
biblio-cli new my-site --type=profile --domain=example.com --theme=profile-default
```

(`--type` is one of `profile`, `log`, or `docs`; omit any flag to be
prompted interactively.) Generates a complete, runnable, themed project:
`content/` (seeded per site type), `content/site.yaml`, `wrangler.jsonc`
(routing pre-filled from `--domain`), `tokens.css`, and — as of this
release — a `.biblio/scaffold-manifest.json` baseline that `update`/
`upgrade` diff against later.

### Configure

Author content as Markdown + YAML frontmatter under `content/`, per the
site type's own shape (profile/experience/skills/projects, dated log
entries, or section/order docs pages). Toggle the in-Worker MCP server via
`mcp_enabled` in `content/site.yaml` — the sole MCP configuration surface.

```sh
biblio-cli index   # regenerate content-index.json after any content change
```

### Update

Bring the project's scaffold-tracked parts (workers/, app/, build config,
CI workflow) up to date with the version of `biblio-cli` you're running,
while your authored content and any customizations to those scaffold files
are preserved untouched:

```sh
biblio-cli update
```

If a file you've customized would otherwise be overwritten by a scaffold
change, the run halts before writing anything and reports the conflict —
resolve it explicitly, then re-run:

```sh
biblio-cli update --resolve workers/entry.ts=keep-author
# or: --resolve workers/entry.ts=use-scaffold
```

`biblio-cli upgrade <version>` is the version-targeted variant — it
succeeds only when `<version>` matches the `biblio-cli` binary you're
running, since each binary embeds exactly one scaffold-template version
(no historical multi-version store).

### Deploy

```sh
biblio-cli build    # content-index -> react-router prerender -> wrangler bundle
biblio-cli deploy   # wrangler deploy, to the route configured in wrangler.jsonc
```

Re-running `deploy` against an already-live site updates it in place. The
generated GitHub Actions workflow (`.github/workflows/deploy.yml`) runs the
same two commands on push, so the CLI stays the single, consistent
lifecycle surface (Principle III) whether run by hand or by CI.

### Destroy

There is no `biblio-cli destroy` command — tearing down a deployed site is
a one-off action, not a recurring lifecycle step, so it stays a direct call
to the underlying tool rather than a speculative wrapper (Principle I):

```sh
npx wrangler delete   # from the project directory; removes the deployed Worker
```

The project directory itself (and your authored `content/`) is untouched —
deleting it (or not) is a separate, ordinary filesystem decision.

### What an AI assistant can do

Unchanged from v0.05 — connect an MCP client to `https://<your-domain>/mcp`
for the shared and per-site-type tools documented in the
[v0.05 guide](../v0.05/README.md#what-an-ai-assistant-can-do).

### Storage escalation, if you ever need it

[docs/storage-escalation-guide.md](../../storage-escalation-guide.md)
documents checkable thresholds for the two conditional paths (D1, Vectorize)
and records real dry-run proofs that both are executable on demand. Neither
is enabled by default — escalate only if your own site's real
`content-index.json` scale crosses a documented threshold.

## Known scope limits in v0.06

- `upgrade <version>` only supports upgrading to the exact version of the
  `biblio-cli` binary you're running — there is no historical
  multi-version scaffold-template store to upgrade *from* an old version
  *to* a specific older or newer one across binaries.
- D1 and Vectorize remain proven-executable-on-demand, not adopted — per
  Principle V, actually escalating either without a real, currently
  crossed threshold would be the speculative build-up-front the constitution
  forbids.
- No authentication on `/mcp`, by design, same as every other site type.
