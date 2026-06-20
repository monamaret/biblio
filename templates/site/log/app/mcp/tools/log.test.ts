import { describe, expect, it } from "vitest";
import { listEntries, getEntry } from "./log";

// Every log-specific tool's content response must also be a clean,
// self-sufficient Markdown/plain-text rendering (FR-005; Principle IV).
describe("log-specific MCP tools content contract", () => {
  it("list_entries returns every entry's title, date, and slug when no filter is given", () => {
    const result = listEntries({});
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("Hello, World");
    expect(text).toContain("2026-06-01");
    expect(text).toContain("/posts/2026-06-01-hello-world");
    expect(text).toContain("Shipped the filter UI");
    expect(text).toContain("Refactored the content index");
  });

  it("list_entries narrows by tag", () => {
    const result = listEntries({ tag: "frontend" });
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("Shipped the filter UI");
    expect(text).not.toContain("Hello, World");
  });

  it("list_entries narrows by project", () => {
    const result = listEntries({ project: "biblio" });
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("Shipped the filter UI");
    expect(text).toContain("Refactored the content index");
    expect(text).not.toContain("Hello, World");
  });

  it("list_entries narrows by date range", () => {
    const result = listEntries({ dateFrom: "2026-06-05", dateTo: "2026-06-10" });
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("Shipped the filter UI");
    expect(text).toContain("Refactored the content index");
    expect(text).not.toContain("Hello, World");
  });

  it("get_entry returns the full Markdown body for a known slug", () => {
    const result = getEntry("posts/2026-06-01-hello-world");
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("Hello, World");
    expect(text).toContain("first entry in your project log");
  });

  it("get_entry rejects a slug that doesn't exist with a clear error, not a silent empty result", () => {
    expect(() => getEntry("does-not-exist")).toThrow();
  });
});
