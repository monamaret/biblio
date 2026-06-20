import { describe, expect, it } from "vitest";
import { toc } from "./toc";
import type { ContentEntry } from "./load";

// Asserts the deliberate absence of tag/project filtering for the docs
// site type (spec Scenario 1.2; contracts/docs-site-and-review.md "Filter
// UI"): the documentation type's hierarchy-based browsing model has no
// filter-query surface at all, unlike the project-log type's
// `filterEntries()`/FilterQuery. Proving this negative is itself a
// required, load-bearing test (Principle VI; data-model.md "Table of
// Contents / Navigation Structure").
describe("no filter surface on the docs site type", () => {
  it("toc() output carries no tags/project fields on any page", () => {
    const entries: ContentEntry[] = [
      { slug: "intro", path: "index.md", body: "", title: "Intro", section: "Overview", order: 1 },
    ];
    const sections = toc(entries);
    for (const section of sections) {
      for (const page of section.pages) {
        expect(page).not.toHaveProperty("tags");
        expect(page).not.toHaveProperty("project");
      }
    }
  });

  it("toc()'s return shape exposes no filter-query parameter (no FilterQuery-like argument)", () => {
    // toc() takes only entries — no second "query"/"filter" argument exists
    // for this site type, unlike templates/site/log/app/content/filters.ts's
    // filterEntries(entries, query).
    expect(toc.length).toBe(1);
  });
});
