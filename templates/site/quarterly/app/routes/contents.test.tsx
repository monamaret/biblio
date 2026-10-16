import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import * as load from "../content/load";
import Contents from "./contents";

// Wrap the accessors so individual tests can substitute index states the
// sample content doesn't have (no issues, an empty issue, a cover plate).
vi.mock("../content/load", async (importOriginal) => {
  const actual = await importOriginal<typeof import("../content/load")>();
  return {
    ...actual,
    currentIssue: vi.fn(actual.currentIssue),
    getIssue: vi.fn(actual.getIssue),
    articlesInIssue: vi.fn(actual.articlesInIssue),
  };
});

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route path="/" element={<Contents />} />
        <Route path="/issues/:issue" element={<Contents />} />
      </Routes>
    </MemoryRouter>,
  );
}

// The contents page: masthead, "No. N · season", cover plate, articles
// grouped by category with numerals, leaders and authors (specs/004
// FR-013, FR-014; contracts/site-surface.md "Contents page").
describe("Contents", () => {
  it("renders the current issue's contents at /", () => {
    const { container } = renderAt("/");
    expect(screen.getByRole("link", { name: "The Quarterly" }).getAttribute("href")).toBe("/");
    expect(screen.getByText("No. 2 · Autumn 2026")).toBeDefined();
    expect(screen.getByRole("heading", { level: 1, name: "Notes from the Margin" })).toBeDefined();
    expect(screen.getByRole("heading", { level: 2, name: "Contents" })).toBeDefined();
    expect(screen.getAllByRole("heading", { level: 3 }).map((h) => h.textContent)).toEqual(["Essays", "Design"]);

    const items = container.querySelectorAll(".contents-item");
    expect(items).toHaveLength(3);
    expect(Array.from(items).map((li) => li.querySelector(".numeral")?.textContent)).toEqual(["I.", "II.", "III."]);
    expect(container.querySelectorAll(".contents-item .leader")).toHaveLength(3);

    const first = screen.getByRole("link", { name: "On Margins" });
    expect(first.getAttribute("href")).toBe("/issues/02-autumn-2026/on-margins");
    expect(items[0].querySelector(".author")?.textContent).toBe("A. Novak");

    expect(screen.getByRole("link", { name: /Previous issues/ }).getAttribute("href")).toBe("/issues");
  });

  it("renders a specific issue at /issues/:issue", () => {
    renderAt("/issues/01-spring-2026");
    expect(screen.getByText("No. 1 · Spring 2026")).toBeDefined();
    expect(screen.getByRole("link", { name: "On Reading Slowly" })).toBeDefined();
  });

  it("falls back to the date when an issue has no season", () => {
    const issue = { ...load.getIssue("02-autumn-2026")!, season: undefined };
    vi.mocked(load.getIssue).mockReturnValueOnce(issue);
    renderAt("/issues/02-autumn-2026");
    expect(screen.getByText("No. 2 · 2026-10-01")).toBeDefined();
  });

  it("shows the cover plate with its alt text when the issue has one", () => {
    const issue = { ...load.getIssue("01-spring-2026")!, cover_plate: "/plates/0123456789abcdef.jpg" };
    vi.mocked(load.getIssue).mockReturnValueOnce(issue);
    renderAt("/issues/01-spring-2026");
    const img = screen.getByRole("img", { name: issue.cover_alt! });
    expect(img.getAttribute("src")).toBe("/plates/0123456789abcdef.jpg");
  });

  it("shows no image when the issue has no plate", () => {
    const { container } = renderAt("/issues/02-autumn-2026");
    expect(container.querySelector("img")).toBeNull();
  });

  it("shows an empty state when no issue has been published", () => {
    vi.mocked(load.currentIssue).mockReturnValueOnce(undefined);
    const { container } = renderAt("/");
    expect(screen.getByText("No issues have been published yet.")).toBeDefined();
    expect(container.querySelector(".contents-item")).toBeNull();
  });

  it("shows the issue but a 'no articles' message for an issue with zero articles", () => {
    const issue = { ...load.getIssue("02-autumn-2026")!, article_count: 0 };
    vi.mocked(load.getIssue).mockReturnValueOnce(issue);
    vi.mocked(load.articlesInIssue).mockReturnValueOnce([]);
    const { container } = renderAt("/issues/02-autumn-2026");
    expect(screen.getByText("No. 2 · Autumn 2026")).toBeDefined();
    expect(screen.getByText("No articles in this issue yet.")).toBeDefined();
    expect(container.querySelector(".contents-item")).toBeNull();
    expect(screen.getByRole("link", { name: /Previous issues/ })).toBeDefined();
  });

  it("renders a not-found view for an unknown issue", () => {
    renderAt("/issues/99-nope");
    expect(screen.getByRole("heading", { name: "Not found" })).toBeDefined();
  });
});
