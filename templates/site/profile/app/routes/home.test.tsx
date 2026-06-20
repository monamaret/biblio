import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Home from "./home";

// Home renders the list/navigation view from the seeded content-index.json
// fixture — every authored entry, linked to its detail route (FR-001, FR-003;
// spec Scenario 1 & 3).
describe("Home", () => {
  it("renders a link for every content-index.json entry", () => {
    render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(
      screen.getByRole("link", { name: "About Me" }).getAttribute("href"),
    ).toBe("/about");
    expect(
      screen
        .getByRole("link", { name: "Sample Project" })
        .getAttribute("href"),
    ).toBe("/projects/sample-project");
  });
});
