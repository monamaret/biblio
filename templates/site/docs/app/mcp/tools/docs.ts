import { z } from "zod";
import type { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { listEntries as queryEntries, getEntryBySlug } from "../query";
import { toc } from "../../content/toc";
import type { ToolResult } from "./shared";

function textResult(text: string): ToolResult {
  return { content: [{ type: "text", text }] };
}

// listSections is docs-type-specific: every section and its pages, in
// hierarchy order — a complete, self-sufficient answer on its own
// (FR-002, FR-004; data-model.md "Table of Contents / Navigation
// Structure"), backed by the same toc() the browser TOC view uses
// (Principle II — one TOC-computation owner).
//
// This site type has no filter-query surface at all (Scenario 1.2) — a
// log-shaped filter argument (e.g. `tag`/`project`/date-range, which only
// make sense for the project-log type) is rejected with a clear,
// structured error rather than silently ignored or accepted as a no-op
// (spec Edge Cases; contracts/docs-site-and-review.md "MCP tool surface" —
// the same documented-default contract as v0.04's type-mismatch handling).
export function listSections(query?: { tag?: string; project?: string; dateFrom?: string; dateTo?: string }): ToolResult {
  if (query && (query.tag !== undefined || query.project !== undefined || query.dateFrom !== undefined || query.dateTo !== undefined)) {
    throw new Error(
      "list_sections: this site type has no tag/project/date filtering — received a log-shaped filter argument",
    );
  }
  const sections = toc(queryEntries());
  if (sections.length === 0) {
    return textResult("No pages are authored yet.");
  }
  const lines: string[] = [];
  for (const section of sections) {
    lines.push(`## ${section.section}`);
    for (const page of section.pages) {
      const title = typeof page.title === "string" ? page.title : page.slug;
      lines.push(`- ${title} (/${page.slug})`);
    }
  }
  return textResult(lines.join("\n"));
}

// getDocPage is docs-type-specific: the full Markdown body of one page by
// slug. Named distinctly from the shared `get_page` tool
// (templates/site/docs/app/mcp/tools/shared.ts, registered on every site
// type) to avoid a duplicate-tool-name registration collision — same
// naming precedent as v0.04's `get_entry`. Rejects a slug that doesn't
// match any page with a clear, structured error rather than a silent
// empty result (spec Edge Cases).
export function getDocPage(slug: string): ToolResult {
  const entry = getEntryBySlug(slug);
  if (!entry) {
    throw new Error(`get_doc_page: no page found for slug "${slug}"`);
  }
  const title = typeof entry.title === "string" ? entry.title : entry.slug;
  return textResult(`# ${title}\n\n${entry.body}`);
}

// registerDocsTools wires every docs-type-specific tool onto the McpServer
// instance (contracts/docs-site-and-review.md "MCP tool surface"). A
// request for a capability this site type doesn't define (e.g. a
// profile/log-only tool) simply isn't registered here — the server
// surfaces MCP's standard "unknown tool" error rather than a silent empty
// result (spec Edge Cases).
export function registerDocsTools(server: McpServer) {
  server.registerTool(
    "list_sections",
    {
      description: "List every documentation section and its pages, in hierarchy order",
      inputSchema: {
        tag: z.string().optional(),
        project: z.string().optional(),
        dateFrom: z.string().optional(),
        dateTo: z.string().optional(),
      },
    },
    async (query) => listSections(query),
  );

  server.registerTool(
    "get_doc_page",
    { description: "Get the full content of a documentation page by its slug", inputSchema: { slug: z.string() } },
    async ({ slug }: { slug: string }) => getDocPage(slug),
  );
}
