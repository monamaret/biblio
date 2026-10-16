import { describe, expect, it, vi } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import * as load from "../content/load";
import Detail from "./detail";

vi.mock("../content/load", async (importOriginal) => {
  const actual = await importOriginal<typeof import("../content/load")>();
  return {
    ...actual,
    getEntryBySlug: vi.fn(actual.getEntryBySlug),
    getColophon: vi.fn(actual.getColophon),
  };
});

function renderAt(path: string) {
  return render(
    <MemoryRouter initialEntries={[path]}>
      <Routes>
        <Route path="*" element={<Detail />} />
      </Routes>
    </MemoryRouter>,
  );
}

// The article page: category label, title, author, dek, lead plate with
// "Fig. 1 — caption", body; the colophon; a not-found view
// (contracts/site-surface.md "Article page").
describe("Detail", () => {
  it("renders an article's category, title, author, dek and body", () => {
    const { container } = renderAt("/issues/02-autumn-2026/type-in-brno");
    expect(container.querySelector(".category-label")?.textContent).toBe("Design");
    expect(screen.getByRole("heading", { level: 1, name: "Type in Brno" })).toBeDefined();
    expect(container.querySelector(".byline")?.textContent).toContain("M. Hora");
    expect(container.querySelector(".dek")?.textContent).toBeTruthy();
    expect(container.querySelector(".article-body")?.textContent).toContain("Brno");
  });

  it("shows the lead plate with alt text and a numbered caption when the article has one", () => {
    const entry: load.Article = {
      ...(load.getEntryBySlug("issues/02-autumn-2026/on-margins") as load.Article),
      plate: "/plates/fedcba9876543210.jpg",
    };
    vi.mocked(load.getEntryBySlug).mockReturnValueOnce(entry);
    renderAt("/issues/02-autumn-2026/on-margins");
    const img = screen.getByRole("img", { name: entry.image_alt! });
    expect(img.getAttribute("src")).toBe("/plates/fedcba9876543210.jpg");
    expect(screen.getByText(`Fig. 1 — ${entry.image_caption!}`)).toBeDefined();
  });

  it("renders no lead image when the article has no plate yet", () => {
    const { container } = renderAt("/issues/02-autumn-2026/on-margins");
    expect(container.querySelector("figure.lead-plate")).toBeNull();
  });

  it("renders the colophon at /colophon", () => {
    renderAt("/colophon");
    expect(screen.getByRole("heading", { level: 1, name: "Colophon" })).toBeDefined();
    expect(screen.getByText(/Archivo/)).toBeDefined();
  });

  it("links the colophon from the masthead only when it exists", () => {
    renderAt("/issues/02-autumn-2026/type-in-brno");
    expect(screen.getByRole("link", { name: "Colophon" }).getAttribute("href")).toBe("/colophon");

    vi.mocked(load.getColophon).mockReturnValue(undefined);
    const { container } = renderAt("/issues/02-autumn-2026/type-in-brno");
    expect(Array.from(container.querySelectorAll("a")).some((a) => a.textContent === "Colophon")).toBe(false);
    vi.mocked(load.getColophon).mockRestore();
  });

  it("renders a not-found view for an unknown slug", () => {
    renderAt("/does-not-exist");
    expect(screen.getByRole("heading", { name: "Not found" })).toBeDefined();
  });
});
