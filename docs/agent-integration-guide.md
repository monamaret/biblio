# Coding agent integration guide

This repo wires up three layers of agent tooling — opencode's project instructions, Claude Code's Spec Kit workflow, and Claude Code's plugin config. None of it is project *application* code; it's all development-workflow tooling, and this doc explains how the pieces fit together and why each one is there. For project conventions themselves (stack, build/test commands, feature process), see [`AGENTS.md`](../AGENTS.md) — this guide is about the tooling that *enforces* and *automates* that process, not a restatement of it.

## Where to look first

| Question | Look here |
|---|---|
| "What's the project, stack, and build/test commands?" | [`AGENTS.md`](../AGENTS.md) |
| "What's the repo layout and current release state?" | [`CLAUDE.md`](../CLAUDE.md) |
| "How do I run the spec-driven workflow for a new feature?" | This doc, [Spec Kit workflow](#spec-kit-workflow) |
| "How do opencode subagents/skills work in this repo?" | This doc, [opencode](#opencode) |
| "Why is plugin X enabled in Claude Code?" | This doc, [Claude Code plugin config](#claude-code-plugin-config) |

## Spec Kit workflow

This repo runs [Spec Kit](https://github.com/github/spec-kit) (`.specify/`) with Claude Code as its integration (`.specify/integration.json` → `"integration": "claude"`). Two extensions are installed (`.specify/extensions.yml`):

- **`git`** — auto-commits at each workflow step and creates feature branches (`speckit.git.feature`, `speckit.git.commit`, `speckit.git.initialize`, `speckit.git.validate`, `speckit.git.remote`).
- **`agent-context`** — refreshes agent-facing context files after `specify`/`plan` (`speckit.agent-context.update`).

A `github` extension also backs the issue-sync commands (`speckit.taskstoissues`, `speckit.github.close-issues`) referenced in `.specify/extensions.yml`'s hooks, even though it's not a separately versioned entry in `.specify/extensions/.registry`.

Each of these is exposed as a Claude Code skill under `.claude/skills/speckit-*/SKILL.md` — invoke them as slash commands (e.g. `/speckit-specify`). The fixed command sequence, matching `AGENTS.md`'s ["Feature development process"](../AGENTS.md#feature-development-process):

```
/speckit-constitution   # once, to ratify this project's principles (.specify/memory/constitution.md)
/speckit-specify        # turn a feature item into a spec
/speckit-clarify        # resolve ambiguities before planning
/speckit-plan           # produce an implementation plan
/speckit-tasks          # break the plan into tasks
/speckit-implement      # execute the tasks
/speckit-analyze        # cross-check spec/plan/tasks for drift
/speckit-checklist      # generate a review checklist
```

`extensions.yml`'s `hooks` section (`auto_execute_hooks: true`) wires the `git`/`agent-context`/`github` extensions to fire automatically around these steps — e.g. `before_specify` creates a feature branch, every `after_*` step auto-commits, `before_implement`/`after_implement` sync GitHub issues from `tasks.md`. Read `.specify/extensions.yml` directly for the exact hook-by-hook wiring rather than relying on this doc to stay byte-for-byte current with it.

Note `AGENTS.md`'s own "Feature development process" section describes a *fuller* pipeline than Spec Kit alone covers — a numbered research doc under `specs/research/`, a feature item under `specs/features/`, roadmap inclusion, and an architecture sync — all of which happen **before** `/speckit-specify` is ever run. Spec Kit formalizes the spec → plan → tasks → implement steps; it doesn't replace the research/feature-item process that precedes it in this repo.

## opencode

`opencode.json` points opencode's `instructions` at `AGENTS.md` — that's the entire wiring; opencode auto-loads it as project-level context on startup.

`.opencode/agents/` holds project-scoped subagent definitions, one Markdown file per agent (filename = agent name, frontmatter + system-prompt body — see [`.opencode/agents/README.md`](../.opencode/agents/README.md) for the exact frontmatter schema). As of this writing the directory only contains that format-reference `README.md` — no project-specific subagents have been created yet. When one is added, it must also be listed under `AGENTS.md`'s ["Subagents"](../AGENTS.md#subagents) section; the README file format doc and the AGENTS.md listing are two different things and both need updating together.

`.opencode/` also vendors opencode's own `package.json`/tooling for *building* opencode plugins/agents — that's infrastructure for authoring `.opencode/agents/*.md` files, not an application dependency of `biblio-cli` or any scaffolded site.

## Claude Code plugin config

`.claude/settings.json` enables the `mcp-apps` plugin (`claude-plugins-official`) — this repo's actual shape is an MCP-server/Cloudflare-Workers project (every scaffolded site mounts `/mcp`, and this repo itself develops the MCP App UI bundle layer under `templates/site/*/app/mcp/ui/`), so the plugin that understands MCP App resources is the one enabled here, not a generic everything-on plugin list.

`.claude/skills/speckit-*/` (one directory per Spec Kit command) are the Claude Code surface for the workflow described above — they're generated/synced by Spec Kit itself, not hand-authored, so treat changes to their content as a Spec Kit version bump rather than a manual edit.

Anything **not** covered above — personal editor settings, `.claude/settings.local.json`, or similar machine-specific config — doesn't belong in this repo. Local agent configuration is only checked in here when it *is* the shared workflow tooling itself (Spec Kit's `.specify/`/`.claude/skills/speckit-*`, opencode's `.opencode/agents/`); anything user-specific stays out, full stop.
