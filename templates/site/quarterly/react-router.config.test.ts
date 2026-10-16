import { describe, expect, it } from "vitest";
import { prerenderPaths } from "./react-router.config";
import contentIndex from "./content-index.json";
import type { QuarterlyIndex } from "./app/content/load";

const index = contentIndex as unknown as QuarterlyIndex;

// Every published page is prerendered — home, archive, each issue, each
// category, each article, the colophon — and nothing for drafts or
// issue.md entries themselves (specs/004 FR-013, SC-002; research.md R12).
describe("prerenderPaths", () => {
  it("lists every published page exactly once", () => {
    const paths = prerenderPaths(index);
    expect(paths).toEqual([
      "/",
      "/issues",
      "/issues/02-autumn-2026",
      "/issues/01-spring-2026",
      "/category/essays",
      "/category/design",
      "/issues/02-autumn-2026/on-margins",
      "/issues/02-autumn-2026/letters-from-lodz",
      "/issues/02-autumn-2026/type-in-brno",
      "/issues/01-spring-2026/on-reading-slowly",
      "/issues/01-spring-2026/the-rigid-grid",
      "/colophon",
    ]);
  });

  it("never prerenders author or comment pages", () => {
    expect(prerenderPaths(index).some((p) => /\/(author|comments?)\b/.test(p))).toBe(false);
  });

  it("skips draft issues and their articles", () => {
    const withDraft: QuarterlyIndex = {
      ...index,
      issues: [{ ...index.issues[0], issue: "03-winter-2026", slug: "issues/03-winter-2026", number: 3, draft: true }, ...index.issues],
      entries: [
        ...index.entries,
        { slug: "issues/03-winter-2026/wip", path: "issues/03-winter-2026/wip.md", body: "", kind: "article", draft: true, title: "WIP" },
      ],
    };
    const paths = prerenderPaths(withDraft);
    expect(paths.some((p) => p.includes("03-winter-2026"))).toBe(false);
  });
});
