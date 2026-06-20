import { describe, expect, it } from "vitest";
import { negotiateVariant, negotiatePageResponse } from "./accept";

// Re-exercises the type-agnostic negotiation logic (copied unchanged from
// templates/site/profile in Phase 2) against the log site's own content
// shape — confirming it generalizes rather than assuming it (FR-005;
// spec Scenario 4).
describe("negotiatePageResponse — log content shape", () => {
  it("returns the Markdown body with Vary header for a known entry slug", async () => {
    const response = await negotiatePageResponse(
      "posts/2026-06-01-hello-world",
      "text/markdown",
      "some-client/1.0",
    );
    expect(response).toBeDefined();
    expect(response!.headers.get("Vary")).toBe("Accept, User-Agent");
    const text = await response!.text();
    expect(text).toContain("Hello, World");
  });

  it("selects the ansi tier for a terminal User-Agent against log content (actual WASM rendering is exercised in ansi.test.ts)", () => {
    expect(negotiateVariant("text/markdown", "curl/8.0")).toBe("ansi");
  });

  it("returns a clear, structured 404 for a slug with no matching entry", async () => {
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
    expect(text).toContain("Hello, World");
  });
});
