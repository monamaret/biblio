import { describe, expect, it } from "vitest";
import { negotiateVariant, negotiatePageResponse } from "./accept";

// Re-exercises the type-agnostic negotiation logic (copied unchanged from
// templates/site/profile in Phase 2) against the docs site's own content
// shape — confirming it generalizes rather than assuming it (FR-004;
// spec Scenario 3).
describe("negotiatePageResponse — docs content shape", () => {
  it("returns the Markdown body with Vary header for a known page slug", async () => {
    const response = await negotiatePageResponse("getting-started", "text/markdown", "some-client/1.0");
    expect(response).toBeDefined();
    expect(response!.headers.get("Vary")).toBe("Accept, User-Agent");
    const text = await response!.text();
    expect(text).toContain("Getting Started");
  });

  it("selects the ansi tier for a terminal User-Agent against docs content (actual WASM rendering is exercised in ansi.test.ts)", () => {
    expect(negotiateVariant("text/markdown", "curl/8.0")).toBe("ansi");
  });

  it("returns a clear, structured 404 for a slug with no matching page", async () => {
    const response = await negotiatePageResponse("does-not-exist", "text/markdown", "some-client/1.0");
    expect(response).toBeDefined();
    expect(response!.status).toBe(404);
    const text = await response!.text();
    expect(text).toContain("does-not-exist");
  });

  it("returns a generated index listing for the home page (empty slug), not a false 404", async () => {
    const response = await negotiatePageResponse("", "text/markdown", "some-client/1.0");
    expect(response).toBeDefined();
    expect(response!.status).toBe(200);
    const text = await response!.text();
    expect(text).toContain("Introduction");
  });
});
