import { defineConfig } from "vite";
import { reactRouter } from "@react-router/dev/vite";
import { llmsTxtPlugin } from "./app/discovery/llms-txt";

// @cloudflare/vite-plugin is deliberately not used here. It was tried again
// for v0.03 (now that wrangler.jsonc has a real Worker `main` entry,
// workers/entry.ts, for the MCP server) but it claims the same Vite "ssr"
// environment React Router's own prerender step uses, pulling
// workers/entry.ts's dependencies (agents, md4x, mimetext — ~2MB) into
// what should be a small, prerender-only SSR bundle, and breaking the
// build entirely ("Server build file not found in manifest"). Removing it
// fixes both: `biblio-cli build` (`vite build`) only needs React Router's
// own prerender step, and `biblio-cli dev`/`deploy` drive Wrangler
// directly as separate processes (see v0.02's same finding in
// internal/devserver), so this plugin's dev-time Workers emulation was
// never actually load-bearing here.
export default defineConfig({
  plugins: [reactRouter(), llmsTxtPlugin()],
});
