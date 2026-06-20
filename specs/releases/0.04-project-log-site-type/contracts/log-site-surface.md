# Contracts: Project-log site type surface (v0.04)

**Type of contract**: Extension of the existing CLI ([v0.01](../../0.01-scaffold-and-theme/contracts/cli-commands.md)/[v0.02](../../0.02-build-deploy-verify/contracts/cli-commands.md)), MCP tool ([v0.03](../../0.03-agent-terminal-layer/contracts/mcp-tools.md)), and content-negotiation contracts to the project-log content shape. No new commands, endpoints, or response formats are introduced — this document specifies how the *existing* contracts behave for `--type=log`.

## CLI surface (unchanged commands, exercised against `--type=log`)

| Command | Contract for the project-log type |
|---|---|
| `biblio-cli new <name> --type=log ...` | Scaffolds a project-log site using the same flag/TUI contract as [v0.01](../../0.01-scaffold-and-theme/contracts/cli-commands.md) — seeds dated/tagged/project-attributed entry templates and a project-log theme preset (FR-001) |
| `biblio-cli index` | Produces the same canonical `content-index.json` shape, now containing `date`/`tags[]`/`project` fields per entry — no schema fork (FR-001, [data-model.md](../data-model.md) Project-Log Entry) |
| `biblio-cli dev` / `build` / `deploy` | Behave identically to [v0.01](../../0.01-scaffold-and-theme/contracts/cli-commands.md)/[v0.02](../../0.02-build-deploy-verify/contracts/cli-commands.md) — any divergence found here is logged and fixed as a defect against the shared command, never special-cased per type (spec Assumptions) |

## Browser-facing surface (entry list + filters)

| Surface | Contract |
|---|---|
| Entry list | Every published entry appears with correctly attributed `date`/`tags`/`project` (FR-002, Scenario 1) |
| Entry detail | Opens correctly from the list (FR-002) |
| Filter UI | Filters by date range, tag, and project; combining dimensions returns the intersection — no false positives or omissions (FR-003, [data-model.md](../data-model.md) Filter Query) |
| New tag/project | Appears correctly in filters immediately after the next build — no extra registration step (spec Edge Cases) |

## MCP tool surface (extends [v0.03 contract](../../0.03-agent-terminal-layer/contracts/mcp-tools.md))

| Aspect | Contract |
|---|---|
| Type-specific tools | Structured queries appropriate to the project-log content shape (timeline/tag/project queries), discoverable via `list-tools` exactly as the v0.03 contract requires (FR-005) |
| Primary response | Same self-sufficient Markdown/plain-text `content` requirement as every tool (Principle IV) — extended to this content shape, not redefined |
| Type-mismatch rejection | A profile-shaped query (e.g. "get résumé") is rejected with the **same** clear, structured "capability not defined for this site type" error the v0.03 contract already specifies — no bespoke error format (spec Edge Cases, research.md) |

## Terminal/negotiation surface (reuses [v0.03 contract](../../0.03-agent-terminal-layer/contracts/mcp-tools.md) verbatim)

| Aspect | Contract |
|---|---|
| Variant selection, `Vary` header, tier-i/ii rendering | Identical mechanism to v0.03, exercised against project-log page/entry content — any rendering gap here is a platform defect to fix, not a reason to fork the negotiation logic |

## Authoring-rhythm contract (the release's central claim)

| Aspect | Contract |
|---|---|
| "Add a new day's entry" | Uses the *same* authoring flow as the previous entry — author a Markdown+frontmatter file, run `biblio-cli index`/`build`/`deploy` — with zero one-off structural or config adjustments, repeatable for "several consecutive published entries" (FR-004, Scenario 3, SC-002) |
