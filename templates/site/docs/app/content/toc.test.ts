import { describe, expect, it } from "vitest";
import { toc } from "./toc";
import type { ContentEntry } from "./load";

function page(overrides: Partial<ContentEntry>): ContentEntry {
  return { slug: "p", path: "p.md", body: "", title: "P", section: "A", order: 1, ...overrides };
}

describe("toc", () => {
  it("groups pages into sections, ordered by order within each section", () => {
    const entries: ContentEntry[] = [
      page({ slug: "intro", title: "Intro", section: "Overview", order: 1 }),
      page({ slug: "getting-started", title: "Getting Started", section: "Overview", order: 2 }),
      page({ slug: "config", title: "Config", section: "Reference", order: 1 }),
    ];
    const sections = toc(entries);
    expect(sections).toHaveLength(2);
    const overview = sections.find((s) => s.section === "Overview");
    expect(overview?.pages.map((p) => p.slug)).toEqual(["intro", "getting-started"]);
    const reference = sections.find((s) => s.section === "Reference");
    expect(reference?.pages.map((p) => p.slug)).toEqual(["config"]);
  });

  it("orders sections deterministically by their lowest-order page", () => {
    const entries: ContentEntry[] = [
      page({ slug: "config", section: "Reference", order: 5 }),
      page({ slug: "intro", section: "Overview", order: 1 }),
    ];
    const sections = toc(entries);
    expect(sections.map((s) => s.section)).toEqual(["Overview", "Reference"]);
  });

  it("resolves two pages with the same section/order via the documented default (slug tiebreak), not an undefined order", () => {
    const entries: ContentEntry[] = [
      page({ slug: "z-page", section: "Overview", order: 1 }),
      page({ slug: "a-page", section: "Overview", order: 1 }),
    ];
    const sections = toc(entries);
    expect(sections[0].pages.map((p) => p.slug)).toEqual(["a-page", "z-page"]);
  });

  it("every page is reachable through the returned structure (SC-001)", () => {
    const entries: ContentEntry[] = [
      page({ slug: "one", section: "A", order: 1 }),
      page({ slug: "two", section: "B", order: 1 }),
      page({ slug: "three", section: "A", order: 2 }),
    ];
    const sections = toc(entries);
    const allSlugs = sections.flatMap((s) => s.pages.map((p) => p.slug));
    expect(allSlugs.sort()).toEqual(["one", "three", "two"]);
  });
});
