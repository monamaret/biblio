import { describe, expect, it } from "vitest";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { InMemoryTransport } from "@modelcontextprotocol/sdk/inMemory.js";
import { createServer } from "./server";

// connectClient links a real MCP Client to the real McpServer the deployed
// site mounts on /mcp* via InMemoryTransport — a genuine connect/list/call
// exchange at the protocol level, not a mock.
async function connectClient() {
  const server = createServer();
  const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair();
  const client = new Client({ name: "test-client", version: "1.0.0" });
  await Promise.all([client.connect(clientTransport), server.connect(serverTransport)]);
  return client;
}

// /mcp* mounting succeeds when mcp_enabled is true: a real client can list
// the tools this quarterly site serves (specs/004-quarterly-site-type
// FR-032; quarterly-specific tools are added in US5, tasks.md T082).
describe("MCP server route mounting", () => {
  it("lists the shared tools and the four quarterly tools", async () => {
    const client = await connectClient();
    const { tools } = await client.listTools();
    const names = tools.map((t) => t.name);

    expect(names).toEqual(
      expect.arrayContaining(["list_pages", "get_page", "search_content", "list_issues", "get_issue", "list_articles", "get_article"]),
    );
    expect(tools.find((t) => t.name === "get_article")?.outputSchema).toBeDefined();
  });

  it("serves get_article's structured notes over the protocol", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "get_article", arguments: { slug: "issues/02-autumn-2026/on-margins" } });
    expect(result.isError).toBeFalsy();
    expect((result.structuredContent as { marginNotes: unknown[] }).marginNotes).toHaveLength(3);
  });

  it("answers an unknown issue in text rather than failing", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "get_issue", arguments: { issue: "7" } });
    expect(result.isError).toBeFalsy();
    expect(JSON.stringify(result.content)).toContain("No issue matches");
  });
});

// A request for a capability this site type doesn't define (a log-only
// tool) is rejected with a structured error, not a silent empty result.
describe("MCP server unsupported-capability handling", () => {
  it("rejects a call to a log-only tool this site type doesn't define", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "list_entries", arguments: {} });
    expect(result.isError).toBe(true);
  });
});
