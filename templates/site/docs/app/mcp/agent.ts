import { McpAgent } from "agents/mcp";
import { createServer } from "./server";

export type Env = {
  MyMCP: DurableObjectNamespace;
};

// MyMCP is the McpAgent subclass that mounts `/mcp*` on the deployed
// Worker (FR-001; data-model.md "MCP Server Mount") — only ever present
// when `mcp_enabled` is true, per internal/scaffold/generator.go's
// conditional wrangler.jsonc wiring. This file (unlike server.ts) depends
// on the real Workers/Durable-Object runtime via `agents/mcp` and is never
// imported by anything under vitest.
export class MyMCP extends McpAgent<Env, unknown, {}> {
  server = createServer();

  async init() {
    // Tool registration already happened in createServer(); no further
    // per-session setup is needed.
  }
}

// mcpHandler is mounted onto /mcp* by the Worker entry script
// (workers/entry.ts) rather than owning the whole top-level fetch itself —
// the entry script also needs to route negotiation/static-asset requests.
export const mcpHandler = MyMCP.serve("/mcp", { binding: "MyMCP" });
