import type { ContentEntry } from "./load";

// FilterQuery is a visitor- or agent-issued request to narrow the entry
// list by one or more dimensions (FR-003; data-model.md "Filter Query").
// All fields are optional; an absent dimension is not filtered on.
export type FilterQuery = {
  dateFrom?: string;
  dateTo?: string;
  tag?: string;
  project?: string;
};

// filterEntries computes the FilterQuery's intersection over entries — the
// single, canonical filter-computation module (Principle II) shared by the
// entry-list UI (home.tsx) and the log-specific MCP tools (mcp/tools/log.ts).
// It reads tag/project values directly off each entry rather than against
// any pre-declared list, so a brand-new tag/project value is immediately
// filterable as soon as it appears in content-index.json — no separate
// registration step (spec Edge Cases). Combining dimensions returns the
// intersection: every entry must match every given dimension, with no
// false positives or omissions (spec Edge Cases).
export function filterEntries(entries: ContentEntry[], query: FilterQuery): ContentEntry[] {
  return entries.filter((entry) => {
    if (query.dateFrom || query.dateTo) {
      const date = typeof entry.date === "string" ? entry.date : undefined;
      if (!date) return false;
      if (query.dateFrom && date < query.dateFrom) return false;
      if (query.dateTo && date > query.dateTo) return false;
    }

    if (query.tag) {
      const tags = Array.isArray(entry.tags) ? (entry.tags as unknown[]) : [];
      if (!tags.includes(query.tag)) return false;
    }

    if (query.project) {
      if (entry.project !== query.project) return false;
    }

    return true;
  });
}
