import { fileURLToPath } from "node:url";
import fs from "node:fs";
import path from "node:path";
import type { Config } from "@react-router/dev/config";

type ContentEntry = { slug: string };

// Reads the unchanged content-index.json produced by `biblio-cli index` —
// the single source of truth this prerender step consumes read-only
// (Principle II; data-model.md "Pre-rendered Page Set" `source-index`). An
// absent index (before the author's first `biblio-cli index` run) yields no
// entries rather than failing the dev server.
function contentEntries(): ContentEntry[] {
  const dir = path.dirname(fileURLToPath(import.meta.url));
  const indexPath = path.join(dir, "content-index.json");
  if (!fs.existsSync(indexPath)) {
    return [];
  }
  const parsed = JSON.parse(fs.readFileSync(indexPath, "utf-8")) as {
    entries?: ContentEntry[];
  };
  return parsed.entries ?? [];
}

export default {
  // Matches the `./dist/client` assets directory `biblio-cli new` already
  // pre-fills in wrangler.jsonc (internal/scaffold/generator.go); React
  // Router's own default is `build/client`.
  buildDirectory: "dist",
  // Workers Static Assets serves pre-rendered HTML directly — no SSR worker
  // route exists yet (that's v0.03's `/mcp` mount; Principle V topology).
  ssr: false,
  async prerender() {
    // One pre-rendered route per content-index.json entry, plus the list/nav
    // home route — the Pre-rendered Page Set completeness invariant
    // (FR-001; data-model.md "Pre-rendered Page Set").
    return ["/", ...contentEntries().map((entry) => `/${entry.slug}`)];
  },
  // Required for the static prerender build to actually run a (transient)
  // SSR pass — without it, `vite build` silently produces only the client
  // bundle and no prerendered HTML, with neither an error nor a warning.
  future: { v8_viteEnvironmentApi: true },
} satisfies Config;
