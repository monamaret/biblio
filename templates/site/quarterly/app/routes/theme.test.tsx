import { existsSync, readFileSync } from "node:fs";
import type { ReactElement } from "react";
import path from "node:path";
import { describe, expect, it } from "vitest";
import { render, cleanup } from "@testing-library/react";
import { MemoryRouter, Routes, Route } from "react-router";
import Contents from "./contents";
import Archive from "./archive";
import Category from "./category";
import Detail from "./detail";

const root = process.cwd();
const css = readFileSync(path.join(root, "public", "quarterly.css"), "utf-8");
const rootSource = readFileSync(path.join(root, "app", "root.tsx"), "utf-8");

// cssOutsideFontFace drops comments and @font-face blocks, leaving the
// rules that must consume tokens only.
function cssOutsideFontFace(source: string): string {
  return source.replace(/\/\*[\s\S]*?\*\//g, "").replace(/@font-face\s*{[^}]*}/g, "");
}

const NAMED_COLORS = /(?<![\w-])(black|white|red|blue|green|gray|grey|yellow|orange|purple|pink|brown|navy|teal|silver|maroon)(?![\w-])/i;

// One styling path: tokens.css (generated from site.yaml) then
// quarterly.css (template-owned, tokens only), no inline styles, fonts
// self-hosted (specs/004-quarterly-site-type FR-026, FR-027;
// contracts/config-and-tokens.md "Template stylesheet").
describe("quarterly theme", () => {
  it("links tokens.css once, before quarterly.css, from the root layout", () => {
    expect(rootSource.match(/href="\/tokens\.css"/g)).toHaveLength(1);
    expect(rootSource.indexOf('href="/tokens.css"')).toBeLessThan(rootSource.indexOf('href="/quarterly.css"'));
  });

  it("uses no color literal or font family name outside @font-face", () => {
    const rules = cssOutsideFontFace(css);
    expect(rules).not.toMatch(/#[0-9a-f]{3,8}\b/i);
    expect(rules).not.toMatch(/\b(rgb|rgba|hsl|hsla)\(/i);
    expect(rules).not.toMatch(NAMED_COLORS);
    expect(rules).not.toMatch(/font-family\s*:(?!\s*var\()/i);
    expect(rules).not.toMatch(/Archivo|Plex/i);
  });

  it("loads fonts only from the site itself, and every @font-face file exists", () => {
    for (const [, url] of css.matchAll(/url\(\s*["']?([^"')]+)["']?\s*\)/g)) {
      expect(url.startsWith("/") && !url.startsWith("//"), url).toBe(true);
    }
    for (const [, href] of rootSource.matchAll(/href="([^"]+)"/g)) {
      expect(/^https?:|^\/\//.test(href), href).toBe(false);
    }
    expect(css).not.toMatch(/@import/);
    const fontFaces = css.match(/@font-face\s*{[^}]*}/g) ?? [];
    expect(fontFaces.length).toBeGreaterThan(0);
    for (const block of fontFaces) {
      for (const [, url] of block.matchAll(/url\(\s*["']?([^"')]+)["']?\s*\)/g)) {
        expect(url.startsWith("/fonts/"), url).toBe(true);
        expect(existsSync(path.join(root, "public", url)), url).toBe(true);
      }
    }
  });

  it("renders every route without inline styles", () => {
    const cases: Array<[string, string, ReactElement]> = [
      ["/", "/", <Contents />],
      ["/issues", "/issues", <Archive />],
      ["/category/essays", "/category/:category", <Category />],
      ["/issues/02-autumn-2026/on-margins", "*", <Detail />],
    ];
    for (const [url, pattern, element] of cases) {
      const { container } = render(
        <MemoryRouter initialEntries={[url]}>
          <Routes>
            <Route path={pattern} element={element} />
          </Routes>
        </MemoryRouter>,
      );
      expect(container.querySelector("[style]"), url).toBeNull();
      expect(container.querySelector("style"), url).toBeNull();
      cleanup();
    }
  });
});
