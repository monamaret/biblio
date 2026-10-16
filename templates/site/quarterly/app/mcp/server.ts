import { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";
import { registerSharedTools } from "./tools/shared";
import { registerQuarterlyTools } from "./tools/quarterly";

// createServer builds a fully tool-registered McpServer — the shared
// (cross-site-type) and quarterly-type-specific tools. Deliberately has no
// dependency on `agents/mcp` (McpAgent) or any `cloudflare:`-scheme
// import, so it's directly testable via the MCP SDK's InMemoryTransport
// in plain Node/vitest — `agent.ts` is the thin Workers-runtime wrapper
// around this same server.
//
// Requesting a capability the deployed site type doesn't define (e.g. a
// log-only tool against this quarterly site) is never registered here in
// the first place — the MCP SDK's own `tools/call` handling returns the
// standard JSON-RPC "Unknown tool" error rather than a silent empty result
// or a crash.
export function createServer(): McpServer {
  const server = new McpServer({ name: "biblio-quarterly-site", version: "1.0.0" });
  registerSharedTools(server);
  registerQuarterlyTools(server);
  return server;
}
