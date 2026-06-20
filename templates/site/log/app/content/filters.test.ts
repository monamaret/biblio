import { describe, expect, it } from "vitest";
import { filterEntries, type FilterQuery } from "./filters";
import type { ContentEntry } from "./load";

const entries: ContentEntry[] = [
  {
    slug: "posts/a",
    path: "posts/a.md",
    body: "",
    title: "A",
    date: "2026-06-01",
    tags: ["meta"],
  },
  {
    slug: "posts/b",
    path: "posts/b.md",
    body: "",
    title: "B",
    date: "2026-06-05",
    tags: ["shipping", "frontend"],
    project: "biblio",
  },
  {
    slug: "posts/c",
    path: "posts/c.md",
    body: "",
    title: "C",
    date: "2026-06-10",
    tags: ["shipping", "backend"],
    project: "biblio",
  },
];

function slugs(result: ContentEntry[]): string[] {
  return result.map((e) => e.slug);
}

// filterEntries computes date-range/tag/project filtering over
// content-index.json entries — the single, canonical filter-computation
// module shared by the entry-list UI and the log-specific MCP tools
// (FR-003; data-model.md "Filter Query"/"Filtered Entry Set").
describe("filterEntries — single-dimension filters", () => {
  it("filters by date range alone", () => {
    const query: FilterQuery = { dateFrom: "2026-06-05", dateTo: "2026-06-10" };
    expect(slugs(filterEntries(entries, query))).toEqual(["posts/b", "posts/c"]);
  });

  it("filters by tag alone", () => {
    const query: FilterQuery = { tag: "frontend" };
    expect(slugs(filterEntries(entries, query))).toEqual(["posts/b"]);
  });

  it("filters by project alone", () => {
    const query: FilterQuery = { project: "biblio" };
    expect(slugs(filterEntries(entries, query))).toEqual(["posts/b", "posts/c"]);
  });

  it("returns every entry when no filter dimension is given", () => {
    expect(slugs(filterEntries(entries, {}))).toEqual(["posts/a", "posts/b", "posts/c"]);
  });
});

// Combining dimensions MUST return the intersection — no false positives,
// no omissions (spec Edge Cases).
describe("filterEntries — combined intersection semantics", () => {
  it("intersects a date range and a tag", () => {
    const query: FilterQuery = { dateFrom: "2026-06-01", dateTo: "2026-06-05", tag: "shipping" };
    // posts/b matches both; posts/a has no "shipping" tag; posts/c is out of range.
    expect(slugs(filterEntries(entries, query))).toEqual(["posts/b"]);
  });

  it("returns an empty result when dimensions don't overlap on any entry", () => {
    const query: FilterQuery = { tag: "meta", project: "biblio" };
    expect(filterEntries(entries, query)).toEqual([]);
  });
});

// A brand-new tag/project value present in content-index.json must be
// immediately filterable with no separate registration step — filtering
// reads directly off whatever values exist on the entries, never a
// pre-declared tag/project list (spec Edge Cases).
describe("filterEntries — new tag/project requires no registration step", () => {
  it("matches a tag never seen before in any prior filter call", () => {
    const freshEntries: ContentEntry[] = [
      ...entries,
      { slug: "posts/d", path: "posts/d.md", body: "", title: "D", date: "2026-06-15", tags: ["brand-new-tag"] },
    ];
    expect(slugs(filterEntries(freshEntries, { tag: "brand-new-tag" }))).toEqual(["posts/d"]);
  });

  it("matches a project never seen before in any prior filter call", () => {
    const freshEntries: ContentEntry[] = [
      ...entries,
      { slug: "posts/e", path: "posts/e.md", body: "", title: "E", date: "2026-06-20", project: "brand-new-project" },
    ];
    expect(slugs(filterEntries(freshEntries, { project: "brand-new-project" }))).toEqual(["posts/e"]);
  });
});
