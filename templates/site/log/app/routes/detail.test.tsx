import { afterEach, describe, expect, it } from "vitest";
import { cleanup, render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Detail from "./detail";

// Detail renders the per-entry view, displaying date/tags/project plus the
// rendered Markdown body (FR-002; spec Scenario 1).
describe("Detail", () => {
  // vitest.config.ts doesn't set `test.globals: true`, so
  // @testing-library/react's auto-cleanup-after-each-test doesn't register —
  // without this, renders from earlier tests stay in the DOM and any two
  // tests asserting the same heading text collide.
  afterEach(cleanup);

  it("renders the matching entry's title, date, tags, and body", () => {
    render(
      <MemoryRouter initialEntries={["/posts/2026-06-01-hello-world"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "Hello, World" })).toBeDefined();
    expect(screen.getByText(/2026-06-01/)).toBeDefined();
    expect(screen.getByText(/meta/)).toBeDefined();
    expect(screen.getByText(/first entry in your project log/)).toBeDefined();
  });

  it("resolves a slug even with the trailing slash Cloudflare's static-asset serving redirects to", () => {
    render(
      <MemoryRouter initialEntries={["/posts/2026-06-01-hello-world/"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "Hello, World" })).toBeDefined();
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
