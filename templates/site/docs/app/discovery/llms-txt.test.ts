import { describe, expect, it } from "vitest";
import { buildLlmsTxt } from "./llms-txt";
import type { ContentEntry } from "../content/load";

const entries: ContentEntry[] = [
  { slug: "index", path: "index.md", body: "...", title: "Introduction", order: 1, section: "Overview" },
  { slug: "getting-started", path: "getting-started.md", body: "...", title: "Getting Started", order: 2, section: "Overview" },
];

// Minimal: site_name/description plus a flat list of key page links,
// derived directly from content-index.json — not a fuller dump
// (research.md "llms.txt investment level"; FR-005). Section/order-aware
// in the sense that it lists every entry regardless of section, exactly as
// every other site type's llms.txt does (no docs-specific grouping needed
// here, since the file is already a flat key-page list by design).
describe("buildLlmsTxt", () => {
  it("includes the site name and description", () => {
    const txt = buildLlmsTxt("My Docs Site", "A short summary", entries);
    expect(txt).toContain("My Docs Site");
    expect(txt).toContain("A short summary");
  });

  it("lists every content-index entry as a flat link, not a full content dump", () => {
    const txt = buildLlmsTxt("My Docs Site", "A short summary", entries);
    expect(txt).toContain("[Introduction](/index)");
    expect(txt).toContain("[Getting Started](/getting-started)");
    // Not a fuller dump: entry bodies must not be inlined.
    expect(txt).not.toContain(entries[0].body);
  });

  it("handles zero entries without throwing", () => {
    expect(() => buildLlmsTxt("Empty Site", "", [])).not.toThrow();
  });
});
