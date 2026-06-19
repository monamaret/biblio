import { readFileSync } from "node:fs";
import path from "node:path";
import { describe, expect, it } from "vitest";
import { render, cleanup } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Home from "./home";
import Detail from "./detail";

// Every route renders through root.tsx's single Layout, which links the
// unchanged v0.01-generated tokens.css — no route module may introduce its
// own competing inline/duplicate styling path (FR-003; spec Scenario 3).
describe("Theme application", () => {
  it("links tokens.css exactly once, from the shared root Layout", () => {
    const rootSource = readFileSync(
      path.join(process.cwd(), "app", "root.tsx"),
      "utf-8",
    );
    const matches = rootSource.match(/href="\/tokens\.css"/g) ?? [];
    expect(matches).toHaveLength(1);
  });

  it("renders Home with no inline/duplicate styling path", () => {
    const { container } = render(
      <MemoryRouter initialEntries={["/"]}>
        <Routes>
          <Route path="/" element={<Home />} />
        </Routes>
      </MemoryRouter>,
    );
    expect(container.querySelector("[style]")).toBeNull();
    expect(container.querySelector("style")).toBeNull();
    cleanup();
  });

  it("renders Detail with no inline/duplicate styling path", () => {
    const { container } = render(
      <MemoryRouter initialEntries={["/about"]}>
        <Routes>
          <Route path="*" element={<Detail />} />
        </Routes>
      </MemoryRouter>,
    );
    expect(container.querySelector("[style]")).toBeNull();
    expect(container.querySelector("style")).toBeNull();
    cleanup();
  });
});
