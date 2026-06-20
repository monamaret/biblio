import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Home from "./home";

// Home renders the TOC/navigation view from the seeded content-index.json
// fixture, grouped by section and ordered within each (FR-002; spec
// Scenario 1.1) — no tag/project filter controls anywhere (Scenario 1.2).
describe("Home", () => {
  it("renders every page as a link, grouped under its section heading", () => {
    render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "Overview" })).toBeDefined();
    expect(screen.getByRole("heading", { name: "Reference" })).toBeDefined();
    expect(screen.getByRole("link", { name: "Introduction" }).getAttribute("href")).toBe("/index");
    expect(screen.getByRole("link", { name: "Getting Started" }).getAttribute("href")).toBe("/getting-started");
    expect(screen.getByRole("link", { name: "Configuration" }).getAttribute("href")).toBe("/configuration");
  });

  it("renders pages within a section in order", () => {
    render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </MemoryRouter>,
    );

    const links = screen.getAllByRole("link").map((el) => el.textContent);
    expect(links.indexOf("Introduction")).toBeLessThan(links.indexOf("Getting Started"));
  });

  it("renders no tag/project filter controls anywhere (spec Scenario 1.2)", () => {
    render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.queryByRole("combobox")).toBeNull();
    expect(screen.queryByLabelText(/tag/i)).toBeNull();
    expect(screen.queryByLabelText(/project/i)).toBeNull();
  });
});
