import { z } from "zod";
import type { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { listEntries as queryEntries, getEntryBySlug } from "../query";
import { filterEntries, type FilterQuery } from "../../content/filters";
import type { ToolResult } from "./shared";

function textResult(text: string): ToolResult {
  return { content: [{ type: "text", text }] };
}

function entryLine(entry: ReturnType<typeof queryEntries>[number]): string {
  const title = typeof entry.title === "string" ? entry.title : entry.slug;
  const date = typeof entry.date === "string" ? entry.date : undefined;
  const suffix = date ? ` — ${date}` : "";
  return `- ${title}${suffix} (/${entry.slug})`;
}

// listEntries is log-type-specific: every entry matching the given
// date-range/tag/project FilterQuery, reusing the same `filterEntries()`
// the browser filter UI uses (Principle II — one filter-computation
// owner; FR-005, data-model.md "Filter Query").
export function listEntries(query: FilterQuery): ToolResult {
  const matches = filterEntries(queryEntries(), query);
  if (matches.length === 0) {
    return textResult("No entries matched the given filter.");
  }
  return textResult(matches.map(entryLine).join("\n"));
}

// getEntry is log-type-specific: the full Markdown body of one entry by
// slug. Rejects a slug that doesn't match any entry with a clear,
// structured error rather than a silent empty result (spec Edge Cases).
export function getEntry(slug: string): ToolResult {
  const entry = getEntryBySlug(slug);
  if (!entry) {
    throw new Error(`get_entry: no entry found for slug "${slug}"`);
  }
  const title = typeof entry.title === "string" ? entry.title : entry.slug;
  return textResult(`# ${title}\n\n${entry.body}`);
}

// registerLogTools wires every log-type-specific tool onto the McpServer
// instance (contracts/log-site-surface.md "MCP tool surface"). A request
// for a capability this site type doesn't define (e.g. a profile-only
// tool) simply isn't registered here — the server surfaces MCP's standard
// "unknown tool" error rather than a silent empty result (spec Edge Cases).
export function registerLogTools(server: McpServer) {
  server.registerTool(
    "list_entries",
    {
      description: "List project-log entries, optionally filtered by date range, tag, or project",
      inputSchema: {
        dateFrom: z.string().optional(),
        dateTo: z.string().optional(),
        tag: z.string().optional(),
        project: z.string().optional(),
      },
    },
    async (query: FilterQuery) => listEntries(query),
  );

  server.registerTool(
    "get_entry",
    { description: "Get the full content of a project-log entry by its slug", inputSchema: { slug: z.string() } },
    async ({ slug }: { slug: string }) => getEntry(slug),
  );
}
