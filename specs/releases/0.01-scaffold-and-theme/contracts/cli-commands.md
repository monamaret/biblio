# Contracts: `biblio-cli` command surface (v0.01)

**Type of contract**: CLI command/flag schema — the interface this release exposes to its single user (the author). Per the plan's Project Type (CLI tool), this is the appropriate contract form; there is no network API in this release (the MCP/`fetch` surface doesn't exist until v0.03).

## `biblio-cli new <project-name>`

Scaffolds a complete, runnable combined site project (FR-001, FR-002).

**Flags**:
| Flag | Required | Values | Behavior when omitted |
|---|---|---|---|
| `--type` | no | `profile` \| `log` \| `docs` | Bubble Tea TUI prompts interactively |
| `--domain` | no | domain string | TUI prompts interactively |
| `--theme` | no | preset id matching `--type` | TUI prompts interactively, offering only presets valid for the chosen `--type` |

**Output (success)**: a populated `<project-name>/` directory containing `content/` (seeded per `--type`), `content/site.yaml`, `wrangler.jsonc`, and `tokens.css` — matching spec Scenario 1, acceptance criteria 1–2.

**Output (failure)**: a clear, actionable error identifying which input was invalid (e.g., unknown `--type`, a `--theme` that doesn't match the chosen `--type`'s presets, or a `project-name` that collides with an existing directory) — never a partially-written project left behind silently.

**Exit codes**: `0` success; non-zero with a human-readable message on any validation or generation failure.

## `biblio-cli index [path]`

Produces the single canonical content index from an authored `content/` directory (FR-003, spec Scenario 2).

**Arguments**: optional `path` (defaults to the current project's `content/` directory).

**Output (success)**: writes `content-index.json` reflecting every authored page/entry with its derived `slug`/`path` and structured fields, in the shape contract defined in [data-model.md](./../data-model.md) (Content Index entity).

**Output (failure)**: surfaces malformed or missing structured frontmatter fields as a clear, per-file error — does not silently emit an incomplete index (spec Edge Cases).

**Exit codes**: `0` success with index written; non-zero if any content item fails schema validation, listing every failing file (not just the first).

## `biblio-cli dev`

Starts local preview of a scaffolded project from one command (FR-005, spec Scenario 3).

**Arguments**: none (operates on the current project directory).

**Behavior**: wraps and supervises `vite dev` + `wrangler dev` as one supervised process pair; surfaces both servers' addresses and forwards their output so the author can preview content live.

**Output (failure)**: if either underlying process fails to start (e.g., port collision, missing `wrangler.jsonc`), reports which one and why — does not leave an orphaned half-started dev session.

**Exit codes**: runs until interrupted (`Ctrl+C`); `0` on clean shutdown of both wrapped processes.

## Cross-cutting contract notes

- All three commands are reachable without flags via the Bubble Tea-driven guided flow (Principle III) — flags are the scripting/repeat-use path, not the only path.
- `biblio-cli` itself is distributed as a single binary (GitHub Releases via `goreleaser`, or `go install`) — there is no contract surface beyond this CLI in v0.01 (FR-006).
