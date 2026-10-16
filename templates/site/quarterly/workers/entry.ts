import { mcpHandler } from "../app/mcp/agent";
import { negotiatePageResponse } from "../app/negotiation/accept";

// The Worker `main` entry, mounted only when `mcp_enabled` is true
// (internal/scaffold/generator.go's `renderWranglerJSONC`). Routes `/mcp*`
// to the in-Worker MCP server, negotiates page requests between HTML and
// Markdown/ANSI, and falls through to Workers Static Assets for
// everything else (Principle V — one Worker, one topology). When
// `mcp_enabled` is false, wrangler.jsonc omits `main` entirely — this file
// is never bundled or invoked at all, which is the actual "zero reachable
// trace" guarantee (FR-002), verified Go-side in
// internal/scaffold/generator_test.go rather than here, since there is no
// live code path to exercise once this script isn't part of the deployed
// artifact in the first place.
export { MyMCP } from "../app/mcp/agent";

export type Env = {
  ASSETS: Fetcher;
  MyMCP: DurableObjectNamespace;
};

export default {
  async fetch(request: Request, env: Env, ctx: ExecutionContext): Promise<Response> {
    const url = new URL(request.url);

    if (url.pathname === "/mcp" || url.pathname.startsWith("/mcp/")) {
      return mcpHandler.fetch(request, env, ctx);
    }

    const slug = url.pathname.replace(/^\/+/, "").replace(/\/+$/, "") || "";
    const negotiated = await negotiatePageResponse(
      slug,
      request.headers.get("Accept"),
      request.headers.get("User-Agent"),
    );
    if (negotiated) {
      return negotiated;
    }

    return env.ASSETS.fetch(request);
  },
} satisfies ExportedHandler<Env>;
