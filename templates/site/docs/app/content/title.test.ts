import { beforeEach, describe, expect, it, vi } from "vitest";
import { siteMeta } from "./load";
import { pageTitle } from "./title";

vi.mock("./load", () => ({ siteMeta: vi.fn() }));

describe("pageTitle", () => {
  beforeEach(() => {
    vi.mocked(siteMeta).mockReturnValue({ name: "502.io", description: "" });
  });

  it("is the site name alone for the home page", () => {
    expect(pageTitle()).toBe("502.io");
  });

  it("puts the page before the site name", () => {
    expect(pageTitle("On Margins")).toBe("On Margins · 502.io");
  });

  it("is just the page when site.yaml has no name", () => {
    vi.mocked(siteMeta).mockReturnValue({ name: "", description: "" });
    expect(pageTitle("On Margins")).toBe("On Margins");
  });
});
