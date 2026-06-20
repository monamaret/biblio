# v0.03 — Make the site usable by AI agents and terminal clients

**Created**: 2026-06-07
**Status**: Draft
**Roadmap deliverables**: [§4 In-Worker MCP server](../research/roadmap.md#4-in-worker-mcp-server) · [§5 MCP App UI layer](../research/roadmap.md#5-mcp-app-ui-layer) · [§7 Terminal/CLI-optimized view](../research/roadmap.md#7-terminalcli-optimized-view-of-the-site)
**Depends on**: [v0.02](0.02-release.md) — needs a real, deployed site to attach the agent-facing layer to
**Feature specs**: [004](../features/004-in-worker-mcp-server.md) · [005](../features/005-mcp-app-ui-layer.md) · [007](../features/007-terminal-cli-view.md)

## Summary

The author's deployed site also becomes a first-class destination for AI assistants and terminal/CLI tools — answering structured questions through a built-in MCP server, rendering rich inline views where supported, and degrading to a clean, readable view for `curl`/terminal clients — so the same deployment serves both human visitors and agent/CLI consumers. The roadmap places this layer only after a real site exists to attach it to (v0.02), and notes every blocking research spike for the terminal view is already resolved, making it execution-ready.

## User Scenarios & Acceptance Criteria

### Scenario 1 — An MCP client queries the live site

As an AI assistant connected via MCP, I can discover and call the site's tools and get correctly structured answers.

1. **Given** a deployed site with the agent-facing capability enabled, **When** an MCP client connects to its MCP endpoint, **Then** the client can list and successfully call the tools appropriate to that site's type and receive correctly structured results.

### Scenario 2 — A graphical host renders a rich inline view

As a user of a graphical MCP host, I see a rich inline view for chart-like results, themed consistently with the site — and a host that can't render it still gets a complete text answer.

1. **Given** a tool result that includes a richer inline view, **When** it is opened in a graphical MCP host, **Then** the view renders correctly, inherits the site's visual theme, and is also accompanied by a clean, self-sufficient text/Markdown rendering for hosts that can't render it.

### Scenario 3 — A terminal client gets a clean Markdown view

As someone using `curl` or a terminal-based reader, I get a readable Markdown/ANSI rendering instead of full browser HTML.

1. **Given** the same deployed site, **When** a terminal or `curl`-based client requests a page while indicating it prefers Markdown, **Then** the response is a clean, readable Markdown or ANSI-rendered version of that page rather than full browser-oriented HTML.

### Scenario 4 — The agent-facing layer can be turned off

As the author, if I disable the agent-facing capability for a site, no trace of it should remain reachable.

1. **Given** a site where the agent-facing capability has been turned off in its configuration, **When** any client requests its MCP endpoint, **Then** no MCP route, tool surface, or related generated artifact is present.

## Functional Requirements

- **FR-001**: The system MUST expose, on a deployed site whose agent-facing setting is enabled, a structured query interface (MCP server) offering both the capabilities shared by every site type and the capabilities specific to that site's type.
- **FR-002**: The system MUST ensure that when a site's agent-facing setting is disabled, no related route, generated artifact, or tool surface is present or reachable on that site.
- **FR-003**: The system MUST, for any tool result that includes a richer inline view, also produce a clean, self-sufficient text/Markdown rendering as the primary result — so agent hosts that can't render the richer view still get a complete answer.
- **FR-004**: The system MUST detect, for any page request, whether the requesting client prefers a Markdown/terminal-friendly rendering or a standard browser rendering, and MUST return the appropriate one — including a graceful, lightly-formatted rendering for terminal clients that can display ANSI styling.
- **FR-005**: The system MUST make every page discoverable in its Markdown-friendly form (e.g., via sibling addresses and/or discoverability hints) and MUST publish a minimal machine-readable summary of the site's purpose and key pages for agent/crawler consumption.

## Edge Cases

- What happens when an MCP client or terminal client requests a page that doesn't exist, or requests a capability that the site's type doesn't define?
- How does the negotiation logic behave for a client that sends ambiguous or conflicting signals (e.g., an `Accept` header that doesn't clearly prefer either Markdown or HTML)?
- What happens when a graphical host opens a rich inline view but its sandbox can't load the theme tokens — does it still fall back to a usable, if unthemed, view?

## Success Criteria

- **SC-001**: A real MCP client can complete a full connect → list-tools → call-tool exchange against the live pilot site on the first attempt.
- **SC-002**: A terminal/`curl` request for the live site returns a clean, readable rendering with no leftover browser-page chrome.
- **SC-003**: Turning the agent-facing setting off for a site results in zero reachable trace of the MCP surface, verified by direct request to its endpoint.
- **SC-004**: This release is independently demonstrable purely on agent/terminal-client evidence — showing it never requires re-proving v0.02's browser-facing behavior, and vice versa.

## Assumptions

- **Single pilot site remains the proving ground**: This release is demonstrated against the same one real site v0.02 deployed — proving the agent-facing layer on the other two site types is deferred to v0.04/v0.05 alongside their own replication.
- **No timeline commitments**: Consistent with the roadmap's explicit choice to avoid effort/timeline estimates, this spec defines what the release contains and how its completion is demonstrated, not when it ships.
- **Research spikes already resolved**: The roadmap records all four blocking research spikes for the terminal/CLI view (§0.2) as resolved — this release is treated as a fully-specified execution task, not one with open design questions of its own.
