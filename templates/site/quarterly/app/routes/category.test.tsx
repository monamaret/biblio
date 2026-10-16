import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Category from "./category";

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route path="/category/:category" element={<Category />} />
      </Routes>
    </MemoryRouter>,
  );
}

// /category/:key lists only that category's articles across issues, newest
// issue first, each showing its issue (FR-011, FR-013).
describe("Category", () => {
  it("lists only that category's articles, each with its issue", () => {
    const { container } = renderAt("/category/essays");
    expect(screen.getByRole("heading", { level: 1, name: "Essays" })).toBeDefined();
    const items = container.querySelectorAll(".category-item");
    expect(items).toHaveLength(3);
    expect(screen.getByRole("link", { name: "On Margins" }).getAttribute("href")).toBe(
      "/issues/02-autumn-2026/on-margins",
    );
    expect(items[2].textContent).toContain("On Reading Slowly");
    expect(screen.getAllByRole("link", { name: "No. 1 · Spring 2026" })[0].getAttribute("href")).toBe(
      "/issues/01-spring-2026",
    );
    expect(container.textContent).not.toContain("Type in Brno");
  });

  it("renders a not-found view for an unknown category", () => {
    renderAt("/category/poetry");
    expect(screen.getByRole("heading", { name: "Not found" })).toBeDefined();
  });
});
