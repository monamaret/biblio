import { describe, expect, it } from "vitest";
import { render, screen, fireEvent } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Home from "./home";

function renderHome() {
  render(
    <MemoryRouter initialEntries={["/"]}>
      <Routes>
        <Route path="/" element={<Home />} />
      </Routes>
    </MemoryRouter>,
  );
}

// Home renders the chronological entry-list view from the seeded
// content-index.json fixture, displaying each entry's date and tags
// (FR-002; spec Scenario 1).
describe("Home", () => {
  it("renders a link with the entry's date and tags for every content-index.json entry", () => {
    renderHome();

    const link = screen.getByRole("link", { name: /Hello, World/ });
    expect(link.getAttribute("href")).toBe("/posts/2026-06-01-hello-world");

    const item = link.closest("li");
    expect(item?.textContent).toContain("2026-06-01");
    expect(item?.textContent).toContain("meta");
  });
});

// Filter controls reflect the visitor's selected date-range/tag/project
// state and narrow the rendered list to only the matching entries
// (FR-003; spec Scenario 2).
describe("Home — filter controls", () => {
  it("shows every seeded entry before any filter is applied", () => {
    renderHome();
    expect(screen.getByRole("link", { name: /Hello, World/ })).toBeDefined();
    expect(screen.getByRole("link", { name: /Shipped the filter UI/ })).toBeDefined();
    expect(screen.getByRole("link", { name: /Refactored the content index/ })).toBeDefined();
  });

  it("narrows the list to only entries matching the selected tag", () => {
    renderHome();

    fireEvent.change(screen.getByRole("combobox", { name: /tag/i }), {
      target: { value: "frontend" },
    });

    expect(screen.getByRole("link", { name: /Shipped the filter UI/ })).toBeDefined();
    expect(screen.queryByRole("link", { name: /Hello, World/ })).toBeNull();
    expect(screen.queryByRole("link", { name: /Refactored the content index/ })).toBeNull();
  });

  it("narrows the list to only entries matching the selected project", () => {
    renderHome();

    fireEvent.change(screen.getByRole("combobox", { name: /project/i }), {
      target: { value: "biblio" },
    });

    expect(screen.getByRole("link", { name: /Shipped the filter UI/ })).toBeDefined();
    expect(screen.getByRole("link", { name: /Refactored the content index/ })).toBeDefined();
    expect(screen.queryByRole("link", { name: /Hello, World/ })).toBeNull();
  });
});
