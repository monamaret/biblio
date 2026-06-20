import { describe, expect, it } from "vitest";
import { listSections, getDocPage } from "./docs";

describe("listSections", () => {
  it("lists every section with its pages, content-complete on its own", () => {
    const result = listSections();
    const text = result.content[0];
    if (text.type !== "text") throw new Error("expected text content");
    expect(text.text).toContain("Overview");
    expect(text.text).toContain("Reference");
    expect(text.text).toContain("Introduction");
    expect(text.text).toContain("Configuration");
  });

  it("rejects a log-shaped tag-filter argument with a clear, structured error rather than silently ignoring it (spec Edge Cases)", () => {
    expect(() => listSections({ tag: "shipping" } as unknown as undefined)).toThrow(/tag/i);
  });
});

describe("getDocPage", () => {
  it("returns the full Markdown body of a page by its slug", () => {
    const result = getDocPage("getting-started");
    const text = result.content[0];
    if (text.type !== "text") throw new Error("expected text content");
    expect(text.text).toContain("Getting Started");
    expect(text.text).toContain("Walk readers through the first steps");
  });

  it("rejects a slug that doesn't match any page with a clear, structured error", () => {
    expect(() => getDocPage("does-not-exist")).toThrow(/no page found/i);
  });
});
