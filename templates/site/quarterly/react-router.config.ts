import { fileURLToPath } from "node:url";
import fs from "node:fs";
import path from "node:path";
import type { Config } from "@react-router/dev/config";
import type { QuarterlyIndex } from "./app/content/load";

// prerenderPaths lists every published page of a quarterly site from its
// content index: the current-issue home, the archive, each issue's
// contents, each category, each article, and the colophon. Drafts are
// skipped (a build-mode index never has them; this is belt and braces).
// Issue entries themselves are served by the contents route, not the splat
// (specs/004-quarterly-site-type FR-013, SC-002; research.md R12).
export function prerenderPaths(index: QuarterlyIndex): string[] {
  const published = index.issues.filter((issue) => !issue.draft);
  const publishedIssues = new Set(published.map((issue) => issue.issue));
  const paths = ["/", "/issues"];
  for (const issue of published) {
    paths.push(`/issues/${issue.issue}`);
  }
  for (const category of index.categories) {
    paths.push(`/category/${category.key}`);
  }
  for (const entry of index.entries) {
    if (entry.kind === "article" && !entry.draft && publishedIssues.has(entry.issue as string)) {
      paths.push(`/${entry.slug}`);
    }
  }
  if (index.entries.some((entry) => entry.kind === "colophon")) {
    paths.push("/colophon");
  }
  return paths;
}

// readIndex reads the content-index.json biblio-cli wrote — the single
// source of truth this prerender step consumes read-only (Principle II).
// An absent index (before the first `biblio-cli build`/`dev`) yields just
// the home and archive pages rather than failing.
function readIndex(): QuarterlyIndex {
  const dir = path.dirname(fileURLToPath(import.meta.url));
  const indexPath = path.join(dir, "content-index.json");
  if (!fs.existsSync(indexPath)) {
    return { entries: [], issues: [], categories: [], current_issue: null, typesetting: { margin_notes: true, numerals: "roman" } };
  }
  return JSON.parse(fs.readFileSync(indexPath, "utf-8")) as QuarterlyIndex;
}

export default {
  // Matches the `./dist/client` assets directory `biblio-cli new` already
  // pre-fills in wrangler.jsonc (internal/scaffold/generator.go).
  buildDirectory: "dist",
  // Workers Static Assets serves pre-rendered HTML directly — no SSR worker
  // route exists beyond the `/mcp` mount (Principle V topology).
  ssr: false,
  async prerender() {
    return prerenderPaths(readIndex());
  },
  // Required for the static prerender build to actually run a (transient)
  // SSR pass — without it, `vite build` silently produces only the client
  // bundle and no prerendered HTML.
  future: { v8_viteEnvironmentApi: true },
} satisfies Config;
