// The content-index-backed query helper shared by every MCP tool and the
// negotiation layer (data-model.md "MCP Tool" `query-source`). Re-reads the
// same content-index.json accessors the prerendered site already uses
// (Principle II — one source of truth, no separate parsing path here).
export { type ContentEntry, listEntries, getEntryBySlug } from "../content/load";

import { listEntries as _listEntries, type ContentEntry as _ContentEntry } from "../content/load";

// searchEntries does a simple case-insensitive keyword match across an
// entry's title/summary/body — keyword search only, still pre-Vectorize
// (Principle V; constitution storage-escalation triggers not yet crossed).
export function searchEntries(query: string): _ContentEntry[] {
  const needle = query.trim().toLowerCase();
  if (!needle) return [];
  return _listEntries().filter((entry) => {
    const haystack = [
      typeof entry.title === "string" ? entry.title : "",
      typeof entry.summary === "string" ? entry.summary : "",
      entry.body,
    ]
      .join(" ")
      .toLowerCase();
    return haystack.includes(needle);
  });
}
