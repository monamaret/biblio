import { describe, expect, it } from "vitest";
import { buildLlmsTxt } from "./llms-txt";
import type { ContentEntry } from "../content/load";

const entries: ContentEntry[] = [
  { slug: "about", path: "about.md", body: "...", title: "About Me" },
  { slug: "projects/sample-project", path: "projects/sample-project.md", body: "...", title: "Sample Project" },
];

// Minimal: site_name/description plus a flat list of key page links,
// derived directly from content-index.json — not a fuller dump
// (research.md "llms.txt investment level"; FR-005).
describe("buildLlmsTxt", () => {
  it("includes the site name and description", () => {
    const txt = buildLlmsTxt("My Profile Site", "A short summary", entries);
    expect(txt).toContain("My Profile Site");
    expect(txt).toContain("A short summary");
  });

  it("lists every content-index entry as a flat link, not a full content dump", () => {
    const txt = buildLlmsTxt("My Profile Site", "A short summary", entries);
    expect(txt).toContain("[About Me](/about)");
    expect(txt).toContain("[Sample Project](/projects/sample-project)");
    // Not a fuller dump: entry bodies must not be inlined.
    expect(txt).not.toContain(entries[0].body);
  });

  it("handles zero entries without throwing", () => {
    expect(() => buildLlmsTxt("Empty Site", "", [])).not.toThrow();
  });
});
