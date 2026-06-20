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
// the tools appropriate to this docs site (FR-004; spec Scenario 1/3).
describe("MCP server route mounting", () => {
  it("lists shared and docs-type-specific tools", async () => {
    const client = await connectClient();
    const { tools } = await client.listTools();
    const names = tools.map((t) => t.name);

    expect(names).toContain("list_pages");
    expect(names).toContain("get_page");
    expect(names).toContain("search_content");
    expect(names).toContain("list_sections");
    expect(names).toContain("get_doc_page");
  });
});

// A request for a capability the docs site type doesn't define (e.g. a
// profile/log-only tool never registered here) is rejected with the same
// clear, structured error v0.03/v0.04 already specify — not a silent empty
// result or a crash (spec Edge Cases; contracts/docs-site-and-review.md
// "MCP tool surface").
describe("MCP server unsupported-capability handling", () => {
  it("rejects a call to a profile-only tool this site type doesn't define", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "get_profile", arguments: {} });
    expect(result.isError).toBe(true);
  });

  it("rejects a call to a log-only tool this site type doesn't define", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "list_entries", arguments: {} });
    expect(result.isError).toBe(true);
  });

  it("rejects get_doc_page for a slug that doesn't exist with a structured error, not a crash", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "get_doc_page", arguments: { slug: "does-not-exist" } });
    expect(result.isError).toBe(true);
  });
});
