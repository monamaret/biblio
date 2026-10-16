import { describe, expect, it } from "vitest";
import { buildLlmsTxt } from "./llms-txt";
import contentIndex from "../../content-index.json";
import type { QuarterlyIndex } from "../content/load";

const index = contentIndex as unknown as QuarterlyIndex;

// llms.txt for a quarterly site: name and description, then the current
// issue, every issue, the current issue's articles, and the colophon —
// links only, never bodies, never drafts (contracts/site-surface.md
// "llms.txt").
describe("buildLlmsTxt (quarterly)", () => {
  const txt = buildLlmsTxt(index);

  it("starts with the site name and description", () => {
    expect(txt.startsWith("# The Quarterly\n\n")).toBe(true);
    expect(txt).toContain(index.site!.description);
  });

  it("lists the current issue, every issue, current articles and the colophon", () => {
    expect(txt).toContain("## Current issue\n\n- No. 2 · Autumn 2026: /issues/02-autumn-2026\n");
    expect(txt).toContain("## Issues\n\n- No. 2 · Autumn 2026: /issues/02-autumn-2026\n- No. 1 · Spring 2026: /issues/01-spring-2026\n");
    expect(txt).toContain("## Articles (current issue)");
    expect(txt).toContain("- On Margins — A. Novak: /issues/02-autumn-2026/on-margins");
    expect(txt).not.toContain("On Reading Slowly");
    expect(txt).toContain("## About\n\n- Colophon: /colophon\n");
  });

  it("never inlines article bodies", () => {
    const article = index.entries.find((e) => e.kind === "article")!;
    expect(txt).not.toContain(article.body.slice(0, 40));
  });

  it("never lists draft issues", () => {
    const withDraft: QuarterlyIndex = {
      ...index,
      issues: [{ ...index.issues[0], issue: "03-winter-2026", slug: "issues/03-winter-2026", number: 3, draft: true }, ...index.issues],
    };
    expect(buildLlmsTxt(withDraft)).not.toContain("03-winter-2026");
  });

  it("handles a site with no issues", () => {
    const empty: QuarterlyIndex = { ...index, entries: [], issues: [], categories: [], current_issue: null };
    expect(() => buildLlmsTxt(empty)).not.toThrow();
    expect(buildLlmsTxt(empty)).not.toContain("## Current issue");
  });
});
