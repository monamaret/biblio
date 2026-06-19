# Contracts: `biblio-cli` command surface additions (v0.02)

**Type of contract**: CLI command/flag schema, extending [v0.01's contract](../../0.01-scaffold-and-theme/contracts/cli-commands.md). No network/API contract is introduced — the deployed site is browser-facing HTML only at this point (the `/mcp` and content-negotiation contracts arrive in v0.03).

## `biblio-cli build`

Pre-renders every content-driven page of the current scaffolded project into static output (FR-001, spec Scenario 1).

**Arguments**: none (operates on the current project directory; thin wrapper over `vite build` with React Router v7's `prerender`).

**Output (success)**: complete static output with one rendered page per `content-index.json` entry — list/detail/navigation views included, matching [data-model.md](./../data-model.md) Pre-rendered Page Set completeness invariant.

**Output (failure)**: a clear, per-file error identifying which authored content failed validation (e.g. malformed frontmatter) — the build MUST fail rather than silently produce a page set with gaps (spec Edge Cases).

**Exit codes**: `0` on a complete, gap-free build; non-zero with the failing-file list on any validation failure.

## `biblio-cli deploy`

Publishes a built site to its configured public address from a single command (FR-002, spec Scenario 2).

**Arguments**: none (thin wrapper over `wrangler deploy` plus the generated CI workflow; operates against the route/custom-domain binding already present in `wrangler.jsonc` from v0.01 scaffolding).

**Output (success)**: the site becomes reachable at its configured public address and serves the pre-rendered pages correctly (FR-003); re-running against an already-live site updates cleanly in place with no visitor-facing disruption (spec Edge Cases).

**Output (failure)**: surfaces the underlying `wrangler deploy`/CI failure with enough context to diagnose (e.g. auth, route-binding, or build-artifact issues) — never a silent partial publish.

**Exit codes**: `0` on a successful publish/update; non-zero with the underlying failure surfaced on any deploy error.

## Cross-cutting contract notes

- Both commands remain thin wrappers — `build` over `vite build`, `deploy` over `wrangler deploy` + generated CI — per the confirmed §0.3 command-surface decision that folds the entire create→dev→build→deploy→update lifecycle into one consistent `biblio-cli` surface (Principle III).
- Neither command introduces a new contract surface beyond what `vite`/`wrangler` already define; `biblio-cli`'s contract here is "wrap cleanly, surface failures clearly, never diverge from the underlying tool's behavior."
- The deployed site itself exposes no contract beyond standard browser-facing HTML in this release — no `/mcp` endpoint, no content-negotiation behavior (both arrive in [v0.03](../../0.03-agent-terminal-layer/contracts/)).
