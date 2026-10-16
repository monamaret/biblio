import { describe, expect, it } from "vitest";
import { getArticle, getIssue, listArticles, listIssues } from "./quarterly";
import { getEntryBySlug } from "../../content/load";

function text(result: { content: Array<{ type: string; text?: string }> }): string {
  return result.content.map((c) => c.text ?? "").join("\n");
}

// The quarterly tools' text answers and structured output, per
// contracts/mcp-tools.md (specs/004-quarterly-site-type FR-032–FR-034).
describe("list_issues", () => {
  it("returns a table of published issues, newest first", () => {
    const out = text(listIssues());
    expect(out).toContain("| No. | Title | Date | Articles | Slug |");
    expect(out).toContain("| 2 | Notes from the Margin | 2026-10-01 | 3 | 02-autumn-2026 |");
    expect(out.indexOf("| 2 |")).toBeLessThan(out.indexOf("| 1 |"));
  });
});

describe("get_issue", () => {
  it("returns an issue's contents grouped by category, by slug or by number", () => {
    const bySlug = text(getIssue("02-autumn-2026"));
    expect(bySlug).toContain("# No. 2 · Autumn 2026 — Notes from the Margin");
    expect(bySlug).toContain("## Essays");
    expect(bySlug).toContain("I. On Margins — A. Novak (/issues/02-autumn-2026/on-margins)");
    expect(bySlug).toContain("III. Type in Brno — M. Hora (/issues/02-autumn-2026/type-in-brno)");
    expect(bySlug.indexOf("## Essays")).toBeLessThan(bySlug.indexOf("## Design"));
    expect(text(getIssue("2"))).toBe(bySlug);
  });

  it("returns a not-found text answer for an unknown issue", () => {
    const result = getIssue("7");
    expect(result.isError).toBeFalsy();
    expect(text(result)).toBe('No issue matches "7". Use list_issues to see published issues.');
  });
});

describe("list_articles", () => {
  it("lists article summaries for the given filters", () => {
    const out = text(listArticles({ category: "design" }));
    expect(out).toBe(
      [
        "- Type in Brno — M. Hora · Design · No. 2 — A short walk through the city's printing houses. (/issues/02-autumn-2026/type-in-brno)",
        "- The Rigid Grid — J. Kral · Design · No. 1 — Why a strict grid frees the designer instead of binding them. (/issues/01-spring-2026/the-rigid-grid)",
      ].join("\n"),
    );
  });

  it("says so when nothing matches", () => {
    expect(text(listArticles({ tag: "nope" }))).toBe("No articles matched the given filter.");
  });

  it("names the field when a date is malformed", () => {
    const result = listArticles({ dateTo: "soon" });
    expect(result.isError).toBeFalsy();
    expect(text(result)).toContain("dateTo");
  });
});

describe("get_article", () => {
  it("returns the full text with every note, plus the notes as structured data", () => {
    const result = getArticle("/issues/02-autumn-2026/on-margins");
    const out = text(result);
    expect(out.startsWith("# On Margins\n")).toBe(true);
    expect(out).toContain("A. Novak · Essays · No. 2 · Autumn 2026");
    expect(out).toContain("[^m1]: cf. Gerstner, *Designing Programmes*, 1964");
    expect(out).toContain("[^1]: Brno, 1971.");
    expect(out).not.toContain("{>>");

    const s = result.structuredContent as Record<string, unknown>;
    expect(s).toMatchObject({
      slug: "issues/02-autumn-2026/on-margins",
      title: "On Margins",
      author: "A. Novak",
      category: "Essays",
      issue: { slug: "02-autumn-2026", number: 2, date: "2026-10-01" },
    });
    const footnotes = s.footnotes as Array<Record<string, unknown>>;
    const marginNotes = s.marginNotes as Array<Record<string, unknown>>;
    expect(footnotes).toHaveLength(3);
    expect(marginNotes).toHaveLength(3);
    expect(footnotes[0]).toEqual({ id: "fn-1", label: "1", number: 1, text: "Brno, 1971.", refOffset: expect.any(Number) });
    expect(marginNotes[0]).toEqual({ id: "m1", text: "cf. Gerstner, *Designing Programmes*, 1964", refOffset: expect.any(Number) });
    expect(JSON.stringify(s)).not.toMatch(/ref_offset|margin_notes/);
  });

  it("points refOffset at the reference marker, in code points", () => {
    const s = getArticle("issues/02-autumn-2026/on-margins").structuredContent as {
      footnotes: Array<{ refOffset: number }>;
      marginNotes: Array<{ refOffset: number }>;
    };
    // Offsets index the article's body as stored in the index.
    const article = Array.from(getEntryBySlug("issues/02-autumn-2026/on-margins")!.body);
    expect(article.slice(s.footnotes[0].refOffset, s.footnotes[0].refOffset + 4).join("")).toBe("[^1]");
    expect(article.slice(s.marginNotes[0].refOffset, s.marginNotes[0].refOffset + 3).join("")).toBe("{>>");
  });

  it("accepts a bare <issue>/<article> slug", () => {
    expect(getArticle("02-autumn-2026/type-in-brno").structuredContent).toMatchObject({ title: "Type in Brno" });
  });

  it("returns a not-found text answer for an unknown article", () => {
    const result = getArticle("issues/02-autumn-2026/nope");
    expect(result.isError).toBe(true);
    expect(result.structuredContent).toBeUndefined();
    expect(text(result)).toBe('No article matches "issues/02-autumn-2026/nope". Use list_articles to find one.');
  });
});
