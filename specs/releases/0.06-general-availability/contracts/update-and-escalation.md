# Contracts: `biblio-cli update`/`upgrade` + storage-escalation guidance (v0.06)

**Type of contract**: CLI command/flag schema (extends [v0.01](../../0.01-scaffold-and-theme/contracts/cli-commands.md)/[v0.02](../../0.02-build-deploy-verify/contracts/cli-commands.md)) plus a documentation/guidance contract for the conditional storage-escalation paths. No network/API or MCP contract changes — this release operates entirely at the CLI/scaffold-management layer (Constitution Check: Principle IV is N/A here).

## `biblio-cli update` / `biblio-cli upgrade`

Carries an existing deployed site forward to a newer scaffold version (FR-001, Scenario 1). `update` advances to the latest compatible version; `upgrade` targets a specific version — both share the same underlying classify/merge/conflict machinery (`internal/migrate`).

**Arguments**: operates on the current project directory; `upgrade` additionally accepts a target-version identifier.

**Behavior contract**:

| Step | Contract |
|---|---|
| 1. Classify | Every region of the existing site is classified as `scaffold-tracked`, `author-customized`, or `ambiguous` against the canonical scaffold-template definitions (Principle II) — *before* any write occurs ([data-model.md](../data-model.md) Scaffold Region Classification) |
| 2. Advance | `scaffold-tracked` regions are brought up to date to the target version |
| 3. Preserve | `author-customized` regions are left untouched, byte-for-byte |
| 4. Surface | `ambiguous` regions become Update Conflicts, presented to the author with both the scaffold-proposed and author-current content — the flow waits for an explicit resolution and never auto-resolves (FR-002) |
| 5. Commit | Only after all conflicts are resolved does the flow apply changes; if interrupted at any step, the site remains in a consistent, recoverable state (spec Edge Cases) |

**Output (success)**: a report of what was preserved, what was advanced, and how each conflict was resolved — the triple [data-model.md](../data-model.md)'s Update Run entity defines as a successful run.

**Output (failure / interruption)**: a clear statement of the site's current (consistent, recoverable) state and how to resume or roll back — never a half-migrated site with no documented path forward.

**Exit codes**: `0` on a clean completion with all conflicts resolved; non-zero (with state-recovery guidance) on any interruption or unresolved-conflict halt.

**Cross-type requirement**: this exact contract MUST hold, unmodified, when run against all three real site types (profile, project log, documentation) — FR-003, Scenario 2, validated by `test/integration/update-three-sites/`.

## Storage-Escalation Guidance (`docs/storage-escalation-guide.md`)

A documentation contract, not a runtime command — but one this release MUST prove against real data (FR-004, FR-005, Scenario 3).

| Aspect | Contract |
|---|---|
| Thresholds documented | Concrete, checkable thresholds for both conditional escalation paths: D1 (storage volume / relational-query needs outgrowing the JSON index) and Vectorize (search *quality*, not volume, demonstrably failing on conceptual queries) — the exact triggers Principle V names |
| Real-scale check | The guidance, run against each of the three real sites' actual `content-index.json`-derived scale, MUST correctly report whether any site has crossed a threshold (FR-004) |
| Executable-on-demand proof | For each path the guidance recommends, its steps MUST be concrete enough to actually carry out — proven via an actual dry-run/walkthrough against (a copy of) real site data, not merely described in principle (FR-005) |
| What GA does NOT require | Actually carrying out either escalation absent a real crossed threshold — doing so would itself violate Principle V's "not something to pre-build just in case" rule (research.md "What GA means") |

## GA Readiness Review (closing artifact, FR-006, Scenario 4)

| Aspect | Contract |
|---|---|
| Scope | Every roadmap deliverable §1–§10, cross-referenced against this release's acceptance bar |
| Pass condition | Every ordinary deliverable is `complete-and-demonstrated`; the two conditional storage-escalation items specifically are `proven-executable-on-demand` (research.md "What GA means", [data-model.md](../data-model.md) GA Readiness Review) |
| Independence | MUST reference, not re-run, v0.01–v0.05's already-closed evidence — SC-003 requires this release's evidence to stand on its own without re-proving prior releases' site-type behaviors, and vice versa |
| Output | A `ga-determination: true/false` record — `true` is the literal, checkable definition of "the platform is Generally Available" this release establishes |
