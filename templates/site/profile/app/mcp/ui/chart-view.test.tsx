import { describe, expect, it } from "vitest";
import { renderChartView } from "./chart-view";

describe("ChartView MCP App UI bundle", () => {
  it("renders correctly and links tokens.css", () => {
    const html = renderChartView([
      { label: "About Me", value: 1 },
      { label: "Sample Project", value: 1 },
    ]);

    expect(html).toContain('<link rel="stylesheet" href="/tokens.css"');
    expect(html).toContain("About Me");
    expect(html).toContain("Sample Project");
  });

  // The sandbox's stylesheet `<link>` may fail to load at runtime (spec
  // Edge Cases) — the view must still be usable. Every themed value uses a
  // CSS custom property with an explicit fallback (`var(--token,
  // fallback)`), which is the browser-native mechanism for exactly this:
  // no JS error handling is needed because a missing stylesheet never
  // throws, it just means the custom properties are undefined and the
  // fallback values apply automatically.
  it("falls back to a usable, unthemed view when tokens.css can't load", () => {
    const html = renderChartView([{ label: "Solo", value: 1 }]);

    expect(html).toMatch(/var\(--color-accent,\s*#2563eb\)/);
    expect(html).toMatch(/var\(--color-fg,\s*#1a1a1a\)/);
    expect(html).toMatch(/var\(--font-body,\s*system-ui/);
  });

  it("renders an empty dataset without throwing", () => {
    expect(() => renderChartView([])).not.toThrow();
  });
});
