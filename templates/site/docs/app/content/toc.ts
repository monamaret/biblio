import type { ContentEntry } from "./load";

// TocSection is one section of the docs site's hierarchy-based browsing
// structure — every page in `pages`, already ordered by its `order` field
// (data-model.md "Table of Contents / Navigation Structure").
export type TocSection = {
  section: string;
  pages: ContentEntry[];
};

function orderOf(entry: ContentEntry): number {
  const order = entry.order;
  return typeof order === "number" ? order : Number.MAX_SAFE_INTEGER;
}

// toc deterministically groups entries into sections[] from their
// section/order fields (FR-002, Scenario 1; data-model.md "Table of
// Contents / Navigation Structure"). It is the canonical TOC-computation
// module shared by the browser UI (routes/home.tsx) and the MCP tools
// (mcp/tools/docs.ts) — a single owner, mirroring
// templates/site/log/app/content/filters.ts's pattern (Principle II).
//
// Resolved default for ambiguous placement (spec Edge Cases): this schema
// only ever carries one `section` string per page (internal/content/schema.go
// validateDocsFields), so a page never literally belongs to two sections —
// each page's `section` value is its one primary-section assignment, full
// stop. Within that, two pages sharing the same section *and* the same
// `order` value resolve their relative order deterministically by slug
// (ascending) rather than leaving it undefined or duplicating either page's
// placement. Sections themselves are ordered by their lowest-order page, so
// the section list itself never depends on the entries array's incidental
// (filesystem-walk) order.
export function toc(entries: ContentEntry[]): TocSection[] {
  const bySection = new Map<string, ContentEntry[]>();
  for (const entry of entries) {
    const section = typeof entry.section === "string" ? entry.section : "";
    const pages = bySection.get(section) ?? [];
    pages.push(entry);
    bySection.set(section, pages);
  }

  const sections: TocSection[] = Array.from(bySection.entries()).map(([section, pages]) => ({
    section,
    pages: [...pages].sort((a, b) => orderOf(a) - orderOf(b) || a.slug.localeCompare(b.slug)),
  }));

  sections.sort((a, b) => orderOf(a.pages[0]) - orderOf(b.pages[0]));

  return sections;
}
