import { renderToStaticMarkup } from "react-dom/server";

export type ChartDatum = { label: string; value: number };

// ChartView is the small, sandboxed MCP App UI bundle for chart-like tool
// results (research.md "MCP App UI bundle approach"). It applies the
// site's tokens.css theme via CSS custom properties with explicit
// fallback values — `var(--token, fallback)` degrades gracefully to a
// usable, unthemed view if the sandbox can't load tokens.css at all,
// without needing any JS-level error handling (spec Edge Cases).
function ChartView({ data }: { data: ChartDatum[] }) {
  const max = Math.max(1, ...data.map((d) => d.value));
  return (
    <ul
      style={{
        listStyle: "none",
        margin: 0,
        padding: 0,
        fontFamily: "var(--font-body, system-ui, sans-serif)",
        color: "var(--color-fg, #1a1a1a)",
        background: "var(--color-bg, #ffffff)",
      }}
    >
      {data.map((d) => (
        <li key={d.label} style={{ display: "flex", alignItems: "center", gap: "0.5rem", margin: "0.25rem 0" }}>
          <span style={{ minWidth: "8rem" }}>{d.label}</span>
          <span
            style={{
              display: "block",
              height: "0.75rem",
              width: `${(d.value / max) * 100}%`,
              background: "var(--color-accent, #2563eb)",
              borderRadius: "0.25rem",
            }}
          />
        </li>
      ))}
    </ul>
  );
}

// renderChartView produces the full, sandboxable HTML document embedded as
// the tool result's `resource` content block (always paired with, never
// replacing, the primary `content` Markdown response — Principle IV,
// FR-004). Loads tokens.css if reachable; the CSS-variable fallbacks above
// keep the view usable even if the sandboxed context can't fetch it.
export function renderChartView(data: ChartDatum[]): string {
  const body = renderToStaticMarkup(<ChartView data={data} />);
  return `<!doctype html>
<html lang="en">
  <head>
    <meta charset="UTF-8" />
    <link rel="stylesheet" href="/tokens.css" />
  </head>
  <body>${body}</body>
</html>`;
}
