import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Home from "./home";

// Home renders the chronological entry-list view from the seeded
// content-index.json fixture, displaying each entry's date and tags
// (FR-002; spec Scenario 1).
describe("Home", () => {
  it("renders a link with the entry's date and tags for every content-index.json entry", () => {
    render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </MemoryRouter>,
    );

    const link = screen.getByRole("link", { name: /Hello, World/ });
    expect(link.getAttribute("href")).toBe("/posts/2026-06-01-hello-world");

    const item = link.closest("li");
    expect(item?.textContent).toContain("2026-06-01");
    expect(item?.textContent).toContain("meta");
  });
});
