import { describe, expect, it } from "vitest";
import { filterArticles } from "./filters";
import { listArticles } from "./load";

function slugs(query: Parameters<typeof filterArticles>[1]): string[] {
  const result = filterArticles(listArticles(), query);
  if ("error" in result) throw new Error(result.error);
  return result.articles.map((a) => a.slug.split("/").pop()!);
}

// list_articles' filters: issue (slug or number), category (key or label,
// any case), author (case-insensitive exact), tag (exact), and an
// inclusive issue-date range — combined with AND (specs/004-quarterly-site-type
// FR-032; contracts/mcp-tools.md "list_articles").
describe("filterArticles", () => {
  it("returns every article, newest issue first, with no filters", () => {
    expect(slugs({})).toEqual(["on-margins", "letters-from-lodz", "type-in-brno", "on-reading-slowly", "the-rigid-grid"]);
  });

  it("filters by issue slug or number", () => {
    expect(slugs({ issue: "1" })).toEqual(["on-reading-slowly", "the-rigid-grid"]);
    expect(slugs({ issue: "01-spring-2026" })).toEqual(["on-reading-slowly", "the-rigid-grid"]);
    expect(slugs({ issue: "9" })).toEqual([]);
  });

  it("filters by category key or label in any case", () => {
    expect(slugs({ category: "Essays" })).toEqual(["on-margins", "letters-from-lodz", "on-reading-slowly"]);
    expect(slugs({ category: "essays" })).toEqual(slugs({ category: " ESSAYS " }));
    expect(slugs({ category: "poetry" })).toEqual([]);
  });

  it("filters by author, case-insensitively", () => {
    expect(slugs({ author: "a. novak" })).toEqual(["on-margins", "on-reading-slowly"]);
  });

  it("filters by tag", () => {
    expect(slugs({ tag: "grid" })).toEqual(["the-rigid-grid"]);
  });

  it("filters by an inclusive issue-date range", () => {
    expect(slugs({ dateFrom: "2026-06-01" })).toEqual(["on-margins", "letters-from-lodz", "type-in-brno"]);
    expect(slugs({ dateTo: "2026-04-01" })).toEqual(["on-reading-slowly", "the-rigid-grid"]);
    expect(slugs({ dateFrom: "2026-10-01", dateTo: "2026-10-01" })).toHaveLength(3);
  });

  it("combines filters with AND", () => {
    expect(slugs({ category: "essays", author: "A. Novak", dateFrom: "2026-06-01" })).toEqual(["on-margins"]);
  });

  it("rejects a malformed date, naming the field", () => {
    const result = filterArticles(listArticles(), { dateFrom: "June" });
    expect("error" in result && result.error).toContain("dateFrom");
  });
});
