import { z } from "zod";
import type { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { listEntries, getEntryBySlug, searchEntries } from "../query";
import { renderChartView } from "../ui/chart-view";

// ToolResult mirrors the MCP SDK's CallToolResult content-block shape this
// release relies on: `content` is always a self-sufficient Markdown/
// plain-text rendering (Principle IV, FR-003) — `content` is what every
// host sees regardless of richer-view support; an embedded `resource`
// content block carries the optional `ui-resource` (data-model.md "MCP
// Tool" `ui-resource`) and is strictly additive, never a substitute.
export type ToolResult = {
  content: Array<
    | { type: "text"; text: string }
    | { type: "resource"; resource: { uri: string; mimeType: string; text: string } }
  >;
};

function textResult(text: string): ToolResult {
  return { content: [{ type: "text", text }] };
}

// listPages is a shared, cross-site-type tool: every authored entry's
// slug/title as a flat Markdown list — a complete, self-sufficient answer
// on its own (FR-003).
export function listPages(): ToolResult {
  const entries = listEntries();
  if (entries.length === 0) {
    return textResult("No pages are authored yet.");
  }
  const lines = entries.map(
    (entry) => `- ${typeof entry.title === "string" ? entry.title : entry.slug} (/${entry.slug})`,
  );
  return textResult(lines.join("\n"));
}

// getPage is a shared tool: the full Markdown body of one entry by slug.
// Rejects a slug that doesn't match any entry with a clear, structured
// error rather than a silent empty result (spec Edge Cases).
export function getPage(slug: string): ToolResult {
  const entry = getEntryBySlug(slug);
  if (!entry) {
    throw new Error(`get_page: no page found for slug "${slug}"`);
  }
  const title = typeof entry.title === "string" ? entry.title : entry.slug;
  return textResult(`# ${title}\n\n${entry.body}`);
}

// searchContent is a shared tool: keyword search across every entry's
// title/summary/body (Principle V — keyword search only, still
// pre-Vectorize). With a richer view attached for graphical hosts
// (data-model.md "MCP Tool" `ui-resource`) — the text answer is always
// present alongside it, never replaced by it (FR-003, Scenario 2).
export function searchContent(query: string): ToolResult {
  const matches = searchEntries(query);
  if (matches.length === 0) {
    return textResult(`No pages matched "${query}".`);
  }
  const lines = matches.map(
    (entry) => `- ${typeof entry.title === "string" ? entry.title : entry.slug} (/${entry.slug})`,
  );
  const text = `${matches.length} page(s) matched "${query}":\n${lines.join("\n")}`;

  const chartHtml = renderChartView(
    matches.map((entry) => ({
      label: typeof entry.title === "string" ? entry.title : entry.slug,
      value: 1,
    })),
  );

  return {
    content: [
      { type: "text", text },
      {
        type: "resource",
        resource: { uri: `ui://chart-view/search/${encodeURIComponent(query)}`, mimeType: "text/html", text: chartHtml },
      },
    ],
  };
}

// registerSharedTools wires every shared tool onto the McpServer instance
// (contracts/mcp-tools.md "Tool contract").
export function registerSharedTools(server: McpServer) {
  server.registerTool(
    "list_pages",
    { description: "List every authored page on this site" },
    async () => listPages(),
  );

  server.registerTool(
    "get_page",
    { description: "Get the full content of a page by its slug", inputSchema: { slug: z.string() } },
    async ({ slug }: { slug: string }) => getPage(slug),
  );

  server.registerTool(
    "search_content",
    { description: "Keyword-search every authored page's title, summary, and body", inputSchema: { query: z.string() } },
    async ({ query }: { query: string }) => searchContent(query),
  );
}
