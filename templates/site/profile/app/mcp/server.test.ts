import { describe, expect, it } from "vitest";
import { Client } from "@modelcontextprotocol/sdk/client/index.js";
import { InMemoryTransport } from "@modelcontextprotocol/sdk/inMemory.js";
import { createServer } from "./server";

// connectClient links a real MCP Client to the real McpServer the deployed
// site mounts on /mcp* via InMemoryTransport — a genuine connect/list/call
// exchange at the protocol level, not a mock (matches the real exchange
// validated against the actual Workers runtime in research spike 002).
async function connectClient() {
  const server = createServer();
  const [clientTransport, serverTransport] = InMemoryTransport.createLinkedPair();
  const client = new Client({ name: "test-client", version: "1.0.0" });
  await Promise.all([client.connect(clientTransport), server.connect(serverTransport)]);
  return client;
}

// /mcp* mounting succeeds when mcp_enabled is true: a real client can list
// the tools appropriate to this profile site (FR-001; spec Scenario 1).
describe("MCP server route mounting", () => {
  it("lists shared and profile-type-specific tools", async () => {
    const client = await connectClient();
    const { tools } = await client.listTools();
    const names = tools.map((t) => t.name);

    expect(names).toContain("list_pages");
    expect(names).toContain("get_page");
    expect(names).toContain("search_content");
    expect(names).toContain("get_profile");
    expect(names).toContain("list_projects");
  });
});

// A request for a capability the profile site type doesn't define (e.g. a
// log/docs-only tool never registered here) is rejected with a clear,
// structured error — not a silent empty result or a crash (spec Edge
// Cases; contracts/mcp-tools.md "Failure-mode contract").
describe("MCP server unsupported-capability handling", () => {
  it("rejects a call to a tool this site type doesn't define", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "list_log_entries", arguments: {} });
    expect(result.isError).toBe(true);
  });

  it("rejects get_page for a slug that doesn't exist with a structured error, not a crash", async () => {
    const client = await connectClient();
    const result = await client.callTool({ name: "get_page", arguments: { slug: "does-not-exist" } });
    expect(result.isError).toBe(true);
  });
});
