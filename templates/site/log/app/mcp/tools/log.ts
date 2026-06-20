import type { McpServer } from "@modelcontextprotocol/sdk/server/mcp.js";

// registerLogTools is a no-op placeholder for now — T031/T032 (Phase 6)
// implement the real project-log-specific tools (`list_entries`,
// `get_entry`) and fill in this function's body. Kept here, rather than
// left unimplemented, so `templates/site/log/app/mcp/server.ts`'s
// `createServer()` (T010) has something real to import and call, and the
// log site builds/typechecks cleanly at every checkpoint between now and
// Phase 6.
export function registerLogTools(_server: McpServer) {}
