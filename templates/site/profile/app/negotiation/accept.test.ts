import { describe, expect, it } from "vitest";
import { negotiateVariant, negotiatePageResponse } from "./accept";

describe("negotiateVariant", () => {
  it("selects html when Accept is absent", () => {
    expect(negotiateVariant(undefined, undefined)).toBe("html");
  });

  it("selects html for a typical browser Accept header", () => {
    expect(
      negotiateVariant("text/html,application/xhtml+xml,application/xml;q=0.9,*/*;q=0.8", "Mozilla/5.0"),
    ).toBe("html");
  });

  it("selects markdown when Accept explicitly prefers text/markdown", () => {
    expect(negotiateVariant("text/markdown", "some-http-client/1.0")).toBe("markdown");
  });

  it("selects ansi when Accept prefers text/markdown and User-Agent is a known terminal client", () => {
    expect(negotiateVariant("text/markdown", "curl/8.0")).toBe("ansi");
    expect(negotiateVariant("text/markdown", "Wget/1.21")).toBe("ansi");
  });

  it("selects markdown (not ansi) for curl-like User-Agent without an explicit markdown preference using just */*", () => {
    // */* alone, with no text/markdown mention at all, must resolve to html
    // regardless of User-Agent — Accept governs variant, User-Agent only
    // breaks the tier-i/tier-ii tie within the markdown family.
    expect(negotiateVariant("*/*", "curl/8.0")).toBe("html");
  });
});

// The contract requires deterministic resolution of an ambiguous/
// conflicting Accept header (spec Edge Cases) — this is the documented
// rule from research.md "Deterministic Accept/User-Agent tie-breaking
// rule": markdown wins on an exact tie.
describe("negotiateVariant ambiguous Accept resolution", () => {
  it("resolves an exact q-value tie between text/markdown and text/html in favor of markdown", () => {
    expect(negotiateVariant("text/markdown;q=0.8, text/html;q=0.8", "some-client/1.0")).toBe("markdown");
  });

  it("resolves an implicit tie (both default to q=1) in favor of markdown", () => {
    expect(negotiateVariant("text/markdown, text/html", "some-client/1.0")).toBe("markdown");
  });

  it("respects an explicit higher q for text/html over text/markdown", () => {
    expect(negotiateVariant("text/markdown;q=0.5, text/html;q=0.9", "some-client/1.0")).toBe("html");
  });
});

describe("negotiatePageResponse", () => {
  it("returns undefined for html (caller falls through to ASSETS)", async () => {
    const response = await negotiatePageResponse("about", "text/html", "Mozilla/5.0");
    expect(response).toBeUndefined();
  });

  it("returns the Markdown body with Vary header for a known slug", async () => {
    const response = await negotiatePageResponse("about", "text/markdown", "some-client/1.0");
    expect(response).toBeDefined();
    expect(response!.headers.get("Vary")).toBe("Accept, User-Agent");
    const text = await response!.text();
    expect(text).toContain("About Me");
  });

  it("returns a clear, structured 404 for a slug with no matching content-index entry, not a silent empty/200", async () => {
    const response = await negotiatePageResponse("does-not-exist", "text/markdown", "some-client/1.0");
    expect(response).toBeDefined();
    expect(response!.status).toBe(404);
    expect(response!.headers.get("Vary")).toBe("Accept, User-Agent");
    const text = await response!.text();
    expect(text).toContain("does-not-exist");
  });

  it("returns a generated index listing for the home page (empty slug), not a false 404", async () => {
    const response = await negotiatePageResponse("", "text/markdown", "some-client/1.0");
    expect(response).toBeDefined();
    expect(response!.status).toBe(200);
    const text = await response!.text();
    expect(text).toContain("About Me");
  });
});
