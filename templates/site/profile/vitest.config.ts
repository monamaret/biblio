import { defineConfig } from "vitest/config";
import react from "@vitejs/plugin-react";

// A separate config from vite.config.ts: the React Router vite plugin
// expects a framework-mode dev/build pipeline, which fights vitest's own
// transform pipeline — tests just need React's plugin for JSX.
export default defineConfig({
  plugins: [react()],
  test: {
    environment: "jsdom",
    setupFiles: ["./vitest.setup.ts"],
  },
});
