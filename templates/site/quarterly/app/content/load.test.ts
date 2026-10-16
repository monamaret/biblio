import { describe, expect, it } from "vitest";
import {
  articlesInCategory,
  articlesInIssue,
  currentIssue,
  getColophon,
  getEntryBySlug,
  getIssue,
  listCategories,
  listIssues,
  siteMeta,
  typesetting,
} from "./load";

// The accessors read content-index.json's quarterly arrays directly — no
// path parsing on this side (Principle II; specs/004-quarterly-site-type
// data-model.md "content-index.json shape"). These run against the
// template's own sample content (two issues, two categories, a colophon).
describe("quarterly index accessors", () => {
  it("lists published issues newest first", () => {
    const issues = listIssues();
    expect(issues.map((i) => i.issue)).toEqual(["02-autumn-2026", "01-spring-2026"]);
    expect(issues[0]).toMatchObject({ number: 2, season: "Autumn 2026", article_count: 3 });
  });

  it("finds an issue by folder slug or by number", () => {
    expect(getIssue("01-spring-2026")?.number).toBe(1);
    expect(getIssue("1")?.issue).toBe("01-spring-2026");
    expect(getIssue("99")).toBeUndefined();
    expect(getIssue("no-such-issue")).toBeUndefined();
  });

  it("returns the current issue", () => {
    expect(currentIssue()?.issue).toBe("02-autumn-2026");
  });

  it("lists categories with their first-seen labels", () => {
    expect(listCategories().map((c) => [c.key, c.label])).toEqual([
      ["essays", "Essays"],
      ["design", "Design"],
    ]);
  });

  it("groups an issue's articles by category in contents order", () => {
    const groups = articlesInIssue("02-autumn-2026");
    expect(groups.map((g) => g.category.key)).toEqual(["essays", "design"]);
    expect(groups[0].articles.map((a) => a.slug)).toEqual([
      "issues/02-autumn-2026/on-margins",
      "issues/02-autumn-2026/letters-from-lodz",
    ]);
    expect(groups[1].articles.map((a) => a.position)).toEqual([3]);
  });

  it("lists a category's articles across issues, newest issue first", () => {
    expect(articlesInCategory("essays").map((a) => a.slug)).toEqual([
      "issues/02-autumn-2026/on-margins",
      "issues/02-autumn-2026/letters-from-lodz",
      "issues/01-spring-2026/on-reading-slowly",
    ]);
    expect(articlesInCategory("nope")).toEqual([]);
  });

  it("finds the colophon and article entries by slug", () => {
    expect(getColophon()?.title).toBe("Colophon");
    expect(getEntryBySlug("issues/02-autumn-2026/on-margins")?.kind).toBe("article");
  });

  it("exposes the site name and the structural typesetting switches", () => {
    expect(siteMeta().name).toBe("The Quarterly");
    expect(typesetting()).toEqual({ margin_notes: true, numerals: "roman" });
  });
});
