import { describe, expect, it } from "vitest";
import { render, screen } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Detail from "./detail";

// Detail renders the per-page view, including its section/order context
// (FR-002; spec Scenario 1.1) plus the rendered Markdown body.
describe("Detail", () => {
  it("renders the matching page's title, section, and body", () => {
    render(
      <MemoryRouter initialEntries={["/getting-started"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );

    expect(screen.getByRole("heading", { name: "Getting Started" })).toBeDefined();
    expect(screen.getByText(/Overview/)).toBeDefined();
    expect(screen.getByText(/Walk readers through the first steps/)).toBeDefined();
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
