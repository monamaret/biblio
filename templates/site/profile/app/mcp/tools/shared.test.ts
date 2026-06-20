import { describe, expect, it } from "vitest";
import { listPages, getPage, searchContent } from "./shared";

// Every shared tool's primary content response must be a clean,
// self-sufficient Markdown/plain-text rendering (Principle IV, FR-003) —
// this is what every tool returns regardless of richer-view support.
describe("shared MCP tools content contract", () => {
  it("list_pages returns a Markdown list of every authored entry", () => {
    const result = listPages();
    expect(result.content[0].type).toBe("text");
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("About Me");
    expect(text).toContain("/about");
  });

  it("get_page returns the full Markdown body for a known slug", () => {
    const result = getPage("about");
    const text = (result.content[0] as { text: string }).text;
    expect(text).toContain("About Me");
    expect(text).toContain("Welcome to your new profile site");
  });

  it("get_page rejects an unknown slug with a clear, structured error rather than a silent empty result", () => {
    expect(() => getPage("does-not-exist")).toThrow(/no page found/);
  });

  it("search_content returns a self-sufficient text answer for matches", () => {
    const result = searchContent("profile");
    const textBlock = result.content.find((c) => c.type === "text");
    expect(textBlock).toBeDefined();
    expect((textBlock as { text: string }).text.length).toBeGreaterThan(0);
  });

  it("search_content returns a self-sufficient text answer even with zero matches", () => {
    const result = searchContent("xyznonexistentquery");
    const textBlock = result.content.find((c) => c.type === "text");
    expect(textBlock).toBeDefined();
    expect((textBlock as { text: string }).text).toContain("No pages matched");
  });
});

// A tool result carrying a `ui-resource` must still produce the unchanged,
// self-sufficient `content` response alongside it — never a substitute for
// it (Principle IV; data-model.md "MCP Tool" `ui-resource`; spec Scenario 2).
describe("ui-resource pairing", () => {
  it("search_content's resource block never replaces its text content", () => {
    const result = searchContent("profile");
    const textBlock = result.content.find((c) => c.type === "text");
    const resourceBlock = result.content.find((c) => c.type === "resource");

    expect(textBlock).toBeDefined();
    expect(resourceBlock).toBeDefined();
    expect((textBlock as { text: string }).text.length).toBeGreaterThan(0);
  });

  it("search_content with zero matches has no resource block to attach (nothing to chart)", () => {
    const result = searchContent("xyznonexistentquery");
    const resourceBlock = result.content.find((c) => c.type === "resource");
    expect(resourceBlock).toBeUndefined();
  });
});
