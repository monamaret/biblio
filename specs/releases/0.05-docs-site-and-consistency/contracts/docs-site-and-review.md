# Contracts: Documentation site type surface + cross-type review (v0.05)

**Type of contract**: Extension of the existing CLI/MCP/negotiation contracts ([v0.01](../../0.01-scaffold-and-theme/contracts/cli-commands.md)–[v0.04](../../0.04-project-log-site-type/contracts/log-site-surface.md)) to the documentation content shape, plus a new structured-review contract for the cross-type consistency check — the first artifact of its kind in the release sequence.

## CLI surface (unchanged commands, exercised against `--type=docs`)

| Command | Contract for the documentation type |
|---|---|
| `biblio-cli new <name> --type=docs ...` | Scaffolds a documentation site using the same flag/TUI contract as every other type — seeds hierarchical (`section`/`order`) page templates and a documentation theme preset (FR-001) |
| `biblio-cli index` | Produces the same canonical `content-index.json` shape, now containing `section`/`order` fields per page — no schema fork, and explicitly **no** `tags[]`/`project` fields ([data-model.md](../data-model.md) Documentation Page) |
| `dev` / `build` / `deploy` | Identical contract to every other type; divergences are platform defects, never special-cased |

## Browser-facing surface (TOC/hierarchy navigation — deliberately filter-free)

| Surface | Contract |
|---|---|
| TOC/navigation | Correctly reflects `section`/`order` structure; every page reachable and renders correctly (FR-002, Scenario 1.1) |
| Filter UI | **Must be absent** — no tag/project filter controls or query params anywhere on this site type (FR-002, Scenario 1.2); a test suite (`test/unit/no-filter-surface/`) asserts this absence as a first-class, required contract, not an unverified assumption |
| Ambiguous section/order mapping | Resolved via a documented, sensible default (e.g. primary-section assignment) — never an undefined or duplicated placement (spec Edge Cases) |

## MCP tool surface (extends [v0.03](../../0.03-agent-terminal-layer/contracts/mcp-tools.md)/[v0.04](../../0.04-project-log-site-type/contracts/log-site-surface.md))

| Aspect | Contract |
|---|---|
| Type-specific tools | Structured queries appropriate to the hierarchical documentation shape (e.g. "list sections," "get page by path"), discoverable via `list-tools` |
| Primary response | Same self-sufficient Markdown/plain-text `content` requirement as every tool (Principle IV) |
| Log-shaped filter param received | Handled via the **same documented, sensible default** as a v0.04-style type mismatch — clearly rejected or cleanly ignored, never silently mishandled (spec Edge Cases, research.md) |

## Terminal/negotiation surface (reuses [v0.03](../../0.03-agent-terminal-layer/contracts/mcp-tools.md) verbatim)

| Aspect | Contract |
|---|---|
| Variant selection, `Vary` header, tier-i/ii rendering | Identical mechanism, exercised against documentation content — returns a rendering "appropriate to its own content shape" (FR-004), not a forced reuse of the project-log or profile templates |

## Cross-Type Consistency Review contract (new — closing artifact of this release)

| Aspect | Contract |
|---|---|
| Scope | Exactly the three live deployments: profile (v0.02), project log (v0.04), documentation (this release) — run only once all three exist (research.md "why documentation is replicated last") |
| Theming check | Each site displays its own distinct, correctly-applied theme; **zero** cross-type bleed-through or inconsistency (FR-003, Scenario 2, SC-002) |
| Terminal-rendering check | Each site returns a clean, correct, type-appropriate terminal/Markdown view on request (FR-004, Scenario 3, SC-002) |
| Inconsistency triage | Any finding is filed as a defect against the *shared* `internal/theme`/negotiation code and fixed there — never patched as a per-site override (spec Edge Cases; [data-model.md](../data-model.md) Cross-Type Consistency Record) |
| Output | A structured review record ([data-model.md](../data-model.md) Cross-Type Consistency Record) showing zero open inconsistencies at close — the artifact SC-002 is checked against |
