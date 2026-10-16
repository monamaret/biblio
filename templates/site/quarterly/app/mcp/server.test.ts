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
  it("lists the shared tools", async () => {
    const client = await connectClient();
    const { tools } = await client.listTools();
    const names = tools.map((t) => t.name);

    expect(names).toContain("list_pages");
    expect(names).toContain("get_page");
    expect(names).toContain("search_content");
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
