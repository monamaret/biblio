# Quickstart Validation: v0.03 — Make the site usable by AI agents and terminal clients

Validates spec Success Criteria SC-001..SC-004, purely via agent/terminal-client evidence — never re-proving v0.02's browser-facing behavior (spec SC-004). Builds directly on the live pilot site deployed in [v0.02](../0.02-build-deploy-verify/quickstart.md).

## Prerequisites

- The live pilot site from v0.02, deployed with `mcp_enabled: true` in `content/site.yaml`
- An MCP client (e.g. an MCP-capable AI assistant) and a terminal/`curl` for negotiation checks

## Scenario 1 — An MCP client queries the live site (spec Scenario 1)

```text
Connect an MCP client to https://<configured-domain>/mcp
list-tools
call-tool <a tool appropriate to this site's type>
```

**Expected outcome**: the client completes a full connect → list-tools → call-tool exchange on the first attempt, receiving correctly structured results appropriate to the site's type ([contracts/mcp-tools.md](./contracts/mcp-tools.md) Tool contract; SC-001).

## Scenario 2 — A graphical host renders a rich inline view (spec Scenario 2)

```text
Open a tool result that includes a richer inline view in a graphical MCP host.
```

**Expected outcome**: the view renders correctly, inherits the site's `tokens.css` theme, and is accompanied by a clean, self-sufficient text/Markdown rendering for hosts that can't render it ([data-model.md](./data-model.md) MCP Tool entity; Principle IV). Simulate a host whose sandbox can't load theme tokens — expect a usable, if unthemed, fallback view (spec Edge Cases).

## Scenario 3 — A terminal client gets a clean Markdown view (spec Scenario 3)

```bash
curl -H "Accept: text/markdown" https://<configured-domain>/<some-page>
curl -A "curl/8.0" https://<configured-domain>/<some-page>
```

**Expected outcome**: a clean, readable Markdown or ANSI-rendered version of the page, with no leftover browser-page chrome — and a `Vary: Accept, User-Agent` header on the response ([contracts/mcp-tools.md](./contracts/mcp-tools.md) Negotiated page-rendering contract; SC-002). Send an ambiguous/conflicting `Accept` header and confirm the negotiation resolves deterministically (spec Edge Cases).

## Scenario 4 — The agent-facing layer can be turned off (spec Scenario 4)

```bash
# Re-scaffold or reconfigure with mcp_enabled: false, rebuild, redeploy, then:
curl -I https://<configured-domain>/mcp
```

**Expected outcome**: zero reachable trace of the MCP surface — no `/mcp*` route, tool surface, or related generated artifact present (FR-002; [data-model.md](./data-model.md) MCP Server Mount entity; SC-003).

## Demonstrating SC-001..SC-004

- **SC-001**: Scenario 1 completes a full connect → list-tools → call-tool exchange on the first attempt against the live pilot site.
- **SC-002**: Scenario 3's terminal/`curl` requests return clean, readable renderings with no leftover browser-page chrome.
- **SC-003**: Scenario 4 confirms zero reachable trace of the MCP surface when `mcp_enabled` is off, verified by direct request to its endpoint.
- **SC-004**: All of the above is demonstrable purely from agent/terminal-client evidence — none of it requires re-proving v0.02's browser-facing behavior, and v0.02's evidence never required this release's capability either.
