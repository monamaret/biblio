import { describe, expect, it } from "vitest";
import { renderMarkdown } from "./markdown";

// Markdown views of the list pages are built from the index; an article's
// view is its title plus body (contracts/site-surface.md "Content
// negotiation"). Numerals follow typesetting.numerals (roman by default).
describe("renderMarkdown", () => {
  it("renders the current issue's contents for the home page", () => {
    const md = renderMarkdown("")!;
    expect(md).toContain("# The Quarterly");
    expect(md).toContain("No. 2 · Autumn 2026 — Notes from the Margin");
    expect(md).toContain("## Essays");
    expect(md).toContain("I. [On Margins](/issues/02-autumn-2026/on-margins) — A. Novak");
    expect(md).toContain("## Design");
    expect(md).toContain("III. [Type in Brno](/issues/02-autumn-2026/type-in-brno) — M. Hora");
    expect(md).toContain("[Previous issues](/issues)");
  });

  it("renders a specific issue's contents", () => {
    const md = renderMarkdown("issues/01-spring-2026")!;
    expect(md).toContain("No. 1 · Spring 2026");
    expect(md).toContain("[On Reading Slowly](/issues/01-spring-2026/on-reading-slowly)");
  });

  it("renders the archive", () => {
    const md = renderMarkdown("issues")!;
    expect(md).toContain("# Issues");
    expect(md.indexOf("No. 2")).toBeLessThan(md.indexOf("No. 1"));
    expect(md).toContain("[No. 1 · Spring 2026 — Margins and Grids](/issues/01-spring-2026) — 2 articles");
  });

  it("renders a category page", () => {
    const md = renderMarkdown("category/design")!;
    expect(md).toContain("# Design");
    expect(md).toContain("[Type in Brno](/issues/02-autumn-2026/type-in-brno) — M. Hora · No. 2 · Autumn 2026");
    expect(md).not.toContain("On Margins");
  });

  it("renders an article as its title plus body", () => {
    const md = renderMarkdown("issues/02-autumn-2026/type-in-brno")!;
    expect(md.startsWith("# Type in Brno\n\n")).toBe(true);
  });

  it("renders the colophon", () => {
    expect(renderMarkdown("colophon")).toContain("# Colophon");
  });

  it("returns undefined for unknown pages", () => {
    expect(renderMarkdown("issues/99-nope")).toBeUndefined();
    expect(renderMarkdown("category/poetry")).toBeUndefined();
    expect(renderMarkdown("nope")).toBeUndefined();
  });
});
