# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project state

This repository is an early-stage scaffold: it currently contains only tooling
configuration and planning templates, no application source code, build
system, or tests yet. There is nothing to build, lint, or test at this point.
When the project gains actual source code, update this file with the real
commands and architecture.

## Repo layout

- `AGENTS.md` — the canonical project instructions file, auto-loaded by
  opencode (see `opencode.json` -> `instructions`). It is currently a TODO
  template (Project, Conventions, Build/Test/Lint, Subagents, Skills, Notes
  sections are unfilled). Treat it as the source of truth for project
  conventions once it is filled in — extend/update it alongside this file
  rather than letting the two drift apart.
- `opencode.json` — opencode config; points `instructions` at `AGENTS.md`.
- `.opencode/agents/` — project-scoped opencode subagent definitions. One
  markdown file per agent; filename = agent name, frontmatter + system-prompt
  body. See `.opencode/agents/README.md` for the exact file format
  (`description` required; optional `model`, `mode`, `permission`, etc.). New
  agents must also be listed under "Subagents" in `AGENTS.md`.
- `.opencode/` also contains opencode's own `package.json`/`node_modules`
  (the `@opencode-ai/plugin` SDK) — this is tooling for building opencode
  plugins/agents, not the project's application dependencies.
- `specs/research/` — research docs for upcoming work, written from
  `research-template.md` (Problem Statement, Background, Libraries & Stack,
  Architecture Patterns, Existing Codebase Evaluation, Security & Compliance,
  Performance/Reliability/Ops, Risks, Open Questions, Decision Log,
  References). Copy this template when starting research on a new feature or
  technical decision and fill in each section before implementation begins.
- `.claude/settings.json` / `.claude/settings.local.json` — enabled Claude
  Code plugins for this repo: `mcp-apps`, `skill-creator`, `feature-dev`,
  `mcp-server-dev`, `cloudflare`. Their presence signals the likely direction
  of this project (building an MCP server / Cloudflare Workers app with
  interactive UI), but no such code exists yet.

## Working in this repo

- Before writing implementation code, check whether `AGENTS.md` has been
  filled in with real conventions/build commands — it is the file opencode
  loads as project instructions and is meant to be the living guide.
- New features or non-trivial technical decisions should start with a
  research doc under `specs/research/`, based on `research-template.md`.

<!-- SPECKIT START -->
For additional context about technologies to be used, project structure,
shell commands, and other important information, read the current plan:
`specs/releases/0.01-scaffold-and-theme/plan.md`
<!-- SPECKIT END -->
