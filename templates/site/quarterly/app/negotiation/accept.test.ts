import { describe, expect, it } from "vitest";
import { negotiateVariant, negotiatePageResponse } from "./accept";

// The type-agnostic negotiation logic (copied unchanged from the log
// template) against the quarterly site's content shape (specs/004 FR-004;
// contracts/site-surface.md "Content negotiation").
describe("negotiatePageResponse — quarterly content shape", () => {
  it("returns an article's Markdown with the Vary header", async () => {
    const response = await negotiatePageResponse("issues/02-autumn-2026/type-in-brno", "text/markdown", "some-client/1.0");
    expect(response).toBeDefined();
    expect(response!.headers.get("Vary")).toBe("Accept, User-Agent");
    expect(await response!.text()).toContain("# Type in Brno");
  });

  it("selects the ansi tier for a terminal User-Agent", () => {
    expect(negotiateVariant("text/markdown", "curl/8.0")).toBe("ansi");
  });

  it("returns a structured 404 for an unknown slug", async () => {
    const response = await negotiatePageResponse("does-not-exist", "text/markdown", "some-client/1.0");
    expect(response!.status).toBe(404);
    expect(await response!.text()).toContain("does-not-exist");
  });

  it("serves the home page (current issue contents) instead of a false 404", async () => {
    const response = await negotiatePageResponse("", "text/markdown", "some-client/1.0");
    expect(response!.status).toBe(200);
    expect(await response!.text()).toContain("Notes from the Margin");
  });
});
