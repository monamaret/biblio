import { getEntryBySlug, listEntries } from "../mcp/query";

// renderMarkdown is tier-i: raw Markdown passthrough for a content-index
// entry by slug — research.md "In-Worker Markdown→ANSI rendering" names
// this the baseline tier every Markdown/terminal-preferring client gets
// when it isn't recognized as ANSI-capable. The empty slug (the home/TOC
// page, "/") has no single backing content-index entry, so it gets a
// generated index listing instead of a false 404 — it's still a real,
// existing page (FR-004 covers every page, not just per-entry detail
// pages).
export function renderMarkdown(slug: string): string | undefined {
  if (slug === "") {
    const lines = listEntries().map((entry) => {
      const title = typeof entry.title === "string" ? entry.title : entry.slug;
      return `- [${title}](/${entry.slug})`;
    });
    return `# Welcome to your biblio docs site\n\n${lines.join("\n")}`;
  }

  const entry = getEntryBySlug(slug);
  if (!entry) return undefined;
  const title = typeof entry.title === "string" ? entry.title : entry.slug;
  return `# ${title}\n\n${entry.body}`;
}
