import { defineConfig } from "vite";
import { reactRouter } from "@react-router/dev/vite";
import { llmsTxtPlugin } from "./app/discovery/llms-txt";

// @cloudflare/vite-plugin is deliberately not used here — see
// templates/site/profile/vite.config.ts for why: it claims the same Vite
// "ssr" environment React Router's own prerender step uses, breaking the
// build once a Worker `main` entry (workers/entry.ts) exists.
export default defineConfig({
  plugins: [reactRouter(), llmsTxtPlugin()],
});
