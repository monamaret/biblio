import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter } from "react-router";
import Archive from "./archive";

// /issues lists every published issue, newest first, with number, title,
// season or date, article count and a link (contracts/site-surface.md
// "Routes").
describe("Archive", () => {
  it("lists every published issue newest first", () => {
    const { container } = render(
      <MemoryRouter initialEntries={["/issues"]}>
        <Archive />
      </MemoryRouter>,
    );
    expect(screen.getByRole("heading", { level: 1, name: "Issues" })).toBeDefined();
    const items = container.querySelectorAll(".archive-item");
    expect(items).toHaveLength(2);

    const newest = screen.getByRole("link", { name: /No\. 2 · Autumn 2026/ });
    expect(newest.getAttribute("href")).toBe("/issues/02-autumn-2026");
    expect(items[0].textContent).toContain("Notes from the Margin");
    expect(items[0].textContent).toContain("3 articles");
    expect(items[1].textContent).toContain("No. 1 · Spring 2026");
    expect(items[1].textContent).toContain("2 articles");
  });
});
