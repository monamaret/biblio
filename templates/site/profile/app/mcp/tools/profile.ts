import type { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { listEntries, getEntryBySlug } from "../query";
import type { ToolResult } from "./shared";

function textResult(text: string): ToolResult {
  return { content: [{ type: "text", text }] };
}

// getProfile is profile-type-specific: the author's "about" page — the
// capability this site type defines that a log/docs site wouldn't
// (data-model.md "MCP Tool" `scope: type-specific`).
export function getProfile(): ToolResult {
  const entry = getEntryBySlug("about");
  if (!entry) {
    throw new Error("get_profile: no \"about\" page is authored for this profile site");
  }
  const title = typeof entry.title === "string" ? entry.title : entry.slug;
  return textResult(`# ${title}\n\n${entry.body}`);
}

// listProjects is profile-type-specific: every entry under the
// "projects/" slug prefix, as a flat Markdown list.
export function listProjects(): ToolResult {
  const projects = listEntries().filter((entry) => entry.slug.startsWith("projects/"));
  if (projects.length === 0) {
    return textResult("No projects are authored yet.");
  }
  const lines = projects.map(
    (entry) => `- ${typeof entry.title === "string" ? entry.title : entry.slug} (/${entry.slug})`,
  );
  return textResult(lines.join("\n"));
}

// registerProfileTools wires every profile-type-specific tool onto the
// McpServer instance (contracts/mcp-tools.md "Tool contract"). A request
// for a capability this site type doesn't define (e.g. a log/docs-only
// tool) simply isn't registered here — the server surfaces MCP's standard
// "unknown tool" error rather than a silent empty result (spec Edge Cases).
export function registerProfileTools(server: McpServer) {
  server.registerTool(
    "get_profile",
    { description: "Get the author's profile/about page" },
    async () => getProfile(),
  );

  server.registerTool(
    "list_projects",
    { description: "List every project on this profile site" },
    async () => listProjects(),
  );
}
