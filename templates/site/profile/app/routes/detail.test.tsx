import { afterEach, describe, expect, it } from "vitest";
import { cleanup, render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Detail from "./detail";

// Detail renders the per-entry view for both a single-segment slug ("about")
// and a multi-segment slug ("projects/sample-project") — the splat route
// must resolve both without one route per nesting depth (FR-001).
describe("Detail", () => {
  // vitest.config.ts doesn't set `test.globals: true`, so
  // @testing-library/react's auto-cleanup-after-each-test doesn't register —
  // without this, renders from earlier tests stay in the DOM and any two
  // tests asserting the same heading text collide.
  afterEach(cleanup);

  it("renders the matching entry's title and body for a single-segment slug", () => {
    render(
      <MemoryRouter initialEntries={["/about"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "About Me" })).toBeDefined();
    expect(screen.getByText(/Welcome to your new profile site/)).toBeDefined();
  });

  it("renders the matching entry's title and body for a multi-segment slug", () => {
    render(
      <MemoryRouter initialEntries={["/projects/sample-project"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(
      screen.getByRole("heading", { name: "Sample Project" }),
    ).toBeDefined();
    expect(screen.getByText(/Describe the project here/)).toBeDefined();
  });

  it("resolves a slug even with the trailing slash Cloudflare's static-asset serving redirects to", () => {
    render(
      <MemoryRouter initialEntries={["/about/"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "About Me" })).toBeDefined();
  });

  it("renders a not-found view for an unknown slug instead of throwing", () => {
    render(
      <MemoryRouter initialEntries={["/does-not-exist"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "Not found" })).toBeDefined();
  });
});
